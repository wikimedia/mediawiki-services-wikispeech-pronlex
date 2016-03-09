package symbolset

import (
	"fmt"
	"regexp"
	"strings"
)

// structs in package symbolset

// SymbolSet: package private struct
type SymbolSet struct {
	Name    string
	Symbols []Symbol

	// derived values computed upon initialization
	Phonemes                  []Symbol
	PhoneticSymbols           []Symbol
	StressSymbols             []Symbol
	Syllabic                  []Symbol
	NonSyllabic               []Symbol
	PhonemeDelimiters         []Symbol
	ExplicitPhonemeDelimiters []Symbol

	PhonemeDelimiter            Symbol
	ExplicitPhonemeDelimiter    Symbol
	HasExplicitPhonemeDelimiter bool

	PhonemeRe          *regexp.Regexp
	SyllabicRe         *regexp.Regexp
	NonSyllabicRe      *regexp.Regexp
	PhonemeDelimiterRe *regexp.Regexp
	SplitRe            *regexp.Regexp
}

// start: SymbolSetMapper
type SymbolSetMapper struct {
	FromName   string
	ToName     string
	SymbolList []SymbolPair

	ipa IPA

	// derived values computed upon initialization
	FromIsIPA bool
	ToIsIPA   bool

	From      SymbolSet
	To        SymbolSet
	SymbolMap map[Symbol]Symbol

	RepeatedPhonemeDelimiters *regexp.Regexp
}

func (ssm SymbolSetMapper) preFilter(trans string, ss SymbolSet) (string, error) {
	switch ssm.FromIsIPA {
	case true:
		return ssm.ipa.filterBeforeMappingFromIpa(trans, ss)
	default:
		return trans, nil
	}
}

func (ssm SymbolSetMapper) postFilter(trans string, ss SymbolSet) (string, error) {
	switch ssm.ToIsIPA {
	case true:
		return ssm.ipa.filterAfterMappingToIpa(trans, ss)
	default:
		return trans, nil
	}
}

func (ssm SymbolSetMapper) mapTranscription(input string) (string, error) {
	res, err := ssm.preFilter(input, ssm.From)
	if err != nil {
		return res, err
	}
	splitted, err := ssm.From.SplitTranscription(res)
	if err != nil {
		return res, err
	}
	mapped := make([]string, 0)
	for _, fromS := range splitted {
		from, err := ssm.From.Get(fromS)
		if err != nil {
			return res, fmt.Errorf("input symbol /%s/ is undefined : %v", fromS, err)
		}
		to := ssm.SymbolMap[from]
		if to.Type == UndefinedSymbol {
			return res, fmt.Errorf("couldn't map symbol /%s/", fromS)
		}
		mapped = append(mapped, to.String)
	}
	mapped, err = ssm.To.FilterAmbiguous(mapped)
	if err != nil {
		return res, err
	}
	res = strings.Join(mapped, ssm.To.PhonemeDelimiter.String)

	// remove repeated phoneme delimiters
	res = ssm.RepeatedPhonemeDelimiters.ReplaceAllString(res, ssm.To.PhonemeDelimiter.String)
	return ssm.postFilter(res, ssm.To)
}

// end: SymbolSetMapper

type SymbolPair struct {
	Sym1 Symbol
	Sym2 Symbol
}

// Sort according to symbol length
type SymbolSlice []Symbol

func (a SymbolSlice) Len() int           { return len(a) }
func (a SymbolSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SymbolSlice) Less(i, j int) bool { return len(a[i].String) > len(a[j].String) }

type SymbolType int

const (
	UndefinedSymbol SymbolType = iota

	Syllabic
	NonSyllabic
	Stress

	PhonemeDelimiter
	ExplicitPhonemeDelimiter
	SyllableDelimiter
	MorphemeDelimiter
	WordDelimiter
)

type Symbol struct {
	String string
	Type   SymbolType
	Desc   string
}

func (ss SymbolSet) Contains(symbol string) bool {
	for _, s := range ss.Symbols {
		if s.String == symbol {
			return true
		}
	}
	return false
}

func (ss SymbolSet) Get(symbol string) (Symbol, error) {
	for _, s := range ss.Symbols {
		if s.String == symbol {
			return s, nil
		}
	}
	return Symbol{}, fmt.Errorf("No symbol /%s/ in symbol set", symbol)
}

func (ss SymbolSet) FilterAmbiguous(trans []string) ([]string, error) {
	potentiallyAmbs := ss.PhoneticSymbols
	phnDel := ss.PhonemeDelimiter.String
	explicitPhnDel := ss.ExplicitPhonemeDelimiter
	res := make([]string, 0)
	for i, phn0 := range trans[0 : len(trans)-1] {
		phn1 := trans[i+1]
		test := phn0 + phnDel + phn1
		if contains(potentiallyAmbs, test) {
			if !ss.HasExplicitPhonemeDelimiter {
				return res, fmt.Errorf("Explicit phoneme delimiter was undefined when needed for [%s] vs [%s] + [%s]", (phn0 + phn1), phn0, phn1)
			} else {
				res = append(res, phn0+explicitPhnDel.String)
			}
		}
	}
	res = append(res, trans[len(trans)])
	return trans, nil
}

func (ss SymbolSet) PreCheckAmbiguous() error {
	allSymbols := ss.Phonemes
	res := make([]string, 0)
	for _, a := range allSymbols {
		for _, b := range allSymbols {
			if len(a.String) > 0 && len(b.String) > 0 {
				res = append(res, a.String)
				res = append(res, b.String)
			}
		}
	}
	_, err := ss.FilterAmbiguous(res)
	return err
}

func (ss SymbolSet) SplitTranscription(input string) ([]string, error) {
	delim := ss.PhonemeDelimiterRe
	if delim.FindStringIndex("") != nil {
		rest := input
		acc := make([]string, 0)
		for len(rest) > 0 {
			match := ss.SplitRe.FindStringIndex(rest)
			switch match {
			case nil:
				return nil, fmt.Errorf("Transcription not splittable (invalid symbols?)! input=/%s/, acc=/%s/, rest=/%s/", input, strings.Join(acc, ss.PhonemeDelimiter.String), rest)

			default:
				acc = append(acc, rest[match[0]:match[1]])
				rest = rest[match[1]:]
			}
		}
		return acc, nil
	} else {
		return delim.Split(input, -1), nil
	}
}

// end: symbol set

// start: IPA
// SYMBOLS: http://www.phon.ucl.ac.uk/home/wells/ipa-unicode.htm#numbers
type IPA struct {
	ipa      string
	accentI  string
	accentII string
}

func (ipa IPA) IsIPA(symbolSetName string) bool {
	return strings.Contains(strings.ToLower(symbolSetName), ipa.ipa)
}

func (ipa IPA) checkFilterRequirements(ss SymbolSet) error {
	if !ss.Contains(ipa.accentI) {
		return fmt.Errorf("No IPA stress symbol in stress symbol list? IPA stress =/%v/, stress symbols=%v", ipa.accentI, ss.StressSymbols)
	} else if !ss.Contains(ipa.accentI + ipa.accentII) {
		return fmt.Errorf("No IPA tone II symbol in stress symbol list? IPA stress =/%s/, stress symbols=%v", ipa.accentI, ss.StressSymbols)
	} else {
		return nil
	}
}

func (ipa IPA) filterBeforeMappingFromIpa(trans string, ss SymbolSet) (string, error) {
	// IPA: ˈba`ŋ.ka => ˈ`baŋ.ka"
	err := ipa.checkFilterRequirements(ss)
	if err != nil {
		return "", err
	}
	s := ipa.accentI + "(" + ss.PhonemeRe.String() + "+)" + ipa.accentII
	repl, err := regexp.Compile(s)
	if err != nil {
		err = fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	res := repl.ReplaceAllString(trans, ipa.accentI+ipa.accentII+"$1")
	return res, nil
}

func (ipa IPA) filterAfterMappingToIpa(trans string, ss SymbolSet) (string, error) {
	conditionForAfterMapping := ipa.accentI + ipa.accentII
	// IPA: /'`pa.pa/ => /'pa`.pa/
	if strings.Contains(trans, conditionForAfterMapping) {
		err := ipa.checkFilterRequirements(ss)
		if err != nil {
			return "", err
		}
		repl, err := regexp.Compile(ipa.accentI + ipa.accentII + "(" + ss.NonSyllabicRe.String() + "*)(" + ss.SyllabicRe.String() + ")")
		res := repl.ReplaceAllString(trans, ipa.accentI+ipa.accentII+"$1")
		return res, nil
	} else {
		return trans, nil
	}
}

// end: IPA
