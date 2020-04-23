package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
)

func main() {

	var cmdName = "exportLex"

	var header = flag.Bool("header", false, "print header")

	var engineFlag = flag.String("db_engine", "sqlite", "db engine (sqlite or mariadb)")
	var dbLocation = flag.String("db_location", "", "db location (folder for sqlite; address for mariadb)")
	var dbName = flag.String("db_name", "", "db name (if empty, a list of available lexicons will be printed)")
	var lexName = flag.String("lex_name", "", "lexicon name")
	var outFile = flag.String("out_file", "", "Output file")

	var fatalError = false
	var dieIfEmptyFlag = func(name string, val *string) {
		if *val == "" {
			fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] flag %s is required", cmdName, name))
			fatalError = true
		}
	}
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "USAGE: exportLex [FLAGS]\n\n")
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
	err := dbm.OpenDB(*dbLocation, dbRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db : %v", err)
		os.Exit(1)
	}

	lexRefs, err := dbm.ListLexicons()
	lexNames := make(map[lex.LexName]bool)
	if err != nil {
		log.Fatalf("darn : %v", err)
	}
	for _, ref := range lexRefs {
		if *lexName == "" {
			fmt.Println(ref.LexRef.LexName)
		}
		lexNames[ref.LexRef.LexName] = true
	}
	if *lexName == "" {
		return
	}

	dieIfEmptyFlag("lex_name", lexName)
	dieIfEmptyFlag("out_file", lexName)
	if fatalError {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] exit from unrecoverable errors", cmdName))
		os.Exit(1)
	}

	lexRef := lex.NewLexRef(*dbName, *lexName)

	if _, ok := lexNames[lexRef.LexName]; !ok {
		log.Fatalf("no such lexicon name '%s'", *lexName)
		return
	}
	q := dbapi.DBMQuery{LexRefs: []lex.LexRef{lexRef}, Query: dbapi.Query{WordLike: "%"}}
	f, err := os.Create(*outFile)
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
