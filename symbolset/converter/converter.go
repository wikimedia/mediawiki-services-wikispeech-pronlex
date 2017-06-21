package converter

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/symbolset"
)

type Converter struct {
	Name  string
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
	invalid, err := c.getInvalidSymbols(res, c.To)
	if err != nil {
		return "", err
	}
	if len(invalid) > 0 {
		return res, fmt.Errorf("Invalid symbol(s) in output transcription /%s/: %v", res, invalid)
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

func (c Converter) getInvalidSymbols(trans string, symbolset symbolset.SymbolSet) ([]string, error) {
	if trans == symbolset.PhonemeDelimiter.String {
		return []string{}, nil
	}
	invalid := []string{}
	splitted, err := symbolset.SplitTranscription(trans)
	if err != nil {
		return invalid, err
	}
	for _, phn := range splitted {
		if !symbolset.ValidSymbol(phn) {
			invalid = append(invalid, phn)
		}
	}
	return invalid, nil
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
