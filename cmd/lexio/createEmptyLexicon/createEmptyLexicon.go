// createEmptyLexicon initialises a database from the schema defining a lexicon database, but empty of data.
package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

func createEmptyLexicon(engine dbapi.DBEngine, dbLocation string, dbRef lex.DBRef, lexRefX lex.LexRefWithInfo, locale string, createDbIfNotExists bool, closeAfter bool) error {

	//fmt.Printf("loc=%s, ref=%s, lexRefX=%s\n", dbLocation, dbRef, lexRefX)
	dbm, err := dbapi.NewDBManager(engine)
	if err != nil {
		return fmt.Errorf("couldn't create db manager : %v", err)
	}

	lexRef := lexRefX.LexRef // db + lexname without ssName

	if closeAfter {
		defer dbm.CloseDB(dbRef)
	}

	dbExists, err := dbm.DBExists(dbLocation, dbRef)
	if err != nil {
		return fmt.Errorf("couldn't check if db exists : %v", err)
	}

	if !dbExists {
		if createDbIfNotExists {
			err := dbm.DefineDB(dbLocation, dbRef)
			if err != nil {
				return fmt.Errorf("couldn't create db %s : %v", dbRef, err)
			}
			//log.Printf("Created db %s\n", dbRef)
		} else {
			return fmt.Errorf("db does not exist: %s", dbRef)
		}
	} else {
		err := dbm.OpenDB(dbLocation, dbRef)
		if err != nil {
			return fmt.Errorf("couldn't open db %s : %v", dbRef, err)
		}
		//log.Printf("Opened db %s\n", dbRef)
	}

	if !dbm.ContainsDB(dbRef) {
		return fmt.Errorf("db should be registered in dbm, but wasn't : %s", dbRef)
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
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon] created lexicon %s in %s\n", lexRefX.LexRef.LexName, dbRef)
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon]  > symbolset: %s\n", lexRefX.SymbolSetName)
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon]  > locale:    %s\n", locale)

	return nil
}

func main() {
	cmdName := "createEmptyLexicon"

	var createDbIfNotExists = flag.Bool("createdb", false, "create db if it doesn't exist")

	var engineFlag = flag.String("db_engine", "sqlite", "db engine (sqlite or mariadb)")
	var dbLocation = flag.String("db_location", "", "db location (folder for sqlite; address for mariadb)")
	var dbName = flag.String("db_name", "", "db name")
	var lexName = flag.String("lexicon", "", "lexicon name")
	var locale = flag.String("locale", "", "lexicon locale")
	var ssName = flag.String("symbolset", "", "lexicon symbolset")

	var help = flag.Bool("help", false, "print help and exit")

	var fatalError = false
	var dieIfEmptyFlag = func(name string, val *string) {
		if *val == "" {
			fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] flag %s is required", cmdName, name))
			fatalError = true
		}
	}

	var printUsage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
     $ createEmptyLexicon [flags]

Flags: 
`)
		flag.PrintDefaults()

	}

	flag.Usage = printUsage

	flag.Parse()

	if *help {
		printUsage()
		os.Exit(1)
	}

	if len(flag.Args()) != 0 {
		printUsage()
		os.Exit(1)
	}

	dieIfEmptyFlag("db_engine", engineFlag)
	dieIfEmptyFlag("db_location", dbLocation)
	dieIfEmptyFlag("db_name", dbName)
	dieIfEmptyFlag("lexicon", lexName)
	dieIfEmptyFlag("locale", locale)
	dieIfEmptyFlag("symbolset", ssName)
	if fatalError {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] exit from unrecoverable errors", cmdName))
		os.Exit(1)
	}

	var dbEngine dbapi.DBEngine
	if *engineFlag == "sqlite" {
		dbEngine = dbapi.Sqlite
	} else if *engineFlag == "mariadb" {
		dbEngine = dbapi.MariaDB
	} else {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] %v", cmdName, "invalid db engine"))
		os.Exit(1)
	}

	lexRefX := lex.NewLexRefWithInfo(*dbName, *lexName, *ssName)
	closeAfter := true

	dbapi.Sqlite3WithRegex()
	err := createEmptyLexicon(dbEngine, *dbLocation, lexRefX.LexRef.DBRef, lexRefX, *locale, *createDbIfNotExists, closeAfter)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] %v", cmdName, err))
		os.Exit(1)
	}

}
