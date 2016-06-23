package validation

import (
	"fmt"

	"github.com/stts-se/pronlex/lex"
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

// Rule interface. To create a validation.Rule, make a struct implementing Validate(lex.Entry) []Result
type Rule interface {
	Validate(lex.Entry) []Result
}

// Validator is a struct containing a slice of rules
type Validator struct {
	Rules []Rule
}

// ValidateEntry is used to validate single entries. Any validation
// errors are added to the entry's EntryValidations field. The
// function returns true if the entry is valid (i.e., no validation
// issues are found), otherwise false.
func (v Validator) ValidateEntry(e *lex.Entry) bool {
	for _, rule := range v.Rules {
		for _, res := range rule.Validate(*e) {
			var ev = lex.EntryValidation{
				RuleName: res.RuleName,
				Level:    res.Level,
				Message:  res.Message,
			}
			e.EntryValidations = append(e.EntryValidations, ev)
		}
	}
	return len(e.EntryValidations) == 0
}

// Validate is used to validate a slice of entries.  Any validation
// errors are added to each entry's EntryValidations field. The
// function returns true if the entry is valid (i.e., no validation
// issues are found), otherwise false.
func (v Validator) Validate(entries []*lex.Entry) bool {
	var result = true
	for _, e := range entries {
		var ok = v.ValidateEntry(e)
		if !ok {
			result = false
		}
	}
	return result
}
