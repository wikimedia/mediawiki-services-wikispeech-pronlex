package dbapi

import (
	"database/sql"
	"log"
	"os"
	"testing"
)


func Test_InsertEntries(t *testing.T) {

	err := os.Remove("./testlex.db")

	db, err := sql.Open("sqlite3", "./testlex.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
	ff("Failed to exec PRAGMA call %v", err)

	defer db.Close()


	_, err = db.Exec(Schema) // Creates new lexicon database
	ff("Failed to create lexicon db: %v", err)

	// TODO Borde returnera error
	//CreateTables(db, cmds)

	l := Lexicon{Name: "test", SymbolSetName: "ZZ"}

	l = InsertLexicon(db, l)

	lxs := ListLexicons(db)
	if len(lxs) != 1 {
		t.Errorf(fs, 1, len(lxs))
	}

	t1 := Transcription{Strn: "A: p a", Language: "Svetsko"}
	t2 := Transcription{Strn: "a pp a", Language: "svinspråket"}

	e1 := Entry{Strn: "apa",
		PartOfSpeech:   "NN",
		WordParts:      "apa",
		Language:       "XYZZ",
		Transcriptions: []Transcription{t1, t2}}

	InsertEntries(db, l, []Entry{e1})

	// Check that there are things in db:
	q := Query{Words: []string{"apa"}, Page: 0, PageLength: 25}


	var entries map[string][]Entry
	entries = GetEntries(db, q)


	if len(entries) != 1 {
		t.Errorf(fs, 1, len(entries))
	}

	for _, e := range entries {
		ts := len(e[0].Transcriptions)
		if ts != 2 {
			t.Errorf(fs, 2, ts)
		}
	}

	le := Lemma{Strn: "apa", Reading: "67t", Paradigm: "7(c)"}
	le2, err := InsertLemma(db, le)
	if le2.Id < 1 {
		t.Errorf(fs, "more than zero", le2.Id)
	}

	le3, err := SetOrGetLemma(db, "apa", "67t", "7(c)")
	if le3.Id < 1 {
		t.Errorf(fs, "more than zero", le3.Id)
	}

	err = AssociateLemma2Entry(db, le3, entries["apa"][0])
	if err != nil {
		t.Error(fs, nil, err)
	}

	ess := GetEntries(db, q)
	if len(ess) != 1 {
		t.Error("ERRRRRRROR")
	}
	lm := ess["apa"][0].Lemma
	if lm.Id < 1 {
		t.Errorf(fs, "id larger than zero", lm.Id)
	}

	if lm.Strn != "apa" {
		t.Errorf(fs, "apa", lm.Strn)
	}
	if lm.Reading != "67t" {
		t.Errorf(fs, "67t", lm.Reading)
	}

}
