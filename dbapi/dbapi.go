/*
Package dbapi contains code wrapped around an SQL(ite3) DB.
It is used for inserting, updating and retrieving lexical entries from
a pronounciation lexicon database. A lexical entry is represented by
the dbapi.Entry struct, that mirrors entries of the entry database
table, along with associated tables such as transcription and lemma.
*/
package dbapi

//go get github.com/mattn/go-sqlite3

import (
	"database/sql"
	"fmt"
	// installs sqlite3 driver
	"github.com/mattn/go-sqlite3"
	"github.com/stts-se/pronlex/lex"
	//"github.com/stts-se/pronlex/validation"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var remem = make(map[string]*regexp.Regexp)
var regexMem = func(re, s string) (bool, error) {
	if r, ok := remem[re]; ok {
		return r.MatchString(s), nil
	}
	r, err := regexp.Compile(re)
	if err != nil {
		return false, err
	}
	remem[re] = r
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

// TODO not tested... but what could possibly go wrong...?:

// DeleteUnusedSymbolSets deletes symbols, from the SYMBOLSET table, that have lexiconIDs that do not exist in the LEXICON table.
func DeleteUnusedSymbolSets(db *sql.DB) error {
	// TODO this should optimally be taken care of by restrictions/dependencies in the db
	var sql = "DELETE FROM symbolset WHERE lexiconid NOT IN (SELECT id from lexicon)"
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("DeleteUnusedSymbolSets failed starting transaction: %v", err)
	}
	defer tx.Commit()

	_, err = tx.Exec(sql)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("DeleteUnusedSymbolSets failed : %v", err)
	}
	return err
}

// ListLexicons returns a list of the lexicons defined in the db
// (i.e., Lexicon structs corresponding to the rows of the lexicon
// table)
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

// GetLexicon returns a Lexicon struct matching a lexicon name in the db.
// Returns error if no such lexicon name in db
func GetLexicon(db *sql.DB, name string) (Lexicon, error) {
	res0, err := GetLexicons(db, []string{name})
	if err != nil {
		return Lexicon{}, fmt.Errorf("failed to retrieve lexicon %s : %v", name, err)
	}
	if len(res0) != 1 {
		return Lexicon{}, fmt.Errorf("failed to retrieve lexicon %s : %v", name, err)
	}
	return res0[0], nil
}

// GetLexicons takes a list of lexicon names and returns a list of
// Lexicon structs corresponding to rows of db lexicon table with those name fields.
func GetLexicons(db *sql.DB, names []string) ([]Lexicon, error) {
	var res []Lexicon
	found := make(map[string]bool)
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
		found[strings.ToLower(lname)] = true
		res = append(res, Lexicon{ID: id, Name: lname, SymbolSetName: symbolsetname})
	}
	err = rows.Err()
	rows.Close()

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
func LexiconFromID(db *sql.DB, id int64) (Lexicon, error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return Lexicon{}, fmt.Errorf("LexiconFromID failed to start db transaction : %v", err)
	}

	return LexiconFromIDTx(tx, id)
}

// LexiconFromIDTx returns a Lexicon struct corresponding to a row in
// the lexicon table with the given id
func LexiconFromIDTx(tx *sql.Tx, id int64) (Lexicon, error) {
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
// It should be impossible to delete the Lexicon table entry if associated to any entries.
func DeleteLexicon(db *sql.DB, id int64) error {
	fmt.Printf("DeleteLexicon called with id %d\n", id)
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return err
	}
	return DeleteLexiconTx(tx, id)
}

// DeleteLexiconTx deletes the lexicon name from the lexicon
// table. Notice that it does not remove the associated entries.
// It should be impossible to delete the Lexicon table entry if associated to any entries.
func DeleteLexiconTx(tx *sql.Tx, id int64) error {
	var n int
	err := tx.QueryRow("select count(*) from entry where entry.lexiconid = ?", id).Scan(&n)
	// must always return a row, no need to check for empty row
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("HEJ DIN FAN")
			return fmt.Errorf("The was no lexicon with id %d : %v", id, err)
		}
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

// SuperDeleteLexicon deletes the lexicon name from the lexicon
// table and also whipes all associated entries out of existence.
// It also deletes all entries from the Symbolset table associated to the lexicon.
// func SuperDeleteLexicon0(db *sql.DB, id int64) error {
// 	fmt.Printf("DeleteLexicon called with id %d\n", id)
// 	tx, err := db.Begin()
// 	defer tx.Commit()
// 	if err != nil {
// 		return err
// 	}
// 	return SuperDeleteLexiconTx(tx, id)
// }

// SuperDeleteLexicon deletes the lexicon name from the lexicon
// table and also whipes all associated entries out of existence.
// It also deletes all entries from the Symbolset table associated to the lexicon.
// TODO Send progress message to client over websocket
func SuperDeleteLexicon(db *sql.DB, id int64) error {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		return fmt.Errorf("SuperDeleteLexicon failed to initiate transaction : %v", err)
	}
	return SuperDeleteLexiconTx(tx, id)
}

// SuperDeleteLexiconTx deletes the lexicon name from the lexicon
// table and also whipes all associated entries out of existence.
// It also deletes all entries from the Symbolset table associated to the lexicon.
func SuperDeleteLexiconTx(tx *sql.Tx, id int64) error {

	fmt.Println("dbapi.superDeleteLexiconTX was called")

	_, err := tx.Exec("DELETE FROM symbolset WHERE lexiconid = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("dbapi.SuperDeleteLexiconTx : failed to delete symbol set : %v", err)
	}

	fmt.Println("dbapi.superDeleteLexiconTX finished deleting from symbol set")

	_, err = tx.Exec("DELETE FROM entry WHERE lexiconid = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("dbapi.SuperDeleteLexiconTx : failed to delete entries : %v", err)
	}

	fmt.Println("dbapi.superDeleteLexiconTX finished deleting from entry set")

	_, err = tx.Exec("DELETE FROM lexicon WHERE id = ?", id)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("dbapi.SuperDeleteLexiconTx : failed to delete lexicon : %v", err)
	}

	fmt.Println("dbapi.superDeleteLexiconTX finished deleting from lexicon set")

	fmt.Printf("Deleting lexicon %d\n", id)

	return nil
}

// InsertOrUpdateLexicon takes a Lexicon struct and either inserts it into the db, if its id = 0, or updates its string fields if the id is greater than 0.
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

	if l.ID <= 0 {
		return InsertLexiconTx(tx, l)
	}

	// else if l.ID > 0
	res, err := LexiconFromIDTx(tx, l.ID)

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

	return res, nil
}

// InsertLexicon saves the name of a new lexicon to the db.
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

// TODO move to function?
var entrySTMT = "insert into entry (lexiconid, strn, language, partofspeech, wordparts) values (?, ?, ?, ?, ?)"
var transAfterEntrySTMT = "insert into transcription (entryid, strn, language, sources) values (?, ?, ?, ?)"

//var statusSetCurrentFalse = "UPDATE entrystatus SET current = 0 WHERE entrystatus.entryid = ?"
var insertStatus = "INSERT INTO entrystatus (entryid, name, source) values (?, ?, ?)"

// InsertEntries saves a list of Entries and associates them to Lexicon
// TODO change input arg to sql.Tx
func InsertEntries(db *sql.DB, l Lexicon, es []lex.Entry) ([]int64, error) {

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

// InsertLemma saves a lex.Lemma to the db, but does not associate it with anlex.Entry
// TODO do we need both InsertLemma and SetOrGetLemma?
func InsertLemma(tx *sql.Tx, l lex.Lemma) (lex.Lemma, error) {
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
func SetOrGetLemma(tx *sql.Tx, strn string, reading string, paradigm string) (lex.Lemma, error) {
	res := lex.Lemma{}

	var id int64
	var strn0, reading0, paradigm0 string
	sqlS := "select id, strn, reading, paradigm from lemma where strn = ? and reading = ?"
	err := tx.QueryRow(sqlS, strn, reading).Scan(&id, &strn0, &reading0, &paradigm0)
	switch {
	case err == sql.ErrNoRows:
		return InsertLemma(tx, lex.Lemma{ID: id, Strn: strn, Reading: reading, Paradigm: paradigm})
	case err != nil:
		return res, fmt.Errorf("SetOrGetLemma failed querying db : %v", err)
	}

	res.ID = id
	res.Strn = strn0
	res.Reading = reading0
	res.Paradigm = paradigm0

	return res, err
}

// AssociateLemma2Entry adds a lex.Lemma to anlex.Entry via a linking table
func AssociateLemma2Entry(db *sql.Tx, l lex.Lemma, e lex.Entry) error {
	sql := "insert into Lemma2Entry (lemmaId, entryId) values (?, ?)"
	_, err := db.Exec(sql, l.ID, e.ID)
	if err != nil {
		err = fmt.Errorf("failed to associate lemma "+l.Strn+" and entry "+e.Strn+":%v", err)
	}
	return err
}

func getLemmaFromEntryIDTx(tx *sql.Tx, id int64) (lex.Lemma, error) {
	res := lex.Lemma{}
	sqlS := "select lemma.id, lemma.strn, lemma.reading, lemma.paradigm from entry, lemma, lemma2entry where " +
		"entry.id = ? and entry.id = lemma2entry.entryid and lemma.id = lemma2entry.lemmaid"
	var lID int64
	var strn, reading, paradigm string
	err := tx.QueryRow(sqlS, id).Scan(&lID, &strn, &reading, &paradigm)
	switch {
	case err == sql.ErrNoRows:
		// TODO No row:
		// Silently return empty lex.Lemma below
	case err != nil:
		//ff("getLemmaFromENtryId: %v", err)
		return res, fmt.Errorf("QueryRow failure : %v", err)
	}

	// TODO Now silently returns empty lemma if nothing returned from db. Ok?
	res.ID = lID
	res.Strn = strn
	res.Reading = reading
	res.Paradigm = paradigm

	return res, nil
}

// LookUp takes a Query struct, searches the lexicon db, and writes the result to the
//lex.EntryWriter.
func LookUp(db *sql.DB, q Query, out lex.EntryWriter) error {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to initialize transaction : %v", err)
	}
	return LookUpTx(tx, q, out)
}

// LookUpTx takes a Query struct, searches the lexicon db, and writes the result to the
// EntryWriter.
// TODO: rewrite to go through the result set before building the result. That is, save all structs corresponding to rows in the scanning run, then build the result structure (so that no identical values are duplicated: a result set may have several rows of repeated data)
func LookUpTx(tx *sql.Tx, q Query, out lex.EntryWriter) error {

	//fmt.Printf("QUWRY %v\n\n", q)

	sqlString, values := SelectEntriesSQL(q)

	//fmt.Printf("SQL %v\n\n", sqlString)

	//fmt.Printf("VALUES %v\n\n", values)

	rows, err := tx.Query(sqlString, values...)
	if err != nil {
		tx.Rollback() // nothing to rollback here, but may have been called from withing another transaction
		return err
	}

	var lexiconID, entryID int64
	var entryStrn, entryLanguage, partOfSpeech, wordParts string

	var transcriptionID, transcriptionEntryID int64
	var transcriptionStrn, transcriptionLanguage, transcriptionSources string

	// Optional/nullable values

	var lemmaID sql.NullInt64
	var lemmaStrn, lemmaReading, lemmaParadigm sql.NullString

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
			&transcriptionSources,

			// Optional, from LEFT JOIN

			&lemmaID,
			&lemmaStrn,
			&lemmaReading,
			&lemmaParadigm,

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
		if lastE != entryID {
			if lastE != -1 {
				out.Write(currE)
			}
			currE = lex.Entry{
				LexiconID:    lexiconID,
				ID:           entryID,
				Strn:         entryStrn,
				Language:     entryLanguage,
				PartOfSpeech: partOfSpeech,
				WordParts:    wordParts,
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
		tx.Rollback() // nothing to rollback here, but may have been called from withing another transaction
		return rows.Err()
	}

	return nil
}

// LookUpIntoSlice is a wrapper around LookUp, returning a slice of Entries
func LookUpIntoSlice(db *sql.DB, q Query) ([]lex.Entry, error) {
	var esw lex.EntrySliceWriter
	err := LookUp(db, q, &esw)
	if err != nil {
		return esw.Entries, fmt.Errorf("failed lookup : %v", err)
	}
	return esw.Entries, nil
}

// LookUpIntoMap is a wrapper around LookUp, returning a map where the
// keys are word forms and the values are slices of Entries. (There may be several entries with the same Strn value.)
func LookUpIntoMap(db *sql.DB, q Query) (map[string][]lex.Entry, error) {
	res := make(map[string][]lex.Entry)
	var esw lex.EntrySliceWriter
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

// GetEntryFromID is a wrapper around LookUp and returns the lex.Entry corresponding to the db id
func GetEntryFromID(db *sql.DB, id int64) (lex.Entry, error) {
	res := lex.Entry{}
	q := Query{EntryIDs: []int64{id}}
	esw := lex.EntrySliceWriter{}
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

// UpdateEntry wraps call to UpdateEntryTx with a transaction, and returns the updated entry, fresh from the db
// TODO Consider how to handle inconsistent input entries
func UpdateEntry(db *sql.DB, e lex.Entry) (res lex.Entry, updated bool, err error) {
	tx, err := db.Begin()
	defer tx.Commit()
	if err != nil {
		tx.Rollback()
		return res, updated, fmt.Errorf("failed starting transaction for updating entry : %v", err)
	}

	updated, err = UpdateEntryTx(tx, e)
	if err != nil {
		tx.Rollback()
		return res, updated, fmt.Errorf("failed updating entry : %v", err)
	}
	tx.Commit()
	res, err = GetEntryFromID(db, e.ID)
	if err != nil {
		tx.Rollback()
		return res, updated, fmt.Errorf("failed getting updated entry : %v", err)
	}

	return res, updated, err
}

// UpdateEntryTx updates the fields of an lex.Entry that do not match the
// corresponding values in the db
func UpdateEntryTx(tx *sql.Tx, e lex.Entry) (updated bool, err error) { // TODO return the updated entry?
	// updated == false
	//dbEntryMap := //GetEntriesFromIDsTx(tx, []int64{(e.ID)})
	var esw lex.EntrySliceWriter
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
	updated5, err := updateEntryStatus(tx, e, dbEntries[0])
	if err != nil {
		return updated5, err
	}

	updated6, err := updateEntryValidation(tx, e, dbEntries[0])
	if err != nil {
		return updated6, err
	}

	return updated1 || updated2 || updated3 || updated4 || updated5 || updated6, err
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
// of the same LexiconID. Silently ingnores the symbol if Symbol.Symbol == "".
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
		//if trm(s.Symbol) != "" { // cannot trim spaces ...
		_, err = tx.Exec("insert into symbolset (lexiconid, symbol, category, description, ipa) values (?, ?, ?, ?, ?)",
			s.LexiconID, s.Symbol, s.Category, s.Description, s.IPA)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed inserting symbol : %v", err)
		}
		//}
	}
	return nil
}

// GetSymbolSet returns the set of Symbols defined for a lexicon with the given db id
func GetSymbolSet(db *sql.DB, lexiconID int64) ([]Symbol, error) {
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
	rows, err := tx.Query("select lexiconid, symbol, category, description, ipa from symbolset where lexiconid = ? order by id", lexiconID)
	if err != nil {
		return res, fmt.Errorf("failed db query : %v", err)
	}

	var lexID int64
	var symbol, category, description, ipa string
	for rows.Next() {
		rows.Scan(&lexID, &symbol, &category, &description, &ipa)
		s := Symbol{
			LexiconID:   lexID,
			Symbol:      symbol,
			Category:    category,
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

// LexiconStats calls the database a number of times, gathering different numbers, e.g. on how many entries there are in a lexicon.
func LexiconStats(db *sql.DB, lexiconID int64) (LexStats, error) {
	res := LexStats{LexiconID: lexiconID}

	tx, err := db.Begin()
	defer tx.Commit()

	if err != nil {
		return res, fmt.Errorf("dbapi.LexiconStats failed opening db transaction : %v", err)
	}

	// number of entries in a lexicon
	var entries int64
	err = tx.QueryRow("SELECT COUNT(*) FROM entry WHERE entry.lexiconid = ?", lexiconID).Scan(&entries)
	if err != nil || err == sql.ErrNoRows {
		return res, fmt.Errorf("dbapi.LexiconStats failed QueryRow : %v", err)
	}

	res.Entries = entries

	// number of each type of entry status

	//select entrystatus.name, count(entrystatus.name) from entry, entrystatus where entry.lexiconid = 3 and entry.id = entrystatus.entryid and entrystatus.current = 1 group by entrystatus.name

	rows, err := tx.Query("select entrystatus.name, count(entrystatus.name) from entry, entrystatus where entry.lexiconid = ? and entry.id = entrystatus.entryid and entrystatus.current = 1 group by entrystatus.name", lexiconID)
	defer rows.Close()
	if err != nil {
		return res, fmt.Errorf("db query failed : %v", err)
	}

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
	return res, err

	// TODO add queries for additional stats

	//return res, nil

}
