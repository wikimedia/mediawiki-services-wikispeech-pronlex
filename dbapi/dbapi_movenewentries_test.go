package dbapi

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/stts-se/pronlex/lex"
)

func Test_MoveNewEntries(t *testing.T) {

	dbFile := "./movetestlex.db"
	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		err0 := os.Remove(dbFile)
		if err0 != nil {
			log.Fatalf("failed to remove %s : %v", dbFile, err0)
		}
	}

	db, err := sql.Open("sqlite3_with_regexp", dbFile)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
	if err != nil {
		log.Fatalf("Failed to exec PRAGMA call %v", err)
	}

	defer db.Close()

	_, err = execSchema(db) // Creates new lexicon database
	if err != nil {
		log.Fatalf("Failed to create lexicon db: %v", err)
	}

	l1 := Lexicon{Name: "test1", SymbolSetName: "ZZ"}
	l1, err = InsertLexicon(db, l1)
	if err != nil {
		t.Errorf("holy cow (1)! : %v", err)
	}

	l2 := Lexicon{Name: "test2", SymbolSetName: "ZZ"}
	l2, err = InsertLexicon(db, l2)
	if err != nil {
		t.Errorf("holy cow (2)! : %v", err)
	}

	t1 := lex.Transcription{Strn: `"" f I N . e . % rl i: . k a`}
	e1 := lex.Entry{
		Strn:           "fingerlika",
		PartOfSpeech:   "JJ",
		Morphology:     "SIN-PLU|IND-DEF|NOM|UTR-NEU|POS",
		Language:       "sv",
		Preferred:      true,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "newEntry", Source: "testSource"},
	}

	// Same entry in both lexica, nothing should be moved
	_, err = InsertEntries(db, l1, []lex.Entry{e1})
	if err != nil {
		t.Errorf("The sky is falling! : %v", err)
	}
	_, err = InsertEntries(db, l2, []lex.Entry{e1})
	if err != nil {
		t.Errorf("The sky is falling! : %v", err)
	}

	res, err := MoveNewEntries(db, l1.Name, l2.Name, "from"+l1.Name, "moved")
	if err != nil {
		t.Errorf("What?! : %v", err)
	}

	if w, g := int64(0), res.n; w != g {
		t.Errorf("Wanted '%d' got '%d'", w, g)
	}

}
