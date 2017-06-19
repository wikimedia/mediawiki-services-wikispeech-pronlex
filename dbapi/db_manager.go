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
		return fmt.Errorf("DBManager.AddDB: db '%s' has already been added (use RemoveDB to remove it)", name)
	}

	dmb.dbs[name] = db

	return nil
}
