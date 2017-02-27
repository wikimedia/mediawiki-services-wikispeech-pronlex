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

func (r RequiredTransRe) Validate(e lex.Entry) (Result, error) {
	var messages = make([]string, 0)
	for _, t := range e.Transcriptions {
		if m, err := r.Re.MatchString(strings.TrimSpace(t.Strn)); !m {
			if err != nil {
				return Result{RuleName: r.Name, Level: r.Level}, err
			} else {
				messages = append(
					messages,
					fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn))
			}
		}
	}
	return Result{RuleName: r.Name, Level: r.Level, Messages: messages}, nil
}

func (r RequiredTransRe) ShouldAccept() []lex.Entry {
	return make([]lex.Entry, 0)
}
func (r RequiredTransRe) ShouldReject() []lex.Entry {
	return make([]lex.Entry, 0)
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
		symbolset.Symbol{String: "a", Cat: symbolset.Syllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "A:", Cat: symbolset.Syllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "b", Cat: symbolset.NonSyllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "p", Cat: symbolset.NonSyllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "N", Cat: symbolset.NonSyllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "n", Cat: symbolset.NonSyllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: " ", Cat: symbolset.PhonemeDelimiter, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: ".", Cat: symbolset.SyllableDelimiter, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "\"", Cat: symbolset.Stress, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "\"\"", Cat: symbolset.Stress, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
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

	var eVals1 []string // := make([]string, 0)
	var es []lex.Entry  // es := make([]lex.Entry, 0)
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
	var eVals2 []string // eVals2 := make([]string, 0)
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
