package rules

import (
	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
)

func newSvSeNstValidator(symbolset symbolset.SymbolSet) (validation.Validator, error) {
	primaryStressRe, err := ProcessTransRe(symbolset, "\"")
	if err != nil {
		return validation.Validator{}, err
	}
	syllabicRe, err := ProcessTransRe(symbolset, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	if err != nil {
		return validation.Validator{}, err
	}

	maxOneSyllabic, err := ProcessTransRe(symbolset, "syllabic[^.+%\"-]*( +syllabic)")
	if err != nil {
		return validation.Validator{}, err
	}

	reFrom, err := regexp2.Compile("(.)\\1[+]\\1", regexp2.None)
	if err != nil {
		return validation.Validator{}, err
	}
	decomp2Orth := Decomp2Orth{CompDelim: "+",
		AcceptEmptyDecomp: true,
		PreFilterWordPartString: func(s string) (string, error) {
			res, err := reFrom.Replace(s, "$1+$1", 0, -1)
			if err != nil {
				return s, err
			}
			return res, nil
		}}

	repeatedPhnRe, err := ProcessTransRe(symbolset, "symbol( +[.~])? +\\1")
	if err != nil {
		return validation.Validator{}, err
	}

	var vali = validation.Validator{
		Name: symbolset.Name,
		Rules: []validation.Rule{
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
			IllegalTransRe{
				Name:    "MaxOneSyllabic",
				Level:   "Fatal",
				Message: "A syllable cannot contain more than one syllabic phoneme",
				Re:      maxOneSyllabic,
			},
			IllegalTransRe{
				Name:    "repeated_phonemes",
				Level:   "Fatal",
				Message: "Repeated phonemes cannot be used within the same morpheme",
				Re:      repeatedPhnRe,
			},
			decomp2Orth,
			SymbolSetRule{
				SymbolSet: symbolset,
			},
		}}
	return vali, nil
}

func newNbNoNstValidator(symbolset symbolset.SymbolSet) (validation.Validator, error) {
	primaryStressRe, err := ProcessTransRe(symbolset, "\"")
	if err != nil {
		return validation.Validator{}, err
	}
	syllabicRe, err := ProcessTransRe(symbolset, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	if err != nil {
		return validation.Validator{}, err
	}

	reFrom, err := regexp2.Compile("(.)\\1[+]\\1", regexp2.None)
	if err != nil {
		return validation.Validator{}, err
	}
	decomp2Orth := Decomp2Orth{CompDelim: "+",
		AcceptEmptyDecomp: true,
		PreFilterWordPartString: func(s string) (string, error) {
			res, err := reFrom.Replace(s, "$1+$1", 0, -1)
			if err != nil {
				return s, err
			}
			return res, nil
		}}

	var vali = validation.Validator{
		Name: symbolset.Name,
		Rules: []validation.Rule{
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

func newEnUsCmuNstValidator(symbolset symbolset.SymbolSet) (validation.Validator, error) {
	exactOnePrimStressRe, err := ProcessTransRe(symbolset, "^[^\"]*\"[^\"]*$")
	if err != nil {
		return validation.Validator{}, err
	}
	maxOneSecStressRe, err := ProcessTransRe(symbolset, "%.*%")
	if err != nil {
		return validation.Validator{}, err
	}

	var vali = validation.Validator{
		Name: symbolset.Name,
		Rules: []validation.Rule{
			MustHaveTrans{},
			NoEmptyTrans{},
			RequiredTransRe{
				Name:    "primary_stress",
				Level:   "Format",
				Message: "Each trans should have one primary stress",
				Re:      exactOnePrimStressRe,
			},
			IllegalTransRe{
				Name:    "secondary_stress",
				Level:   "Format",
				Message: "Each trans can have max one secondary stress",
				Re:      maxOneSecStressRe,
			},
			SymbolSetRule{
				SymbolSet: symbolset,
			},
		}}
	return vali, nil
}
