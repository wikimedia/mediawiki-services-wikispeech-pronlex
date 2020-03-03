package rules

import (
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/symbolset"
)

func NewRequiredTransRe(ss symbolset.SymbolSet, name string, level string, re string, msg string, accept []lex.Entry, reject []lex.Entry) (validation.Rule, error) {
	rex, err := ProcessTransRe(ss, re)
	if err != nil {
		return RequiredTransRe{}, err
	}
	return RequiredTransRe{
		NameStr:  name,
		LevelStr: level,
		Message:  msg,
		Re:       rex,
		Accept:   accept,
		Reject:   reject,
	}, nil
}

func NewIllegalTransRe(ss symbolset.SymbolSet, name string, level string, re string, msg string, accept []lex.Entry, reject []lex.Entry) (validation.Rule, error) {
	rex, err := ProcessTransRe(ss, re)
	if err != nil {
		return IllegalTransRe{}, err
	}
	return IllegalTransRe{
		NameStr:  name,
		LevelStr: level,
		Message:  msg,
		Re:       rex,
		Accept:   accept,
		Reject:   reject,
	}, nil
}

func NewRequiredOrthRe(ss symbolset.SymbolSet, name string, level string, re string, msg string, accept []lex.Entry, reject []lex.Entry) (validation.Rule, error) {
	rex, err := ProcessRe(re)
	if err != nil {
		return RequiredOrthRe{}, err
	}
	return RequiredOrthRe{
		NameStr:  name,
		LevelStr: level,
		Message:  msg,
		Re:       rex,
		Accept:   accept,
		Reject:   reject,
	}, nil
}

func NewIllegalOrthRe(ss symbolset.SymbolSet, name string, level string, re string, msg string, accept []lex.Entry, reject []lex.Entry) (validation.Rule, error) {
	rex, err := ProcessRe(re)
	if err != nil {
		return IllegalOrthRe{}, err
	}
	return IllegalOrthRe{
		NameStr:  name,
		LevelStr: level,
		Message:  msg,
		Re:       rex,
		Accept:   accept,
		Reject:   reject,
	}, nil
}
