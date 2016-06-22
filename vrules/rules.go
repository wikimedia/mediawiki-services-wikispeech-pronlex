package vrules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
)

// SymbolSetRule is a general rule for verifying that each phoneme is a legal symbol
type SymbolSetRule struct {
	SymbolSet symbolset.SymbolSet
}

func (r SymbolSetRule) Validate(e lex.Entry) []validation.Result {
	var result = make([]validation.Result, 0)
	for _, t := range e.Transcriptions {
		splitted, err := r.SymbolSet.SplitTranscription(t.Strn)
		if err != nil {
			result = append(result, validation.Result{"SymbolSet", "Fatal", fmt.Sprintf("Couldn't split transcription: /%s/", t.Strn)})
		} else {
			for _, symbol := range splitted {
				if !r.SymbolSet.ValidSymbol(symbol) {
					result = append(result, validation.Result{"SymbolSet", "Fatal", fmt.Sprintf("Invalid transcription symbol: %s in /%s/", symbol, t.Strn)})
				}
			}
		}
	}
	return result
}

/*
ProcessTransRe converts pre-defined entities to the appropriate symbols. Strings replaced are: syllabic, nonsyllabic, phoneme, symbol.
*/
func ProcessTransRe(SymbolSet symbolset.SymbolSet, Regexp string) (*regexp.Regexp, error) {
	Regexp = strings.Replace(Regexp, "nonsyllabic", SymbolSet.NonSyllabicRe.String(), -1)
	Regexp = strings.Replace(Regexp, "syllabic", SymbolSet.SyllabicRe.String(), -1)
	Regexp = strings.Replace(Regexp, "phoneme", SymbolSet.PhonemeRe.String(), -1)
	Regexp = strings.Replace(Regexp, "symbol", SymbolSet.SymbolRe.String(), -1)
	return regexp.Compile(Regexp)
}

// IllegalTransRe is a general rule type to check for illegal transcriptions by regexp
type IllegalTransRe struct {
	Name    string
	Level   string
	Message string
	Re      *regexp.Regexp
}

func (r IllegalTransRe) Validate(e lex.Entry) []validation.Result {
	var result = make([]validation.Result, 0)
	for _, t := range e.Transcriptions {
		if r.Re.MatchString(strings.TrimSpace(t.Strn)) {
			result = append(result, validation.Result{r.Name, r.Level, fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn)})
		}
	}
	return result
}

// RequiredTransRe is a general rule type used to defined basic transcription requirements using regexps
type RequiredTransRe struct {
	Name    string
	Level   string
	Message string
	Re      *regexp.Regexp
}

func (r RequiredTransRe) Validate(e lex.Entry) []validation.Result {
	var result = make([]validation.Result, 0)
	for _, t := range e.Transcriptions {
		if !r.Re.MatchString(strings.TrimSpace(t.Strn)) {
			result = append(result, validation.Result{r.Name, r.Level, fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn)})
		}
	}
	return result
}

// MustHaveTrans is a general rule to make sure each entry has at least one transcription
type MustHaveTrans struct {
}

func (r MustHaveTrans) Validate(e lex.Entry) []validation.Result {
	name := "MustHaveTrans"
	level := "Format"
	var result = make([]validation.Result, 0)
	if len(e.Transcriptions) == 0 {
		result = append(result, validation.Result{name, level, "At least one transcription is required"})
	}
	return result
}

// NoEmptyTrans is a general rule to make sure no transcriptions are be empty
type NoEmptyTrans struct {
}

func (r NoEmptyTrans) Validate(e lex.Entry) []validation.Result {
	name := "NoEmptyTrans"
	level := "Format"
	var result = make([]validation.Result, 0)
	for _, t := range e.Transcriptions {
		if len(strings.TrimSpace(t.Strn)) == 0 {
			result = append(result, validation.Result{name, level, "Empty transcriptions are not allowed"})
		}
	}
	return result
}

// Decomp2Orth is a general rule type to validate the word parts vs. the orthography. A filter is used to control the filtering, typically how to treat triple consonants at boundaries.
type Decomp2Orth struct {
	compDelim               string
	preFilterWordPartString func(string) (string, error)
}

func (r Decomp2Orth) Validate(e lex.Entry) []validation.Result {
	name := "Decomp2Orth"
	level := "Fatal"
	var result = make([]validation.Result, 0)
	filteredWordParts, err := r.preFilterWordPartString(e.WordParts)
	if err != nil {
		result = append(result, validation.Result{name, level, fmt.Sprintf("decomp/orth rule returned error on replace call: %v", err)})
	}
	expectOrth := strings.Replace(filteredWordParts, r.compDelim, "", -1)
	if expectOrth != e.Strn {
		result = append(result, validation.Result{name, level, fmt.Sprintf("decomp/orth mismatch: %s/%s", e.WordParts, e.Strn)})
	}
	return result
}
