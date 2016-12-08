package symbolset2

import (
	"fmt"
	"regexp"
	"strings"
)

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

	PhonemeRe     *regexp.Regexp
	SyllabicRe    *regexp.Regexp
	NonSyllabicRe *regexp.Regexp
	SymbolRe      *regexp.Regexp

	phonemeDelimiter          Symbol
	phonemeDelimiterRe        *regexp.Regexp
	repeatedPhonemeDelimiters *regexp.Regexp
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

// GetFromIPA searches the SymbolSet for a symbol with the given IPA symbol string
func (ss SymbolSet) GetFromIPA(symbol string) (Symbol, error) {
	for _, s := range ss.Symbols {
		if s.IPA.String == symbol {
			return s, nil
		}
	}
	return Symbol{}, fmt.Errorf("no symbol /%s/ in symbol set", symbol)
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

// func (ss SymbolSet) preFilter(trans string, fromType SymbolSetType) (string, error) {
// 	if fromType == IPA {
// 		return ipaFilter.filterBeforeMappingFromIpa(trans, ss)
// 	} else if fromType == CMU {
// 		return cmuFilter.filterBeforeMappingFromCMU(trans, ss)
// 	}
// 	return trans, nil
// }

// func (ss SymbolSet) postFilter(trans string, toType SymbolSetType) (string, error) {
// 	if toType == IPA {
// 		return ipaFilter.filterAfterMappingToIpa(trans, ss)
// 	} else if toType == CMU {
// 		return cmuFilter.filterAfterMappingToCMU(trans, ss)
// 	}
// 	return trans, nil
// }

// ConvertToIPA maps one input transcription string into an IPA transcription
func (ss SymbolSet) ConvertToIPA(trans string) (string, error) {
	res := trans
	//res, err := ss.preFilter(input, ss.From)
	// if err != nil {
	// 	return "", err
	// }
	splitted, err := ss.SplitTranscription(res)
	if err != nil {
		return "", err
	}
	var mapped = make([]string, 0)
	for _, fromS := range splitted {
		symbol, err := ss.Get(fromS)
		if err != nil {
			return "", fmt.Errorf("input symbol /%s/ is undefined : %v", fromS, err)
		}
		to := symbol.IPA.String
		if len(to) > 0 {
			mapped = append(mapped, to)
		}
	}
	res = strings.Join(mapped, ss.phonemeDelimiter.IPA.String)

	//res, err = ss.postFilter(res, ss.To)
	return res, err
}

// ConvertFromIPA maps one input IPA transcription into the current symbol set
func (ss SymbolSet) ConvertFromIPA(trans string) (string, error) {
	res := trans
	//res, err := ss.preFilter(trans, ss.From)
	// if err != nil {
	// 	return "", err
	// }
	splitted, err := ss.SplitTranscription(res)
	if err != nil {
		return "", err
	}
	var mapped = make([]string, 0)
	for _, fromS := range splitted {
		symbol, err := ss.GetFromIPA(fromS)
		if err != nil {
			return "", fmt.Errorf("input symbol /%s/ is undefined : %v", fromS, err)
		}
		to := symbol.String
		if len(to) > 0 {
			mapped = append(mapped, to)
		}
	}
	res = strings.Join(mapped, ss.phonemeDelimiter.String)

	// remove repeated phoneme delimiters, if any
	res = ss.repeatedPhonemeDelimiters.ReplaceAllString(res, ss.phonemeDelimiter.IPA.String)
	//res, err = ss.postFilter(res, ss.To)
	return res, err
}
