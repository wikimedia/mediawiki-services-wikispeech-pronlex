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

var sucTags = map[string]bool{
	"AB": true,
	"DT": true,
	"HA": true,
	"HD": true,
	"HP": true,
	"HS": true,
	"IE": true,
	"IN": true,
	"JJ": true,
	"KN": true,
	"NN": true,
	"PC": true,
	"PF": true, // ???
	"PL": true,
	"PM": true,
	"PN": true,
	"PP": true,
	"PS": true,
	"RG": true,
	"RO": true,
	"SN": true,
	"UO": true,
	"VB": true,
}

var langCodes = map[string]string{
	"swe": "sv-se",
	"sfi": "sv-fi",
	"nor": "nb-no",
	"nno": "nn-no",
	"eng": "en",
	"fin": "fi-fi",
	"ger": "de-de",
	"fre": "fr-fr",
	"rus": "ru-ru",
	"lat": "la",
	"ita": "it-it",
	"for": "foreign",
	"dan": "da-dk",
	"spa": "es-es",
}

var upperCase = regexp.MustCompile("^[A-ZÅÄÖ]+$")
var garbLine = regexp.MustCompile(".*;GARB;.*")

func removableLine(origOrth string, line string, e lex.Entry) bool {
	if upperCase.MatchString(origOrth) && e.PartOfSpeech == "RG" {
		return true
	} else if garbLine.MatchString(line) {
		return true
	}
	return false
}

func validPos(pos string) bool {
	if pos == "" {
		return true
	}
	_, ok := sucTags[pos]
	if ok {
		return true
	}
	return false
}

func mapLanguage(lang string) (string, error) {
	if lang == "" {
		return lang, nil
	}
	l, ok := langCodes[strings.ToLower(lang)]
	if ok {
		return l, nil
	}
	return lang, fmt.Errorf("couldn't map language <%v>", lang)
}

func mapTransLanguages(e *lex.Entry) error {
	var newTs []lex.Transcription
	for _, t := range e.Transcriptions {
		l, err := mapLanguage(t.Language)
		if err != nil {
			return err
		}
		t.Language = l
		newTs = append(newTs, t)
	}
	e.Transcriptions = newTs
	return nil
}

func mapTranscription(t0 string) string {
	t := t0
	long2From := regexp.MustCompile("2:( *[.]? *[%\"]* *[r])")
	long2To := "9:$1"
	short2From := regexp.MustCompile("2( *[.]? *[%\"]* *[r])")
	short2To := "9$1"
	longEFrom := regexp.MustCompile("E:( *[.]? *[%\"]* *[r])")
	longETo := "{:$1"
	shortEFrom := regexp.MustCompile("E( *[.]? *[%\"]* *[r])")
	shortETo := "{$1"
	t = long2From.ReplaceAllString(t, long2To)
	t = short2From.ReplaceAllString(t, short2To)
	t = longEFrom.ReplaceAllString(t, longETo)
	t = shortEFrom.ReplaceAllString(t, shortETo)
	if len(t0) != len(t) {
		panic(fmt.Sprintf("mapTranscription | Conversion error %s => %s", t0, t))
	}
	return t
}

func testMapTranscription(input string, expect string) string {
	result := mapTranscription(input)
	if result != expect {
		return fmt.Sprintf("input: %s\nexpected : %s\ngot      : %s", input, expect, result)
	}
	return ""
}

func testMapTranscriptions() {
	var res []string
	res = append(res, testMapTranscription(`"" b A: rn . k a . m a . % rd 2 . r e n`, `"" b A: rn . k a . m a . % rd 9 . r e n`))
	res = append(res, testMapTranscription(`"" b A: rn . k a . m a . % rn 2: . rd e n`, `"" b A: rn . k a . m a . % rn 9: . rd e n`))
	res = append(res, testMapTranscription(`"" m a s . f 2 . % rs rt 2: . r e l . s e n s`, `"" m a s . f 9 . % rs rt 9: . r e l . s e n s`))
	res = append(res, testMapTranscription(`m 2 rd . 2: r . 2 rt`, `m 9 rd . 9: r . 9 rt`))
	res = append(res, testMapTranscription(`k a . p I . "" t A: l . f 2 . % rs 2 r j . n I N . e n s`, `k a . p I . "" t A: l . f 9 . % rs 9 r j . n I N . e n s`))
	res = append(res, testMapTranscription(`k a . p I . "" t A: l . f 2: . % rs 2 r j . n I N . e n s`, `k a . p I . "" t A: l . f 9: . % rs 9 r j . n I N . e n s`))
	res = append(res, testMapTranscription(`" l E: n`, `" l E: n`))
	res = append(res, testMapTranscription(`" l E: rn`, `" l {: rn`))
	res = append(res, testMapTranscription(`"" E rt . % v E k s t`, `"" { rt . % v E k s t`))
	var errs []string
	for _, s := range res {
		if s != "" {
			errs = append(errs, s)
		}
	}
	if len(errs) > 0 {
		panic(fmt.Sprintf("%v", errs))
	}
}

func mapTranscriptions(e *lex.Entry, mapper symbolset.Mapper) error {

	err := mapper.MapTranscriptions(e)
	if err != nil {
		return err
	}
	var newTs []lex.Transcription
	for _, t := range e.Transcriptions {

		t.Strn = mapTranscription(t.Strn)
		newTs = append(newTs, t)
	}
	e.Transcriptions = newTs
	return nil
}

func main() {

	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "<INPUT NST LEX FILE> <NST-SAMPA SYMBOLSET> <WS-SAMPA SYMBOLSET>")
		fmt.Fprintln(os.Stderr, "\tsample invokation: svSeNST2WS swe030224NST.pron.utf8 sv-se_nst-xsampa.sym sv-se_ws-sampa.sym ")
		return
	}

	nstFileName := os.Args[1]
	ssFileName1 := os.Args[2]
	ssFileName2 := os.Args[3]

	mapper, err := symbolset.LoadMapperFromFile("SAMPA", "SYMBOL", ssFileName1, ssFileName2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mappers: %v\n", err)
		return
	}

	testMapTranscriptions()

	ssRuleTo := rules.SymbolSetRule{SymbolSet: mapper.SymbolSet2}

	nstFile, err := os.Open(nstFileName)
	defer nstFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open lexicon file: %v\n", err)
		return
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
		hasError := false
		if err := nst.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "general error	failed reading line %v : %v\n", n, err)
			hasError = true
		}
		line := nst.Text()

		e, origOrth, err := nstFmt.ParseToEntry(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "general error	failed to convert line %v to entry : %v\n", n, err)
			fmt.Fprintf(os.Stderr, "general error	failing line: %v\n", line)
			hasError = true
		}

		if removableLine(origOrth, line, e) {
			fmt.Fprintf(os.Stderr, "skipping line	%v\n", line)
			continue
		}

		e.EntryStatus.Name = "imported"
		e.EntryStatus.Source = "nst"
		e.Language, err = mapLanguage(e.Language)
		if err != nil {
			fmt.Fprintf(os.Stderr, "entry language error	%v\n", err)
			hasError = true
		}
		err = mapTransLanguages(&e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "trans language error	%v\n", err)
			hasError = true
		}
		if !validPos(e.PartOfSpeech) {
			fmt.Fprintf(os.Stderr, "pos error	invalid pos tag <%v>\n", e.PartOfSpeech)
			hasError = true
		}

		err = mapTranscriptions(&e, mapper)
		if err != nil {
			fmt.Fprintf(os.Stderr, "transcription error	failed to map transcription symbols : %v\n", err)
			hasError = true
		}

		if !hasError {
			valres, err := ssRuleTo.Validate(e)
			if err != nil {
				panic(err) // shouldn't happen
			}
			for _, r := range valres.Messages {
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

	_ = lex.Entry{}
}
