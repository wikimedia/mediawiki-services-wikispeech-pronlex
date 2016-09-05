package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/vrules"
)

func main() {

	sampleInvocation := `go run importLexToDB.go sv-se.nst sv-se.nst-SAMPA symbolset/static/sv-se_ws-sampa_maptable.csv pronlex.db swe030224NST.pron_utf8.txt`

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
	if err != nil {
		log.Fatal(err)
	}

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
		return
	}

	lexicon := dbapi.Lexicon{Name: lexName, SymbolSetName: symbolSetName}
	lexicon, err = dbapi.InsertLexicon(db, lexicon)

	if err != nil {
		log.Fatal(err)
	}

	logger := StderrLogger{}
	dbapi.ImportLexiconFile(dn, logger, lexName, inFile, symbolSet)
}
