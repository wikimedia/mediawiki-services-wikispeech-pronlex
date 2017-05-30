package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"github.com/stts-se/pronlex/decompounder/svnst"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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
		fmt.Fprintf(os.Stderr, "skipping input: different len: %v vs %v", wps0, trans)
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

func saveToDB(m map[string]map[WP]int, dbFile, lexiconName string) error {

	return nil
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

func main() {

	if len(os.Args) != 2 && len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, filepath.Base(os.Args[0]), "<lexicon file> <N errors before exit>?")
		os.Exit(1)
	}

	fn := os.Args[1]
	fh, err := os.Open(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Mabel Tainter Memorial Building! : %v", err)
	}

	var s *bufio.Scanner
	if strings.HasSuffix(fn, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Streptomyces tsukubaensis! : %v", err)
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
		fs := strings.Split(l, "\t")
		morf := fs[2]
		// disregard items withs spurious morphology:
		if strings.Contains(morf, "|||") || strings.HasPrefix(morf, " |") {
			continue
		}

		//orth := fs[0]
		pos := strings.TrimSpace(fs[1])
		morph := strings.TrimSpace(fs[2])
		decomp0 := fs[3]
		decomp := cleanUpDecomp(decomp0)
		firstTrans := fs[8]
		if strings.Contains(decomp, "+") {
			//fmt.Printf("%s %s %s\n", orth, decomp, firstTrans)

			rez, err := svnst.HerusticSvNSTTransDecomp(decomp, firstTrans)
			if err != nil {
				fmt.Fprintf(os.Stderr, "FAIL: %v : %s\t%s\t%#v\n\n", err, decomp, firstTrans, rez)
				fails++
				if fails >= exitAfter {
					os.Exit(1)
				}
				continue
			}

			addWordParts(decomp, rez, pos, morph)

			//fmt.Printf("%s\t%s\n", decomp, strings.Join(rez, "	<+>	"))
		}
	}

	dump(suffixLex)

	//fmt.Printf("%v\n", suffixLex)
	fmt.Fprintf(os.Stderr, "lines failed to split: %d\n", fails)
}
