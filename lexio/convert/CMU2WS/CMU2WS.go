package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation/rules"
)

var ws, errrr = line.NewWS()

func out(e lex.Entry) {
	e.EntryStatus.Name = "imported"
	e.EntryStatus.Source = "cmu"

	if errrr != nil {
		log.Fatalf("instantiation failure : %v", errrr)
	}

	s, err := ws.Entry2String(e)
	if err != nil {
		log.Fatalf("entry2string failed : %v", err)
	}

	fmt.Printf("%s\n", s)
}

func main() {

	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "<INPUT CMU LEX FILE> <CMU SYMBOLSET> <WS-SAMPA SYMBOLSET>")
		fmt.Fprintln(os.Stderr, "\tsample invokation:  go run convert.go cmudict-0.7b.utf8 en-us_cmu.tab en_us_sampa_mary.tab")
		return
	}

	cmuFileName := os.Args[1]
	ssFileName1 := os.Args[2]
	ssFileName2 := os.Args[3]

	mapper, err := symbolset.LoadMapperFromFile("CMU", "SYMBOL", ssFileName1, ssFileName2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mappers: %v\n", err)
		return
	}
	ssRuleTo := rules.SymbolSetRule{SymbolSet: mapper.SymbolSet2}

	cmuFile, err := os.Open(cmuFileName)
	defer cmuFile.Close()
	if err != nil {

		log.Fatalf("couldn't open input file : %v", err)
	}

	var variant = regexp.MustCompile("\\([0-9]\\)")
	lastEntry := lex.Entry{}

	s := bufio.NewScanner(cmuFile)
	for s.Scan() {
		if err = s.Err(); err != nil {
			log.Fatalf("scanner failure : %v", err)
		}
		l := s.Text()

		fs := strings.Split(l, "  ")
		if len(fs) != 2 || strings.HasPrefix(l, ";;") {
			//Print non-entries (prefix and license) as comments
			l = "# " + l
			fmt.Println(l)
			continue
		}
		w := strings.ToLower(fs[0])
		t0 := strings.Replace(fs[1], "AH0", "AX", -1)
		t0 = strings.Replace(t0, "0", "", -1)

		t, err := mapper.MapTranscription(t0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to map transcription symbols : %v\n", err)
		} else {

			// Variant transcription, not a new entry
			if variant.MatchString(w) {
				t0 := lex.Transcription{Strn: t}
				lastEntry.Transcriptions = append(lastEntry.Transcriptions, t0)
			} else {
				if lastEntry.Strn != "" {
					for _, r := range ssRuleTo.Validate(lastEntry) {
						panic(r) // shouldn't happen
					}
					out(lastEntry)
				}
				ts := []lex.Transcription{lex.Transcription{Strn: t}}
				// Hard-wired (ISO 639-1 language name + ISO 3166-1 alpha 2 country)
				lastEntry = lex.Entry{Strn: w, Transcriptions: ts, Language: "en-us"}
			}
		}
	}
	// Flush
	if lastEntry.Strn != "" {
		for _, r := range ssRuleTo.Validate(lastEntry) {
			panic(r) // shouldn't happen
		}

		out(lastEntry)
	}

}
