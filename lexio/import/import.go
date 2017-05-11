package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

func printStats(stats dbapi.LexStats, validate bool) error {
	var fstr = "%-16s %6d\n"
	println("\nLEXICON STATISTICS")
	fmt.Printf(fstr, "entries", stats.Entries)
	for _, s2f := range stats.StatusFrequencies {
		fs := strings.Split(s2f, "\t")
		if len(fs) != 2 {
			return fmt.Errorf("couldn't parse status-freq from string: %s", s2f)
		}
		var freq, err = strconv.ParseInt(fs[1], 10, 64)
		if err != nil {
			return fmt.Errorf("couldn't parse status-freq from string: %s", s2f)
		}
		var status = "status:" + fs[0]
		fmt.Printf(fstr, status, freq)
	}
	if validate {
		fmt.Printf(fstr, "invalid entries", stats.ValStats.InvalidEntries)
		fmt.Printf(fstr, "validation msgs", stats.ValStats.TotalValidations)
	}
	return nil
}

func main() {

	var f = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	var validate = f.Bool("validate", false, "validate each entry, and save the validation in the database (default: false)")
	var force = f.Bool("force", false, "force loading of lexicon even if the symbolset is undefined (default: false)")
	var replace = f.Bool("replace", false, "if the lexicon already exists, delete it before importing the new input data (default: false)")
	var help = f.Bool("help", false, "print help message")

	usage := `USAGE:
  go run import.go <FLAGS> <DB FILE> <LEXICON NAME> <LEXICON FILE> <SYMBOLSET NAME> <SYMBOLSET FOLDER>

FLAGS:
   -validate bool  validate each entry, and save the validation in the database (default: false)
   -force    bool  force loading of lexicon even if the symbolset is undefined (default: false)
   -replace  bool  if the lexicon already exists, delete it before importing the new input data (default: false)
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

	lexiconExists := false
	lexicon, err := dbapi.GetLexicon(db, lexName)
	if err == nil {
		lexiconExists = true
		if *replace {
			log.Printf("Running SuperDelete on lexicon %s. This may take some time. Please do not abort during deletion.\n", lexName)
			err := dbapi.SuperDeleteLexicon(db, lexicon.ID)
			if err != nil {
				log.Fatalf("Couldn't super delete lexicon %s : %s", lexName, err)
				return
			}
			log.Printf("Deleted lexicon %s\n", lexName)
			lexicon = dbapi.Lexicon{Name: lexName, SymbolSetName: symbolSetName}
			lexicon, err = dbapi.InsertLexicon(db, lexicon)
			if err != nil {
				log.Fatal(err)
				return
			}

		} else {
			log.Fatalf("Nothing will be added. Lexicon already exists in database: %s. Use the -replace switch if you want to replace the old lexicon.", lexName)
			return
		}
	}

	if !lexiconExists {
		lexicon = dbapi.Lexicon{Name: lexName, SymbolSetName: symbolSetName}
		lexicon, err = dbapi.InsertLexicon(db, lexicon)
		if err != nil {
			log.Fatal(err)
			return
		}
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

	fmt.Fprintf(os.Stderr, "\n")
	logger.Write("finished importing lexicon file")
	logger.Write("dbFile=" + dbFile)
	logger.Write("lexName=" + lexName)
	logger.Write("lexFile=" + inFile)
	logger.Write("symbolSet=" + symbolSetName)
	logger.Write("symbolSetFolder=" + symsetDirName)
	logger.Write("validate=" + strconv.FormatBool(*validate))
	fmt.Fprintf(os.Stderr, "\n")

	stats, err := dbapi.LexiconStats(db, lexicon.ID)
	if err != nil {
		logger.Write(fmt.Sprintf("failed to retreive statistics : %v", err))
		return
	}
	err = printStats(stats, *validate)
	if err != nil {
		logger.Write(fmt.Sprintf("failed to print statistics : %v", err))
		return
	}
}
