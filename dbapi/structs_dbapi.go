package dbapi

import (
	"fmt"
	"io"
	"strings"
)

// TODO Lägga till bolska fält för 'not'?
// Kunna sätta sortering eller ej?

// Query represents an sql search query to the lexicon database
type Query struct {
	// Lexicons to be searched. Empty means 'all' (TODO I think)
	Lexicons []Lexicon `json:"lexicons"`
	// list of words to get corresponding entries for
	Words []string `json:"words"`
	// a 'like' db search expression matching words
	WordLike   string `json:"wordLike"`
	WordRegexp string `json:"wordRegexp"`
	// a slice of Entry.IDs to search for
	EntryIDs []int64 `json:"entryIds"`
	// a 'like' db search expression matching transcriptions
	TranscriptionLike   string `json:"transcriptionLike"`
	TranscriptionRegexp string `json:"transcriptionRegexp"`
	// a 'like' db search expression matching part of speech strings
	PartOfSpeechLike   string `json:"partOfSpeechLike"`
	PartOfSpeechRegexp string `json:"partOfSpeechRegexp"`
	// list of lemma forms to get corresponding entries for
	Lemmas []string `json:"lemmas"`
	// an SQL 'like' expression to match lemma forms
	LemmaLike   string `json:"lemmaLike"`
	LemmaRegexp string `json:"lemmaRegexp"`
	// an SQL 'like' expression to match lemma readings
	ReadingLike   string `json:"readingLike"`
	ReadingRegexp string `json:"readingRegexp"`
	// an SQL 'like' expression to match lemma paradigms
	ParadigmLike   string `json:"paradigmLike"`
	ParadigmRegexp string `json:"paradigmRegexp"`
	// the page returned by the SQL query's 'LIMIT' (starts at 1)
	Page int64 `json:"page"`
	// the page length of the SQL query's 'LIMIT'
	PageLength int64 `json:"pageLength"`
}

// Empty returns true if there are not search criteria values
func (q Query) Empty() bool {
	switch {
	case len(q.Words) > 0:
		return false
	case strings.TrimSpace(q.WordLike) != "":
		return false
	case strings.TrimSpace(q.WordRegexp) != "":
		return false
	case len(q.EntryIDs) > 0:
		return false
	case strings.TrimSpace(q.TranscriptionLike) != "":
		return false
	case strings.TrimSpace(q.TranscriptionRegexp) != "":
		return false
	case strings.TrimSpace(q.PartOfSpeechLike) != "":
		return false
	case strings.TrimSpace(q.PartOfSpeechRegexp) != "":
		return false
	case len(q.Lemmas) > 0:
		return false
	case strings.TrimSpace(q.LemmaLike) != "":
		return false
	case strings.TrimSpace(q.LemmaRegexp) != "":
		return false
	case strings.TrimSpace(q.ReadingLike) != "":
		return false
	case strings.TrimSpace(q.ReadingRegexp) != "":
		return false
	case strings.TrimSpace(q.ParadigmLike) != "":
		return false
	case strings.TrimSpace(q.ParadigmRegexp) != "":
		return false
	}

	//log.Printf("dbapi.EmptyQuery: query struct appears to lack any search constraint: %v", q)
	return true
}

// NewQuery returns a Query instance where PageLength: 25
func NewQuery() Query {
	return Query{PageLength: 25}
}

// LexiconIDs returns a list of db IDs of the Lexicons of the Query
func (q Query) LexiconIDs() []int64 {
	var ids []int64
	for _, l := range q.Lexicons {
		ids = append(ids, l.ID)
	}
	return ids
}

// Lexicon corresponds to the lexicon db table, to which Entries are
// associated
type Lexicon struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	SymbolSetName string `json:"symbolSetName"`
}

// Transcription corresponds to the transcription db table
type Transcription struct {
	ID       int64  `json:"id"`
	EntryID  int64  `json:"entryId"`
	Strn     string `json:"strn"`
	Language string `json:"language"`
}

// TranscriptionSlice is used for
// soring according to ascending id
type TranscriptionSlice []Transcription

func (a TranscriptionSlice) Len() int           { return len(a) }
func (a TranscriptionSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TranscriptionSlice) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Entry defines a lexical entry. It does not correspond one-to-one to
// the entry db table, since it contains data also from associated
// tabled (Lemma, Transcription)
type Entry struct {
	ID             int64           `json:"id"`
	LexiconID      int64           `json:"lexiconId"`
	Strn           string          `json:"strn"`
	Language       string          `json:"language"`
	PartOfSpeech   string          `json:"partOfSpeech"`
	WordParts      string          `json:"wordParts"`
	Lemma          Lemma           `json:"lemma"`
	Transcriptions []Transcription `json:"transcriptions"`
}

// EntryWriter is an interface defining things to which one can write an Entry.
// See EntrySliceWriter, for returning i sice of Entry, and EntryFileWriter, for writing Entries to file.
type EntryWriter interface {
	Write(Entry) error
}

// EntryFileWriter outputs formated entries to an io.Writer.
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

// EntrySliceWriter is a container for returning Entries from a LookUp call to the db
// Example usage:
//	var q := dbapi.Query{ ... }
//	var esw dbapi.EntrySliceWriter
//	err := dbapi.LookUp(db, q, &esw)
//	[...] esw.Entries // process Entries
type EntrySliceWriter struct {
	Entries []Entry
}

func (w *EntrySliceWriter) Write(e Entry) error {
	w.Entries = append(w.Entries, e)
	return nil // fmt.Errorf("not implemented")
}

// Lemma corresponds to a row of the lemma db table
type Lemma struct {
	ID       int64  `json:"id"` // Är noll ett pålitligt 'None'-värde? Dvs börjar databaser alltid räkna från 1?
	Strn     string `json:"strn"`
	Reading  string `json:"reading"`
	Paradigm string `json:"paradigm"`
}

// Symbol corresponds to the symbol db table, and holds a phonetic
// symbol
type Symbol struct {
	LexiconID   int64  `json:"lexiconId"`
	Symbol      string `json:"symbol"`
	Category    string `json:"category"`
	Subcat      string `json:"subcat"`
	Description string `json:"description"`
	IPA         string `json:"ipa"`
}
