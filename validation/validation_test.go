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
func test_processTransRe(SymbolSet symbolset.SymbolSet, Regexp string) (*regexp2.Regexp, error) {
	Regexp = strings.Replace(Regexp, "nonsyllabic", SymbolSet.NonSyllabicRe.String(), -1)
	Regexp = strings.Replace(Regexp, "syllabic", SymbolSet.SyllabicRe.String(), -1)
	Regexp = strings.Replace(Regexp, "phoneme", SymbolSet.PhonemeRe.String(), -1)
	Regexp = strings.Replace(Regexp, "symbol", SymbolSet.SymbolRe.String(), -1)
	return regexp2.Compile(Regexp, regexp2.None)
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

func (r RequiredTransRe) Validate(e lex.Entry) (Result, error) {
	var messages = make([]string, 0)
	for _, t := range e.Transcriptions {
		if m, err := r.Re.MatchString(strings.TrimSpace(t.Strn)); !m {
			if err != nil {
				return Result{RuleName: r.Name(), Level: r.Level()}, err
			}
			messages = append(
				messages,
				fmt.Sprintf("%s. Found: /%s/", r.Message, t.Strn))
		}
	}
	return Result{RuleName: r.Name(), Level: r.Level(), Messages: messages}, nil
}

func (r RequiredTransRe) ShouldAccept() []lex.Entry {
	return r.Accept
}
func (r RequiredTransRe) ShouldReject() []lex.Entry {
	return r.Reject
}
func (r RequiredTransRe) Name() string {
	return r.NameStr
}
func (r RequiredTransRe) Level() string {
	return r.LevelStr
}

func test_createEntry(orth string, transes []string) lex.Entry {
	ts := []lex.Transcription{}
	for _, t0 := range transes {
		ts = append(ts, lex.Transcription{Strn: t0, Language: ""})
	}

	e := lex.Entry{Strn: orth,
		PartOfSpeech:   "",
		Morphology:     "",
		WordParts:      orth,
		Language:       "",
		Transcriptions: ts,
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	return e
}

func test_createEntries() []lex.Entry {
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

func test_createValidator() Validator {
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

	primaryStressRe, err := test_processTransRe(ss, "\"")
	ff("%v", err)

	syllabicRe, err := test_processTransRe(ss, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	ff("%v", err)

	var v = Validator{
		Name: ss.Name,
		Rules: []Rule{
			RequiredTransRe{
				NameStr:  "primary_stress",
				LevelStr: "Fatal",
				Message:  "Primary stress required",
				Re:       primaryStressRe,
			},
			RequiredTransRe{
				NameStr:  "syllabic",
				LevelStr: "Format",
				Message:  "Each syllable needs a syllabic phoneme",
				Re:       syllabicRe,
			},
		}}
	return v
}

func Test_ValidateEntry1(t *testing.T) {
	v := test_createValidator()
	es0 := test_createEntries()

	var eVals1 []string // := make([]string, 0)
	var es []lex.Entry  // es := make([]lex.Entry, 0)
	for _, e := range es0 {
		v.ValidateEntry(&e)
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

func test_createValidatorForTestingTestSuite() Validator {
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

	primaryStressRe, err := test_processTransRe(ss, "\"")
	ff("%v", err)

	syllabicRe, err := test_processTransRe(ss, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	ff("%v", err)

	var v = Validator{
		Name: ss.Name,
		Rules: []Rule{
			RequiredTransRe{
				NameStr:  "primary_stress",
				LevelStr: "Fatal",
				Message:  "Primary stress required",
				Re:       primaryStressRe,
				Accept:   []lex.Entry{test_createEntry("apa", []string{"\" A: . p a"})},
				Reject:   []lex.Entry{test_createEntry("apa", []string{"A: . p a"})},
			},
			RequiredTransRe{
				NameStr:  "syllabic",
				LevelStr: "Format",
				Message:  "Each syllable needs a syllabic phoneme",
				Re:       syllabicRe,
				Accept:   []lex.Entry{test_createEntry("apa", []string{"\" A: . p a"})},
				Reject:   []lex.Entry{test_createEntry("apa", []string{"\" A: . p . a"})},
			},
		}}
	return v
}

func test_createValidatorForTestingTestSuite_Invalid() Validator {
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

	primaryStressRe, err := test_processTransRe(ss, "\"")
	ff("%v", err)

	syllabicRe, err := test_processTransRe(ss, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	ff("%v", err)

	var v = Validator{
		Name: ss.Name,
		Rules: []Rule{
			RequiredTransRe{
				NameStr:  "primary_stress",
				LevelStr: "Fatal",
				Message:  "Primary stress required",
				Re:       primaryStressRe,
				Accept:   []lex.Entry{test_createEntry("apa", []string{"\" A: . p a"})},
				Reject:   []lex.Entry{test_createEntry("apa", []string{"A: . p a"})},
			},
			RequiredTransRe{
				NameStr:  "syllabic",
				LevelStr: "Format",
				Message:  "Each syllable needs a syllabic phoneme",
				Re:       syllabicRe,
				Accept:   []lex.Entry{test_createEntry("apa", []string{"A: . p a"})},
				Reject:   []lex.Entry{test_createEntry("apa", []string{"\" A: . p a"})},
			},
		}}
	return v
}

func Test_TestSuite_Valid(t *testing.T) {
	v := test_createValidatorForTestingTestSuite()
	res, err := v.RunTests()
	ff("%v", err)
	if res.Size() > 0 {
		t.Errorf("Expected validator suite to test without errors. Found: %#v", res)
	}
}

func Test_TestSuite_Invalid(t *testing.T) {
	v := test_createValidatorForTestingTestSuite_Invalid()
	res, err := v.RunTests()
	ff("%v", err)
	if len(res.RejectErrors) != 1 || len(res.CrossErrors) != 1 || res.Size() == 0 {
		t.Errorf("Didn't expect validator suite to test without errors. Found %#v", res)
	}
}

func test_createInvalidEntries() []lex.Entry {
	t1a := lex.Transcription{Strn: "\" A: p a", Language: "sv-se"}
	t1b := lex.Transcription{Strn: "\" a p a", Language: "sv-se"}

	e1 := lex.Entry{Strn: "apa",
		ID:             1,
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apa",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t1a, t1b},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	t2a := lex.Transcription{Strn: "\" A: p a n", Language: "sv-se"}
	t2b := lex.Transcription{Strn: "A p a n", Language: "sv-se"}

	e2 := lex.Entry{Strn: "apan",
		ID:             2,
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apan",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t2a, t2b},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	t3a := lex.Transcription{Strn: "\" . p a n", Language: "sv-se"}
	e3 := lex.Entry{Strn: "appan",
		ID:             3,
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "appan",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t3a},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	return []lex.Entry{e1, e2, e3}
}

func Test_ValidateEntriesWithPointer(t *testing.T) {
	v := test_createValidator()
	es := test_createInvalidEntries()

	var resVals []string
	res, _ := v.validateEntriesWithPointer(es)
	for _, e := range res {
		for _, v := range e.EntryValidations {
			resVals = append(resVals, v.String())
		}
	}

	var expectVals = []string{
		`Fatal|primary_stress: Primary stress required. Found: /A p a n/`,
		`Format|syllabic: Each syllable needs a syllabic phoneme. Found: /A p a n/`,
		`Format|syllabic: Each syllable needs a syllabic phoneme. Found: /" . p a n/`,
	}

	if !reflect.DeepEqual(resVals, expectVals) {
		t.Errorf(fs, expectVals, resVals)
	}

	if len(resVals) != 3 {
		t.Errorf(fs, "3", len(resVals))
	}

}
