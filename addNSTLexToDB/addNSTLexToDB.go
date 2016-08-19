package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	//	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/vrules"
)

func main() {

	sampleInvocation := `go run addNSTLexToDB.go sv.se.nst sv.se.nst-SAMPA symbolset/static/sv_ws-sampa_maptable.csv pronlex.db swe030224NST.pron_utf8.txt`

	if len(os.Args) != 6 {
		log.Fatal("Expected <DB LEXICON NAME> <SYMBOLSET NAME> <SYMBOLSET FILE> <DB FILE> <NST INPUT FILE>", "\n\tSample invocation: ", sampleInvocation)
	}

	lexName := os.Args[1]
	symbolSetName := os.Args[2]
	ssFileName := os.Args[3]
	dbFile := os.Args[4]
	inFile := os.Args[5]

	_, err := os.Stat(dbFile)
	if err != nil {
		log.Fatalf("Cannot find db file. %v", err)
	}

	ssMapper, err := symbolset.LoadMapper(symbolSetName, ssFileName, "SYMBOL", "IPA")
	ssRule := vrules.SymbolSetRule{ssMapper.From}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	_, err = dbapi.GetLexicon(db, lexName)
	if err == nil {
		log.Fatalf("Nothing will be added. Lexicon already exists in database: %s", lexName)
	}

	lexicon := dbapi.Lexicon{Name: lexName, SymbolSetName: symbolSetName}
	lexicon, err = dbapi.InsertLexicon(db, lexicon)
	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.Open(inFile)
	defer fh.Close()
	if err != nil {
		log.Fatal(err)
	}

	nstFmt, err := line.NewNST()
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(fh)
	n := 0
	var eBuf []lex.Entry
	for s.Scan() {
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
		l := s.Text()
		e, err := nstFmt.ParseToEntry(l)
		if err != nil {
			log.Fatal(err)
		}

		// if no space in transcription, add these using symbolset.SplitIntoPhonemes utility

		for _, r := range ssRule.Validate(e) {
			panic(r) // shouldn't happen
		}

		e.EntryStatus = lex.EntryStatus{Name: "imported", Source: "nst"}
		eBuf = append(eBuf, e)
		n++
		if n%10000 == 0 {
			_, err = dbapi.InsertEntries(db, lexicon, eBuf)
			if err != nil {
				log.Fatal(err)
			}
			eBuf = make([]lex.Entry, 0)
			fmt.Printf("\rLines read: %d               \r", n)
		}
	}
	dbapi.InsertEntries(db, lexicon, eBuf) // flushing the buffer

	_, err = db.Exec("ANALYZE")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Lines read:\t%d", n)
}
