/*
The dbapi package contains code wrapped around an SQL(ite3) DB.
*/

package dbapi

//go get github.com/mattn/go-sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	// installs sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"io"
	"strconv"
	"strings"
)

// ListLexicons returns a list of the lexicons defined in the db
// (i.e., the rows of the lexicon table)
func ListLexicons(db *sql.DB) ([]Lexicon, error) {
	var res []Lexicon
	sql := "select id, name, symbolsetname from lexicon"
	rows, err := db.Query(sql)
	defer rows.Close()
	if err != nil {
		return res, fmt.Errorf("db query failed : %v", err)
	}

	for rows.Next() {
		var id int64
		var name string
		var symbolSetName string
		err = rows.Scan(&id, &name, &symbolSetName)
		if err != nil {
			return res, fmt.Errorf("scanning row failed : %v", err)
		}
		l := Lexicon{ID: id, Name: name, SymbolSetName: symbolSetName}
		res = append(res, l)
	}
	err = rows.Err()
	return res, err
}

func GetLexicons(db *sql.DB, names []string) ([]Lexicon, error) {
	var res []Lexicon
	if 0 == len(names) {
		return res, nil
	}

	var id int64
	var lname string
	var symbolsetname string

	rows, err := db.Query("select id, name, symbolsetname from lexicon where name in "+nQs(len(names)), convS(names)...)
	defer rows.Close()
	if err != nil {
		return res, fmt.Errorf("failed db select on lexicon table : %v", err)
	}
	for rows.Next() {
		err := rows.Scan(&id, &lname, &symbolsetname)
		if err != nil {
			return res, fmt.Errorf("failed rows scan : %v", err)
		}
		res = append(res, Lexicon{ID: id, Name: lname, SymbolSetName: symbolsetname})
	}
	err = rows.Err()
	rows.Close()

	return res, err
}

func GetLexicon(db *sql.DB, name string) (Lexicon, error) {
	var id int64
	var lname string
	var symbolsetname string

	if "" == strings.TrimSpace(name) {
		return Lexicon{}, errors.New("lexicon name must not be the empty string")
	}

	err := db.QueryRow("select id, name, symbolsetname from lexicon where name = ?", strings.ToLower(name)).Scan(&id, &lname, &symbolsetname)
	if err != nil {
		//log.Fatalf("DISASTER: %s", err)
		return Lexicon{}, fmt.Errorf("db query failed : %v", err)
	}

	return Lexicon{ID: id, Name: lname, SymbolSetName: symbolsetname}, nil
}

// LexiconFromID returns a Lexicon struct corresponding to the db row
// with that ID
func LexiconFromID(tx *sql.Tx, id int64) (Lexicon, error) {
	res := Lexicon{}
	var dbID int64
	var name, symbolSetName string
	err := tx.QueryRow("select id, name, symbolsetname from lexicon where id = ?", id).Scan(&dbID, &name, &symbolSetName)
	if err == sql.ErrNoRows {
		return res, fmt.Errorf("no lexicon with id %d : %v", id, err)
	}
	if err != nil {
		return res, fmt.Errorf("query failed %v", err)
	}

	res.ID = dbID
	res.Name = name
	res.SymbolSetName = symbolSetName

	return res, err
}

// DeleteLexicon deletes the lexicon name from the lexicon
// table. Notice that it does not remove the associated entries.
func DeleteLexicon(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return err
	}
	return DeleteLexiconTx(tx, id)
}

// DeleteLexiconTx deletes the lexicon name from the lexicon
// table. Notice that it does not remove the associated entries.
func DeleteLexiconTx(tx *sql.Tx, id int64) error {
	var n int
	err := tx.QueryRow("select count(*) from entry where entry.lexiconid = ?", id).Scan(&n)
	// must always return a row, no need to check for empty row
	if err != nil {
		return err
	}

	if n > 0 {
		return fmt.Errorf("delete all its entries before deleting a lexicon (number of entries: " + strconv.Itoa(n) + ")")
	}

	_, err = tx.Exec("delete from lexicon where id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete lexicon : %v", err)
	}

	return nil
}

func InsertOrUpdateLexicon(db *sql.DB, l Lexicon) (Lexicon, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return Lexicon{}, fmt.Errorf("failed begin transaction %v", err)
	}
	return insertOrUpdateLexiconTx(tx, l)
}
func insertOrUpdateLexiconTx(tx *sql.Tx, l Lexicon) (Lexicon, error) {

	res := Lexicon{}

	if l.ID == 0 {
		return InsertLexiconTx(tx, l)
	}

	if l.ID > 0 {
		res, err := LexiconFromID(tx, l.ID)
		if err != nil {
			return res, fmt.Errorf("faild get lexicon : %v", err)
		}
		if l != res {
			_, err := tx.Exec("update lexicon set name = ?, symbolsetname = ? where id = ?", strings.ToLower(l.Name), l.SymbolSetName, res.ID)
			if err != nil {
				tx.Rollback()
				return res, fmt.Errorf("failed to update lex : %v", err)
			}
		}
	}
	return res, nil
}

// InsertLexicon saves the name of a new lexicon to the db.
// TODO change input arg to sql.Tx ?
func InsertLexicon(db *sql.DB, l Lexicon) (Lexicon, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	res, err := InsertLexiconTx(tx, l)
	tx.Commit()

	return res, err
}

// InsertLexiconTx saves the name of a new lexicon to the db.
func InsertLexiconTx(tx *sql.Tx, l Lexicon) (Lexicon, error) {

	res, err := tx.Exec("insert into lexicon (name, symbolsetname) values (?, ?)", strings.ToLower(l.Name), l.SymbolSetName)
	if err != nil {
		tx.Rollback()
		return l, fmt.Errorf("failed to insert lexicon name + symbolset name : %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return l, fmt.Errorf("failed to get last insert id : %v", err)
	}

	//tx.Commit()

	return Lexicon{ID: id, Name: strings.ToLower(l.Name), SymbolSetName: l.SymbolSetName}, err
}

// InsertEntries saves a list of Entries and associates them to Lexicon
// TODO change input arg to sql.Tx
func InsertEntries(db *sql.DB, l Lexicon, es []Entry) ([]int64, error) {

	// TODO move to function
	var entrySTMT = "insert into entry (lexiconid, strn, language, partofspeech, wordparts) values (?, ?, ?, ?, ?)"
	var transAfterEntrySTMT = "insert into transcription (entryid, strn, language) values (?, ?, ?)"

	var ids []int64
	// Transaction -->
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return ids, fmt.Errorf("begin transaction faild : %v", err)
	}

	stmt1, err := tx.Prepare(entrySTMT)
	if err != nil {
		return ids, fmt.Errorf("failed prepare : %v", err)
	}
	stmt2, err := tx.Prepare(transAfterEntrySTMT)
	if err != nil {
		return ids, fmt.Errorf("failed prepare : %v", err)
	}

	for _, e := range es {
		res, err := tx.Stmt(stmt1).Exec(
			l.ID,
			strings.ToLower(e.Strn),
			e.Language,
			e.PartOfSpeech,
			e.WordParts)
		if err != nil {
			tx.Rollback()
			return ids, fmt.Errorf("failed exec : %v", err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return ids, fmt.Errorf("failed last insert id : %v", err)
		}
		// We want the Entry to have the right id for inserting lemma assocs below
		e.ID = id

		ids = append(ids, id)

		// res.Close()

		for _, t := range e.Transcriptions {
			_, err := tx.Stmt(stmt2).Exec(id, t.Strn, t.Language)
			if err != nil {
				tx.Rollback()
				return ids, fmt.Errorf("failed exec : %v", err)
			}
		}

		//log.Printf("%v", e)
		if "" != e.Lemma.Strn && "" != e.Lemma.Reading {
			lemma, err := SetOrGetLemma(tx, e.Lemma.Strn, e.Lemma.Reading, e.Lemma.Paradigm)
			if err != nil {
				tx.Rollback()
				return ids, fmt.Errorf("failed set or get lemma : %v", err)
			}
			err = AssociateLemma2Entry(tx, lemma, e)
			if err != nil {
				tx.Rollback()
				return ids, fmt.Errorf("Failed lemma to entry assoc: %v", err)
			}
		}
	}

	tx.Commit()
	// <- transaction

	return ids, err
}

// AssociateLemma2Entry adds a Lemma to an Entry via a linking table
func AssociateLemma2Entry(db *sql.Tx, l Lemma, e Entry) error {
	sql := "insert into Lemma2Entry (lemmaId, entryId) values (?, ?)"
	_, err := db.Exec(sql, l.ID, e.ID)
	if err != nil {
		err = fmt.Errorf("failed to associate lemma "+l.Strn+" and entry "+e.Strn+":%v", err)
	}
	return err
}

// TODO do we need both InsertLemma and SetOrGetLemma?
// InsertLemma saves a Lemma to the db, but does not associate it with an Entry
func InsertLemma(tx *sql.Tx, l Lemma) (Lemma, error) {
	sql := "insert into lemma (strn, reading, paradigm) values (?, ?, ?)"
	res, err := tx.Exec(sql, l.Strn, l.Reading, l.Paradigm)
	if err != nil {
		err = fmt.Errorf("failed insert lemma "+l.Strn+": %v", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("failed last LastInsertId after insert lemma "+l.Strn+": %v", err)
	}
	l.ID = id
	return l, err
}

// TODO do we need both InsertLemma and SetOrGetLemma?
// SetOrGetLemma saves a new Lemma to the db, or returns a matching already existing one
func SetOrGetLemma(tx *sql.Tx, strn string, reading string, paradigm string) (Lemma, error) {
	res := Lemma{}

	var id int64
	var strn0, reading0, paradigm0 string
	sqlS := "select id, strn, reading, paradigm from lemma where strn = ? and reading = ?"
	err := tx.QueryRow(sqlS, strn, reading).Scan(&id, &strn0, &reading0, &paradigm0)
	switch {
	case err == sql.ErrNoRows:
		return InsertLemma(tx, Lemma{ID: id, Strn: strn, Reading: reading, Paradigm: paradigm})
	case err != nil:
		return res, fmt.Errorf("SetOrGetLemma failed querying db : %v", err)
	}

	res.ID = id
	res.Strn = strn0
	res.Reading = reading0
	res.Paradigm = paradigm0

	return res, err
}

// TODO return error
func getLemmaFromEntryIDTx(tx *sql.Tx, id int64) (Lemma, error) {
	res := Lemma{}
	sqlS := "select lemma.id, lemma.strn, lemma.reading, lemma.paradigm from entry, lemma, lemma2entry where " +
		"entry.id = ? and entry.id = lemma2entry.entryid and lemma.id = lemma2entry.lemmaid"
	var lID int64
	var strn, reading, paradigm string
	err := tx.QueryRow(sqlS, id).Scan(&lID, &strn, &reading, &paradigm)
	switch {
	case err == sql.ErrNoRows:
		// TODO No row:
		// Silently return empty Lemma below
	case err != nil:
		//ff("getLemmaFromENtryId: %v", err)
		return res, fmt.Errorf("QueryRow failure : %v", err)
	}

	// TODO Now silently returns empty lemma if nothing returned from db. Ok?
	// Return err when there is an err
	res.ID = lID
	res.Strn = strn
	res.Reading = reading
	res.Paradigm = paradigm

	return res, nil
}

func entryMapToEntrySlice(em map[string][]Entry) []Entry {
	var res []Entry
	for _, v := range em {
		res = append(res, v...)
	}
	return res
}

type EntryWriter interface {
	Write(Entry) error
}

// EntriesFileWriter outputs formated entries to an io.Writer.
// Exmaple usage:
//	bf := bufio.NewWriter(f)
//	defer bf.Flush()
//	bfx := dbapi.EntriesFileWriter{bf}
//	dbapi.LookUp(db, q, bfx)
type EntryFileWriter struct {
	Writer io.Writer
}

func (w EntryFileWriter) Write(e Entry) error {
	// TODO call to line formatting of Entry
	_, err := fmt.Fprintf(w.Writer, "%v\n", e)
	return err
}

type EntrySliceWriter struct {
	Entries []Entry
}

func (w *EntrySliceWriter) Write(e Entry) error {
	w.Entries = append(w.Entries, e)
	return nil // fmt.Errorf("not implemented")
}

func LookUp(db *sql.DB, q Query, out EntryWriter) error {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to initialize transaction : %v", err)
	}
	return LookUpTx(tx, q, out)
}

func LookUpTx(tx *sql.Tx, q Query, out EntryWriter) error {

	sqlString, values := SelectEntriesSQL(q)
	rows, err := tx.Query(sqlString, values...)
	if err != nil {
		tx.Rollback() // nothing to rollback here, but may have been called from withing another transaction
		return err
	}

	var lexiconID, entryID int64
	var entryStrn, entryLanguage, partOfSpeech, wordParts string
	var transcriptionID, transcriptionEntryID int64
	var transcriptionStrn, transcriptionLanguage string
	var lemmaID sql.NullInt64
	var lemmaStrn, lemmaReading, lemmaParadigm sql.NullString

	var currE Entry
	var lastE int64
	lastE = -1
	for rows.Next() {
		rows.Scan(
			&lexiconID,
			&entryID,
			&entryStrn,
			&entryLanguage,
			&partOfSpeech,
			&wordParts,

			&transcriptionID,
			&transcriptionEntryID,
			&transcriptionStrn,
			&transcriptionLanguage,

			&lemmaID,
			&lemmaStrn,
			&lemmaReading,
			&lemmaParadigm,
		)
		// new entry starts here.
		// print last entry.
		// all rows with same entryID belongs to the same entry.
		// rows ordered by entryID
		if lastE != entryID {
			if lastE != -1 {
				out.Write(currE)
			}
			currE = Entry{
				LexiconID:    lexiconID,
				ID:           entryID,
				Strn:         entryStrn,
				Language:     entryLanguage,
				PartOfSpeech: partOfSpeech,
				WordParts:    wordParts,
			}
			if lemmaStrn.Valid {
				l := Lemma{Strn: lemmaStrn.String}
				if lemmaID.Valid {
					l.ID = lemmaID.Int64
				}
				if lemmaReading.Valid {
					l.Reading = lemmaReading.String
				}
				if lemmaParadigm.Valid {
					l.Paradigm = lemmaParadigm.String
				}
				currE.Lemma = l
			}
		}
		// transcriptions ordered by id so they will be added
		// in correct order
		currT := Transcription{
			ID:       transcriptionID,
			EntryID:  transcriptionEntryID,
			Strn:     transcriptionStrn,
			Language: transcriptionLanguage,
		}

		currE.Transcriptions = append(currE.Transcriptions, currT)

		lastE = entryID
	}

	// mustn't forget last entry, or lexicon will shrink by one
	// entry for each export/import...
	//	fmt.Fprintf(out, "%v\n", currE)

	out.Write(currE)

	if rows.Err() != nil {
		tx.Rollback() // nothing to rollback here, but may have been called from withing another transaction
		return rows.Err()
	}
	return nil
}

func LookUpIntoSlice(db *sql.DB, q Query) ([]Entry, error) {
	var esw EntrySliceWriter
	err := LookUp(db, q, &esw)
	if err != nil {
		return esw.Entries, fmt.Errorf("failed lookup : %v", err)
	}
	return esw.Entries, nil
}
func LookUpIntoMap(db *sql.DB, q Query) (map[string][]Entry, error) {
	res := make(map[string][]Entry)
	var esw EntrySliceWriter
	err := LookUp(db, q, &esw)
	if err != nil {
		return res, fmt.Errorf("failed lookup : %v", err)
	}
	for _, e := range esw.Entries {
		es := res[e.Strn]
		es = append(es, e)
		res[e.Strn] = es
	}
	return res, err
}

func GetEntryFromID(db *sql.DB, id int64) (Entry, error) {
	res := Entry{}
	q := Query{EntryIDs: []int64{id}}
	esw := EntrySliceWriter{}
	err := LookUp(db, q, &esw)
	if err != nil {
		return res, fmt.Errorf("LookUp failed : %v", err)
	}

	if len(esw.Entries) == 0 {
		return res, fmt.Errorf("no entry found with id %d", id)
	}
	if len(esw.Entries) > 1 {
		return res, fmt.Errorf("LookUp resulted in more than one entry")
	}
	return esw.Entries[0], nil

}

// UpdateEntry wraps call to UpdateEntryTx with a transaction
func UpdateEntry(db *sql.DB, e Entry) (updated bool, err error) { // TODO return the updated entry?
	tx, err := db.Begin()
	if err != nil {
		return updated, fmt.Errorf("failed updating entry : %v", err)
	}
	defer tx.Commit()
	return UpdateEntryTx(tx, e)
}

// UpdateEntryTx updates the fields of an Entry that do not match the
// corresponding values in the db
func UpdateEntryTx(tx *sql.Tx, e Entry) (updated bool, err error) { // TODO return the updated entry?
	// updated == false
	//dbEntryMap := //GetEntriesFromIDsTx(tx, []int64{(e.ID)})
	var esw EntrySliceWriter
	err = LookUpTx(tx, Query{EntryIDs: []int64{e.ID}}, &esw) //entryMapToEntrySlice(dbEntryMap)
	dbEntries := esw.Entries
	if len(dbEntries) == 0 {

		return updated, fmt.Errorf("no entry with id '%d'", e.ID)
	}
	if len(dbEntries) > 1 {

		return updated, fmt.Errorf("very bad error, more than one entry with id '%d'", e.ID)
	}

	updated1, err := updateTranscriptions(tx, e, dbEntries[0])
	if err != nil {
		return updated1, err
	}
	updated2, err := updateLemma(tx, e, dbEntries[0])
	if err != nil {
		return updated2, err
	}

	updated3, err := updateWordParts(tx, e, dbEntries[0])
	if err != nil {
		return updated3, err
	}
	updated4, err := updateLanguage(tx, e, dbEntries[0])
	if err != nil {
		return updated4, err
	}

	return updated1 || updated2 || updated3 || updated4, err
}

func getTIDs(ts []Transcription) []int64 {
	var res []int64
	for _, t := range ts {
		res = append(res, t.ID)
	}
	return res
}

func equal(ts1 []Transcription, ts2 []Transcription) bool {
	if len(ts1) != len(ts2) {
		return false
	}
	for i := range ts1 {
		if ts1[i] != ts2[i] {
			return false
		}
	}

	return true
}

func updateLanguage(tx *sql.Tx, e Entry, dbE Entry) (bool, error) {
	if e.ID != dbE.ID {
		return false, fmt.Errorf("new and old entries have different ids")
	}
	if e.Language == dbE.Language {
		return false, nil
	}
	_, err := tx.Exec("update entry set language = ? where entry.id = ?", e.Language, e.ID)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed language update : %v", err)
	}
	return true, nil
}

func updateWordParts(tx *sql.Tx, e Entry, dbE Entry) (bool, error) {
	if e.ID != dbE.ID {
		return false, fmt.Errorf("new and old entries have different ids")
	}
	if e.WordParts == dbE.WordParts {
		return false, nil
	}
	_, err := tx.Exec("update entry set wordparts = ? where entry.id = ?", e.WordParts, e.ID)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed worparts update : %v", err)
	}
	return true, nil
}

func updateLemma(tx *sql.Tx, e Entry, dbE Entry) (updated bool, err error) {
	if e.Lemma == dbE.Lemma {
		return false, nil
	}
	// If e.Lemma uninitialized, and different from dbE, then wipe
	// old lemma from db
	if e.Lemma.ID == 0 && e.Lemma.Strn == "" {
		_, err = tx.Exec("delete from lemma where lemma.id = ?", dbE.Lemma.ID)
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed to delete old lemma : %v", err)
		}
	}
	// Only one alternative left, to update old lemma with new values
	_, err = tx.Exec("update lemma set strn = ?, reading = ?, paradigm = ? where lemma.id = ?", e.Lemma.Strn, e.Lemma.Reading, e.Lemma.Paradigm, dbE.Lemma.ID)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to update lemma : %v", err)
	}
	return true, nil
}

// TODO move to function
var transSTMT = "insert into transcription (entryid, strn, language) values (?, ?, ?)"

func updateTranscriptions(tx *sql.Tx, e Entry, dbE Entry) (updated bool, err error) {
	if e.ID != dbE.ID {
		return false, fmt.Errorf("update and db entry id differ")
	}

	// the easy way would be to simply nuke any transcriptions for
	// the entry and substitute for the new ones, but we only want
	// to save new transcriptions if there are changes

	// If the new and old transcriptions differ, remove the old
	// and inser the new ones
	if !equal(e.Transcriptions, dbE.Transcriptions) {
		transIDs := getTIDs(dbE.Transcriptions)
		// TODO move to a function
		_, err := tx.Exec("delete from transcription where transcription.id in "+nQs(len(transIDs)), convI(transIDs)...)
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed transcription delete : %v", err)
		}
		for _, t := range e.Transcriptions {
			_, err := tx.Exec(transSTMT, e.ID, t.Strn, t.Language)
			if err != nil {
				tx.Rollback()
				return false, fmt.Errorf("failed transcription update : %v", err)
			}
		}
		// different sets of transcription, new ones inserted
		return true, nil
	}
	// Nothing happened
	return false, err
}

func unique(ns []int64) []int64 {
	tmpMap := make(map[int64]int)
	var res []int64
	for _, n := range ns {
		if _, ok := tmpMap[n]; !ok {
			res = append(res, n)
			tmpMap[n]++
		}
	}
	return res
}
func uniqIDs(ss []Symbol) []int64 {
	res := make([]int64, len(ss))
	for i, s := range ss {
		res[i] = s.LexiconID
	}
	return unique(res)
}

// SaveSymbolSet saves list of symbols that share the same LexiconID
// to the db. Prior to saving the list, it removes all current Symbols
// of the same LexiconID
func SaveSymbolSet(db *sql.DB, symbolSet []Symbol) error {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return fmt.Errorf("failed begin db transaction : %v", err)
	}
	return SaveSymbolSetTx(tx, symbolSet)
}

// SaveSymbolSetTx saves list of symbols that share the same LexiconID
// to the db. Prior to saving the list, it removes all current Symbols
// of the same LexiconID
func SaveSymbolSetTx(tx *sql.Tx, symbolSet []Symbol) error {
	if len(symbolSet) == 0 {
		return nil //li vanilli
	}
	unqIDs := uniqIDs(symbolSet)
	if len(unqIDs) != 1 {
		tx.Rollback()
		return fmt.Errorf("cannot save set of symbols with different lexiconIDs %v : ", unqIDs)
	}

	// Nuke current symbol set for lexicon of ID id:
	id := unqIDs[0]
	_, err := tx.Exec("delete from symbolset where lexiconid = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed deleting current symbol set : %v", err)
	}

	for _, s := range symbolSet {
		// TODO prepared statement?
		_, err = tx.Exec("insert into symbolset (lexiconid, symbol, category, subcat, description, ipa) values (?, ?, ?, ?, ?, ?)",
			s.LexiconID, s.Symbol, s.Category, s.Subcat, s.Description, s.IPA)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed inserting symbol : %v", err)
		}
	}

	return nil
}

// SymbolSet returns the set of Symbols defined for a lexicon with the given db id
func SymbolSet(db *sql.DB, lexiconID int64) ([]Symbol, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return []Symbol{}, fmt.Errorf("failed to start db transaction : %v", err)
	}
	return SymbolSetTx(tx, lexiconID)
}

// SymbolSetTx returns the set of Symbols defined for a lexicon with the given db id
func SymbolSetTx(tx *sql.Tx, lexiconID int64) ([]Symbol, error) {
	var res []Symbol
	rows, err := tx.Query("select lexiconid, symbol, category, subcat, description, ipa from symbolset where lexiconid = ?", lexiconID)
	if err != nil {
		return res, fmt.Errorf("failed db query : %v", err)
	}

	var lexID int64
	var symbol, category, subcat, description, ipa string
	for rows.Next() {
		rows.Scan(&lexID, &symbol, &category, &subcat, &description, &ipa)
		s := Symbol{
			LexiconID:   lexID,
			Symbol:      symbol,
			Category:    category,
			Subcat:      subcat,
			Description: description,
			IPA:         ipa,
		}
		res = append(res, s)
	}
	if rows.Err() != nil {
		return res, fmt.Errorf("error while reading db query result : %v", rows.Err())
	}

	return res, nil
}
