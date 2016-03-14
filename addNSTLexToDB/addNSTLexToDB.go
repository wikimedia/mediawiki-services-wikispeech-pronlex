package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/stts-se/pronlex/dbapi"
	"log"
	"os"
	"strings"
)

const (
	orth     = 0
	pos      = 1
	morph    = 2
	decomp   = 3
	wordLang = 6
	//aCrAbbr        = 9
	trans1         = 11
	translang1     = 14
	trans2         = 15
	translang2     = 18
	trans3         = 19
	translang3     = 22
	trans4         = 23
	translang4     = 26
	lemma          = 32
	inflectionRule = 33
	//morphLabel     = 34
)

func appendTrans(ts []dbapi.Transcription, t string, l string) []dbapi.Transcription {
	if "" == strings.TrimSpace(t) {
		return ts
	}
	ts = append(ts, dbapi.Transcription{Strn: t, Language: l})
	return ts
}

func getTranses(fs []string) []dbapi.Transcription {
	t1, l1, t2, l2, t3, l3, t4, l4 := fs[trans1], fs[translang1], fs[trans2], fs[translang2], fs[trans3], fs[translang3], fs[trans4], fs[translang4]

	res := make([]dbapi.Transcription, 0)
	res = appendTrans(res, t1, l1)
	res = appendTrans(res, t2, l2)
	res = appendTrans(res, t3, l3)
	res = appendTrans(res, t4, l4)

	return res
}

func nstLine2Entry(l string) dbapi.Entry {

	fs := strings.Split(l, ";")

	res := dbapi.Entry{
		Strn:           strings.ToLower(fs[orth]),
		Language:       fs[wordLang],
		PartOfSpeech:   fs[pos] + " " + fs[morph],
		WordParts:      fs[decomp],
		Transcriptions: getTranses(fs),
	}

	lemmaReading := strings.SplitN(fs[lemma], "|", 2)
	lemma := ""
	reading := ""
	if len(lemmaReading) == 2 {
		lemma = lemmaReading[0]
		reading = lemmaReading[1]
	}
	if len(lemmaReading) == 1 {
		lemma = lemmaReading[0]
	}
	paradigm := fs[inflectionRule]
	lemmaStruct := dbapi.Lemma{Strn: lemma, Reading: reading, Paradigm: paradigm}

	if "" != lemmaStruct.Strn {
		res.Lemma = lemmaStruct
	}

	return res
}

func main() {

	sampleInvocation := `go run addNSTLexToDB.go sv.se.nst pronlex.db swe030224NST.pron_utf8.txt`

	if len(os.Args) != 4 {
		log.Fatal("Expected <DB LEXICON NAME> <DB FILE> <NST INPUT FILE>", "\n\tSample invocation: ", sampleInvocation)
	}

	lexName := os.Args[1]
	dbFile := os.Args[2]
	inFile := os.Args[3]

	_, err := os.Stat(dbFile)
	if err != nil {
		log.Fatalf("Cannot find db file. %v", err)
	}

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

	// TODO hard coded symbol set name
	lex := dbapi.Lexicon{Name: lexName, SymbolSetName: "nst-sv-SAMPA"}
	lex, err = dbapi.InsertLexicon(db, lex)
	if err != nil {
		log.Fatal(err)
	}

	fh, err := os.Open(inFile)
	defer fh.Close()
	if err != nil {
		log.Fatal(err)
	}

	s := bufio.NewScanner(fh)
	n := 0
	eBuf := make([]dbapi.Entry, 0)
	for s.Scan() {
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
		l := s.Text()
		e := nstLine2Entry(l)
		eBuf = append(eBuf, e)
		n++
		if n%10000 == 0 {
			_, err = dbapi.InsertEntries(db, lex, eBuf)
			if err != nil {
				log.Fatal(err)
			}
			eBuf = make([]dbapi.Entry, 0)
			fmt.Printf("\rLines read: %d               \r", n)
		}
	}
	dbapi.InsertEntries(db, lex, eBuf) // flushing the buffer

	log.Printf("Lines read:\t%d", n)
}
