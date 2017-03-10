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
		fs := strings.Split(l, "\t")
		decomp := strings.ToLower(fs[3])
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

	if len(s) < 4 {
		return false
	}

	return true
}
