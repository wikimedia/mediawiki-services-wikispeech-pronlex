package dbapi

import (
	"database/sql"
	"fmt"
	"os"
)

/*
func updateInsertPrefTrigger(tx *sql.Tx) error {

	return nil
}
*/

func dropTrigger(tx *sql.Tx, triggerName string) error {

	triggs0, err := listNamesOfTriggersTx(tx) // Defined in dbapi.go
	if err != nil {
		//fmt.Fprintf(os.Stderr, "What? : %v\n", err)
		return fmt.Errorf("dbapi.dropTrigger : %v", err)
	}
	triggs := make(map[string]bool)
	for _, t := range triggs0 {
		//fmt.Println(">>> " + t)
		triggs[t] = true
	}

	if _, ok := triggs[triggerName]; ok {
		rez, err := tx.Exec("DROP TRIGGER " + triggerName)
		if err != nil {
			msg := fmt.Sprintf("dbapi.UpdateSchema failed when dropping trigger : %v", err)
			err2 := tx.Rollback()
			if err2 != nil {
				msg = fmt.Sprintf("%s : dropTrigger failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)
		}
		_, err = rez.RowsAffected()
		if err != nil {
			msg := fmt.Sprintf("dbapi.UpdateSchema failed when calling RowsAffected : %v", err)
			err2 := tx.Rollback()
			if err2 != nil {
				msg = fmt.Sprintf("%s : dropTrigger failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)
		}
		//fmt.Println("DROPPED TRIGGER " + triggerName)
	} //else {
	//	fmt.Fprintf(os.Stderr, "dbapi.dropTrigger: No such trigger in DB: '%s'\n", triggerName)
	//}

	return nil
}

// UpdateSchema migrates a 'live' pronlex db to a new schema
// version. The dbFile argument is the path to an Sqlite db file.
func UpdateSchema(dbFile string) error {
	//var err error

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		return fmt.Errorf("dbapi.UpdateSchema: %v", err)
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("dbapi.UpdateSchema: %v", err)
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("dbapi.UpdateSchema: %v", err)
	}

	defer tx.Commit()

	var userVersion int
	err = tx.QueryRow("PRAGMA user_version").Scan(&userVersion)
	if err != nil {

		return fmt.Errorf("dbapi.UpdateSchema: %v", err)
	}
	//fmt.Fprintf(os.Stderr, "dbapi.UpdateSchema: current user_version: %d\n", userVersion)

	if userVersion < 1 {
		// Substitute faulty version of trigger

		//Defined in dbapi.go
		err := dropTrigger(tx, "insertPref")
		if err != nil {
			msg := fmt.Sprintf("drop trigger updatetPref failed : %v", err)
			err2 := tx.Rollback()
			if err2 != nil {
				msg = fmt.Sprintf("%s : failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)
		}

		// Misspelled name of trigger in some version of schema
		err = dropTrigger(tx, "updatetPref")
		if err != nil {
			msg := fmt.Sprintf("drop trigger updatetPref failed : %v", err)
			err2 := tx.Rollback()
			if err2 != nil {
				msg = fmt.Sprintf("%s : failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)
		}

		err = dropTrigger(tx, "updatePref")
		if err != nil {
			msg := fmt.Sprintf("drop trigger updatetPref failed : %v", err)
			err2 := tx.Rollback()
			if err2 != nil {
				msg = fmt.Sprintf("%s : failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)
		}

		err = dropTrigger(tx, "insertEntryStatus")
		if err != nil {
			msg := fmt.Sprintf("drop trigger updatetPref failed : %v", err)
			err2 := tx.Rollback()
			if err2 != nil {
				msg = fmt.Sprintf("%s : failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)
		}

		err = dropTrigger(tx, "updateEntryStatus")
		if err != nil {
			msg := fmt.Sprintf("drop trigger updatetPref failed : %v", err)
			err2 := tx.Rollback()
			if err2 != nil {
				msg = fmt.Sprintf("%s : failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)
		}

		//fmt.Fprintf(os.Stderr, "%s: user_version %d\n", dbFile, userVersion)

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
			err2 := tx.Rollback()
			msg := fmt.Sprintf("failed to create trigger : %v", err)
			if err2 != nil {
				msg = fmt.Sprintf("%s : failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)
		}
		_, err = tx.Exec("PRAGMA user_version = 1")
		if err != nil {
			msg := fmt.Sprintf("failed to set user_version variable : %v", err)
			err2 := tx.Rollback()
			if err2 != nil {
				msg = fmt.Sprintf("%s : failed rollback : %v", msg, err2)
			}
			return fmt.Errorf(msg)

		}

		//fmt.Fprintf(os.Stderr, "dbapi.UpdateSchema: Created triggers\n")
	}

	err = tx.QueryRow("PRAGMA user_version").Scan(&userVersion)
	if err != nil {
		return fmt.Errorf("UpdateSchema: %v", err)
	} // else {
	// 	fmt.Printf("%s user_version %d\n", dbFile, userVersion)
	// }

	return nil
}
