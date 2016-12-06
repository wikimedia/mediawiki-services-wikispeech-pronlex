package symbolset2

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// symbolSlice is used for sorting slices of symbols according to symbol length. Symbols with equal length will be sorted alphabetically.
type symbolSlice []Symbol

func (a symbolSlice) Len() int      { return len(a) }
func (a symbolSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a symbolSlice) Less(i, j int) bool {
	if len(a[i].String) != len(a[j].String) {
		return len(a[i].String) > len(a[j].String)
	}
	return a[i].String < a[j].String
}

var SymbolSetSuffix = ".tab"

// func LoadSymbolSetsFromDir(dirName string) (map[string]SymbolSet, error) {
// 	// list files in symbol set dir
// 	fileInfos, err := ioutil.ReadDir(dirName)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed reading symbol set dir : %v", err)
// 	}
// 	var fErrs error
// 	var symSets []SymbolSet
// 	for _, fi := range fileInfos {
// 		if strings.HasSuffix(fi.Name(), SymbolSetSuffix) {
// 			symset, err := LoadSymbolSet(filepath.Join(dirName, fi.Name()))
// 			if err != nil {
// 				if fErrs != nil {
// 					fErrs = fmt.Errorf("%v : %v", fErrs, err)
// 				} else {
// 					fErrs = err
// 				}
// 			} else {
// 				symSets = append(symSets, symset)
// 			}
// 		}
// 	}

// 	if fErrs != nil {
// 		return nil, fmt.Errorf("failed to load symbol set : %v", fErrs)
// 	}

// 	var symbolSetsMap = make(map[string]SymbolSet)
// 	for _, z := range symSets {
// 		// TODO check that x.Name doesn't already exist
// 		symbolSetsMap[z.Name] = z
// 	}
// 	return symbolSetsMap, nil
// }

// FilterSymbolsByCat is used to filter out specific symbol types from the symbol set (syllabic, non syllabic, etc)
func FilterSymbolsByCat(symbols []Symbol, types []SymbolCat) []Symbol {
	var res = make([]Symbol, 0)
	for _, s := range symbols {
		if containsCat(types, s.Cat) {
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
	regexp.MustCompile(s)
	re, err := regexp.Compile(s)
	if err != nil {
		err = fmt.Errorf("couldn't compile regexp from string '%s' : %v", s, err)
		return nil, err
	}
	return re, nil
}

func containsCat(types []SymbolCat, t SymbolCat) bool {
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
