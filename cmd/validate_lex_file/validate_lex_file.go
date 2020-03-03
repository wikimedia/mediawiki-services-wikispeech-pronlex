package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/pronlex/validation/validators"
	"github.com/stts-se/symbolset"
)

var vServ = validators.ValidatorService{Validators: make(map[string]*validation.Validator)}

func loadValidators(symsetDirName string) error {
	symbolSets, err := symbolset.LoadSymbolSetsFromDir(symsetDirName)
	if err != nil {
		return err
	}
	err = vServ.Load(symbolSets, symsetDirName)
	return err
}

func main() {

	usage := `USAGE:
  validate_lex_file <LEXICON FILE> <SYMBOLSET NAME> <SYMBOLSET FOLDER> <PRINTMODE>

PRINTMODE: valid/invalid/all

SAMPLE INVOCATION:
  validate_lex_file [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8 sv-se_ws-sampa [SYMBOLSET FOLDER] valid`

	var args = os.Args
	if len(args) != 5 {
		fmt.Println(usage)
		os.Exit(1)
	}

	inFile := args[1]
	symbolSetName := args[2]
	symsetDirName := args[3]
	printModeS := strings.ToLower(args[4])
	var printMode dbapi.PrintMode
	if printModeS == "all" {
		printMode = dbapi.PrintAll
	} else if printModeS == "valid" {
		printMode = dbapi.PrintValid
	} else if printModeS == "invalid" {
		printMode = dbapi.PrintInvalid
	} else {
		msg := fmt.Sprintf("invalid print mode : %s", printModeS)
		log.Fatal(msg)
		return
	}

	//validator := &validation.Validator{}
	err := loadValidators(symsetDirName)
	if err != nil {
		msg := fmt.Sprintf("failed to load validators : %v", err)
		log.Fatal(msg)
		return
	}
	validator, err := vServ.ValidatorForName(symbolSetName)
	//validator = vdat
	if err != nil {
		msg := fmt.Sprintf("failed to get validator for symbol set %v : %v", symbolSetName, err)
		log.Fatal(msg)
		return
	}
	log.Println("Validator created for " + validator.Name)

	logger := dbapi.StdoutLogger{}
	err = dbapi.ValidateLexiconFile(logger, inFile, validator, printMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to validate lexicon file : %v", err)
	}
}
