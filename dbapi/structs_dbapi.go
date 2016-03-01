package dbapi

import (
	"strings"
)

// TODO Lägga till bolska fält för 'not'?
// Kunna sätta sortering eller ej?
// Nu är det en inbyggd order by entry.strn, men det kanske segar ner stora sökningar?

// Query represents an sql search query to the lexicon database
type Query struct {
	Lexicons          []Lexicon `json:"lexicons"`          // Lexicons to be searched. Empty means 'all' (TODO I think)
	Words             []string  `json:"words"`             // list of words to get corresponding entries for
	WordLike          string    `json:"wordLike"`          // a 'like' db search expression matching words
	TranscriptionLike string    `json:"transcriptionLike"` // a 'like' db search expression matching transcriptions
	PartOfSpeechLike  string    `json:"partOfSpeechLike"`  // a 'like' db search expression matching part of speech strings
	Lemmas            []string  `json:"lemmas"`            // list of lemma forms to get corresponding entries for
	LemmaLike         string    `json:"lemmaLike"`
	ReadingLike       string    `json:"readingLike"`
	ParadigmLike      string    `json:"paradigmLike"`
	Page              int64     `json:"page"`
	PageLength        int64     `json:"pageLength"`
}

func (q Query) Empty() bool {
	switch {
	case len(q.Words) > 0:
		return false
	case strings.TrimSpace(q.WordLike) != "":
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

// ??? Klumpigt försök med defaultvärden
func NewQuery() Query {
	return Query{PageLength: 25}
}

func (q Query) LexiconIds() []int64 {
	ids := make([]int64, 0)
	for _, l := range q.Lexicons {
		ids = append(ids, l.Id)
	}
	return ids
}

// TODO Id int64 Är noll ett pålitligt 'None'-värde? Dvs börjar databaser alltid räkna från 1?
// TODO kolla efter motsvarighet till Option el dyl. Kolla "New[Struct...]"

type Lexicon struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	SymbolSetName string `json:"symbolSetName"`
}

type Transcription struct {
	Id       int64  `json:"id"`
	EntryId  int64  `json:"entryId"`
	Strn     string `json:"strn"`
	Language string `json:"language"`
}

// Sort according to ascending id
type TranscriptionSlice []Transcription

func (a TranscriptionSlice) Len() int           { return len(a) }
func (a TranscriptionSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TranscriptionSlice) Less(i, j int) bool { return a[i].Id < a[j].Id }

type Entry struct {
	Id             int64           `json:"id"`
	LexiconId      int64           `json:"lexiconId"`
	Strn           string          `json:"strn"`
	Language       string          `json:"language"`
	PartOfSpeech   string          `json:"partOfSpeech"`
	WordParts      string          `json:"wordParts"`
	Lemma          Lemma           `json:"lemma"`
	Transcriptions []Transcription `json:"transcriptions"`
}

type Lemma struct {
	Id       int64  `json:"id"` // Är noll ett pålitligt 'None'-värde? Dvs börjar databaser alltid räkna från 1?
	Strn     string `json:"strn"`
	Reading  string `json:"reading"`
	Paradigm string `json:"paradigm"`
}
