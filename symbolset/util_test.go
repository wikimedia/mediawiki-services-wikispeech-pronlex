package symbolset

import "testing"

var fsExp = "Expected: '%v' got: '%v'"

func testEqStrings(t *testing.T, expect []string, result []string) {
	if len(expect) != len(result) {
		t.Errorf(fsExp, expect, result)
		return
	}
	for i, ex := range expect {
		re := result[i]
		if ex != re {
			t.Errorf(fsExp, expect, result)
			return
		}
	}
}

func testEqSymbols(t *testing.T, expect []Symbol, result []Symbol) {
	if len(expect) != len(result) {
		t.Errorf(fsExp, expect, result)
		return
	}
	for i, ex := range expect {
		re := result[i]
		if ex != re {
			t.Errorf(fsExp, expect, result)
			return
		}
	}
}

func Test_buildRegexp1(t *testing.T) {
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"t", NonSyllabic, ""},
		Symbol{".", SyllableDelimiter, ""},
		Symbol{"s", NonSyllabic, ""},
		Symbol{"t_s", NonSyllabic, ""},
		Symbol{"", PhonemeDelimiter, ""},
		Symbol{" ", PhonemeDelimiter, ""},
		Symbol{"e", Syllabic, ""},
	}
	re, err := buildRegexp(symbols)
	if err != nil {
		t.Errorf("buildRegexp() didn't expect error here : %v", err)
	}
	expect := "(?:t_s| |\\.|a|e|s|t|)"
	if expect != re.String() {
		t.Errorf(fsExp, expect, re.String())
	}
}

func Test_buildRegexp2(t *testing.T) {
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"t", NonSyllabic, ""},
		Symbol{".", SyllableDelimiter, ""},
		Symbol{"s", NonSyllabic, ""},
		Symbol{"t_s", NonSyllabic, ""},
		Symbol{"", PhonemeDelimiter, ""},
		Symbol{" ", PhonemeDelimiter, ""},
		Symbol{"e", Syllabic, ""},
	}
	re, err := buildRegexpWithGroup(symbols, true, false)
	if err != nil {
		t.Errorf("buildRegexp() didn't expect error here : %v", err)
	}
	expect := "(t_s| |\\.|a|e|s|t)"
	if expect != re.String() {
		t.Errorf(fsExp, expect, re.String())
	}
}

func Test_buildRegexp3(t *testing.T) {
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"t", NonSyllabic, ""},
		Symbol{".", SyllableDelimiter, ""},
		Symbol{"s", NonSyllabic, ""},
		Symbol{"t_s", NonSyllabic, ""},
		Symbol{"$", PhonemeDelimiter, ""},
		Symbol{" ", PhonemeDelimiter, ""},
		Symbol{"", PhonemeDelimiter, ""},
		Symbol{"e", Syllabic, ""},
	}
	re, err := buildRegexpWithGroup(symbols, false, false)
	if err != nil {
		t.Errorf("buildRegexp() didn't expect error here : %v", err)
	}
	expect := "(t_s| |\\$|\\.|a|e|s|t|)"
	if expect != re.String() {
		t.Errorf(fsExp, expect, re.String())
	}
}

func Test_FilterSymbolsByCat(t *testing.T) {
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"t", NonSyllabic, ""},
		Symbol{"%", Stress, ""},
		Symbol{".", SyllableDelimiter, ""},
		Symbol{"s", NonSyllabic, ""},
		Symbol{"t_s", NonSyllabic, ""},
		Symbol{"$", PhonemeDelimiter, ""},
		Symbol{"\"", Stress, ""},
		Symbol{"e", Syllabic, ""},
		Symbol{"-", ExplicitPhonemeDelimiter, ""},
		Symbol{"", PhonemeDelimiter, ""},
		Symbol{"+", MorphemeDelimiter, ""},
	}
	stressE := []Symbol{
		Symbol{"%", Stress, ""},
		Symbol{"\"", Stress, ""},
	}
	stressR := FilterSymbolsByCat(symbols, []SymbolCat{Stress})
	testEqSymbols(t, stressE, stressR)

	delimE := []Symbol{
		Symbol{".", SyllableDelimiter, ""},
		Symbol{"$", PhonemeDelimiter, ""},
		Symbol{"-", ExplicitPhonemeDelimiter, ""},
		Symbol{"", PhonemeDelimiter, ""},
		Symbol{"+", MorphemeDelimiter, ""},
	}
	delimR := FilterSymbolsByCat(symbols, []SymbolCat{SyllableDelimiter, PhonemeDelimiter, ExplicitPhonemeDelimiter, MorphemeDelimiter})
	testEqSymbols(t, delimE, delimR)
}

func Test_contains(t *testing.T) {
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"t", NonSyllabic, ""},
		Symbol{".", SyllableDelimiter, ""},
		Symbol{"s", NonSyllabic, ""},
		Symbol{"t_s", NonSyllabic, ""},
		Symbol{"$", PhonemeDelimiter, ""},
		Symbol{" ", PhonemeDelimiter, ""},
		Symbol{"", PhonemeDelimiter, ""},
		Symbol{"e", Syllabic, ""},
	}
	var s string

	s = "t_s"
	if !contains(symbols, s) {
		t.Errorf("contains() Expected true for symbol %s in %v", s, symbols)
	}
	s = "_"
	if contains(symbols, s) {
		t.Errorf("contains() Expected false for symbol %s in %v", s, symbols)
	}
}
