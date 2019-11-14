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

const usage = `USAGE: lexlookup <Sqlite3 pronlex DB file> (<words...> | <stdin>)

If a single word is supplied, it may contain the characters '%%' and '_' for LIKE string search.


Print missing words: 
lexlookup <Sqlite3 pronlex DB file> -missing <words>

Deleting a DB entry:
lexlookup <Sqlite3 pronlex DB file> -delete -id <int> -db_ref <string> -db_name <string>

`

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

func lookUp0(dbFileName string, q dbapi.DBMQuery, dbm *dbapi.DBManager) ([]lex.Entry, error) {
	var res []lex.Entry

	lexica, err := dbm.ListLexicons()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to list lexicons in db file '%s' %v\n", dbFileName, err)
		os.Exit(1)
	}
	lexRefs := []lex.LexRef{}
	for _, l := range lexica {
		lexRef := lex.NewLexRef(dbFileName, string(l.LexRef.LexName))
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

func lookUp(words []string, dbFileName string, dbm *dbapi.DBManager) ([]lex.Entry, error) {
	var res []lex.Entry

	q := dbapi.NewQuery()
	// PageLength defaults to 25
	q.PageLength = 20000

	// The single word supplied is a 'LIKE' match query
	if len(words) == 1 && isLikeExpression(words[0]) {
		q.WordLike = words[0]
		dbmq := dbapi.DBMQuery{Query: q}
		return lookUp0(dbFileName, dbmq, dbm)
	}

	var chunk []string
	for i, w := range words {
		//fmt.Fprintf(os.Stderr, "n: %d\n", i+1)
		chunk = append(chunk, w)
		if (i >= lookUpChunk && i%lookUpChunk == 0) || i == len(words)-1 {

			q.Words = chunk
			l, err := lookUp0(dbFileName, dbapi.DBMQuery{Query: q}, dbm)
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

	//defer profile.Start().Stop()

	verb := true

	var flags = flag.NewFlagSet("lexlookup", flag.ExitOnError)

	deleteFlag := flags.Bool("delete", false, "Delete entry. Required flags: -id <int> -db_ref <string> -lex_name <string>")
	idFlag := flags.Int("id", 0, "DB entry id")
	dbFlag := flags.String("db_ref", "", "DB reference name")
	lexFlag := flags.String("lex_name", "", "Lexicon name")

	printMissingFlag := flags.Bool("missing", false, "Print the words not found in the lexicon. Required flags: -id <int> -db_ref <string> -lex_name <string>")

	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(0)
	}

	flags.Parse(os.Args[2:])

	prettyPrint := true

	dbPath := os.Args[1]

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "ERROR: could not find pronlex db file '%s'\n", dbPath)
		os.Exit(1)
	}

	dbFileName := path.Base(dbPath)

	dbm := dbapi.NewDBManager()

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to open pronlex Sqlite3 di file '%s' : %v\n", dbFileName, err)
		os.Exit(1)
	}

	err = dbm.AddDB(lex.DBRef(dbFileName), db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to initialise db manager : %v\n", err)
		os.Exit(1)
	}

	// Delete entry
	if *deleteFlag {
		err := deleteEntry(dbm, int64(*idFlag), *dbFlag, *lexFlag)

		if err != nil {

			fmt.Fprintf(os.Stderr, "Failed to delete entry : %v\n", err)
		}

		return
	}

	// Look up

	words := []string{}
	// read from stdin
	if len(flags.Args()) == 0 {
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
		words = flags.Args() //os.Args[2:]
	}

	if len(words) == 0 {
		fmt.Fprintf(os.Stderr, "lexlookup: No input words supplied\n")
		os.Exit(0)
	}

	// Only look up same string once
	words = remDupes(words)

	entries, err := lookUp(words, dbFileName, dbm)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed look-up : '%v'\n", err)
		os.Exit(1)
	}

	// No match
	if len(entries) == 0 && !*printMissingFlag {
		if verb {
			fmt.Fprintf(os.Stderr, "lexlookup: no matching entry in db '%s'\n", dbFileName)
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
