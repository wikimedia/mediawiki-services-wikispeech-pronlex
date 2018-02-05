package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"github.com/stts-se/pronlex/dbapi"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// SQL DUMP:
/*
sqlite3 <dbFile> .dump | gzip -c > <sqlDumpFile>
**/

// SQL LOAD:
// gunzip -c <dumpFile> | sqlite3 <dbFile>

// TODO:
// * BUG(hanna) Tests: simple sanity checks after import (count entries etc...)
// * BUG(hanna) For later: Compare sql dump's schema version to the current schema.go (dbapi/schema.go). Refuse import if they don't match. Requires db change.

const sqlitePath = "sqlite3"

// change this regexp in order to check for other definitions of schema version
var schemaVersionRe = regexp.MustCompile("^\\s*PRAGMA user_version = ([0-9]+);\\s*$")

func deepCompare(file1, file2 string) bool {
	f1 := getFileReader(file1)
	f2 := getFileReader(file2)

	sscan := bufio.NewScanner(f1)
	dscan := bufio.NewScanner(f2)

	for sscan.Scan() {
		dscan.Scan()
		if !bytes.Equal(sscan.Bytes(), dscan.Bytes()) {
			return true
		}
	}

	return false
}

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
	log.Printf("Exported %s into db %s\n", dbFile, outFile)
	return nil
}

func runPostTests(dbFile string, sqlDumpFile string) {
	var testDump = "_test_" + sqlDumpFile
	var ext = filepath.Ext(testDump)
	if ext == ".gz" {
		testDump = testDump[0 : len(testDump)-len(ext)]
	}
	sqlDump(dbFile, testDump)
	defer os.Remove(testDump)
	if deepCompare(sqlDumpFile, testDump) {
		fmt.Printf("Imported db %s seems to match the input sql dump file %s\n", dbFile, sqlDumpFile)
	} else {
		log.Fatalf("Imported db %s does not match the input sql dump file %s", dbFile, sqlDumpFile)
	}
}

func validateSchemaVersion(fName string) error {
	file := getFileReader(fName)

	scanner := bufio.NewScanner(file)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		l := scanner.Text()
		if schemaVersionRe.MatchString(l) {
			matches := schemaVersionRe.FindStringSubmatch(l)
			if len(matches) >= 2 {
				schemaVersion := matches[1]
				if schemaVersion == dbapi.SchemaVersion {
					fmt.Printf("Valid schema version: %s\n", schemaVersion)
					return nil
				} else {
					log.Fatalf("Invalid schema in file %s: found %s, expected %s", fName, schemaVersion, dbapi.SchemaVersion)
				}
			} else {
				log.Fatalf("Error parsing schema version on line %d in file %s", lineNo, fName)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("No schema version in file: %s\n", fName)
	return nil
}

func getFileReader(fName string) io.Reader {
	var fileStream io.Reader

	if strings.HasSuffix(fName, ".sql") {
		// cat <sqlDumpFile> | sqlite3 <dbFile>
		fs, err := os.Open(fName)
		if err != nil {
			log.Fatalf("Couldn't open sql dump file %s for reading : %v\n", fName, err)
		}
		fileStream = io.Reader(fs)

	} else if strings.HasSuffix(fName, ".sql.gz") {
		// zcat <sqlDumpFile> | gunzip -c | sqlite3 <dbFile>

		fh, err := os.Open(fName)
		if err != nil {
			var msg = fmt.Sprintf("Couldn't open file : %v", err)
			fmt.Println(msg)
		}

		if strings.HasSuffix(fName, ".gz") {
			gz, err := gzip.NewReader(fh)
			if err != nil {
				log.Fatalf("Couldn't to open gz reader : %v", err)
			}
			fileStream = io.Reader(gz)
		}

	} else {
		log.Fatalf("Unknown file type: %s. Expected .sql or .sql.gz", fName)
	}
	return fileStream
}

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "USAGE:\nimportSql <SQL DUMP FILE (.sql or .sql.gz)> <NEW DB FILE>\n")
		os.Exit(1)
	}

	var sqlDumpFile = os.Args[1]
	var dbFile = os.Args[2]

	fmt.Printf("Input file: %s\n", sqlDumpFile)
	fmt.Printf("Output db: %s\n", dbFile)

	_, err := os.Stat(sqlDumpFile)
	if err != nil {
		log.Fatalf("Input file does not exist : %v\n", err)
	}

	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		log.Fatalf("Db file already exists: %s\n", dbFile)
	}

	err = validateSchemaVersion(sqlDumpFile)
	if err != nil {
		log.Fatalf("Couldn't read validate schema version in file %s : %v\n", sqlDumpFile, err)
	}

	sqliteCmd := exec.Command(sqlitePath, dbFile)
	stdin := sqlDumpFile
	sqliteCmd.Stdin = getFileReader(stdin)
	var sqliteOut bytes.Buffer
	sqliteCmd.Stdout = &sqliteOut
	sqliteCmd.Stderr = os.Stderr
	err = sqliteCmd.Run()
	if err != nil {
		log.Fatalf("Couldn't load sql dump %s into db %s : %v\n", sqlDumpFile, dbFile, err)
	}
	if len(sqliteOut.String()) > 0 {
		fmt.Println(sqliteOut.String())
	}

	log.Printf("Imported %s into db %s\n", sqlDumpFile, dbFile)

	runPostTests(dbFile, sqlDumpFile)
}
