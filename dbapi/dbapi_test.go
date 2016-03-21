package dbapi

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"log"
	"os"
	"regexp"
	"testing"
)

// ff is a place holder to be replaced by proper error handling
func ff(f string, err error) {
	if err != nil {
		log.Fatalf(f, err)
	}
}

func Test_InsertEntries(t *testing.T) {

	err := os.Remove("./testlex.db")
	ff("failed to remove testlex.db : %v", err)

	regex := func(re, s string) (bool, error) {
		return regexp.MatchString(re, s)
	}
	sql.Register("sqlite3_with_go_func",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				return conn.RegisterFunc("regexp", regex, true)
			},
		})

	db, err := sql.Open("sqlite3_with_go_func", "./testlex.db")
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

	l, err = InsertLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	lxs, err := ListLexicons(db)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
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
	entries, err = LookUpIntoMap(db, q) // GetEntries(db, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
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
	tx0, err := db.Begin()
	defer tx0.Commit()
	ff("transaction failed : %v", err)
	le2, err := InsertLemma(tx0, le)
	tx0.Commit()
	if le2.ID < 1 {
		t.Errorf(fs, "more than zero", le2.ID)
	}

	tx00, err := db.Begin()
	ff("tx failed : %v", err)
	defer tx00.Commit()

	le3, err := SetOrGetLemma(tx00, "apa", "67t", "7(c)")
	if le3.ID < 1 {
		t.Errorf(fs, "more than zero", le3.ID)
	}
	tx00.Commit()

	tx01, err := db.Begin()
	ff("tx failed : %v", err)
	defer tx01.Commit()
	err = AssociateLemma2Entry(tx01, le3, entries["apa"][0])
	if err != nil {
		t.Error(fs, nil, err)
	}
	tx01.Commit()

	//ess, err := GetEntries(db, q)
	//var esw EntrySliceWriter
	ess, err := LookUpIntoMap(db, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	if len(ess) != 1 {
		t.Error("ERRRRRRROR")
	}
	lm := ess["apa"][0].Lemma
	if lm.ID < 1 {
		t.Errorf(fs, "id larger than zero", lm.ID)
	}

	if lm.Strn != "apa" {
		t.Errorf(fs, "apa", lm.Strn)
	}
	if lm.Reading != "67t" {
		t.Errorf(fs, "67t", lm.Reading)
	}

	//ees := GetEntriesFromIDs(db, []int64{ess["apa"][0].ID})
	ees, err := LookUpIntoMap(db, Query{EntryIDs: []int64{ess["apa"][0].ID}})
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(ees) != 1 {
		t.Errorf(fs, 1, len(ees))
	}

	// Change transcriptions and update db
	ees0 := ees["apa"][0]
	t10 := Transcription{Strn: "A: p A:", Language: "Apo"}
	t20 := Transcription{Strn: "a p a", Language: "Sweinsprach"}
	t30 := Transcription{Strn: "a pp a", Language: "Mysko"}

	ees0.Transcriptions = []Transcription{t10, t20, t30}

	updated, err := UpdateEntry(db, ees0)

	if !updated {
		t.Errorf(fs, true, updated)
	}
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	eApa, err := GetEntryFromID(db, ees0.ID)
	if len(eApa.Transcriptions) != 3 {
		t.Errorf(fs, 3, len(eApa.Transcriptions))
	}
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	eApa.Lemma.Strn = "tjubba"
	eApa.WordParts = "fin+krog"
	eApa.Language = "gummiapa"
	updated, err = UpdateEntry(db, eApa)
	if err != nil {
		t.Errorf(fs, "nil", err)
	}
	if !updated {
		t.Errorf(fs, true, updated)
	}

	eApax, err := GetEntryFromID(db, ees0.ID)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if eApax.Lemma.Strn != "tjubba" {
		t.Errorf(fs, "tjubba", eApax.Lemma.Strn)
	}
	if eApax.WordParts != "fin+krog" {
		t.Errorf(fs, "fin+krog", eApax.WordParts)
	}
	if eApax.Language != "gummiapa" {
		t.Errorf(fs, "gummiapa", eApax.Language)
	}

	rezz, err := db.Query("select entry.strn from entry where strn regexp '^a'")
	if err != nil {
		log.Fatalf("Agh: %v", err)
	}
	var strn string
	for rezz.Next() {
		rezz.Scan(&strn)
		log.Printf(">>> %s", strn)
	}

}

func Test_unique(t *testing.T) {
	in := []int64{1, 2, 3}

	res := unique(in)
	if len(res) != 3 {
		t.Errorf(fs, 3, len(res))
	}

	in = []int64{3, 3, 3}

	res = unique(in)
	if len(res) != 1 {
		t.Errorf(fs, 1, len(res))
	}
	if res[0] != 3 {
		t.Errorf(fs, 3, res[0])
	}
}
