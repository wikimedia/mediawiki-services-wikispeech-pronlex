package main

import (
	"bufio"
	"fmt"
	//"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
)

var ws, errrr = line.NewWS()

func out(e lex.Entry) {
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

	cmuFileName := os.Args[1]

	cmuFile, err := os.Open(cmuFileName)
	defer cmuFile.Close()
	if err != nil {

		log.Fatal("couldn't open input file : %v", err)
	}

	var variant = regexp.MustCompile("\\([0-9]\\)")
	lastEntry := lex.Entry{}

	s := bufio.NewScanner(cmuFile)
	for s.Scan() {
		if err = s.Err(); err != nil {
			log.Fatal("scanner failure : %v", err)
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
		t := fs[1]

		// Variant transcription, not a new entry
		if variant.MatchString(w) {
			t0 := lex.Transcription{Strn: t}
			lastEntry.Transcriptions = append(lastEntry.Transcriptions, t0)
		} else {
			if lastEntry.Strn != "" {
				out(lastEntry)
			}
			ts := []lex.Transcription{lex.Transcription{Strn: t}}
			// Hard-wired (ISO 639-1 language name + ISO 3166-1 alpha 2 country)
			lastEntry = lex.Entry{Strn: w, Transcriptions: ts, Language: "en-us"}
		}
	}
	// Flush
	if lastEntry.Strn != "" {
		out(lastEntry)
	}

}
