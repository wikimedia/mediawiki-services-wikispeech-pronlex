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
		return fmt.Errorf("UpdateSchema: %v", err)
	}
	defer db.Close()

	var userVersion int
	err = db.QueryRow("PRAGMA user_version").Scan(&userVersion)
	if err != nil {

		return fmt.Errorf("UpdateSchema: %v", err)
	}
	fmt.Fprintf(os.Stderr, "UpdateSchema: current user_version: %d\n", userVersion)

	if userVersion < 1 {
		// Substitute faulty version of trigger

		rez, err := db.Exec("DROP TRIGGER insertPref")
		if err != nil {
			return fmt.Errorf("UpdateSchema failed when dropping trigger : %v", err)
		}
		n, err := rez.RowsAffected()
		if err != nil {
			return fmt.Errorf("UpdateSchema failed when calling RowsAffected : %v", err)
		}
	}

	return nil
}
