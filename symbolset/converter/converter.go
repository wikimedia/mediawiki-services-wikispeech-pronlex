package converter

import (
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
	return res, nil
}

func (c Converter) Test() error {
	// TODO:
	// * run pre-defined tests (from the rule file)!
	// * check that all symbols in From are covered by the rules
	// * check that all input symbols in the rules are legal symbols in From
	// * check that all generated output symbols are legal symbols in To
	return nil
}

type Rule interface {
	Convert(trans string, symbolset symbolset.SymbolSet) (string, error)
}

type SymbolRule struct {
	From string
	To   string
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
	From regexp2.Regexp
	To   string
}

func (r RegexpRule) Convert(trans string, symbolset symbolset.SymbolSet) (string, error) {
	res, err := r.From.Replace(trans, r.To, -1, -1)
	if err != nil {
		return "", err
	}
	return res, nil
}
