package dbapi

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
)

type DBManager struct {
	sync.RWMutex
	dbs map[string]*sql.DB
}

func NewDBManager() DBManager {
	return DBManager{dbs: make(map[string]*sql.DB)}
}

func (dmb DBManager) AddDB(name string, db *sql.DB) error {
	name = strings.TrimSpace(name)
	if "" == name {
		return fmt.Errorf("DBManager.AddDB: illegal argument: name must not be empty")
	}
	if strings.Contains(name, ":") {
		return fmt.Errorf("DBManager.AddDB: illegal argument: name must not contain ':'")
	}
	if nil == db {
		return fmt.Errorf("DBManager.AddDB: illegal argument: db must not be nil")
	}

	dmb.Lock()
	defer dmb.Unlock()

	if _, ok := dmb.dbs[name]; ok {
		return fmt.Errorf("DBManager.AddDB: db already exists: '%s'", name)
	}

	dmb.dbs[name] = db

	return nil
}

func (dbm DBManager) RemoveDB(name string) error {
	name = strings.TrimSpace(name)
	dbm.Lock()
	defer dbm.Unlock()

	if _, ok := dbm.dbs[name]; !ok {
		return fmt.Errorf("DBManager.RemoveDB: no such db '%s'", name)
	}

	delete(dbm.dbs, name)

	return nil
}

func (dbm DBManager) DefineLexicon(dbName string, lexes ...string) error {
	for _, l := range lexes {
		db, ok := dbm.dbs[dbName]
		if !ok {
			return fmt.Errorf("DBManager.DefineLexicon: No such db: '%s'", dbName)
		}
		_, err := InsertLexicon(db, Lexicon{Name: l})
		if err != nil {
			return fmt.Errorf("DBManager.DefineLexicon: failed to add '%s:%s' : %v", dbName, l, err)
		}
	}

	return nil
}

func (dbm DBManager) ListDBNames() []string {
	var res []string

	dbm.RLock()
	defer dbm.RUnlock()

	for k, _ := range dbm.dbs {
		res = append(res, k)
	}

	return res
}

type lexRes struct {
	dbName string
	lexs   []Lexicon
	err    error
}

// Warning: this is maybe my first attempt at concurrency using a channel in Go

func (dbm DBManager) ListLexicons() ([]string, error) {
	var res []string

	dbm.RLock()
	defer dbm.RUnlock()

	ch := make(chan lexRes)
	defer close(ch) // ?
	dbs := dbm.dbs
	// Go ask each db instance in its own Go-routine
	for dbName, db := range dbs {
		go func(dbName string, db *sql.DB, ch chan lexRes) {
			lexs, err := ListLexicons(db)
			r := lexRes{dbName: dbName, lexs: lexs, err: err}
			ch <- r
		}(dbName, db, ch)
	}

	// Read result from channel
	for i := 0; i < len(dbs); i++ {
		var r lexRes
		r = <-ch // Blocks until there is a result (I
		// think). Can we be stuck here forever, if
		// db call hangs?

		// If we encounter an error, just bail out
		if r.err != nil {
			return res, fmt.Errorf("DBManager.ListLexicons freak out : %v", r.err)

		}

		// A full lexicon name consists of the name of
		// the db and the name of the lexicon in the
		// db joined by ':'
		for _, l := range r.lexs {
			fullName := r.dbName + ":" + strings.TrimSpace(l.Name)
			res = append(res, fullName)
		}
	}

	return res, nil
}

// func (dbm DBManager) ListLexicons0() ([]string, error) {
// 	var res []string

// 	dbm.RLock()
// 	defer dbm.RUnlock()

// 	for dbName, db := range dbm.dbs {
// 		lexs, err := ListLexicons(db)
// 		if err != nil {
// 			return res, fmt.Errorf("DBManager.ListLexicons : %v", err)
// 		}

// 		// l is a dbapi.Lexicon struct
// 		for _, l := range lexs {
// 			// A full lexicon name consists of the name of
// 			// the db and the name of the lexicon in the
// 			// db joined by ':'
// 			fullName := dbName + ":" + strings.TrimSpace(l.Name)
// 			res = append(res, fullName)
// 		}
// 	}

// 	return res, nil
// }
