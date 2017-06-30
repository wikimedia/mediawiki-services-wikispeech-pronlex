package dbapi

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
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
		return fmt.Errorf("DBManager.DeleteLexicon: couldn't delete '%s' : %v", lexRef, err)
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
		return LexStats{}, fmt.Errorf("DBManager.LexiconStats: couldn't get stats '%s' : %v", lexRef, err)
	}

	return stats, nil
}

func (dbm DBManager) DefineLexicons(dbRef lex.DBRef, symbolSetName string, lexes ...lex.LexName) error {

	dbm.RLock()
	defer dbm.RUnlock()

	for _, l := range lexes {
		db, ok := dbm.dbs[dbRef]
		if !ok {
			return fmt.Errorf("DBManager.DefineLexicon: No such db: '%s'", dbRef)
		}
		_, err := defineLexicon(db, lexicon{name: string(l), symbolSetName: symbolSetName})
		if err != nil {
			return fmt.Errorf("DBManager.DefineLexicon: failed to add '%s:%s' : %v", dbRef, l, err)
		}
	}

	return nil
}

func (dbm DBManager) DefineLexicon(lexRef lex.LexRef, symbolSetName string) error {

	dbm.RLock()
	defer dbm.RUnlock()

	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return fmt.Errorf("DBManager.DefineLexicon: No such db: '%s'", lexRef.DBRef)
	}
	_, err := defineLexicon(db, lexicon{name: string(lexRef.LexName), symbolSetName: symbolSetName})
	if err != nil {
		return fmt.Errorf("DBManager.DefineLexicon: failed to add '%s' : %v", lexRef.String(), err)
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

type lookUpRes struct {
	dbRef   lex.DBRef // TODO: move to lex.Entry!!
	entries []lex.Entry
	err     error
}

func (dbm DBManager) LookUpIntoSlice(q DBMQuery) ([]lex.Entry, error) {
	var res = []lex.Entry{}
	writer := lex.EntrySliceWriter{}
	err := dbm.LookUp(q, &writer)
	if err != nil {
		return res, fmt.Errorf("DBManager.LookUp failed : %v", err)
	}
	for _, e := range writer.Entries {
		res = append(res, e)
	}
	return res, nil
}

func (dbm DBManager) LookUpIntoMap(q DBMQuery) (map[lex.DBRef][]lex.Entry, error) {
	var res = make(map[lex.DBRef][]lex.Entry)
	writer := lex.EntrySliceWriter{}
	err := dbm.LookUp(q, &writer)
	if err != nil {
		return res, fmt.Errorf("DBManager.LookUp failed : %v", err)
	}
	for _, e := range writer.Entries {
		es := res[e.LexRef.DBRef]
		es = append(es, e)
		res[e.LexRef.DBRef] = es
	}
	return res, nil
}

func (dbm DBManager) LookUp(q DBMQuery, out lex.EntryWriter) error {
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
				rez.entries = append(rez.entries, e)
			}

			ch <- rez
		}(db, dbR, lexs)
	}

	for i := 0; i < len(dbz); i++ {
		lkUp := <-ch
		if lkUp.err != nil {
			return fmt.Errorf("DBManager.LookUp failed : %v", lkUp.err)
		}

		for _, e := range lkUp.entries {
			out.Write(e)
		}
	}

	return nil
}

type lexRes struct {
	lexes []lex.LexRefWithInfo
	err   error
}

// Warning: this is maybe my first attempt at concurrency using a channel in Go

func (dbm DBManager) ListLexicons() ([]lex.LexRefWithInfo, error) {
	var res = []lex.LexRefWithInfo{}

	dbm.RLock()
	defer dbm.RUnlock()

	ch := make(chan lexRes)
	defer close(ch) // ?
	dbs := dbm.dbs
	// Go ask each db instance in its own Go-routine
	for dbRef, db := range dbs {
		go func(dbRef lex.DBRef, db *sql.DB, ch0 chan lexRes) {
			lexs, err := listLexicons(db)
			lexList := []lex.LexRefWithInfo{}
			for _, ln := range lexs {
				lexRef := lex.LexRef{dbRef, lex.LexName(ln.name)}
				withInfo := lex.LexRefWithInfo{lexRef, ln.symbolSetName}
				lexList = append(lexList, withInfo)
			}
			r := lexRes{lexes: lexList, err: err}
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
		for _, lex := range r.lexes {
			res = append(res, lex)
		}
	}

	return res, nil
}

func (dbm DBManager) LexiconExists(lexRef lex.LexRef) (bool, error) {
	lexInfo, err := dbm.ListLexicons()
	if err != nil {
		return false, fmt.Errorf("DBManager.LexiconExists failed : %v", err)
	}
	for _, lf := range lexInfo {
		if lf.LexRef == lexRef {
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

func (dbm DBManager) ImportLexiconFile(lexRef lex.LexRef, logger Logger, lexiconFileName string, validator *validation.Validator) error {
	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return fmt.Errorf("DBManager.ImportLexiconFile: no such db '%s'", lexRef.DBRef)
	}
	return ImportLexiconFile(db, lexRef.LexName, logger, lexiconFileName, validator)
}

func (dbm DBManager) EntryCount(lexRef lex.LexRef) (int64, error) {
	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return 0, fmt.Errorf("DBManager.ImportLexiconFile: no such db '%s'", lexRef.DBRef)
	}
	return entryCount(db, string(lexRef.LexName))
}

func (dbm DBManager) ListCurrentEntryStatuses(lexRef lex.LexRef) ([]string, error) {
	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return []string{}, fmt.Errorf("DBManager.ListCurrentEntryStatuses: no such db '%s'", lexRef.DBRef)
	}
	return listCurrentEntryStatuses(db, string(lexRef.LexName))
}

func (dbm DBManager) ListAllEntryStatuses(lexRef lex.LexRef) ([]string, error) {
	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return []string{}, fmt.Errorf("DBManager.ListAllEntryStatuses: no such db '%s'", lexRef.DBRef)
	}
	return listAllEntryStatuses(db, string(lexRef.LexName))
}

func (dbm DBManager) GetLexicon(lexRef lex.LexRef) (lex.LexRefWithInfo, error) {
	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return lex.LexRefWithInfo{}, fmt.Errorf("DBManager.GetLexicon: no such db '%s'", lexRef.DBRef)
	}
	l, err := getLexicon(db, string(lexRef.LexName))
	if err != nil {
		return lex.LexRefWithInfo{}, err
	}
	return lex.LexRefWithInfo{
		LexRef:        lexRef,
		SymbolSetName: l.symbolSetName,
	}, nil
}

// MoveNewEntries moves lexical entries from the lexicon named
// fromLexicon to the lexicon named toLexicon.  The 'newSource' string is
// the name of the new source of the entries to be moved, and 'newStatus' is
// the name of the new status to set on the moved entries.  Currently,
// source and/or status may not be the empty string. TODO: Maybe it
// should be possible to skip source and status values?
//
// Only "new" entries are moved, i.e., entries with lex.Entry.Strn
// values found in fromLexicon but *not* found in toLexicon.  The
// rationale behind this function is to first create a small
// additional lexicon with new entries (the fromLexicon), that can
// later be appended to the master lexicon (the toLexicon).
func (dbm DBManager) MoveNewEntries(dbRef lex.DBRef, fromLex, toLex lex.LexName, newSource, newStatus string) (MoveResult, error) {
	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[dbRef]
	if !ok {
		return MoveResult{}, fmt.Errorf("DBManager.MoveNewEntries: no such db '%s'", dbRef)
	}
	return moveNewEntries(db, string(fromLex), string(toLex), newSource, newStatus)
}

func (dbm DBManager) Validate(lexRef lex.LexRef, logger Logger, vd validation.Validator, q Query) (ValStats, error) {
	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return ValStats{}, fmt.Errorf("DBManager.Validate: no such db '%s'", lexRef.DBRef)
	}
	return Validate(db, []lex.LexName{lexRef.LexName}, logger, vd, q)
}

func (dbm DBManager) ValidationStats(lexRef lex.LexRef) (ValStats, error) {
	dbm.Lock()
	defer dbm.Unlock()
	db, ok := dbm.dbs[lexRef.DBRef]
	if !ok {
		return ValStats{}, fmt.Errorf("DBManager.Validate: no such db '%s'", lexRef.DBRef)
	}
	return validationStats(db, string(lexRef.LexName))
}

func (dbm DBManager) ListDatabases() ([]lex.DBRef, error) {
	dbm.Lock()
	defer dbm.Unlock()
	res := []lex.DBRef{}
	for dbRef, _ := range dbm.dbs {
		res = append(res, dbRef)
	}
	return res, nil
}
