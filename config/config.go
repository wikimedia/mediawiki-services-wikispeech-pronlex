package config

import (
	"fmt"
	"strings"

	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/validation"
)

type Json struct {
	Lexicons    []string `json:"lexicons"`    // lexicon name(s) in lex db
	StatusNames []string `json:"statusNames"` // should be in lex db
	PosTags     []string `json:"posTags"`     // should be in lex db
	Languages   []string `json:"languages"`   // should be in lex db
	Encoding    string   `json:"encoding"`    // should be in lex db
	OrthCharsRe string   `json:"orthCharsRe"` // should be in lex db
	Validator   string   `json:"validator"`   // should be in lex db?
	LineParsers []string `json:"lineParsers"` // should NOT be in lex db
}

type Config struct {
	Lexicons    []string             // lexicon name(s) in lex db
	StatusNames []string             // should be in lex db
	PosTags     []string             // should be in lex db
	Languages   []string             // should be in lex db
	Encoding    string               // should be in lex db
	OrthCharsRe string               // should be in lex db
	Validator   validation.Validator // should be in lex db?
	LineParsers []line.Parser        // should NOT be in lex db
}

func New(json Json) (Config, error) {

	v, err := validatorForName(json.Validator)
	if err != nil {
		return Config{}, err
	}
	parsers, err := lineParsersForNames(json.LineParsers)
	if err != nil {
		return Config{}, err
	}

	result := Config{
		Lexicons:    json.Lexicons,
		StatusNames: json.StatusNames,
		PosTags:     json.PosTags,
		Languages:   json.Languages,
		Encoding:    json.Encoding,
		OrthCharsRe: json.OrthCharsRe,
		Validator:   v,
		LineParsers: parsers,
	}
	return result, nil
}

func validatorForName(name string) (validation.Validator, error) {
	if v, ok := validators()[name]; ok {
		return v, nil
	} else {
		return validation.Validator{}, fmt.Errorf("No validator for name: %s", name)
	}
}

func lineParsersForNames(names []string) ([]line.Parser, error) {
	var result = make([]line.Parser, 0)
	for _, name := range names {
		if lineFmt, err := lineParser(name); err != nil {
			result = append(result, lineFmt)
		} else {
			return []line.Parser{}, err
		}
	}
	return result, nil
}

func validators() map[string]validation.Validator {
	return map[string]validation.Validator{
	//"NlbDemoValidator": validator.NewNLBDemoValidator(),
	}

}
func lineParser(name string) (line.Parser, error) {
	switch strings.ToLower(name) {
	case "nlb":
		return line.NewNST()
	default:
		return line.Format{}, fmt.Errorf("no line format for name %s", name)
	}
}
