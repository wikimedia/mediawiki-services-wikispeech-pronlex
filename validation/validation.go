// Package validation is used to validate entries (transcriptions, language labels, pos tags, etc)
package validation

import "github.com/stts-se/pronlex/dbapi"

/*
Result is a validation result with an arbitrary
	Name - arbitrary string
	Level - typically indicating severity (e.g. Info/Warning/Fatal/Format)
	Message - arbitrary string
*/
type Result struct {
	Name    string
	Level   string
	Message string
}

// Rule interface. To create a validation.Rule, make a struct implementing Validate(dbapi.Entry) []Result
type Rule interface {
	Validate(dbapi.Entry) []Result
}

// Validator is a struct containing a slice of rules
type Validator struct {
	Rules []Rule
}

// Validate is used to validate an entry using a slice of rules
func (v Validator) Validate(e dbapi.Entry) []Result {
	var result []Result
	for _, rule := range v.Rules {
		for _, res := range rule.Validate(e) {
			result = append(result, res)
		}
	}
	return result
}
