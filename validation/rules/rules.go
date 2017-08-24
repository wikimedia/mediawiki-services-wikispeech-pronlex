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
func (r SymbolSetRule) Validate(e lex.Entry) (validation.Result, error) {
	var messages = make([]string, 0)
	for _, t := range e.Transcriptions {
		splitted, err := r.SymbolSet.SplitTranscription(t.Strn)
		if err != nil {
			return validation.Result{RuleName: "SymbolSet", Level: "Fatal"}, err
		}
		for _, symbol := range splitted {
			if !r.SymbolSet.ValidSymbol(symbol) {
				messages = append(
					messages,
					fmt.Sprintf("Invalid transcription symbol '%s' in /%s/", symbol, t.Strn))
			}
		}
	}
	return validation.Result{RuleName: r.Name(), Level: r.Level(), Messages: messages}, nil
}

// ShouldAccept returns a slice of entries that the rule should accept
func (r SymbolSetRule) ShouldAccept() []lex.Entry {
	return make([]lex.Entry, 0)
}

// ShouldReject returns a slice of entries that the rule should reject
func (r SymbolSetRule) ShouldReject() []lex.Entry {
	return make([]lex.Entry, 0)
}

// Name is the name of this rule
func (r SymbolSetRule) Name() string {
	return "SymbolSet"
}

// Level is the rule level (typically format, fatal, warning, info)
func (r SymbolSetRule) Level() string {
	return "Fatal"
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
	NameStr  string
	LevelStr string
	Message  string
	Re       *regexp2.Regexp
	Accept   []lex.Entry
	Reject   []lex.Entry
}

// Validate a lex.Entry
func (r IllegalTransRe) Validate(e lex.Entry) (validation.Result, error) {
	var messages = make([]string, 0)
	for _, t := range e.Transcriptions {
		if m, err := r.Re.MatchString(strings.TrimSpace(t.Strn)); m {
			if err != nil {
				return validation.Result{RuleName: r.Name(), Level: r.Level()}, err
			}
			messages = append(
				messages,
				fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn))
		}
	}
	return validation.Result{RuleName: r.Name(), Level: r.Level(), Messages: messages}, nil
}

// ShouldAccept returns a slice of entries that the rule should accept
func (r IllegalTransRe) ShouldAccept() []lex.Entry {
	return r.Accept
}

// ShouldReject returns a slice of entries that the rule should reject
func (r IllegalTransRe) ShouldReject() []lex.Entry {
	return r.Reject
}

// Name is the name of this rule
func (r IllegalTransRe) Name() string {
	return r.NameStr
}

// Level is the rule level (typically format, fatal, warning, info)
func (r IllegalTransRe) Level() string {
	return r.LevelStr
}

// RequiredTransRe is a general rule type used to defined basic transcription requirements using regexps
type RequiredTransRe struct {
	NameStr  string
	LevelStr string
	Message  string
	Re       *regexp2.Regexp
	Accept   []lex.Entry
	Reject   []lex.Entry
}

// Validate a lex.Entry
func (r RequiredTransRe) Validate(e lex.Entry) (validation.Result, error) {
	var messages = make([]string, 0)
	for _, t := range e.Transcriptions {
		if m, err := r.Re.MatchString(strings.TrimSpace(t.Strn)); !m {
			if err != nil {
				return validation.Result{RuleName: r.Name(), Level: r.Level()}, err
			}
			messages = append(
				messages,
				fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn))
		}
	}
	return validation.Result{RuleName: r.Name(), Level: r.Level(), Messages: messages}, nil
}

// ShouldAccept returns a slice of entries that the rule should accept
func (r RequiredTransRe) ShouldAccept() []lex.Entry {
	return r.Accept
}

// ShouldReject returns a slice of entries that the rule should reject
func (r RequiredTransRe) ShouldReject() []lex.Entry {
	return r.Reject
}

// Name is the name of this rule
func (r RequiredTransRe) Name() string {
	return r.NameStr
}

// Level is the rule level (typically format, fatal, warning, info)
func (r RequiredTransRe) Level() string {
	return r.LevelStr
}

// MustHaveTrans is a general rule to make sure each entry has at least one transcription
type MustHaveTrans struct {
}

// Validate a lex.Entry
func (r MustHaveTrans) Validate(e lex.Entry) (validation.Result, error) {
	var messages = make([]string, 0)
	if len(e.Transcriptions) == 0 {
		messages = append(messages, "At least one transcription is required")
	}
	return validation.Result{RuleName: r.Name(), Level: r.Level(), Messages: messages}, nil
}

// ShouldAccept returns a slice of entries that the rule should accept
func (r MustHaveTrans) ShouldAccept() []lex.Entry {
	return make([]lex.Entry, 0)
}

// ShouldReject returns a slice of entries that the rule should reject
func (r MustHaveTrans) ShouldReject() []lex.Entry {
	return make([]lex.Entry, 0)
}

// Name is the name of this rule
func (r MustHaveTrans) Name() string {
	return "MustHaveTrans"
}

// Level is the rule level (typically format, fatal, warning, info)
func (r MustHaveTrans) Level() string {
	return "Format"
}

// NoEmptyTrans is a general rule to make sure no transcriptions are be empty
type NoEmptyTrans struct {
}

// Validate a lex.Entry
func (r NoEmptyTrans) Validate(e lex.Entry) (validation.Result, error) {
	var messages = make([]string, 0)
	for _, t := range e.Transcriptions {
		if len(strings.TrimSpace(t.Strn)) == 0 {
			messages = append(messages, "Empty transcriptions are not allowed")
		}
	}
	return validation.Result{RuleName: r.Name(), Level: r.Level(), Messages: messages}, nil
}

// ShouldAccept returns a slice of entries that the rule should accept
func (r NoEmptyTrans) ShouldAccept() []lex.Entry {
	return make([]lex.Entry, 0)
}

// ShouldReject returns a slice of entries that the rule should reject
func (r NoEmptyTrans) ShouldReject() []lex.Entry {
	return make([]lex.Entry, 0)
}

// Name is the name of this rule
func (r NoEmptyTrans) Name() string {
	return "NoEmptyTrans"
}

// Level is the rule level (typically format, fatal, warning, info)
func (r NoEmptyTrans) Level() string {
	return "Format"
}

// Decomp2Orth is a general rule type to validate the word parts vs. the orthography. A filter is used to control the filtering, typically how to treat triple consonants at boundaries.
type Decomp2Orth struct {
	CompDelim               string
	AcceptEmptyDecomp       bool
	PreFilterWordPartString func(string) (string, error)
	Accept                  []lex.Entry
	Reject                  []lex.Entry
}

// Validate a lex.Entry
func (r Decomp2Orth) Validate(e lex.Entry) (validation.Result, error) {
	var messages = make([]string, 0)
	if r.AcceptEmptyDecomp && len(strings.TrimSpace(e.WordParts)) == 0 {
		return validation.Result{RuleName: r.Name(), Level: r.Level(), Messages: messages}, nil
	}
	filteredWordParts, err := r.PreFilterWordPartString(e.WordParts)
	if err != nil {
		return validation.Result{RuleName: "SymbolSet", Level: "Fatal"}, err
	}
	expectOrth := strings.Replace(filteredWordParts, r.CompDelim, "", -1)
	if expectOrth != e.Strn {
		messages = append(
			messages,
			fmt.Sprintf("decomp/orth mismatch: %s/%s", e.WordParts, e.Strn))
	}
	return validation.Result{RuleName: r.Name(), Level: r.Level(), Messages: messages}, nil
}

// ShouldAccept returns a slice of entries that the rule should accept
func (r Decomp2Orth) ShouldAccept() []lex.Entry {
	return r.Accept
}

// ShouldReject returns a slice of entries that the rule should reject
func (r Decomp2Orth) ShouldReject() []lex.Entry {
	return r.Accept
}

// Name is the name of this rule
func (r Decomp2Orth) Name() string {
	return "Decomp2Orth"
}

// Level is the rule level (typically format, fatal, warning, info)
func (r Decomp2Orth) Level() string {
	return "Fatal"
}
