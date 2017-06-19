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

func (DBManager) New() DBManager {
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

func (dbm DBManager) ListDBNames() []string {
	var res []string

	dbm.RLock()
	defer dbm.RUnlock()

	for k, _ := range dbm.dbs {
		res = append(res, k)
	}

	return res
}

func (dbm DBManager) ListLexicons() ([]string, error) {
	var res []string

	dbm.RLock()
	defer dbm.RUnlock()

	for dbName, db := range dbm.dbs {
		lexs, err := ListLexicons(db)
		if err != nil {
			return res, fmt.Errorf("DBManager.ListLexicons : %v", err)
		}

		// l is a dbapi.Lexicon struct
		for _, l := range lexs {
			// A full lexicon name consists of the name of
			// the db and the name of the lexicon in the
			// db joined by ':'
			fullName := dbName + ":" + strings.TrimSpace(l.Name)
			res = append(res, fullName)
		}
	}

	return res, nil
}
