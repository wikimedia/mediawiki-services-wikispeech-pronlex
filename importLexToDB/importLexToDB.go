package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/pronlex/vrules"

	"fmt"
)

var vServ = vrules.ValidatorService{Validators: make(map[string]*validation.Validator)}

func loadValidators(symsetDirName string) error {
	symbolSets, err := symbolset.LoadSymbolSetsFromDir(symsetDirName)
	if err != nil {
		return err
	}
	err = vServ.Load(symbolSets)
	return err
}

func main() {

	sampleInvocation := `go run importLexToDB.go pronlex.db sv-se.nst [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8 sv-se_ws-sampa [SYMBOLSET FOLDER]`

	if len(os.Args) != 5 && len(os.Args) != 6 {
		log.Fatal("Expected <DB FILE> <LEXICON NAME> <LEXICON FILE> <SYMBOLSET NAME> <SYMBOLSET FOLDER> (optional)", "\n\t if <SYMBOLSET FOLDER> is specified, all entries will be validated upon import, and the validation result will be d in the database\n\tSample invocation: ", sampleInvocation)
	}

	dbFile := os.Args[1]
	lexName := os.Args[2]
	inFile := os.Args[3]
	symbolSetName := os.Args[4]

	validator := &validation.Validator{}
	if len(os.Args) == 6 {
		symsetDirName := os.Args[5]

		err := loadValidators(symsetDirName)
		if err != nil {
			msg := fmt.Sprintf("failed to load validators : %v", err)
			log.Fatal(msg)
			return
		}
		vdat, err := vServ.ValidatorForName(symbolSetName)
		validator = vdat
		if err != nil {
			msg := fmt.Sprintf("failed to get validator for symbol set %v : %v", symbolSetName, err)
			log.Fatal(msg)
			return
		}
		log.Println("Validator created for " + validator.Name)
	}

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
	err = dbapi.ImportLexiconFile(db, logger, lexName, inFile, validator)

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

}
