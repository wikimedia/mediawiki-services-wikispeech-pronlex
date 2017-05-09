package symbolset

// symbol set filters for accent/stress placement

import (
	"fmt"
	"regexp"
	"strings"
)

func preFilter(ss SymbolSet, trans string, fromType Type) (string, error) {
	if fromType == IPA {
		return filterBeforeMappingFromIPA(ss, trans)
	} else if fromType == CMU {
		return filterBeforeMappingFromCMU(ss, trans)
	}
	return trans, nil
}

func postFilter(ss SymbolSet, trans string, toType Type) (string, error) {
	if toType == IPA {
		return filterAfterMappingToIPA(ss, trans)
	} else if toType == CMU {
		return filterAfterMappingToCMU(ss, trans)
	}
	return trans, nil
}

var ipaAccentI = "\u02C8"
var ipaAccentII = "\u0300"
var ipaSecStress = "\u02CC"
var ipaLength = "\u02D0"
var cmuString = "cmu"

func filterBeforeMappingFromIPA(ss SymbolSet, trans string) (string, error) {
	// IPA: ˈba`ŋ.ka => ˈ`baŋ.ka"
	// IPA: ˈɑ̀ː.pa => ˈ`ɑː.pa
	trans = strings.Replace(trans, ipaAccentII+ipaLength, ipaLength+ipaAccentII, -1)
	s := ipaAccentI + "(" + ss.ipaPhonemeRe.String() + "+)" + ipaAccentII
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	res := repl.ReplaceAllString(trans, ipaAccentI+ipaAccentII+"$1")
	return res, nil
}

func filterAfterMappingToIPA(ss SymbolSet, trans string) (string, error) {
	// IPA: /ə.ba⁀ʊˈt/ => /ə.ˈba⁀ʊt/
	s := "(" + ss.ipaNonSyllabicRe.String() + "*)(" + ss.ipaSyllabicRe.String() + ")" + ipaSecStress
	repl, err := regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	trans = repl.ReplaceAllString(trans, ipaSecStress+"$1$2")

	// Move sec stress to consonant cluster before the vowel
	s = "(" + ss.ipaNonSyllabicRe.String() + "*)(" + ss.ipaSyllabicRe.String() + ")" + ipaAccentI
	repl, err = regexp.Compile(s)
	if err != nil {
		return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
	}
	trans = repl.ReplaceAllString(trans, ipaAccentI+"$1$2")

	// IPA: əs.ˈ̀̀e ...
	// IPA: /'`pa.pa/ => /'pa`.pa/
	accentIIConditionForAfterMapping := ipaAccentI + ipaAccentII
	if strings.Contains(trans, accentIIConditionForAfterMapping) {
		s := ipaAccentI + ipaAccentII + "(" + ss.ipaNonSyllabicRe.String() + "*)(" + ss.ipaSyllabicRe.String() + ")"
		repl, err := regexp.Compile(s)
		if err != nil {
			return "", fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
		}
		res := repl.ReplaceAllString(trans, ipaAccentI+"$1$2"+ipaAccentII)
		trans = res
	}
	// IPA: /'paː`.pa/ => /'pa`ː.pa/
	trans = strings.Replace(trans, ipaLength+ipaAccentII, ipaAccentII+ipaLength, -1)
	return trans, nil
}

func filterBeforeMappingFromCMU(ss SymbolSet, trans string) (string, error) {
	re, err := regexp.Compile("(.)([012])")
	if err != nil {
		return "", err
	}
	trans = re.ReplaceAllString(trans, "$1 $2")
	return trans, nil
}

func filterAfterMappingToCMU(ss SymbolSet, trans string) (string, error) {
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
