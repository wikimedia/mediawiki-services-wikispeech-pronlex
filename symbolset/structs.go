package symbolset

import (
	"fmt"
	"regexp"
	"strings"
)

// structs in package symbolset

// SymbolSet struct for package private usage.
// To create a new SymbolSet, use NewSymbolSet
type SymbolSet struct {
	Name    string
	Symbols []Symbol

	// derived values computed upon initialization
	phonemes                  []Symbol
	phoneticSymbols           []Symbol
	stressSymbols             []Symbol
	syllabic                  []Symbol
	nonSyllabic               []Symbol
	phonemeDelimiters         []Symbol
	explicitPhonemeDelimiters []Symbol

	phonemeDelimiter            Symbol
	explicitPhonemeDelimiter    Symbol
	hasExplicitPhonemeDelimiter bool

	phonemeRe          *regexp.Regexp
	syllabicRe         *regexp.Regexp
	nonSyllabicRe      *regexp.Regexp
	phonemeDelimiterRe *regexp.Regexp
	symbolRe           *regexp.Regexp
}

// SymbolSetMapper struct for package private usage.
// To create a new SymbolSet, use NewSymbolSet.
type SymbolSetMapper struct {
	FromName   string
	ToName     string
	SymbolList []SymbolPair

	ipa IPA

	// derived values computed upon initialization
	fromIsIPA bool
	toIsIPA   bool

	from      SymbolSet
	to        SymbolSet
	symbolMap map[Symbol]Symbol

	repeatedPhonemeDelimiters *regexp.Regexp
}

func (ssm SymbolSetMapper) preFilter(trans string, ss SymbolSet) (string, error) {
	switch ssm.fromIsIPA {
	case true:
		return ssm.ipa.filterBeforeMappingFromIpa(trans, ss)
	default:
		return trans, nil
	}
}

func (ssm SymbolSetMapper) postFilter(trans string, ss SymbolSet) (string, error) {
	switch ssm.toIsIPA {
	case true:
		return ssm.ipa.filterAfterMappingToIpa(trans, ss)
	default:
		return trans, nil
	}
}

func (ssm SymbolSetMapper) mapTranscription(input string) (string, error) {
	res, err := ssm.preFilter(input, ssm.from)
	if err != nil {
		return "", err
	}
	splitted, err := ssm.from.SplitTranscription(res)
	if err != nil {
		return "", err
	}
	mapped := make([]string, 0)
	for _, fromS := range splitted {
		from, err := ssm.from.Get(fromS)
		if err != nil {
			return "", fmt.Errorf("input symbol /%s/ is undefined : %v", fromS, err)
		}
		to := ssm.symbolMap[from]
		if to.Type == UndefinedSymbol {
			return "", fmt.Errorf("couldn't map symbol /%s/", fromS)
		}
		if len(to.String) > 0 {
			mapped = append(mapped, to.String)
		}
	}
	mapped, err = ssm.to.FilterAmbiguous(mapped)
	if err != nil {
		return "", err
	}
	res = strings.Join(mapped, ssm.to.phonemeDelimiter.String)

	// remove repeated phoneme delimiters
	res = ssm.repeatedPhonemeDelimiters.ReplaceAllString(res, ssm.to.phonemeDelimiter.String)
	return ssm.postFilter(res, ssm.to)
}

type SymbolPair struct {
	Sym1 Symbol
	Sym2 Symbol
}

// SymbolSlice for sorting symbol slices according to symbol length
type SymbolSlice []Symbol

func (a SymbolSlice) Len() int      { return len(a) }
func (a SymbolSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SymbolSlice) Less(i, j int) bool {
	if len(a[i].String) != len(a[j].String) {
		return len(a[i].String) > len(a[j].String)
	} else {
		return a[i].String < a[j].String
	}
}

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
	return Symbol{}, fmt.Errorf("no symbol /%s/ in symbol set", symbol)
}

func (ss SymbolSet) FilterAmbiguous(trans []string) ([]string, error) {
	potentiallyAmbs := ss.phoneticSymbols
	phnDel := ss.phonemeDelimiter.String
	explicitPhnDel := ss.explicitPhonemeDelimiter
	res := make([]string, 0)
	for i, phn0 := range trans[0 : len(trans)-1] {
		phn1 := trans[i+1]
		test := phn0 + phnDel + phn1
		if contains(potentiallyAmbs, test) {
			if !ss.hasExplicitPhonemeDelimiter {
				return nil, fmt.Errorf("explicit phoneme delimiter in %s was undefined when needed for [%s] vs [%s] + [%s]", ss.Name, (phn0 + phn1), phn0, phn1)
			} else {
				res = append(res, phn0+explicitPhnDel.String)
			}
		} else {
			res = append(res, phn0)
		}
	}
	res = append(res, trans[len(trans)-1])
	return res, nil
}

func (ss SymbolSet) PreCheckAmbiguous() error {
	allSymbols := ss.phonemes
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
	delim := ss.phonemeDelimiterRe
	if delim.FindStringIndex("") != nil {
		rest := input
		acc := make([]string, 0)
		for len(rest) > 0 {
			match := ss.symbolRe.FindStringIndex(rest)
			switch match {
			case nil:
				return nil, fmt.Errorf("transcription not splittable (invalid symbols?)! input=/%s/, acc=/%s/, rest=/%s/", input, strings.Join(acc, ss.phonemeDelimiter.String), rest)

			default:
				if match[0] != 0 {
					return nil, fmt.Errorf("couldn't parse transcription /%s/, it may contain undefined symbols!", input)
				}
				acc = append(acc, rest[match[0]:match[1]])
				rest = rest[match[1]:]
			}
		}
		return acc, nil
	} else {
		return delim.Split(input, -1), nil
	}
}

// IPA utilility functions with struct for package private usage.
// To create a new SymbolSet, use NewSymbolSet.
// Symbols and codes: http://www.phon.ucl.ac.uk/home/wells/ipa-unicode.htm#numbers
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
		return fmt.Errorf("no IPA stress symbol in stress symbol list? IPA stress =/%v/, stress symbols=%v", ipa.accentI, ss.stressSymbols)
	} else if !ss.Contains(ipa.accentI + ipa.accentII) {
		return fmt.Errorf("no IPA tone II symbol in stress symbol list? IPA stress =/%s/, stress symbols=%v", ipa.accentI, ss.stressSymbols)
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
	s := ipa.accentI + "(" + ss.phonemeRe.String() + "+)" + ipa.accentII
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
		repl, err := regexp.Compile(ipa.accentI + ipa.accentII + "(" + ss.nonSyllabicRe.String() + "*)(" + ss.syllabicRe.String() + ")")
		res := repl.ReplaceAllString(trans, ipa.accentI+"$1$2"+ipa.accentII)
		return res, nil
	} else {
		return trans, nil
	}
}
