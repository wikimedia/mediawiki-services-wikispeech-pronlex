package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
)

func main() {

	// TODO read lexicon name from command line
	// + db look up

	if len(os.Args) != 4 {
		log.Println("go run export.go <DB_FILE> <LEXICON_NAME> <OUTPUT_FILE_NAME>")
		return
	}

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatalf("darn : %v", err)
	}

	lexName := os.Args[2]
	if "" == lexName {
		log.Fatalf("invalid lexicon name '%s'", lexName)
		return
	}
	ls := []lex.LexName{lex.LexName(lexName)}
	q := dbapi.Query{}
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
	writer := line.FileWriter{Parser: wsFmt, Writer: bf}
	dbapi.LookUp(db, ls, q, writer)

}
