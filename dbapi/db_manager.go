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
	if "" == strings.TrimSpace(name) {
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

func (dmb DBManager) RemoveDB(name string) error {
	dmb.Lock()
	defer dmb.Unlock()

	if _, ok := dmb.dbs[name]; !ok {
		return fmt.Errorf("DBManager.RemoveDB: no such db '%s'", name)
	}

	delete(dmb.dbs, name)

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

	for name, db := range dbm.dbs {
		lexs, err := ListLexicons(db)
		if err != nil {
			return res, fmt.Errorf("DBManager.ListLexicons : %v", err)
		}

		for _, l := range lexs {
			res = append(res, name+":"+l.Name)
		}
	}

	return res, nil
}
