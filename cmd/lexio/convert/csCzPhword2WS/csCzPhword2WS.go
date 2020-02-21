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
	"github.com/stts-se/pronlex/symbolset/mapper"
	"github.com/stts-se/pronlex/validation/rules"
)

var posTags = map[string]string{
	"1": "noun",
	"2": "adjective",
	"3": "pronoun",
	"4": "numeral",
	"5": "verb",
	"6": "adverb",
	"7": "preposition",
	"8": "conjunction",
	"9": "particle",
	"0": "interjection",
}

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

func mapPos(pos string) (string, error) {
	if pos == "" {
		return "", fmt.Errorf("empty pos: %s", pos)
	}
	res, ok := posTags[pos]
	if ok {
		return res, nil
	}
	return "", fmt.Errorf("invalid pos: %s", pos)
}

var initialSyllabicRE = regexp.MustCompile("^([aeiouāēīōūöäE])")

func mapTranscription(t0 string, mapper mapper.Mapper) (string, error) {
	t := t0
	downcase := []string{"T", "K", "S", "Š", "F", "P", "Ť"}
	for _, ch := range downcase {
		t = strings.Replace(t, ch, strings.ToLower(ch), -1)
	}
	t = strings.Replace(t, "+", "", -1)
	t = initialSyllabicRE.ReplaceAllString(t, "?$1")
	t = "'" + t
	t, err := mapper.MapTranscription(t)
	if err != nil {
		return "", err
	}
	return t, nil
}

func main() {

	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "<INPUT PHWORD LEX FILE> <PHWORD-SAMPA SYMBOLSET> <WS-SAMPA SYMBOLSET>")
		fmt.Fprintln(os.Stderr, "\tsample invokation: svSeNST2WS swe030224NST.pron.utf8 sv-se_nst-xsampa.sym sv-se_ws-sampa.sym ")
		return
	}

	nstFileName := os.Args[1]
	ssFileName1 := os.Args[2]
	ssFileName2 := os.Args[3]

	mapper, err := mapper.LoadMapperFromFile("SAMPA", "SYMBOL", ssFileName1, ssFileName2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mappers: %v\n", err)
		return
	}

	ssRuleTo := rules.SymbolSetRule{SymbolSet: mapper.SymbolSet2}

	nstFile, err := os.Open(filepath.Clean(nstFileName))
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open lexicon file: %v\n", err)
		return
	}
	defer nstFile.Close()

	sc := bufio.NewScanner(nstFile)
	n := 0
	for sc.Scan() {
		n++
		hasError := false
		if err := sc.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "general error	failed reading line %v : %v\n", n, err)
			hasError = true
		}
		l := sc.Text()
		fs := strings.Split(l, ";")
		if len(fs) < 14 {
			//Print non-entries (prefix and license) as comments
			l = "# " + l
			fmt.Println(l)
			continue
		}

		w := strings.ToLower(fs[0])
		t0 := fs[2]
		t, err := mapTranscription(t0, mapper)
		if err != nil {
			fmt.Fprintf(os.Stderr, "transcription error	failed to map transcription symbols: %v\n", err)
			hasError = true
		}

		pos0 := fs[13]
		pos, err := mapPos(pos0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "pos error	invalid pos tag <%v>\n", pos0)
			hasError = true
		}

		e := lex.Entry{Strn: w, Transcriptions: []lex.Transcription{{Strn: t}}}

		e.EntryStatus.Name = "imported"
		e.EntryStatus.Source = "phword"
		e.Language = "cs-cz"
		e.PartOfSpeech = pos

		if !hasError {
			valres, err := ssRuleTo.Validate(e)
			if err != nil {
				panic(err) // shouldn't happen
			}
			for _, r := range valres.Messages {
				panic(r) // shouldn't happen
			}

			out(e)
		}
	}

}
