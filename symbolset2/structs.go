// symbolset2 work in progress -- refactoring symbolset package
package symbolset2

// refactoring, to replace package symbolset in the future

// structs in package symbolset

type SymbolSetType int

const (
	Cmu SymbolSetType = iota
	Sampa
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

// Ipa symbol with Unicode representation
type Ipa struct {
	String  string
	Unicode string
}

// Symbol represent a phoneme, stress or delimiter symbol used in transcriptions, including the IPA symbol with unicode
type Symbol struct {
	String string
	Cat    SymbolCat
	Desc   string
	Ipa    Ipa
}
