package symbolset2

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// inits.go Initialization functions for structs in package symbolset

// NewSymbolSet is a constructor for 'symbols' with built-in error checks
func NewSymbolSet(name string, symbols []Symbol) (SymbolSet, error) {
	return NewSymbolSetWithTests(name, symbols, true)
}

// NewSymbolSetWithTests is a constructor for 'symbols' with built-in error checks
func NewSymbolSetWithTests(name string, symbols []Symbol, checkForDups bool) (SymbolSet, error) {
	var nilRes SymbolSet

	// filtered lists
	phonemes := filterSymbolsByCat(symbols, []SymbolCat{Syllabic, NonSyllabic, Stress})
	phoneticSymbols := filterSymbolsByCat(symbols, []SymbolCat{Syllabic, NonSyllabic})
	stressSymbols := filterSymbolsByCat(symbols, []SymbolCat{Stress})
	syllabic := filterSymbolsByCat(symbols, []SymbolCat{Syllabic})
	nonSyllabic := filterSymbolsByCat(symbols, []SymbolCat{NonSyllabic})
	phonemeDelimiters := filterSymbolsByCat(symbols, []SymbolCat{PhonemeDelimiter})

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

	// IPA regexps
	ipaSyllabicRe, err := buildIPARegexp(syllabic)
	if err != nil {
		return nilRes, err
	}
	ipaNonSyllabicRe, err := buildIPARegexp(nonSyllabic)
	if err != nil {
		return nilRes, err
	}
	ipaPhonemeRe, err := buildIPARegexp(phonemes)
	if err != nil {
		return nilRes, err
	}

	// compare ipa string vs unicode
	for _, symbol := range symbols {
		uFromString := string2unicode(symbol.IPA.String)

		if len(symbol.IPA.String) == 0 {
			uFromString = ""
		}
		if symbol.IPA.Unicode != uFromString {
			return nilRes, fmt.Errorf("ipa symbol /%s/ does not match unicode '%s' -- expected '%s'", symbol.IPA.String, symbol.IPA.Unicode, uFromString)
		}
		if strings.Contains(symbol.IPA.String, " ") {
			return nilRes, fmt.Errorf("ipa symbols cannot contain white space -- found /%s/", symbol.IPA.String)
		}
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

	repeatedPhonemeDelimiters, err := regexp.Compile(phonemeDelimiterRe.String() + "+")
	if err != nil {
		return nilRes, err
	}

	ssType := Other
	nameLC := strings.ToLower(name)
	if strings.Contains(nameLC, "ipa") {
		ssType = IPA
	} else if strings.Contains(nameLC, "sampa") {
		ssType = SAMPA
	} else if strings.Contains(nameLC, "cmu") {
		ssType = CMU
	}

	res := SymbolSet{
		Name:    name,
		Type:    ssType,
		Symbols: symbols,

		isInit: true,

		phonemes:        phonemes,
		phoneticSymbols: phoneticSymbols,
		stressSymbols:   stressSymbols,
		syllabic:        syllabic,
		nonSyllabic:     nonSyllabic,

		PhonemeRe:     phonemeRe,
		SyllabicRe:    syllabicRe,
		NonSyllabicRe: nonSyllabicRe,
		SymbolRe:      symbolRe,

		ipaSyllabicRe:    ipaSyllabicRe,
		ipaNonSyllabicRe: ipaNonSyllabicRe,
		ipaPhonemeRe:     ipaPhonemeRe,

		phonemeDelimiter:          phonemeDelimiter,
		phonemeDelimiterRe:        phonemeDelimiterRe,
		repeatedPhonemeDelimiters: repeatedPhonemeDelimiters,
	}
	return res, nil

}

// LoadSymbolSet loads a SymbolSet from file
func LoadSymbolSet(fName string) (SymbolSet, error) {
	name := filepath.Base(fName)
	var extension = filepath.Ext(name)
	name = name[0 : len(name)-len(extension)]
	return loadSymbolSet0(name, fName)
}

var header = "DESCRIPTION	SYMBOL	IPA	IPA UNICODE	CATEGORY"

// loadSymbolSet_ loads a SymbolSet from file
func loadSymbolSet0(name string, fName string) (SymbolSet, error) {
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
				ipaSym := IPASymbol{String: ipa, Unicode: ipaUnicode}
				sym := Symbol{
					String: symbol,
					Cat:    symCat,
					Desc:   desc,
					IPA:    ipaSym,
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

// LoadSymbolSetsFromDir loads a all symbol sets from the specified folder (all files with .tab extension)
func LoadSymbolSetsFromDir(dirName string) (map[string]SymbolSet, error) {
	// list files in symbol set dir
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, fmt.Errorf("failed reading symbol set dir : %v", err)
	}
	var fErrs error
	var symSets []SymbolSet
	for _, fi := range fileInfos {
		if strings.HasSuffix(fi.Name(), SymbolSetSuffix) {
			symset, err := LoadSymbolSet(filepath.Join(dirName, fi.Name()))
			if err != nil {
				if fErrs != nil {
					fErrs = fmt.Errorf("%v : %v", fErrs, err)
				} else {
					fErrs = err
				}
			} else {
				symSets = append(symSets, symset)
			}
		}
	}

	if fErrs != nil {
		return nil, fmt.Errorf("failed to load symbol set : %v", fErrs)
	}

	var symbolSetsMap = make(map[string]SymbolSet)
	for _, z := range symSets {
		// TODO checks that x.Name doesn't already exist ?
		if _, ok := symbolSetsMap[z.Name]; ok {
			// do nothing
		} else {
			symbolSetsMap[z.Name] = z
		}
	}
	return symbolSetsMap, nil
}

// LoadMapper loads a symbol set mapper from two SymbolSet instances
func LoadMapper(s1 SymbolSet, s2 SymbolSet) (Mapper, error) {
	fromName := s1.Name
	toName := s2.Name
	name := fromName + "_2_" + toName

	mapper := Mapper{name, s1, s2}

	var errs []string

	for _, symbol := range s1.Symbols {
		if len(symbol.String) > 0 {
			mapped, err := mapper.MapTranscription(symbol.String)
			if len(mapped) > 0 {
				if err != nil {
					return mapper, fmt.Errorf("couldn't test mapper: %v\n", err)
				}
			}
		}
	}
	if len(errs) > 0 {
		return mapper, fmt.Errorf("mapper initialization tests failed : %v", strings.Join(errs, "; "))
	}

	return mapper, nil
}

// LoadMapperFromFile loads two SymbolSet instances from files.
func LoadMapperFromFile(fromName string, toName string, fName1 string, fName2 string) (Mapper, error) {

	if fromName == toName {
		return Mapper{}, fmt.Errorf("should not load symbol sets with the same name: %s", fromName)
	}
	if fName1 == fName2 {
		return Mapper{}, fmt.Errorf("should not load both symbol sets from the same file: %s", fName1)
	}

	m1, err := loadSymbolSet0(fromName, fName1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mapper: %v\n", err)
		return Mapper{}, err
	}
	s2, err := loadSymbolSet0(toName, fName2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mapper: %v\n", err)
		return Mapper{}, err
	}
	return LoadMapper(m1, s2)
}