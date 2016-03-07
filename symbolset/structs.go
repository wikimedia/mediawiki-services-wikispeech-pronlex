package symbolset

import (
	//"regexp"
	"strings"
)

type PhonemeType int
type DelimiterType int

const (
	SyllabicT PhonemeType = iota
	NonSyllabicT
	StressT
)
const (
	PhonemeDelimiterT DelimiterType = iota
	ExplicitPhonemeDelimiterT
	SyllableDelimiterT
	MorphemeDelimiterT
	WordDelimiterT
)

type Symbol interface {
	S() string
	Desc() string
}

type Phoneme interface {
	S() string
	Desc() string
	Type() PhonemeType
}

type Delimiter interface {
	S() string
	Desc() string
	Type() DelimiterType
}

type Syllabic struct {
	symbol string
	desc   string
}

func (s Syllabic) S() string {
	return s.symbol
}
func (s Syllabic) Desc() string {
	return s.desc
}
func (s Syllabic) Type() PhonemeType {
	return SyllabicT
}

type NonSyllabic struct {
	symbol string
	desc   string
}

func (s NonSyllabic) S() string {
	return s.symbol
}
func (s NonSyllabic) Desc() string {
	return s.desc
}
func (s NonSyllabic) Type() PhonemeType {
	return NonSyllabicT
}

type Stress struct {
	symbol string
	desc   string
}

func (s Stress) S() string {
	return s.symbol
}
func (s Stress) Desc() string {
	return s.desc
}
func (s Stress) Type() PhonemeType {
	return StressT
}

type PhonemeDelimiter struct {
	symbol string
	desc   string
}

func (s PhonemeDelimiter) S() string {
	return s.symbol
}
func (s PhonemeDelimiter) Desc() string {
	return s.desc
}
func (s PhonemeDelimiter) Type() DelimiterType {
	return PhonemeDelimiterT
}

type ExplicitPhonemeDelimiter struct {
	symbol string
	desc   string
}

func (s ExplicitPhonemeDelimiter) S() string {
	return s.symbol
}
func (s ExplicitPhonemeDelimiter) Desc() string {
	return s.desc
}
func (s ExplicitPhonemeDelimiter) Type() DelimiterType {
	return ExplicitPhonemeDelimiterT
}

type SyllableDelimiter struct {
	symbol string
	desc   string
}

func (s SyllableDelimiter) S() string {
	return s.symbol
}
func (s SyllableDelimiter) Desc() string {
	return s.desc
}
func (s SyllableDelimiter) Type() DelimiterType {
	return SyllableDelimiterT
}

type MorphemeDelimiter struct {
	symbol string
	desc   string
}

func (s MorphemeDelimiter) S() string {
	return s.symbol
}
func (s MorphemeDelimiter) Desc() string {
	return s.desc
}
func (s MorphemeDelimiter) Type() DelimiterType {
	return MorphemeDelimiterT
}

type WordDelimiter struct {
	symbol string
	desc   string
}

func (s WordDelimiter) S() string {
	return s.symbol
}
func (s WordDelimiter) Desc() string {
	return s.desc
}
func (s WordDelimiter) Type() DelimiterType {
	return WordDelimiterT
}

type SymbolSet struct {
	Name    string
	Symbols []Symbol
}

// SYMBOLS: http://www.phon.ucl.ac.uk/home/wells/ipa-unicode.htm#numbers
type IPA struct {
	ipa      string
	accentI  string
	accentII string
}

func NewIPA() IPA {
	return IPA{
		ipa:      "ipa",
		accentI:  "\u02C8",
		accentII: "\u0300",
	}
}
func (ipa IPA) IsIPA(symbolSetName string) bool {
	return strings.Contains(strings.ToLower(symbolSetName), ipa.ipa)
}

//   def filterBeforeMapping(trans: String, ss: SymbolSet): String = {
//     // IPA: ˈba`ŋ.ka => ˈ`baŋ.ka"
//     require(ss.stressSymbols.contains(IPA.accentI), "No IPA stress symbol in stress symbol list? IPA stress =/" + IPA.accentI + "/, stress symbols=" + ss.stressSymbols)
//     require(ss.stressSymbols.contains(IPA.accentI+IPA.accentII), "No IPA tone II symbol in stress symbol list? IPA stress =/" + IPA.accentI + "/, stress symbols=" + ss.stressSymbols)
//     val replacementRe = (IPA.accentI + "(" + ss.phonemeRe +"+)" + IPA.accentII).r
//     replacementRe.replaceAllIn(trans, IPA.accentI + IPA.accentII + "$1")
//   }

//   val conditionForAfterMapping = IPA.accentI + IPA.accentII
//   def filterAfterMapping(trans: String, ss: SymbolSet): String = {
//     // IPA: /'`pa.pa/ => /'pa`.pa/
//     if (trans.contains(conditionForAfterMapping)) {
//       require(ss.stressSymbols.contains(IPA.accentI), "No IPA stress symbol in stress symbol list? IPA stress =/" + IPA.accentI + "/, stress symbols=" + ss.stressSymbols)
//       require(ss.stressSymbols.contains(IPA.accentI+IPA.accentII), "No IPA tone II symbol in stress symbol list? IPA stress =/" + IPA.accentI + "/, stress symbols=" + ss.stressSymbols)
//       val replacementRe = (IPA.accentI + IPA.accentII +"(" + ss.nonSyllabicRe +"*)(" + ss.syllabicRe + ")").r
//       replacementRe.replaceAllIn(trans, IPA.accentI + "$1$2" + IPA.accentII)
//     }
//     else trans
//   }

// }
