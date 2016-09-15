// Package config is used for lexicon configuration settings
package config

import (
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/validation"
)

// WORK IN PROGRESS, NOT READY FOR PROPER USE

// JSON representation of a configuration used with a lexicon (string values only)
type JSON struct {
	Lexicons    []string `json:"lexicons"`
	StatusNames []string `json:"statusNames"`
	PosTags     []string `json:"posTags"`
	Languages   []string `json:"languages"`
	Users       []string `json:"users"`
	Sources     []string `json:"string"`
	Encoding    string   `json:"encoding"`
	OrthCharsRe string   `json:"orthCharsRe"`
	Validator   string   `json:"validator"`
	LineParsers []string `json:"lineParsers"`
	TTSMapper   string   `json:"ttsMapper"`
}

// Config configuration used with a lexicon
type Config struct {
	Lexicons    []string             // lexicon name(s) in some db
	StatusNames []string             // should be in some (the lex?) db
	PosTags     []string             // should be in some (the lex?) db
	Languages   []string             // should be in some (the lex?) db
	Users       []string             // should be in some (the lex?) db
	Sources     []string             // should be in some (the lex?) db
	Encoding    string               // should be in some (the lex?) db
	OrthCharsRe string               // should be in some (the lex?) db
	Validator   validation.Validator // should be in some (the lex?) db?
	LineParsers []line.Parser        // should be in some (the lex?) db?
	TTSMapper   string               //??
}

// New creates a new Config instance from a JSON representation
// func New(json JSON) (Config, error) {
// 	v, err := validatorForName(json.Validator)
// 	if err != nil {
// 		return Config{}, err
// 	}
// 	parsers, err := lineParsersForNames(json.LineParsers)
// 	if err != nil {
// 		return Config{}, err
// 	}

// 	result := Config{
// 		Lexicons:    json.Lexicons,
// 		StatusNames: json.StatusNames,
// 		PosTags:     json.PosTags,
// 		Languages:   json.Languages,
// 		Users:       json.Users,
// 		Sources:     json.Sources,
// 		Encoding:    json.Encoding,
// 		OrthCharsRe: json.OrthCharsRe,
// 		Validator:   v,
// 		LineParsers: parsers,
// 		TTSMapper:   json.TTSMapper,
// 	}
// 	return result, nil
// }

// func validatorForName(name string) (validation.Validator, error) {
// 	validators, err := validators()
// 	if err != nil {
// 		return validation.Validator{}, err
// 	}
// 	if v, ok := validators[name]; ok {
// 		return v, nil
// 	}
// 	return validation.Validator{}, fmt.Errorf("No validator for name: %s", name)
// }

// func lineParsersForNames(names []string) ([]line.Parser, error) {
// 	var result = make([]line.Parser, 0)
// 	for _, name := range names {
// 		if lineFmt, err := lineParser(name); err != nil {
// 			result = append(result, lineFmt)
// 		} else {
// 			return []line.Parser{}, err
// 		}
// 	}
// 	return result, nil
// }

// func validators() (map[string]validation.Validator, error) {
// 	nstDemo, err := vrules.NewSvSeNstValidator()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return map[string]validation.Validator{
// 		"SvSeNstValidator": nstDemo,
// 	}, nil
//}

// func lineParser(name string) (line.Parser, error) {
// 	switch strings.ToLower(name) {
// 	case "nst":
// 		return line.NewNST()
// 	default:
// 		return nil, fmt.Errorf("no line format for name %s", name)
// 	}
// }
