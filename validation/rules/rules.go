package rules

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
)

// SymbolSetRule is a general rule for verifying that each phoneme is a legal symbol
type SymbolSetRule struct {
	SymbolSet symbolset.SymbolSet
}

// Validate a lex.Entry
func (r SymbolSetRule) Validate(e lex.Entry) []validation.Result {
	var result []validation.Result
	for _, t := range e.Transcriptions {
		splitted, err := r.SymbolSet.SplitTranscription(t.Strn)
		if err != nil {
			result = append(result, validation.Result{
				RuleName: "SymbolSet",
				Level:    "Format",
				Message:  fmt.Sprintf("Couldn't split transcription: /%s/", t.Strn)})
		} else {
			for _, symbol := range splitted {
				if !r.SymbolSet.ValidSymbol(symbol) {
					result = append(result, validation.Result{
						RuleName: "SymbolSet",
						Level:    "Fatal",
						Message:  fmt.Sprintf("Invalid transcription symbol '%s' in /%s/", symbol, t.Strn)})
				}
			}
		}
	}
	return result
}

/*
ProcessTransRe converts pre-defined entities to the appropriate symbols. Strings replaced are: syllabic, nonsyllabic, phoneme, symbol.
*/
func ProcessTransRe(SymbolSet symbolset.SymbolSet, Regexp string) (*regexp2.Regexp, error) {
	Regexp = strings.Replace(Regexp, "nonsyllabic", SymbolSet.NonSyllabicRe.String(), -1)
	Regexp = strings.Replace(Regexp, "syllabic", SymbolSet.SyllabicRe.String(), -1)
	Regexp = strings.Replace(Regexp, "phoneme", SymbolSet.PhonemeRe.String(), -1)
	Regexp = strings.Replace(Regexp, "symbol", SymbolSet.SymbolRe.String(), -1)
	return regexp2.Compile(Regexp, regexp2.None)
}

// IllegalTransRe is a general rule type to check for illegal transcriptions by regexp
type IllegalTransRe struct {
	Name    string
	Level   string
	Message string
	Re      *regexp2.Regexp
}

// Validate a lex.Entry
func (r IllegalTransRe) Validate(e lex.Entry) []validation.Result {
	var result = make([]validation.Result, 0)
	for _, t := range e.Transcriptions {
		if m, err := r.Re.MatchString(strings.TrimSpace(t.Strn)); m {
			if err != nil {
				result = append(result, validation.Result{
					RuleName: "System",
					Level:    "Format",
					Message:  fmt.Sprintf("error when validating rule %s on transcription string /%s/ : %v", r.Name, t.Strn, err)})
			} else {
				result = append(result, validation.Result{
					RuleName: r.Name,
					Level:    r.Level,
					Message:  fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn)})
			}
		}
	}
	return result
}

// RequiredTransRe is a general rule type used to defined basic transcription requirements using regexps
type RequiredTransRe struct {
	Name    string
	Level   string
	Message string
	Re      *regexp2.Regexp
}

// Validate a lex.Entry
func (r RequiredTransRe) Validate(e lex.Entry) []validation.Result {
	var result = make([]validation.Result, 0)
	for _, t := range e.Transcriptions {
		if m, err := r.Re.MatchString(strings.TrimSpace(t.Strn)); !m {
			if err != nil {
				result = append(result, validation.Result{
					RuleName: "System",
					Level:    "Format",
					Message:  fmt.Sprintf("error when validating rule %s on transcription string /%s/ : %v", r.Name, t.Strn, err)})
			} else {
				result = append(result, validation.Result{
					RuleName: r.Name,
					Level:    r.Level,
					Message:  fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn)})
			}
		}
	}
	return result
}

// MustHaveTrans is a general rule to make sure each entry has at least one transcription
type MustHaveTrans struct {
}

// Validate a lex.Entry
func (r MustHaveTrans) Validate(e lex.Entry) []validation.Result {
	name := "MustHaveTrans"
	level := "Format"
	var result = make([]validation.Result, 0)
	if len(e.Transcriptions) == 0 {
		result = append(result, validation.Result{
			RuleName: name,
			Level:    level,
			Message:  "At least one transcription is required"})
	}
	return result
}

// NoEmptyTrans is a general rule to make sure no transcriptions are be empty
type NoEmptyTrans struct {
}

// Validate a lex.Entry
func (r NoEmptyTrans) Validate(e lex.Entry) []validation.Result {
	name := "NoEmptyTrans"
	level := "Format"
	var result = make([]validation.Result, 0)
	for _, t := range e.Transcriptions {
		if len(strings.TrimSpace(t.Strn)) == 0 {
			result = append(result, validation.Result{
				RuleName: name,
				Level:    level,
				Message:  "Empty transcriptions are not allowed"})
		}
	}
	return result
}

// Decomp2Orth is a general rule type to validate the word parts vs. the orthography. A filter is used to control the filtering, typically how to treat triple consonants at boundaries.
type Decomp2Orth struct {
	CompDelim               string
	AcceptEmptyDecomp       bool
	PreFilterWordPartString func(string) (string, error)
}

// Validate a lex.Entry
func (r Decomp2Orth) Validate(e lex.Entry) []validation.Result {
	name := "Decomp2Orth"
	level := "Fatal"
	var result = make([]validation.Result, 0)
	if r.AcceptEmptyDecomp && len(strings.TrimSpace(e.WordParts)) == 0 {
		return result
	}
	filteredWordParts, err := r.PreFilterWordPartString(e.WordParts)
	if err != nil {
		result = append(result, validation.Result{
			RuleName: name,
			Level:    level,
			Message:  fmt.Sprintf("decomp/orth rule returned error on replace call: %v", err)})
	}
	expectOrth := strings.Replace(filteredWordParts, r.CompDelim, "", -1)
	if expectOrth != e.Strn {
		result = append(result, validation.Result{
			RuleName: name,
			Level:    level,
			Message:  fmt.Sprintf("decomp/orth mismatch: %s/%s", e.WordParts, e.Strn)})
	}
	return result
}