package symbolset

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/lex"
)

// structs in package symbolset

// Symbols is a struct for package private usage.
// To create a new 'Symbols' instance, use NewSymbols
type Symbols struct {
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

// Mapper is a struct for package private usage. To create a new instance of Mapper, use LoadMapper.
type Mapper struct {
	Name       string
	SymbolSet1 SymbolSet
	SymbolSet2 SymbolSet
}

// SymbolSet is a struct for package private usage.
// To create a new SymbolSet, use NewSymbolSet.
type SymbolSet struct {
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

	From      Symbols
	To        Symbols
	symbolMap map[string]Symbol

	repeatedPhonemeDelimiters *regexp.Regexp
}

func (m SymbolSet) reverse(newName string) (SymbolSet, error) {
	var symbols = make([]SymbolPair, 0)

	for _, pair := range m.Symbols {
		s1 := pair.Sym1
		s2 := pair.Sym2
		symbols = append(symbols, SymbolPair{s2, s1})
	}
	return NewSymbolSet(newName, m.ToName, m.FromName, symbols)
}

func (m SymbolSet) preFilter(trans string, ss Symbols) (string, error) {
	if m.fromIsIPA {
		return m.ipa.filterBeforeMappingFromIpa(trans, ss)
	} else if m.fromIsCMU {
		return m.cmu.filterBeforeMappingFromCMU(trans, ss)
	}
	return trans, nil
}

func (m SymbolSet) postFilter(trans string, ss Symbols) (string, error) {
	if m.toIsIPA {
		return m.ipa.filterAfterMappingToIpa(trans, ss)
	} else if m.toIsCMU {
		return m.cmu.filterAfterMappingToCMU(trans, ss)
	}
	return trans, nil
}

// MapTranscriptions maps the input entry's transcriptions (in-place)
func (m SymbolSet) MapTranscriptions(e *lex.Entry) error {
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

// MapSymbol maps one symbol into the corresponding symbol in the new symbol set
func (m SymbolSet) MapSymbol(symbol Symbol) (Symbol, error) {
	res, ok := m.symbolMap[symbol.String]
	if !ok {
		return symbol, fmt.Errorf("unknown input symbol %v", symbol)
	}
	return res, nil
}

// MapSymbolString maps one symbol into the corresponding symbol in the new symbol set
func (m SymbolSet) MapSymbolString(symbol string) (string, error) {
	sym, err := m.From.Get(symbol)
	if err != nil {
		return symbol, err
	}
	res, err := m.MapSymbol(sym)
	return res.String, nil
}

// MapTranscription maps one input transcription string into the new symbol set.
func (m SymbolSet) MapTranscription(input string) (string, error) {
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
		to := m.symbolMap[from.String]
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
	res, err = m.postFilter(res, m.To)
	return res, err
}

// SymbolPair is a tuple inside the SymbolSet used to store the symbol mappings in a list in preserved order
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

//go:generate stringer -type=SymbolCat

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
	//IPA    string // ! TODO: should be like this
	Cat  SymbolCat
	Desc string
}

// Get searches the Symbols for a symbol with the given string
func (ss Symbols) Get(symbol string) (Symbol, error) {
	for _, s := range ss.Symbols {
		if s.String == symbol {
			return s, nil
		}
	}
	return Symbol{}, fmt.Errorf("no symbol /%s/ in symbol set", symbol)
}

func (ss Symbols) filterAmbiguous(trans []string) ([]string, error) {
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

func (ss Symbols) preCheckAmbiguous() error {
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
func (ss Symbols) ValidSymbol(symbol string) bool {
	return contains(ss.Symbols, symbol)
}

// SplitTranscription splits the input transcription into separate symbols
func (ss Symbols) SplitTranscription(input string) ([]string, error) {
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

func (ipa ipa) filterBeforeMappingFromIpa(trans string, ss Symbols) (string, error) {
	// IPA: ˈba`ŋ.ka => ˈ`baŋ.ka"
	s := ipa.accentI + "(" + ss.PhonemeRe.String() + "+)" + ipa.accentII
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	res := repl.ReplaceAllString(trans, ipa.accentI+ipa.accentII+"$1")
	return res, nil
}

func (ipa ipa) filterAfterMappingToIpa(trans string, ss Symbols) (string, error) {
	// IPA: /ə.ba⁀ʊˈt/ => /ə.ˈba⁀ʊt/
	s := "(" + ss.NonSyllabicRe.String() + "*)(" + ss.SyllabicRe.String() + ")" + ipa.accentI
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	trans = repl.ReplaceAllString(trans, ipa.accentI+"$1$2")

	// IPA: əs.ˈ̀̀e ...
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

func (cmu cmu) filterBeforeMappingFromCMU(trans string, ss Symbols) (string, error) {
	re, err := regexp.Compile("(.)([012])")
	if err != nil {
		return "", err
	}
	trans = re.ReplaceAllString(trans, "$1 $2")
	return trans, nil
}

func (cmu cmu) filterAfterMappingToCMU(trans string, ss Symbols) (string, error) {
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

// MapTranscription maps one input transcription string into the new symbol set.
func (m Mapper) MapTranscription(input string) (string, error) {
	res, err := m.SymbolSet1.MapTranscription(input)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription (1) : %v", err)
	}
	res, err = m.SymbolSet2.MapTranscription(res)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription (2) : %v", err)
	}
	return res, nil
}

// MapSymbol maps one input transcription symbol into the new symbol set.
func (m Mapper) MapSymbol(input Symbol) (Symbol, error) {
	res, err := m.SymbolSet1.MapSymbol(input)
	if err != nil {
		return Symbol{}, fmt.Errorf("couldn't map symbol (1) : %v", err)
	}
	res, err = m.SymbolSet2.MapSymbol(res)
	if err != nil {
		return Symbol{}, fmt.Errorf("couldn't map symbol (2) : %v", err)
	}
	return res, nil
}

// MapSymbolString maps one input transcription symbol into the new symbol set.
func (m Mapper) MapSymbolString(input string) (string, error) {
	res, err := m.SymbolSet1.MapSymbolString(input)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription : %v", err)
	}
	res, err = m.SymbolSet2.MapSymbolString(res)
	if err != nil {
		return "", fmt.Errorf("couldn't map transcription : %v", err)
	}
	return res, nil
}

// MapTranscriptions maps the input entry's transcriptions (in-place)
func (m Mapper) MapTranscriptions(e *lex.Entry) error {
	err := m.SymbolSet1.MapTranscriptions(e)
	if err != nil {
		return fmt.Errorf("couldn't map transcription : %v", err)
	}
	err = m.SymbolSet2.MapTranscriptions(e)
	if err != nil {
		return fmt.Errorf("couldn't map transcription : %v", err)
	}
	return nil
}
