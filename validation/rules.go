package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/symbolset"
)

type SymbolSetRule struct {
	Name      string
	Level     string
	Message   string
	SymbolSet symbolset.SymbolSet
}

func (r SymbolSetRule) Validate(e dbapi.Entry) []Result {
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		for _, symbol := range r.SymbolSet.SplitTranscription(t.Strn) {
			panic("not implemented")
			//if r.Re.MatchString(strings.TrimSpace(t.Strn)) {
			//	result = append(result, Result{r.Name, r.Level, r.Message})
			//}
		}
	}
	return result
}

type IllegalTransRe struct {
	Name    string
	Level   string
	Message string
	Re      regexp.Regexp
}

func (r IllegalTransRe) Validate(e dbapi.Entry) []Result {
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		if r.Re.MatchString(strings.TrimSpace(t.Strn)) {
			result = append(result, Result{r.Name, r.Level, r.Message})
		}
	}
	return result
}

type RequiredTransRe struct {
	Name    string
	Level   string
	Message string
	Re      regexp.Regexp
}

func (r RequiredTransRe) Validate(e dbapi.Entry) []Result {
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		if !r.Re.MatchString(strings.TrimSpace(t.Strn)) {
			result = append(result, Result{r.Name, r.Level, r.Message})
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

type Decomp2Orth struct {
}

func (r Decomp2Orth) Validate(e dbapi.Entry) []Result {
	name := "Decomp2Orth"
	level := "Fatal"
	var result = make([]Result, 0)
	expectOrth := strings.Replace(e.WordParts, "+", "", -1)
	if expectOrth != e.Strn {
		result = append(result, Result{name, level, fmt.Sprintf("decomp/orth mismatch: %s/%s", e.WordParts, e.Strn)})
	}
	return result
}
