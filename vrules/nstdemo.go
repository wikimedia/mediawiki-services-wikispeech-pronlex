package vrules

import (
	"fmt"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
)

// // ValidatorForLexicon Place holder function. This should be handeled via the database (T138479, T138480).
// func ValidatorForLexicon(lexName string) (validation.Validator, error) {
// 	if lexName == "sv.se.nst" {
// 		return NewNSTDemoValidator()
// 	}
// 	return validation.Validator{}, fmt.Errorf("no validator is defined for lexicon: %v", lexName)
// }

// After some thinking, we decided that a validator should be linked to a SymbolSet rather than a Lexicon.
// (Different lexica could have same symbol set.)

// TODO add validation of legal transcription symbols

//TODO NL: A bit confusing that this lives in package vrules, when it returns a validation.Validator...?

// ValidatorForSymbolSet Place holder function. This should be handeled via the database (T138479, T138480).
func ValidatorForSymbolSet(symbolSetName string) (validation.Validator, error) {
	if symbolSetName == "sv.se.nst-SAMPA" {
		return NewNSTDemoValidator()
	}
	return validation.Validator{}, fmt.Errorf("no validator is defined for symbol set: %v", symbolSetName)
}

// NewNSTDemoValidator is used for testing
func NewNSTDemoValidator() (validation.Validator, error) {
	symbolset, err := NewNSTSvHardWired()
	if err != nil {
		return validation.Validator{}, err
	}
	primaryStressRe, err := ProcessTransRe(symbolset, "\"")
	if err != nil {
		return validation.Validator{}, err
	}
	syllabicRe, err := ProcessTransRe(symbolset, "^(\"\"|\"|%)? *(nonsyllabic )*syllabic( nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic )*syllabic( nonsyllabic)*)*$")
	if err != nil {
		return validation.Validator{}, err
	}

	reFrom, err := regexp2.Compile("(.)\\1[+]\\1", regexp2.None)
	if err != nil {
		return validation.Validator{}, err
	}
	decomp2Orth := Decomp2Orth{"+", func(s string) (string, error) {
		res, err := reFrom.Replace(s, "$1+$1", 0, -1)
		if err != nil {
			return s, err
		}
		return res, nil
	}}

	var vali = validation.Validator{[]validation.Rule{
		MustHaveTrans{},
		NoEmptyTrans{},
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
		decomp2Orth,
		SymbolSetRule{
			SymbolSet: symbolset,
		},
	}}
	return vali, nil
}

// SvNSTHardWired is a temporary function that should not be used in production
func NewNSTSvHardWired() (symbolset.Symbols, error) {
	name := "sv.se.nst-SAMPA"

	syms := []symbolset.Symbol{
		symbolset.Symbol{Desc: "sil", String: "i:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "sill", String: "I", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "full", String: "u0", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "ful", String: "}:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "matt", String: "a", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "mat", String: "A:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "bot", String: "u:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "bott", String: "U", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "häl", String: "E:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "häll", String: "E", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "aula", String: "au", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "syl", String: "y:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "syll", String: "Y", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "hel", String: "e:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "herr,hett", String: "e", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "nöt", String: "2:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "mött,förra", String: "9", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "mål", String: "o:", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "moll,håll", String: "O", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "bättre", String: "@", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "europa", String: "eu", Cat: symbolset.Syllabic},
		symbolset.Symbol{Desc: "pol", String: "p", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "bok", String: "b", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "tok", String: "t", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "bort", String: "rt", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "mod", String: "m", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "nod", String: "n", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "dop", String: "d", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "bord", String: "rd", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "fot", String: "k", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "våt", String: "g", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "lång", String: "N", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "forna", String: "rn", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "fot", String: "f", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "våt", String: "v", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "kjol", String: "C", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "fors", String: "rs", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "rov", String: "r", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "lov", String: "l", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "sot", String: "s", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "sjok", String: "x", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "hot", String: "h", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "porla", String: "rl", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "jord", String: "j", Cat: symbolset.NonSyllabic},
		symbolset.Symbol{Desc: "syllable delimiter", String: ".", Cat: symbolset.SyllableDelimiter},
		symbolset.Symbol{Desc: "accent I", String: `"`, Cat: symbolset.Stress},
		symbolset.Symbol{Desc: "accent II", String: `""`, Cat: symbolset.Stress},
		symbolset.Symbol{Desc: "secondary stress", String: "%", Cat: symbolset.Stress},
		symbolset.Symbol{Desc: "phoneme delimiter", String: " ", Cat: symbolset.PhonemeDelimiter},
		symbolset.Symbol{"+", symbolset.CompoundDelimiter, ""},
	}

	return symbolset.NewSymbols(name, syms)

}

func NewNSTNbvHardWired() (symbolset.Symbols, error) {
	name := "NST nob sampa"
	symbols := []symbolset.Symbol{
		symbolset.Symbol{"@", symbolset.Syllabic, ""},
		symbolset.Symbol{"A", symbolset.Syllabic, ""},
		symbolset.Symbol{"E", symbolset.Syllabic, ""},
		symbolset.Symbol{"I", symbolset.Syllabic, ""},
		symbolset.Symbol{"O", symbolset.Syllabic, ""},
		symbolset.Symbol{"U", symbolset.Syllabic, ""},
		symbolset.Symbol{"u0", symbolset.Syllabic, ""},
		symbolset.Symbol{"Y", symbolset.Syllabic, ""},
		symbolset.Symbol{"\\{", symbolset.Syllabic, ""},
		symbolset.Symbol{"9", symbolset.Syllabic, ""},
		symbolset.Symbol{"A:", symbolset.Syllabic, ""},
		symbolset.Symbol{"e:", symbolset.Syllabic, ""},
		symbolset.Symbol{"i:", symbolset.Syllabic, ""},
		symbolset.Symbol{"o:", symbolset.Syllabic, ""},
		symbolset.Symbol{"u:", symbolset.Syllabic, ""},
		symbolset.Symbol{"\\}:", symbolset.Syllabic, ""},
		symbolset.Symbol{"y:", symbolset.Syllabic, ""},
		symbolset.Symbol{"{:", symbolset.Syllabic, ""},
		symbolset.Symbol{"2:", symbolset.Syllabic, ""},
		symbolset.Symbol{"9:", symbolset.Syllabic, ""},
		symbolset.Symbol{"\\{\\*I", symbolset.Syllabic, ""},
		symbolset.Symbol{"9\\*Y", symbolset.Syllabic, ""},
		symbolset.Symbol{"A\\*I", symbolset.Syllabic, ""},
		symbolset.Symbol{"E\\*\\}", symbolset.Syllabic, ""},
		symbolset.Symbol{"O\\*Y", symbolset.Syllabic, ""},
		symbolset.Symbol{"o~", symbolset.Syllabic, ""},
		symbolset.Symbol{"n=", symbolset.Syllabic, ""},
		symbolset.Symbol{"l=", symbolset.Syllabic, ""},
		symbolset.Symbol{"n`=", symbolset.Syllabic, ""},
		symbolset.Symbol{"l`=", symbolset.Syllabic, ""},
		symbolset.Symbol{"}*I", symbolset.Syllabic, ""},
		symbolset.Symbol{"a\\*U", symbolset.Syllabic, ""},
		symbolset.Symbol{"@\\*U", symbolset.Syllabic, ""},
		symbolset.Symbol{"e~", symbolset.Syllabic, ""},
		symbolset.Symbol{"3:", symbolset.Syllabic, ""},
		symbolset.Symbol{"a", symbolset.Syllabic, ""},
		symbolset.Symbol{"a:", symbolset.Syllabic, ""},
		symbolset.Symbol{"U:", symbolset.Syllabic, ""},
		symbolset.Symbol{"V", symbolset.Syllabic, ""},
		symbolset.Symbol{"U4", symbolset.Syllabic, ""},
		symbolset.Symbol{"I@", symbolset.Syllabic, ""},

		symbolset.Symbol{"p", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"t", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"k", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"b", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"d", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"g", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"f", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"v", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"h", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"j", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"s", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"l", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"r", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"n", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"m", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"N", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"t`", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"d`", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"s`", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"n`", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"l`", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"S", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"C", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"tS", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"dZ", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"w", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"x", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"T", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"D", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"r3", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"Z", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"z", symbolset.NonSyllabic, ""},

		symbolset.Symbol{"%", symbolset.Stress, ""},
		symbolset.Symbol{"\"\"", symbolset.Stress, ""},
		symbolset.Symbol{"\"", symbolset.Stress, ""},
		symbolset.Symbol{" ", symbolset.PhonemeDelimiter, ""},
		symbolset.Symbol{"$", symbolset.SyllableDelimiter, ""},
		symbolset.Symbol{"-", symbolset.CompoundDelimiter, ""},
	}
	return symbolset.NewSymbols(name, symbols)
}
