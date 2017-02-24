package validators

import (
	"testing"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
)

var fsExp = "Expected: '%v' got: '%v'"

// SvNSTHardWired is a temporary function that should not be used in production
func newNSTSvHardWired_ForTesting() (symbolset.SymbolSet, error) {
	name := "sv.se.nst-SAMPA"

	syms := []symbolset.Symbol{
		symbolset.Symbol{Desc: "sil", String: "i:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "sill", String: "I", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "full", String: "u0", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "ful", String: "}:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "matt", String: "a", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "mat", String: "A:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "bot", String: "u:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "bott", String: "U", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "häl", String: "E:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "häll", String: "E", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "aula", String: "au", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "syl", String: "y:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "syll", String: "Y", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "hel", String: "e:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "herr,hett", String: "e", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "nöt", String: "2:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "mött,förra", String: "9", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "mål", String: "o:", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "moll,håll", String: "O", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "bättre", String: "@", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "europa", String: "eu", Cat: symbolset.Syllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "pol", String: "p", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "bok", String: "b", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "tok", String: "t", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "bort", String: "rt", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "mod", String: "m", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "nod", String: "n", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "dop", String: "d", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "bord", String: "rd", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "fot", String: "k", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "våt", String: "g", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "lång", String: "N", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "forna", String: "rn", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "fot", String: "f", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "våt", String: "v", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "kjol", String: "C", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "fors", String: "rs", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "rov", String: "r", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "lov", String: "l", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "sot", String: "s", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "sjok", String: "x", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "hot", String: "h", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "porla", String: "rl", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "jord", String: "j", Cat: symbolset.NonSyllabic, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "syllable delimiter", String: ".", Cat: symbolset.SyllableDelimiter, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "accent I", String: `"`, Cat: symbolset.Stress, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "accent II", String: `""`, Cat: symbolset.Stress, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "secondary stress", String: "%", Cat: symbolset.Stress, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{Desc: "phoneme delimiter", String: " ", Cat: symbolset.PhonemeDelimiter, IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"+", symbolset.CompoundDelimiter, "", symbolset.IPASymbol{String: "", Unicode: ""}},
	}

	return symbolset.NewSymbolSet(name, syms)

}

func newNSTNbvHardWired_ForTesting() (symbolset.SymbolSet, error) {
	name := "NST nob sampa"
	symbols := []symbolset.Symbol{
		symbolset.Symbol{"@", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"A", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"E", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"I", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"O", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"U", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"u0", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"Y", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"\\{", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"9", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"A:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"e:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"i:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"o:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"u:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"\\}:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"y:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"{:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"2:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"9:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"\\{\\*I", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"9\\*Y", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"A\\*I", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"E\\*\\}", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"O\\*Y", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"o~", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"n=", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"l=", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"n`=", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"l`=", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"}*I", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"a\\*U", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"@\\*U", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"e~", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"3:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"a", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"a:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"U:", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"V", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"U4", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"I@", symbolset.Syllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},

		symbolset.Symbol{"p", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"t", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"k", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"b", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"d", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"g", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"f", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"v", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"h", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"j", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"s", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"l", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"r", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"n", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"m", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"N", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"t`", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"d`", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"s`", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"n`", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"l`", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"S", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"C", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"tS", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"dZ", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"w", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"x", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"T", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"D", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"r3", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"Z", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"z", symbolset.NonSyllabic, "", symbolset.IPASymbol{String: "", Unicode: ""}},

		symbolset.Symbol{"%", symbolset.Stress, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"\"\"", symbolset.Stress, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"\"", symbolset.Stress, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{" ", symbolset.PhonemeDelimiter, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"$", symbolset.SyllableDelimiter, "", symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{"-", symbolset.CompoundDelimiter, "", symbolset.IPASymbol{String: "", Unicode: ""}},
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
