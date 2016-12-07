package symbolset2

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// inits.go Initialization functions for structs in package symbolset

func trimIfNeeded(s string) string {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) > 0 {
		return trimmed
	}
	return s
}

// NewSymbolSet is a constructor for 'symbols' with built-in error checks
func NewSymbolSet(name string, symbols []Symbol) (SymbolSet, error) {
	return NewSymbolSetWithTests(name, symbols, true)
}

// NewsymbolsWithTests is a constructor for 'symbols' with built-in error checks
func NewSymbolSetWithTests(name string, symbols []Symbol, checkForDups bool) (SymbolSet, error) {
	var nilRes SymbolSet

	// filtered lists
	phonemes := FilterSymbolsByCat(symbols, []SymbolCat{Syllabic, NonSyllabic, Stress})
	phoneticSymbols := FilterSymbolsByCat(symbols, []SymbolCat{Syllabic, NonSyllabic})
	stressSymbols := FilterSymbolsByCat(symbols, []SymbolCat{Stress})
	syllabic := FilterSymbolsByCat(symbols, []SymbolCat{Syllabic})
	nonSyllabic := FilterSymbolsByCat(symbols, []SymbolCat{NonSyllabic})
	phonemeDelimiters := FilterSymbolsByCat(symbols, []SymbolCat{PhonemeDelimiter})

	// specific symbol initialization
	if len(phonemeDelimiters) < 1 {
		return nilRes, fmt.Errorf("no phoneme delimiters defined in symbol set %s", name)
	}
	phonemeDelimiter := phonemeDelimiters[0]

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
	symbolRe, err := buildRegexpWithGroup(symbols, true, false)

	if err != nil {
		return nilRes, err
	}

	if checkForDups {
		seenSymbols := make(map[string]Symbol)
		var dupSymbols []string
		for _, symbol := range symbols {
			if _, exists := seenSymbols[symbol.String]; exists {
				dupSymbols = append(dupSymbols, symbol.String)
			}
			seenSymbols[symbol.String] = symbol
		}
		if len(dupSymbols) > 0 {
			return nilRes, fmt.Errorf("input symbol set contains duplicates of phoneme %v. All symbols: %v", dupSymbols, symbols)
		}
	}

	res := SymbolSet{
		Name:    name,
		Symbols: symbols,

		isInit: true,

		phonemes:        phonemes,
		phoneticSymbols: phoneticSymbols,
		stressSymbols:   stressSymbols,
		syllabic:        syllabic,
		nonSyllabic:     nonSyllabic,

		phonemeDelimiter: phonemeDelimiter,

		PhonemeRe:          phonemeRe,
		SyllabicRe:         syllabicRe,
		NonSyllabicRe:      nonSyllabicRe,
		SymbolRe:           symbolRe,
		phonemeDelimiterRe: phonemeDelimiterRe,
	}
	return res, nil

}

// LoadSymbolSet loads a SymbolSet from file
func LoadSymbolSet(fName string) (SymbolSet, error) {
	name := filepath.Base(fName)
	var extension = filepath.Ext(name)
	name = name[0 : len(name)-len(extension)]
	return loadSymbolSet_(name, fName)
}

var header = "DESCRIPTION	SYMBOL	IPA	IPA UNICODE	CATEGORY"

// loadSymbolSet_ loads a SymbolSet from file
func loadSymbolSet_(name string, fName string) (SymbolSet, error) {
	var nilRes SymbolSet
	fh, err := os.Open(fName)
	defer fh.Close()
	if err != nil {
		return nilRes, err
	}
	s := bufio.NewScanner(fh)
	n := 0
	var descIndex = 0
	var symbolIndex = 1
	var ipaIndex = 2
	var ipaUnicodeIndex = 3
	var typeIndex = 4
	var symbols = make([]Symbol, 0)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nilRes, err
		}
		n++
		l := s.Text()
		if len(strings.TrimSpace(l)) > 0 && !strings.HasPrefix(strings.TrimSpace(l), "#") {
			if n == 1 { // header
				if l != header {
					return nilRes, fmt.Errorf("expected header '%s', found '%s'", header, l)
				}
			} else {
				fs := strings.Split(l, "\t")
				symbol := trimIfNeeded(fs[symbolIndex])
				ipa := trimIfNeeded(fs[ipaIndex])
				ipaUnicode := trimIfNeeded(fs[ipaUnicodeIndex])
				desc := fs[descIndex]
				typeS := fs[typeIndex]
				var symCat SymbolCat
				switch typeS {
				case "syllabic":
					symCat = Syllabic
				case "non syllabic":
					symCat = NonSyllabic
				case "stress":
					symCat = Stress
				case "phoneme delimiter":
					symCat = PhonemeDelimiter
				case "syllable delimiter":
					symCat = SyllableDelimiter
				case "morpheme delimiter":
					symCat = MorphemeDelimiter
				case "compound delimiter":
					symCat = CompoundDelimiter
				case "word delimiter":
					symCat = WordDelimiter
				default:
					return nilRes, fmt.Errorf("unknown symbol type on line:\t" + l)
				}
				sym := Symbol{String: symbol, Cat: symCat, Desc: desc,
					IPA: IPA{String: ipa, Unicode: ipaUnicode},
				}
				symbols = append(symbols, sym)
			}
		}
	}

	m, err := NewSymbolSet(name, symbols)
	if err != nil {
		return nilRes, fmt.Errorf("couldn't load mapper from file %v : %v", fName, err)
	}
	return m, nil
}
