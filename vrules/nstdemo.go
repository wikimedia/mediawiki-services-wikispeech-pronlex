package vrules

import (
	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
)

// NewNSTDemoValidator is used for testing
func NewNSTDemoValidator() (validation.Validator, error) {
	symbolset, err := symbolset.NewNSTSymbolSet()
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
