package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
)

func main() {

	var usage = "USAGE: exportLex [-header] <DB_FILE> <LEXICON_NAME> <OUTPUT_FILE_NAME>\n" +
		" if only <DB_FILE> is specified, a list of available lexicons will be printed\n" +
		" optional flag: header (print header in output file)\n"

	var header = flag.Bool("header", false, "print header")
	flag.Usage = func() {
		fmt.Println(strings.TrimSpace(usage))
	}
	flag.Parse()

	var args = flag.Args()

	if len(args) != 3 && len(args) != 1 {
		fmt.Fprint(os.Stderr, usage)
		return
	}

	dbFile := args[0]
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("darn : %v", err)
	}

	dbm := dbapi.NewDBManager()
	dbRef := lex.DBRef(dbFile)
	err = dbm.AddDB(dbRef, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to add db to db manager : %v", err)
		os.Exit(1)
	}

	lexRefs, err := dbm.ListLexicons()
	lexNames := make(map[lex.LexName]bool)
	if err != nil {
		log.Fatalf("darn : %v", err)
	}
	for _, ref := range lexRefs {
		if len(args) == 1 {
			fmt.Println(ref.LexRef.LexName)
		}
		lexNames[ref.LexRef.LexName] = true
	}
	if len(args) == 1 {
		return
	}

	lexName := args[1]
	lexRef := lex.NewLexRef(dbFile, lexName)

	if lexName == "" {
		log.Fatalf("invalid lexicon name '%s'", lexName)
		return
	}
	if _, ok := lexNames[lexRef.LexName]; !ok {
		log.Fatalf("no such lexicon name '%s'", lexName)
		return
	}
	q := dbapi.DBMQuery{LexRefs: []lex.LexRef{lexRef}, Query: dbapi.Query{WordLike: "%"}}
	f, err := os.Create(args[2])
	if err != nil {
		log.Fatalf("aouch : %v", err)
	}

	bf := bufio.NewWriter(f)
	defer bf.Flush()

	wsFmt, err := line.NewWS()
	if err != nil {
		log.Fatal(err)
	}
	if *header {
		_, err := bf.Write([]byte(fmt.Sprintf("#%s\n", wsFmt.Header())))
		if err != nil {
			fmt.Fprintf(os.Stderr, "write error : %v\n", err)
			os.Exit(1)
		}
	}

	writer := line.FileWriter{Parser: wsFmt, Writer: bf}
	err = dbm.LookUp(q, writer)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to do lexicon lookup : %v\n", err)
	}
}
