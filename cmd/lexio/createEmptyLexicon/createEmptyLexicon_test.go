package main

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

//var fsExpTrans = "Expected: '%v' got: %v'"

func TestProperCloseDontRemoveLockFiles(t *testing.T) {

	// 1. SETUP
	tmpDir, err := ioutil.TempDir(os.TempDir(), "pronlex-tmp-A")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	dbPath := path.Join(tmpDir, "createEmptyLexicon1.db")
	lexName := "lex_test1"
	ssName := "ss_test1"
	locale := "en_US"
	lexRefX := lex.NewLexRefWithInfo(dbPath, lexName, ssName)
	closeAfter := true
	createDB := true // if not exists

	// 2. CREATE A NEW LEXICON IN A NEW DB
	err = createEmptyLexicon(dbPath, lexRefX, locale, createDB, closeAfter)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	// 3. DO NOT REMOVE shm/wal FILES

	// 4. OPEN THE SAME DB
	createDB = false
	lexRefX = lex.NewLexRefWithInfo(dbPath, lexName+"-2", ssName)
	err = createEmptyLexicon(dbPath, lexRefX, locale, createDB, closeAfter)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}
}

func TestProperCloseRemoveLockFiles(t *testing.T) {

	// 1. SETUP
	tmpDir, err := ioutil.TempDir(os.TempDir(), "pronlex-tmp-B")
	//defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	dbPath := path.Join(tmpDir, "createEmptyLexicon2.db")
	lexName := "lex_test1"
	ssName := "ss_test1"
	locale := "en_US"
	lexRefX := lex.NewLexRefWithInfo(dbPath, lexName, ssName)
	closeAfter := false
	createDB := true // if not exists

	// 2. CREATE A NEW LEXICON IN A NEW DB
	err = createEmptyLexicon(dbPath, lexRefX, locale, createDB, closeAfter)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	// 3. REMOVE shm/wal FILES if they exist (they should not exist if the db has been properly closed)
	walFile := dbPath + "-wal"
	shmFile := dbPath + "-shm"
	err = os.RemoveAll(walFile)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}
	err = os.RemoveAll(shmFile)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	// 4. OPEN THE SAME DB
	createDB = false
	lexRefX = lex.NewLexRefWithInfo(dbPath, lexName+"-2", ssName)
	err = createEmptyLexicon(dbPath, lexRefX, locale, createDB, closeAfter)
	if err == nil {
		t.Errorf("expected error here, found : %v", err)
		return
	}
}

func init() {
	dbapi.Sqlite3WithRegex()
}
