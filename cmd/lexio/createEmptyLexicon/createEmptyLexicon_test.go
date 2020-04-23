package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

//var fsExpTrans = "Expected: '%v' got: %v'"

func TestProperCloseDontRemoveLockFilesSqlite(t *testing.T) {

	// 1. SETUP
	tmpDir, err := ioutil.TempDir(os.TempDir(), "pronlex-tmp-A")
	fmt.Printf("TMPDIR: %s\n", tmpDir)
	//defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	dbName := "createEmptyLexicon2"
	lexName := "lex_test1"
	ssName := "ss_test1"
	locale := "en_US"
	lexRefX := lex.NewLexRefWithInfo(dbName, lexName, ssName)
	dbRef := lexRefX.LexRef.DBRef
	closeAfter := true
	createDB := true // if not exists

	// 2. CREATE A NEW LEXICON IN A NEW DB
	err = createEmptyLexicon(dbapi.Sqlite, tmpDir, dbRef, lexRefX, locale, createDB, closeAfter)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	// 3. DO NOT REMOVE shm/wal FILES

	// 4. OPEN THE SAME DB
	createDB = false
	lexRefX = lex.NewLexRefWithInfo(string(dbRef), lexName+"-2", ssName)
	err = createEmptyLexicon(dbapi.Sqlite, tmpDir, dbRef, lexRefX, locale, createDB, closeAfter)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}
}

func TestProperCloseRemoveLockFilesSqlite(t *testing.T) {

	// 1. SETUP
	tmpDir, err := ioutil.TempDir(os.TempDir(), "pronlex-tmp-B")
	//fmt.Printf("TMPDIR: %s\n", tmpDir)
	//defer os.RemoveAll(tmpDir)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	dbName := "createEmptyLexicon2"
	lexName := "lex_test1"
	ssName := "ss_test1"
	locale := "en_US"
	lexRefX := lex.NewLexRefWithInfo(dbName, lexName, ssName)
	dbRef := lexRefX.LexRef.DBRef
	closeAfter := false
	createDB := true // if not exists

	// 2. CREATE A NEW LEXICON IN A NEW DB
	err = createEmptyLexicon(dbapi.Sqlite, tmpDir, dbRef, lexRefX, locale, createDB, closeAfter)
	if err != nil {
		t.Errorf("didn't expect error here, found : %v", err)
		return
	}

	// 3. REMOVE shm/wal FILES if they exist (they should not exist if the db has been properly closed)
	dbPath := path.Join(tmpDir, string(dbRef)+".db")
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
	lexRefX = lex.NewLexRefWithInfo(string(dbRef), lexName+"-2", ssName)
	err = createEmptyLexicon(dbapi.Sqlite, tmpDir, dbRef, lexRefX, locale, createDB, closeAfter)
	if err == nil {
		t.Errorf("expected error here, found : %v", err)
		return
	}
}

func init() {
	dbapi.Sqlite3WithRegex()
}
