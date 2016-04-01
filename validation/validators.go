package validation

import "github.com/stts-se/pronlex/symbolset"

func NewNSTDemoValidator() (Validator, error) {
	symbolset, err := symbolset.NewNSTSymbolSet()
	if err != nil {
		return Validator{}, err
	}
	var vali = Validator{[]Rule{
		MustHaveTrans{},
		NoEmptyTrans{},
		IllegalTransRe{
			Name:    "final_nostress_nolong",
			Level:   "Warning",
			Message: "final syllable should normally be unstressed with short vowel",
			Re:      ProcessTransRe(symbolset, "\\$ (syllabic )*(@|A|E|I|O|U|u0|Y|{|9|n=|l=|n`=|l`=)( syllabic)*$"),
		},
		RequiredTransRe{
			Name:    "primary_stress",
			Level:   "Fatal",
			Message: "Primary stress required",
			Re:      ProcessTransRe(symbolset, "\""),
		},
		SymbolSetRule{
			SymbolSet: symbolset,
		},
	}}
	return vali, nil
}
