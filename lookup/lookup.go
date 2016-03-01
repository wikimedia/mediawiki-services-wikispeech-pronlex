package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"morf.se/wsgo/pronlex/dbapi"
	"os"
)

// TODO replace calls to ff() with proper error handling
// ff is a place holder to easily find places lacking sane error handling 
func ff(f string, err error) {
	if err != nil {
		log.Fatalf(f, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Println("<SQLITE DB FILE> <WORDS ...>")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	l, err := dbapi.GetLexicon(db, "sv.se.nst")
	ff("Failed to get lexicon from db: %v", err)
	ls := []dbapi.Lexicon{l}

	q := dbapi.Query{Lexicons: ls,
		Words:      os.Args[2:],
		PageLength: 100}

	//log.Printf("VOFF: %v", q)

	res := dbapi.GetEntries(db, q)
	if err != nil {
		log.Fatal(err)
	}
	//for _, e := range res {
	res0, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", res0)
	//}
}
