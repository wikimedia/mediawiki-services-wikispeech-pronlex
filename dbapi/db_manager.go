package dbapi

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/stts-se/pronlex/lex"
)

type DBManager struct {
	sync.RWMutex
	dbs map[lex.DBRef]*sql.DB
}

func NewDBManager() DBManager {
	return DBManager{dbs: make(map[lex.DBRef]*sql.DB)}
}

func (dbm DBManager) AddDB(dbRef lex.DBRef, db *sql.DB) error {
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

	dbm.Lock()
	defer dbm.Unlock()

	if _, ok := dbm.dbs[dbRef]; ok {
		return fmt.Errorf("DBManager.AddDB: db already exists: '%s'", name)
	}

	dbm.dbs[dbRef] = db

	return nil
}

func (dbm DBManager) RemoveDB(dbRef lex.DBRef) error {
	name := string(dbRef)
	dbm.Lock()
	defer dbm.Unlock()

	if _, ok := dbm.dbs[dbRef]; !ok {
		return fmt.Errorf("DBManager.RemoveDB: no such db '%s'", name)
	}

	delete(dbm.dbs, dbRef)

	return nil
}

func (dbm DBManager) SuperDeleteLexicon(lexRef lex.LexRef) error {
	dbm.Lock()
	defer dbm.Unlock()

	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return fmt.Errorf("DBManager.SuperDeleteLexicon: no such db '%s'", lexRef.DBRef)
	}

	err := superDeleteLexicon(db, string(lexRef.LexName))
	if err != nil {
		return fmt.Errorf("DBManager.SuperDeleteLexicon: couldn't delete '%s'", lexRef)
	}

	return nil
}

func (dbm DBManager) DeleteLexicon(lexRef lex.LexRef) error {
	dbm.Lock()
	defer dbm.Unlock()

	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return fmt.Errorf("DBManager.DeleteLexicon: no such db '%s'", lexRef.DBRef)
	}

	err := deleteLexicon(db, string(lexRef.LexName))
	if err != nil {
		return fmt.Errorf("DBManager.SuperDeleteLexicon: couldn't delete '%s'", lexRef)
	}

	return nil
}

func (dbm DBManager) LexiconStats(lexRef lex.LexRef) (LexStats, error) {
	dbm.Lock()
	defer dbm.Unlock()

	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return LexStats{}, fmt.Errorf("DBManager.LexiconStats: no such db '%s'", lexRef.DBRef)
	}

	stats, err := lexiconStats(db, string(lexRef.LexName))
	if err != nil {
		return LexStats{}, fmt.Errorf("DBManager.LexiconStats: couldn't get stats '%s'", lexRef)
	}

	return stats, nil
}

func (dbm DBManager) DefineLexicon(dbRef lex.DBRef, symbolSetName string, lexes ...lex.LexName) error {

	dbm.RLock()
	defer dbm.RUnlock()

	for _, l := range lexes {
		db, ok := dbm.dbs[dbRef]
		if !ok {
			return fmt.Errorf("DBManager.DefineLexicon: No such db: '%s'", dbRef)
		}
		_, err := defineLexicon(db, Lexicon{Name: string(l), SymbolSetName: symbolSetName}) // TODO lex.LexName
		if err != nil {
			return fmt.Errorf("DBManager.DefineLexicon: failed to add '%s:%s' : %v", dbRef, l, err)
		}
	}

	return nil
}

func (dbm DBManager) ListDBNames() []lex.DBRef {
	var res []lex.DBRef

	dbm.RLock()
	defer dbm.RUnlock()

	for k := range dbm.dbs {
		res = append(res, k)
	}

	return res
}

// func splitFullLexiconName(fullLexName string) (string, string, error) {
// 	nameSplit := strings.SplitN(strings.TrimSpace(fullLexName), ":", 2)
// 	if len(nameSplit) != 2 {
// 		return "", "", fmt.Errorf("DBManager.splitFullLexiconName: failed to split full lexicon name into two colon separated parts: '%s'", fullLexName)
// 	}
// 	db := nameSplit[0]
// 	if "" == db {
// 		return "", "", fmt.Errorf("DBManager.splitFullLexiconName: db part of lexicon name empty: '%s'", fullLexName)
// 	}
// 	lex := nameSplit[1]
// 	if "" == lex {
// 		return "", "", fmt.Errorf("DBManager.splitFullLexiconName: lexicon part of full lexicon name empty: '%s'", fullLexName)
// 	}

// 	return db, lex, nil
// }

type lookUpRes struct {
	dbRef   lex.DBRef // TODO: move to lex.Entry!!
	entries []lex.Entry
	err     error
}

func (dbm DBManager) LookUp(q DBMQuery) (map[lex.DBRef][]lex.Entry, error) {
	var res = make(map[lex.DBRef][]lex.Entry)

	dbz := make(map[lex.DBRef][]lex.LexName)
	for _, l := range q.LexRefs {
		lexList := dbz[l.DBRef]
		dbz[l.DBRef] = append(lexList, l.LexName)
	}

	dbm.RLock()
	defer dbm.RUnlock()

	ch := make(chan lookUpRes)
	for dbR, lexs := range dbz {
		db, ok := dbm.dbs[dbR]
		if !ok {
			return res, fmt.Errorf("DBManager.LookUp: no db of name '%s'", dbR)
		}

		go func(db0 *sql.DB, dbRef lex.DBRef, lexNames []lex.LexName) {
			rez := lookUpRes{}
			rez.dbRef = dbRef
			ew := lex.EntrySliceWriter{}
			err := lookUp(db0, lexNames, q.Query, &ew)
			if err != nil {
				rez.err = fmt.Errorf("DBManager.LookUp dbapi.LookUp failed : %v", err)
				ch <- rez
				return
			}
			for _, e := range ew.Entries {
				e.LexRef.DBRef = dbRef
				rez.entries = append(rez.entries, e)
			}

			ch <- rez
		}(db, dbR, lexs)
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

func (dbm DBManager) LookUpIntoWriter(q DBMQuery, out lex.EntryWriter) error {
	var res = make(map[lex.DBRef][]lex.Entry)

	dbz := make(map[lex.DBRef][]lex.LexName)
	for _, l := range q.LexRefs {
		lexList := dbz[l.DBRef]
		dbz[l.DBRef] = append(lexList, l.LexName)
	}

	dbm.RLock()
	defer dbm.RUnlock()

	ch := make(chan lookUpRes)
	for dbR, lexs := range dbz {
		db, ok := dbm.dbs[dbR]
		if !ok {
			return fmt.Errorf("DBManager.LookUp: no db of name '%s'", dbR)
		}

		go func(db0 *sql.DB, dbRef lex.DBRef, lexNames []lex.LexName) {
			rez := lookUpRes{}
			rez.dbRef = dbRef
			ew := lex.EntrySliceWriter{}
			err := lookUp(db0, lexNames, q.Query, &ew)
			if err != nil {
				rez.err = fmt.Errorf("DBManager.LookUp dbapi.LookUp failed : %v", err)
				ch <- rez
				return
			}
			for _, e := range ew.Entries {
				e.LexRef.DBRef = dbRef
				out.Write(e)
			}

			ch <- rez
		}(db, dbR, lexs)
	}

	for i := 0; i < len(dbz); i++ {
		lkUp := <-ch
		if lkUp.err != nil {
			return fmt.Errorf("DBManager.LookUp failed : %v", lkUp.err)
		}

		if _, ok := res[lkUp.dbRef]; ok {
			return fmt.Errorf("DBManage.LookUp: returned several result for single DB '%s'", lkUp.dbRef)
		}
		res[lkUp.dbRef] = lkUp.entries
	}

	return nil
}

type lexRes struct {
	dbRef lex.DBRef
	lexs  []lex.LexName
	err   error
}

// Warning: this is maybe my first attempt at concurrency using a channel in Go

func (dbm DBManager) ListLexicons() ([]lex.LexRef, error) {
	var res []lex.LexRef

	dbm.RLock()
	defer dbm.RUnlock()

	ch := make(chan lexRes)
	defer close(ch) // ?
	dbs := dbm.dbs
	// Go ask each db instance in its own Go-routine
	for dbRef, db := range dbs {
		go func(dbRef lex.DBRef, db *sql.DB, ch0 chan lexRes) {
			lexNames, err := listLexicons(db)
			lexs := []lex.LexName{}
			for _, ln := range lexNames {
				lexs = append(lexs, lex.LexName(ln.Name))
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
			res = append(res, lex.LexRef{r.dbRef, l})
		}
	}

	return res, nil
}

func (dbm DBManager) LexiconExists(lexRef lex.LexRef) (bool, error) {
	lexRefs, err := dbm.ListLexicons()
	if err != nil {
		return false, fmt.Errorf("DBManager.LexiconExists failed : %v", err)
	}
	for _, lf := range lexRefs {
		if lf == lexRef {
			return true, nil
		}
	}
	return false, nil
}

func (dbm DBManager) InsertEntries(lexRef lex.LexRef, entries []lex.Entry) ([]int64, error) {

	var res []int64

	dbm.Lock()
	defer dbm.Unlock()

	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return res, fmt.Errorf("DBManager.InsertEntries: unknown db '%s'", lexRef.DBRef)
	}

	//_ = db
	//_ = lexName
	l, err := getLexicon(db, string(lexRef.LexName))
	//fmt.Printf("%v\n", l)
	if err != nil {
		return res, fmt.Errorf("DBManager.InsertEntries failed call to getLexicons : %v", err)
	}
	//fmt.Println(lexName)
	res, err = insertEntries(db, l, entries)
	if err != nil {
		return res, fmt.Errorf("DBManager.InsertEntries failed: %v", err)
	}
	return res, err
}

//func (dbm DBManager) UpdateEntry(dbRef lex.DBRef, e lex.Entry) (lex.Entry, bool, error) {
func (dbm DBManager) UpdateEntry(e lex.Entry) (lex.Entry, bool, error) {
	var res lex.Entry

	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[e.LexRef.DBRef]
	if !ok {
		return res, false, fmt.Errorf("DBManager.UpdateEntry: no such db '%s'", e.LexRef.DBRef)
	}

	return updateEntry(db, e)
}
