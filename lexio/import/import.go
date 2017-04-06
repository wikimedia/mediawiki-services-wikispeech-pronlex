package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/pronlex/validation/validators"

	"fmt"
)

var vServ = validators.ValidatorService{Validators: make(map[string]*validation.Validator)}

func loadValidators(symsetDirName string) error {
	symbolSets, err := symbolset.LoadSymbolSetsFromDir(symsetDirName)
	if err != nil {
		return err
	}
	err = vServ.Load(symbolSets)
	return err
}

func main() {

	var f = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var validate = f.Bool("validate", false, "validate each entry, and save the validation in the database (default: false)")
	var force = f.Bool("force", false, "force loading of lexicon even if the symbolset is undefined (default: false)")
	var help = f.Bool("help", false, "print help message")

	usage := `USAGE:
  go run import.go <FLAGS> <DB FILE> <LEXICON NAME> <LEXICON FILE> <SYMBOLSET NAME> <SYMBOLSET FOLDER>

FLAGS:
   -validate bool  validate each entry when loading, and save validation in the database (default: false)
   -force    bool  force loading of lexicon even if the symbolset is undefined (default: false)
   -help     bool  print help message

SAMPLE INVOCATION:
  go run import.go -validate pronlex.db sv-se.nst [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8 sv-se_ws-sampa [SYMBOLSET FOLDER]`

	f.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}

	var args = os.Args
	if strings.HasSuffix(args[0], "import") {
		args = args[1:] // remove first argument if it's the program name
	}
	err := f.Parse(args)
	if err != nil {
		os.Exit(1)
	}

	args = f.Args()

	if *help {
		fmt.Println(usage)
		os.Exit(1)
	}

	if len(args) != 5 {
		fmt.Println(usage)
		os.Exit(1)
	}

	dbFile := args[0]
	lexName := args[1]
	inFile := args[2]
	symbolSetName := args[3]

	validator := &validation.Validator{}
	symsetDirName := args[4]
	if *validate {

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
	} else if !*force { // do not validate but still check symbolset
		_, err := symbolset.LoadSymbolSet(filepath.Join(symsetDirName, symbolSetName+symbolset.SymbolSetSuffix))
		if err != nil {
			msg := fmt.Sprintf("failed to load symbol set %v : %v", symbolSetName, err)
			log.Fatal(msg)
			return
		}
	}

	_, err = os.Stat(dbFile)
	if err != nil {
		log.Fatalf("Cannot find db file: %v", err)
	}

	_, err = os.Stat(inFile)
	if err != nil {
		log.Fatalf("Cannot find lexicon file: %v", err)
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
