package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/vrules"
)

func main() {

	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "<INPUT NST LEX FILE> <LEX2IPA MAPPER> <IPA2SAMPA MAPPER>")
		fmt.Fprintln(os.Stderr, "\tsample invokation:  go run convertNST2WS.go swe030224NST.pron.utf8 sv-se_nst-xsampa.csv sv-se_ws-sampa.csv ")
		return
	}

	nstFileName := os.Args[1]
	ssFileName1 := os.Args[2]
	ssFileName2 := os.Args[3]

	ssMapper1, err := symbolset.LoadMapper("LEX2IPA", ssFileName1, "SAMPA", "IPA")
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mapper: %v\n", err)
	}
	ssMapper2, err := symbolset.LoadMapper("IPA2SAMPA", ssFileName2, "IPA", "SAMPA")
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mapper: %v\n", err)
	}
	ssRuleTo := vrules.SymbolSetRule{ssMapper2.To}

	nstFile, err := os.Open(nstFileName)
	defer nstFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open lexicon file: %v\n", err)
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

		// todo: multimapper call direct
		err = ssMapper1.MapTranscriptions(&e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert entry to IPA : %v\n", err)
		} else {
			//fmt.Fprintf(os.Stderr, "%v\n", e.Transcriptions)
			err = ssMapper2.MapTranscriptions(&e)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to convert entry from IPA to SAMPA : %v\n", err)
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
	}

	_ = lex.Entry{}
}
