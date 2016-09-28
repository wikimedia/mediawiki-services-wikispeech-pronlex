package dbapi

import (
	"database/sql"
	"fmt"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
)

type ValStats struct {
	Values map[string]int
}

func (v ValStats) increment(key string, incr int) {
	if _, ok := v.Values[key]; !ok {
		v.Values[key] = 0
	}
	v.Values[key] = v.Values[key] + incr

}

func Validate(db *sql.DB, logger Logger, vd validation.Validator, q Query) (ValStats, error) {

	stats := ValStats{Values: make(map[string]int)}

	logger.Write(fmt.Sprintf("query: %v", q))

	q.PageLength = 0 //todo?
	q.Page = 0       //todo?

	ew := lex.EntrySliceWriter{}
	logger.Write("Fetching entries from lexicon ... ")
	err := LookUp(db, q, &ew)
	if err != nil {
		return stats, fmt.Errorf("couldn't lookup for validation : %s", err)
	}
	stats.Values["Total"] = len(ew.Entries)
	logger.Write(fmt.Sprintf("Found %d entries", len(ew.Entries)))

	n := 0

	for _, e := range ew.Entries {
		n = n + 1
		oldVal := e.EntryValidations
		vd.ValidateEntry(&e)
		newVal := e.EntryValidations

		// if entry is updated with validations, update entry in db
		if len(oldVal) > 0 || len(newVal) > 0 {
			UpdateEntry(db, e)
		}
		stats.increment("Validated", 1)
		if len(newVal) > 0 {
			stats.increment("Invalid", 1)
		}
		for _, v := range newVal {
			stats.increment("Level:"+v.Level, 1)
			stats.increment("Rule:"+v.RuleName, 1)
		}
		if n%10 == 0 {
			msg := fmt.Sprintf("%v", stats)
			logger.Write(msg)
		}
	}

	return stats, nil
}
