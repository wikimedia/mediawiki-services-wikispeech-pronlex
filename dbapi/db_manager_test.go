package dbapi

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/stts-se/pronlex/lex"
)

func Test_splitFullLexiconName(t *testing.T) {

	n1 := "sv_se-nst:full_words:v1.0"
	db1, l1, err := splitFullLexiconName(n1)
	if w, g := "sv_se-nst", db1; w != g {
		t.Errorf("wanted %s got '%s'", w, g)
	}
	if w, g := "full_words:v1.0", l1; w != g {
		t.Errorf("wanted %s got '%s'", w, g)
	}
	if err != nil {
		t.Errorf("Auch! %v", err)
	}

	// Invalid db name
	n2 := ":full_words:v1.0"
	db2, l2, err := splitFullLexiconName(n2)
	if err == nil {
		t.Errorf("wanted error, got nil")
	}
	if w, g := "", db2; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}
	if w, g := "", l2; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	// Invalid db name
	n2 = "full_wordsv1.0:"
	db2, l2, err = splitFullLexiconName(n2)
	if err == nil {
		t.Errorf("wanted error, got nil")
	}
	if w, g := "", db2; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}
	if w, g := "", l2; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

}

func Test_DBManager(t *testing.T) {

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

	err = dbm.DefineLexicon("db1", "sv_sampa", l1_1, l1_2, l1_3)
	if err != nil {
		t.Errorf("Quack! %v", err)
	}
	err = dbm.DefineLexicon("db2", "sv_sampa", l2_1, l2_2, l2_3)
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
	if w := "db1:zuperlex2"; !lexsM[w] {
		t.Errorf("expected db not found: '%s'", w)
	}

	if w := "db1:zuperlex3"; !lexsM[w] {
		t.Errorf("expected db not found: '%s'", w)
	}

	if w := "db2:zuperlex1"; !lexsM[w] {
		t.Errorf("expected db not found: '%s'", w)
	}
	if w := "db2:zuperlex2"; !lexsM[w] {
		t.Errorf("expected db not found: '%s'", w)
	}

	if w := "db2:zuperduperlex"; !lexsM[w] {
		t.Errorf("expected db not found: '%s'", w)
	}

	//e1 := lex.Entry{Strn: "hus", Transcriptions: []lex.Transcription{lex.Transcription{Strn: `" h u: s`}}}
	t1 := lex.Transcription{Strn: "A: p a", Language: "Svetsko"}
	t2 := lex.Transcription{Strn: "a pp a", Language: "svinspråket"}

	e1 := lex.Entry{Strn: "apa",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apa",
		Language:       "XYZZ",
		Preferred:      true,
		Transcriptions: []lex.Transcription{t1, t2},
		EntryStatus:    lex.EntryStatus{Name: "old1", Source: "tst"}}

	ids, err := dbm.InsertEntries("db2:zuperduperlex", []lex.Entry{e1})
	if w, g := 1, len(ids); w != g {
		t.Errorf("Wanted %v got %v", w, g)
	}
	if err != nil {
		t.Errorf("dbm.InsertEntries: %v", err)
	}

	q := Query{Words: []string{"apa"}}
	lookRes, err := dbm.LookUp([]string{"db2:zuperduperlex"}, q)
	if w, g := 1, len(lookRes); w != g {
		t.Errorf("wanted %d got %d", w, g)
	}
	ents := lookRes["db2"]
	//fmt.Printf("%v\n", ents)
	if w, g := 1, len(ents); w != g {
		t.Errorf("wanted %d got %d", w, g)
	}

	//fmt.Printf("%v\n", lexs)
}
