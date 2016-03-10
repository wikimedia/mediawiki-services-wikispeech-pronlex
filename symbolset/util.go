package symbolset

import (
	"fmt"
	"regexp"
	"strings"
)

// start: general util functions

func filterSymbolsByType(symbols []Symbol, types []SymbolType) []Symbol {
	res := make([]Symbol, 0)
	for _, s := range symbols {
		if containsType(types, s.Type) {
			res = append(res, s)
		}
	}
	return res
}

func buildRegexp(symbols []Symbol) (*regexp.Regexp, error) {
	res := make([]string, 0)
	for _, s := range symbols {
		res = append(res, regexp.QuoteMeta(s.String))
	}
	s := "(?:" + strings.Join(res, "|") + ")"
	re, err := regexp.Compile(s)
	if err != nil {
		err = fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
		return nil, err
	} else {
		return re, nil
	}
}

func containsType(types []SymbolType, t SymbolType) bool {
	for _, t0 := range types {
		if t0 == t {
			return true
		}
	}
	return false
}

func contains(symbols []Symbol, symbol string) bool {
	for _, s := range symbols {
		if s.String == symbol {
			return true
		}
	}
	return false
}

func indexOf(elements []string, element string) int {
	for i, s := range elements {
		if s == element {
			return i
		}
	}
	return -1
}

// end: general util general util functions
