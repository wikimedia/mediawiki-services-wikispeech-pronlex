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

// ValidatorForSymbolSet Place holder function. This should be handeled via the database (T138479, T138480).
func ValidatorForSymbolSet(symbolSetName string) (validation.Validator, error) {
	if symbolSetName == "sv.se.nst-SAMPA" {
		return NewNSTDemoValidator()
	}
	return validation.Validator{}, fmt.Errorf("no validator is defined for symbol set: %v", symbolSetName)
}

// NewNSTDemoValidator is used for testing
func NewNSTDemoValidator() (validation.Validator, error) {
	symbolset, err := NewNSTSymbolSet()
	if err != nil {
		return validation.Validator{}, err
	}
	finalNostressNolongRe, err := ProcessTransRe(symbolset, "\\$ (nonsyllabic )*(@|A|E|I|O|U|u0|Y|{|9|n=|l=|n`=|l`=)( nonsyllabic)*$")
	if err != nil {
		return validation.Validator{}, err
	}
	primaryStressRe, err := ProcessTransRe(symbolset, "\"")
	if err != nil {
		return validation.Validator{}, err
	}
	syllabicRe, err := ProcessTransRe(symbolset, "^(\"\"|\"|%)? *(nonsyllabic )*syllabic( nonsyllabic)*( (\\$|-) (\"\"|\"|%)? *(nonsyllabic )*syllabic( nonsyllabic)*)*$")
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
			Name:    "final_nostress_nolong",
			Level:   "Warning",
			Message: "final syllable should normally be unstressed with short vowel",
			Re:      finalNostressNolongRe,
		},
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

// NewNSTSymbolSet is used for testing
func NewNSTSymbolSet() (symbolset.SymbolSet, error) {
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
	return symbolset.NewSymbolSet(name, symbols)
}
