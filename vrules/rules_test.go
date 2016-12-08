package vrules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
)

var fsExp = "Expected: '%v' got: '%v'"

type testMustHaveTrans struct {
}

func (r testMustHaveTrans) Validate(e lex.Entry) []validation.Result {
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

type testNoEmptyTrans struct {
}

func (r testNoEmptyTrans) Validate(e lex.Entry) []validation.Result {
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

type testDecomp2Orth struct {
}

func (r testDecomp2Orth) Validate(e lex.Entry) []validation.Result {
	name := "Decomp2Orth"
	level := "Fatal"
	var result = make([]validation.Result, 0)
	expectOrth := strings.Replace(e.WordParts, "+", "", -1)
	if expectOrth != e.Strn {
		result = append(result, validation.Result{
			RuleName: name,
			Level:    level,
			Message:  fmt.Sprintf("decomp/orth mismatch: %s/%s", e.WordParts, e.Strn)})
	}
	return result
}

func Test1(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testMustHaveTrans{}, testNoEmptyTrans{}}}

	var e = lex.Entry{
		Strn:         "anka",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "anka",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" a N . k a",
				Language: "swe",
			},
		},
	}

	var _, result = vali.ValidateEntries([]lex.Entry{e})

	if result != true {
		t.Errorf(fsExp, make([]validation.Result, 0), result)
	}
}

func Test2(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testMustHaveTrans{}, testNoEmptyTrans{}}}

	var e = lex.Entry{
		Strn:           "anka",
		Language:       "swe",
		PartOfSpeech:   "NN",
		WordParts:      "anka",
		Transcriptions: []lex.Transcription{},
	}

	e, _ = vali.ValidateEntry(e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "MustHaveTrans",
			Level:    "Format",
			Message:  "[...]",
		},
	}

	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	} else {
		if result[0].RuleName != "MustHaveTrans" {
			t.Errorf(fsExp, expect, result)
		}
	}
}

func Test3(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testMustHaveTrans{}, testNoEmptyTrans{}, testDecomp2Orth{}}}

	var e = lex.Entry{
		Strn:         "ankstjärt",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "ank+sjärt",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" a N k + % x { rt",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "Decomp2Orth",
			Level:    "Fatal",
			Message:  "[...]",
		},
	}
	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	} else {
		if result[0].RuleName != "Decomp2Orth" {
			t.Errorf(fsExp, expect, result)
		}
	}
}

func Test4(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testMustHaveTrans{}, testNoEmptyTrans{}, testDecomp2Orth{}}}

	var e = lex.Entry{
		Strn:         "ankstjärtsbad",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "ank+stjärts+bad",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" a N k + x { rt rs + % b A: d",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}
}

func TestNst1(t *testing.T) {
	symbolset, err := newNSTSvHardWired_ForTesting()
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	vali, err := newSvSeNstValidator(symbolset)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	var e = lex.Entry{
		Strn:         "banen",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "banen",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\" b A: . n @ n",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = lex.Entry{
		Strn:         "bantorget",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "ban+torget",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" b A: n + % t O r . j @ t",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	result = e.EntryValidations

	expect = []lex.EntryValidation{}

	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = lex.Entry{
		Strn:         "battorget",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "bat+torget",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" b A: t + % t O r . j @ t",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	result = e.EntryValidations

	expect = []lex.EntryValidation{}

	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = lex.Entry{
		Strn:         "battorget",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "batt+torget",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" b a t + % t O r . j @ t",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	result = e.EntryValidations

	expect = []lex.EntryValidation{}

	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = lex.Entry{
		Strn:         "batttorget",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "batt+torget",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" b a t + % t O r . j @ t",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	result = e.EntryValidations

	expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "Decomp2Orth",
			Level:    "Fatal",
			Message:  "[...]"},
	}

	if len(result) != len(expect) || result[0].RuleName != expect[0].RuleName {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = lex.Entry{
		Strn:         "apnos",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "ap+nos",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "Aa: p n u: s",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	result = e.EntryValidations

	expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "primary_stress",
			Level:    "Fatal",
			Message:  "[...]"},
		lex.EntryValidation{
			RuleName: "syllabic",
			Level:    "Format",
			Message:  "[...]"},
		lex.EntryValidation{
			RuleName: "symbolset",
			Level:    "Format",
			Message:  "[...]"},
		lex.EntryValidation{
			RuleName: "MaxOneSyllabic",
			Level:    "Fatal",
			Message:  "[...]"},
	}

	if len(result) != len(expect) || result[0].RuleName != expect[0].RuleName || result[1].RuleName != expect[1].RuleName {
		t.Errorf(fsExp, expect, result)
	}
}

func TestRegexp2Backrefs(t *testing.T) {
	reFrom, err := regexp2.Compile("(.)\\1[+]\\1", regexp2.None)
	if err != nil {
		t.Errorf("%v", err)
	}
	//

	reTo := "$1+$1"
	input := "hatt+torget"
	expect := "hat+torget"
	result, err := reFrom.Replace(input, reTo, 0, -1)
	if err != nil {
		t.Errorf("%v", err)
	}
	if result != expect {
		t.Errorf(fsExp, expect, result)
	}

	//

	input = "hats+torget"
	expect = "hats+torget"
	result, err = reFrom.Replace(input, reTo, 0, -1)
	if err != nil {
		t.Errorf("%v", err)
	}

	if result != expect {
		t.Errorf(fsExp, expect, result)
	}

}

// SvNSTHardWired is a temporary function that should not be used in production
func newNSTSvHardWired_ForTesting() (symbolset.SymbolSet, error) {
	name := "sv.se.nst-SAMPA"

	syms := []symbolset.Symbol{
		symbolset.Symbol{Desc: "sil", String: "i:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "sill", String: "I", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "full", String: "u0", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "ful", String: "}:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "matt", String: "a", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "mat", String: "A:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "bot", String: "u:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "bott", String: "U", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "häl", String: "E:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "häll", String: "E", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "aula", String: "au", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "syl", String: "y:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "syll", String: "Y", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "hel", String: "e:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "herr,hett", String: "e", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "nöt", String: "2:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "mött,förra", String: "9", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "mål", String: "o:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "moll,håll", String: "O", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "bättre", String: "@", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "europa", String: "eu", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "pol", String: "p", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "bok", String: "b", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "tok", String: "t", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "bort", String: "rt", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "mod", String: "m", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "nod", String: "n", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "dop", String: "d", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "bord", String: "rd", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "fot", String: "k", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "våt", String: "g", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "lång", String: "N", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "forna", String: "rn", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "fot", String: "f", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "våt", String: "v", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "kjol", String: "C", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "fors", String: "rs", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "rov", String: "r", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "lov", String: "l", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "sot", String: "s", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "sjok", String: "x", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "hot", String: "h", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "porla", String: "rl", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "jord", String: "j", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "syllable delimiter", String: ".", Cat: symbolset.SyllableDelimiter, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "accent I", String: `"`, Cat: symbolset.Stress, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "accent II", String: `""`, Cat: symbolset.Stress, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "secondary stress", String: "%", Cat: symbolset.Stress, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{Desc: "phoneme delimiter", String: " ", Cat: symbolset.PhonemeDelimiter, IPA: symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"+", symbolset.CompoundDelimiter, "", symbolset.IPASymbol{"", ""}},
	}

	return symbolset.NewSymbolSet(name, syms)

}

func newNSTNbvHardWired_ForTesting() (symbolset.SymbolSet, error) {
	name := "NST nob sampa"
	symbols := []symbolset.Symbol{
		symbolset.Symbol{"@", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"A", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"E", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"I", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"O", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"U", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"u0", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"Y", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"\\{", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"9", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"A:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"e:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"i:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"o:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"u:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"\\}:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"y:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"{:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"2:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"9:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"\\{\\*I", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"9\\*Y", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"A\\*I", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"E\\*\\}", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"O\\*Y", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"o~", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"n=", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"l=", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"n`=", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"l`=", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"}*I", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"a\\*U", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"@\\*U", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"e~", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"3:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"a", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"a:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"U:", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"V", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"U4", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"I@", symbolset.Syllabic, "", symbolset.IPASymbol{"", ""}},

		symbolset.Symbol{"p", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"t", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"k", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"b", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"d", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"g", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"f", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"v", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"h", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"j", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"s", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"l", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"r", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"n", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"m", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"N", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"t`", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"d`", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"s`", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"n`", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"l`", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"S", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"C", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"tS", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"dZ", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"w", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"x", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"T", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"D", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"r3", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"Z", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"z", symbolset.NonSyllabic, "", symbolset.IPASymbol{"", ""}},

		symbolset.Symbol{"%", symbolset.Stress, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"\"\"", symbolset.Stress, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"\"", symbolset.Stress, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{" ", symbolset.PhonemeDelimiter, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"$", symbolset.SyllableDelimiter, "", symbolset.IPASymbol{"", ""}},
		symbolset.Symbol{"-", symbolset.CompoundDelimiter, "", symbolset.IPASymbol{"", ""}},
	}
	return symbolset.NewSymbolSet(name, symbols)
}

func TestWhitespace(t *testing.T) {
	symbolset, err := newNSTSvHardWired_ForTesting()
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	vali, err := newSvSeNstValidator(symbolset)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	var e = lex.Entry{
		Strn:         "banen",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "banen",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\" b A: . n @ n",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = lex.Entry{
		Strn:         "banen",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "banen",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\" b A: . n @ n ",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	result = e.EntryValidations

	expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "SymbolSet",
			Level:    "Fatal",
			Message:  "[...]"},
	}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = lex.Entry{
		Strn:         "banen",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "banen",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\" b A: . n @  n",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	result = e.EntryValidations

	expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "SymbolSet",
			Level:    "Fatal",
			Message:  "[...]"},
	}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//

}

func TestRepeated(t *testing.T) {
	symbolset, err := newNSTSvHardWired_ForTesting()
	if err != nil {
		t.Errorf("%s", err)
		return
	}
	vali, err := newSvSeNstValidator(symbolset)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	var e = lex.Entry{
		Strn:         "banen",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "banen",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\" b A: . n @ n",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = lex.Entry{
		Strn:         "banen",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "banen",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\" b A: n . n @ n",
				Language: "swe",
			},
		},
	}

	e, _ = vali.ValidateEntry(e)
	result = e.EntryValidations

	expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "repeated_phonemes",
			Level:    "Fatal",
			Message:  "[...]"},
	}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//
}
