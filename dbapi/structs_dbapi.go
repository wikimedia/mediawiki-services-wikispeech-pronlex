package dbapi

import (
	"strings"
)

// TODO Lägga till bolska fält för 'not'?
// Kunna sätta sortering eller ej?
// Nu är det en inbyggd order by entry.strn, men det kanske segar ner stora sökningar?

// Query represents an sql search query to the lexicon database
type Query struct {
	Lexicons          []Lexicon `json:"lexicons"` // Lexicons to be searched. Empty means 'all' (TODO I think)
	Words             []string  `json:"words"`    // list of words to get corresponding entries for
	WordLike          string    `json:"wordLike"` // a 'like' db search expression matching words
	EntryIDs          []int64   `json:"entryIds"`
	TranscriptionLike string    `json:"transcriptionLike"` // a 'like' db search expression matching transcriptions
	PartOfSpeechLike  string    `json:"partOfSpeechLike"`  // a 'like' db search expression matching part of speech strings
	Lemmas            []string  `json:"lemmas"`            // list of lemma forms to get corresponding entries for
	LemmaLike         string    `json:"lemmaLike"`
	ReadingLike       string    `json:"readingLike"`
	ParadigmLike      string    `json:"paradigmLike"`
	Page              int64     `json:"page"`
	PageLength        int64     `json:"pageLength"`
}

// Empty returns true if there are not search criteria values
func (q Query) Empty() bool {
	switch {
	case len(q.Words) > 0:
		return false
	case strings.TrimSpace(q.WordLike) != "":
		return false
	case len(q.EntryIDs) > 0:
		return false
	case strings.TrimSpace(q.TranscriptionLike) != "":
		return false
	case strings.TrimSpace(q.PartOfSpeechLike) != "":
		return false
	case len(q.Lemmas) > 0:
		return false
	case strings.TrimSpace(q.LemmaLike) != "":
		return false
	case strings.TrimSpace(q.ReadingLike) != "":
		return false
	case strings.TrimSpace(q.ParadigmLike) != "":
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
