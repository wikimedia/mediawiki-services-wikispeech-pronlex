package dbapi

import (
	"fmt"
	"io"
	"strings"
)

// TODO Lägga till bolska fält för 'not'?
// Kunna sätta sortering eller ej?

// Query represents an sql search query to the lexicon database
// TODO Change to list(s) of search critieria.
// TODO add boolean for include/exclude (i.e., "NOT" in the generated SQL).
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
// This is no longer a sane way to do it, since the number of search criteria has grown.
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

// Lemma corresponds to a row of the lemma db table
type Lemma struct {
	ID       int64  `json:"id"`
	Strn     string `json:"strn"`
	Reading  string `json:"reading"`
	Paradigm string `json:"paradigm"`
}

// EntryStatus associates a status to an Entry. The status has a name (such as 'ok') and a source (a string identifying who or what generated the status)
type EntryStatus struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Source string `json:"source"`
	//EntryID int64  `json:"entryId"`
	//Timestamp int64  `json:"timestamp"`
	Timestamp string `json:"timestamp"`
	Current   bool   `json:"current"`
}

// EntryValidation associates a validation result to an Entry
type EntryValidation struct {
	ID int64 `json:"id"`

	// Lower case name of level of severity
	Level     string `json:"level"`
	Name      string `json:"name"`
	Message   string `json:"Message"`
	Timestamp string `json:"timestamp"`
}

// SourceDelimiter is used to split a string of sevaral sources into a slice
var SourceDelimiter = " : "

// Transcription corresponds to the transcription db table
type Transcription struct {
	ID       int64    `json:"id"`
	EntryID  int64    `json:"entryId"`
	Strn     string   `json:"strn"`
	Language string   `json:"language"`
	Sources  []string `json:"sources"`
}

// AddSource ... adds a source string at the beginning of the
// Transcription.Sources slice. If the source is already present,
// AddSource silently ignores to add the already existing
// source. AddSource returns an error when the input string contains the
// SourceDelimiter string.
func (t *Transcription) AddSource(s string) error {
	sDC := strings.ToLower(strings.TrimSpace(s))
	if strings.Contains(sDC, SourceDelimiter) {
		return fmt.Errorf("cannot add source containing the source delimiter : '%v'", SourceDelimiter)
	}

	for i := range t.Sources {
		if sDC == t.Sources[i] {
			return nil // source already there
		}
	}

	t.Sources = append([]string{sDC}, t.Sources...)

	return nil
}

// SourcesString returns the []string items of Transcription.Sources as a string, where the items are delimited by SourceDelimiter
func (t Transcription) SourcesString() string {
	return strings.Join(t.Sources, SourceDelimiter)
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
	ID               int64             `json:"id"`
	LexiconID        int64             `json:"lexiconId"`
	Strn             string            `json:"strn"`
	Language         string            `json:"language"`
	PartOfSpeech     string            `json:"partOfSpeech"`
	WordParts        string            `json:"wordParts"`
	Lemma            Lemma             `json:"lemma"`
	Transcriptions   []Transcription   `json:"transcriptions"`
	EntryStatus      EntryStatus       `json:"status"` // TODO Probably should be a slice of statuses?
	EntryValidations []EntryValidation `json:"entryValidations"`
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

// Symbol corresponds to the symbol db table, and holds a phonetic
// symbol
type Symbol struct {
	LexiconID   int64  `json:"lexiconId"`
	Symbol      string `json:"symbol"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IPA         string `json:"ipa"`
}

// TODO add fields for additional stats
type LexStats struct {
	LexiconID int64 `json:lexiconId`
	// The number of entries in the lexicon corresponding to database id LexiconID
	Entries int64 `json:entries`
}
