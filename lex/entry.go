package lex

import (
	"fmt"
	"io"
	"strings"
)

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
	RuleName  string `json:"ruleName"`
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

// Lemma corresponds to a row of the lemma db table
type Lemma struct {
	ID       int64  `json:"id"`
	Strn     string `json:"strn"`
	Reading  string `json:"reading"`
	Paradigm string `json:"paradigm"`
}

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
// See EntrySliceWriter, for returning a slice of Entry, and EntryFileWriter, for writing Entries to file.
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
