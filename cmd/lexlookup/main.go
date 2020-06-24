package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	// TODO: Look at the mysql import
	// Why isn't the below import needed?
	// The Sqlite3 driver is needed.
	//_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"

	//"github.com/pkg/profile"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

// TODO: Better command line flags structure if adding more options.

// TODO: Sort output entries according to input order

// TODO: update, set preferred, output preferences
// (prettyprint, JSON or text, verbosity, etc)

// returns true if string contains '_' or '%'
func isLikeExpression(s string) bool {
	if strings.Contains(s, "_") {
		return true
	}

	if strings.Contains(s, "%") {
		return true
	}

	return false
}

func deleteEntry(dbm *dbapi.DBManager, entryID int64, dbRef, lexName string) error {
	lexRef := lex.LexRef{DBRef: lex.DBRef(dbRef), LexName: lex.LexName(lexName)}
	_, err := dbm.DeleteEntry(entryID, lexRef)

	return err
}

func remDupes(words []string) []string {
	var res []string
	found := make(map[string]bool)
	for _, w := range words {
		w0 := strings.ToLower(w)
		if !found[w0] {
			res = append(res, w)
			found[w0] = true
		}
	}

	return res
}

// Max number of words in each lookup call
const lookUpChunk = 100

func lookUp0(dbRef lex.DBRef, q dbapi.DBMQuery, dbm *dbapi.DBManager) ([]lex.Entry, error) {
	var res []lex.Entry

	lexica, err := dbm.ListLexicons()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to list lexicons in db '%s' %v\n", dbRef, err)
		os.Exit(1)
	}
	lexRefs := []lex.LexRef{}
	for _, l := range lexica {
		lexRef := lex.NewLexRef(string(dbRef), string(l.LexRef.LexName))
		lexRefs = append(lexRefs, lexRef)
	}

	q.LexRefs = lexRefs //[]lex.LexRef{lex.NewLexRef("db_file", "sv-se.nst")}

	resWriter := lex.EntrySliceWriter{}

	err = dbm.LookUp(q, &resWriter)
	if err != nil {
		return res, fmt.Errorf("failed database look up : %v", err)
	}

	return resWriter.Entries, nil
}

func lookUp(words []string, dbRef lex.DBRef, dbm *dbapi.DBManager) ([]lex.Entry, error) {
	var res []lex.Entry

	q := dbapi.NewQuery()
	// PageLength defaults to 25
	q.PageLength = 20000

	// The single word supplied is a 'LIKE' match query
	if len(words) == 1 && isLikeExpression(words[0]) {
		q.WordLike = words[0]
		dbmq := dbapi.DBMQuery{Query: q}
		return lookUp0(dbRef, dbmq, dbm)
	}

	var chunk []string
	for i, w := range words {
		//fmt.Fprintf(os.Stderr, "n: %d\n", i+1)
		chunk = append(chunk, w)
		if (i >= lookUpChunk && i%lookUpChunk == 0) || i == len(words)-1 {

			q.Words = chunk
			l, err := lookUp0(dbRef, dbapi.DBMQuery{Query: q}, dbm)
			if err != nil {
				return res, err
			}
			res = append(res, l...)

			//fmt.Fprintf(os.Stderr, "SO FAR: %d\n", len(chunk))

			chunk = []string{}
		}
	}

	//fmt.Fprintf(os.Stderr, "TOTO: %d\n", len(words))

	return res, nil
}

func main() {

	var err error
	var cmdName = "lexlookup"
	//defer profile.Start().Stop()

	verb := true

	deleteFlag := flag.Bool("delete", false, "Delete entry. Required flags: -id <int> -db_engine <string> -db_location <string> -db_name <string> -lex_name <string>")
	idFlag := flag.Int("id", 0, "DB entry id")

	printMissingFlag := flag.Bool("missing", false, "Print the words not found in the lexicon. Required flags: -db_engine <string> -db_location <string> -db_name <string> -lex_name <string>")

	engineFlag := flag.String("db_engine", "sqlite", "db engine (sqlite or mariadb)")
	dbLocation := flag.String("db_location", "", "DB location (folder for sqlite; address for mariadb)")
	dbName := flag.String("db_name", "", "DB reference name (for sqlite, it should be without the .db suffix")
	lexName := flag.String("lexicon", "", "Lexicon name")

	var fatalError = false
	var dieIfEmptyFlag = func(name string, val *string) {
		if *val == "" {
			fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] flag %s is required", cmdName, name))
			fatalError = true
		}
	}

	var printUsage = func() {
		fmt.Fprintln(os.Stderr, `USAGE: lexlookup (<words...> | <stdin>)

If a single word is supplied, it may contain the characters '%%' and '_' for LIKE string search.

Lookup (list all words starting with 'k'):
lexlookup -db_engine sqlite -db_location ~/wikispeech/sqlite/ -db_name wikispeech_lexserver_demo 'k%'

Lookup (specific words):
lexlookup -db_engine sqlite -db_location ~/wikispeech/sqlite/ -db_name wikispeech_lexserver_demo 'hunden'

Print missing words: 
lexlookup -db_engine sqlite -db_location ~/wikispeech/sqlite/ -db_name wikispeech_lexserver_demo -missing <words>

Deleting a DB entry:
lexlookup -db_engine sqlite -db_location ~/wikispeech/sqlite/ -db_name wikispeech_lexserver_demo -delete -id <int> -db_ref <string> -lex_name <string>

Flags:
`)
		flag.PrintDefaults()
	}

	flag.Usage = func() {
		printUsage()
		os.Exit(0)
	}
	flag.Parse()
	// if len(flag.Args()) < 1 {
	// 	printUsage()
	// 	os.Exit(0)
	// }

	dieIfEmptyFlag("db_engine", engineFlag)
	dieIfEmptyFlag("db_location", dbLocation)
	dieIfEmptyFlag("db_name", dbName)
	if *deleteFlag {
		dieIfEmptyFlag("lexicon", lexName)
	}
	if fatalError {
		fmt.Fprintln(os.Stderr, fmt.Errorf("[%s] exit from unrecoverable errors", cmdName))
		os.Exit(1)
	}

	prettyPrint := true

	var db *sql.DB
	var dbEngine dbapi.DBEngine
	var dbm *dbapi.DBManager
	if *engineFlag == "mariadb" {
		dbEngine = dbapi.MariaDB
		dbm = dbapi.NewMariaDBManager()
	} else if *engineFlag == "sqlite" {
		dbEngine = dbapi.Sqlite
		dbm = dbapi.NewSqliteDBManager()
	} else {
		fmt.Fprintf(os.Stderr, "invalid db engine : %s\n", *engineFlag)
		os.Exit(1)
	}

	if dbEngine == dbapi.Sqlite { // Sqlite

		dbPath := path.Join(*dbLocation, *dbName+".db")
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "ERROR: could not find pronlex db file '%s'\n", dbPath)
			os.Exit(1)
		}

		//dbPath := path.Base(dbPath)

		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to open pronlex Sqlite3 db file '%s' : %v\n", dbPath, err)
			os.Exit(1)
		}
	} else if dbEngine == dbapi.MariaDB { // MySQL/MariaDB
		dbPath := path.Join(*dbLocation, *dbName)
		db, err = sql.Open("mysql", dbPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to connect to MySQL/MariaDB db '%s' : %v\n", dbPath, err)
			os.Exit(1)

		}

	}

	err = dbm.AddDB(lex.DBRef(*dbName), db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to initialise db manager : %v\n", err)
		os.Exit(1)
	}

	// Delete entry
	if *deleteFlag {
		err := deleteEntry(dbm, int64(*idFlag), *dbName, *lexName)

		if err != nil {

			fmt.Fprintf(os.Stderr, "Failed to delete entry : %v\n", err)
		}

		return
	}

	// Look up

	words := []string{}
	// read from stdin
	if len(flag.Args()) == 0 {
		s := bufio.NewScanner(os.Stdin)
		r, err := regexp.Compile(`\s+`)
		if err != nil {
			fmt.Fprintf(os.Stderr, "decomper: split regexp failure : %v\n", err)
			os.Exit(1)
		}
		for s.Scan() {
			l := s.Text()
			lWds := r.Split(l, -1)

			// Don't think there can be empty strings here, but let's throw them away anyway...
			for _, w := range lWds {
				if strings.TrimSpace(w) == "" {
					continue
				}
				words = append(words, w)
			}
		}
	} else {
		words = flag.Args() //os.Args[2:]
	}

	if len(words) == 0 {
		fmt.Fprintf(os.Stderr, "lexlookup: No input words supplied\n")
		os.Exit(0)
	}

	// Only look up same string once
	words = remDupes(words)

	entries, err := lookUp(words, lex.DBRef(*dbName), dbm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed look-up : '%v'\n", err)
		os.Exit(1)
	}

	// No match
	if len(entries) == 0 && !*printMissingFlag {
		if verb {
			fmt.Fprintf(os.Stderr, "lexlookup: no matching entry in db '%s'\n", *dbName)
		}
		return
	}

	// Print input words *not* found in the lexicon
	if *printMissingFlag {
		foundWords := make(map[string]bool)
		for _, e := range entries {
			//fmt.Printf("%#v\n", e)
			f := strings.ToLower(e.Strn)
			foundWords[f] = true
		}

		for _, w := range words {
			w = strings.ToLower(w)
			if !foundWords[w] {
				fmt.Println(w)
			}

		}

		return
	}

	//for _, e := range resWriter.Entries {
	var jsn []byte

	if prettyPrint {
		jsn, err = json.MarshalIndent(entries, "", "    ")
	} else {
		jsn, err = json.Marshal(entries)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to produce JSON of database result : %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", jsn)
	//}

}
