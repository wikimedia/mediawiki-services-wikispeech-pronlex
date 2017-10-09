/*
Package dbapi contains code wrapped around an SQL(ite3) DB.
It is used for inserting, updating and retrieving lexical entries from
a pronunciation lexicon database. A lexical entry is represented by
the dbapi.Entry struct, that mirrors entries of the entry database
table, along with associated tables such as transcription and lemma.
*/
package dbapi

//go get github.com/mattn/go-sqlite3

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"

	// installs sqlite3 driver
	"github.com/mattn/go-sqlite3"
	"github.com/stts-se/pronlex/lex"
	//"github.com/stts-se/pronlex/validation"
)

var remem = struct {
	sync.Mutex
	re map[string]*regexp.Regexp
}{
	re: make(map[string]*regexp.Regexp),
}

var regexMem = func(re, s string) (bool, error) {

	remem.Lock()
	defer remem.Unlock()
	if r, ok := remem.re[re]; ok {
		return r.MatchString(s), nil
	}

	r, err := regexp.Compile(re)
	if err != nil {
		return false, err
	}
	remem.re[re] = r
	return r.MatchString(s), nil
}

// Sqlite3WithRegex registers an Sqlite3 driver with regexp support. (Unfortunately quite slow regexp matching)
func Sqlite3WithRegex() {
	// regex := func(re, s string) (bool, error) {
	// 	//return regexp.MatchString(re, s)
	// 	return regexp.MatchString(re, s)
	// }
	sql.Register("sqlite3_with_regexp",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				return conn.RegisterFunc("regexp", regexMem, true)
			},
		})
}

func listNamesOfTriggers(db *sql.DB) ([]string, error) {
	tx, err := db.Begin()
	if err != nil {
		return []string{}, fmt.Errorf("dbapi.listNamesOfTriggers : %v", err)
	}
	defer tx.Commit()
	return listNamesOfTriggersTx(tx)
}

func listNamesOfTriggersTx(tx *sql.Tx) ([]string, error) {
	var res []string

	q := "select name from sqlite_master where type = 'trigger'"
	rows, err := tx.Query(q)
	if err != nil {
		tx.Rollback()
		log.Printf("dbapi.listNamesOfTriggersTx : %v\n", err)
		return res, fmt.Errorf("dbapi.listNamesOfTriggers : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			tx.Rollback()
			return res, fmt.Errorf("dbapi.listNamesOfTriggers : %v", err)
		}
		res = append(res, name)
	}

	return res, nil
}

// ListEntryTableColumnNames is a meta-function that returns the names of the columns of the 'entry' lexicon database table.
// It can be used for checking that the entry table has the expected columns.
func listEntryTableColumnNames(db *sql.DB) ([]string, error) {
	q := "SELECT * FROM entry LIMIT 0"

	rows, err := db.Query(q)
	if err != nil {
		return []string{}, fmt.Errorf("ListEntryTableColumnNames : %v", err)
	}
	defer rows.Close()
	colNames, err := rows.Columns()
	if err != nil {
		return colNames, fmt.Errorf("ListEntryTableColumnNames : %v", err)
	}

	return colNames, err
}

// ListLexicons returns a list of the lexicons defined in the db
// (i.e., Lexicon structs corresponding to the rows of the lexicon
// table).
//
// TODO: Create a DB struct, and move functions of type
// funcName(db *sql.DB, ...) into methods of the new struct.
func listLexicons(db *sql.DB) ([]lexicon, error) {
	var res []lexicon
	sql := "select id, name, symbolsetname from lexicon"
	rows, err := db.Query(sql)
	if err != nil {
		return res, fmt.Errorf("db query failed : %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		l := lexicon{}
		err = rows.Scan(&l.id, &l.name, &l.symbolSetName)
		if err != nil {
			return res, fmt.Errorf("scanning row failed : %v", err)
		}
		res = append(res, l)
	}
	err = rows.Err()
	return res, err
}

// GetLexicon returns a Lexicon struct matching a lexicon name in the db.
// Returns error if no such lexicon name in db
func getLexicon(db *sql.DB, name string) (lexicon, error) {
	tx, err := db.Begin()
	if err != nil {
		return lexicon{}, fmt.Errorf("failed to create transaction : %v", err)
	}
	defer tx.Commit()
	return getLexiconTx(tx, name)
}

func getLexiconMapTx(tx *sql.Tx) (map[string]bool, error) {
	res := make(map[string]bool)

	rows, err := tx.Query("select name from lexicon")
	if err != nil {
		return res, fmt.Errorf("failed db select on lexicon table : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		res[name] = true
		if err != nil {
			return res, fmt.Errorf("failed db select on lexicon table : %v", err)
		}
	}
	return res, err

}

func getLexiconTx(tx *sql.Tx, name string) (lexicon, error) {
	res := lexicon{}
	name0 := strings.ToLower(name)
	var err error

	row := tx.QueryRow("select id, name, symbolsetname from lexicon where name = ? ", name0).Scan(&res.id, &res.name, &res.symbolSetName)
	if row == sql.ErrNoRows {
		return res, fmt.Errorf("couldn't find lexicon '%s'", name)
	}

	return res, err

}

// GetLexicons takes a list of lexicon names and returns a list of
// Lexicon structs corresponding to rows of db lexicon table with those name fields.
func getLexicons(db *sql.DB, names []string) ([]lexicon, error) {
	var res []lexicon
	found := make(map[string]bool)
	if 0 == len(names) {
		return res, nil
	}

	rows, err := db.Query("select id, name, symbolsetname from lexicon where name in "+nQs(len(names)), convS(names)...)
	if err != nil {
		return res, fmt.Errorf("failed db select on lexicon table : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		l := lexicon{}
		err := rows.Scan(&l.id, &l.name, &l.symbolSetName)
		if err != nil {
			return res, fmt.Errorf("failed rows scan : %v", err)
		}
		found[strings.ToLower(l.name)] = true
		res = append(res, l)
	}

	err = rows.Err()
	//rows.Close()

	if len(res) != len(names) {
		var missing []string
		for _, n := range names {
			if _, ok := found[strings.ToLower(n)]; !ok {
				missing = append(missing, n)
			}
		}

		err0 := fmt.Errorf("unknown lexicon(s): %v", strings.Join(missing, ", "))
		if err != nil {
			err = fmt.Errorf("%v : %v", err, err0)
		} else {
			err = err0
		}
	}

	return res, err
}

// LexiconFromID returns a Lexicon struct corresponding to a row in
// the lexicon table with the given id
func lexiconFromID(db *sql.DB, id int64) (lexicon, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return lexicon{}, fmt.Errorf("lexiconFromID failed to start db transaction : %v", err)
	}

	return lexiconFromIDTx(tx, id)
}

// LexiconFromIDTx returns a Lexicon struct corresponding to a row in
// the lexicon table with the given id
func lexiconFromIDTx(tx *sql.Tx, id int64) (lexicon, error) {
	res := lexicon{}
	err := tx.QueryRow("select id, name, symbolsetname from lexicon where id = ?", id).Scan(&res.id, &res.name, &res.symbolSetName)
	if err == sql.ErrNoRows {
		return res, fmt.Errorf("no lexicon with id %d : %v", id, err)
	}
	if err != nil {
		return res, fmt.Errorf("query failed %v", err)
	}

	return res, err
}

// DeleteLexicon deletes the lexicon name from the lexicon
// table. Notice that it does not remove the associated entries.
// It should be impossible to delete the Lexicon table entry if associated to any entries.
func deleteLexicon(db *sql.DB, lexName string) error {
	log.Printf("deleteLexicon called with lexicon name %s\n", lexName)
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return err
	}
	return deleteLexiconTx(tx, lexName)
}

// DeleteLexiconTx deletes the lexicon name from the lexicon
// table. Notice that it does not remove the associated entries.
// It should be impossible to delete the Lexicon table entry if associated to any entries.
func deleteLexiconTx(tx *sql.Tx, lexName string) error {
	// does it exist?
	lexExists, err := lexiconExists(tx, lexName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("dbapi.DeleteLexiconTx : failed to lookup lexicon from name %s : %v", lexName, err)
	}
	if !lexExists {
		return fmt.Errorf("dbapi.DeleteLexiconTx : no lexicon exists with name : %s", lexName)
	}

	n := 0
	err = tx.QueryRow("select count(*) from entry, lexicon where lexicon.name = ? and entry.lexiconid = lexicon.id", lexName).Scan(&n)
	// must always return a row, no need to check for empty row
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("HEJ DIN FAN")
			return fmt.Errorf("The was no lexicon with name %s : %v", lexName, err)
		}
		return err
	}
	if n > 0 {
		return fmt.Errorf("delete all its entries before deleting a lexicon (number of entries: " + strconv.Itoa(n) + ")")
	}

	_, err = tx.Exec("delete from lexicon where name = ?", lexName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete lexicon : %v", err)
	}
	return nil
}

func lexiconExists(tx *sql.Tx, lexName string) (bool, error) {

	var id int64
	err := tx.QueryRow("SELECT id FROM lexicon WHERE name = ?", lexName).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		return false, fmt.Errorf("dbapi.lexiconExists db query failed : %v", err)
	default:
		return true, nil
	}
}

// SuperDeleteLexicon deletes the lexicon name from the lexicon
// table and also whipes all associated entries out of existence.
// TODO Send progress message to client over websocket
func superDeleteLexicon(db *sql.DB, lexName string) error {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return fmt.Errorf("superDeleteLexicon failed to initiate transaction : %v", err)
	}
	return superDeleteLexiconTx(tx, lexName)
}

// SuperDeleteLexiconTx deletes the lexicon name from the lexicon
// table and also whipes all associated entries out of existence.
func superDeleteLexiconTx(tx *sql.Tx, lexName string) error {

	log.Println("dbapi.superDeleteLexiconTX was called")

	// does it exist?
	lexExists, err := lexiconExists(tx, lexName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("dbapi.SuperDeleteLexiconTx : failed to lookup lexicon from name %s : %v", lexName, err)
	}
	if !lexExists {
		return fmt.Errorf("dbapi.SuperDeleteLexiconTx : no lexicon exists with name : %s", lexName)
	}

	// delete entries
	_, err = tx.Exec("DELETE FROM entry WHERE lexiconid IN (SELECT id FROM lexicon WHERE name = ?)", lexName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("dbapi.SuperDeleteLexiconTx : failed to delete entries : %v", err)
	}
	log.Println("dbapi.superDeleteLexiconTX finished deleting from entry set")

	// delete lexicon
	_, err = tx.Exec("DELETE FROM lexicon WHERE name = ?", lexName)

	if err != nil {
		tx.Rollback()
		return fmt.Errorf("dbapi.SuperDeleteLexiconTx : failed to delete lexicon : %v", err)
	}

	log.Println("dbapi.superDeleteLexiconTX finished deleting from lexicon set")

	log.Printf("Deleted lexicon named %s\n", lexName)

	return nil
}

// DefineLexicon saves the name of a new lexicon to the db.
func defineLexicon(db *sql.DB, l lexicon) (lexicon, error) {
	tx, err := db.Begin()
	if err != nil {
		return lexicon{}, fmt.Errorf("failed to get db transaction : %v", err)
	}
	defer tx.Commit()
	res, err := defineLexiconTx(tx, l)
	//tx.Commit()

	return res, err
}

// DefineLexiconTx saves the name of a new lexicon to the db.
func defineLexiconTx(tx *sql.Tx, l lexicon) (lexicon, error) {

	res, err := tx.Exec("insert into lexicon (name, symbolsetname) values (?, ?)", strings.ToLower(l.name), l.symbolSetName)
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

	return lexicon{id: id, name: strings.ToLower(l.name), symbolSetName: l.symbolSetName}, err
}

// MoveResult is returned from the MoveNewEntries function.
// TODO Since it only contains a single int64, this struct is probably not needed.
// Only useful if more info is to be returned.
type MoveResult struct {
	N int64
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
func moveNewEntries(db *sql.DB, fromLexicon, toLexicon, newSource, newStatus string) (MoveResult, error) {
	if strings.TrimSpace(newSource) == "" {
		msg := "MoveNewEntries called with the empty 'newSource' argument"
		return MoveResult{}, fmt.Errorf(msg)
	}
	if strings.TrimSpace(newStatus) == "" {
		msg := "MoveNewEntries called with the empty 'newStatus' argument"
		return MoveResult{}, fmt.Errorf(msg)
	}

	tx, err := db.Begin()
	if err != nil {
		return MoveResult{}, fmt.Errorf("failed to get db transaction : %v", err)
	}
	defer tx.Commit()

	return moveNewEntriesTx(tx, fromLexicon, toLexicon, newSource, newStatus)
}

// moveNewEntriesTx is documented under MoveNewEntries
func moveNewEntriesTx(tx *sql.Tx, fromLexicon, toLexicon, newSource, newStatus string) (MoveResult, error) {
	if strings.TrimSpace(newSource) == "" {
		msg := "moveNewEntriesTx called with the empty 'newSource' argument"
		tx.Rollback()
		return MoveResult{}, fmt.Errorf(msg)
	}
	if strings.TrimSpace(newStatus) == "" {
		msg := "moveNewEntriesTx called with the empty 'newStatus' argument"
		tx.Rollback()
		return MoveResult{}, fmt.Errorf(msg)
	}

	res := MoveResult{}
	var err error
	fromLex, err := getLexiconTx(tx, fromLexicon)
	if err != nil {
		err := fmt.Errorf("couldn't find lexicon %s : %v", fromLexicon, err)
		tx.Rollback()
		return res, err
	}
	toLex, err := getLexiconTx(tx, toLexicon)
	if err != nil {
		err := fmt.Errorf("couldn't find lexicon %s : %v", toLexicon, err)
		tx.Rollback()
		return res, err
	}

	//TODO This may be computationally expensive for large lexica?
	// query :=
	// 	`SELECT a.id FROM entry a WHERE a.lexiconid = ? AND  NOT EXISTS(SELECT strn FROM entry WHERE lexiconid = ? AND strn = a.strn)`

	// movable, err := tx.Query(query, fromLex.ID, toLex.ID)
	// if err != nil {
	// 	tx.Rollback()
	// 	return res, fmt.Errorf("select query failed : %v", err)
	// }
	// defer movable.Close()

	// var movableIDs []int64
	// for movable.Next() {
	// 	var id int64
	// 	if err := movable.Scan(&id); err != nil {
	// 		return res, fmt.Errorf("failed retrieving query result : %v", err)
	// 	}
	// 	movableIDs = append(movableIDs, id)
	// }

	// if err0 := movable.Err(); err0 != nil {
	// 	newErr := fmt.Errorf("MoveNewEntriesTx suffered a db mishap : %v", err0)
	// 	if err == nil {
	// 		err = newErr
	// 	} else {

	// 		err = fmt.Errorf("%v : %v", err, newErr)
	// 	}
	// }

	// if len(movableIDs) == 0 {
	// 	return res, err
	// }

	// var idsStrn []string
	// for _, n := range movableIDs {
	// 	idsStrn = append(idsStrn, strconv.Itoa(int(n)))
	// }

	where := `WHERE entry.id IN (SELECT a.id FROM entry a WHERE a.lexiconid = ?
                       AND NOT EXISTS(SELECT strn FROM entry WHERE lexiconid = ? AND strn = a.strn))`

	insertQuery := `INSERT INTO entrystatus (name, source, entryid, current) SELECT ?, ?, entry.id, '1' FROM entry ` + where

	// updateQuery0 := `UPDATE entrystatus SET current = 1 AND source = ? AND name = ? ` + where + ` AND entrystatus.entryid = entry.id`
	q0Rez, err := tx.Exec(insertQuery, newStatus, newSource, fromLex.id, toLex.id)
	if err != nil {
		tx.Rollback()
		return res, fmt.Errorf("failed to update entrystatus : %v", err)
	}

	//_ = q0Rez

	updateQuery := `UPDATE entry SET lexiconid = ? ` + where

	//log.Printf("Q: %s\n", updateQuery)

	qRez, err := tx.Exec(updateQuery, toLex.id, fromLex.id, toLex.id)
	if err != nil {
		tx.Rollback()
		return res, fmt.Errorf("failed to update lexiconids : %v", err)
	}

	//TODO Should this result in an error and rollback?
	q0N, _ := q0Rez.RowsAffected()
	qN, _ := qRez.RowsAffected()
	if q0N != qN {
		log.Printf("dbapi.moveNewEntriesTx: UPDATE and INSERT queries affected different number of rows: %v and %v", q0N, qN)
	}

	if n, err := qRez.RowsAffected(); err == nil {
		res.N = n
	}
	return res, err
}

// TODO move to function?
var entrySTMT = "insert into entry (lexiconid, strn, language, partofspeech, morphology, wordparts, preferred) values (?, ?, ?, ?, ?, ?, ?)"
var transAfterEntrySTMT = "insert into transcription (entryid, strn, language, sources) values (?, ?, ?, ?)"

//var statusSetCurrentFalse = "UPDATE entrystatus SET current = 0 WHERE entrystatus.entryid = ?"
var insertStatus = "INSERT INTO entrystatus (entryid, name, source) values (?, ?, ?)"

// InsertEntries saves a list of Entries and associates them to Lexicon
// TODO: Change second input argument to string (lexicon name) instead of Lexicon struct.
// TODO change input arg to sql.Tx
func insertEntries(db *sql.DB, l lexicon, es []lex.Entry) ([]int64, error) {

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
		// convert 'Preferred' into DB integer value
		var pref int64
		if e.Preferred {
			pref = 1
		}
		res, err := tx.Stmt(stmt1).Exec(
			l.id,
			strings.ToLower(e.Strn),
			e.Language,
			e.PartOfSpeech,
			e.Morphology,
			e.WordParts,
			pref)
		if err != nil {
			tx.Rollback()
			return ids, fmt.Errorf("failed exec : %v", err)
		}

		id, err := res.LastInsertId()
		if err != nil {
			tx.Rollback()
			return ids, fmt.Errorf("failed last insert id : %v", err)
		}
		// We want thelex.Entry to have the right id for inserting lemma assocs below
		e.ID = id

		ids = append(ids, id)

		// res.Close()

		for _, t := range e.Transcriptions {
			_, err := tx.Stmt(stmt2).Exec(id, t.Strn, t.Language, t.SourcesString())
			if err != nil {
				tx.Rollback()
				return ids, fmt.Errorf("failed exec : %v", err)
			}
		}

		//log.Printf("%v", e)
		if "" != e.Lemma.Strn { // && "" != e.Lemma.Reading {
			lemma, err := setOrGetLemma(tx, e.Lemma.Strn, e.Lemma.Reading, e.Lemma.Paradigm)
			if err != nil {
				tx.Rollback()
				return ids, fmt.Errorf("failed set or get lemma : %v", err)
			}
			err = associateLemma2Entry(tx, lemma, e)
			if err != nil {
				tx.Rollback()
				return ids, fmt.Errorf("Failed lemma to entry assoc: %v", err)
			}
		}

		if "" != e.Tag {
			err = insertEntryTagTx(tx, e.ID, e.Tag)
			if err != nil {
				tx.Rollback()
				return ids, fmt.Errorf("Failed to insert entry tag '%s' for '%s': %v", e.Tag, e.Strn, err)
			}
		}

		if trm(e.EntryStatus.Name) != "" {
			//var statusSetCurrentFalse = "UPDATE entrystatus SET current = 0 WHERE entrystatus.entryid = ?"
			//var insertStatus = "INSERT INTO entrystatus (entryid, name, source) values (?, ?, ?)"

			// _, err := tx.Exec(statusSetCurrentFalse, e.ID)
			// if err != nil {
			// 	tx.Rollback()
			// 	return ids, fmt.Errorf("updating lex.EntryStatus.Current failed : %v", err)
			// }
			_, err = tx.Exec(insertStatus, e.ID, strings.ToLower(e.EntryStatus.Name), strings.ToLower(e.EntryStatus.Source)) //, e.EntryStatus.Current) // TODO?
			if err != nil {
				tx.Rollback()
				return ids, fmt.Errorf("inserting EntryStatus failed : %v", err)
			}
		}

		err = insertEntryValidations(tx, e, e.EntryValidations)
		if err != nil {
			tx.Rollback()
			return ids, fmt.Errorf("inserting EntryValidations failed : %v", err)
		}
	}

	tx.Commit()
	// <- transaction

	return ids, err
}

var insertEntryTag = "INSERT INTO EntryTag (entryId, tag) values (?, ?)"

//TODO Add tests

//TODO Add db look-up to see if db uniqueness constraints are
//violated, to return more gentle error (or no error) instead of
//failing and rolling back.
func insertEntryTagTx(tx *sql.Tx, entryID int64, tag string) error {

	tag = strings.TrimSpace(strings.ToLower(tag))

	var eId int64
	var eTag, wordForm string
	// Check if it is already there, then silently do nuttin'
	chkResErr := tx.QueryRow("SELECT entryId, tag, wordForm FROM EntryTag WHERE entryId = ?", entryID).Scan(&eId, &eTag, &wordForm)
	_ = chkResErr
	//sdkjjks :=

	// Entry already has wanted tag, silently accept the fact
	if eId == entryID && eTag == tag {
		return nil
	}

	// Entry had different tag, report error but do nothing
	if eTag != "" && tag != eTag {
		return fmt.Errorf("insertEntryTag: failed to insert tag '%s' because entry already had tag '%s'", tag, eTag)
	}

	// TODO Check that another entry of the same wordform has this tag

	//err == sql.ErrNoRows

	insert, err := tx.Prepare(insertEntryTag)
	if err != nil {
		// Let caller be responsible for rollback
		//tx.Rollback()
		return fmt.Errorf("failed prepare : %v", err)
	}

	_, err = tx.Stmt(insert).Exec(entryID, tag)
	if err != nil {
		// TODO Maybe no rollback?
		// Let caller be responsible for rollback
		//tx.Rollback()
		return fmt.Errorf("failed insert entry tag : %v", err)
	}

	return nil
}

// InsertLemma saves a lex.Lemma to the db, but does not associate it with anlex.Entry
// TODO do we need both InsertLemma and SetOrGetLemma?
func insertLemma(tx *sql.Tx, l lex.Lemma) (lex.Lemma, error) {
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

// SetOrGetLemma saves a new lex.Lemma to the db, or returns a matching already existing one
// TODO do we need both InsertLemma and SetOrGetLemma?
func setOrGetLemma(tx *sql.Tx, strn string, reading string, paradigm string) (lex.Lemma, error) {
	res := lex.Lemma{}

	sqlS := "select id, strn, reading, paradigm from lemma where strn = ? and reading = ?"
	err := tx.QueryRow(sqlS, strn, reading).Scan(&res.ID, &res.Strn, &res.Reading, &res.Paradigm)
	switch {
	case err == sql.ErrNoRows:
		return insertLemma(tx, lex.Lemma{Strn: strn, Reading: reading, Paradigm: paradigm})
	case err != nil:
		return res, fmt.Errorf("setOrGetLemma failed querying db : %v", err)
	}

	return res, err
}

// AssociateLemma2Entry adds a lex.Lemma to anlex.Entry via a linking table
func associateLemma2Entry(db *sql.Tx, l lex.Lemma, e lex.Entry) error {
	sql := "insert into Lemma2Entry (lemmaId, entryId) values (?, ?)"
	_, err := db.Exec(sql, l.ID, e.ID)
	if err != nil {
		err = fmt.Errorf("failed to associate lemma "+l.Strn+" and entry "+e.Strn+":%v", err)
	}
	return err
}

// NOT USED AS OF 2017-02-24
// func getLemmaFromEntryIDTx(tx *sql.Tx, id int64) (lex.Lemma, error) {
// 	res := lex.Lemma{}
// 	sqlS := "select lemma.id, lemma.strn, lemma.reading, lemma.paradigm from entry, lemma, lemma2entry where " +
// 		"entry.id = ? and entry.id = lemma2entry.entryid and lemma.id = lemma2entry.lemmaid"
// 	err := tx.QueryRow(sqlS, id).Scan(&res.ID, &res.Strn, &res.Reading, &res.Paradigm)
// 	switch {
// 	case err == sql.ErrNoRows:
// 		// TODO No row:
// 		// Silently return empty lex.Lemma below
// 	case err != nil:
// 		//ff("getLemmaFromENtryId: %v", err)
// 		return res, fmt.Errorf("QueryRow failure : %v", err)
// 	}

// 	// TODO Now silently returns empty lemma if nothing returned from db. Ok?
// 	return res, nil
// }

// LookUpIds takes a Query struct, searches the lexicon db, and writes the result to a slice of ids
func lookUpIds(db *sql.DB, lexNames []lex.LexName, q Query) ([]int64, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to initialize transaction : %v", err)
	}
	return lookUpIdsTx(tx, lexNames, q)
}

// LookUpIdsTx takes a Query struct, searches the lexicon db, and returns a slice of ids
func lookUpIdsTx(tx *sql.Tx, lexNames []lex.LexName, q Query) ([]int64, error) {
	var result []int64

	err := validateInputLexicons(tx, lexNames, q)
	if err != nil {
		return result, err
	}

	sqlStmt := selectEntryIdsSQL(lexNames, q)

	rows, err := tx.Query(sqlStmt.sql, sqlStmt.values...)
	if err != nil {
		tx.Rollback() // nothing to rollback here, but may have been called from within another transaction
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var entryID int64
		rows.Scan(
			&entryID,
		)
		result = append(result, entryID)
	}
	if rows.Err() != nil {
		tx.Rollback() // nothing to rollback here, but may have been called from within another transaction
		return nil, rows.Err()
	}

	return result, nil
}

// LookUp takes a Query struct, searches the lexicon db, and writes the result to the
//lex.EntryWriter.
func lookUp(db *sql.DB, lexNames []lex.LexName, q Query, out lex.EntryWriter) error {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to initialize transaction : %v", err)
	}
	return lookUpTx(tx, lexNames, q, out)
}

func validateInputLexicons(tx *sql.Tx, lexNames []lex.LexName, q Query) error {
	if len(lexNames) == 0 && len(q.EntryIDs) == 0 { // if entry id is specified, we can do the search without the lexicon name
		tx.Rollback()
		return fmt.Errorf("cannot perform a search without at least one lexicon specified")
	}

	lexiconMap, err := getLexiconMapTx(tx)
	if err != nil {
		tx.Rollback() // nothing to rollback here, but may have been called from within another transaction
		return err
	}
	for _, lexName := range lexNames {
		_, ok := lexiconMap[string(lexName)]
		if !ok {
			tx.Rollback()
			return fmt.Errorf("no lexicon exists with name: %s", lexName)
		}
	}
	return nil
}

// LookUpTx takes a Query struct, searches the lexicon db, and writes the result to the
// EntryWriter.
// TODO: rewrite to go through the result set before building the result. That is, save all structs corresponding to rows in the scanning run, then build the result structure (so that no identical values are duplicated: a result set may have several rows of repeated data)
func lookUpTx(tx *sql.Tx, lexNames []lex.LexName, q Query, out lex.EntryWriter) error {

	//log.Printf("QUWRY %v\n\n", q)

	sqlStmt := selectEntriesSQL(lexNames, q)

	//log.Printf("SQL %v\n\n", sqlStmt)

	//log.Printf("VALUES %v\n\n", values)

	err := validateInputLexicons(tx, lexNames, q)
	if err != nil {
		return err
	}

	rows, err := tx.Query(sqlStmt.sql, sqlStmt.values...)
	if err != nil {
		tx.Rollback() // nothing to rollback here, but may have been called from within another transaction
		return err
	}
	defer rows.Close()

	var entryID, preferred int64
	var lexiconName, entryStrn, entryLanguage, partOfSpeech, morphology, wordParts string

	var transcriptionID, transcriptionEntryID int64
	var transcriptionStrn, transcriptionLanguage, transcriptionSources string

	// Optional/nullable values

	var lemmaID sql.NullInt64
	var lemmaStrn, lemmaReading, lemmaParadigm, entryTag sql.NullString

	var entryStatusID sql.NullInt64
	var entryStatusName, entryStatusSource sql.NullString
	var entryStatusTimestamp sql.NullString //sql.NullInt64
	var entryStatusCurrent sql.NullBool

	var entryValidationID sql.NullInt64
	var entryValidationLevel, entryValidationName, entryValidationMessage, entryValidationTimestamp sql.NullString

	// transcription ids read so far, in order not to add same trans twice
	transIDs := make(map[int64]int)
	// entry validation ids read so far, in order not to add same validation twice
	valiIDs := make(map[int64]int)

	var currE lex.Entry
	var lastE int64
	lastE = -1
	for rows.Next() {
		rows.Scan(
			&lexiconName,
			&entryID,
			&entryStrn,
			&entryLanguage,
			&partOfSpeech,
			&morphology,
			&wordParts,
			&preferred,

			&transcriptionID,
			&transcriptionEntryID,
			&transcriptionStrn,
			&transcriptionLanguage,
			&transcriptionSources,

			// Optional, from LEFT JOIN

			&lemmaID,
			&lemmaStrn,
			&lemmaReading,
			&lemmaParadigm,

			&entryTag,

			&entryStatusID,
			&entryStatusName,
			&entryStatusSource,
			&entryStatusTimestamp,
			&entryStatusCurrent,

			&entryValidationID,
			&entryValidationLevel,
			&entryValidationName,
			&entryValidationMessage,
			&entryValidationTimestamp,
		)
		// new entry starts here.
		//
		// all rows with same entryID belongs to the same entry.
		// rows ordered by entryID
		var pref bool // convert 'preferred' value from DB integer value
		if preferred == 1 {
			pref = true
		}
		if lastE != entryID {
			if lastE != -1 {
				out.Write(currE)
			}
			currE = lex.Entry{
				LexRef:       lex.NewLexRef("", lexiconName), // DBRef is not set here (will be set by DBManager)
				ID:           entryID,
				Strn:         entryStrn,
				Language:     entryLanguage,
				PartOfSpeech: partOfSpeech,
				Morphology:   morphology,
				WordParts:    wordParts,
				Preferred:    pref,
				Tag:          entryTag.String,
			}
			// max one lemma per entry
			if lemmaStrn.Valid && trm(lemmaStrn.String) != "" {
				l := lex.Lemma{Strn: lemmaStrn.String}
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

			// TODO Only add once per lex.Entry as long as single status?
			// TODO probably should be a slice of statuses?
			// TODO now checks for current = true before adding to lex.Entry
			if entryStatusID.Valid && entryStatusName.Valid && trm(entryStatusName.String) != "" {
				es := lex.EntryStatus{ID: entryStatusID.Int64, Name: entryStatusName.String}
				if entryStatusSource.Valid {
					es.Source = entryStatusSource.String
				}
				if entryStatusTimestamp.Valid {
					es.Timestamp = entryStatusTimestamp.String
				}
				if entryStatusCurrent.Valid {
					es.Current = entryStatusCurrent.Bool
				}
				// only update the lex.Entry with status if current = true
				if entryStatusCurrent.Valid && entryStatusCurrent.Bool {
					currE.EntryStatus = es
				}
			}
		}
		// Things that may appear in several rows of a single lex.Entry below:

		// transcriptions ordered by id so they will be added
		// in correct order
		// Only add transcriptions that are !ok, i.e. not added already
		if _, ok := transIDs[transcriptionID]; !ok {
			currT := lex.Transcription{
				ID:       transcriptionID,
				EntryID:  transcriptionEntryID,
				Strn:     transcriptionStrn,
				Language: transcriptionLanguage,
				//Sources:  strings.Split(transcriptionSources, SourceDelimiter),
			}
			// Sources may be empty string in db
			if trm(transcriptionSources) == "" {
				currT.Sources = make([]string, 0)
			} else {
				// strings.Split returns the empty string if input the empty string
				currT.Sources = strings.Split(transcriptionSources, lex.SourceDelimiter)
			}

			currE.Transcriptions = append(currE.Transcriptions, currT)
			transIDs[transcriptionID]++
		}

		if currE.EntryValidations == nil {
			currE.EntryValidations = []lex.EntryValidation{}
		}
		// zero or more lex.EntryValidations
		if entryValidationID.Valid && entryValidationLevel.Valid && entryValidationName.Valid && entryValidationMessage.Valid && entryValidationTimestamp.Valid {
			if _, ok := valiIDs[entryValidationID.Int64]; !ok {
				currV := lex.EntryValidation{
					ID:        entryValidationID.Int64,
					Level:     entryValidationLevel.String,
					RuleName:  entryValidationName.String,
					Message:   entryValidationMessage.String,
					Timestamp: entryValidationTimestamp.String,
				}
				currE.EntryValidations = append(currE.EntryValidations, currV)
				valiIDs[entryValidationID.Int64]++
			}
		}
		lastE = entryID
	}

	// mustn't forget last entry, or lexicon will shrink by one
	// entry for each export/import...
	//	fmt.Fprintf(out, "%v\n", currE)
	// but only print last entry if there were any entries...
	if lastE > -1 {
		out.Write(currE)
	}
	if rows.Err() != nil {
		tx.Rollback() // nothing to rollback here, but may have been called from within another transaction
		return rows.Err()
	}

	return nil
}

// LookUpIntoSlice is a wrapper around LookUp, returning a slice of Entries
func lookUpIntoSlice(db *sql.DB, lexNames []lex.LexName, q Query) ([]lex.Entry, error) {
	var esw lex.EntrySliceWriter
	err := lookUp(db, lexNames, q, &esw)
	if err != nil {
		return esw.Entries, fmt.Errorf("failed lookup : %v", err)
	}
	return esw.Entries, nil
}

// LookUpIntoMap is a wrapper around LookUp, returning a map where the
// keys are word forms and the values are slices of Entries. (There may be several entries with the same Strn value.)
func lookUpIntoMap(db *sql.DB, lexNames []lex.LexName, q Query) (map[string][]lex.Entry, error) {
	res := make(map[string][]lex.Entry)
	var esw lex.EntrySliceWriter
	err := lookUp(db, lexNames, q, &esw)
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

// GetEntryFromID is a wrapper around LookUp and returns the lex.Entry corresponding to the db id
func getEntryFromID(db *sql.DB, id int64) (lex.Entry, error) {
	res := lex.Entry{}
	q := Query{EntryIDs: []int64{id}}
	esw := lex.EntrySliceWriter{}
	err := lookUp(db, []lex.LexName{}, q, &esw)
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

// GetEntriesFromIDs is a wrapper around LookUp and returns the lex.Entry corresponding to the db id
func getEntriesFromIDs(db *sql.DB, ids []int64, out lex.EntryWriter) error {
	q := Query{EntryIDs: ids}
	err := lookUp(db, []lex.LexName{}, q, out)
	if err != nil {
		return fmt.Errorf("LookUp failed : %v", err)
	}

	if out.Size() != len(ids) {
		return fmt.Errorf("got %d input ids, but found %d entries", len(ids), out.Size())
	}
	return nil

}

// UpdateEntry wraps call to UpdateEntryTx with a transaction, and returns the updated entry, fresh from the db
// TODO Consider how to handle inconsistent input entries
// TODO Full name of DB as input param?
func updateEntry(db *sql.DB, e lex.Entry) (res lex.Entry, updated bool, err error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		tx.Rollback()
		return res, updated, fmt.Errorf("failed starting transaction for updating entry : %v", err)
	}

	updated, err = updateEntryTx(tx, e)
	if err != nil {
		tx.Rollback()
		return res, updated, fmt.Errorf("failed updating entry : %v", err)
	}
	tx.Commit()
	res, err = getEntryFromID(db, e.ID)
	if err != nil {
		tx.Rollback()
		return res, updated, fmt.Errorf("failed getting updated entry : %v", err)
	}

	return res, updated, err
}

// UpdateEntryTx updates the fields of an lex.Entry that do not match the
// corresponding values in the db
func updateEntryTx(tx *sql.Tx, e lex.Entry) (updated bool, err error) { // TODO return the updated entry?
	// updated == false
	//dbEntryMap := //GetEntriesFromIDsTx(tx, []int64{(e.ID)})
	var esw lex.EntrySliceWriter
	err = lookUpTx(tx, []lex.LexName{e.LexRef.LexName}, Query{EntryIDs: []int64{e.ID}}, &esw) //entryMapToEntrySlice(dbEntryMap)
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
	updated5, err := updateEntryStatus(tx, e, dbEntries[0])
	if err != nil {
		return updated5, err
	}

	updated6, err := updateEntryValidation(tx, e, dbEntries[0])
	if err != nil {
		return updated6, err
	}

	updated7, err := updatePreferred(tx, e, dbEntries[0])
	if err != nil {
		return updated7, err
	}

	updated8, err := updateEntryTag(tx, e, dbEntries[0])
	if err != nil {
		return updated8, err
	}

	return updated1 || updated2 || updated3 || updated4 || updated5 || updated6 || updated7 || updated8, err
}

func getTIDs(ts []lex.Transcription) []int64 {
	var res []int64
	for _, t := range ts {
		res = append(res, t.ID)
	}
	return res
}

func equal(ts1 []lex.Transcription, ts2 []lex.Transcription) bool {
	if len(ts1) != len(ts2) {
		return false
	}
	for i := range ts1 {
		//if ts1[i] != ts2[i] {
		// TODO? Define Equal(Transcription) on lex.Transcription
		if !reflect.DeepEqual(ts1[i], ts2[i]) {
			return false
		}
	}

	return true
}

func updateLanguage(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error) {
	if e.ID != dbE.ID {
		tx.Rollback()
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

func updateWordParts(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error) {
	if e.ID != dbE.ID {
		tx.Rollback()
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

func updatePreferred(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error) {
	if e.ID != dbE.ID {
		tx.Rollback()
		return false, fmt.Errorf("new and old entries have different ids")
	}
	if e.Preferred == dbE.Preferred {
		return false, nil
	}

	// convert bool into DB integer
	var pref int64
	if e.Preferred {
		pref = 1
	}
	_, err := tx.Exec("update entry set preferred = ? where entry.id = ?", pref, e.ID)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed preferred update : %v", err)
	}
	return true, nil
}

func updateLemma(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (updated bool, err error) {
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

func updateEntryTag(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error) {
	if e.ID != dbE.ID {
		tx.Rollback()
		return false, fmt.Errorf("updateEntryTag: new and old entries have different ids")
	}
	if e.Tag == dbE.Tag {
		return false, nil
	}

	_, err := tx.Exec("UPDATE entrytag SET tag = ? WHERE entryid = ?", e.Tag, e.ID)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("updateEntryTag failed : %v", err)
	}

	return true, nil
}

// TODO move to function
var transSTMT = "insert into transcription (entryid, strn, language, sources) values (?, ?, ?, ?)"

func updateTranscriptions(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (updated bool, err error) {
	if e.ID != dbE.ID {
		return false, fmt.Errorf("update and db entry id differ")
	}

	if len(e.Transcriptions) == 0 {
		return false, fmt.Errorf("cannot update to an empty list of transcriptions")
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
			_, err := tx.Exec(transSTMT, e.ID, t.Strn, t.Language, t.SourcesString())
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

// TODO add DB trigger to set current = false for old statuses at update/insert
var statusSetCurrentFalse = "UPDATE entrystatus SET current = 0 WHERE entryid = ?"

//var insertStatus = "INSERT INTO entrystatus (entryid, name, source) values (?, ?, ?)"

// TODO always insert new status, or only when name and source have changed. Or...?
func updateEntryStatus(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (updated bool, err error) {
	if trm(e.EntryStatus.Name) != "" {
		_, err := tx.Exec(statusSetCurrentFalse, dbE.ID)
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed EntryStatus.Current update : %v", err)
		}
		_, err = tx.Exec(insertStatus, dbE.ID, strings.ToLower(e.EntryStatus.Name), strings.ToLower(e.EntryStatus.Source))
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed EntryStatus update : %v", err)
		}

		return true, nil
	}

	return false, nil
}

type vali struct {
	level string
	name  string
	msg   string
}

// TODO test me
// returns the validations of e1 not found in e2 as fist return arg
// returns the validations found only in e2, to be removed, as second arg
func newValidations(e1 lex.Entry, e2 lex.Entry) ([]lex.EntryValidation, []lex.EntryValidation) {
	var res1 []lex.EntryValidation
	var res2 []lex.EntryValidation
	e1M := make(map[vali]int)
	e2M := make(map[vali]int)
	for _, v := range e1.EntryValidations {
		e1M[vali{level: v.Level, name: v.RuleName, msg: v.Message}]++
	}
	for _, v := range e2.EntryValidations {
		e2M[vali{level: v.Level, name: v.RuleName, msg: v.Message}]++
	}

	for _, v := range e1.EntryValidations {
		if _, ok := e2M[vali{level: v.Level, name: v.RuleName, msg: v.Message}]; !ok {
			res1 = append(res1, v) // only in e1
		}
	}
	for _, v := range e2.EntryValidations {
		if _, ok := e1M[vali{level: v.Level, name: v.RuleName, msg: v.Message}]; !ok {
			res2 = append(res2, v) // only in e2
		}
	}

	return res1, res2
}

var insValiSQL = "INSERT INTO entryvalidation (entryid, level, name, message) values (?, ?, ?, ?)"

func insertEntryValidations(tx *sql.Tx, e lex.Entry, eValis []lex.EntryValidation) error {
	for _, v := range eValis {
		_, err := tx.Exec(insValiSQL, e.ID, strings.ToLower(v.Level), v.RuleName, v.Message)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert EntryValidation : %v", err)
		}
	}
	return nil
}

func updateValidation(db *sql.DB, entries []lex.Entry) error {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed starting transaction for updating validation : %v", err)
	}

	err = updateValidationTx(tx, entries)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed updating validation : %v", err)
	}
	tx.Commit()
	return nil
}

func updateValidationTx(tx *sql.Tx, entries []lex.Entry) error {
	for _, e := range entries {
		_, err := updateEntryValidationForce(tx, e)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateEntryValidationForce(tx *sql.Tx, e lex.Entry) (bool, error) {
	_, err := tx.Exec("DELETE FROM entryvalidation WHERE entryid = ?", e.ID)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed deleting EntryValidation : %v", err)
	}

	err = insertEntryValidations(tx, e, e.EntryValidations)
	if err != nil {
		return false, err
	}

	return true, nil
}

func updateEntryValidation(tx *sql.Tx, e lex.Entry, dbE lex.Entry) (bool, error) {
	newValidations, removeValidations := newValidations(e, dbE)
	if len(newValidations) == 0 && len(removeValidations) == 0 {
		return false, nil
	}

	err := insertEntryValidations(tx, dbE, newValidations)
	if err != nil {
		return false, err
	}

	for _, v := range removeValidations {
		_, err := tx.Exec("DELETE FROM entryvalidation WHERE id = ?", v.ID)
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed deleting EntryValidation : %v", err)
		}
	}

	return true, nil
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

func entryCount(db *sql.DB, lexiconName string) (int64, error) {
	tx, err := db.Begin()
	defer tx.Commit()

	if err != nil {
		return -1, fmt.Errorf("dbapi.EntryCount failed opening db transaction : %v", err)
	}

	// number of entries in a lexicon
	var entries int64
	err = tx.QueryRow("SELECT COUNT(*) FROM entry, lexicon WHERE entry.lexiconid = lexicon.id and lexicon.name = ?", lexiconName).Scan(&entries)
	if err != nil || err == sql.ErrNoRows {
		return -1, fmt.Errorf("dbapi.EntryCount failed QueryRow : %v", err)
	}
	return entries, nil
}

// EntryCount counts the number of lines in a lexicon
// func entryCount(db *sql.DB, lexiconID int64) (int64, error) {
// 	tx, err := db.Begin()
// 	defer tx.Commit()

// 	if err != nil {
// 		return -1, fmt.Errorf("dbapi.EntryCount failed opening db transaction : %v", err)
// 	}

// 	// number of entries in a lexicon
// 	var entries int64
// 	err = tx.QueryRow("SELECT COUNT(*) FROM entry WHERE entry.lexiconid = ?", lexiconID).Scan(&entries)
// 	if err != nil || err == sql.ErrNoRows {
// 		return -1, fmt.Errorf("dbapi.EntryCount failed QueryRow : %v", err)
// 	}
// 	return entries, nil
// }

// ListCurrentEntryStatuses returns a list of all names EntryStatuses marked 'current' (i.e., the most recent status).
func listCurrentEntryStatuses(db *sql.DB, lexiconName string) ([]string, error) {
	return listEntryStatuses(db, lexiconName, true)
}

// ListAllEntryStatuses returns a list of all names EntryStatuses, also those that are not 'current'  (i.e., the most recent status).
// In other words, this list potentially includes statuses not in use, but that have been used.
func listAllEntryStatuses(db *sql.DB, lexiconName string) ([]string, error) {
	return listEntryStatuses(db, lexiconName, false)
}

func listEntryStatuses(db *sql.DB, lexiconName string, onlyCurrent bool) ([]string, error) {
	var res []string

	tx, err := db.Begin()
	if err != nil {
		return res, fmt.Errorf("dbapi.ListCurrentEntryStatuses : %v", err)
	}
	defer tx.Commit()

	// TODO This query seems a bit slow?
	q := "SELECT DISTINCT entryStatus.name FROM lexicon, entry, entryStatus WHERE lexicon.name = ? AND lexicon.id = entry.lexiconID and entry.id = entryStatus.entryId"
	qOnlyCurrent := " AND entryStatus.current = 1"
	if onlyCurrent {
		q += qOnlyCurrent
	}

	rows, err := tx.Query(q, lexiconName)
	if err != nil {
		tx.Rollback() // ?
		return res, fmt.Errorf("ListCurrentEntryStatuses : %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var statusName string
		rows.Scan(&statusName)
		res = append(res, statusName)
	}

	err = rows.Err()
	return res, err
}

// LexiconStats calls the database a number of times, gathering different numbers, e.g. on how many entries there are in a lexicon.
func lexiconStats(db *sql.DB, lexName string) (LexStats, error) {
	res := LexStats{Lexicon: lexName}

	tx, err := db.Begin()
	if err != nil {
		return res, fmt.Errorf("dbapi.LexiconStats failed opening db transaction : %v", err)
	}
	defer tx.Commit()

	lex, err := getLexiconTx(tx, lexName)
	if err != nil {
		return res, fmt.Errorf("dbapi.LexiconStats failed getting lexicon id : %v", err)
	}
	lexiconID := lex.id

	// t1 := time.Now()

	// number of entries in a lexicon
	var entries int64
	err = tx.QueryRow("SELECT COUNT(*) FROM entry WHERE entry.lexiconid = ?", lexiconID).Scan(&entries)
	if err != nil || err == sql.ErrNoRows {
		return res, fmt.Errorf("dbapi.LexiconStats failed QueryRow : %v", err)
	}
	res.Entries = entries

	// number of each type of entry status

	// t2 := time.Now()
	// log.Printf("dbapi.LexiconStats TOTAL COUNT TOOK %v\n", t2.Sub(t1))

	rows, err := tx.Query("select entrystatus.name, count(entrystatus.name) from entry, entrystatus where entry.lexiconid = ? and entry.id = entrystatus.entryid and entrystatus.current = 1 group by entrystatus.name", lexiconID)
	if err != nil {
		return res, fmt.Errorf("db query failed : %v", err)
	}

	defer rows.Close()
	for rows.Next() {
		var status string
		var freq string
		err = rows.Scan(&status, &freq)
		if err != nil {
			return res, fmt.Errorf("scanning row failed : %v", err)
		}

		res.StatusFrequencies = append(res.StatusFrequencies, status+"\t"+freq)
	}
	err = rows.Err()
	if err != nil {
		return res, err
	}

	// t3 := time.Now()
	// log.Printf("dbapi.LexiconStats COUNT PER STATUS TOOK %v\n", t3.Sub(t2))

	valStats, err := validationStatsTx(tx, lexiconID)
	res.ValStats = valStats

	// t4 := time.Now()
	// log.Printf("dbapi.LexiconStats VAL STATS TOOK %v\n", t4.Sub(t3))

	// _ = t1
	// _ = t4
	// log.Printf("dbapi.LexiconStats STATS TOOK %v\n", t4.Sub(t1))

	return res, err

}

func validationStats(db *sql.DB, lexName string) (ValStats, error) {
	tx, err := db.Begin()
	defer tx.Commit()

	if err != nil {
		return ValStats{}, fmt.Errorf("dbapi.ValidationStats failed opening db transaction : %v", err)
	}
	lex, err := getLexiconTx(tx, lexName)
	if err != nil {
		return ValStats{}, fmt.Errorf("dbapi.LexiconStats failed getting lexicon id : %v", err)
	}
	lexID := lex.id
	return validationStatsTx(tx, lexID)
}

func validationStatsTx(tx *sql.Tx, lexiconID int64) (ValStats, error) {

	res := ValStats{Rules: make(map[string]int), Levels: make(map[string]int)}

	// number of entries in the lexicon
	err := tx.QueryRow("SELECT COUNT(*) FROM entry WHERE entry.lexiconid = ?", lexiconID).Scan(&res.TotalEntries)
	if err != nil || err == sql.ErrNoRows {
		return res, fmt.Errorf("dbapi.ValidationStats failed QueryRow : %v", err)
	}

	res.ValidatedEntries = res.TotalEntries

	// number of invalid entries
	err = tx.QueryRow("SELECT COUNT (DISTINCT entryvalidation.entryid) FROM entry, entryvalidation WHERE entry.id = entryvalidation.entryid AND entry.lexiconid = ?", lexiconID).Scan(&res.InvalidEntries)
	if err != nil || err == sql.ErrNoRows {
		return res, fmt.Errorf("dbapi.ValidationStats failed QueryRow : %v", err)
	}

	// number of validations
	err = tx.QueryRow("SELECT COUNT (DISTINCT entryvalidation.id) FROM entry, entryvalidation WHERE entry.id = entryvalidation.entryid AND entry.lexiconid = ?", lexiconID).Scan(&res.TotalValidations)
	if err != nil || err == sql.ErrNoRows {
		return res, fmt.Errorf("dbapi.ValidationStats failed QueryRow : %v", err)
	}

	levels, err := tx.Query("select entryvalidation.level, count(entryvalidation.level) from entry, entryvalidation where entry.lexiconid = ? and entry.id = entryvalidation.entryid group by entryvalidation.level", lexiconID)
	if err != nil {
		return res, fmt.Errorf("db query failed : %v", err)
	}

	defer levels.Close()
	for levels.Next() {
		var name string
		var count int
		err = levels.Scan(&name, &count)
		if err != nil {
			return res, fmt.Errorf("scanning row failed : %v", err)
		}

		res.Levels[strings.ToLower(name)] = count
	}
	err = levels.Err()
	if err != nil {
		return res, err
	}

	names, err := tx.Query("select entryvalidation.level, entryvalidation.name, count(entryvalidation.name) from entry, entryvalidation where entry.lexiconid = ? and entry.id = entryvalidation.entryid group by entryvalidation.name", lexiconID)
	if err != nil {
		return res, fmt.Errorf("db query failed : %v", err)
	}

	defer names.Close()
	for names.Next() {
		var name string
		var level string
		var count int
		err = names.Scan(&level, &name, &count)
		nameWithLevel := fmt.Sprintf("%s (%s)", strings.ToLower(name), strings.ToLower(level))
		if err != nil {
			return res, fmt.Errorf("scanning row failed : %v", err)
		}

		res.Rules[nameWithLevel] = count
	}
	err = names.Err()
	if err != nil {
		return res, err
	}

	// finally
	return res, err

}
