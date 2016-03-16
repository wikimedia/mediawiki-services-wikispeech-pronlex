package main

import (
	"bufio"
	"database/sql"
	"github.com/stts-se/pronlex/dbapi"
	"log"
	"os"
)

func main() {

	// TODO help message
	// TODO command line arg processing

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatalf("darn : %v", err)
	}

	// TODO read lexicon name from command line
	l := dbapi.Lexicon{ID: 1}

	f, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalf("aouch : %v", err)
	}

	f.Sync() // ??
	bf := bufio.NewWriter(f)
	defer bf.Flush()

	dbapi.Export(db, l, bf)

}
