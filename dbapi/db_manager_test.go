package dbapi

import (
	"database/sql"
	//"fmt"
	"log"
	"os"
	"testing"
)

func Test_ListLexicon(t *testing.T) {

	dbPath1 := "./testlex_listlex1.db"
	dbPath2 := "./testlex_listlex2.db"

	if _, err := os.Stat(dbPath1); !os.IsNotExist(err) {
		err := os.Remove(dbPath1)

		if err != nil {
			log.Fatalf("failed to remove '%s' : %v", dbPath1, err)
		}
	}

	if _, err := os.Stat(dbPath2); !os.IsNotExist(err) {
		err := os.Remove(dbPath2)

		if err != nil {
			log.Fatalf("failed to remove '%s' : %v", dbPath2, err)
		}
	}

	db1, err := sql.Open("sqlite3_with_regexp", dbPath1)
	if err != nil {
		log.Fatal(err)
	}
	db2, err := sql.Open("sqlite3_with_regexp", dbPath2)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db1.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db2.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db1.Exec("PRAGMA case_sensitive_like=ON")
	ff("Failed to exec PRAGMA call %v", err)
	_, err = db2.Exec("PRAGMA case_sensitive_like=ON")
	ff("Failed to exec PRAGMA call %v", err)

	defer db1.Close()
	defer db2.Close()

	_, err = execSchema(db1) // Creates new lexicon database
	if err != nil {
		log.Fatalf("NO! %v", err)
	}
	_, err = execSchema(db2) // Creates new lexicon database
	if err != nil {
		log.Fatalf("NO! %v", err)
	}

	dbm := NewDBManager()
	dbm.AddDB("db1", db1)
	dbm.AddDB("db2", db2)

	l1_1 := "zuperlex1"
	l1_2 := "zuperlex2"
	l1_3 := "zuperlex3"

	l2_1 := "zuperlex1"
	l2_2 := "zuperlex2"
	l2_3 := "zuperduperlex"

	err = dbm.DefineLexicon("db1", l1_1, l1_2, l1_3)
	if err != nil {
		t.Errorf("Quack! %v", err)
	}
	err = dbm.DefineLexicon("db2", l2_1, l2_2, l2_3)
	if err != nil {
		t.Errorf("Quack! %v", err)
	}

	lexs, err := dbm.ListLexicons()
	if err != nil {
		t.Errorf("Quack! %v", err)
	}

	if w, g := 6, len(lexs); w != g {
		t.Errorf("wanted %v got %v", w, g)
	}

	lexsM := make(map[string]bool)
	for _, l := range lexs {
		lexsM[l] = true
	}
	if w := "db1:zuperlex1"; !lexsM[w] {
		t.Errorf("expected db not found: '%s'", w)
	}
	if w := "db2:zuperlex1"; !lexsM[w] {
		t.Errorf("expected db not found: '%s'", w)
	}

	//fmt.Printf("%v\n", lexs)
}
