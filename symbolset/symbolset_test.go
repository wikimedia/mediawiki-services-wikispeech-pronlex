package symbolset

import (
	"testing"
)

var fsExp = "Expected: '%v' got: '%v'"
var fsDidntExp = "Didn't expect: '%v'"

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

func Test_NewSymbolSet_WithoutPhonemeDelimiter(t *testing.T) {
	name := "ss"
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"t", NonSyllabic, ""},
	}
	_, err := NewSymbolSet(name, symbols)
	if err == nil {
		t.Errorf("NewSymbolSet() should fail if no phoneme delimiter is defined")
	}
}

func Test_SplitTranscription_Normal1(t *testing.T) {
	name := "ss"
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"t", NonSyllabic, ""},
		Symbol{"s", NonSyllabic, ""},
		Symbol{"t_s", NonSyllabic, ""},
		Symbol{" ", PhonemeDelimiter, "phn delim"},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here : %v", err)
	}

	input := "a t s t_s s"
	expect := []string{"a", "t", "s", "t_s", "s"}
	result, err := ss.SplitTranscription(input)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here")
	}
	testEqStrings(t, expect, result)
}

func Test_SplitTranscription_EmptyPhonemeDelmiter1(t *testing.T) {
	name := "ss"
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"t", NonSyllabic, ""},
		Symbol{"s", NonSyllabic, ""},
		Symbol{"t_s", NonSyllabic, ""},
		Symbol{"", PhonemeDelimiter, ""},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here")
	}

	input := "atst_ss"
	expect := []string{"a", "t", "s", "t_s", "s"}
	result, err := ss.SplitTranscription(input)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here")
	}
	testEqStrings(t, expect, result)
}

func Test_SplitTranscription_FailWithUnknownSymbols_EmptyDelim(t *testing.T) {
	name := "sampa"
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"b", NonSyllabic, ""},
		Symbol{"N", NonSyllabic, ""},
		Symbol{"", PhonemeDelimiter, ""},
		Symbol{".", SyllableDelimiter, ""},
		Symbol{"\"", Stress, ""},
		Symbol{"\"\"", Stress, ""},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"baN.ka"
	//expect := []string{"\"\"", "b", "a", "N", ".", "k", "a"}
	result, err := ss.SplitTranscription(input)
	if err == nil {
		t.Errorf("SplitTranscription() expected error here, but got %s", result)
	}
}

func Test_SplitTranscription_NoFailWithUnknownSymbols_NonEmptyDelim(t *testing.T) {
	name := "sampa"
	symbols := []Symbol{
		Symbol{"a", Syllabic, ""},
		Symbol{"b", NonSyllabic, ""},
		Symbol{"N", NonSyllabic, ""},
		Symbol{" ", PhonemeDelimiter, ""},
		Symbol{".", SyllableDelimiter, ""},
		Symbol{"\"", Stress, ""},
		Symbol{"\"\"", Stress, ""},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here : %v", err)
	}
	input := "\"\" b a N . k a"
	expect := []string{"\"\"", "b", "a", "N", ".", "k", "a"}
	result, err := ss.SplitTranscription(input)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here : %v", err)
	}
	testEqStrings(t, expect, result)
}
