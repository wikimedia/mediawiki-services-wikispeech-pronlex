package dbapi

// For validating a lexicon db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
)

func processChunk(db *sql.DB, chunk []int64, vd validation.Validator, stats ValStats) (ValStats, error) {
	q := Query{EntryIDs: chunk}
	var w lex.EntrySliceWriter

	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		tx.Rollback()
		return stats, fmt.Errorf("failed to initialize transaction : %v", err)
	}

	err = lookUp(db, []lex.LexName{}, q, &w)
	if err != nil {
		tx.Rollback()
		return stats, fmt.Errorf("couldn't lookup from ids : %s", err)
	}
	if w.Size() != len(chunk) {
		tx.Rollback()
		return stats, fmt.Errorf("got %d input ids, but found %d entries", len(chunk), w.Size())
	}

	validated, _ := vd.ValidateEntries(w.Entries)
	origMap := make(map[int64]lex.Entry)
	for _, orig := range w.Entries {
		origMap[orig.ID] = orig
	}
	updated := []lex.Entry{}
	for _, e := range validated {
		stats.ValidatedEntries++
		newVal := e.EntryValidations
		oldVal := origMap[e.ID].EntryValidations
		if len(newVal) == 0 && len(oldVal) == 0 {
			// no change
		} else {
			updated = append(updated, e)
		}
		if len(newVal) > 0 {
			stats.InvalidEntries++
			for _, v := range newVal {
				stats.TotalValidations++
				stats.Levels[strings.ToLower(v.Level)]++
				stats.Rules[strings.ToLower(v.RuleName+" ("+v.Level+")")]++
			}
		}
	}

	err = updateValidationTx(tx, updated)
	if err != nil {
		tx.Rollback()
		return stats, fmt.Errorf("couldn't update validation : %s", err)
	}
	tx.Commit()
	return stats, nil
}

func Validate(db *sql.DB, lexNames []lex.LexName, logger Logger, vd validation.Validator, q Query) (ValStats, error) {

	start := time.Now()

	stats := ValStats{Levels: make(map[string]int), Rules: make(map[string]int)}

	logger.Write(fmt.Sprintf("query: %v", q))

	q.PageLength = 0 //todo?
	q.Page = 0       //todo?

	logger.Write("Fetching entries from lexicon ... ")
	ids, err := lookUpIds(db, lexNames, q)
	if err != nil {
		return stats, fmt.Errorf("couldn't lookup for validation : %s", err)
	}
	total := len(ids)
	stats.TotalEntries = total
	stats.ValidatedEntries = 0
	stats.InvalidEntries = 0
	stats.TotalValidations = 0
	logger.Write(fmt.Sprintf("Found %d entries", total))

	n := 0

	chunkSize := 500
	var chunk []int64

	for _, id := range ids {
		n = n + 1
		chunk = append(chunk, id)

		if n%chunkSize == 0 {
			stats, err = processChunk(db, chunk, vd, stats)
			if err != nil {
				return stats, err
			}
			chunk = []int64{}
		}

		if n%10 == 0 {
			js, err := json.Marshal(stats)
			if err != nil {
				return stats, fmt.Errorf("couldn't marshal validation stats : %s", err)
			}
			msg := fmt.Sprintf("%s", js)
			logger.Write(msg)
		}
	}
	if len(chunk) > 0 {
		stats, err = processChunk(db, chunk, vd, stats)
		if err != nil {
			return stats, err
		}
		chunk = []int64{}
	}
	end := time.Now()
	log.Printf("dbapi/validation.go Validate took %v\n", end.Sub(start))

	return stats, nil
}
