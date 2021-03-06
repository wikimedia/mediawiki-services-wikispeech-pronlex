package validators

import (
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/pronlex/validation/rules"
	"github.com/stts-se/symbolset"
)

func newNbNoNstValidator(symbolset symbolset.SymbolSet) (validation.Validator, error) {
	primaryStressRe, err := rules.ProcessTransRe(symbolset, "\"")
	if err != nil {
		return validation.Validator{}, err
	}
	syllabicRe, err := rules.ProcessTransRe(symbolset, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	if err != nil {
		return validation.Validator{}, err
	}

	stressFirst, err := rules.ProcessTransRe(symbolset, "[^.!+ ] +(\"\"|\"|%)")
	if err != nil {
		return validation.Validator{}, err
	}

	reFrom, err := regexp2.Compile("(.)\\1[+]\\1", regexp2.None)
	if err != nil {
		return validation.Validator{}, err
	}
	decomp2Orth := rules.Decomp2Orth{CompDelim: "+",
		AcceptEmptyDecomp: true,
		PreFilterWordPartString: func(s string) (string, error) {
			res, err := reFrom.Replace(s, "$1+$1", 0, -1)
			res = strings.ToLower(strings.Replace(res, "!", "", -1))
			if err != nil {
				return s, err
			}
			return res, nil
		}}

	var vali = validation.Validator{
		Name: symbolset.Name,
		Rules: []validation.Rule{
			rules.MustHaveTrans{},
			rules.NoEmptyTrans{},
			rules.RequiredTransRe{
				NameStr:  "primary_stress",
				LevelStr: "Fatal",
				Message:  "Primary stress required",
				Re:       primaryStressRe,
			},
			rules.IllegalTransRe{
				NameStr:  "stress_first",
				LevelStr: "Fatal",
				Message:  "Stress can only be used in syllable initial position",
				Re:       stressFirst,
			},
			rules.RequiredTransRe{
				NameStr:  "syllabic",
				LevelStr: "Format",
				Message:  "Each syllable needs a syllabic phoneme",
				Re:       syllabicRe,
			},
			decomp2Orth,
			rules.SymbolSetRule{
				SymbolSet: symbolset,
			},
		}}
	return vali, nil
}
