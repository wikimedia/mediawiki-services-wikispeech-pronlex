package dbapi

import (
	"database/sql"
	"flag"
	//"fmt"
	"time"

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

func execSchema(db *sql.DB) (sql.Result, error) {
	ti := time.Now()
	res, err := db.Exec(Schema)

	_ = ti
	//fmt.Printf("[dbapi_test] db.Exec(Schema) took %v\n", time.Since(ti))
	return res, err
}

func TestMain(m *testing.M) {
	flag.Parse() // should be here
	Sqlite3WithRegex()
	os.Exit(m.Run()) // should be here
}

func Test_SuperDeleteLexicon(t *testing.T) {

	dbPath := "./testlex_superdelete.db"

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

	_, err = execSchema(db) // Creates new lexicon database

	ff("Failed to create lexicon db: %v", err)

	// TODO Borde returnera error
	//CreateTables(db, cmds)

	l := lexicon{name: "test", symbolSetName: "ZZ", locale: "ll"}

	l, err = defineLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	lxs, err := listLexicons(db)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(lxs) != 1 {
		t.Errorf(fs, 1, len(lxs))
	}
	if lxs[0].name != "test" {
		t.Errorf(fs, "test", lxs[0].name)
	}
	if lxs[0].id <= 0 {
		t.Errorf(fs, ">0", lxs[0].id)
	}
	if lxs[0].symbolSetName != "ZZ" {
		t.Errorf(fs, "ZZ", lxs[0].symbolSetName)
	}

	lx, err := getLexicon(db, "test")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if w, g := "test", lx.name; w != g {
		t.Errorf("Wanted %s got %s", w, g)
	}
	if w, g := "ZZ", lx.symbolSetName; w != g {
		t.Errorf("Wanted %s got %s", w, g)
	}
	lx, err = getLexicon(db, "xyzzhga_skdjdj")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if w, g := "", lx.name; w != g {
		t.Errorf("Wanted empty string, got '%s'", g)
	}

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

	t1 = lex.Transcription{Strn: "A: p a n", Language: "Svetsko"}
	t2 = lex.Transcription{Strn: "a pp a n", Language: "svinspråket"}

	e2 := lex.Entry{Strn: "apan",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apa",
		Language:       "XYZZ",
		Preferred:      true,
		Transcriptions: []lex.Transcription{t1, t2},
		EntryStatus:    lex.EntryStatus{Name: "old1", Source: "tst"}}

	_, errx := insertEntries(db, l, []lex.Entry{e1, e2})
	if errx != nil {
		t.Errorf(fs, "nil", errx)
	}

	// Check that there are things in db:
	q := Query{Page: 0, PageLength: 25}

	entries, err := lookUpIntoMap(db, []lex.LexName{lex.LexName(l.name)}, q) // GetEntries(db, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if got, want := len(entries), 2; got != want {
		t.Errorf(fs, got, want)
	}
	lexes, err := listLexicons(db)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if got, want := len(lexes), 1; got != want {
		t.Errorf(fs, got, want)
	}

	superDeleteLexicon(db, "test")

	// check that the lexicon named 'test' is deleted
	lexes, err = listLexicons(db)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if got, want := len(lexes), 0; got != want {
		t.Errorf(fs, got, want)
	}

}

func Test_insertEntries(t *testing.T) {

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

	_, err = execSchema(db) // Creates new lexicon database

	ff("Failed to create lexicon db: %v", err)

	// TODO Borde returnera error
	//CreateTables(db, cmds)

	l := lexicon{name: "test", symbolSetName: "ZZ", locale: "ll"}

	l, err = defineLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	lxs, err := listLexicons(db)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(lxs) != 1 {
		t.Errorf(fs, 1, len(lxs))
	}
	if lxs[0].name != "test" {
		t.Errorf(fs, "test", lxs[0].name)
	}
	if lxs[0].id <= 0 {
		t.Errorf(fs, ">0", lxs[0].id)
	}
	if lxs[0].symbolSetName != "ZZ" {
		t.Errorf(fs, "ZZ", lxs[0].symbolSetName)
	}

	lx, err := getLexicon(db, "test")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	if w, g := "test", lx.name; w != g {
		t.Errorf("Wanted %s got %s", w, g)
	}
	if w, g := "ZZ", lx.symbolSetName; w != g {
		t.Errorf("Wanted %s got %s", w, g)
	}
	lx, err = getLexicon(db, "xyzzhga_skdjdj")
	if err == nil {
		t.Error("Expected error, got nil")
	}
	if w, g := "", lx.name; w != g {
		t.Errorf("Wanted empty string, got '%s'", g)
	}

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

	_, errx := insertEntries(db, l, []lex.Entry{e1})
	if errx != nil {
		t.Errorf(fs, "nil", errx)
	}
	// Check that there are things in db:
	q := Query{Words: []string{"apa"}, Page: 0, PageLength: 25}

	var entries map[string][]lex.Entry
	entries, err = lookUpIntoMap(db, []lex.LexName{lex.LexName(l.name)}, q) // GetEntries(db, q)
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
	if got, want := ea.Preferred, true; got != want {
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
	le2, err := insertLemma(tx0, le)
	tx0.Commit()
	if le2.ID < 1 {
		t.Errorf(fs, "more than zero", le2.ID)
	}

	que := Query{TranscriptionLike: "%pp%"}
	var queRez lex.EntrySliceWriter
	err = lookUp(db, []lex.LexName{lex.LexName(l.name)}, que, &queRez)
	if err != nil {
		t.Errorf("Wanted nil, got %v", err)
	}
	if got, want := len(queRez.Entries), 1; got != want {
		t.Errorf("Got: %v Wanted: %v", got, want)
	}

	tx00, err := db.Begin()
	ff("tx failed : %v", err)
	defer tx00.Commit()

	le3, err := setOrGetLemma(tx00, "apa", "67t", "7(c)")
	if le3.ID < 1 {
		t.Errorf(fs, "more than zero", le3.ID)
	}
	tx00.Commit()

	tx01, err := db.Begin()
	ff("tx failed : %v", err)
	defer tx01.Commit()
	err = associateLemma2Entry(tx01, le3, entries["apa"][0])
	if err != nil {
		t.Error(fs, nil, err)
	}
	tx01.Commit()

	//ess, err := GetEntries(db, q)
	//var esw lex.EntrySliceWriter
	ess, err := lookUpIntoMap(db, []lex.LexName{lex.LexName(l.name)}, q)
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
	ees, err := lookUpIntoMap(db, []lex.LexName{lex.LexName(l.name)}, Query{EntryIDs: []int64{ess["apa"][0].ID}})
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(ees) != 1 {
		t.Errorf(fs, 1, len(ees))
	}

	// Check that no entries with entryvalidation exist
	noev, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, Query{HasEntryValidation: true})
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
	ees0.EntryValidations = []lex.EntryValidation{{Level: "severe", RuleName: "barf", Message: "it hurts"}}

	newE, updated, err := updateEntry(db, ees0)

	if err != nil {
		t.Errorf(fs, nil, err)
	}

	if !updated {
		t.Errorf(fs, true, updated)
	}

	if want, got := true, newE.Strn == ees0.Strn; !got {
		t.Errorf(fs, got, want)
	}

	eApa, err := getEntryFromID(db, ees0.ID)
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
	noev, err = lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, Query{HasEntryValidation: true})
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
	newE2, updated, err := updateEntry(db, eApa)
	if err != nil {
		t.Errorf(fs, "nil", err)
	}
	if !updated {
		t.Errorf(fs, true, updated)
	}
	if want, got := true, newE2.Strn == eApa.Strn; !got {
		t.Errorf(fs, got, want)
	}

	eApax, err := getEntryFromID(db, ees0.ID)
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
	noev, err = lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, Query{HasEntryValidation: true})
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

	// Add another entry with same str as existing one, to test preferred
	e1b := lex.Entry{Strn: "apa",
		PartOfSpeech:   "XX",
		Morphology:     "NEU UTR",
		WordParts:      "apa",
		Language:       "XYZZ",
		Preferred:      true,
		Transcriptions: []lex.Transcription{t1, t2},
		EntryStatus:    lex.EntryStatus{Name: "old2", Source: "tst"}}

	_, errxb := insertEntries(db, l, []lex.Entry{e1b})
	if errxb != nil {
		t.Errorf("Failed to insert entry: %v", errxb)
	}

	// Check that only the new entry has Preferred == 1
	//q := Query{Words: []string{"apa"}, Page: 0, PageLength: 25}

	var entries2 []lex.Entry
	entries2, err = lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	//fmt.Printf("%#v\n", entries2[0])
	//fmt.Printf("%#v\n", entries2[1])
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if got, want := len(entries2), 2; got != want {
		t.Errorf(fs, want, got)
	}

	if entries2[0].PartOfSpeech == "XX" && !entries2[0].Preferred {
		t.Errorf(fs, 1, entries2[0].Preferred)
	}
	if entries2[1].PartOfSpeech == "XX" && !entries2[1].Preferred {
		t.Errorf(fs, 1, entries2[1].Preferred)
	}
	if entries2[0].PartOfSpeech != "XX" && entries2[0].Preferred {
		t.Errorf(fs, 0, entries2[0].Preferred)
	}
	if entries2[1].PartOfSpeech != "XX" && entries2[1].Preferred {
		t.Errorf(fs, 0, entries2[1].Preferred)
	}

	// TODO should be in a test of its own
	eStatsus, err := listCurrentEntryStatuses(db, l.name)
	if err != nil {
		t.Errorf("%v", err)
	}
	if w, g := 2, len(eStatsus); w != g {
		t.Errorf(fs, w, g)
	}
	eStatsus2, err := listAllEntryStatuses(db, l.name)
	if err != nil {
		t.Errorf("%v", err)
	}
	if w, g := 3, len(eStatsus2); w != g {
		t.Errorf(fs, w, g)
	}

	stat, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, Query{EntryStatus: []string{"new"}})
	if err != nil {
		t.Errorf("%v", err)
	}
	if w, g := 1, len(stat); w != g {
		t.Errorf(fs, w, g)
	}
	stat1, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, Query{EntryStatus: []string{"dkhfkhekjeh"}})
	if err != nil {
		t.Errorf("%v", err)
	}
	if w, g := 0, len(stat1); w != g {
		t.Errorf(fs, w, g)
	}
	stat2, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, Query{EntryStatus: []string{"new", "old2"}})
	if err != nil {
		t.Errorf("%v", err)
	}
	if w, g := 2, len(stat2); w != g {
		t.Errorf(fs, w, g)
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

func Test_ImportLexiconFile(t *testing.T) {

	symbolSet, err := symbolset.LoadSymbolSet("./../symbolset/test_data/sv-se_ws-sampa.sym")
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
	//defer db.Commit()
	defer db.Close()

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
	ff("Failed to exec PRAGMA call %v", err)

	defer db.Close()

	_, err = execSchema(db) // Creates new lexicon database
	ff("Failed to create lexicon db: %v", err)

	logger := StderrLogger{}
	l := lexicon{name: "test", symbolSetName: symbolSet.Name, locale: "ll"}

	l, err = defineLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	// actual tests start here
	err = ImportLexiconFile(db, lex.LexName(l.name), logger, "./sv-lextest.txt", &validation.Validator{})
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	q := Query{Words: []string{"sprängstoff"}}

	res, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o := res[0].Strn
	if o != "sprängstoff" {
		t.Errorf(fs, "sprängstoff", o)
	}

	q = Query{Words: []string{"sittriktiga"}}
	res, err = lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o = res[0].Strn
	if o != "sittriktiga" {
		t.Errorf(fs, "sittriktiga", o)
	}

	//Let's throw in a test of deleteEntry as well:
	eX := res[0]
	deleteEntry(db, eX.ID, l.name)

	// Run same query again, efter deleting Entry
	resX, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(resX) != 0 {
		t.Errorf(fs, "0", len(res))
	}

}

func Test_ImportLexiconFileWithDupLines(t *testing.T) {

	symbolSet, err := symbolset.LoadSymbolSet("./../symbolset/test_data/sv-se_ws-sampa.sym")
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

	_, err = execSchema(db) // Creates new lexicon database
	ff("Failed to create lexicon db: %v", err)

	logger := StderrLogger{}
	l := lexicon{name: "test", symbolSetName: symbolSet.Name, locale: "ll"}

	l, err = defineLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	// actual tests start here
	err = ImportLexiconFile(db, lex.LexName(l.name), logger, "./sv-lextest-dups.txt", &validation.Validator{})
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	q := Query{Words: []string{"sprängstoff"}}

	res, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o := res[0].Strn
	if o != "sprängstoff" {
		t.Errorf(fs, "sprängstoff", o)
	}

	q = Query{Words: []string{"sittriktiga"}}
	res, err = lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o = res[0].Strn
	if o != "sittriktiga" {
		t.Errorf(fs, "sittriktiga", o)
	}

	q = Query{Words: []string{"vadare"}}
	res, err = lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o = res[0].Strn
	if o != "vadare" {
		t.Errorf(fs, "vadare", o)
	}

	q = Query{}
	res, err = lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(res) != 19 {
		t.Errorf(fs, "19", len(res))
	}
}

func Test_ImportLexiconFileInvalid(t *testing.T) {

	symbolSet, err := symbolset.LoadSymbolSet("./../symbolset/test_data/sv-se_ws-sampa.sym")
	//ssMapper, err := symbolset.LoadSymbolSet("./../symbolset/test_data/sv-se_ws-sampa.sym")
	if err != nil {
		log.Fatal(err)
	}

	//symbolSet := ssMapper.From

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

	_, err = execSchema(db) // Creates new lexicon database
	ff("Failed to create lexicon db: %v", err)

	logger := StderrLogger{}
	l := lexicon{name: "test", symbolSetName: symbolSet.Name, locale: "ll"}

	l, err = defineLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	// actual tests start here
	err = ImportLexiconFile(db, lex.LexName(l.name), logger, "./sv-lextest-invalid-no-fields.txt", &validation.Validator{})
	if err == nil {
		t.Errorf("Expected errors, but got nil")
	}

}

func Test_ImportLexiconFileGz(t *testing.T) {

	symbolSet, err := symbolset.LoadSymbolSet("./../symbolset/test_data/sv-se_ws-sampa.sym")
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

	_, err = execSchema(db) // Creates new lexicon database
	ff("Failed to create lexicon db: %v", err)

	logger := StderrLogger{}
	l := lexicon{name: "test", symbolSetName: symbolSet.Name, locale: "ll"}

	l, err = defineLexicon(db, l)
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	// actual tests start here
	err = ImportLexiconFile(db, lex.LexName(l.name), logger, "./sv-lextest.txt.gz", &validation.Validator{})
	if err != nil {
		t.Errorf(fs, nil, err)
	}

	q := Query{Words: []string{"sprängstoff"}}

	res, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o := res[0].Strn
	if o != "sprängstoff" {
		t.Errorf(fs, "sprängstoff", o)
	}

	q = Query{Words: []string{"sittriktiga"}}
	res, err = lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(res) != 1 {
		t.Errorf(fs, "1", len(res))
	}
	o = res[0].Strn
	if o != "sittriktiga" {
		t.Errorf(fs, "sittriktiga", o)
	}

	//Let's throw in a test of deleteEntry as well:
	eX := res[0]
	deleteEntry(db, eX.ID, l.name)

	// Run same query again, efter deleting Entry
	resX, err := lookUpIntoSlice(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf(fs, nil, err)
	}
	if len(resX) != 0 {
		t.Errorf(fs, "0", len(res))
	}

}
