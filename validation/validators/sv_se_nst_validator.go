package validators

import (
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/pronlex/validation/rules"
	"github.com/stts-se/symbolset"
)

func newSvSeNstValidator(symbolset symbolset.SymbolSet) (validation.Validator, error) {
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

	maxOneSyllabic, err := rules.ProcessTransRe(symbolset, "syllabic[^.+%\"-]*( +syllabic)")
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

	repeatedPhnRe, err := rules.ProcessTransRe(symbolset, "symbol( +[.~])? +\\1( |$)")
	if err != nil {
		return validation.Validator{}, err
	}

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
				Accept: []lex.Entry{
					{Transcriptions: []lex.Transcription{
						{Strn: "\" A: . p a"}}},
					{Transcriptions: []lex.Transcription{
						{Strn: "p O . \" E N"}}},
				},
				Reject: []lex.Entry{
					{Transcriptions: []lex.Transcription{
						{Strn: "A: \" . p a"}}},
					{Transcriptions: []lex.Transcription{
						{Strn: "s k r \" A: . p a"}}},
				},
			},
			rules.RequiredTransRe{
				NameStr:  "syllabic",
				LevelStr: "Format",
				Message:  "Each syllable needs a syllabic phoneme",
				Re:       syllabicRe,
			},
			rules.IllegalTransRe{
				NameStr:  "MaxOneSyllabic",
				LevelStr: "Fatal",
				Message:  "A syllable cannot contain more than one syllabic phoneme",
				Re:       maxOneSyllabic,
			},
			rules.IllegalTransRe{
				NameStr:  "repeated_phonemes",
				LevelStr: "Fatal",
				Message:  "Repeated phonemes cannot be used within the same morpheme",
				Re:       repeatedPhnRe,
			},
			decomp2Orth,
			rules.SymbolSetRule{
				SymbolSet: symbolset,
			},
		}}
	return vali, nil
}
