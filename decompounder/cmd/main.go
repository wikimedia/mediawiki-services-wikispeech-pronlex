package main

import (
	//"bufio"
	"fmt"
	//"io"
	//"compress/gzip"
	//"encoding/json"
	//"log"
	"os"
	"strings"

	"github.com/stts-se/pronlex/decompounder"
)

// func p() { fmt.Println() }

// func xtractCompField(s string) string {
// 	var res string
// 	fs := strings.Split(s, "\t")

// 	if len(fs) >= 4 {
// 		res = fs[3]
// 	}

// 	return res
// }

// func decompParts(s string) []string {
// 	var res []string

// 	for _, s := range strings.Split(s, "+") {
// 		res = append(res, strings.ToLower(s))
// 	}

// 	return res
// }

func main() {

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "go run main.go <WORD PARTS FILE> <LOWER CASE word word word...>\n")
		return
	}

	fn := os.Args[1]

	decomp, err := decompounder.NewDecompounderFromFile(fn)
	if err != nil {
		fmt.Printf("failed loading file '%s' : %v", fn, err)
		os.Exit(0)
	}

	for _, w := range os.Args[2:] {
		decomps := decomp.Decomp(w)
		//fmt.Printf("%s\t", w)
		for _, d := range decomps {
			fmt.Printf("%s\n", strings.Join(d, "+"))
		}
		//fmt.Println()
	}

}
