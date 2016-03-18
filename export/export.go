package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/line"
)

func main() {

	// TODO help message
	// TODO command line arg processing

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatalf("darn : %v", err)
	}

	// TODO read lexicon name from command line
	// + db look up
	ls := []dbapi.Lexicon{dbapi.Lexicon{ID: 1}}
	q := dbapi.Query{Lexicons: ls, Words: []string{"hundar"}}
	f, err := os.Create(os.Args[2])
	if err != nil {
		log.Fatalf("aouch : %v", err)
	}

	bf := bufio.NewWriter(f)
	defer bf.Flush()
	// bfx := dbapi.EntryFileWriter{bf}
	// dbapi.LookUp(db, q, bfx)

	//fmt.Println()

	// ew := &dbapi.EntriesSliceWriter{}
	// dbapi.LookUp(db, q, ew)
	// for _, v := range ew.Entries {
	// 	fmt.Printf("%v\n", v)
	// }

	// HL 20160318 : TEMPLATE CODE TO GENERATE NST OUTPUT:
	nstFmt, err := line.NewNST()
	if err != nil {
		log.Fatal(err)
	}
	nstW := line.NSTFileWriter{nstFmt, bf}
	dbapi.LookUp(db, q, nstW)

}
