package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/symbolset"
)

type SymbolSetRule struct {
	SymbolSet symbolset.SymbolSet
}

func (r SymbolSetRule) Validate(e dbapi.Entry) []Result {
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		splitted, err := r.SymbolSet.SplitTranscription(t.Strn)
		if err != nil {
			result = append(result, Result{"SymbolSet", "Fatal", fmt.Sprintf("Couldn't split transcription: /%s/", t.Strn)})
		} else {
			for _, symbol := range splitted {
				if !r.SymbolSet.ValidSymbol(symbol) {
					result = append(result, Result{"SymbolSet", "Fatal", fmt.Sprintf("Invalid transcription symbol: %s in /%s/", symbol, t.Strn)})
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

type IllegalTransRe struct {
	Name    string
	Level   string
	Message string
	Re      *regexp.Regexp
}

func (r IllegalTransRe) Validate(e dbapi.Entry) []Result {
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		if r.Re.MatchString(strings.TrimSpace(t.Strn)) {
			result = append(result, Result{r.Name, r.Level, fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn)})
		}
	}
	return result
}

type RequiredTransRe struct {
	Name    string
	Level   string
	Message string
	Re      *regexp.Regexp
}

func (r RequiredTransRe) Validate(e dbapi.Entry) []Result {
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		if !r.Re.MatchString(strings.TrimSpace(t.Strn)) {
			result = append(result, Result{r.Name, r.Level, fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn)})
		}
	}
	return result
}

type MustHaveTrans struct {
}

func (r MustHaveTrans) Validate(e dbapi.Entry) []Result {
	name := "MustHaveTrans"
	level := "Format"
	var result = make([]Result, 0)
	if len(e.Transcriptions) == 0 {
		result = append(result, Result{name, level, "At least one transcription is required"})
	}
	return result
}

type NoEmptyTrans struct {
}

func (r NoEmptyTrans) Validate(e dbapi.Entry) []Result {
	name := "NoEmptyTrans"
	level := "Format"
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		if len(strings.TrimSpace(t.Strn)) == 0 {
			result = append(result, Result{name, level, "Empty transcriptions are not allowed"})
		}
	}
	return result
}

func NewDecomp2Orth(SS symbolset.SymbolSet, PreFilterWordPartString func(string) string) (Decomp2Orth, error) {
	compDelims := symbolset.FilterSymbolsByType(SS.Symbols, []symbolset.SymbolType{symbolset.CompoundDelimiter})
	if len(compDelims) > 0 {
		compDelim := compDelims[0].String
		return Decomp2Orth{compDelim, PreFilterWordPartString}, nil
	}
	return Decomp2Orth{}, fmt.Errorf("no compound delimiter in symbol set")
}

type Decomp2Orth struct {
	compDelim               string
	preFilterWordPartString func(string) string
}

func (r Decomp2Orth) Validate(e dbapi.Entry) []Result {
	name := "Decomp2Orth"
	level := "Fatal"
	var result = make([]Result, 0)
	filteredWordParts := r.preFilterWordPartString(e.WordParts)
	expectOrth := strings.Replace(filteredWordParts, r.compDelim, "", -1)
	if expectOrth != e.Strn {
		result = append(result, Result{name, level, fmt.Sprintf("decomp/orth mismatch: %s/%s", e.WordParts, e.Strn)})
	}
	return result
}
