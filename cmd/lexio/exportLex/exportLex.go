package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
)

func main() {

	// TODO read lexicon name from command line
	// + db look up

	if len(os.Args) != 4 && len(os.Args) != 2 {
		log.Println("exportLex <DB_FILE> <LEXICON_NAME> <OUTPUT_FILE_NAME>")
		log.Println(" if only <DB_FILE> is specified, a list of available lexicons will be printed")
		return
	}

	dbFile := os.Args[1]
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("darn : %v", err)
	}

	dbm := dbapi.NewDBManager()
	dbRef := lex.DBRef(dbFile)
	dbm.AddDB(dbRef, db)

	lexRefs, err := dbm.ListLexicons()
	lexNames := make(map[lex.LexName]bool)
	if err != nil {
		log.Fatalf("darn : %v", err)
	}
	for _, ref := range lexRefs {
		if len(os.Args) == 2 {
			fmt.Println(ref.LexRef.LexName)
		}
		lexNames[ref.LexRef.LexName] = true
	}
	if len(os.Args) == 2 {
		return
	}

	lexName := os.Args[2]
	lexRef := lex.NewLexRef(dbFile, lexName)

	if "" == lexName {
		log.Fatalf("invalid lexicon name '%s'", lexName)
		return
	}
	if _, ok := lexNames[lexRef.LexName]; !ok {
		log.Fatalf("no such lexicon name '%s'", lexName)
		return
	}
	q := dbapi.DBMQuery{LexRefs: []lex.LexRef{lexRef}}
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
	dbm.LookUp(q, writer)
}
