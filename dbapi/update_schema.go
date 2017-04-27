package dbapi

import (
	"database/sql"
	"fmt"
	"os"
)

// UpdateSchema migrates a 'live' pronlex db to a new schema
// version. The dbFile argument is the path to an Sqlite db file.
func UpdateSchema(dbFile string) error {
	//var err error

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return fmt.Errorf("UpdateSchema: %v\n", err)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("UpdateSchama: %v", err)
	}
	defer db.Close()

	var userVersion int
	err = db.QueryRow("PRAGMA user_verion").Scan(&userVersion)

	return nil
}
