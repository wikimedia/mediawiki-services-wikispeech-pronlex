package line

import (
	"fmt"
	"strings"
)

// Bug(hanna) Add field property to indicate whether a field is required or optional (or just use two maps: RequiredFields and OptionalFields).

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

// FieldName is a generated function to derive the name of a Field
func FieldName(f Field) (string, error) {
	switch f {
	case Orth:
		return "Orth", nil
	case Pos:
		return "Pos", nil
	case Morph:
		return "Morph", nil
	case Decomp:
		return "Decomp", nil
	case WordLang:
		return "WordLang", nil
	case Trans1:
		return "Trans1", nil
	case Translang1:
		return "Translang1", nil
	case Trans2:
		return "Trans2", nil
	case Translang2:
		return "Translang2", nil
	case Trans3:
		return "Trans3", nil
	case Translang3:
		return "Translang3", nil
	case Trans4:
		return "Trans4", nil
	case Translang4:
		return "Translang4", nil
	case Trans5:
		return "Trans5", nil
	case Translang5:
		return "Translang5", nil
	case Trans6:
		return "Trans6", nil
	case Translang6:
		return "Translang6", nil
	case Lemma:
		return "Lemma", nil
	case InflectionRule:
		return "InflectionRule", nil
	default:
		return "", fmt.Errorf("undefined field: %v", f)
	}
}

// Format a lexicon's line format definition
type Format struct {
	FieldSep string
	Fields   map[Field]int
}

// Parse function used for parsing input lines
func (f Format) Parse(line string) (map[Field]string, error) {
	fields := strings.Split(line, f.FieldSep)
	var res = make(map[Field]string)
	for field, i := range f.Fields {
		if i < len(fields) {
			res[field] = fields[i]
		} else {
			return make(map[Field]string), fmt.Errorf("Couldn't find field %v with index %v in line %v", field, i, line)
		}
	}
	return res, nil
}
