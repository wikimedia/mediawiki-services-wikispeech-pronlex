package symbolset2

import (
	"fmt"
	"regexp"
	"strings"
)

// refactoring, to replace package symbolset in the future

// structs in package symbolset

type SymbolSetType int

const (
	CMU SymbolSetType = iota
	SAMPA
	IPA
	Other
)

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

// IPA symbol with Unicode representation
type IPASymbol struct {
	String  string
	Unicode string
}

// Symbol represent a phoneme, stress or delimiter symbol used in transcriptions, including the IPA symbol with unicode
type Symbol struct {
	String string
	Cat    SymbolCat
	Desc   string
	IPA    IPASymbol
}

var ipaAccentI = "\u02C8"
var ipaAccentII = "\u0300"
var ipaLength = "\u02D0"

// ipa utilility functions with struct for package private usage.
// Symbols and codes: http://www.phon.ucl.ac.uk/home/wells/ipa-unicode.htm#numbers
type ipaFilter struct {
	accentI  string
	accentII string
	length   string
}

func (ipa ipaFilter) filterBeforeMappingFromIpa(trans string, ss SymbolSet) (string, error) {
	// IPA: ˈba`ŋ.ka => ˈ`baŋ.ka"
	// IPA: ˈɑ̀ː.pa => ˈ`ɑː.pa
	trans = strings.Replace(trans, ipa.accentII+ipa.length, ipa.length+ipa.accentII, -1)
	s := ipa.accentI + "(" + ss.PhonemeRe.String() + "+)" + ipa.accentII
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	res := repl.ReplaceAllString(trans, ipa.accentI+ipa.accentII+"$1")
	return res, nil
}

func (ipa ipaFilter) filterAfterMappingToIpa(trans string, ss SymbolSet) (string, error) {
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
		trans = res
	}
	trans = strings.Replace(trans, ipa.length+ipa.accentII, ipa.accentII+ipa.length, -1)
	return trans, nil
}

var cmuString = "cmu"

// cmu utilility functions with struct for package private usage.
type cmuFilter struct {
	cmu string
}

func (cmu cmuFilter) isCMU(symbolSetName string) bool {
	return strings.Contains(strings.ToLower(symbolSetName), cmu.cmu)
}

func (cmu cmuFilter) filterBeforeMappingFromCMU(trans string, ss SymbolSet) (string, error) {
	re, err := regexp.Compile("(.)([012])")
	if err != nil {
		return "", err
	}
	trans = re.ReplaceAllString(trans, "$1 $2")
	return trans, nil
}

func (cmu cmuFilter) filterAfterMappingToCMU(trans string, ss SymbolSet) (string, error) {
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
