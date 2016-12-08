package dbapi

import (
	"database/sql"
	"flag"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
	//"github.com/mattn/go-sqlite3"
	"log"
	"os"
	//"regexp"
	"testing"
)

// ff is a place holder to be replaced by proper error handling
func ff(f string, err error) {
	if err != nil {
		log.Fatalf(f, err)
	}
}

func TestMain(m *testing.M) {
	flag.Parse() // should be here
	Sqlite3WithRegex()
	os.Exit(m.Run()) // should be here
}

func Test_InsertEntries(t *testing.T) {

	dbPath := "./testlex.db"

	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		err := os.Remove(dbPath)
		ff("failed to remove "+dbPath+" : %v", err)
	}

	db, err := sql.Open("sqlite3_with_regexp", dbPath)
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

	t1 := lex.Transcription{Strn: "A: p a", Language: "Svetsko"}
	t2 := lex.Transcription{Strn: "a pp a", Language: "svinspråket"}

	e1 := lex.Entry{Strn: "apa",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apa",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t1, t2},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	_, errx := InsertEntries(db, l, []lex.Entry{e1})
	if errx != nil {
		t.Errorf(fs, "nil", errx)
	}
	// Check that there are things in db:
	q := Query{Words: []string{"apa"}, Page: 0, PageLength: 25}

	var entries map[string][]lex.Entry
	entries, err = LookUpIntoMap(db, q) // GetEntries(db, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if got, want := len(entries), 1; got != want {
		t.Errorf(fs, got, want)
	}

	ea := entries["apa"][0]
	if got, want := ea.Morphology, "NEU UTR"; got != want {
		t.Errorf(fs, got, want)
	}

	for _, e := range entries {
		ts := len(e[0].Transcriptions)
		if ts != 2 {
			t.Errorf(fs, 2, ts)
		}
	}

	le := lex.Lemma{Strn: "apa", Reading: "67t", Paradigm: "7(c)"}
	tx0, err := db.Begin()
	defer tx0.Commit()
	ff("transaction failed : %v", err)
	le2, err := InsertLemma(tx0, le)
	tx0.Commit()
	if le2.ID < 1 {
		t.Errorf(fs, "more than zero", le2.ID)
	}

	que := Query{TranscriptionLike: "%pp%"}
	var queRez lex.EntrySliceWriter
	err = LookUp(db, que, &queRez)
	if err != nil {
		t.Errorf("Wanted nil, got %v", err)
	}
	if got, want := len(queRez.Entries), 1; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
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
	//var esw lex.EntrySliceWriter
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

	// Check that no entries with entryvalidation exist
	noev, err := LookUpIntoSlice(db, Query{HasEntryValidation: true})
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if got, want := len(noev), 0; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	// Change transcriptions and update db
	ees0 := ees["apa"][0]
	t10 := lex.Transcription{Strn: "A: p A:", Language: "Apo"}
	t10.AddSource("orangu1")
	t20 := lex.Transcription{Strn: "a p a", Language: "Sweinsprach"}
	t20.AddSource("orangu2")
	t30 := lex.Transcription{Strn: "a pp a", Language: "Mysko"}
	t30.AddSource("orangu3")
	t30.AddSource("orangu4")
	ees0.Transcriptions = []lex.Transcription{t10, t20, t30}
	// add new lex.EntryStatus
	ees0.EntryStatus = lex.EntryStatus{Name: "new", Source: "tst"}
	// new validation
	ees0.EntryValidations = []lex.EntryValidation{lex.EntryValidation{Level: "severe", RuleName: "barf", Message: "it hurts"}}

	newE, updated, err := UpdateEntry(db, ees0)

	if !updated {
		t.Errorf(fs, true, updated)
	}
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	if want, got := true, newE.Strn == ees0.Strn; !got {
		t.Errorf(fs, got, want)
	}

	eApa, err := GetEntryFromID(db, ees0.ID)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(eApa.Transcriptions) != 3 {
		t.Errorf(fs, 3, len(eApa.Transcriptions))
	}

	if got, want := eApa.Transcriptions[0].Sources[0], "orangu1"; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	if got, want := len(eApa.Transcriptions[2].Sources), 2; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}
	if got, want := eApa.Transcriptions[2].Sources[0], "orangu4"; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	if got, want := eApa.EntryStatus.Name, "new"; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	if got, want := len(eApa.EntryValidations), 1; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}
	if got, want := eApa.EntryValidations[0].Level, "severe"; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	if got, want := eApa.EntryValidations[0].RuleName, "barf"; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}
	if got, want := eApa.EntryValidations[0].Message, "it hurts"; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	// Check that one entry with entryvalidation exists
	noev, err = LookUpIntoSlice(db, Query{HasEntryValidation: true})
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if got, want := len(noev), 1; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	eApa.Lemma.Strn = "tjubba"
	eApa.WordParts = "fin+krog"
	eApa.Language = "gummiapa"
	eApa.EntryValidations = []lex.EntryValidation{}
	newE2, updated, err := UpdateEntry(db, eApa)
	if err != nil {
		t.Errorf(fs, "nil", err)
	}
	if !updated {
		t.Errorf(fs, true, updated)
	}
	if want, got := true, newE2.Strn == eApa.Strn; !got {
		t.Errorf(fs, got, want)
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
	if got, want := len(eApax.EntryValidations), 0; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	// Check that no entries with entryvalidation exist
	noev, err = LookUpIntoSlice(db, Query{HasEntryValidation: true})
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if got, want := len(noev), 0; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	// rezz, err := db.Query("select entry.strn from entry where strn regexp '^a'")
	// if err != nil {
	// 	log.Fatalf("Agh: %v", err)
	// }
	// var strn string
	// for rezz.Next() {
	// 	rezz.Scan(&strn)
	// 	log.Printf(">>> %s", strn)
	// }

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

func Test_ImportLexiconFile(t *testing.T) {

	symbolSet, err := symbolset.LoadSymbolSet("./../symbolset/static/sv-se_ws-sampa.tab")
	if err != nil {
		log.Fatal(err)
	}

	dbFile := "./iotestlex.db"
	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		err := os.Remove(dbFile)
		ff("failed to remove iotestlex.db : %v", err)
	}

	db, err := sql.Open("sqlite3_with_regexp", "./iotestlex.db")
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

	logger := StderrLogger{}
	l := Lexicon{Name: "test", SymbolSetName: symbolSet.Name}

	l, err = InsertLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	// actual tests start here
	err = ImportLexiconFile(db, logger, l.Name, "./sv-lextest.txt", &validation.Validator{})
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	q := Query{Words: []string{"sprängstoff"}}

	res, err := LookUpIntoSlice(db, q)
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o := res[0].Strn
	if o != "sprängstoff" {
		t.Errorf(fs, "sprängstoff", o)
	}

	q = Query{Words: []string{"sittriktiga"}}
	res, err = LookUpIntoSlice(db, q)
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o = res[0].Strn
	if o != "sittriktiga" {
		t.Errorf(fs, "sittriktiga", o)
	}

}

// func Test_ImportLexiconFileInvalid(t *testing.T) {

// 	ssMapper, err := symbolset.LoadSymbolSet("./../symbolset/static/sv-se_ws-sampa.tab")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	symbolSet := ssMapper.From

// 	dbFile := "./iotestlex.db"
// 	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
// 		err := os.Remove(dbFile)
// 		ff("failed to remove iotestlex.db : %v", err)
// 	}

// 	db, err := sql.Open("sqlite3_with_regexp", "./iotestlex.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	_, err = db.Exec("PRAGMA foreign_keys = ON")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
// 	ff("Failed to exec PRAGMA call %v", err)

// 	defer db.Close()

// 	_, err = db.Exec(Schema) // Creates new lexicon database
// 	ff("Failed to create lexicon db: %v", err)

// 	logger := StderrLogger{}
// 	l := Lexicon{Name: "test", SymbolSetName: symbolSet.Name}

// 	l, err = InsertLexicon(db, l)
// 	if err != nil {
// 		t.Errorf(fs, nil, err)
// 	}

// 	// actual tests start here
// 	errs := ImportLexiconFile(db, logger, l.Name, "./sv-lextest-invalid.txt")
// 	if len(errs) != 1 {
// 		t.Errorf(fs, nil, errs)
// 	}

// 	q := Query{Words: []string{"sprängstoff"}}

// 	res, err := LookUpIntoSlice(db, q)
// 	if len(res) != 1 {
// 		t.Errorf(fs, "1", len(res))
// 	}
// 	o := res[0].Strn
// 	if o != "sprängstoff" {
// 		t.Errorf(fs, "sprängstoff", o)
// 	}

// 	q = Query{Words: []string{"sittriktigas"}}
// 	res, err = LookUpIntoSlice(db, q)
// 	if len(res) != 1 {
// 		t.Errorf(fs, "1", len(res))
// 		return
// 	}
// 	o = res[0].Strn
// 	if o != "sittriktigas" {
// 		t.Errorf(fs, "sittriktigas", o)
// 	}

// }

func Test_ImportLexiconFileGz(t *testing.T) {

	symbolSet, err := symbolset.LoadSymbolSet("./../symbolset/static/sv-se_ws-sampa.tab")
	if err != nil {
		log.Fatal(err)
	}

	dbFile := "./iotestlex.db"
	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		err := os.Remove(dbFile)
		ff("failed to remove iotestlex.db : %v", err)
	}

	db, err := sql.Open("sqlite3_with_regexp", "./iotestlex.db")
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

	logger := StderrLogger{}
	l := Lexicon{Name: "test", SymbolSetName: symbolSet.Name}

	l, err = InsertLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	// actual tests start here
	err = ImportLexiconFile(db, logger, l.Name, "./sv-lextest.txt.gz", &validation.Validator{})
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	q := Query{Words: []string{"sprängstoff"}}

	res, err := LookUpIntoSlice(db, q)
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o := res[0].Strn
	if o != "sprängstoff" {
		t.Errorf(fs, "sprängstoff", o)
	}

	q = Query{Words: []string{"sittriktiga"}}
	res, err = LookUpIntoSlice(db, q)
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o = res[0].Strn
	if o != "sittriktiga" {
		t.Errorf(fs, "sittriktiga", o)
	}

}
