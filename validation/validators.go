package validation

import (
	"regexp"

	"github.com/stts-se/pronlex/symbolset"
)

//"fmt"

//"github.com/stts-se/pronlex/dbapi"

func NewNSTDemoValidator() (Validator, error) {
	symbolset, err := symbolset.NewNSTSymbolSet()
	if err != nil {
		return err
	}
	var vali = Validator{[]Rule{
		MustHaveTrans{},
		NoEmptyTrans{},
		IllegalTransRe{
			Name:    "final_nostress_nolong",
			Level:   "Warning",
			Message: "final syllable should normally be unstressed with short vowel",
			Re:      regexp.MustCompile("\\$ (consonant )*(@|A|E|I|O|U|u0|Y|{|9|n=|l=|n`=|l`=)( consonant)*$"),
		},
		IllegalTransRe{
			Name:    "primary_stress",
			Level:   "Fatal",
			Message: "Primary stress required",
			Re:      regexp.MustCompile("\""),
		},
		SymbolSetRule{
			SymbolSet: symbolset,
		},
	}}
	return vali, nil
}
