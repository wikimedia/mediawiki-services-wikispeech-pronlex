package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation/rules"
	"github.com/stts-se/rbg2p"
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

func syllabify(syller rbg2p.Syllabifier, phnSet rbg2p.PhonemeSet, trans string) (string, error) {
	phonemes, err := phnSet.SplitTranscription(trans)
	if err != nil {
		return "", err
	}
	sylled := syller.SyllabifyFromPhonemes(phonemes)
	return sylled, nil
}

var moveStressRe = regexp.MustCompile("([^ ]+)([012])")

func main() {

	if len(os.Args) != 5 {
		fmt.Fprintln(os.Stderr, "<INPUT CMU LEX FILE> <CMU SYMBOLSET> <WS-SAMPA SYMBOLSET> <SYLLDEF FILE>")
		fmt.Fprintln(os.Stderr, "\tsample invokation: CMU2WS cmudict-0.7b.utf8 en-us_cmu.sym en_us_sampa_mary.sym enu_cmu.syll")
		return
	}

	cmuFileName := os.Args[1]
	ssFileName1 := os.Args[2]
	ssFileName2 := os.Args[3]
	syllRuleFile := os.Args[4]

	mapper, err := symbolset.LoadMapperFromFile("CMU", "SYMBOL", ssFileName1, ssFileName2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mappers: %v\n", err)
		return
	}
	ssRuleTo := rules.SymbolSetRule{SymbolSet: mapper.SymbolSet2}

	syller, err := rbg2p.LoadSyllFile(syllRuleFile)
	if err != nil {
		log.Printf("couldn't load rule file %s : %s", syllRuleFile, err)
		os.Exit(1)
	}
	phnSet := syller.PhonemeSet

	cmuFile, err := os.Open(filepath.Clean(cmuFileName))
	if err != nil {

		log.Fatalf("couldn't open input file : %v", err)
	}
	defer cmuFile.Close()

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
		t0 = moveStressRe.ReplaceAllString(t0, "$2 $1")

		t0, err = syllabify(syller, phnSet, t0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to syllabify transcription : %v\n", err)
			continue
		}

		t, err := mapper.MapTranscription(t0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to map transcription symbols : %v\n", err)
			continue
		}
		// Variant transcription, not a new entry
		if variant.MatchString(w) {
			t0 := lex.Transcription{Strn: t}
			lastEntry.Transcriptions = append(lastEntry.Transcriptions, t0)
		} else {
			if lastEntry.Strn != "" {
				valres, err := ssRuleTo.Validate(lastEntry)
				if err != nil {
					panic(err) // shouldn't happen
				}
				for _, r := range valres.Messages {
					panic(r) // shouldn't happen
				}
				out(lastEntry)
			}
			ts := []lex.Transcription{{Strn: t}}
			// Hard-wired (ISO 639-1 language name + ISO 3166-1 alpha 2 country)
			lastEntry = lex.Entry{Strn: w, Transcriptions: ts, Language: "en-us"}
		}
	}
	// Flush
	if lastEntry.Strn != "" {
		valres, err := ssRuleTo.Validate(lastEntry)
		if err != nil {
			panic(err) // shouldn't happen
		}
		for _, r := range valres.Messages {
			panic(r) // shouldn't happen
		}

		out(lastEntry)
	}

}
