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
		log.Fatalf("How is this misery possible? : %v", err)
	}

	fmt.Printf("%s\n", s)

}

func main() {

	cmuFileName := os.Args[1]

	cmuFile, err := os.Open(cmuFileName)
	defer cmuFile.Close()
	if err != nil {

		log.Fatal("Auch! : %v", err)
	}

	var variant = regexp.MustCompile("\\([0-9]\\)")
	lastEntry := lex.Entry{}

	s := bufio.NewScanner(cmuFile)
	for s.Scan() {
		if err = s.Err(); err != nil {
			log.Fatal("Funk! : %v", err)
		}
		l := s.Text()
		fs := strings.Split(l, "  ")
		if len(fs) != 2 || strings.HasPrefix(l, ";;") {
			fmt.Fprintf(os.Stderr, "skipping line: %v\n", l)
			continue
		} //else {
		w := strings.ToLower(fs[0])
		t := fs[1]

		// Variant transcription, not a new entry
		if variant.MatchString(w) {
			//fmt.Println(l)
			t0 := lex.Transcription{Strn: t}
			lastEntry.Transcriptions = append(lastEntry.Transcriptions, t0)
		} else {
			//fmt.Println("EN APA")
			if lastEntry.Strn != "" {
				out(lastEntry)
				//fmt.Printf("%v\n", lastEntry)
			}
			ts := []lex.Transcription{lex.Transcription{Strn: t}}
			// TODO Hard-wired language name
			lastEntry = lex.Entry{Strn: w, Transcriptions: ts, Language: "EN"}
		}
		// }
	}
	// Flush
	if lastEntry.Strn != "" {
		//fmt.Printf("%v\n", lastEntry)
		out(lastEntry)
	}

}
