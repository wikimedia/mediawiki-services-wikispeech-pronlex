package validation

import "github.com/stts-se/pronlex/symbolset"

// NewNSTDemoValidator is used for testing
func NewNSTDemoValidator() (Validator, error) {
	symbolset, err := symbolset.NewNSTSymbolSet()
	if err != nil {
		return Validator{}, err
	}
	finalNostressNolongRe, err := ProcessTransRe(symbolset, "\\$ (nonsyllabic )*(@|A|E|I|O|U|u0|Y|{|9|n=|l=|n`=|l`=)( nonsyllabic)*$")
	if err != nil {
		return Validator{}, err
	}
	primaryStressRe, err := ProcessTransRe(symbolset, "\"")
	if err != nil {
		return Validator{}, err
	}
	syllabicRe, err := ProcessTransRe(symbolset, "^(\"\"|\"|%)? *(nonsyllabic )*syllabic( nonsyllabic)*( (\\$|-) (\"\"|\"|%)? *(nonsyllabic )*syllabic( nonsyllabic)*)*$")
	if err != nil {
		return Validator{}, err
	}

	var vali = Validator{[]Rule{
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
		SymbolSetRule{
			SymbolSet: symbolset,
		},
	}}
	return vali, nil
}
