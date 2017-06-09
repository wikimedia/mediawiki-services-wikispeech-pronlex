// NB: One-off hack that is not intended to be maintained.

package main

import (
	"bufio"
	"compress/gzip"
	"fmt"

	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/stts-se/pronlex/decompounder/svnst"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line" // For formatting output file
)

// Takes a lexfile in WS format and tries to split the (first) transcription of each word into compound parts according to the word parts field, where the entry is decompounded.

func cleanUpDecomp(d string) string {
	res := strings.ToLower(d)
	res = strings.Replace(res, "!", "", -1)
	res = strings.Replace(res, "+s+", "s+", -1)
	res = strings.Replace(res, "+i+", "i+", -1) // Nikola+i+
	res = strings.Replace(res, "+o+", "o+", -1)
	res = strings.Replace(res, "+a+", "a+", -1)
	res = strings.Replace(res, "+u+", "u+", -1)
	res = strings.Replace(res, "+g+", "g+", -1) // armerin+g+
	res = strings.Replace(res, "+ar+", "ar+", -1)
	res = strings.Replace(res, "+r+", "r+", -1)
	res = strings.Replace(res, "+ra+", "ra+", -1)
	res = strings.Replace(res, "+es+", "es+", -1)
	res = strings.Replace(res, "+na+", "na+", -1)
	res = strings.Replace(res, "+on+", "on+", -1)
	res = strings.Replace(res, "+e+", "e+", -1)
	res = strings.Replace(res, "+-+", "+", -1)
	if strings.HasPrefix(res, "+-+") {
		res = strings.Replace(res, "+-+", "", 1)
	}
	if strings.HasPrefix(res, "+") {
		res = strings.Replace(res, "+", "", 1)
	}

	if strings.HasSuffix(res, "+s") {
		res = strings.Replace(res, "+s", "s", 1)
	}

	return res
}

type WP struct {
	trans string
	pos   string
	morph string
}

var prefixLex = make(map[string]map[WP]int)
var infixLex = make(map[string]map[WP]int)
var suffixLex = make(map[string]map[WP]int)

func add(strn, trans, pos, morph string, lex map[string]map[WP]int) {
	if v, ok := lex[strn]; ok {
		v[WP{trans: trans, pos: pos, morph: morph}]++
	} else {
		m := make(map[WP]int)
		wp := WP{trans: trans, pos: pos, morph: morph}
		m[wp]++
		lex[strn] = m
	}
}

func addWordParts(wps string, trans []string, pos, morph string) {
	var wps0 = strings.Split(wps, "+")

	if len(wps0) != len(trans) {
		fmt.Fprintf(os.Stderr, "skipping input: different len: %v vs %v\n", wps0, trans)
	}

	if len(wps0) < 2 {
		return
	}

	if len(wps0) > 1 {
		add(wps0[0], trans[0], "", "", prefixLex)
	}

	if len(wps0) > 2 {
		//add(wps0[0], trans[0], prefixLex)
		for i := 1; i < len(wps0)-1; i++ {
			add(wps0[i], trans[i], "", "", infixLex)
		}
	}

	last := len(wps0) - 1
	add(wps0[last], trans[last], pos, morph, suffixLex)

	//fmt.Println(wps)
}

type Freq struct {
	Word WP
	Freq int
}

func freqSort(m map[WP]int) []Freq {
	var res []Freq
	for k, v := range m {
		res = append(res, Freq{k, v})
	}

	sort.Slice(res, func(i, j int) bool { return res[i].Freq > res[j].Freq })

	return res
}

func totFreq(fs []Freq) int {
	var res int
	for _, f := range fs {
		res += f.Freq
	}
	return res
}

// Filter that returns unique transcriptions over frequency min, and discards the rest.
// Starting with the highest frequency transcription.
func minFreqAndDifferentTrans(fs []Freq, min int) []Freq {
	var res []Freq
	seenTrans := make(map[string]bool)

	sort.Slice(fs, func(i, j int) bool { return fs[i].Freq > fs[j].Freq })

	for _, f := range fs {
		trans := f.Word.trans
		if f.Freq >= min && !seenTrans[trans] {
			seenTrans[trans] = true
			res = append(res, f)
		}
	}

	return res
}

// Filter that collapses entries of the same POS and MORPH values. The
// most frequent transcription is placed first in a slice.
func collapseEntriesWithSamePOS(w string, f []Freq) []lex.Entry {
	var res []lex.Entry

	// Keeps track of already seen POS+MORPH
	seen := make(map[string]bool)

	for _, s := range f {

		posMorph := s.Word.pos + s.Word.morph

		t := lex.Transcription{Strn: s.Word.trans}
		e := lex.Entry{
			Strn:           w,
			WordParts:      w,
			PartOfSpeech:   s.Word.pos,
			Morphology:     s.Word.morph,
			Transcriptions: []lex.Transcription{t},
		}

		if ok := seen[posMorph]; !ok {
			res = append(res, e)
			seen[posMorph] = true
		} else { // Add transcription variant to already existing Entry
			for i, ex := range res {
				if ex.PartOfSpeech == s.Word.pos && ex.Morphology == s.Word.morph {
					ts := res[i].Transcriptions
					ts = append(ts, t)
					res[i].Transcriptions = ts
				}
			}
		}
	}
	return res
}

func dump(m map[string]map[WP]int) {
	for k, v := range m {
		srt := freqSort(v)
		tot := totFreq(srt)
		min := 4
		if tot > 50 {
			min = 20
		}
		fltr := minFreqAndDifferentTrans(srt, min)
		if len(fltr) == 0 && len(srt) > 0 {
			fltr = srt[0:1]
		}
		if tot > 5 {
			for _, s := range fltr {
				fmt.Printf("%d\t%s\t%v\n", tot, k, s)
			}
		}
	}
}

func toFile(m map[string]map[WP]int, fileName string) error {

	fh, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {

		return fmt.Errorf("lexfile2decompPronlex.toFile failed : %v", err)
	}
	defer fh.Close()

	fmt.Fprintf(os.Stderr, "printing compound parts to file '%s'\n", fileName)

	wsFmt, err := line.NewWS()
	if err != nil {
		return fmt.Errorf("lexfile2decompPronlen.toFile failed initializing : %v", err)
	}

	for k, v := range m {
		srt := freqSort(v) // []Freq
		tot := totFreq(srt)
		min := 4
		if tot > 50 {
			min = 20
		}
		fltr := minFreqAndDifferentTrans(srt, min)
		//fltr := collapseEntriesWithSamePOS(fltr0)
		if len(fltr) == 0 && len(srt) > 0 { // Nothing survived the filtering. Put back first element if any.
			fltr = srt[0:1]
		}
		if tot > 5 {
			// TODO: If same POS+MORPH but different trans, collapse into one lex.Entry
			entries := collapseEntriesWithSamePOS(k, fltr)
			for _, e := range entries {

				es, err := wsFmt.Entry2String(e)
				if err != nil {
					return fmt.Errorf("lexfile2decompPronlex.toFile failed formatting : %v", err)
				}

				//fmt.Printf("%d\t%s\t%s\n", tot, k, es)
				fmt.Fprintf(fh, "%s\n", es)
			}
		}
	}

	//fh.Flush()

	return nil
}

func main() {

	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, filepath.Base(os.Args[0]), "<lexicon file> <N errors before exit>?")
		os.Exit(1)
	}

	fn := os.Args[1]
	fh, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Mabel Tainter Memorial Building! : %v\n", err)
	}

	var s *bufio.Scanner
	if strings.HasSuffix(fn, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Streptomyces tsukubaensis! : %v\n", err)
		}
		s = bufio.NewScanner(gz)
	} else {
		s = bufio.NewScanner(fh)
	}

	exitAfter := 0
	if len(os.Args) == 3 {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil {
			msg := fmt.Sprintf("Second optional argument should be an integer, got '%s'", os.Args[2])
			fmt.Fprintf(os.Stderr, "%s\n", msg)
			os.Exit(1)
		} //else {
		exitAfter = i
		//}
	}

	fails := 0

	for s.Scan() {
		l := s.Text()

		wsFmt, err := line.NewWS()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to initialize line formatter : %v\n", err)
			os.Exit(0)
		}

		e, err := wsFmt.ParseToEntry(l)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse line : %v\n", err)
			os.Exit(0)
		}

		// disregard items withs spurious morphology:
		if strings.Contains(e.Morphology, "|||") || strings.HasPrefix(e.Morphology, " |") {
			continue
		}

		pos := e.PartOfSpeech
		morph := e.Morphology
		decomp0 := e.WordParts
		decomp := cleanUpDecomp(decomp0)
		firstTrans := e.Transcriptions[0].Strn // fs[8]
		if strings.Contains(decomp, "+") {
			//fmt.Printf("%s %s %s\n", orth, decomp, firstTrans)

			rez, err := svnst.HerusticSvNSTTransDecomp(decomp, firstTrans)
			if err != nil {
				fmt.Fprintf(os.Stderr, "FAIL: %v : %s\t%s\t%#v\n", err, decomp, firstTrans, rez)
				fails++
				if fails >= exitAfter {
					os.Exit(1)
				}
				continue
			}

			addWordParts(decomp, rez, pos, morph)

		}
	}

	// TODO: File base name as command line arg instead
	outFileBaseName := strings.Replace(filepath.Base(fn), ".gz", "", -1)
	err = toFile(suffixLex, outFileBaseName+"_sufflex.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAILURE: %v\n", err)
		os.Exit(0)
	}
	err = toFile(prefixLex, outFileBaseName+"_preflex.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAILURE: %v\n", err)
		os.Exit(0)
	}
	err = toFile(infixLex, outFileBaseName+"_infixlex.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAILURE: %v\n", err)
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "lines failed to split: %d\n", fails)
}
