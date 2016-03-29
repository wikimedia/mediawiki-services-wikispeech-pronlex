// Package validation is used to validate entries (transcriptions, language labels, pos tags, etc)
package validation

import (
	"fmt"

	"github.com/stts-se/pronlex/dbapi"
)

/*
Result is a validation result with the following fields:
	RuleName - arbitrary string
	Level - typically indicating severity (e.g. Info/Warning/Fatal/Format)
	Message - arbitrary string
*/
type Result struct {
	RuleName string
	Level    string
	Message  string
}

// String returns a simple string representation of the Result instance
func (r Result) String() string {
	return fmt.Sprintf("%s (%s): %s", r.RuleName, r.Level, r.Message)
}

// Rule interface. To create a validation.Rule, make a struct implementing Validate(dbapi.Entry) []Result
type Rule interface {
	Validate(dbapi.Entry) []Result
}

// Validator is a struct containing a slice of rules
type Validator struct {
	Rules []Rule
}

// Validate is used to validate entries using a slice of rules
func (v Validator) Validate(entries []dbapi.Entry) []Result {
	var result []Result
	for _, e := range entries {
		for _, rule := range v.Rules {
			for _, res := range rule.Validate(e) {
				result = append(result, res)
			}
		}
	}
	return result
}
