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
	if nil == db {
		return fmt.Errorf("DBManager.AddDB: illegal argument: db must not be nil")
	}

	dmb.Lock()
	defer dmb.Unlock()
	if _, ok := dmb.dbs[name]; ok {
		return fmt.Errorf("DBManager.AddDB: db already exists (use RemoveDB to remove it): '%s'", name)
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
