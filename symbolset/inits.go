package symbolset

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Initialization functions for structs in package symbolset

func NewIPA() IPA {
	return IPA{
		ipa:      "ipa",
		accentI:  "\u02C8",
		accentII: "\u0300",
	}
}

// newSymbolSet: public constructor for SymbolSet with built-in error checks
func NewSymbolSet(name string, symbols []Symbol) (SymbolSet, error) {
	var nilRes SymbolSet

	// filtered lists
	phonemes := filterSymbolsByType(symbols, []SymbolType{Syllabic, NonSyllabic, Stress})
	phoneticSymbols := filterSymbolsByType(symbols, []SymbolType{Syllabic, NonSyllabic})
	stressSymbols := filterSymbolsByType(symbols, []SymbolType{Stress})
	syllabic := filterSymbolsByType(symbols, []SymbolType{Syllabic})
	nonSyllabic := filterSymbolsByType(symbols, []SymbolType{NonSyllabic})
	phonemeDelimiters := filterSymbolsByType(symbols, []SymbolType{PhonemeDelimiter})
	explicitPhonemeDelimiters := filterSymbolsByType(symbols, []SymbolType{ExplicitPhonemeDelimiter})

	// specific symbol initialization
	if len(phonemeDelimiters) < 1 {
		return nilRes, fmt.Errorf("No phoneme delimiters defined in symbol set %s", name)
	}
	phonemeDelimiter := phonemeDelimiters[0]

	var explicitPhonemeDelimiter Symbol
	if len(explicitPhonemeDelimiters) < 1 {
		explicitPhonemeDelimiter = Symbol{"", ExplicitPhonemeDelimiter, ""}
	} else {
		explicitPhonemeDelimiter = explicitPhonemeDelimiters[0]
	}

	// regexps
	phonemeRe, err := buildRegexp(phonemes)
	if err != nil {
		return nilRes, err
	}
	syllabicRe, err := buildRegexp(syllabic)
	if err != nil {
		return nilRes, err
	}
	nonSyllabicRe, err := buildRegexp(nonSyllabic)
	if err != nil {
		return nilRes, err
	}
	phonemeDelimiterRe, err := buildRegexp(phonemeDelimiters)
	if err != nil {
		return nilRes, err
	}

	// start: split re
	sort.Sort(SymbolSlice(symbols))
	acc := make([]string, 0)
	for _, s := range symbols {
		if len(s.String) > 0 {
			acc = append(acc, regexp.QuoteMeta(s.String))
		}
	}
	splitRe, err := regexp.Compile("(" + strings.Join(acc, "|") + ")")
	if err != nil {
		return nilRes, err
	}
	// end: split re

	res := SymbolSet{
		Name:                      name,
		Symbols:                   symbols,
		Phonemes:                  phonemes,
		PhoneticSymbols:           phoneticSymbols,
		StressSymbols:             stressSymbols,
		Syllabic:                  syllabic,
		NonSyllabic:               nonSyllabic,
		PhonemeDelimiters:         phonemeDelimiters,
		ExplicitPhonemeDelimiters: explicitPhonemeDelimiters,

		PhonemeDelimiter:            phonemeDelimiter,
		ExplicitPhonemeDelimiter:    explicitPhonemeDelimiter,
		HasExplicitPhonemeDelimiter: len(explicitPhonemeDelimiter.String) > 0,

		PhonemeRe:          phonemeRe,
		SyllabicRe:         syllabicRe,
		NonSyllabicRe:      nonSyllabicRe,
		PhonemeDelimiterRe: phonemeDelimiterRe,
		SplitRe:            splitRe,
	}
	return res, nil

}

func NewSymbolSetMapper(fromName string, toName string, symbolList []SymbolPair) (SymbolSetMapper, error) {
	var nilRes SymbolSetMapper

	ipa := NewIPA()

	toIsIPA := ipa.IsIPA(toName)
	fromIsIPA := ipa.IsIPA(fromName)

	fromSymbols := make([]Symbol, 0)
	toSymbols := make([]Symbol, 0)
	symbolMap := make(map[Symbol]Symbol)

	for _, pair := range symbolList {
		symbolMap[pair.Sym1] = pair.Sym2
		fromSymbols = append(fromSymbols, pair.Sym1)
		toSymbols = append(toSymbols, pair.Sym2)
	}

	from, err := NewSymbolSet(fromName, fromSymbols)
	if err != nil {
		return nilRes, err
	}
	to, err := NewSymbolSet(toName, toSymbols)
	if err != nil {
		return nilRes, err
	}
	if from.Name == to.Name {
		return nilRes, fmt.Errorf("Both phoneme sets cannot have the same name: %s", from.Name)
	}
	err = from.PreCheckAmbiguous()
	if err != nil {
		return nilRes, err
	}
	err = to.PreCheckAmbiguous()
	if err != nil {
		return nilRes, err
	}

	repeatedPhonemeDelimiters, err := regexp.Compile(to.PhonemeDelimiterRe.String() + "+")
	if err != nil {
		return nilRes, err
	}

	ssm := SymbolSetMapper{
		FromName:                  fromName,
		ToName:                    toName,
		SymbolList:                symbolList,
		FromIsIPA:                 fromIsIPA,
		ToIsIPA:                   toIsIPA,
		From:                      from,
		To:                        to,
		ipa:                       ipa,
		SymbolMap:                 symbolMap,
		RepeatedPhonemeDelimiters: repeatedPhonemeDelimiters,
	}
	return ssm, nil

}

func LoadSymbolSetMapper(fName string, fromName string, toName string) (SymbolSetMapper, error) {
	var nilRes SymbolSetMapper
	fh, err := os.Open(fName)
	defer fh.Close()
	if err != nil {
		return nilRes, err
	}

	s := bufio.NewScanner(fh)
	n := 0
	var descIndex, fromIndex, toIndex, typeIndex int
	maptable := make([]SymbolPair, 0)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nilRes, err
		}
		n++
		l := s.Text()
		fs := strings.Split(l, "\t")
		if n == 1 { // header
			descIndex = indexOf(fs, "DESC/EXAMPLE")
			fromIndex = indexOf(fs, fromName)
			toIndex = indexOf(fs, toName)
			typeIndex = indexOf(fs, "TYPE")

		} else {
			from := fs[fromIndex]
			to := fs[toIndex]
			desc := fs[descIndex]
			typeS := fs[typeIndex]
			var symType SymbolType
			switch typeS {
			case "syllabic":
				symType = Syllabic
			case "non syllabic":
				symType = NonSyllabic
			case "stress":
				symType = Stress
			case "phoneme delimiter":
				symType = PhonemeDelimiter
			case "explicit phoneme delimiter":
				symType = ExplicitPhonemeDelimiter
			case "syllable delimiter":
				symType = SyllableDelimiter
			case "morpheme delimiter":
				symType = MorphemeDelimiter
			case "word delimiter":
				symType = WordDelimiter
			default:
				return nilRes, fmt.Errorf("Unknown symbol type on line:\t" + l)
			}
			symFrom := Symbol{String: from, Type: symType, Desc: desc}
			symTo := Symbol{String: to, Type: symType, Desc: desc}
			maptable = append(maptable, SymbolPair{symFrom, symTo})
		}
	}
	ssm, err := NewSymbolSetMapper(fromName, toName, maptable)
	if err != nil {
		return ssm, err
	} else {
		return ssm, nil
	}
}

// end: initialization
