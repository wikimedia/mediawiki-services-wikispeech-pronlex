package converter

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/symbolset"
)

type Converter struct {
	From  symbolset.SymbolSet
	To    symbolset.SymbolSet
	Rules []Rule
}

func (c Converter) Convert(trans string) (string, error) {
	var res = trans
	var err error
	for _, r := range c.Rules {
		res, err = r.Convert(res, c.From)
		if err != nil {
			return "", err
		}
	}
	invalid, err := c.getInvalidPhonemes(res, c.From)
	if err != nil {
		return "", err
	}
	if len(invalid) > 0 {
		return res, fmt.Errorf("Invalid symbols in output transcription /%s/: %v", trans, invalid)
	}
	return res, nil
}

type test struct {
	from string
	to   string
}

type TestResult struct {
	OK     bool
	Errors []string
}

func (c Converter) getInvalidPhonemes(trans string, symbolset symbolset.SymbolSet) ([]string, error) {
	invalid := []string{}
	splitted, err := symbolset.SplitTranscription(trans)
	if err != nil {
		return invalid, err
	}
	for _, phn := range splitted {
		if !c.To.ValidSymbol(phn) {
			invalid = append(invalid, phn)
		}
	}
	return invalid, nil
}

func (c Converter) Test(tests []test) (TestResult, error) {
	res1, err := c.testExamples(tests)
	if err != nil {
		return TestResult{}, err
	}
	res2, err := c.testInternals()
	if err != nil {
		return TestResult{}, err
	}
	if res1.OK && res2.OK {
		return TestResult{}, nil
	}
	return TestResult{OK: false, Errors: append(res1.Errors, res2.Errors...)}, nil
}

// runs pre-defined tests (defined in the input file)
func (c Converter) testExamples(tests []test) (TestResult, error) {
	errors := []string{}
	for _, test := range tests {
		result, err := c.Convert(test.from)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s", err))
			//return TestResult{}, err
		}
		if result != test.to {
			msg := fmt.Sprintf("From /%s/ expected /%s/, but got /%s/", test.from, test.to, result)
			errors = append(errors, msg)
		}
		invalid, err := c.getInvalidPhonemes(result, c.To)
		if err != nil {
			return TestResult{}, err
		}
		if len(invalid) > 0 {
			errors = append(errors, fmt.Sprintf("Invalid symbols in output transcription: %v", invalid))
		}
	}
	ok := (len(errors) == 0)
	return TestResult{OK: ok, Errors: errors}, nil
}

// runs internal tests
func (c Converter) testInternals() (TestResult, error) {
	errors := []string{}
	for _, phn := range c.From.Symbols {
		// check that all input symbols can be converted without errors
		res, err := c.Convert(phn.String)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s", err))
			//return TestResult{}, err
		}
		// check that all output symbols are valid as defined in c.To
		invalid, err := c.getInvalidPhonemes(res, c.To)
		if err != nil {
			return TestResult{}, err
		}
		if len(invalid) > 0 {
			errors = append(errors, fmt.Sprintf("Invalid symbols in output transcription: %v", invalid))
		}
	}
	// for each symbol rule, check that input is defined in c.From, and output is defined in c.To
	for _, rule := range c.Rules {
		if reflect.TypeOf(rule).Name() == "SymbolRule" {
			var sr SymbolRule = rule.(SymbolRule)
			invalid, err := c.getInvalidPhonemes(sr.From, c.From)
			if err != nil {
				return TestResult{}, err
			}
			if len(invalid) > 0 {
				errors = append(errors, fmt.Sprintf("Invalid symbols in input transcription for rule %s: %v", rule, invalid))
			}
			invalid, err = c.getInvalidPhonemes(sr.To, c.To)
			if err != nil {
				return TestResult{}, err
			}
			if len(invalid) > 0 {
				errors = append(errors, fmt.Sprintf("Invalid symbols in output transcription rule %s: %v", rule, invalid))
			}
		} else if reflect.TypeOf(rule).Name() == "RegexpRule" {
			var rr RegexpRule = rule.(RegexpRule)
			invalid, err := c.getInvalidPhonemes(rr.To, c.To)
			if err != nil {
				return TestResult{}, err
			}
			if len(invalid) > 0 {
				errors = append(errors, fmt.Sprintf("Invalid symbols in output transcription for rule %s: %v", rule, invalid))
			}
		}
	}
	ok := (len(errors) == 0)
	return TestResult{OK: ok, Errors: errors}, nil
}

type Rule interface {
	Convert(trans string, symbolset symbolset.SymbolSet) (string, error)
	String() string
}

type SymbolRule struct {
	From string
	To   string
}

func (r SymbolRule) String() string {
	return fmt.Sprintf("%s\t%s\t%s", "SYMBOL", r.From, r.To)
}

func (r SymbolRule) Convert(trans string, symbolset symbolset.SymbolSet) (string, error) {
	splitted, err := symbolset.SplitTranscription(trans)
	if err != nil {
		return "", err
	}
	res := []string{}
	for _, phn := range splitted {
		if phn == r.From {
			res = append(res, r.To)
		} else {
			res = append(res, phn)
		}

	}
	return strings.Join(res, symbolset.PhonemeDelimiter.String), nil
}

type RegexpRule struct {
	From *regexp2.Regexp
	To   string
}

func (r RegexpRule) String() string {
	return fmt.Sprintf("%s\t%s\t%s", "RE", r.From, r.To)
}

func (r RegexpRule) Convert(trans string, symbolset symbolset.SymbolSet) (string, error) {
	res, err := r.From.Replace(trans, r.To, -1, -1)
	if err != nil {
		return "", err
	}
	return res, nil
}
