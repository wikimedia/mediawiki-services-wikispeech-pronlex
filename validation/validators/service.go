package validators

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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
	return &validation.Validator{}, fmt.Errorf("no validator is defined for symbolset %s", symbolSetName)

}

// HasValidator is used to check whether a validator exists for the given symbol set
func (vs ValidatorService) HasValidator(symbolSetName string) bool {
	_, ok := vs.Validators[symbolSetName]
	return ok
}

func (vs ValidatorService) testValidator(v validation.Validator) error {
	tr, err := v.RunTests()
	if err != nil {
		return err
	}
	if tr.Size() > 0 {
		msg := fmt.Sprintf("init tests failed for validator %s", v.Name)
		log.Print(msg)
		for _, e := range tr.AllErrors() {
			log.Printf("%v", e)
		}
		return fmt.Errorf("%s, see log for details", msg)
	}
	return nil
}

// Load is used to load validators for the input symbol sets
func (vs ValidatorService) Load(symbolsets map[string]symbolset.SymbolSet, symsetDirName string) error {
	if ss, ok := symbolsets["sv-se_ws-sampa"]; ok {
		v, err := newSvSeNstValidator(ss)
		if err != nil {
			return fmt.Errorf("couldn't initialize validator from code : %v", err)
		}
		err = vs.testValidator(v)
		if err != nil {
			return fmt.Errorf("couldn't initialize validator from code : %v", err)
		}
		log.Printf("Loaded validator from code: %s", v.Name)
		vs.Validators[ss.Name] = &v
	}
	if ss, ok := symbolsets["sv-se_ws-sampa-DEMO"]; ok { // FOR DEMO DB
		v, err := newSvSeNstValidator(ss)
		if err != nil {
			return fmt.Errorf("couldn't initialize validator from code : %v", err)
		}
		err = vs.testValidator(v)
		if err != nil {
			return fmt.Errorf("couldn't initialize validator from code : %v", err)
		}
		log.Printf("Loaded validator from code: %s", v.Name)
		vs.Validators[ss.Name] = &v
	}
	if ss, ok := symbolsets["nb-no_ws-sampa"]; ok {
		v, err := newNbNoNstValidator(ss)
		if err != nil {
			return fmt.Errorf("couldn't initialize validator from code : %v", err)
		}
		err = vs.testValidator(v)
		if err != nil {
			return fmt.Errorf("couldn't initialize validator from code : %v", err)
		}
		log.Printf("Loaded validator from code: %s", v.Name)
		vs.Validators[ss.Name] = &v
	}
	if ss, ok := symbolsets["en-us_ws-sampa"]; ok {
		v, err := newEnUsCmuValidator(ss)
		if err != nil {
			return fmt.Errorf("couldn't initialize validator from code : %v", err)
		}
		err = vs.testValidator(v)
		if err != nil {
			return fmt.Errorf("couldn't initialize validator from code : %v", err)
		}
		log.Printf("Loaded validator from code: %s", v.Name)
		vs.Validators[ss.Name] = &v
	}
	for _, ss := range symbolsets {
		fName := filepath.Join(symsetDirName, ss.Name) + ".vd" // validator suffix
		if _, err := os.Stat(fName); !os.IsNotExist(err) {
			v, err := LoadValidatorFromFile(ss, fName)
			if err != nil {
				return fmt.Errorf("couldn't initialize validator from file %s : %v", fName, err)
			}
			if v0, ok := vs.Validators[ss.Name]; ok {
				// merge two validators!
				v.Rules = append(v.Rules, v0.Rules...)
			}
			err = vs.testValidator(v)
			if err != nil {
				return fmt.Errorf("couldn't initialize validator from file %s : %v", fName, err)
			}
			vs.Validators[ss.Name] = &v
			log.Printf("Loaded validator from file: %s", v.Name)
		}
	}
	return nil
}
