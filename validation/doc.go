// Package validation is used to define entry validators and rules. It contains a few general rule types (see rules.go) and you can also create new rules using the Rule interface.
//
// To create a validating rule suite, initalize the Validator struct using a slice or Rule instances.
// Use the Validator.Validate (or ValidateEntry) function to have one or more entries validated.
//
package validation
