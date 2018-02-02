package main

import (
	// "database/sql"
	"fmt"
	//_ "github.com/mattn/go-sqlite3"
	"bytes"
	"os"
	"os/exec"
)

// SQL DUMP:
// sqlite3 en_am_cmu_lex.db
// sqlite> .output "en_am_cmu_lex.sql"
// sqlite> .dump
// sqlite> .exit

// TODO: check PRAGMA user_version to ensure matches the current schema.go
// TODO: dump new db and compare to original dump file

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "USAGE:\nimportSql <SQL DUMP FILE> <NEW DB FILE>\n")
		os.Exit(1)
	}

	var sqlDumpFile = os.Args[1]
	var dbFile = os.Args[2]

	fmt.Println(sqlDumpFile)
	fmt.Println(dbFile)

	_, err := os.Stat(sqlDumpFile)
	if err != nil {
		fmt.Printf("Cannot find db file: %v\n", err)
		return
	}

	if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		fmt.Printf("Db file already exists: %s\n", dbFile)
		return
	}

	// cat <sqlDumpFile> | sqlite3 <dbFile>
	cmd := exec.Command("sqlite3", dbFile)
	stream, err := os.Open(sqlDumpFile)
	if err != nil {
		fmt.Printf("Couldn't open sql dump file %s for reading : %v\n", sqlDumpFile, err)
	}
	cmd.Stdin = stream
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Couldn't load sql dump %s into db %s : %v\n", sqlDumpFile, dbFile, err)
	}
	//fmt.Println(out.String())

	fmt.Printf("Imported %s into db %s\n", sqlDumpFile, dbFile)

}
