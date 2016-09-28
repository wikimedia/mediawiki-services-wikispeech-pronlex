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

	stats := ValStats{}

	logger.Write(fmt.Sprintf("query: %v", q))

	q.PageLength = 0 //todo?
	q.Page = 0       //todo?

	ew := lex.EntrySliceWriter{}
	err := LookUp(db, q, &ew)
	if err != nil {
		return stats, fmt.Errorf("couldn't lookup for validation : %s", err)
	}
	stats.Values["Total"] = len(ew.Entries)

	for _, e := range ew.Entries {
		oldVal := e.EntryValidations
		vd.ValidateEntry(&e)
		newVal := e.EntryValidations

		// if entry is updated with validations, update entry in db
		if len(oldVal) > 0 || len(newVal) > 0 {
			UpdateEntry(db, e)
		}
		if len(newVal) > 0 {
			stats.increment("Invalid", 1)
		}
		for _, v := range newVal {
			stats.increment("Level:"+v.Level, 1)
			stats.increment("Rule:"+v.RuleName, 1)
		}
	}

	return stats, nil
}
