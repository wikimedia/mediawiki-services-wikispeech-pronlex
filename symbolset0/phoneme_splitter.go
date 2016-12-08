package symbolset0

// For testing or standalone use only! In production use symbolset.SymbolSet.SplitTranscription

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/stts-se/pronlex/lex"
)

// Sort slice of strings according to len, longest string first
// TODO Should this live in a util lib?
type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) > len(s[j])
}

func splitIntoPhonemes(knownPhonemes []string, transcription string) (phonemes []string, unknown []string) {

	var known []string
	// start by discarding any phoneme strings not substrings of transcription
	for _, ph := range knownPhonemes {
		if strings.Index(transcription, ph) > -1 {
			known = append(known, ph)
		}
	}

	sort.Sort(byLength(known))
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
	} //else {
	//	return phs, unk
	//}

	return phs, unk
}

func consume(srtd *[]string, trans string) (string, string, bool) {
	var resPref string
	var resSuffix string

	var prefixFound bool

	notInTrans := make(map[string]bool)
	for _, ph := range *srtd {

		ind := strings.Index(trans, ph)
		if ind == 0 { // bingo
			resPref = ph
			resSuffix = trans[len(ph):]
			prefixFound = true

			break
		}

		if ind < 0 {
			notInTrans[ph] = true
		}
	}

	// Discard phonemes we know are not in trans
	if len(notInTrans) > 0 {
		var tmp []string

		for _, ph := range *srtd {
			if !notInTrans[ph] {
				tmp = append(tmp, ph)
			}
		}

		*srtd = tmp
	}
	// no known phoneme prefix, separate first rune
	if prefixFound {
		return resPref, resSuffix, prefixFound
	} else {

		t := []rune(trans)
		return string(t[0]), string(t[1:]), false
	}
}

// splitTrans applies splitIntoPhonemes to the transcription strings of a lex.Entry
func splitTrans(e *lex.Entry, symbols []string) {
	var newTs []lex.Transcription
	for _, t := range e.Transcriptions {
		t2, u2 := splitIntoPhonemes(symbols, t.Strn)
		newT := strings.Join(t2, " ")
		if len(u2) > 0 {
			fmt.Fprintf(os.Stderr, "%s > %v --> %v\n", t.Strn, t2, u2)
		}
		newTs = append(newTs, lex.Transcription{ID: t.ID, Strn: newT, EntryID: t.EntryID, Language: t.Language, Sources: t.Sources})
	}

	e.Transcriptions = newTs
}
