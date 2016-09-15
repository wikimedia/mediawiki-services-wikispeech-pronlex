package vrules

import (
	"fmt"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
)

type ValidatorService struct {
	Validators map[string]*validation.Validator
}

func (vs ValidatorService) ValidatorForName(symbolSetName string) (*validation.Validator, error) {
	if vv, ok := vs.Validators[symbolSetName]; ok {
		return vv, nil
	}
	return &validation.Validator{}, fmt.Errorf("no validator loaded for symbolset %s", symbolSetName)

}

func (vs ValidatorService) Load(symbolsets map[string]symbolset.SymbolSet) error {
	if ss, ok := symbolsets["sv-se_ws-sampa"]; ok {
		v, err := newSvSeNstValidator(ss.From)
		if err != nil {
			return fmt.Errorf("couldn't initialize symbol set : %v", err)
		}
		vs.Validators[ss.Name] = &v
	}
	// if ss, ok := symbolsets["nb-no_ws-sampa"]; ok {
	// v, err := NewNbNoNstValidator(ss.From)
	// if err != nil {
	// 	return fmt.Errorf("couldn't initialize symbol set : %v", err)
	// }
	// vs.validators[ss.Name] = v
	// }
	// if ss, ok := symbolsets["en-us_sampa_mary"]; ok {
	// v, err := NewEnUsCmuNstValidator(ss.From)
	// if err != nil {
	// 	return fmt.Errorf("couldn't initialize symbol set : %v", err)
	// }
	// vs.validators[ss.Name] = v
	// }
	return nil
}

func newSvSeNstValidator(symbolset symbolset.Symbols) (validation.Validator, error) {
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
