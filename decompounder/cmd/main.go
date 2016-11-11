package main

import (
	"bufio"
	"fmt"
	//"io"
	"log"
	"os"
	"strings"

	"github.com/stts-se/pronlex/decompounder"
)

func p() { fmt.Println() }

func xtractCompField(s string) string {
	var res string
	fs := strings.Split(s, "\t")

	if len(fs) >= 4 {
		res = fs[3]
	}

	return res
}

func decompParts(s string) []string {
	var res []string

	// For now, nuke compounding 's', etc
	s = strings.Replace(s, "+s+", "s+", -1)
	s = strings.Replace(s, "+n+", "n+", -1)
	s = strings.Replace(s, "+e+", "e+", -1)
	s = strings.Replace(s, "+a+", "a+", -1)
	s = strings.Replace(s, "+o+", "o+", -1)
	s = strings.Replace(s, "+r+", "r+", -1)
	s = strings.Replace(s, "+u+", "u+", -1)

	for _, s := range strings.Split(s, "+") {
		res = append(res, strings.ToLower(s))
	}

	return res
}

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "go run main.go <LEXFILE (decomps in 4th field)> <LOWER CASE word word word...>\n")
		return
	}

	fn := os.Args[1]
	fh, err := os.Open(fn)
	if err != nil {
		log.Fatalf("failed tyo open file: %v", err)
		return
	}

	decomp := decompounder.NewDecompounder()
	n, m := 0, 0
	s := bufio.NewScanner(fh)

	fmt.Fprint(os.Stderr, "loading compound parts...")

	for s.Scan() {

		l := s.Text()
		if err := s.Err(); err != nil {
			log.Fatalf("freaked out reading file: %v", err)
			return
		}
		d := xtractCompField(l)

		ps := decompParts(d)
		if len(ps) == 2 { // two part compounds
			decomp.Prefixes.Add(ps[0])
			decomp.Suffixes.Add(ps[1])
			n++
		}
		// fullwords
		if len(ps) == 1 {
			decomp.Suffixes.Add(ps[0])
			n++
		}
		m++
	}

	fmt.Fprint(os.Stderr, " done\n")

	fmt.Fprintf(os.Stderr, "loaded %d compounds of %d input lines\n", n, m)

	for _, w := range os.Args[2:] {
		decomps := decomp.Decomp(w)
		//fmt.Printf("%s", w)
		for _, d := range decomps {
			fmt.Printf("%s\n", strings.Join(d, "+"))
		}
		//fmt.Println()
	}
}
