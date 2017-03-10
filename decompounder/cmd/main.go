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
		if len(decomps) == 0 {
			fmt.Printf("%s\t\n", w)
		}
		for _, d := range decomps {
			fmt.Printf("%s\t%s\n", w, strings.Join(d, "+"))
		}
	}
}
