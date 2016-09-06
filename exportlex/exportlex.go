package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/line"
)

func main() {

	// TODO read lexicon name from command line
	// + db look up

	if len(os.Args) != 4 {
		log.Println("exportlex <DB_FILE> <LEXICON_DB_ID> <OUTPUT_FILE_NAME>")
		return
	}

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatalf("darn : %v", err)
	}

	dbIDstr := os.Args[2]
	dbID, err := strconv.ParseInt(dbIDstr, 10, 64)
	if err != nil {
		log.Fatalf("failed to convert command line option %s into int : %v", dbIDstr, err)
		return
	}
	ls := []dbapi.Lexicon{dbapi.Lexicon{ID: dbID}}
	q := dbapi.Query{Lexicons: ls}
	f, err := os.Create(os.Args[3])
	if err != nil {
		log.Fatalf("aouch : %v", err)
	}

	bf := bufio.NewWriter(f)
	defer bf.Flush()

	wsFmt, err := line.NewWS()
	if err != nil {
		log.Fatal(err)
	}
	writer := line.FileWriter{wsFmt, bf}
	dbapi.LookUp(db, q, writer)

}
