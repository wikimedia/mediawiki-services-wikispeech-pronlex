package dbapi

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/stts-se/pronlex/lex"
)

type DBRef string
type LexName string

type LexRef struct {
	DBRef   DBRef
	LexName LexName
}

func NewLexRef(lexDB string, lexName string) LexRef {
	return LexRef{DBRef: DBRef(strings.ToLower(strings.TrimSpace(lexDB))),
		LexName: LexName(strings.ToLower(strings.TrimSpace(lexName))),
	}
}

type DBManager struct {
	sync.RWMutex
	dbs map[DBRef]*sql.DB
}

func NewDBManager() DBManager {
	return DBManager{dbs: make(map[DBRef]*sql.DB)}
}

func (dmb DBManager) AddDB(dbRef DBRef, db *sql.DB) error {
	name := string(dbRef)
	if "" == name {
		return fmt.Errorf("DBManager.AddDB: illegal argument: name must not be empty")
	}
	if strings.Contains(name, ":") {
		return fmt.Errorf("DBManager.AddDB: illegal argument: name must not contain ':'")
	}
	if nil == db {
		return fmt.Errorf("DBManager.AddDB: illegal argument: db must not be nil")
	}

	dmb.Lock()
	defer dmb.Unlock()

	if _, ok := dmb.dbs[dbRef]; ok {
		return fmt.Errorf("DBManager.AddDB: db already exists: '%s'", name)
	}

	dmb.dbs[dbRef] = db

	return nil
}

func (dbm DBManager) RemoveDB(dbRef DBRef) error {
	name := string(dbRef)
	dbm.Lock()
	defer dbm.Unlock()

	if _, ok := dbm.dbs[dbRef]; !ok {
		return fmt.Errorf("DBManager.RemoveDB: no such db '%s'", name)
	}

	delete(dbm.dbs, dbRef)

	return nil
}

func (dbm DBManager) DefineLexicon(dbRef DBRef, symbolSetName string, lexes ...LexName) error {

	dbm.RLock()
	defer dbm.RUnlock()

	for _, l := range lexes {
		db, ok := dbm.dbs[dbRef]
		if !ok {
			return fmt.Errorf("DBManager.DefineLexicon: No such db: '%s'", dbRef)
		}
		_, err := DefineLexicon(db, Lexicon{Name: string(l), SymbolSetName: symbolSetName}) // TODO LexName
		if err != nil {
			return fmt.Errorf("DBManager.DefineLexicon: failed to add '%s:%s' : %v", dbRef, l, err)
		}
	}

	return nil
}

func (dbm DBManager) ListDBNames() []DBRef {
	var res []DBRef

	dbm.RLock()
	defer dbm.RUnlock()

	for k, _ := range dbm.dbs {
		res = append(res, k)
	}

	return res
}

func splitFullLexiconName(fullLexName string) (string, string, error) {
	nameSplit := strings.SplitN(strings.TrimSpace(fullLexName), ":", 2)
	if len(nameSplit) != 2 {
		return "", "", fmt.Errorf("DBManager.splitFullLexiconName: failed to split full lexicon name into two colon separated parts: '%s'", fullLexName)
	}
	db := nameSplit[0]
	if "" == db {
		return "", "", fmt.Errorf("DBManager.splitFullLexiconName: db part of lexicon name empty: '%s'", fullLexName)
	}
	lex := nameSplit[1]
	if "" == lex {
		return "", "", fmt.Errorf("DBManager.splitFullLexiconName: lexicon part of full lexicon name empty: '%s'", fullLexName)
	}

	return db, lex, nil
}

type lookUpRes struct {
	dbRef   DBRef // TODO: move to lex.Entry!!
	entries []lex.Entry
	err     error
}

// TODO This turned out somewhat ugly: the Query.Lexicon field is
// overwritten by the full (DB+lexicon name) lexicon names. The Query
// will be copied and instantiated with the Lexicon field for each DB.
// ??? How to handle this in a neater way ???
func (dbm DBManager) LookUp(lexRefs []LexRef, q Query) (map[DBRef][]lex.Entry, error) {
	var res = make(map[DBRef][]lex.Entry)

	dbz := make(map[DBRef][]LexName)
	for _, l := range lexRefs {
		lexList := dbz[l.DBRef]
		dbz[l.DBRef] = append(lexList, l.LexName)
	}

	dbm.RLock()
	defer dbm.RUnlock()

	ch := make(chan lookUpRes)
	for dbN, lexs := range dbz {
		db, ok := dbm.dbs[dbN]
		if !ok {
			return res, fmt.Errorf("DBManager.LookUp: no db of name '%s'", dbN)
		}

		go func(db0 *sql.DB, dbRef DBRef, lexNames []LexName) {
			lexNameStrings := []string{}
			for _, ln := range lexNames {
				lexNameStrings = append(lexNameStrings, string(ln)) // TODO: LexName
			}
			rez := lookUpRes{}
			rez.dbRef = dbRef
			lexs0, err := GetLexicons(db0, lexNameStrings)
			if err != nil {
				rez.err = fmt.Errorf("DBManager.LookUp: failed db query : %v", err)
				ch <- rez
				return
			}
			q.Lexicons = lexs0
			ew := lex.EntrySliceWriter{}
			err = LookUp(db, q, &ew)
			if err != nil {
				rez.err = fmt.Errorf("DBManager.LookUp dbapi.LookUp failed : %v", err)
				ch <- rez
				return
			}
			rez.entries = ew.Entries
			ch <- rez
			//res[dbN] = ew.Entries
		}(db, dbN, lexs)
	}

	for i := 0; i < len(dbz); i++ {
		lkUp := <-ch
		if lkUp.err != nil {
			return res, fmt.Errorf("DBManager.LookUp failed : %v", lkUp.err)
		}

		if _, ok := res[lkUp.dbRef]; ok {
			return res, fmt.Errorf("DBManage.LookUp: returned several result for single DB '%s'", lkUp.dbRef)
		}
		res[lkUp.dbRef] = lkUp.entries
	}

	return res, nil
}

type lexRes struct {
	dbRef DBRef
	lexs  []LexName
	err   error
}

// Warning: this is maybe my first attempt at concurrency using a channel in Go

func (dbm DBManager) ListLexicons() ([]LexRef, error) {
	var res []LexRef

	dbm.RLock()
	defer dbm.RUnlock()

	ch := make(chan lexRes)
	defer close(ch) // ?
	dbs := dbm.dbs
	// Go ask each db instance in its own Go-routine
	for dbRef, db := range dbs {
		go func(dbRef DBRef, db *sql.DB, ch0 chan lexRes) {
			lexNames, err := ListLexicons(db)
			lexs := []LexName{}
			for _, ln := range lexNames {
				lexs = append(lexs, LexName(ln.Name))
			}
			r := lexRes{dbRef: dbRef, lexs: lexs, err: err}
			ch0 <- r
		}(dbRef, db, ch)
	}

	// Read result from channel
	for i := 0; i < len(dbs); i++ {
		var r lexRes
		r = <-ch // Blocks until there is a result (I
		// think). Can we be stuck here forever, if
		// db call hangs?

		// If we encounter an error, just bail out
		if r.err != nil {
			return res, fmt.Errorf("DBManager.ListLexicons freak out : %v", r.err)

		}

		// A full lexicon name consists of the name of
		// the db and the name of the lexicon in the
		// db joined by ':'
		for _, l := range r.lexs {
			res = append(res, LexRef{r.dbRef, l})
		}
	}

	return res, nil
}

func (dbm DBManager) InsertEntries(lexRef LexRef, entries []lex.Entry) ([]int64, error) {

	var res []int64

	dbm.Lock()
	defer dbm.Unlock()

	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return res, fmt.Errorf("DBManager.InsertEntries: unknown db '%s'", lexRef.DBRef)
	}

	//_ = db
	//_ = lexName
	l, err := GetLexicon(db, string(lexRef.LexName))
	//fmt.Printf("%v\n", l)
	if err != nil {
		return res, fmt.Errorf("DBManager.InsertEntries failed call to GetLexicons : %v", err)
	}
	//fmt.Println(lexName)
	res, err = InsertEntries(db, l, entries)
	if err != nil {
		return res, fmt.Errorf("DBManager.InsertEntries failed: %v", err)
	}
	return res, err
}

func (dbm DBManager) UpdateEntry(dbRef DBRef, e lex.Entry) (lex.Entry, bool, error) {
	var res lex.Entry

	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[dbRef]
	if !ok {
		return res, false, fmt.Errorf("DBManager.UpdateEntry: no such db '%s'", dbRef)
	}

	return UpdateEntry(db, e)
}
