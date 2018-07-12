package dbapi

import (
	"strings"

	"github.com/stts-se/pronlex/lex"
)

// DBMQuery is a query used by the DBManager, containing lexicon referenes (db+lex name) and a dbapi.Query
type DBMQuery struct {
	LexRefs []lex.LexRef
	Query   Query
}

// type LexiconQuery struct {
// 	Lexicons []string
// 	Query    Query
// }

// TODO Lägga till bolska fält för 'not'?
// Kunna sätta sortering eller ej?

// Query represents an sql search query to the lexicon database
// TODO Change to list(s) of search critieria.
// TODO add boolean for include/exclude (i.e., "NOT" in the generated SQL).
type Query struct {
	// list of words to get corresponding entries for
	Words []string `json:"words"`
	// a 'like' db search expression matching words
	WordLike   string `json:"wordLike"`
	WordRegexp string `json:"wordRegexp"`

	WordParts       []string `json:"wordParts"`
	WordPartsLike   string   `json:"wordPartsLike"`
	WordPartsRegexp string   `json:"wordPartsRegexp"`

	// a slice of Entry.IDs to search for
	EntryIDs []int64 `json:"entryIds"`
	// a 'like' db search expression matching transcriptions
	TranscriptionLike   string `json:"transcriptionLike"`
	TranscriptionRegexp string `json:"transcriptionRegexp"`
	// a 'like' db search expression matching part of speech strings
	PartOfSpeechLike   string `json:"partOfSpeechLike"`
	PartOfSpeechRegexp string `json:"partOfSpeechRegexp"`

	MorphologyLike string `json:"morphologyLike"`

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

	TagLike      string `json:"tagLike"`
	LanguageLike string `json:"languageLike"`

	CommentLabelLike  string `json:"commentLabelLike"`
	CommentSourceLike string `json:"commenSourceLike"`
	CommentLike       string `json:"commentLike"`

	// A list of entry statuses to match
	EntryStatus []string `json:"entryStatus"`

	// A list of users to match
	Users []string `json:"user"`

	// Select entries with one or more EntryValidations
	HasEntryValidation bool `json:"hasEntryValidation"`

	// // Search for Entries with EntryValidations with the listed
	// // validation rule names (such as 'Decomp2Orth', etc)
	// EntryValidations []string `json:"entryValidations"`

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
	case len(q.WordParts) > 0:
		return false
	case strings.TrimSpace(q.WordPartsLike) != "":
		return false
	case strings.TrimSpace(q.WordPartsRegexp) != "":
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
	case strings.TrimSpace(q.LanguageLike) != "":
		return false
	case strings.TrimSpace(q.MorphologyLike) != "":
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
	case strings.TrimSpace(q.TagLike) != "":
		return false
	case strings.TrimSpace(q.CommentLabelLike) != "":
		return false
	case strings.TrimSpace(q.CommentSourceLike) != "":
		return false
	case strings.TrimSpace(q.CommentLike) != "":
		return false
	case len(q.EntryStatus) > 0:
		return false
	case len(q.Users) > 0:
		return false
	case q.HasEntryValidation:
		return false

	}

	//log.Printf("dbapi.EmptyQuery: query struct appears to lack any search constraint: %v", q)
	return true
}

// NewQuery returns a Query instance where PageLength: 0
func NewQuery() Query {
	//return Query{PageLength: 25}
	return Query{PageLength: 0}
}

// LexiconIDs returns a list of db IDs of the Lexicons of the Query
// func (q Query) LexiconIDs() []int64 {
// 	var ids []int64
// 	for _, l := range q.Lexicons {
// 		ids = append(ids, l.ID)
// 	}
// 	return ids
// }

// Lexicon corresponds to the lexicon db table, to which Entries are
// associated
type lexicon struct {
	id            int64  // `json:"id"`
	name          string // `json:"name"`
	symbolSetName string // `json:"symbolSetName"`
	locale        string // `json:"locale"`
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

type StatusFreq struct {
	Status string `json:"status"`
	Freq   int64  `json:"freq"`
}

// LexStats holds the result of a call to the dbapi.LexiconStats function.
type LexStats struct {
	Lexicon string `json:"lexicon"`
	// The number of entries in the lexicon corresponding to database id LexiconID
	Entries int64 `json:"entries"`

	StatusFrequencies []StatusFreq `json:"statusFrequencies"`

	ValStats ValStats
}

// QueryStats holds the result of a call to the dbapi.LexiconStats function.
// TODO add fields for additional stats
type QueryStats struct {
	Query   Query `json:"query"`
	Entries int64 `json:"entries"`
}

// ValStats is used to incrementally give statistics during a validation process, or to just represent a final validation statistics.
type ValStats struct {
	// TotalEntries is the total entries to be validated
	TotalEntries int

	// ValidatedEntries is the total validated entries so far
	ValidatedEntries int

	// TotalValidations is the total number of validation messages so far
	TotalValidations int

	// InvalidEntries is the number of invalid entries so far
	InvalidEntries int

	Levels map[string]int `json:"levels"`
	Rules  map[string]int `json:"rules"`
}
