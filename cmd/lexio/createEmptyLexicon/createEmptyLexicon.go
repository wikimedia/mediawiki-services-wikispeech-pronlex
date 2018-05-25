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

func main() {

	var createDbIfNotExists1 = flag.Bool("create_db", false, "create db if it doesn't exist")
	var createDbIfNotExists2 = flag.Bool("c", false, "create db if it doesn't exist")
	var help = flag.Bool("help", false, "print help and exit")

	usage := `Usage:
     $ createEmptyLexicon [flags] <DB PATH> <LEX NAME> <SYMBOL SET NAME> <LOCALE>

Flags:
     -create_db   bool    create db if it doesn't exist
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
	lexRef := lex.NewLexRef(dbName, lexName)

	dbapi.Sqlite3WithRegex()
	dbm := dbapi.NewDBManager()

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if createDbIfNotExists {
			err := dbm.DefineSqliteDB(lexRef.DBRef, dbPath)
			if err != nil {
				msg := fmt.Sprintf("[createEmptyLexicon] couldn't create db %s : %v", dbPath, err)
				fmt.Fprintln(os.Stderr, msg)
				os.Exit(1)
			}
			fmt.Fprintf(os.Stderr, "[createEmptyLexicon] created db %s\n", dbPath)
		} else {
			msg := fmt.Sprintf("[createEmptyLexicon] db does not exist %s : %v", dbPath, err)
			fmt.Fprintln(os.Stderr, msg)
			os.Exit(1)
		}
	}

	exists, err := dbm.LexiconExists(lexRef)
	if err != nil {
		msg := fmt.Sprintf("[createEmptyLexicon] couldn't lookup lexicon reference %s : %v", lexRef.String(), err)
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
	if exists {
		msg := fmt.Sprintf("[createEmptyLexicon] couldn't create lexicon that already exists: %s", lexRef.String())
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}

	err = dbm.DefineLexicon(lexRef, ssName, locale)
	if err != nil {
		msg := fmt.Sprintf("[createEmptyLexicon] couldn't create lexicon : %v", err)
		fmt.Fprintln(os.Stderr, msg)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon] created lexicon %s in %s\n", lexName, dbPath)
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon]  > symbolset: %s\n", ssName)
	fmt.Fprintf(os.Stderr, "[createEmptyLexicon]  > locale:    %s\n", locale)

}
