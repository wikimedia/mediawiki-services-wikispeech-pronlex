package symbolset2

// refactoring, to replace package symbolset in the future

import (
	"fmt"
	"regexp"
)

// structs in package symbolset

type SymbolSetType int

const (
	CMU SymbolSetType = iota
	SAMPA
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
type IPA struct {
	String  string
	Unicode string
}

// Symbol represent a phoneme, stress or delimiter symbol used in transcriptions, including the IPA symbol with unicode
type Symbol struct {
	String string
	Cat    SymbolCat
	Desc   string
	IPA    IPA
}

// SymbolSet is a struct for package private usage.
// To create a new 'SymbolSet' instance, use NewSymbolSet
type SymbolSet struct {
	Name    string
	Type    SymbolSetType
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

// SplitTranscription splits the input transcription into separate symbols
func (ss SymbolSet) SplitTranscription(input string) ([]string, error) {
	if !ss.isInit {
		panic("symbolSet " + ss.Name + " has not been initialized properly!")
	}
	delim := ss.phonemeDelimiterRe
	if delim.FindStringIndex("") != nil {
		splitted, unknown, err := splitIntoPhonemes(ss.Symbols, input)
		if err != nil {
			return []string{}, err
		}
		if len(unknown) > 0 {
			return []string{}, fmt.Errorf("found unknown phonemes in transcription /%s/: %v\n", input, unknown)
		}
		return splitted, nil
	}
	return delim.Split(input, -1), nil
}

// SplitIPATranscription splits the input transcription into separate symbols
func (ss SymbolSet) SplitIPATranscription(input string) ([]string, error) {
	if !ss.isInit {
		panic("symbolSet " + ss.Name + " has not been initialized properly!")
	}
	symbols := []Symbol{}
	for _, s := range ss.Symbols {
		ipa := s
		ipa.String = ipa.IPA.String
		symbols = append(symbols, ipa)
	}
	splitted, unknown, err := splitIntoPhonemes(symbols, input)
	if err != nil {
		return []string{}, err
	}
	if len(unknown) > 0 {
		return []string{}, fmt.Errorf("found unknown phonemes in transcription /%s/: %v\n", input, unknown)
	}
	return splitted, nil
}

// func (ss SymbolSet) preFilter(trans string, ss Symbols) (string, error) {
// 	if ss.fromIsIPA {
// 		return ss.ipa.filterBeforeMappingFromIpa(trans, ss)
// 	} else if ss.fromIsCMU {
// 		return ss.cmu.filterBeforeMappingFromCMU(trans, ss)
// 	}
// 	return trans, nil
// }

// func (m SymbolSet) postFilter(trans string, ss Symbols) (string, error) {
// 	if m.toIsIPA {
// 		return m.ipa.filterAfterMappingToIpa(trans, ss)
// 	} else if m.toIsCMU {
// 		return m.cmu.filterAfterMappingToCMU(trans, ss)
// 	}
// 	return trans, nil
// }

// // MapTranscription maps one input transcription string into the new symbol set.
// func (ss SymbolSet) MapTranscriptionToIpa(input string) (string, error) {
// 	res, err := ss.preFilter(input, ss.From)
// 	if err != nil {
// 		return "", err
// 	}
// 	splitted, err := ss.SplitTranscription(res)
// 	if err != nil {
// 		return "", err
// 	}
// 	var mapped = make([]string, 0)
// 	for _, fromS := range splitted {
// 		from, err := m.From.Get(fromS)
// 		if err != nil {
// 			return "", fmt.Errorf("input symbol /%s/ is undefined : %v", fromS, err)
// 		}
// 		to := m.symbolMap[from.String]
// 		if len(to.String) > 0 {
// 			mapped = append(mapped, to.String)
// 		}
// 	}
// 	res = strings.Join(mapped, m.To.phonemeDelimiter.String)

// 	// remove repeated phoneme delimiters
// 	res = m.repeatedPhonemeDelimiters.ReplaceAllString(res, m.To.phonemeDelimiter.String)
// 	res, err = m.postFilter(res, m.To)
// 	return res, err
// }

// // MapTranscription maps one input transcription string into the new symbol set.
// func (ss SymbolSet) MapTranscriptionFromIpa(input string) (string, error) {
// 	res, err := ss.preFilter(input, ss.From)
// 	if err != nil {
// 		return "", err
// 	}
// 	splitted, err := ss.SplitTranscription(res)
// 	if err != nil {
// 		return "", err
// 	}
// 	var mapped = make([]string, 0)
// 	for _, fromS := range splitted {
// 		from, err := m.From.Get(fromS)
// 		if err != nil {
// 			return "", fmt.Errorf("input symbol /%s/ is undefined : %v", fromS, err)
// 		}
// 		to := m.symbolMap[from.String]
// 		if len(to.String) > 0 {
// 			mapped = append(mapped, to.String)
// 		}
// 	}
// 	res = strings.Join(mapped, m.To.phonemeDelimiter.String)

// 	// remove repeated phoneme delimiters
// 	res = m.repeatedPhonemeDelimiters.ReplaceAllString(res, m.To.phonemeDelimiter.String)
// 	res, err = m.postFilter(res, m.To)
// 	return res, err
// }
