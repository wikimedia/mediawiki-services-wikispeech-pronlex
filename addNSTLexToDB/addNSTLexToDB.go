package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/line"
)

func loadLineFmt() (line.Format, error) {
	tests := []line.FormatTest{
		line.FormatTest{"storstaden;NN;SIN|DEF|NOM|UTR;stor+staden;JJ+NN;LEX|INFL;SWE;;;;;\"\"stu:$%s`t`A:$den;1;STD;SWE;;;;;;;;;;;;;;18174;enter_se|inflector;;INFLECTED;storstad|95522;s111n, a->ä, stad;s111;;;;;;;;;;;;;storstaden;;;88748",
			map[line.Field]string{
				line.Orth:           "storstaden",
				line.Pos:            "NN",
				line.Morph:          "SIN|DEF|NOM|UTR",
				line.Decomp:         "stor+staden",
				line.WordLang:       "SWE",
				line.Trans1:         "\"\"stu:$%s`t`A:$den",
				line.Translang1:     "SWE",
				line.Trans2:         "",
				line.Translang2:     "",
				line.Trans3:         "",
				line.Translang3:     "",
				line.Trans4:         "",
				line.Translang4:     "",
				line.Lemma:          "storstad|95522",
				line.InflectionRule: "s111n, a->ä, stad",
			},
			"storstaden;NN;SIN|DEF|NOM|UTR;stor+staden;;;SWE;;;;;\"\"stu:$%s`t`A:$den;;;SWE;;;;;;;;;;;;;;;;;;storstad|95522;s111n, a->ä, stad;;;;;;;;;;;;;;;;;",
		},
	}
	f, err := line.NewFormat(
		"NST",
		";",
		map[line.Field]int{
			line.Orth:           0,
			line.Pos:            1,
			line.Morph:          2,
			line.Decomp:         3,
			line.WordLang:       6,
			line.Trans1:         11,
			line.Translang1:     14,
			line.Trans2:         15,
			line.Translang2:     18,
			line.Trans3:         19,
			line.Translang3:     22,
			line.Trans4:         23,
			line.Translang4:     26,
			line.Lemma:          32,
			line.InflectionRule: 33,
		},
		51,
		tests,
	)
	if err != nil {
		return f, err
	}
	return f, nil
}

func appendTrans(ts []dbapi.Transcription, t string, l string) []dbapi.Transcription {
	if "" == strings.TrimSpace(t) {
		return ts
	}
	ts = append(ts, dbapi.Transcription{Strn: t, Language: l})
	return ts
}

func getTranses(fs map[line.Field]string) []dbapi.Transcription {
	res := make([]dbapi.Transcription, 0)
	res = appendTrans(res, fs[line.Trans1], fs[line.Translang1])
	res = appendTrans(res, fs[line.Trans2], fs[line.Translang2])
	res = appendTrans(res, fs[line.Trans3], fs[line.Translang3])
	res = appendTrans(res, fs[line.Trans4], fs[line.Translang4])
	return res
}

func nstLine2Entry(lineFmt line.Format, l string) (dbapi.Entry, error) {
	fs, err := lineFmt.Parse(l)
	if err != nil {
		return dbapi.Entry{}, err
	}

	res := dbapi.Entry{
		Strn:           strings.ToLower(fs[line.Orth]),
		Language:       fs[line.WordLang],
		PartOfSpeech:   fs[line.Pos] + " " + fs[line.Morph],
		WordParts:      fs[line.Decomp],
		Transcriptions: getTranses(fs),
	}

	lemmaReading := strings.SplitN(fs[line.Lemma], "|", 2)
	lemma := ""
	reading := ""
	if len(lemmaReading) == 2 {
		lemma = lemmaReading[0]
		reading = lemmaReading[1]
	}
	if len(lemmaReading) == 1 {
		lemma = lemmaReading[0]
	}
	paradigm := fs[line.InflectionRule]
	lemmaStruct := dbapi.Lemma{Strn: lemma, Reading: reading, Paradigm: paradigm}

	if "" != lemmaStruct.Strn {
		res.Lemma = lemmaStruct
	}

	return res, nil
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

	lineFmt, err := loadLineFmt()
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
		e, err := nstLine2Entry(lineFmt, l)
		if err != nil {
			log.Fatal(err)
		}
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
