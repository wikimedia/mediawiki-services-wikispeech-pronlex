package main

import (
	"bytes"
	"compress/gzip"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// SQL DUMP:
// sqlite3 <dbFile> .dump | gzip -c > <sqlDumpFile>

// SQL LOAD:
// gunzip -c <dumpFile> | sqlite3 <dbFile>

const sqlitePath = "sqlite3"

func sqlDump(dbFile string, outFile string) error {
	sqliteCmd := exec.Command(sqlitePath, dbFile, ".dump")
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatal("couldn't create output file %s : %v", outFile, err)
	}
	defer out.Close()
	sqliteCmd.Stdout = out
	sqliteCmd.Stderr = os.Stderr
	err = sqliteCmd.Run()
	if err != nil {
		return err
	}
	log.Printf("Exported %s from db %s\n", dbFile, outFile)
	return nil
}

func validateSchemaVersion(db *sql.DB) error {
	ver, err := dbapi.GetSchemaVersion(db)
	if err != nil {
		log.Fatalf("Couldn't retrive schema version : %v", err)
	}
	if ver != dbapi.SchemaVersion {
		log.Fatalf("Mismatching schema versions -- in input file: %s, in dbapi.Schema: %s ", ver, dbapi.SchemaVersion)
	}
	log.Printf("Schema version matching dbapi.Schema: %s\n", dbapi.SchemaVersion)
	return nil
}

func runPostTests(dbFile string, sqlDumpFile string) {

	db, dbm := defineDB(dbFile)
	lexes, err := dbm.ListLexicons()

	// (1) check schema version
	err = validateSchemaVersion(db)
	if err != nil {
		log.Fatalf("Couldn't read validate schema version in file %s : %v\n", sqlDumpFile, err)
	}

	// (2) output statistics
	lexes, err = dbm.ListLexicons()
	if err != nil {
		log.Fatalf("Couldn't list lexicons : %v", err)
	}
	for _, lex := range lexes {
		stats, err := dbm.LexiconStats(lex.LexRef)
		if err != nil {
			log.Fatalf("Failed to retrieve statistics : %v", err)
		}
		err = printStats(stats, true)
		if err != nil {
			log.Fatalf("Failed to print statistics : %v", err)
		}
	}

}

func getFileReader(fName string) io.Reader {
	if strings.HasSuffix(fName, ".sql") {
		// cat <sqlDumpFile> | sqlite3 <dbFile>
		fs, err := os.Open(fName)
		if err != nil {
			log.Fatalf("Couldn't open sql dump file %s for reading : %v\n", fName, err)
		}
		return io.Reader(fs)

	} else if strings.HasSuffix(fName, ".sql.gz") {
		// zcat <sqlDumpFile> | gunzip -c | sqlite3 <dbFile>

		fh, err := os.Open(fName)
		if err != nil {
			log.Fatalf("Couldn't open file : %v", err)
		}

		if strings.HasSuffix(fName, ".gz") {
			gz, err := gzip.NewReader(fh)
			if err != nil {
				log.Fatalf("Couldn't to open gz reader : %v", err)
			}
			return io.Reader(gz)
		}

	}
	log.Fatalf("Unknown file type: %s. Expected .sql or .sql.gz", fName)
	return nil
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

func defineDB(dbFile string) (*sql.DB, *dbapi.DBManager) {
	var db *sql.DB
	var dbm = dbapi.NewDBManager()
	var err error

	dbapi.Sqlite3WithRegex()
	db, err = sql.Open("sqlite3_with_regexp", dbFile)
	if err != nil {
		log.Fatalf("Failed to open dbfile %v", err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatalf("Failed to exec PRAGMA call %v", err)
	}
	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
	if err != nil {
		log.Fatalf("Failed to exec PRAGMA call %v", err)
	}
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		log.Fatalf("Failed to exec PRAGMA call %v", err)
	}

	dbName := filepath.Base(dbFile)
	var extension = filepath.Ext(dbName)
	dbName = dbName[0 : len(dbName)-len(extension)]
	dbRef := lex.DBRef(dbName)
	err = dbm.AddDB(dbRef, db)
	if err != nil {
		log.Fatalf("Failed to add db: %v", err)
	}
	return db, dbm
}

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, `USAGE:
      importSql <SQL DUMP FILE> <NEW DB FILE>

      <SQL DUMP FILE> - sql dump of a lexicon database (.sql or .sql.gz)
      <NEW DB FILE>   - new (non-existing) db file to import into (<DBNAME>.db)
     
     SAMPLE INVOCATION:
       importSql [LEX FILE FOLDER]/swe030224NST.pron-ws.utf8.sql.gz sv_se_nst_lex.db

`)
		os.Exit(1)
	}

	var sqlDumpFile = os.Args[1]
	var dbFile = os.Args[2]

	log.Printf("Input file: %s\n", sqlDumpFile)
	log.Printf("Output db: %s\n", dbFile)

	_, err := os.Stat(sqlDumpFile)
	if err != nil {
		log.Fatalf("Input file does not exist : %v\n", err)
	}

	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		log.Fatalf("Cannot import sql dump into pre-existing database. Db file already exists: %s\n", dbFile)
	}

	sqliteCmd := exec.Command(sqlitePath, dbFile)
	stdin := sqlDumpFile
	sqliteCmd.Stdin = getFileReader(stdin)
	var sqliteOut bytes.Buffer
	sqliteCmd.Stdout = &sqliteOut
	sqliteCmd.Stderr = os.Stderr
	err = sqliteCmd.Run()
	if len(sqliteOut.String()) > 0 {
		log.Println(sqliteOut.String())
	}
	if err != nil {
		log.Fatalf("Couldn't load sql dump %s into db %s : %v\n", sqlDumpFile, dbFile, err)
	}

	log.Printf("Imported %s into db %s\n", sqlDumpFile, dbFile)

	runPostTests(dbFile, sqlDumpFile)
}
