package dbapi

import (
	"database/sql"
	"os"
	"testing"

	"github.com/stts-se/pronlex/lex"
)

func TestPreferred1Sqlite(t *testing.T) {

	dbPath := "./testlex_preferredtag.db"
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		err := os.Remove(dbPath)
		if err != nil {
			t.Errorf("failed to remove %s : %v", dbPath, err)
		}
	}

	db, err := sql.Open("sqlite3_with_regexp", dbPath)
	if err != nil {
		t.Errorf("Failed to open db file %s : %v", dbPath, err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Errorf("Failed to call PRAGMA on db : %v", err)
	}
	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
	if err != nil {
		t.Errorf("Failed to exec PRAGMA call %v", err)
	}
	defer db.Close()

	_, err = execSchemaSqlite(db) // Creates new lexicon database
	if err != nil {

		t.Errorf("Failed to create lexicon db: %v", err)
	}

	l := lexicon{name: "preferred_test", symbolSetName: "ZZ", locale: "ll"}
	l, err = sqliteDBIF{}.defineLexicon(db, l)
	if err != nil {
		t.Errorf("Ooops! : %v", err)
	}

	tx, err := db.Begin()

	defer tx.Commit()
	defer db.Close()
	if err != nil {
		t.Errorf("Failed to start transaction : %v", err)
	}

	// Insert tag for entry that doesn't exist
	err = sqliteDBIF{}.insertEntryTagTx(tx, 0, "ohno")
	//t.Errorf("Error : %v", err)
	if err == nil {
		t.Errorf("Expected error for nonexisting entry id, but got nil")
	}

	tx.Rollback()

	t1 := "city"
	t2 := "drink"
	e1 := lex.Entry{Strn: "rom",
		PartOfSpeech:   "PM",
		WordParts:      "rom",
		Language:       "",
		Preferred:      true,
		Tag:            t1,
		Transcriptions: []lex.Transcription{{Strn: "\" r u m", Language: ""}},
		EntryStatus:    lex.EntryStatus{Name: "unchecked", Source: "imported"}}

	e2 := lex.Entry{Strn: "rom",
		PartOfSpeech:   "NN",
		WordParts:      "rom",
		Language:       "",
		Preferred:      false,
		Tag:            t2,
		Transcriptions: []lex.Transcription{{Strn: "\" r o m", Language: ""}},
		EntryStatus:    lex.EntryStatus{Name: "unchecked", Source: "imported"}}

	// Insert entries
	_, err = sqliteDBIF{}.insertEntries(db, l, []lex.Entry{e1, e2})
	if err != nil {
		t.Errorf("failed to insert entries : %v", err)
	}

	// Fetch inserted entries
	q := Query{Words: []string{"rom"}, Page: 0, PageLength: 25}

	entries, err := sqliteDBIF{}.lookUpIntoMap(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf("Nooo! : %v", err)
	}
	if w, g := 1, len(entries); w != g {
		t.Errorf("Expected '%d' got '%d'", w, g)
	}

	var ent1 lex.Entry
	var ent2 lex.Entry
	seenTag1 := false
	seenTag2 := false
	for _, e := range entries["rom"] {
		if e.Tag == t1 {
			ent1 = e // Save for update test
			seenTag1 = true
		}
		if e.Tag == t2 {
			ent2 = e // Save for update test
			seenTag2 = true
		}
	}
	if !seenTag1 {
		t.Errorf("Couldn't find tag %s in looked up entries %#v", t1, entries)
	}
	if !seenTag2 {
		t.Errorf("Couldn't find tag %s in looked up entries %#v", t2, entries)
	}

	// Verify correct initial pref tags
	if !ent1.Preferred {
		t.Errorf("Expected preferred:true for ent1, found %#v", ent1)
	}
	if ent2.Preferred {
		t.Errorf("Expected preferred:false for ent2, found %#v", ent2)
	}

	// Set language
	ent1.Language = "sv"
	ent1.EntryStatus = lex.EntryStatus{Name: "ok", Source: "hanna"}

	// Update entry with new language
	entUpdate, updated, err := sqliteDBIF{}.updateEntry(db, ent1)
	if err != nil {
		t.Errorf("updateEntry failed : %v", err)
	}
	if !updated {
		t.Errorf("Expected entry to be updated, but nothing happened")
	}
	if !entUpdate.Preferred {
		t.Errorf("Expected updated entry have preferred tag true, but found %#v", entUpdate)
	}

	// Fetch entries again
	q = Query{Words: []string{"rom"}, Page: 0, PageLength: 25}

	entries, err = sqliteDBIF{}.lookUpIntoMap(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf("Nooo! : %v", err)
	}
	if w, g := 1, len(entries); w != g {
		t.Errorf("Expected '%d' got '%d'", w, g)
	}

	seenTag1 = false
	seenTag2 = false
	for _, e := range entries["rom"] {
		if e.Tag == t1 {
			ent1 = e // Save for update test
			seenTag1 = true
		}
		if e.Tag == t2 {
			ent2 = e // Save for update test
			seenTag2 = true
		}
	}
	// Verify that we still have correct pref tags
	if !ent1.Preferred {
		t.Errorf("Expected preferred:true for ent1, found %#v", ent1)
	}
	if ent2.Preferred {
		t.Errorf("Expected preferred:false for ent2, found %#v", ent2)
	}

	// Change preferred
	ent2.Preferred = true

	// Update entry with new language
	entUpdate, updated, err = sqliteDBIF{}.updateEntry(db, ent2)
	if err != nil {
		t.Errorf("updateEntry failed : %v", err)
	}
	if !updated {
		t.Errorf("Expected entry to be updated, but nothing happened")
	}
	if !entUpdate.Preferred {
		t.Errorf("Expected updated entry have preferred tag true, but found %#v", entUpdate)
	}

	// Fetch entries again
	q = Query{Words: []string{"rom"}, Page: 0, PageLength: 25}

	entries, err = sqliteDBIF{}.lookUpIntoMap(db, []lex.LexName{lex.LexName(l.name)}, q)
	if err != nil {
		t.Errorf("Nooo! : %v", err)
	}
	if w, g := 1, len(entries); w != g {
		t.Errorf("Expected '%d' got '%d'", w, g)
	}

	seenTag1 = false
	seenTag2 = false
	for _, e := range entries["rom"] {
		if e.Tag == t1 {
			ent1 = e // Save for update test
			seenTag1 = true
		}
		if e.Tag == t2 {
			ent2 = e // Save for update test
			seenTag2 = true
		}
	}
	// Verify that we have new corrected pref tags
	if ent1.Preferred {
		t.Errorf("Expected preferred:false for ent1, found %#v", ent1)
	}
	if !ent2.Preferred {
		t.Errorf("Expected preferred:true for ent2, found %#v", ent2)
	}

}
