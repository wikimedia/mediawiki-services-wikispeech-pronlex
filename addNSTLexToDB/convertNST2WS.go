package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "<INPUT NST LEX FILE> <SYMBOL SET FILE>")
		fmt.Fprintln(os.Stderr, "\tsample invokation:  go run convertNST2WS.go swe030224NST.pron_utf8.txt sv_nst-xsampa_maptable.csv")
		return
	}

	// Lexicon file
	nstFileName := os.Args[1]

	// Symbol set file
	ssFile, err := os.Open(os.Args[2])
	defer ssFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Foul file fail: %v\n", err)
		return
	}
	sss := bufio.NewScanner(ssFile)

	var symbols []string
	// For now, just chop out second field of symbol set file
	for sss.Scan() {
		if err := sss.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Scandalous scanning scam: %v", err)
			return
		}
		l := sss.Text()
		if strings.HasPrefix(l, "#") || strings.HasPrefix(l, "DESCRIPTION") {
			fmt.Fprintf(os.Stderr, "Skipping:\t%s\n", l)
			continue
		}
		fs := strings.Split(l, "\t")
		if len(fs) != 4 {
			fmt.Fprintf(os.Stderr, "Symbol set file trauma: wrong number of fields in line '%s'. Bailing out.\n", l)
			return
		}
		// TODO Hardwired filed no.
		sym := fs[1]
		symbols = append(symbols, sym)

	}

	fmt.Fprintf(os.Stderr, "Symbols for splitting transcriptions: %v\n", symbols)

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

		// TODO SplitIntoPhonemes on Transcription, based on symbol set for the specific language

		e, err := nstFmt.ParseToEntry(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to conver line %v to entry : %v\n", n, err)
			fmt.Fprintf(os.Stderr, "failing line: %v\n", line)
		}
		symbolset.SplitTrans(&e, symbols)
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
