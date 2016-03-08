package symbolset

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"
)

// start: general util stuff
func buildRegexp(symbols []Symbol) (*regexp.Regexp, error) {
	re := ""
	for _, s := range symbols {
		re = re + regexp.QuoteMeta(s.String)
	}
	return regexp.Compile(re)
}
func contains(symbols []Symbol, symbol string) bool {
	for _, s := range symbols {
		if s.String == symbol {
			return true
		}
	}
	return false
}

// end: general util stuff

// SymbolSet: package private struct
type SymbolSet struct {
	name    string
	symbols []Symbol
	//phonemeRe *regexp.Regexp
}

// NewSymbolSet: public constructor for SymbolSet with build-in error checks
func NewSymbolSet(name string, symbols []Symbol) (SymbolSet, error) {
	res := SymbolSet{name, symbols}
	if len(res.PhonemeDelimiters()) < 1 {
		return res, fmt.Errorf("No phoneme delimiters defined in symbol set")
	} else {
		return res, nil
	}
	panic("Move Scala-style lazy fields to contructor!")
}

// start: symbol set
// Sort according to symbol length
type SymbolSlice []Symbol

func (a SymbolSlice) Len() int           { return len(a) }
func (a SymbolSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SymbolSlice) Less(i, j int) bool { return len(a[i].String) > len(a[j].String) }

type SymbolType int

const (
	Syllabic SymbolType = iota
	NonSyllabic
	Stress

	PhonemeDelimiter
	ExplicitPhonemeDelimiter
	SyllableDelimiter
	MorphemeDelimiter
	WordDelimiter
)

type Symbol struct {
	Desc   string
	String string
	Type   SymbolType
}

func (ss SymbolSet) Symbols() []Symbol {
	return ss.symbols
}

func (ss SymbolSet) Name() string {
	return ss.name
}

func (ss SymbolSet) Contains(symbol string) bool {
	for _, s := range ss.Symbols() {
		if s.String == symbol {
			return true
		}
	}
	return false
}

func (ss SymbolSet) StressSymbols() []Symbol {
	res := make([]Symbol, 0)
	for _, s := range ss.Symbols() {
		if s.Type == Stress {
			res = append(res, s)
		}
	}
	return res
}

func (ss SymbolSet) NonSyllabic() []Symbol {
	res := make([]Symbol, 0)
	for _, s := range ss.Symbols() {
		if s.Type == NonSyllabic {
			res = append(res, s)
		}
	}
	return res
}

func (ss SymbolSet) Syllabic() []Symbol {
	res := make([]Symbol, 0)
	for _, s := range ss.Symbols() {
		if s.Type == Syllabic {
			res = append(res, s)
		}
	}
	return res
}

func (ss SymbolSet) Phonemes() []Symbol {
	res := make([]Symbol, 0)
	for _, s := range ss.Symbols() {
		if s.Type == Syllabic || s.Type == NonSyllabic || s.Type == Stress {
			res = append(res, s)
		}
	}
	return res
}

func (ss SymbolSet) PhoneticSymbols() []Symbol {
	res := make([]Symbol, 0)
	for _, s := range ss.Symbols() {
		if s.Type == Syllabic || s.Type == NonSyllabic {
			res = append(res, s)
		}
	}
	return res
}

func (ss SymbolSet) PhonemeDelimiters() []Symbol {
	res := make([]Symbol, 0)
	for _, s := range ss.Symbols() {
		if s.Type == PhonemeDelimiter {
			res = append(res, s)
		}
	}
	return res
}
func (ss SymbolSet) PhonemeDelimiter() Symbol {
	return ss.PhonemeDelimiters()[0]
}

func (ss SymbolSet) ExplicitPhonemeDelimiter() (Symbol, error) {
	exp := ss.ExplicitPhonemeDelimiters()
	switch len(exp) {
	case 0:
		return Symbol{}, fmt.Errorf("No explicit phoneme delimiter in symbol set")
	default:
		return exp[0], nil
	}
}

func (ss SymbolSet) ExplicitPhonemeDelimiters() []Symbol {
	res := make([]Symbol, 0)
	for _, s := range ss.Symbols() {
		if s.Type == ExplicitPhonemeDelimiter {
			res = append(res, s)
		}
	}
	return res
}

func (ss SymbolSet) PhonemeRe() *regexp.Regexp {
	re, err := buildRegexp(ss.Phonemes())
	if err != nil {
		log.Fatal(err) // TODO
	}
	return re
}

func (ss SymbolSet) NonSyllabicRe() *regexp.Regexp {
	re, err := buildRegexp(ss.NonSyllabic())
	if err != nil {
		log.Fatal(err) // TODO
	}
	return re
}

func (ss SymbolSet) SyllabicRe() *regexp.Regexp {
	re, err := buildRegexp(ss.Syllabic())
	if err != nil {
		log.Fatal(err) // TODO
	}
	return re
}

func (ss SymbolSet) PhonemeDelimiterRe() *regexp.Regexp {
	re, err := buildRegexp(ss.PhonemeDelimiters())
	if err != nil {
		log.Fatal(err) // TODO
	}
	return re
}

func (ss SymbolSet) Get(symbol string) (Symbol, error) {
	for _, s := range ss.symbols {
		if s.String == symbol {
			return s, nil
		}
	}
	return Symbol{}, fmt.Errorf("No symbol /%s/ in symbol set", symbol)
}

func (ss SymbolSet) SplitRe() *regexp.Regexp {
	symbols := ss.Symbols()
	sort.Sort(SymbolSlice(symbols))
	acc := make([]string, 0)
	for _, s := range symbols {
		if len(s.String) > 0 {
			acc = append(acc, regexp.QuoteMeta(s.String))
		}
	}
	re, err := regexp.Compile("(" + strings.Join(acc, "|") + ")")
	if err != nil {
		log.Fatal(err) // TODO
	}
	return re
}

func (ss SymbolSet) FilterAmbiguous(trans []string) ([]string, error) {
	potentiallyAmbs := ss.PhoneticSymbols()
	phnDel := ss.PhonemeDelimiter().String
	explicitPhnDel, explicitPhnDelErr := ss.ExplicitPhonemeDelimiter()
	res := make([]string, 0)
	for i, phn0 := range trans[0 : len(trans)-1] {
		phn1 := trans[i+1]
		test := phn0 + phnDel + phn1
		if contains(potentiallyAmbs, test) {
			if explicitPhnDelErr != nil {
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
	allSymbols := ss.Phonemes()
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
	delim := ss.PhonemeDelimiterRe()
	if delim.FindStringIndex("") != nil {
		rest := input
		acc := make([]string, 0)
		for len(rest) > 0 {
			match := ss.SplitRe().FindStringIndex(rest)
			switch match {
			case nil:
				return nil, fmt.Errorf("Transcription not splittable (invalid symbols?)! input=/%s/, acc=/%s/, rest=/%s/", input, strings.Join(acc, ss.PhonemeDelimiter().String), rest)

			default:
				acc = append(acc, rest[match[0]:match[1]])
				rest = rest[match[1]+1:]
			}
		}
		return nil, nil
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

func NewIPA() IPA {
	return IPA{
		ipa:      "ipa",
		accentI:  "\u02C8",
		accentII: "\u0300",
	}
}
func (ipa IPA) IsIPA(symbolSetName string) bool {
	return strings.Contains(strings.ToLower(symbolSetName), ipa.ipa)
}

func (ipa IPA) checkFilterRequirements(ss SymbolSet) error {
	if !ss.Contains(ipa.accentI) {
		return fmt.Errorf("No IPA stress symbol in stress symbol list? IPA stress =/%v/, stress symbols=%v", ipa.accentI, ss.StressSymbols())
	} else if !ss.Contains(ipa.accentI + ipa.accentII) {
		return fmt.Errorf("No IPA tone II symbol in stress symbol list? IPA stress =/%s/, stress symbols=%v", ipa.accentI, ss.StressSymbols())
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
	repl, err := regexp.Compile(ipa.accentI + "(" + ss.PhonemeRe().String() + "+)" + ipa.accentII)
	if err != nil {
		return "", fmt.Errorf("Couldn't parse regexp : %v", err)
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
		repl, err := regexp.Compile(ipa.accentI + ipa.accentII + "(" + ss.NonSyllabicRe().String() + "*)(" + ss.SyllabicRe().String() + ")")
		res := repl.ReplaceAllString(trans, ipa.accentI+ipa.accentII+"$1")
		return res, nil
	} else {
		return trans, nil
	}
}

// end: IPA
