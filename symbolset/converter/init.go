package converter

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/symbolset"
)

func isTest(s string) bool {
	return strings.HasPrefix(s, "TEST\t")
}

var testRe = regexp.MustCompile("^TEST\t([^\t]+)\t([^\t]+)$")

func parseTest(s string) (test, error) {
	var matchRes []string
	matchRes = testRe.FindStringSubmatch(s)
	if matchRes == nil {
		return test{}, fmt.Errorf("invalid symbol set definition: " + s)
	}
	return test{from: matchRes[1], to: matchRes[2]}, nil
}

var commentAtEndRe = regexp.MustCompile("^(.*[^/]+)//+.*$")

func trimComment(s string) string {
	return strings.TrimSpace(commentAtEndRe.ReplaceAllString(s, "$1"))
}

func isComment(s string) bool {
	return strings.HasPrefix(s, "//")
}

func isFrom(s string) bool {
	return strings.HasPrefix(s, "FROM\t")
}

func isTo(s string) bool {
	return strings.HasPrefix(s, "TO\t")
}

func isRegexpRule(s string) bool {
	return strings.HasPrefix(s, "RE\t")
}

var regexpRuleRe = regexp.MustCompile("^RE\t([^\t]+)\t([^\t]+)$")

func parseRegexpRule(s string) (Rule, error) {
	var matchRes []string
	matchRes = regexpRuleRe.FindStringSubmatch(s)
	if matchRes == nil {
		return RegexpRule{}, fmt.Errorf("invalid regexp rule definition: " + s)
	}
	from, err := regexp2.Compile(matchRes[1], regexp2.None)
	if err != nil {
		return RegexpRule{}, err
	}
	to := matchRes[2]
	return RegexpRule{From: from, To: to}, nil
}

func isSymbolRule(s string) bool {
	return strings.HasPrefix(s, "SYMBOL\t")
}

var symbolRuleRe = regexp.MustCompile("^SYMBOL\t([^\t]+)\t([^\t]+)$")

func parseSymbolRule(s string) (Rule, error) {
	var matchRes []string
	matchRes = symbolRuleRe.FindStringSubmatch(s)
	if matchRes == nil {
		return SymbolRule{}, fmt.Errorf("invalid symbol rule definition: " + s)
	}
	from := matchRes[1]
	to := matchRes[2]
	return SymbolRule{From: from, To: to}, nil
}

func isBlankLine(s string) bool {
	return len(s) == 0
}

var symbolSetRe = regexp.MustCompile("^(FROM|TO)\t([^\t]+)$")

func parseSymbolSet(s string) (string, error) {
	var matchRes []string
	matchRes = symbolSetRe.FindStringSubmatch(s)
	if matchRes == nil {
		return "", fmt.Errorf("invalid symbol set definition: " + s)
	}
	return matchRes[2], nil
}

func LoadFile(symbolSets map[string]symbolset.SymbolSet, fName string) (Converter, TestResult, error) {
	var converter = Converter{}
	var err error
	fh, err := os.Open(fName)
	defer fh.Close()
	if err != nil {
		return Converter{}, TestResult{}, err
	}
	n := 0
	s := bufio.NewScanner(fh)
	var testLines []test
	for s.Scan() {
		if err := s.Err(); err != nil {
			return Converter{}, TestResult{}, err
		}
		n++
		l := trimComment(strings.TrimSpace(s.Text()))
		if isBlankLine(l) || isComment(l) {
		} else if isFrom(l) {
			ss, err := parseSymbolSet(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			if val, ok := symbolSets[ss]; ok {
				converter.From = val
			} else {
				return Converter{}, TestResult{}, fmt.Errorf("Symbolset not defined: %s", ss)
			}
		} else if isTo(l) {
			ss, err := parseSymbolSet(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			if val, ok := symbolSets[ss]; ok {
				converter.To = val
			} else {
				return Converter{}, TestResult{}, fmt.Errorf("Symbolset not defined: %s", ss)
			}
		} else if isSymbolRule(l) {
			rule, err := parseSymbolRule(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			converter.Rules = append(converter.Rules, rule)
		} else if isRegexpRule(l) {
			rule, err := parseRegexpRule(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			converter.Rules = append(converter.Rules, rule)
		} else if isTest(l) {
			test, err := parseTest(l)
			if err != nil {
				return Converter{}, TestResult{}, err
			}
			testLines = append(testLines, test)
		}
	}
	testRes, err := converter.Test(testLines)
	if err != nil {
		return Converter{}, TestResult{}, err
	}
	return converter, testRes, nil
}

var Suffix = ".txt"

// LoadFromDir loads a all symbol sets from the specified folder (all files with .tab extension)
func LoadFromDir(symbolSets map[string]symbolset.SymbolSet, dirName string) (map[string]Converter, TestResult, error) {
	// list files in dir
	fileInfos, err := ioutil.ReadDir(dirName)
	if err != nil {
		return nil, TestResult{}, fmt.Errorf("failed reading symbol set dir : %v", err)
	}
	var fErrs error
	var testResult = TestResult{OK: true}
	var convs []Converter
	for _, fi := range fileInfos {
		if strings.HasSuffix(fi.Name(), Suffix) {
			conv, testRes, err := LoadFile(symbolSets, filepath.Join(dirName, fi.Name()))
			if err != nil {
				thisErr := fmt.Errorf("could't load converter from file %s : %v", fi.Name(), err)
				if fErrs != nil {
					fErrs = fmt.Errorf("%v : %v", fErrs, thisErr)
				} else {
					fErrs = thisErr
				}
			} else {
				if !testRes.OK {
					testResult.OK = false
				}
				testResult.Errors = append(testResult.Errors, testRes.Errors...)
				convs = append(convs, conv)
			}
		}
	}

	if fErrs != nil {
		return nil, TestResult{}, fErrs
	}

	var cMap = make(map[string]Converter)
	for _, c := range convs {
		name := c.From.Name + " to " + c.To.Name
		// TODO checks that x.Name doesn't already exist ?
		if _, ok := cMap[name]; ok {
			// do nothing
		} else {
			cMap[name] = c
		}
	}
	return cMap, testResult, nil
}
