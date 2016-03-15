package line

import (
	"fmt"
	"strings"
)

// BUG(hanna) Add field property to indicate whether a field is required or optional (or just use two maps: RequiredFields and OptionalFields).
// BUG(hanna) Should unknown content generate an error? Depends on the required/optional property (above).

// Field is a simple const for line field definition types
type Field int

const (
	// Orth orthography
	Orth Field = iota

	// Pos part-of-speech (noun, verb, NN, VB, etc)
	Pos

	// Morph morphological tags (case, gender, tense, etc)
	Morph

	// Decomp decompounded orthography field (for compounds)
	Decomp

	// WordLang the word's language
	WordLang

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

	// InflectionRule rule reference (id) for generating inflected forms from lemma
	InflectionRule
)

// Format a lexicon's line format definition
type Format struct {
	FieldSep string
	Fields   map[Field]int
	NFields  int
}

func (f Format) valueOf(index int) Field {
	for field, i := range f.Fields {
		if i == index {
			return field
		}
	}
	return -1
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

func (f Format) rightmostFieldIndex() int {
	var res = -1
	for _, i := range f.Fields {
		if i > res {
			res = i
		}
	}
	return res
}

// String is used to generate an output line from a set of fields
func (f Format) String(fields map[Field]string) (string, error) {
	max := f.rightmostFieldIndex()
	var res = make([]string, max+1)
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
