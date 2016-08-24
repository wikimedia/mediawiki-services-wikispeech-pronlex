package symbolset

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/lex"
)

// structs in package symbolset

// SymbolSet is a struct for package private usage.
// To create a new SymbolSet, use NewSymbolSet
type SymbolSet struct {
	Name    string
	Symbols []Symbol

	// to check if the struct has been initialized properly
	isInit bool

	// derived values computed upon initialization
	phonemes        []Symbol
	phoneticSymbols []Symbol
	stressSymbols   []Symbol
	syllabic        []Symbol
	nonSyllabic     []Symbol

	phonemeDelimiter Symbol

	PhonemeRe          *regexp.Regexp
	SyllabicRe         *regexp.Regexp
	NonSyllabicRe      *regexp.Regexp
	SymbolRe           *regexp.Regexp
	phonemeDelimiterRe *regexp.Regexp
}

// Mapper is a struct for package private usage.
// To create a new Mapper, use NewMapper.
type Mapper struct {
	Name     string
	FromName string
	ToName   string
	Symbols  []SymbolPair

	ipa ipa
	cmu cmu

	// derived values computed upon initialization
	fromIsIPA bool
	toIsIPA   bool
	fromIsCMU bool
	toIsCMU   bool

	From      SymbolSet
	To        SymbolSet
	symbolMap map[Symbol]Symbol

	repeatedPhonemeDelimiters *regexp.Regexp
}

func (m Mapper) preFilter(trans string, ss SymbolSet) (string, error) {
	if m.fromIsIPA {
		return m.ipa.filterBeforeMappingFromIpa(trans, ss)
	} else if m.fromIsCMU {
		return m.cmu.filterBeforeMappingFromCMU(trans, ss), nil
	}
	return trans, nil
}

func (m Mapper) postFilter(trans string, ss SymbolSet) (string, error) {
	if m.toIsIPA {
		return m.ipa.filterAfterMappingToIpa(trans, ss)
	} else if m.toIsCMU {
		return m.cmu.filterAfterMappingToCMU(trans, ss)
	}
	return trans, nil
}

// MapTranscriptions maps the input entry's transcriptions (in-place)
func (m Mapper) MapTranscriptions(e *lex.Entry) error {
	var newTs []lex.Transcription
	var errs []string
	for _, t := range e.Transcriptions {
		newT, err := m.MapTranscription(t.Strn)
		if err != nil {
			errs = append(errs, err.Error())
		}
		newTs = append(newTs, lex.Transcription{ID: t.ID, Strn: newT, EntryID: t.EntryID, Language: t.Language, Sources: t.Sources})
	}
	e.Transcriptions = newTs
	if len(errs) > 0 {
		return fmt.Errorf("%v", strings.Join(errs, "; "))
	}
	return nil
}

// MapTranscription maps one input transcription string into the new symbol set.
func (m Mapper) MapTranscription(input string) (string, error) {
	res, err := m.preFilter(input, m.From)
	if err != nil {
		return "", err
	}
	splitted, err := m.From.SplitTranscription(res)
	if err != nil {
		return "", err
	}
	var mapped = make([]string, 0)
	for _, fromS := range splitted {
		from, err := m.From.Get(fromS)
		if err != nil {
			return "", fmt.Errorf("input symbol /%s/ is undefined : %v", fromS, err)
		}
		to := m.symbolMap[from]
		//if to.Cat == UndefinedSymbol {
		//	return "", fmt.Errorf("couldn't map symbol /%s/", fromS)
		//}
		if len(to.String) > 0 {
			mapped = append(mapped, to.String)
		}
	}
	//mapped, err = m.to.filterAmbiguous(mapped)
	//if err != nil {
	//	return "", err
	//}
	res = strings.Join(mapped, m.To.phonemeDelimiter.String)

	// remove repeated phoneme delimiters
	res = m.repeatedPhonemeDelimiters.ReplaceAllString(res, m.To.phonemeDelimiter.String)
	return m.postFilter(res, m.To)
}

// SymbolPair is a tuple inside the Mapper used to store the symbol mappings in a list in preserved order
type SymbolPair struct {
	Sym1 Symbol
	Sym2 Symbol
}

// symbolSlice is used for sorting slices of symbols according to symbol length. Symbols with equal length will be sorted alphabetically.
type symbolSlice []Symbol

func (a symbolSlice) Len() int      { return len(a) }
func (a symbolSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a symbolSlice) Less(i, j int) bool {
	if len(a[i].String) != len(a[j].String) {
		return len(a[i].String) > len(a[j].String)
	}
	return a[i].String < a[j].String
}

// SymbolCat is used to categorize transcription symbols.
type SymbolCat int

const (
	// Syllabic is used for syllabic phonemes (typically vowels and syllabic consonants)
	Syllabic SymbolCat = iota

	// NonSyllabic is used for non-syllabic phonemes (typically consonants)
	NonSyllabic

	// Stress is used for stress and accent symbols (primary, secondary, tone accents, etc)
	Stress

	// PhonemeDelimiter is used for phoneme delimiters (white space, empty string, etc)
	PhonemeDelimiter

	// SyllableDelimiter is used for syllable delimiters
	SyllableDelimiter

	// MorphemeDelimiter is used for morpheme delimiters that need not align with
	// morpheme boundaries in the decompounded orthography
	MorphemeDelimiter

	// CompoundDelimiter is used for compound delimiters that should be aligned
	// with compound boundaries in the decompounded orthography
	CompoundDelimiter

	// WordDelimiter is used for word delimiters
	WordDelimiter
)

// Symbol represent a phoneme, stress or delimiter symbol used in transcriptions
type Symbol struct {
	String string
	Cat    SymbolCat
	Desc   string
}

// Contains checks if the SymbolSet contains a certain symbol string
func (ss SymbolSet) Contains(symbol string) bool {
	return contains(ss.Symbols, symbol)
}

// Get searches the SymbolSet for a symbol with the given string
func (ss SymbolSet) Get(symbol string) (Symbol, error) {
	for _, s := range ss.Symbols {
		if s.String == symbol {
			return s, nil
		}
	}
	return Symbol{}, fmt.Errorf("no symbol /%s/ in symbol set", symbol)
}

func (ss SymbolSet) filterAmbiguous(trans []string) ([]string, error) {
	potentiallyAmbs := ss.phoneticSymbols
	phnDel := ss.phonemeDelimiter.String
	var res = make([]string, 0)
	for i, phn0 := range trans[0 : len(trans)-1] {
		phn1 := trans[i+1]
		test := phn0 + phnDel + phn1
		if contains(potentiallyAmbs, test) {
			return nil, fmt.Errorf("symbol set %s contains ambiguous symbols: [%s] vs [%s] + [%s]", ss.Name, (phn0 + phn1), phn0, phn1)
		}
		res = append(res, phn0)
	}
	res = append(res, trans[len(trans)-1])
	return res, nil
}

func (ss SymbolSet) preCheckAmbiguous() error {
	allSymbols := ss.phonemes
	var res = make([]string, 0)
	for _, a := range allSymbols {
		for _, b := range allSymbols {
			if len(a.String) > 0 && len(b.String) > 0 {
				res = append(res, a.String)
				res = append(res, b.String)
			}
		}
	}
	_, err := ss.filterAmbiguous(res)
	return err
}

// ValidSymbol checks if a string is a valid symbol or not
func (ss SymbolSet) ValidSymbol(symbol string) bool {
	for _, s := range ss.Symbols {
		if s.String == symbol {
			return true
		}
	}
	return false
}

// SplitTranscription splits the input transcription into separate symbols
func (ss SymbolSet) SplitTranscription(input string) ([]string, error) {
	if !ss.isInit {
		panic("symbolSet " + ss.Name + " has not been initialized properly!")
	}
	delim := ss.phonemeDelimiterRe
	if delim.FindStringIndex("") != nil {
		rest := input
		var acc = make([]string, 0)
		for len(rest) > 0 {
			match := ss.SymbolRe.FindStringIndex(rest)
			switch match {
			case nil:
				return nil, fmt.Errorf("transcription not splittable (invalid symbols?)! input=/%s/, acc=/%s/, rest=/%s/", input, strings.Join(acc, ss.phonemeDelimiter.String), rest)

			default:
				if match[0] != 0 {
					return nil, fmt.Errorf("couldn't parse transcription /%s/, it may contain undefined symbols", input)
				}
				acc = append(acc, rest[match[0]:match[1]])
				rest = rest[match[1]:]
			}
		}
		return acc, nil
	}
	return delim.Split(input, -1), nil
}

// ipa utilility functions with struct for package private usage.
// Symbols and codes: http://www.phon.ucl.ac.uk/home/wells/ipa-unicode.htm#numbers
type ipa struct {
	ipa      string
	accentI  string
	accentII string
}

func (ipa ipa) isIPA(symbolSetName string) bool {
	return strings.Contains(strings.ToLower(symbolSetName), ipa.ipa)
}

func (ipa ipa) filterBeforeMappingFromIpa(trans string, ss SymbolSet) (string, error) {
	// IPA: ˈba`ŋ.ka => ˈ`baŋ.ka"
	s := ipa.accentI + "(" + ss.PhonemeRe.String() + "+)" + ipa.accentII
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	res := repl.ReplaceAllString(trans, ipa.accentI+ipa.accentII+"$1")
	return res, nil
}

func (ipa ipa) filterAfterMappingToIpa(trans string, ss SymbolSet) (string, error) {
	// IPA: /ə.ba⁀ʊˈt/ => /ə.ˈba⁀ʊt/
	s := "(" + ss.NonSyllabicRe.String() + "*)(" + ss.SyllabicRe.String() + ")" + ipa.accentI
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	trans = repl.ReplaceAllString(trans, ipa.accentI+"$1$2")

	// IPA: /'`pa.pa/ => /'pa`.pa/
	accentIIConditionForAfterMapping := ipa.accentI + ipa.accentII
	if strings.Contains(trans, accentIIConditionForAfterMapping) {
		s := ipa.accentI + ipa.accentII + "(" + ss.NonSyllabicRe.String() + "*)(" + ss.SyllabicRe.String() + ")"
		repl, err := regexp.Compile(s)
		if err != nil {
			return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
		}
		res := repl.ReplaceAllString(trans, ipa.accentI+"$1$2"+ipa.accentII)
		return res, nil
	}
	return trans, nil
}

// cmu utilility functions with struct for package private usage.
type cmu struct {
	cmu string
}

func (cmu cmu) isCMU(symbolSetName string) bool {
	return strings.Contains(strings.ToLower(symbolSetName), cmu.cmu)
}

func (cmu cmu) filterBeforeMappingFromCMU(trans string, ss SymbolSet) string {
	trans = strings.Replace(trans, "1", " 1", -1)
	trans = strings.Replace(trans, "2", " 2", -1)
	trans = strings.Replace(trans, "0", " 0", -1)
	return trans
}

func (cmu cmu) filterAfterMappingToCMU(trans string, ss SymbolSet) (string, error) {
	s := "([012]) ((?:" + ss.NonSyllabicRe.String() + " )*)(" + ss.SyllabicRe.String() + ")"
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	trans = repl.ReplaceAllString(trans, "$2$3$1")

	trans = strings.Replace(trans, " 1", "1", -1)
	trans = strings.Replace(trans, " 2", "2", -1)
	trans = strings.Replace(trans, " 0", "0", -1)
	return trans, nil
}
