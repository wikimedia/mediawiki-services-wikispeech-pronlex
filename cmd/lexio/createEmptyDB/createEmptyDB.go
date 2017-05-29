// createEmptyDB initialises an Sqlite3 relational database from the schema defining a lexicon database, but empty of data.
// See dbapi.Schema.
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stts-se/pronlex/dbapi"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "createEmptyDB <OUTPUT FILE NAME>")
		os.Exit(1)
	}

	fOut := os.Args[1]
	if _, err := os.Stat(fOut); !os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Cannot create file that already exists:", fOut)
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", fOut)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(dbapi.Schema)
	if err != nil {
		log.Fatal(err)
	}

}
