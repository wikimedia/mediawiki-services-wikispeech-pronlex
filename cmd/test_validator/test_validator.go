package main

import (
	"fmt"
	"os"

	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation/validators"
)

func main() {
	usage := `USAGE:
  test_validator <SYMBOLSET FILE> <VALIDATOR FILE>`

	if len(os.Args) != 3 {
		fmt.Println(usage)
		os.Exit(1)
	}

	ssFile := os.Args[1]
	vFile := os.Args[2]

	ss, err := symbolset.LoadSymbolSet(ssFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	v, err := validators.LoadValidatorFromFile(ss, vFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Loaded %d rules containing %d tests from file %s\n", len(v.Rules), v.NumberOfTests(), vFile)
	tr, err := v.RunTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if tr.Size() > 0 {
		errs := tr.AllErrors()
		msg := fmt.Sprintf("%d tests failed for validator %s", len(errs), v.Name)
		fmt.Fprint(os.Stderr, msg)
		for _, e := range tr.AllErrors() {
			fmt.Fprintf(os.Stderr, "%v\n", e)
		}
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "No errors found, all %d tests succeeded.\n", v.NumberOfTests())

}
