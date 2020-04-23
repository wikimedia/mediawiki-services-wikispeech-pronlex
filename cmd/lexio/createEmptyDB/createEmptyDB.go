// createEmptyDB initialises an Sqlite3 relational database from the schema defining a lexicon database, but empty of data.
// See dbapi.Schema.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

func main() {
	var cmdName = "exportLex"

	var engineFlag = flag.String("db_engine", "sqlite", "db engine (sqlite or mariadb)")
	var dbLocation = flag.String("db_location", "", "db location (folder for sqlite; address for mariadb)")
	var dbName = flag.String("db_name", "", "db name (if empty, a list of available lexicons will be printed)")

	var fatalError = false
	var dieIfEmptyFlag = func(name string, val *string) {
		if *val == "" {
			fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] flag %s is required", cmdName, name))
			fatalError = true
		}
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: createEmptyDB [FLAGS]\n\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) != 0 {
		flag.Usage()
		os.Exit(1)
	}

	dieIfEmptyFlag("db_engine", engineFlag)
	dieIfEmptyFlag("db_location", dbLocation)
	dieIfEmptyFlag("db_name", dbName)
	if fatalError {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] exit from unrecoverable errors", cmdName))
		os.Exit(1)
	}

	dbapi.Sqlite3WithRegex()
	var dbm *dbapi.DBManager
	if *engineFlag == "mariadb" {
		dbm = dbapi.NewMariaDBManager()
	} else if *engineFlag == "sqlite" {
		dbm = dbapi.NewSqliteDBManager()
	} else {
		fmt.Fprintf(os.Stderr, "invalid db engine : %s\n", *engineFlag)
		os.Exit(1)
	}
	dbRef := lex.DBRef(*dbName)

	dbExists, err := dbm.DBExists(*dbLocation, dbRef)
	if err != nil {
		log.Fatal(err)
	}

	if dbExists {
		fmt.Fprintln(os.Stderr, "[createEmptyDB] Cannot create a db that already exists:", *dbName)
		os.Exit(1)
	}

	err = dbm.DefineDB(*dbLocation, dbRef)
	if err != nil {
		log.Fatalf("Couldn't create db: %v", err)
	}
	log.Printf("Created database %s", *dbName)
}
