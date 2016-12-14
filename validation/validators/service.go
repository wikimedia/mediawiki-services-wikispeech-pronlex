package validators

import (
	"fmt"

	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
)

// ValidatorService is a container for maintaining 'cached' mappers and their symbol sets. Please note that currently, ValidatorService need to be used as mutex, see lexserver/validation.go
type ValidatorService struct {
	Validators map[string]*validation.Validator
}

// ValidatorForName returns the validator with the specified symbol set name. If it's not loaded yet, an error is returned.
func (vs ValidatorService) ValidatorForName(symbolSetName string) (*validation.Validator, error) {
	if vv, ok := vs.Validators[symbolSetName]; ok {
		return vv, nil
	}
	return &validation.Validator{}, fmt.Errorf("no validator loaded for symbolset %s", symbolSetName)

}

// Load is used to load validators for the input symbol sets
func (vs ValidatorService) Load(symbolsets map[string]symbolset.SymbolSet) error {
	if ss, ok := symbolsets["sv-se_ws-sampa"]; ok {
		v, err := newSvSeNstValidator(ss)
		if err != nil {
			return fmt.Errorf("couldn't initialize symbol set : %v", err)
		}
		vs.Validators[ss.Name] = &v
	}
	if ss, ok := symbolsets["nb-no_ws-sampa"]; ok {
		v, err := newNbNoNstValidator(ss)
		if err != nil {
			return fmt.Errorf("couldn't initialize symbol set : %v", err)
		}
		vs.Validators[ss.Name] = &v
	}
	if ss, ok := symbolsets["en-us_sampa_mary"]; ok {
		v, err := newEnUsCmuValidator(ss)
		if err != nil {
			return fmt.Errorf("couldn't initialize symbol set : %v", err)
		}
		vs.Validators[ss.Name] = &v
	}
	return nil
}
