package dbapi

import (
	"database/sql"
	"log"
	"os"
	"testing"
	//"time"

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

	l3 := Lexicon{Name: "test3", SymbolSetName: "ZZ"}
	l3, err = InsertLexicon(db, l3)
	if err != nil {
		t.Errorf("holy cow (3)! : %v", err)
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

	if w, g := int64(0), res.N; w != g {
		t.Errorf("Wanted '%d' got '%d'", w, g)
	}

	// Add entry unique to l1, and this should be movable

	t2 := lex.Transcription{Strn: `"" f I N . e . % rl i: . k a`}
	e2 := lex.Entry{
		Strn:           "fingerlikas",
		PartOfSpeech:   "JJ",
		Morphology:     "SIN-PLU|IND-DEF|NOM|UTR-NEU|POS|GEN",
		Language:       "sv",
		Preferred:      true,
		Transcriptions: []lex.Transcription{t2},
		EntryStatus:    lex.EntryStatus{Name: "newEntry", Source: "testSource"},
	}

	_, err = InsertEntries(db, l1, []lex.Entry{e2})
	if err != nil {
		t.Errorf("The horror, the horror : %v", err)
	}

	// Insert the same entry in "unrelated" third lexicon, to or from which nothing should be moved
	_, err = InsertEntries(db, l3, []lex.Entry{e2})
	if err != nil {
		t.Errorf("Unbelievable! : %v", err)
	}

	res2, err := MoveNewEntries(db, l1.Name, l2.Name, "from:"+l1.Name, "moved")
	if err != nil {
		t.Errorf("No fun : %v", err)
	}
	if w, g := int64(1), res2.N; w != g {
		t.Errorf("wanted %v got %v", w, g)
	}

	statsL1, err := LexiconStats(db, l1.ID)
	if err != nil {
		t.Errorf("didn't expect that : %v", err)
	}
	if w, g := int64(1), statsL1.Entries; w != g {
		t.Errorf("wanted %v got %v", w, g)
	}
	statsL2, err := LexiconStats(db, l2.ID)
	if err != nil {
		t.Errorf("didn't expect that : %v", err)
	}
	if w, g := int64(2), statsL2.Entries; w != g {
		t.Errorf("wanted %v got %v", w, g)
	}

	// Move back again
	res3, err := MoveNewEntries(db, l2.Name, l1.Name, "from:"+l2.Name, "moved_back")
	if err != nil {
		t.Errorf("No fun : %v", err)
	}
	if w, g := int64(1), res3.N; w != g {
		t.Errorf("wanted %v got %v", w, g)
	}

	statsL1b, err := LexiconStats(db, l1.ID)
	if err != nil {
		t.Errorf("didn't expect that : %v", err)
	}
	if w, g := int64(2), statsL1b.Entries; w != g {
		t.Errorf("wanted %v got %v", w, g)
	}
	statsL2b, err := LexiconStats(db, l2.ID)
	if err != nil {
		t.Errorf("didn't expect that : %v", err)
	}
	if w, g := int64(1), statsL2b.Entries; w != g {
		t.Errorf("wanted %v got %v", w, g)
	}

	statsL3, err := LexiconStats(db, l3.ID)
	if err != nil {
		t.Errorf("didn't expect that : %v", err)
	}
	if w, g := int64(1), statsL3.Entries; w != g {
		t.Errorf("wanted %v got %v", w, g)
	}

}