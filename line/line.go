// Package line is used to define lexicon line format for parsing input and printing output
package line

import (
	"fmt"
	"strings"
)

// BUG(hanna) Add field property to indicate whether a field is required or optional (or just use two maps: RequiredFields and OptionalFields).
// BUG(hanna) Should unknown content generate an error? Depends on the required/optional property (above).

// Field is a simple const for line field definition types
type Field int

//go:generate stringer -type=Field

const (
	// Orth orthography
	Orth Field = iota

	// Pos part-of-speech (noun, verb, NN, VB, etc)
	Pos

	// Morph morphological tags (case, gender, tense, etc)
	Morph

	// WordParts decompounded orthography field (for compounds)
	WordParts

	// Lang the word's language
	Lang

	// Trans1 the primary transcription
	Trans1

	// Translang1 the language of the primary transcription
	Translang1

	// Trans2 transcription variant
	Trans2

	// Translang2 language for Trans2
	Translang2

	// Trans3 transcription variant
	Trans3

	// Translang3 language for Trans3
	Translang3

	// Trans4 transcription variant
	Trans4

	// Translang4 language for Trans4
	Translang4

	// Trans5 transcription variant
	Trans5

	// Translang5 language for Trans5
	Translang5

	// Trans6 transcription variant
	Trans6

	// Translang6 language for Trans6
	Translang6

	// Lemma the lemma form. Ttypically orthographic lemmma + some kind of (disambiguation) identifier, eg., wind_01.
	Lemma

	// Paradigm rule reference (id) for generating inflected forms from lemma
	Paradigm
)

// FormatTest defines a test to run upon initialization of Format (using NewFormat)
type FormatTest struct {
	InputLine  string
	Fields     map[Field]string
	OutputLine string
}

// Parser is used to define a lexicon's line parser.
// To implement your own parser, make sure to implement functions Parse(string) and String(map[Field]string)
type Parser interface {

	// Parse is used for parsing input lines
	Parse(string) (map[Field]string, error)

	// String is used to generate an output line from a set of fields
	String(map[Field]string) (string, error)
}

// Format is used to define a lexicon's line.
// This a struct for package private usage.
// To create a new Format instance, use NewFormat.
type Format struct {
	Name     string
	FieldSep string
	Fields   map[Field]int
	NFields  int
}

// NewFormat is a public constructor for Format with built-in error checks and tests
func NewFormat(name string, fieldSep string, fields map[Field]int, nFields int, tests []FormatTest) (Format, error) {
	f := Format{name, fieldSep, fields, nFields}
	var errs = make([]string, 0)
	for _, t := range tests {
		fieldsRes, err := f.Parse(t.InputLine)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%v", err))
		} else if !equals(fieldsRes, t.Fields) {
			errs = append(errs, fmt.Sprintf("Format.Parse: expected %v, found %v", t.Fields, fieldsRes))
		}
		lineRes, err := f.String(t.Fields)
		if err != nil {
			errs = append(errs, fmt.Sprintf("%v", err))
		} else if lineRes != t.OutputLine {
			errs = append(errs, fmt.Sprintf("Format.String: expected %v, found %v", t.OutputLine, lineRes))
		}
	}
	if len(errs) > 0 {
		return Format{}, fmt.Errorf(strings.Join(errs, " : "))
	}
	return f, nil
}

// Parse is used for parsing input lines
func (f Format) Parse(line string) (map[Field]string, error) {
	inputFields := strings.Split(line, f.FieldSep)
	if len(inputFields) != f.NFields {
		return make(map[Field]string), fmt.Errorf("expected %v fields, found %v, %s", f.NFields, len(inputFields), line)
	}
	var res = make(map[Field]string)
	for field, i := range f.Fields {
		if i < len(inputFields) {
			res[field] = inputFields[i]
		} else {
			return make(map[Field]string), fmt.Errorf("couldn't find field %v with index %v in line %v", field, i, line)
		}
	}
	return res, nil
}

// String is used to generate an output line from a set of fields
func (f Format) String(fields map[Field]string) (string, error) {
	var res = make([]string, f.NFields)
	for field, s := range fields {
		i, ok := f.Fields[field]
		if ok {
			res[i] = s
		} else {
			return "", fmt.Errorf("undefined field name: %v ", field.String())
		}
	}
	return strings.Join(res, f.FieldSep), nil
}
