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
	Messages []string
}

// Strings returns a slice of simple string representations of the Result
func (r Result) Strings() []string {
	var res = make([]string, 0)
	for _, msg := range r.Messages {
		res = append(res, fmt.Sprintf("%s (%s): %s", r.RuleName, r.Level, msg))
	}
	return res
}

// Rule interface. To create a validation.Rule, make a struct implementing Validate, ShouldAccept and ShouldReject as defined in this interface.
type Rule interface {
	Validate(lex.Entry) (Result, error)
	ShouldAccept() []lex.Entry
	ShouldReject() []lex.Entry
}

// Validator is a struct containing a slice of rules
type Validator struct {
	Name  string
	Rules []Rule
}

// IsDefined is used to check if the validator is initialized (by checking that the validator has a non-empty name).
func (v Validator) IsDefined() bool {
	return v.Name != ""
}

// ValidateEntry is used to validate single entries. Any validation
// errors are added to the entry's EntryValidations field. The
// function returns true if the entry is valid (i.e., no validation
// issues are found), otherwise false.
func (v Validator) ValidateEntry(e lex.Entry) (lex.Entry, bool) {
	e.EntryValidations = make([]lex.EntryValidation, 0)
	for _, rule := range v.Rules {
		res, err := rule.Validate(e)
		if err != nil {
			var ev = lex.EntryValidation{
				RuleName: "System",
				Level:    "Error",
				Message:  fmt.Sprintf("error when validating word '%s' with rule %s : %v", e.Strn, res.RuleName, err),
			}
			e.EntryValidations = append(e.EntryValidations, ev)
		} else {
			for _, msg := range res.Messages {
				var ev = lex.EntryValidation{
					RuleName: res.RuleName,
					Level:    res.Level,
					Message:  msg,
				}
				e.EntryValidations = append(e.EntryValidations, ev)
			}
		}
	}
	return e, len(e.EntryValidations) == 0
}

// ValidateEntries is used to validate a slice of entries.  Any validation
// errors are added to each entry's EntryValidations field. The
// function returns true if the entry is valid (i.e., no validation
// issues are found), otherwise false.
func (v Validator) ValidateEntries(entries []lex.Entry) ([]lex.Entry, bool) {
	var res []lex.Entry
	var valid = true
	for _, e0 := range entries {
		var e, ok = v.ValidateEntry(e0)
		if !ok {
			valid = false
		}
		res = append(res, e)
	}
	return res, valid
}
