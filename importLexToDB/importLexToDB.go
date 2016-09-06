package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/symbolset"

	"fmt"
)

func main() {

	sampleInvocation := `go run importLexToDB.go pronlex.db sv-se.nst [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8 sv-se_ws-sampa [SYMBOL SET FOLDER]/sv-se_ws-sampa.csv`

	if len(os.Args) != 6 {
		log.Fatal("Expected <DB FILE> <LEXICON NAME> <LEXICON FILE> <SYMBOLSET NAME> <SYMBOLSET FILE>", "\n\tSample invocation: ", sampleInvocation)
	}

	dbFile := os.Args[1]
	lexName := os.Args[2]
	inFile := os.Args[3]
	symbolSetName := os.Args[4]
	ssFileName := os.Args[5]

	_, err := os.Stat(dbFile)
	if err != nil {
		log.Fatalf("Cannot find db file. %v", err)
	}

	ssMapper, err := symbolset.LoadMapper(symbolSetName, ssFileName, "SYMBOL", "IPA")
	if err != nil {
		log.Fatal(err)
	}
	symbolSet := ssMapper.From

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

	logger := dbapi.StderrLogger{}
	// TODO handle errors! Does it make sent to return array of error...?
	var errs []error
	errs = dbapi.ImportLexiconFile(db, logger, lexName, inFile, symbolSet)

	if len(errs) == 0 {
		logger.Write("running the Sqlite3 ANALYZE command. It may take a little while...")
		_, err = db.Exec("ANALYZE")
		if err != nil {
			logger.Write(fmt.Sprintf("failed to run ANALYZE command : %v", err))
			return
		}
	}

	// Loop over the symbols of the symbolset file given as a command line argument.
	// For each such symbol, convert it to a dbapi.Symbol, and finally add all symbols to the db in one go.
	var dbSymSet []dbapi.Symbol
	for _, sym := range symbolSet.Symbols {
		s := sym.String
		cat := sym.Cat.String()
		desc := sym.Desc
		//TODO Add function to obtain map for single symbol...?
		ipa, err := ssMapper.MapTranscription(s)
		if err != nil {
			logger.Write(fmt.Sprintf("failed to obtain IPA carachter for '%v' : %v", s, err))
			//return
		}
		dbSym := dbapi.Symbol{LexiconID: lexicon.ID, Symbol: s, Category: cat, Description: desc, IPA: ipa}

		dbSymSet = append(dbSymSet, dbSym)
	}

	err = dbapi.SaveSymbolSet(db, dbSymSet)
	if err != nil {
		msg := fmt.Sprintf("dbapi.SaveSymbolSet returned error : %v", err)
		logger.Write(msg)
		//log.Println(msg)
	}
}
