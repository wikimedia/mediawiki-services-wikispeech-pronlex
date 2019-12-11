// TODO: Remove this file.
package main

import (
	"compress/gzip"
	"fmt"
	//"io"
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func SvNSTCharCanMatch(c rune) string {

	switch c {
	case 'a':
		return "(A:|a)"
	case 'b':
		return "b"
	default:
		return "UNKNOWN: " + string(c)
	}

}

func SvNSTSplitTrans(wordPart1 string, wordPart2 string, trans string) (string, string) {
	//fmt.Println("'" + wordPart1)
	wp1LastChar := wordPart1[len(wordPart1)-1]
	wp2FirstChar := wordPart2[0]

	wp1RE := SvNSTCharCanMatch(rune(wp1LastChar))
	wp2RE := SvNSTCharCanMatch(rune(wp2FirstChar))

	fmt.Println(wp1RE)
	fmt.Println(wp2RE)

	var trans1 string
	var trans2 string

	return trans1, trans2

}

func SvNSTTransAlign(decomps []string, trans string) []string {
	var res []string

	if len(decomps) == 0 {
		return res
	}
	if len(decomps) == 1 {
		return []string{trans}
	}

	for i := 0; i < len(decomps)-1; i++ {
		SvNSTSplitTrans(decomps[i], decomps[i+1], trans)
	}

	return res
}

// TODO Only care about first transcription variant
func SvNSTDecompTransAlign(lexiconLine string) ([]string, []string) {
	fs := strings.Split(lexiconLine, "\t")
	decomps := strings.Split(strings.Trim(fs[3], "+"), "+")
	fmt.Printf("%#v\n", decomps)
	firstTrans := fs[8]

	//var res []string
	res := SvNSTTransAlign(decomps, firstTrans)
	return decomps, res
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println(filepath.Base(os.Args[0]), "<lexicon file>")
		os.Exit(1)
	}

	fn := os.Args[1]
	fh, err := os.Open(fn)
	if err != nil {
		fmt.Printf("ketonic monosaccharide! : %v", err)
	}

	var s *bufio.Scanner
	if strings.HasSuffix(fn, ".gz") {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			fmt.Printf("high-fructose corn syrup! : %v", err)
		}
		s = bufio.NewScanner(gz)
	} else {
		s = bufio.NewScanner(fh)
	}

	n := 0

	pref := make(map[string]int)
	suf := make(map[string]int)

	for s.Scan() {
		l := s.Text()
		SvNSTDecompTransAlign(l)
		fs := strings.Split(l, "\t")
		decomp := strings.ToLower(strings.TrimSpace(fs[3]))
		wordParts := strings.Split(decomp, "+")
		if len(wordParts) < 2 {
			continue
		}

		//fmt.Printf("%v\n", wordParts)
		for i := 0; i < len(wordParts)-1; i++ {
			wp := clean(wordParts[i])
			//fmt.Printf("%v\n", wordParts[i])
			if ok(wp) {
				pref[wp]++
			}
		}

		wp := clean(wordParts[len(wordParts)-1])
		//fmt.Printf("%v > %s\n", wordParts, wp)
		//fmt.Println()
		if ok(wp) {
			suf[wp]++
		}
		n++
	}

	minFreq := 3
	for s, n := range pref {
		if n >= minFreq {
			fmt.Printf("PREFIX:%s\n", s)
		}

	}

	for s, n := range suf {
		if n >= minFreq {
			fmt.Printf("SUFFIX:%s\n", s)
		}
	}

	fmt.Fprintf(os.Stderr, "Decomp words: %d\n", n)
}

func clean(s string) string {
	s = strings.Trim(s, "!")
	return s
}

func ok(s string) bool {

	return len(s) >= 4

}
