package symbolset

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// inits.go Initialization functions for structs in package symbolset

// newIPA is a package private contructor for the ipa struct with fixed-value fields
func newIPA() ipa {
	return ipa{
		ipa:      "ipa",
		accentI:  "\u02C8",
		accentII: "\u0300",
	}
}

// newCMU is a package private contructor for the ipa struct with fixed-value fields
func newCMU() cmu {
	return cmu{
		cmu: "cmu",
	}
}

// NewSymbols is a constructor for 'symbols' with built-in error checks
func NewSymbols(name string, symbols []Symbol) (Symbols, error) {
	return NewSymbolsWithTests(name, symbols, true)
}

// NewsymbolsWithTests is a constructor for 'symbols' with built-in error checks
func NewSymbolsWithTests(name string, symbols []Symbol, checkForDups bool) (Symbols, error) {
	var nilRes Symbols

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

	res := Symbols{
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

// NewSymbolSet is a public constructor for SymbolSet with built-in error checks
func NewSymbolSet(name string, fromName string, toName string, symbolList []SymbolPair) (SymbolSet, error) {
	var nilRes SymbolSet

	ipa := newIPA()
	cmu := newCMU()

	toIsIPA := ipa.isIPA(toName)
	fromIsIPA := ipa.isIPA(fromName)
	toIsCMU := cmu.isCMU(toName)
	fromIsCMU := cmu.isCMU(fromName)

	if !(fromIsIPA || toIsIPA) {
		return nilRes, fmt.Errorf("one of the defined symbol sets must always be IPA -- found %v - %v", fromName, toName)
	}
	var fromSymbols = make([]Symbol, 0)
	var toSymbols = make([]Symbol, 0)
	var symbolMap = make(map[Symbol]Symbol)

	for _, pair := range symbolList {
		s1 := pair.Sym1
		s2 := pair.Sym2
		symbolMap[s1] = s2
		fromSymbols = append(fromSymbols, s1)
		toSymbols = append(toSymbols, s2)
	}

	from, err := NewSymbolsWithTests(fromName, fromSymbols, true) // check for duplicates in input symbol set
	if err != nil {
		return nilRes, err
	}
	to, err := NewSymbolsWithTests(toName, toSymbols, false) // do not check for duplicates in output phoneme set
	if err != nil {
		return nilRes, err
	}
	if from.Name == to.Name {
		return nilRes, fmt.Errorf("both phoneme sets cannot have the same name: %s", from.Name)
	}
	err = from.preCheckAmbiguous()
	if err != nil {
		return nilRes, err
	}
	err = to.preCheckAmbiguous()
	if err != nil {
		return nilRes, err
	}

	repeatedPhonemeDelimiters, err := regexp.Compile(to.phonemeDelimiterRe.String() + "+")
	if err != nil {
		return nilRes, err
	}

	m := SymbolSet{
		Name:                      name,
		FromName:                  fromName,
		ToName:                    toName,
		Symbols:                   symbolList,
		fromIsIPA:                 fromIsIPA,
		toIsIPA:                   toIsIPA,
		fromIsCMU:                 fromIsCMU,
		toIsCMU:                   toIsCMU,
		From:                      from,
		To:                        to,
		ipa:                       ipa,
		symbolMap:                 symbolMap,
		repeatedPhonemeDelimiters: repeatedPhonemeDelimiters,
	}
	return m, nil

}

// LoadMapper loads a 'mapper pair' from two SymbolSet instances
func LoadMapper(m1 SymbolSet, m2 SymbolSet) (Mapper, error) {
	fromName := m1.Name
	toName := m2.Name
	name := fromName + "2" + toName
	mappers := Mapper{name, m1, m2}

	// for testing:
	m1rev, err := m1.reverse(fromName + "2IPA")
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mapper: %v\n", err)
	}

	m2rev, err := m2.reverse("IPA2" + toName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mapper: %v\n", err)
	}
	mappersrev := Mapper{toName + "2" + fromName, m2rev, m1rev}

	var errs []string

	for _, symbol := range m1.From.Symbols {
		if len(symbol.String) > 0 {
			mapped, err := mappers.MapTranscription(symbol.String)
			if len(mapped) > 0 {
				if err != nil {
					return mappers, fmt.Errorf("couldn't test mapper: %v\n", err)
				}
				mapped, err = mappersrev.MapTranscription(mapped)
				if err != nil {
					return mappers, fmt.Errorf("couldn't test mapper: %v\n", err)
				}
				if mapped != symbol.String {
					errs = append(errs, "couldn't map /"+symbol.String+"/ back and forth -- got /"+mapped+"/")
				}
			}
		}
	}
	if len(errs) > 0 {
		return mappers, fmt.Errorf("Mapper initialization tests failed %v", strings.Join(errs, "; "))
	}

	return mappers, nil
}

// LoadMapperFromFile loads two SymbolSet instances from files.
func LoadMapperFromFile(fromName string, toName string, fName1 string, fName2 string) (Mapper, error) {
	m1, err := loadSymbolSet_(fromName+"2IPA", fName1, fromName, "IPA")
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mapper: %v\n", err)
		return Mapper{"", SymbolSet{}, SymbolSet{}}, err
	}
	m2, err := loadSymbolSet_("IPA2"+toName, fName2, "IPA", toName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mapper: %v\n", err)
		return Mapper{"", SymbolSet{}, SymbolSet{}}, err
	}
	return LoadMapper(m1, m2)
}

// LoadSymbolSet loads a SymbolSet from file
func LoadSymbolSet(fName string) (SymbolSet, error) {
	name := filepath.Base(fName)
	return loadSymbolSet_(name, fName, "", "")
}

// loadSymbolSet_ loads a SymbolSet from file
func loadSymbolSet_(name string, fName string, fromColumn string, toColumn string) (SymbolSet, error) {
	var nilRes SymbolSet
	fh, err := os.Open(fName)
	defer fh.Close()
	if err != nil {
		return nilRes, err
	}
	s := bufio.NewScanner(fh)
	n := 0
	var descIndex = -1
	var fromIndex = -1
	var toIndex = -1
	var typeIndex = -1
	if fromColumn == "" {
		fromIndex = 1
		fromColumn = "SYMBOL"
	}
	if toColumn == "" {
		toIndex = 2
		fromColumn = "IPA"
	}
	var maptable = make([]SymbolPair, 0)
	for s.Scan() {
		if err := s.Err(); err != nil {
			return nilRes, err
		}
		n++
		l := s.Text()
		if len(strings.TrimSpace(l)) > 0 && !strings.HasPrefix(strings.TrimSpace(l), "#") {
			fs := strings.Split(l, "\t")
			if n == 1 { // header
				descIndex = indexOf(fs, "DESCRIPTION")
				if fromIndex == -1 {
					fromIndex = indexOf(fs, fromColumn)
					if fromIndex == -1 {
						return nilRes, fmt.Errorf("from index %v undefined", fromColumn)
					}
				}
				if toIndex == -1 {
					toIndex = indexOf(fs, toColumn)
					if toIndex == -1 {
						return nilRes, fmt.Errorf("to index %v undefined", toColumn)
					}
				}
				typeIndex = indexOf(fs, "CATEGORY")

			} else {
				if descIndex == -1 {
					return nilRes, fmt.Errorf("%v", "description index unset")
				}
				if fromIndex == -1 {
					return nilRes, fmt.Errorf("%v", "from index unset")
				}
				if toIndex == -1 {
					return nilRes, fmt.Errorf("%v", "to index unset")
				}
				if typeIndex == -1 {
					return nilRes, fmt.Errorf("%v", "type index unset")
				}
				from := fs[fromIndex]
				to := fs[toIndex]
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
				symFrom := Symbol{String: from, Cat: symCat, Desc: desc}
				symTo := Symbol{String: to, Cat: symCat, Desc: desc}
				maptable = append(maptable, SymbolPair{symFrom, symTo})
			}
		}
	}
	fromName := ""
	toName := ""
	if fromColumn == "SYMBOL" {
		fromName = name
	} else {
		fromName = fromColumn
	}
	if toColumn == "SYMBOL" {
		toName = name
	} else {
		toName = toColumn
	}
	m, err := NewSymbolSet(name, fromName, toName, maptable)
	if err != nil {
		return nilRes, fmt.Errorf("couldn't load mapper from file %v : %v", fName, err)
	}
	return m, nil
}

// end: initialization
