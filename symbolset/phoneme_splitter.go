package symbolset

import (
	"sort"
	"strings"
)

// Sort slice of strings according to len, longest string first
// TODO Should this live in a util lib?
type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func SplitIntoPhonemes(knownPhonemes []string, transcription string) (phonemes []string, unknown []string) {

	sort.Sort(ByLength(knownPhonemes))

	known := make([]string, len(knownPhonemes))
	copy(known, knownPhonemes)
	return splurt(&known, transcription, []string{}, []string{})
}

func splurt(srted *[]string, trans string, phs []string, unk []string) ([]string, []string) {

	if len(trans) > 0 {
		pre, rest, ok := consume(srted, trans)
		if ok { // known phoneme is prefix if trans
			return splurt(srted, rest, append(phs, pre), unk)
		} else { // unknown prefix, chopped off first rune
			return splurt(srted, rest, append(phs, pre), append(unk, pre))
		}
	} else {
		return phs, unk
	}

	// just to make compiler happy
	//fmt.Printf("SHOULD NEVER HAPPEN splurt: trans '%v' phs '%v' unk '%v'\n", trans, phs, unk)
	return phs, unk
}

func consume(srtd *[]string, trans string) (string, string, bool) {
	//fmt.Printf("consume trans '%v'\n", trans)
	var resPref string
	var resSuffix string

	var inTrans []string

	for _, ph := range *srtd {
		ind := strings.Index(trans, ph)
		if ind == 0 { // bingo
			resPref = ph
			resSuffix = trans[len(ph):]
			return resPref, resSuffix, true
			break
		}

		// These are the phonemes that are _somewhere_ in the input string.
		// Keep these, discard the rest
		if ind > -1 {
			inTrans = append(inTrans, ph)
		}
	}

	// Only keep substrings that we know are somewhere in trans
	*srtd = inTrans

	// no known phoneme prefix, separate first rune
	t := []rune(trans)
	return string(t[0]), string(t[1:]), false
}
