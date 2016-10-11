package dbapi

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
)

type ValStats struct {
	Values map[string]int `json:"values"`
}

func (v ValStats) increment(key string, incr int) {
	if _, ok := v.Values[key]; !ok {
		v.Values[key] = 0
	}
	v.Values[key] = v.Values[key] + incr

}

func processChunk(db *sql.DB, chunk []int64, vd validation.Validator, stats ValStats) error {
	q := Query{EntryIDs: chunk}
	var w lex.EntrySliceWriter

	err := LookUp(db, q, &w)
	if err != nil {
		return fmt.Errorf("couldn't lookup from ids : %s", err)
	}
	updated := []lex.Entry{}
	for _, e := range w.Entries {
		oldVal := e.EntryValidations
		vd.ValidateEntry(&e)
		newVal := e.EntryValidations
		stats.increment("Validated", 1)
		if len(newVal) > 0 {
			stats.increment("Invalid", 1)
			for _, v := range newVal {
				stats.increment("Level:"+v.Level, 1)
				stats.increment("Rule:"+v.RuleName, 1)
			}
		}
		if len(oldVal) > 0 && len(newVal) > 0 {
			updated = append(updated, e)
		}
	}
	err = UpdateValidation(db, updated)
	if err != nil {
		return fmt.Errorf("couldn't update validation : %s", err)
	}
	return nil
}

func Validate(db *sql.DB, logger Logger, vd validation.Validator, q Query) (ValStats, error) {

	stats := ValStats{Values: make(map[string]int)}

	logger.Write(fmt.Sprintf("query: %v", q))

	q.PageLength = 0 //todo?
	q.Page = 0       //todo?

	logger.Write("Fetching entries from lexicon ... ")
	ids, err := LookUpIds(db, q)
	if err != nil {
		return stats, fmt.Errorf("couldn't lookup for validation : %s", err)
	}
	total := len(ids)
	stats.Values["Total"] = total
	logger.Write(fmt.Sprintf("Found %d entries", total))

	n := 0

	chunkSize := 500
	var chunk []int64

	for _, id := range ids {
		n = n + 1
		chunk = append(chunk, id)

		if n%chunkSize == 0 {
			err = processChunk(db, chunk, vd, stats)
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
		processChunk(db, chunk, vd, stats)
		if err != nil {
			return stats, err
		}
		chunk = []int64{}
	}

	return stats, nil
}
