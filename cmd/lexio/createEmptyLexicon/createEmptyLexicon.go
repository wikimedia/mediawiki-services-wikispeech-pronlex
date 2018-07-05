// createEmptyLexicon initialises an Sqlite3 relational database from the schema defining a lexicon database, but empty of data.
// See dbapi.Schema.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

func createEmptyLexicon(dbPath string, lexRefX lex.LexRefWithInfo, locale string, createDbIfNotExists bool, closeAfter bool) error {
	dbm := dbapi.NewDBManager()
	lexRef := lexRefX.LexRef // db + lexname without ssName

	if closeAfter {
		defer dbm.CloseDB(lexRef.DBRef)
	}

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if createDbIfNotExists {
			err := dbm.DefineSqliteDB(lexRef.DBRef, dbPath)
			if err != nil {
				return fmt.Errorf("couldn't create db %s : %v", dbPath, err)
			}
			fmt.Fprintf(os.Stderr, "[createEmptyLexicon] created db %s\n", dbPath)
		} else {
			return fmt.Errorf("db does not exist %s : %v", dbPath, err)
		}
	} else {
		err := dbm.OpenDB(lexRef.DBRef, dbPath)
		if err != nil {
			return fmt.Errorf("couldn't open db %s : %v", dbPath, err)
		}
		fmt.Fprintf(os.Stderr, "[createEmptyLexicon] opened db %s\n", dbPath)
	}

	if !dbm.ContainsDB(lexRef.DBRef) {
		return fmt.Errorf("db should be registered in dbm, but wasn't : %s", lexRef.DBRef)
	}

	exists, err := dbm.LexiconExists(lexRef)
	if err != nil {
		return fmt.Errorf("couldn't lookup lexicon reference %s : %v", lexRef.String(), err)
	}
	if exists {
		return fmt.Errorf("couldn't create lexicon that already exists: %s", lexRef.String())
	}

	err = dbm.DefineLexicon(lexRefX.LexRef, lexRefX.SymbolSetName, locale)
	if err != nil {
		return fmt.Errorf("couldn't create lexicon : %v", err)
	}
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon] created lexicon %s in %s\n", lexRefX.LexRef.LexName, dbPath)
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon]  > symbolset: %s\n", lexRefX.SymbolSetName)
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon]  > locale:    %s\n", locale)

	return nil
}

func main() {

	var createDbIfNotExists1 = flag.Bool("createdb", false, "create db if it doesn't exist")
	var createDbIfNotExists2 = flag.Bool("c", false, "create db if it doesn't exist")
	var help = flag.Bool("help", false, "print help and exit")

	usage := `Usage:
     $ createEmptyLexicon [flags] <DB PATH> <LEX NAME> <SYMBOL SET NAME> <LOCALE>

Flags:
     -createdb   bool    create db if it doesn't exist (alias -c)
     -help        bool    print help and exit
`

	if *help {
		fmt.Println(usage)
		os.Exit(1)
	}

	flag.Parse()

	if len(flag.Args()) != 4 {
		fmt.Println(usage)
		os.Exit(1)
	}

	createDbIfNotExists := (*createDbIfNotExists1 || *createDbIfNotExists2)

	dbPath := flag.Args()[0]
	lexName := flag.Args()[1]
	ssName := flag.Args()[2]
	locale := flag.Args()[3]

	_, dbFile := filepath.Split(dbPath)
	dbExt := filepath.Ext(dbFile)
	dbName := dbFile[0 : len(dbFile)-len(dbExt)]
	lexRefX := lex.NewLexRefWithInfo(dbName, lexName, ssName)
	closeAfter := true

	dbapi.Sqlite3WithRegex()
	err := createEmptyLexicon(dbPath, lexRefX, locale, createDbIfNotExists, closeAfter)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] %v", "createEmptyLexicon", err))
		os.Exit(1)
	}

}
