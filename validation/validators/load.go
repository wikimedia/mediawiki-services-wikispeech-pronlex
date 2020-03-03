package validators

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
	rs "github.com/stts-se/pronlex/validation/rules"
	"github.com/stts-se/symbolset"
)

// RequiredTransRe	primary_stress	Fatal	Primary stress required	"

var typeIndex = 0
var nameIndex = 1
var rLevelIndex = 2
var rMessageIndex = 3
var rReIndex = 4

func find(rules []validation.Rule, rName string) (validation.Rule, int, bool) {
	for i, r := range rules {
		if r.Name() == rName {
			return r, i, true
		}
	}
	return nil, 0, false
}

func buildRegexpRule(ss symbolset.SymbolSet, rType string, rName string, fs []string, acc []lex.Entry, rej []lex.Entry) (validation.Rule, error) {
	nilRes := rs.RequiredTransRe{NameStr: rName}
	if len(fs) < 5 {
		return nilRes, fmt.Errorf("invalid line input for rule: %s", strings.Join(fs, "\t"))
	}
	level := fs[rLevelIndex]
	re := fs[rReIndex]
	msg := fs[rMessageIndex]

	switch rType {
	case "RequiredTransRe":
		{
			r, err := rs.NewRequiredTransRe(ss, rName, level, re, msg, acc, rej)
			if err != nil {
				return nilRes, err
			}
			return r, nil
		}
	case "IllegalTransRe":
		{
			r, err := rs.NewIllegalTransRe(ss, rName, level, re, msg, acc, rej)
			if err != nil {
				return nilRes, err
			}
			return r, nil
		}
	case "RequiredOrthRe":
		{
			r, err := rs.NewRequiredOrthRe(ss, rName, level, re, msg, acc, rej)
			if err != nil {
				return nilRes, err
			}
			return r, nil
		}
	case "IllegalOrthRe":
		{
			r, err := rs.NewIllegalOrthRe(ss, rName, level, re, msg, acc, rej)
			if err != nil {
				return nilRes, err
			}
			return r, nil
		}
	}
	return nilRes, fmt.Errorf("invalid rule type %s for input: %s", rType, strings.Join(fs, "\t"))

}

//ACCEPT	RequiredTransRe	hEst	\" h E s t

func parseEntry(testType string, rName string, fs []string) (string, lex.Entry, error) {
	if len(fs) < 3 {
		return "", lex.Entry{}, fmt.Errorf("invalid line input for %s test: %s", testType, strings.Join(fs, "\t"))
	}
	orth := fs[2]

	e := lex.Entry{Strn: orth}
	for _, ts := range strings.Split(fs[3], "#") {
		ts = strings.TrimSpace(ts)
		if len(ts) > 0 {
			t := lex.Transcription{Strn: ts}
			e.Transcriptions = append(e.Transcriptions, t)
		}
	}
	if len(fs) > 4 {
		e.Language = fs[4]
	}
	if len(fs) > 5 {
		e.PartOfSpeech = fs[5]
	}

	return rName, e, nil

}

var commentRe = regexp.MustCompile("^ *[#/].*")

func LoadValidatorFromFile(ss symbolset.SymbolSet, fName string) (validation.Validator, error) {
	nilRes := validation.Validator{}
	rules := []validation.Rule{}
	fh, err := os.Open(filepath.Clean(fName))
	if err != nil {
		return nilRes, err
	}
	defer fh.Close()
	s := bufio.NewScanner(fh)
	accept := make(map[string][]lex.Entry)
	reject := make(map[string][]lex.Entry)

	rLines := [][]string{}

	for s.Scan() {
		l := s.Text()
		if commentRe.MatchString(l) {
			continue
		}
		if strings.TrimSpace(l) == "" {
			continue
		}
		fs := strings.Split(l, "\t")
		lType := fs[typeIndex]
		rName := lType
		if len(fs) > 1 {
			rName = fs[nameIndex]
		}
		switch lType {
		case "ACCEPT":
			rName, entry, err := parseEntry(lType, rName, fs)
			if err != nil {
				return nilRes, err
			}
			if _, ok := accept[rName]; !ok {
				accept[rName] = []lex.Entry{}
			}
			accept[rName] = append(accept[rName], entry)
		case "REJECT":
			rName, entry, err := parseEntry(lType, rName, fs)
			if err != nil {
				return nilRes, err
			}
			if _, ok := reject[rName]; !ok {
				reject[rName] = []lex.Entry{}
			}
			reject[rName] = append(reject[rName], entry)

		default:
			rLines = append(rLines, fs)
		}
	}

	for _, fs := range rLines {
		lType := fs[typeIndex]
		rName := lType
		if len(fs) > 1 {
			rName = fs[nameIndex]
		}
		acc := accept[rName]
		rej := reject[rName]
		switch lType {
		case "MustHaveTrans":
			rule := rs.MustHaveTrans{Accept: acc, Reject: rej}
			rules = append(rules, rule)
		case "NoEmptyTrans":
			rule := rs.NoEmptyTrans{Accept: acc, Reject: rej}
			rules = append(rules, rule)
		// case "Decomp2Orth" | // TODO: set compDelim + acceptEmpty
		default:
			rule, err := buildRegexpRule(ss, lType, rName, fs, acc, rej)
			if err != nil {
				return nilRes, err
			}
			rules = append(rules, rule)
		}
	}

	rules = append(rules, rs.SymbolSetRule{SymbolSet: ss})

	inputNTests := 0
	for _, entries := range accept {
		inputNTests = inputNTests + len(entries)
	}
	for _, entries := range reject {
		inputNTests = inputNTests + len(entries)
	}

	rNames := make(map[string]bool)
	for _, r := range rules {
		if _, ok := rNames[r.Name()]; ok {
			return nilRes, fmt.Errorf("duplicate rules named %s", r.Name())
		}
		rNames[r.Name()] = true
	}

	for rName := range accept {
		_, _, ok := find(rules, rName)
		if !ok {
			return nilRes, fmt.Errorf("no rule named %s is defined (found in accept example)", rName)
		}
	}
	for rName := range reject {
		_, _, ok := find(rules, rName)
		if !ok {
			return nilRes, fmt.Errorf("no rule named %s is defined (found in reject example)", rName)
		}
	}
	v := validation.Validator{Name: ss.Name, Rules: rules}

	outputNTests := v.NumberOfTests()

	if inputNTests != outputNTests {
		a, r := v.AllTests()

		for _, e := range a {
			fmt.Printf("ACCEPT %v\n", e)
		}
		for _, e := range r {
			fmt.Printf("REJECT %v\n", e)
		}
		return nilRes, fmt.Errorf("file contains %d test lines, validator contains %v tests", inputNTests, outputNTests)
	}

	return v, nil
}
