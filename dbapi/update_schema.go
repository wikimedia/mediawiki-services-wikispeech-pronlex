package dbapi

import (
	"database/sql"
	"fmt"
	"os"
)

func updateInsertPrefTrigger(tx *sql.Tx) error {

	return nil
}

// UpdateSchema migrates a 'live' pronlex db to a new schema
// version. The dbFile argument is the path to an Sqlite db file.
func UpdateSchema(dbFile string) error {
	//var err error

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return fmt.Errorf("dbapi.UpdateSchema: %v\n", err)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("dbapi.UpdateSchema: %v", err)
	}
	defer db.Close()
	tx, err := db.Begin()
	defer tx.Commit()

	var userVersion int
	err = tx.QueryRow("PRAGMA user_version").Scan(&userVersion)
	if err != nil {

		return fmt.Errorf("UpdateSchema: %v", err)
	}
	fmt.Fprintf(os.Stderr, "dbapi.UpdateSchema: current user_version: %d\n", userVersion)

	if userVersion < 1 {
		// Substitute faulty version of trigger

		//Defined in dbapi.go
		triggs0, err := ListNamesOfTriggers(db)
		triggs := make(map[string]bool)
		for _, t := range triggs0 {
			fmt.Println(t)
			triggs[t] = true
		}
		if _, ok := triggs["insertPref"]; ok {
			rez, err := tx.Exec("DROP TRIGGER insertPref")
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("dbapi.UpdateSchema failed when dropping trigger : %v", err)
			}
			_, err = rez.RowsAffected()
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("dbapi.UpdateSchema failed when calling RowsAffected : %v", err)
			}
			fmt.Println("DROPPED TRIGGER insertPref")
		}
		// Misspelled name of trigger in some version of schema
		if _, ok := triggs["updatetPref"]; ok {
			rez, err := tx.Exec("DROP TRIGGER updatetPref")
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("dbapi.UpdateSchema failed when dropping trigger : %v", err)
			}
			_, err = rez.RowsAffected()
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("dbapi.UpdateSchema failed when calling RowsAffected : %v", err)
			}
			fmt.Println("DROPPED TRIGGER updatetPref")
		}

		if _, ok := triggs["updatePref"]; ok {
			rez, err := tx.Exec("DROP TRIGGER updatePref")
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("dbapi.UpdateSchema failed when dropping trigger : %v", err)
			}
			_, err = rez.RowsAffected()
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("dbapi.UpdateSchema failed when calling RowsAffected : %v", err)
			}
			fmt.Println("DROPPED TRIGGER updatePref")
		}

		fmt.Fprintf(os.Stderr, "%s: user_version %d\n", dbFile, userVersion)

		createTriggers := `CREATE TRIGGER insertPref BEFORE INSERT ON ENTRY
                                   BEGIN
                                     UPDATE entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
                                   END;
                                   CREATE TRIGGER updatePref BEFORE UPDATE ON ENTRY
                                   BEGIN
                                     UPDATE entry SET preferred = 0 WHERE strn = NEW.strn AND NEW.preferred <> 0 AND lexiconid = NEW.lexiconid;
                                   END;
                                   
                                   CREATE TRIGGER insertEntryStatus BEFORE INSERT ON ENTRYSTATUS
                                   BEGIN 
                                     UPDATE entrystatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
                                   END;
                                   CREATE TRIGGER updateEntryStatus BEFORE UPDATE ON ENTRYSTATUS
                                   BEGIN
                                     UPDATE entrystatus SET current = 0 WHERE entryid = NEW.entryid AND NEW.current <> 0;
                                   END;`

		_, err = tx.Exec(createTriggers)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create trigger : %v", err)
		}
		_, err = tx.Exec("PRAGMA user_version = 1")
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to set user_version variable : %v", err)
		}

		fmt.Fprintf(os.Stderr, "dbapi.UpdateSchema: Created triggers\n")
	}

	err = tx.QueryRow("PRAGMA user_version").Scan(&userVersion)
	if err != nil {
		return fmt.Errorf("UpdateSchema: %v", err)
	} else {
		fmt.Printf("%s user_version %d\n", dbFile, userVersion)
	}

	return nil
}
