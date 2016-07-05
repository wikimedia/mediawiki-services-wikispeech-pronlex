package main

import (
	"bufio"
	"fmt"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"log"
	"os"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "<INPUT NST LEX FILE>")
		return
	}

	nstFileName := os.Args[1]
	//outFile := os.Args[2]

	nstFile, err := os.Open(nstFileName)
	defer nstFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open file '%v' : %v\n", nstFileName, err)
	}

	nstFmt, err := line.NewNST()
	if err != nil {
		log.Fatal(err)
	}
	wsFmt, err := line.NewWS()
	if err != nil {
		log.Fatal(err)
	}

	nst := bufio.NewScanner(nstFile)
	n := 0
	for nst.Scan() {
		n++
		if err := nst.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "failed reading line %v : %v\n", n, err)
		}
		line := nst.Text()

		// TODO SplitIntoPhonemes on Transcription, based on symbol set for the specific language

		e, err := nstFmt.ParseToEntry(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to conver line %v to entry : %v\n", n, err)
			fmt.Fprintf(os.Stderr, "failing line: %v\n", line)
		}

		res, err := wsFmt.Entry2String(e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert entry to string : %v\n", err)
		} else {
			fmt.Printf("%v\n", res)
		}
	}
	//_ = nstFile

	_ = lex.Entry{}
	//_ = line.NST{}
}
