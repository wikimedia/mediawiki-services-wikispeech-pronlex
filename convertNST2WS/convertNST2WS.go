package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/vrules"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "<INPUT NST LEX FILE> <SYMBOL SET FILE>")
		fmt.Fprintln(os.Stderr, "\tsample invokation:  go run convertNST2WS.go swe030224NST.pron_utf8.txt sv_nst2ws-sampa_maptable.csv")
		return
	}

	// Lexicon file
	nstFileName := os.Args[1]

	ssFileName := os.Args[2]

	reFrom, err := regexp2.Compile("[.][^.]+$", regexp2.None)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Grandmaster Fail and the Furious File: %v\n", err)
	}
	ssName, err := reFrom.Replace(ssFileName, "", 0, -1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Grandmaster Fail and the Furious File: %v\n", err)
	}
	ssMapper, err := symbolset.LoadMapper(ssName, ssFileName, "NST-XSAMPA", "WS-SAMPA")
	ssRuleTo := vrules.SymbolSetRule{ssMapper.To}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Grandmaster Fail and the Furious File: %v\n", err)
	}

	nstFile, err := os.Open(nstFileName)
	defer nstFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Grandmaster Fail and the Furious File: %v\n", err)
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

		e, err := nstFmt.ParseToEntry(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert line %v to entry : %v\n", n, err)
			fmt.Fprintf(os.Stderr, "failing line: %v\n", line)
		}

		e.EntryStatus.Name = "imported"
		e.EntryStatus.Source = "nst"

		err = ssMapper.MapTranscriptions(&e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert entry to string : %v\n", err)
		} else {

			for _, r := range ssRuleTo.Validate(e) {
				panic(r) // shouldn't happen
			}

			res, err := wsFmt.Entry2String(e)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to convert entry to string : %v\n", err)
			} else {
				fmt.Printf("%v\n", res)
			}
		}
	}
	// }
	//_ = nstFile

	_ = lex.Entry{}
	//_ = line.NST{}
}
