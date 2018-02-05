package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// SQL DUMP:
// sqlite3 <dbFile> .dump | gzip -c > <dumpFile>

// SQL LOAD:
// gunzip -c <dumpFile> | sqlite3 <dbFile>

// TODO:
// * Jämför dumpens PRAGMA user_version med aktuellt schema (dbapi/schema.go). Vägra importera om det inte matchar.
// * Test: dump new db and compare to original dump file
// * Test: köra några enkla sanity check-anrop, t.ex. select count entries osv..

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
		fmt.Printf("Input file does not exist : %v\n", err)
		return
	}

	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		fmt.Printf("Db file already exists: %s\n", dbFile)
		return
	}

	var fileStream io.Reader

	if strings.HasSuffix(sqlDumpFile, ".sql") {
		// cat <sqlDumpFile> | sqlite3 <dbFile>
		fs, err := os.Open(sqlDumpFile)
		if err != nil {
			fmt.Printf("Couldn't open sql dump file %s for reading : %v\n", sqlDumpFile, err)
			return
		}
		fileStream = io.Reader(fs)

	} else if strings.HasSuffix(sqlDumpFile, ".sql.gz") {
		// zcat <sqlDumpFile> | gunzip -c | sqlite3 <dbFile>

		fh, err := os.Open(sqlDumpFile)
		defer fh.Close()
		if err != nil {
			var msg = fmt.Sprintf("Couldn't open file : %v", err)
			fmt.Println(msg)
			return
		}

		if strings.HasSuffix(sqlDumpFile, ".gz") {
			gz, err := gzip.NewReader(fh)
			if err != nil {
				var msg = fmt.Sprintf("Couldn't to open gz reader : %v", err)
				fmt.Println(msg)
				return
			}
			fileStream = io.Reader(gz)
		}

	} else {
		fmt.Println("Unknown file type: %s. Expected .sql or .sql.gz", sqlDumpFile)
		return
	}

	sqliteCmd := exec.Command("sqlite3", dbFile)
	sqliteCmd.Stdin = fileStream
	var sqliteOut bytes.Buffer
	sqliteCmd.Stdout = &sqliteOut
	sqliteCmd.Stderr = os.Stderr
	err = sqliteCmd.Run()
	if err != nil {
		fmt.Printf("Couldn't load sql dump %s into db %s : %v\n", sqlDumpFile, dbFile, err)
		return
	}
	if len(sqliteOut.String()) > 0 {
		fmt.Println(sqliteOut.String())
	}

	fmt.Printf("Imported %s into db %s\n", sqlDumpFile, dbFile)

}
