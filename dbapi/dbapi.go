package dbapi

//go get github.com/mattn/go-sqlite3

import (
	"database/sql"
	"errors"
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

func ListLexicons(db *sql.DB) []Lexicon {
	res := make([]Lexicon, 0)
	sql := "select id, name, symbolsetname from lexicon"
	rows, err := db.Query(sql)
	ff("Query failed %v", err)
	defer rows.Close()

	for rows.Next() {
		var id int64
		var name string
		var symbolSetName string
		rows.Scan(&id, &name, &symbolSetName)
		l := Lexicon{Id: id, Name: name, SymbolSetName: symbolSetName}
		res = append(res, l)
	}

	return res
}

func GetLexicons(db *sql.DB, names []string) []Lexicon {
	res := make([]Lexicon, 0)
	if 0 == len(names) {
		return res
	}

	var id int64
	var lname string
	var symbolsetname string

	rows, err := db.Query("select id, name, symbolsetname from lexicon where name in "+nQs(len(names)), convS(names)...)
	ff("Failed DB query for lexicon names:\t%v", err)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &lname, &symbolsetname)
		ff("Scanning rows went wrong: %v", err)
		res = append(res, Lexicon{Id: id, Name: lname, SymbolSetName: symbolsetname})
	}

	return res
}

func GetLexicon(db *sql.DB, name string) (Lexicon, error) {
	var id int64
	var lname string
	var symbolsetname string

	if "" == strings.TrimSpace(name) {
		return Lexicon{}, errors.New("FAN!")
	}

	err := db.QueryRow("select id, name, symbolsetname from lexicon where name = ?", strings.ToLower(name)).Scan(&id, &lname, &symbolsetname)
	if err != nil {
		//log.Fatalf("DISASTER: %s", err)
		return Lexicon{}, errors.New("ZATAN")
	}

	return Lexicon{Id: id, Name: lname, SymbolSetName: symbolsetname}, nil

}

// TODO change input arg to sql.Tx ?
func InsertLexicon(db *sql.DB, l Lexicon) Lexicon {
	tx, err := db.Begin()
	defer tx.Commit()

	res, err := tx.Exec("insert into lexicon (name, symbolsetname) values (?, ?)", strings.ToLower(l.Name), l.SymbolSetName)
	if err != nil {
		log.Fatal("FAILED TO INSERT Lexicon (name + symbolset name): ", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	tx.Commit()

	return Lexicon{Id: id, Name: l.Name, SymbolSetName: l.SymbolSetName}
}

func InsertOrGetLexicon(db *sql.DB, l Lexicon) Lexicon {
	s := "select id, symbolsetname from lexicon where name = ?"

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Commit()

	rows, err := tx.Query(s, l.Name)
	defer rows.Close()
	if err != nil {
		log.Fatal("HEJ: ", err)
	}

	datum := make([]Lexicon, 0)
	for rows.Next() {
		var id int64
		var symbolsetname string
		if err := rows.Scan(&id); err != nil {
			log.Fatal(err)
		}
		if err := rows.Scan(&symbolsetname); err != nil {
			log.Fatal(err)
		}
		datum = append(datum, Lexicon{Id: id, Name: l.Name, SymbolSetName: symbolsetname})
	}

	tx.Commit()

	return l
}

// TODO change input arg to sql.Tx ?
func InsertEntries(db *sql.DB, l Lexicon, es []Entry) []int64 {

	var entrySTMT = "insert into entry (lexiconid, strn, language, partofspeech, wordparts) values (?, ?, ?, ?, ?)"
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
			log.Fatal(err) // TODO rollback?
		}

		id, err := res.LastInsertId()
		if err != nil {
			log.Fatal(err) // TODO rollback?
		}
		// We want the Entry to have the right id for inserting lemma assocs below
		e.Id = id

		ids = append(ids, id)

		// res.Close()

		for _, t := range e.Transcriptions {
			_, err := tx.Stmt(stmt2).Exec(id, t.Strn, t.Language)
			if err != nil {
				log.Fatal(err) // TODO rollback?
			}
		}

		//log.Printf("%v", e)
		if "" != e.Lemma.Strn && "" != e.Lemma.Reading {
			lemma, err := SetOrGetLemmaTx(tx, e.Lemma.Strn, e.Lemma.Reading, e.Lemma.Paradigm)
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
	ff("Wha? %v", err)
	return err
}

func SetOrGetLemmaTx(tx *sql.Tx, strn string, reading string, paradigm string) (Lemma, error) {
	res := Lemma{}

	var id int64
	var strn0, reading0, paradigm0 string
	sqlS := "select id, strn, reading, paradigm from lemma where strn = ? and reading = ?"
	err := tx.QueryRow(sqlS, strn, reading).Scan(&id, &strn0, &reading0, &paradigm0)
	switch {
	case err == sql.ErrNoRows:
		return InsertLemma(tx, Lemma{Id: id, Strn: strn, Reading: reading, Paradigm: paradigm})
	case err != nil:
		ff("SetOrGetLemma failed: %v", err)
	}

	res.Id = id
	res.Strn = strn0
	res.Reading = reading0
	res.Paradigm = paradigm0

	return res, err
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
		ff("SetOrGetLemma failed: %v", err)
	}

	res.Id = id
	res.Strn = strn0
	res.Reading = reading0
	res.Paradigm = paradigm0

	return res, err
}

func getLemmaFromEntryId(tx *sql.Tx, id int64) Lemma {
	res := Lemma{}
	sql := "select lemma.id, lemma.strn, lemma.reading, lemma.paradigm from entry, lemma, lemma2entry where " +
		"entry.id = ? and entry.id = lemma2entry.entryid and lemma.id = lemma2entry.lemmaid"
	var lId int64
	var strn, reading, paradigm string
	err := tx.QueryRow(sql, id).Scan(&lId, &strn, &reading, &paradigm)
	// TODO: why doesn't this work? Must be some silly mistake
	// "./dbapi.go:281: sql.ErrNoRows undefined (type string has no field or method ErrNoRows)"
	// switch {
	// case err == sql.ErrNoRows:
	// 	// No row:
	// 	// Silently return empty Lemma below
	// case err != nil:
	// 	ff("getLemmaFromENtryId: %v", err)
	// }
	_ = err

	// TODO Now silently returns empty lemma if nothing returned from tx. Ok?
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
	ff("InsertLemma tx.Exec: %v", err)
	id, err := res.LastInsertId()
	ff("InsertLemma LastInsertId: %v", err)
	l.Id = id
	return l, err
}

// func InsertLemma(db *sql.DB, l Lemma) (Lemma, error) {
// 	sql := "insert into lemma (strn, reading, paradigm) values (?, ?, ?)"
// 	res, err := db.Exec(sql, l.Strn, l.Reading, l.Paradigm)
// 	ff("InsertLemma db.Exec: %v", err)
// 	id, err := res.LastInsertId()
// 	ff("InsertLemma LastInsertId: %v", err)
// 	l.Id = id
// 	return l, err
// }

// TODO Map gör så att ordnigen blir galen

func EntriesFromIds(tx *sql.Tx, entryIds []int64) map[string][]Entry {
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
				Lemma:        getLemmaFromEntryId(tx, entryId),
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

func GetEntries(tx *sql.Tx, q Query) map[string][]Entry {
	res := make(map[string][]Entry)
	if q.Empty() { // TODO report to client?
		log.Printf("dbapi.GetEntries: Query empty of search constraints: %v", q)
		return res
	}

	qString, vs := idiotSql(q)

	rows, err := tx.Query(qString, vs...)
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

	res = EntriesFromIds(tx, ids)
	return res
}

func UpdateEntry(tx *sql.Tx, e Entry) error {
	return nil
}

func Nothing() {}
