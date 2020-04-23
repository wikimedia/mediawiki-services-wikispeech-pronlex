// createEmptyLexicon initialises a database from the schema defining a lexicon database, but empty of data.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"

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
		defer dbm.CloseDB(lexRef.DBRef)
	}

	sqliteDBPath := path.Join(dbLocation, string(dbRef)+".db")
	if _, err := os.Stat(sqliteDBPath); engine == dbapi.Sqlite && os.IsNotExist(err) {
		if createDbIfNotExists {
			err := dbm.DefineDB(dbLocation, lexRef.DBRef)
			if err != nil {
				return fmt.Errorf("couldn't create db %s : %v", sqliteDBPath, err)
			}
			fmt.Fprintf(os.Stderr, "[createEmptyLexicon] created db %s\n", sqliteDBPath)
		} else {
			return fmt.Errorf("db does not exist %s : %v", sqliteDBPath, err)
		}
	} else {
		err := dbm.OpenDB(dbLocation, lexRef.DBRef)
		if err != nil {
			return fmt.Errorf("couldn't open db %s : %v", sqliteDBPath, err)
		}
		fmt.Fprintf(os.Stderr, "[createEmptyLexicon] opened db %s\n", sqliteDBPath)
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
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon] created lexicon %s in %s\n", lexRefX.LexRef.LexName, dbRef)
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon]  > symbolset: %s\n", lexRefX.SymbolSetName)
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon]  > locale:    %s\n", locale)

	return nil
}

func main() {

	var createDbIfNotExists = flag.Bool("createdb", false, "create db if it doesn't exist")
	var dbEngine = flag.String("db_engine", "sqlite", "db engine (sqlite or mariadb)")
	var dbLocation = flag.String("db_location", "", "db location (folder for sqlite; address for mariadb)")
	var help = flag.Bool("help", false, "print help and exit")

	var printUsage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
     $ createEmptyLexicon [flags] <DB NAME> <LEX NAME> <SYMBOL SET NAME> <LOCALE>

Flags: 
 `)
		flag.PrintDefaults()

	}

	if *help {
		printUsage()
		os.Exit(1)
	}

	flag.Parse()

	if len(flag.Args()) != 4 {
		printUsage()
		os.Exit(1)
	}

	dbRef := lex.DBRef(flag.Args()[0])
	lexName := flag.Args()[1]
	ssName := flag.Args()[2]
	locale := flag.Args()[3]

	var engine dbapi.DBEngine
	if *dbEngine == "sqlite" {
		engine = dbapi.Sqlite
	} else if *dbEngine == "mariadb" {
		engine = dbapi.MariaDB
	} else {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] %v", "createEmptyLexicon", "invalid db engine"))
		os.Exit(1)
	}

	lexRefX := lex.NewLexRefWithInfo(string(dbRef), lexName, ssName)
	closeAfter := true

	dbapi.Sqlite3WithRegex()
	err := createEmptyLexicon(engine, *dbLocation, dbRef, lexRefX, locale, *createDbIfNotExists, closeAfter)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] %v", "createEmptyLexicon", err))
		os.Exit(1)
	}

}
