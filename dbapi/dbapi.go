package dbapi

//go get github.com/mattn/go-sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sort"
	"strings"
)

// TODO
// f is a place holder to be replaced by proper error handling
func f(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// TODO
// ff is a place holder to be replaced by proper error handling
func ff(f string, err error) {
	if err != nil {
		log.Fatalf(f, err)
	}
}

func ListLexicons(db *sql.DB) ([]Lexicon, error) {
	res := make([]Lexicon, 0)
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
		l := Lexicon{Id: id, Name: name, SymbolSetName: symbolSetName}
		res = append(res, l)
	}

	return res, nil
}

func GetLexicons(db *sql.DB, names []string) ([]Lexicon, error) {
	res := make([]Lexicon, 0)
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
		res = append(res, Lexicon{Id: id, Name: lname, SymbolSetName: symbolsetname})
	}
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

	return Lexicon{Id: id, Name: lname, SymbolSetName: symbolsetname}, nil
}

// TODO change input arg to sql.Tx ?
func InsertLexicon(db *sql.DB, l Lexicon) (Lexicon, error) {
	tx, err := db.Begin()
	defer tx.Commit()

	res, err := tx.Exec("insert into lexicon (name, symbolsetname) values (?, ?)", strings.ToLower(l.Name), l.SymbolSetName)
	if err != nil {
		return l, fmt.Errorf("failed to insert lexicon + symbolset name : %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return l, fmt.Errorf("failed to get last insert id : %v", err)
	}

	tx.Commit()

	return Lexicon{Id: id, Name: l.Name, SymbolSetName: l.SymbolSetName}, err
}

// TODO return error
// TODO change input arg to sql.Tx
func InsertEntries(db *sql.DB, l Lexicon, es []Entry) []int64 {

	var entrySTMT = "insert into entry (lexiconid, strn, language, partofspeech, wordparts) values (?, ?, ?, ?, ?)"
	// TODO move to function
	var transAfterEntrySTMT = "insert into transcription (entryid, strn, language) values (?, ?, ?)"

	ids := make([]int64, 0)

	// Transaction -->
	tx, err := db.Begin()
	defer tx.Commit()

	stmt1, err := tx.Prepare(entrySTMT)
	if err != nil {
		log.Fatal(err)
	}
	stmt2, err := tx.Prepare(transAfterEntrySTMT)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	for _, e := range es {
		res, err := tx.Stmt(stmt1).Exec(
			l.Id,
			strings.ToLower(e.Strn),
			e.Language,
			e.PartOfSpeech,
			e.WordParts)
		if err != nil {
			log.Fatal(err) // TODO rollback
		}

		id, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err) // TODO rollback
		}
		// We want the Entry to have the right id for inserting lemma assocs below
		e.Id = id

		ids = append(ids, id)

		// res.Close()

		for _, t := range e.Transcriptions {
			_, err := tx.Stmt(stmt2).Exec(id, t.Strn, t.Language)
			if err != nil {
				log.Fatal(err) // TODO rollback
			}
		}

		//log.Printf("%v", e)
		if "" != e.Lemma.Strn && "" != e.Lemma.Reading {
			lemma, err := SetOrGetLemma(tx, e.Lemma.Strn, e.Lemma.Reading, e.Lemma.Paradigm)
			ff("Failed to insert lemma: %v", err)
			err = AssociateLemma2Entry(tx, lemma, e)
			ff("Failed lemma to entry assoc: %v", err)
		}
	}

	tx.Commit()
	// <- transaction

	return ids
}

func AssociateLemma2Entry(db *sql.Tx, l Lemma, e Entry) error {
	sql := "insert into Lemma2Entry (lemmaId, entryId) values (?, ?)"
	_, err := db.Exec(sql, l.Id, e.Id)
	if err != nil {
		err = fmt.Errorf("failed to associate lemma "+l.Strn+" and entry "+e.Strn+":%v", err)
	}
	return err
}

func SetOrGetLemma(tx *sql.Tx, strn string, reading string, paradigm string) (Lemma, error) {
	res := Lemma{}

	var id int64
	var strn0, reading0, paradigm0 string
	sqlS := "select id, strn, reading, paradigm from lemma where strn = ? and reading = ?"
	err := tx.QueryRow(sqlS, strn, reading).Scan(&id, &strn0, &reading0, &paradigm0)
	switch {
	case err == sql.ErrNoRows:
		return InsertLemma(tx, Lemma{Id: id, Strn: strn, Reading: reading, Paradigm: paradigm})
	case err != nil:
		return res, fmt.Errorf("SetOrGetLemma failed querying db : %v", err)
	}

	res.Id = id
	res.Strn = strn0
	res.Reading = reading0
	res.Paradigm = paradigm0

	return res, err
}

// TODO return error
func getLemmaFromEntryIdTx(tx *sql.Tx, id int64) Lemma {
	res := Lemma{}
	sqlS := "select lemma.id, lemma.strn, lemma.reading, lemma.paradigm from entry, lemma, lemma2entry where " +
		"entry.id = ? and entry.id = lemma2entry.entryid and lemma.id = lemma2entry.lemmaid"
	var lId int64
	var strn, reading, paradigm string
	err := tx.QueryRow(sqlS, id).Scan(&lId, &strn, &reading, &paradigm)
	switch {
	case err == sql.ErrNoRows:
		// TODO No row:
		// Silently return empty Lemma below
	case err != nil:
		ff("getLemmaFromENtryId: %v", err)
	}

	// TODO Now silently returns empty lemma if nothing returned from db. Ok?
	// Return err when there is an err
	res.Id = lId
	res.Strn = strn
	res.Reading = reading
	res.Paradigm = paradigm

	return res
}

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
	l.Id = id
	return l, err
}

func entryMapToEntrySlice(em map[string][]Entry) []Entry {
	res := make([]Entry, 0)
	for _, v := range em {
		res = append(res, v...)
	}
	return res
}

// TODO return error
// TODO this should return []Entry rather than map[string][]Entry?
func GetEntriesFromIdsTx(tx *sql.Tx, entryIds []int64) map[string][]Entry {
	res := make(map[string][]Entry)
	if len(entryIds) == 0 {
		return res
	}

	qString, values := entriesFromIdsSelect(entryIds)
	rows, err := tx.Query(qString, values...)
	if err != nil {
		log.Fatalf("EntriesFromIds: %s", err)
	}
	defer rows.Close()

	// entries map Entries so far. One entry may have several transcriptions, resulting in multiple rows for a single entry
	entries := make(map[int64]Entry)
	transes := make(map[int64]Transcription)
	for rows.Next() {
		var entryId, entryLexiconId int64
		var entryStrn, entryLanguage, entryPartofSpeech, entryWordParts string
		var transcriptionId, transcriptionEntryId int64
		var transcriptionStrn, transcriptionLanguage string

		if err := rows.Scan(&entryId,
			&entryLexiconId,
			&entryStrn,
			&entryLanguage,
			&entryPartofSpeech,
			&entryWordParts,
			&transcriptionId,
			&transcriptionEntryId,
			&transcriptionStrn,
			&transcriptionLanguage); err != nil {
			log.Fatal(err)
		}

		// collect Entry and Transcription separately. insert []Transcription into Entry below

		// collect unique Entries
		if _, ok := entries[entryId]; !ok {
			e := Entry{
				Id:           entryId,
				LexiconId:    entryLexiconId,
				Strn:         entryStrn,
				Language:     entryLanguage,
				PartOfSpeech: entryPartofSpeech,
				WordParts:    entryWordParts,
				Lemma:        getLemmaFromEntryIdTx(tx, entryId),
			}

			entries[entryId] = e
		}

		// collect unique Transcriptions
		if _, ok := transes[transcriptionId]; !ok {
			t := Transcription{
				Id:       transcriptionId,
				EntryId:  transcriptionEntryId,
				Strn:     transcriptionStrn,
				Language: transcriptionLanguage,
			}
			transes[transcriptionId] = t
		}

	} // rows.Next()

	// map entry ids to transcriptions
	eId2ts := make(map[int64][]Transcription)
	for _, t := range transes {
		eId2ts[t.EntryId] = append(eId2ts[t.EntryId], t)
	}

	// Put together Entries and Transcriptions and build up return map res
	for id, e := range entries {
		var ts []Transcription
		var ok bool
		ts, ok = eId2ts[id]
		if !ok {
			log.Fatal("EntriesFromIds: Entry id unknown")
		}
		sort.Sort(TranscriptionSlice(ts)) // Sort according to id. See structs_dbapi.go.
		e.Transcriptions = ts
		res[e.Strn] = append(res[e.Strn], e)
	}

	return res
}

// TODO return error
// TODO error handling!!!
func GetEntryFromId(db *sql.DB, id int64) Entry {
	res := GetEntriesFromIds(db, []int64{id})
	// if len(res) != 1 {
	// 	return nil
	// }
	//return
	for _, v := range res {
		return v[0]
	}

	return Entry{}
}

// TODO return error
// TODO this should return []Entry rather than map[string][]Entry?
func GetEntriesFromIds(db *sql.DB, entryIds []int64) map[string][]Entry {
	tx, err := db.Begin()
	f(err)
	defer tx.Commit()
	return GetEntriesFromIdsTx(tx, entryIds)
}

// TODO return error
// TODO should be a wrapper to GetEntriesTx
func GetEntries(db *sql.DB, q Query) map[string][]Entry {
	res := make(map[string][]Entry)
	if q.Empty() { // TODO report to client?
		log.Printf("dbapi.GetEntries: Query empty of search constraints: %v", q)
		return res
	}

	qString, vs := idiotSql(q)

	rows, err := db.Query(qString, vs...)
	if err != nil {
		log.Fatalf("dbapi.GetEntries:\t%s", err)
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			log.Fatalf("GetEntries(2):\t%s", err)
		}
		ids = append(ids, id)
	}

	// TODO return map shuold be built here rather than by GetEntriesFromIds?
	res = GetEntriesFromIds(db, ids)
	return res
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

func getTIds(ts []Transcription) []int64 {
	res := make([]int64, 0)
	for _, t := range ts {
		res = append(res, t.Id)
	}
	return res
}

func equal(ts1 []Transcription, ts2 []Transcription) bool {
	if len(ts1) != len(ts2) {
		return false
	}
	for i, _ := range ts1 {
		if ts1[i] != ts2[i] {
			return false
		}
	}

	return true
}

func updateLanguage(tx *sql.Tx, e Entry, dbE Entry) (bool, error) {
	if e.Id != dbE.Id {
		return false, fmt.Errorf("new and old entries have different ids")
	}
	if e.Language == dbE.Language {
		return false, nil
	}
	_, err := tx.Exec("update entry set language = ? where entry.id = ?", e.Language, e.Id)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed language update : %v", err)
	}
	return true, nil
}

func updateWordParts(tx *sql.Tx, e Entry, dbE Entry) (bool, error) {
	if e.Id != dbE.Id {
		return false, fmt.Errorf("new and old entries have different ids")
	}
	if e.WordParts == dbE.WordParts {
		return false, nil
	}
	_, err := tx.Exec("update entry set wordparts = ? where entry.id = ?", e.WordParts, e.Id)
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
	if e.Lemma.Id == 0 && e.Lemma.Strn == "" {
		_, err = tx.Exec("delete from lemma where lemma.id = ?", dbE.Lemma.Id)
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed to delete old lemma : %v", err)
		}
	}
	// Only one alternative left, to update old lemma with new values
	_, err = tx.Exec("update lemma set strn = ?, reading = ?, paradigm = ? where lemma.id = ?", e.Lemma.Strn, e.Lemma.Reading, e.Lemma.Paradigm, dbE.Lemma.Id)
	if err != nil {
		tx.Rollback()
		return false, fmt.Errorf("failed to update lemma : %v", err)
	}
	return true, nil
}

// TODO move to function
var transSTMT = "insert into transcription (entryid, strn, language) values (?, ?, ?)"

func updateTranscriptions(tx *sql.Tx, e Entry, dbE Entry) (updated bool, err error) {
	if e.Id != dbE.Id {
		return false, fmt.Errorf("update and db entry id differ")
	}

	// the easy way would be to simply nuke any transcriptions for
	// the entry and substitute for the new ones, but we only want
	// to save new transcriptions if there are changes

	// If the new and old transcriptions differ, remove the old
	// and inser the new ones
	if !equal(e.Transcriptions, dbE.Transcriptions) {
		transIds := getTIds(dbE.Transcriptions)
		// TODO move to a function
		_, err := tx.Exec("delete from transcription where transcription.id in "+nQs(len(transIds)), convI(transIds)...)
		if err != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed transcription delete : %v", err)
		}
		for _, t := range e.Transcriptions {
			_, err := tx.Exec(transSTMT, e.Id, t.Strn, t.Language)
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

func UpdateEntryTx(tx *sql.Tx, e Entry) (updated bool, err error) { // TODO return the updated entry?
	// updated == false
	dbEntryMap := GetEntriesFromIdsTx(tx, []int64{(e.Id)})
	dbEntries := entryMapToEntrySlice(dbEntryMap)
	if len(dbEntries) == 0 {

		return updated, fmt.Errorf("no entry with id '%d'", e.Id)
	}
	if len(dbEntries) > 1 {

		return updated, fmt.Errorf("very bad error, more than one entry with id '%d'", e.Id)
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

func Nothing() {}
