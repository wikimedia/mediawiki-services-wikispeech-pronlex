package validation

import (
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
)

var fs = "Wanted: '%v' got: '%v'"

// ff is a place holder to be replaced by proper error handling
func ff(f string, err error) {
	if err != nil {
		log.Fatalf(f, err)
	}
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

// RequiredTransRe is a general rule type used to defined basic transcription requirements using regexps
type RequiredTransRe struct {
	Name    string
	Level   string
	Message string
	Re      *regexp2.Regexp
}

func (r RequiredTransRe) Validate(e lex.Entry) []Result {
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		if m, err := r.Re.MatchString(strings.TrimSpace(t.Strn)); !m {
			if err != nil {
				result = append(result, Result{
					RuleName: "System",
					Level:    "Format",
					Message:  fmt.Sprintf("error when validating rule %s on transcription string /%s/ : %v", r.Name, t.Strn, err)})
			} else {
				result = append(result, Result{
					RuleName: r.Name,
					Level:    r.Level,
					Message:  fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn)})
			}
		}
	}
	return result
}

func createEntries() []lex.Entry {
	t1a := lex.Transcription{Strn: "\" A: p a", Language: "sv-se"}
	t1b := lex.Transcription{Strn: "\" a p a", Language: "sv-se"}

	e1 := lex.Entry{Strn: "apa",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apa",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t1a, t1b},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	t2a := lex.Transcription{Strn: "\" A: p a n", Language: "sv-se"}
	t2b := lex.Transcription{Strn: "A: p a n", Language: "sv-se"}

	e2 := lex.Entry{Strn: "apan",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apan",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t2a, t2b},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	t3a := lex.Transcription{Strn: "\" A . p a n", Language: "sv-se"}
	e3 := lex.Entry{Strn: "appan",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "appan",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t3a},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	return []lex.Entry{e1, e2, e3}
}

func createValidator() Validator {
	name := "sampa"
	symbols := []symbolset.Symbol{
		symbolset.Symbol{"a", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"A:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"b", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"p", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"N", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"n", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{" ", symbolset.PhonemeDelimiter, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{".", symbolset.SyllableDelimiter, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"\"", symbolset.Stress, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"\"\"", symbolset.Stress, "", symbolset.IPASymbol{"", ""}},
	}
	ss, err := symbolset.NewSymbolSet(name, symbols)
	ff("failed to init symbols : %v", err)

	primaryStressRe, err := ProcessTransRe(ss, "\"")
	ff("%v", err)

	syllabicRe, err := ProcessTransRe(ss, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	ff("%v", err)

	var v = Validator{
		Name: ss.Name,
		Rules: []Rule{
			RequiredTransRe{
				Name:    "primary_stress",
				Level:   "Fatal",
				Message: "Primary stress required",
				Re:      primaryStressRe,
			},
			RequiredTransRe{
				Name:    "syllabic",
				Level:   "Format",
				Message: "Each syllable needs a syllabic phoneme",
				Re:      syllabicRe,
			},
		}}
	return v
}

func Test_ValidateEntry1(t *testing.T) {
	v := createValidator()
	es0 := createEntries()

	eVals1 := make([]string, 0)
	es := make([]lex.Entry, 0)
	for _, e := range es0 {
		e, _ = v.ValidateEntry(e)
		for _, v := range e.EntryValidations {
			eVals1 = append(eVals1, v.String())
		}
		es = append(es, e)
	}
	sort.Strings(eVals1)
	if len(eVals1) < 1 {
		t.Errorf(fs, ">1", len(eVals1))
	}
	eVals2 := make([]string, 0)
	for _, e := range es {
		for _, v := range e.EntryValidations {
			eVals2 = append(eVals2, v.String())
		}
	}
	sort.Strings(eVals2)

	if !reflect.DeepEqual(eVals1, eVals2) {
		t.Errorf(fs, eVals1, eVals2)
	}

}
