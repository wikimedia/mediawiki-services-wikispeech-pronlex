package dbapi

import (
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

// Symbol corresponds to the symbol db table, and holds a phonetic
// symbol
type Symbol struct {
	LexiconID   int64  `json:"lexiconId"`
	Symbol      string `json:"symbol"`
	Category    string `json:"category"`
	Description string `json:"description"`
	IPA         string `json:"ipa"`
}

// LexStats holds the result of a call to the dbapi.LexiconStats function.
// TODO add fields for additional stats
type LexStats struct {
	LexiconID int64 `json:"lexiconId"`
	// The number of entries in the lexicon corresponding to database id LexiconID
	Entries int64 `json:"entries"`

	// Status frequencies, as strings: StatusName<TAB>Frequency
	// TODO better structure for status/freq (string/int)
	StatusFrequencies []string `json:"statusFrequencies"`

	ValStats ValStats
}

// QueryStats holds the result of a call to the dbapi.LexiconStats function.
// TODO add fields for additional stats
type QueryStats struct {
	Query   Query `json:"query"`
	Entries int64 `json:"entries"`
}

type ValStats struct {
	InvalidEntries   int
	TotalEntries     int
	TotalValidations int
	Levels           map[string]int `json:levels`
	Rules            map[string]int `json:rules`
}
