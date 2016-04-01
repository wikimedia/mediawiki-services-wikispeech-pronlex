package symbolset

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// FilterSymbolsByType is used to filter out specific symbol types from the symbol set (syllabic, non syllabic, etc)
func FilterSymbolsByType(symbols []Symbol, types []SymbolType) []Symbol {
	var res = make([]Symbol, 0)
	for _, s := range symbols {
		if containsType(types, s.Type) {
			res = append(res, s)
		}
	}
	return res
}

func buildRegexp(symbols []Symbol) (*regexp.Regexp, error) {
	return buildRegexpWithGroup(symbols, false, true)
}

func buildRegexpWithGroup(symbols []Symbol, removeEmpty bool, anonGroup bool) (*regexp.Regexp, error) {
	sorted := make([]Symbol, len(symbols))
	copy(sorted, symbols)
	sort.Sort(symbolSlice(sorted))
	var acc = make([]string, 0)
	for _, s := range sorted {
		if removeEmpty {
			if len(s.String) > 0 {
				acc = append(acc, regexp.QuoteMeta(s.String))
			}
		} else {
			acc = append(acc, regexp.QuoteMeta(s.String))
		}
	}
	prefix := "(?:"
	if !anonGroup {
		prefix = "("
	}
	s := prefix + strings.Join(acc, "|") + ")"
	re, err := regexp.Compile(s)
	if err != nil {
		err = fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
		return nil, err
	}
	return re, nil
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
