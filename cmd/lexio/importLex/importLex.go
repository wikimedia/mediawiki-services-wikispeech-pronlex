package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/symbolset"
	//loc "github.com/stts-se/pronlex/validation/locale"
	"github.com/stts-se/pronlex/validation/validators"
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

func printStats(stats dbapi.LexStats, validate bool) error {
	var fstr = "%-16s %6d\n"
	println("\nLEXICON STATISTICS")
	fmt.Printf(fstr, "entries", stats.Entries)
	for _, s2f := range stats.StatusFrequencies {
		//fs := strings.Split(s2f, "\t")
		//if len(fs) != 2 {
		//	return fmt.Errorf("couldn't parse status-freq from string: %s", s2f)
		//}
		//var freq, err = strconv.ParseInt(fs[1], 10, 64)
		//if err != nil {
		//	return fmt.Errorf("couldn't parse status-freq from string: %s", s2f)
		//}
		var status = "status:" + s2f.Status
		var freq = s2f.Freq
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
	var quiet = f.Bool("quiet", false, "mute information logging (default: false)")
	var help = f.Bool("help", false, "print help message")
	var createDb = f.Bool("createdb", false, "create db if it doesn't exist")

	usage := `USAGE:
 importLex [FLAGS] <DB FILE> <LEXICON NAME> <LOCALE> <LEXICON FILE> <SYMBOLSET FILE>

FLAGS:
   -validate bool  validate each entry, and save the validation in the database (default: false)
   -force    bool  force loading of lexicon even if the symbolset is undefined (default: false)
   -replace  bool  if the lexicon already exists, delete it before importing the new input data (default: false)
   -createdb bool  create db if it doesn't exist (default: false)
   -quiet    bool  mute information logging (default: false)
   -help     bool  print help message

SAMPLE INVOCATION:
  importLex -validate pronlex.db sv-se.nst sv_SE [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8.gz [SYMBOLSET FOLDER]/sv-se_ws-sampa.sym`

	f.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	var args = os.Args
	if strings.HasSuffix(args[0], "importLex") {
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

	dbapi.Sqlite3WithRegex()
	dbFile := args[0]
	lexName := args[1]
	locale := args[2]
	inFile := args[3]
	symbolSetFile := args[4]

	// if _, err := loc.LookUp(locale); err != nil {
	// 	log.Fatalf("Invalid locale: %v", locale)
	// }

	symbolSetDir, ssFileName := filepath.Split(symbolSetFile)
	ext := filepath.Ext(ssFileName)
	symbolSetName := ssFileName[0 : len(ssFileName)-len(ext)]

	validator := &validation.Validator{}
	if *validate {

		err := loadValidators(symbolSetDir)
		if err != nil {
			msg := fmt.Sprintf("Failed to load validators : %v", err)
			log.Fatal(msg)
			return
		}
		vdat, err := vServ.ValidatorForName(symbolSetName)
		validator = vdat
		if err != nil {
			msg := fmt.Sprintf("Failed to get validator for symbol set %v : %v", symbolSetName, err)
			log.Fatal(msg)
			return
		}
		log.Println("Validator created for " + validator.Name)
	} else if !*force { // do not validate but still check symbolset
		_, err := symbolset.LoadSymbolSet(filepath.Join(symbolSetDir, symbolSetName+symbolset.SymbolSetSuffix))
		if err != nil {
			msg := fmt.Sprintf("Failed to load symbol set %v : %v", symbolSetName, err)
			log.Fatal(msg)
			return
		}
	}

	dbm := dbapi.NewDBManager()
	lexRef := lex.NewLexRef(dbFile, lexName)
	dbRef := lexRef.DBRef

	defer dbm.CloseDB(dbRef)

	if _, err := os.Stat(dbFile); os.IsNotExist(err) {
		if *createDb {
			err := dbm.DefineDB(dbRef, dbFile)
			if err != nil {
				log.Fatalf("couldn't create db %s : %v", dbFile, err)
				return
			}
			log.Printf("Created db %s\n", dbFile)
		} else {
			log.Fatalf("DB does not exist %s : %v", dbFile, err)
			return
		}
	} else {
		err := dbm.OpenDB(dbRef, dbFile)
		if err != nil {
			log.Fatalf("Couldn't open db %s : %v", dbFile, err)
			return
		}
		log.Printf("Opened db %s\n", dbFile)
	}

	if !dbm.ContainsDB(dbRef) {
		log.Fatalf("DB should be registered in dbm, but wasn't : %s", dbRef)
		return
	}

	lexExists, err := dbm.LexiconExists(lexRef)
	if err != nil {
		log.Fatalf("Couldn't super delete lexicon %s: %v", lexRef, err)
	}
	if lexExists {
		if *replace {
			log.Printf("Running SuperDelete on lexicon %s. This may take some time. Please do not abort during deletion.\n", lexRef)
			err := dbm.SuperDeleteLexicon(lexRef)
			if err != nil {
				log.Fatalf("Couldn't super delete lexicon %s : %v", lexRef, err)
				return
			}
			log.Printf("Deleted lexicon %s\n", lexRef)

		} else {
			log.Fatalf("Nothing will be added. Lexicon already exists in database: %s. Use the -replace switch if you want to replace the old lexicon.", lexRef)
			return
		}
	}

	err = dbm.DefineLexicon(lexRef, symbolSetName, locale)
	if err != nil {
		log.Fatal(err)
		return
	}

	var logger dbapi.Logger
	var stderrLogger = dbapi.StderrLogger{20000}

	if *quiet {
		logger = dbapi.SilentLogger{}
	} else {
		logger = stderrLogger
	}
	// TODO handle errors? Does it make sent to return array of error...?
	stderrLogger.Write(fmt.Sprintf("importing lexicon file %s ...", inFile))
	err = dbm.ImportLexiconFile(lexRef, logger, inFile, validator)

	if err != nil {
		log.Fatal(err)
		return
	}

	// stderrLogger.Write("running the Sqlite3 ANALYZE command. It may take a little while...")
	// _, err = db.Exec("ANALYZE")
	// if err != nil {
	// 	stderrLogger.Write(fmt.Sprintf("failed to run ANALYZE command : %v", err))
	// 	return
	// }

	fmt.Fprintf(os.Stderr, "\n")
	stderrLogger.Write("finished importing lexicon file")
	stderrLogger.Write("dbFile=" + dbFile)
	stderrLogger.Write("lexName=" + lexName)
	stderrLogger.Write("lexFile=" + inFile)
	stderrLogger.Write("symbolSet=" + symbolSetName)
	stderrLogger.Write("symbolSetFolder=" + symbolSetDir)
	stderrLogger.Write("validate=" + strconv.FormatBool(*validate))
	fmt.Fprintf(os.Stderr, "\n")

	stats, err := dbm.LexiconStats(lexRef)
	if err != nil {
		stderrLogger.Write(fmt.Sprintf("failed to retrieve statistics : %v", err))
		return
	}
	err = printStats(stats, *validate)
	if err != nil {
		stderrLogger.Write(fmt.Sprintf("failed to print statistics : %v", err))
		return
	}
}
