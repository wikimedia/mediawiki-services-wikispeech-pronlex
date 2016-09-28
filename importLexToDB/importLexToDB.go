package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/stts-se/pronlex/dbapi"

	"fmt"
)

// TODO: Function to be called from server
//func ImportLexicon(dbapi dbapi, lexName string, inFile string, symbolSetName string, ssFileName string) error {
//}

func main() {

	sampleInvocation := `go run importLexToDB.go pronlex.db sv-se.nst [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8 sv-se_ws-sampa`

	if len(os.Args) != 5 {
		log.Fatal("Expected <DB FILE> <LEXICON NAME> <LEXICON FILE> <SYMBOLSET NAME>", "\n\tSample invocation: ", sampleInvocation)
	}

	dbFile := os.Args[1]
	lexName := os.Args[2]
	inFile := os.Args[3]
	symbolSetName := os.Args[4]

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
		return
	}

	lexicon := dbapi.Lexicon{Name: lexName, SymbolSetName: symbolSetName}
	lexicon, err = dbapi.InsertLexicon(db, lexicon)

	if err != nil {
		log.Fatal(err)
	}

	logger := dbapi.StderrLogger{}
	// TODO handle errors! Does it make sent to return array of error...?
	err = dbapi.ImportLexiconFile(db, logger, lexName, inFile)

	if err != nil {
		log.Fatal(err)
		return
	}

	logger.Write("running the Sqlite3 ANALYZE command. It may take a little while...")
	_, err = db.Exec("ANALYZE")
	if err != nil {
		logger.Write(fmt.Sprintf("failed to run ANALYZE command : %v", err))
		return
	}

	logger.Write("finished importing lexicon file")

	// // Loop over the symbols of the symbolset file given as a command line argument.
	// // For each such symbol, convert it to a dbapi.Symbol, and finally add all symbols to the db in one go.
	// var dbSymSet []dbapi.Symbol
	// for _, sym := range symbolSet.Symbols {
	// 	s := sym.String
	// 	cat := sym.Cat.String()
	// 	desc := sym.Desc
	// 	ipa, err := ss.MapSymbol(sym)
	// 	if err != nil {
	// 		logger.Write(fmt.Sprintf("failed to obtain IPA character for '%v' : %v", s, err))
	// 	}
	// 	dbSym := dbapi.Symbol{LexiconID: lexicon.ID, Symbol: s, Category: cat, Description: desc, IPA: ipa.String}

	// 	dbSymSet = append(dbSymSet, dbSym)
	// }

	// err = dbapi.SaveSymbolSet(db, dbSymSet)
	// if err != nil {
	// 	msg := fmt.Sprintf("dbapi.SaveSymbolSet returned error : %v", err)
	// 	logger.Write(msg)
	// 	//log.Println(msg)
	// }
	// logger.Write("finished loading symbol set")
}
