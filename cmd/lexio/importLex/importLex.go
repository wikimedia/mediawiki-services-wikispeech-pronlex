package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

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

	var cmdName = "importLex"

	var validate = flag.Bool("validate", false, "validate each entry, and save the validation in the database (default: false)")
	var force = flag.Bool("force", false, "force loading of lexicon even if the symbolset is undefined (default: false)")
	//var replace = flag.Bool("replace", false, "if the lexicon already exists, delete it before importing the new input data (default: false)")
	var quiet = flag.Bool("quiet", false, "mute information logging (default: false)")
	var help = flag.Bool("help", false, "print help message")
	var createDb = flag.Bool("createdb", false, "create db if it doesn't exist")

	var engineFlag = flag.String("db_engine", "sqlite", "db engine (sqlite or mariadb)")
	var dbLocation = flag.String("db_location", "", "db location (folder for sqlite; address for mariadb)")
	var dbName = flag.String("db_name", "", "db name")
	var lexName = flag.String("lex_name", "", "lexicon name")
	var lexFile = flag.String("lex_file", "", "lexicon file")
	var locale = flag.String("locale", "", "lexicon locale")
	var ssFile = flag.String("symbolset", "", "lexicon symbolset file")

	var fatalError = false
	var dieIfEmptyFlag = func(name string, val *string) {
		if *val == "" {
			fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] flag %s is required", cmdName, name))
			fatalError = true
		}
	}

	var printUsage = func() {
		fmt.Fprintf(os.Stderr, `USAGE:
 importLex [FLAGS]

FLAGS:
`)
		flag.PrintDefaults()

		fmt.Fprintf(os.Stderr, `
SAMPLE INVOCATION:
  importLex -db_engine mariadb -db_location 'speechoid:@tcp(127.0.0.1:3306)' -lex_name sv-se.nst -locale sv_SE -lex_file [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8.gz -db_name svtest -symbolset [SYMBOLSET FOLDER]/sv-se_ws-sampa.sym 
  importLex -db_engine sqlite -db_location ~/wikispeech/sqlite -lex_name sv-se.nst -locale sv_SE -lex_file [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8.gz -db_name svtest -symbolset [SYMBOLSET FOLDER]/sv-se_ws-sampa.sym 

`)
	}

	flag.Usage = printUsage

	flag.Parse()

	args := flag.Args()

	if *help {
		printUsage()
		os.Exit(1)
	}

	if len(args) != 0 {
		printUsage()
		os.Exit(1)
	}
	dieIfEmptyFlag("db_engine", engineFlag)
	dieIfEmptyFlag("db_location", dbLocation)
	dieIfEmptyFlag("db_name", dbName)
	dieIfEmptyFlag("lex_name", lexName)
	dieIfEmptyFlag("lex_file", lexName)
	dieIfEmptyFlag("locale", locale)
	dieIfEmptyFlag("symbolset", ssFile)
	if fatalError {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] exit from unrecoverable errors", cmdName))
		os.Exit(1)
	}

	dbapi.Sqlite3WithRegex()
	dbRef := lex.DBRef(*dbName)

	// if _, err := loc.LookUp(locale); err != nil {
	// 	log.Fatalf("Invalid locale: %v", locale)
	// }

	symbolSetDir, ssFileName := filepath.Split(*ssFile)
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

	var dbm *dbapi.DBManager
	if *engineFlag == "mariadb" {
		dbm = dbapi.NewMariaDBManager()
	} else if *engineFlag == "sqlite" {
		dbm = dbapi.NewSqliteDBManager()
	} else {
		fmt.Fprintf(os.Stderr, "invalid db engine : %s\n", *engineFlag)
		os.Exit(1)
	}

	lexRef := lex.LexRef{dbRef, lex.LexName(*lexName)}

	defer dbm.CloseDB(dbRef)

	dbExists, err := dbm.DBExists(*dbLocation, dbRef)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] %v", cmdName, err))
		os.Exit(1)
	}

	if !dbExists {
		if *createDb {
			err := dbm.DefineDB(*dbLocation, dbRef)
			if err != nil {
				log.Fatalf("couldn't create db %s : %v", dbRef, err)
				return
			}
			log.Printf("Created db %s\n", dbRef)
		} else {
			log.Fatalf("DB does not exist: %s", dbRef)
			return
		}
	} else {
		err := dbm.OpenDB(*dbLocation, dbRef)
		if err != nil {
			log.Fatalf("Couldn't open db %s : %v", dbRef, err)
			return
		}
		log.Printf("Opened db %s\n", dbRef)
	}

	if !dbm.ContainsDB(dbRef) {
		log.Fatalf("DB should be registered in dbm, but wasn't : %s", dbRef)
		return
	}

	lexExists, err := dbm.LexiconExists(lexRef)
	if err != nil {
		log.Fatalf("Couldn't check if lexicon %s exists: %v", lexRef, err)
	}
	if lexExists {
		// if *replace {
		// 	log.Printf("Running SuperDelete on lexicon %s. This may take some time. Please do not abort during deletion.\n", lexRef)
		// 	err := dbm.SuperDeleteLexicon(lexRef)
		// 	if err != nil {
		// 		log.Fatalf("Couldn't super delete lexicon %s : %v", lexRef, err)
		// 		return
		// 	}
		// 	log.Printf("Deleted lexicon %s\n", lexRef)

		// } else {
		//log.Fatalf("Nothing will be added. Lexicon already exists in database: %s. Use the -replace switch if you want to replace the old lexicon.", lexRef)
		//return
		log.Fatalf("Nothing will be added. Lexicon already exists in database: %s.", lexRef)
		return
		// }
	}

	err = dbm.DefineLexicon(lexRef, symbolSetName, *locale)
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
	stderrLogger.Write(fmt.Sprintf("importing lexicon file %s ...", *lexFile))
	err = dbm.ImportLexiconFile(lexRef, logger, *lexFile, validator)

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
	stderrLogger.Write("dbLocation=" + *dbLocation)
	stderrLogger.Write("dbName=" + string(dbRef))
	stderrLogger.Write("lexName=" + *lexName)
	stderrLogger.Write("lexFile=" + *lexFile)
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
