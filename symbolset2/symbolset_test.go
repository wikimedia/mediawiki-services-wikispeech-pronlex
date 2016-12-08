package symbolset2

import "testing"

var fsExpTrans = "Expected: /%v/ got: /%v/"

func Test_NewSymbolSet_WithoutPhonemeDelimiter(t *testing.T) {
	name := "ss"
	symbols := []Symbol{
		Symbol{"a", Syllabic, "", IPASymbol{"", ""}},
		Symbol{"t", NonSyllabic, "", IPASymbol{"", ""}},
	}
	_, err := NewSymbolSet(name, symbols)
	if err == nil {
		t.Errorf("NewSymbolSet() should fail if no phoneme delimiter is defined")
	}
}

func Test_NewSymbolSet_FailIfInputContainsDuplicates(t *testing.T) {
	name := "ss"
	symbols := []Symbol{
		Symbol{"a", Syllabic, "", IPASymbol{"", ""}},
		Symbol{"a", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"t", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{" ", PhonemeDelimiter, "phn delim", IPASymbol{"", ""}},
	}
	_, err := NewSymbolSet(name, symbols)
	if err == nil {
		t.Errorf("NewSymbolSet() expected error here")
	}
}

func Test_SplitTranscription_Normal1(t *testing.T) {
	name := "ss"
	symbols := []Symbol{
		Symbol{"a", Syllabic, "", IPASymbol{"", ""}},
		Symbol{"t", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"s", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"t_s", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{" ", PhonemeDelimiter, "phn delim", IPASymbol{"", ""}},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here : %v", err)
		return
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
		Symbol{"a", Syllabic, "", IPASymbol{"", ""}},
		Symbol{"t", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"s", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"t_s", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here")
		return
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
		Symbol{"a", Syllabic, "", IPASymbol{"", ""}},
		Symbol{"b", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"N", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		Symbol{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		Symbol{"\"", Stress, "", IPASymbol{"", ""}},
		Symbol{"\"\"", Stress, "", IPASymbol{"", ""}},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here : %v", err)
		return
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
		Symbol{"a", Syllabic, "", IPASymbol{"", ""}},
		Symbol{"b", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"N", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
		Symbol{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		Symbol{"\"", Stress, "", IPASymbol{"", ""}},
		Symbol{"\"\"", Stress, "", IPASymbol{"", ""}},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here : %v", err)
		return
	}
	input := "\"\" b a N . k a"
	expect := []string{"\"\"", "b", "a", "N", ".", "k", "a"}
	result, err := ss.SplitTranscription(input)
	if err != nil {
		t.Errorf("SplitTranscription() didn't expect error here : %v", err)
	}
	testEqStrings(t, expect, result)
}

func Test_ValidSymbol1(t *testing.T) {
	name := "sampa"
	symbols := []Symbol{
		Symbol{"a", Syllabic, "", IPASymbol{"", ""}},
		Symbol{"b", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{"N", NonSyllabic, "", IPASymbol{"", ""}},
		Symbol{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
		Symbol{".", SyllableDelimiter, "", IPASymbol{"", ""}},
		Symbol{"\"", Stress, "", IPASymbol{"", ""}},
		Symbol{"\"\"", Stress, "", IPASymbol{"", ""}},
	}
	ss, err := NewSymbolSet(name, symbols)
	if err != nil {
		t.Errorf("didn't expect error here : %v", err)
		return
	}

	var phn = ""

	phn = "a"
	if !ss.ValidSymbol(phn) {
		t.Errorf("expected phoneme %v to be valid", phn)
	}

	phn = "."
	if !ss.ValidSymbol(phn) {
		t.Errorf("expected phoneme %v to be valid", phn)
	}

	phn = "x"
	if ss.ValidSymbol(phn) {
		t.Errorf("expected phoneme %v to be invalid", phn)
	}

}

func Test_ConvertToIPA(t *testing.T) {
	symbols := []Symbol{
		Symbol{"a", Syllabic, "", IPASymbol{"a", "U+0061"}},
		Symbol{"b", NonSyllabic, "", IPASymbol{"b", "U+0062"}},
		Symbol{"r", NonSyllabic, "", IPASymbol{"r", "U+0072"}},
		Symbol{"k", NonSyllabic, "", IPASymbol{"k", "U+006B"}},
		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		Symbol{"A:", Syllabic, "", IPASymbol{"ɑː", "U+0251U+02D0"}},
		Symbol{"$", SyllableDelimiter, "", IPASymbol{".", "U+002E"}},
		Symbol{"\"", Stress, "", IPASymbol{"\u02C8", "U+02C8"}},
		Symbol{"\"\"", Stress, "", IPASymbol{"\u02C8\u0300", "U+02C8U+0300"}},
	}
	ss, err := NewSymbolSet("sampa", symbols)
	if err != nil {
		t.Errorf("NewSymbolSet() didn't expect error here : %v", err)
		return
	}

	// --
	input := "\"\"brA:$ka"
	expect := "\u02C8brɑ\u0300ː.ka"
	result, err := ss.ConvertToIPA(input)
	if err != nil {
		t.Errorf("ConvertToIPA() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}

	// --
	input = "\"brA:$ka"
	expect = "\u02C8brɑː.ka"
	result, err = ss.ConvertToIPA(input)
	if err != nil {
		t.Errorf("ConvertToIPA() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_ConvertFromIPA(t *testing.T) {
	symbols := []Symbol{
		Symbol{"a", Syllabic, "", IPASymbol{"a", "U+0061"}},
		Symbol{"b", NonSyllabic, "", IPASymbol{"b", "U+0062"}},
		Symbol{"r", NonSyllabic, "", IPASymbol{"r", "U+0072"}},
		Symbol{"k", NonSyllabic, "", IPASymbol{"k", "U+006B"}},
		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
		Symbol{"A:", Syllabic, "", IPASymbol{"ɑː", "U+0251U+02D0"}},
		Symbol{"$", SyllableDelimiter, "", IPASymbol{".", "U+002E"}},
		Symbol{"\"", Stress, "", IPASymbol{"\u02C8", "U+02C8"}},
		Symbol{"\"\"", Stress, "", IPASymbol{"\u02C8\u0300", "U+02C8U+0300"}},
	}
	ss, err := NewSymbolSet("sampa", symbols)
	if err != nil {
		t.Errorf("NewSymbolSet() didn't expect error here : %v", err)
		return
	}

	// --
	input := "\u02C8brɑ\u0300ː.ka"
	expect := "\"\"brA:$ka"
	result, err := ss.ConvertFromIPA(input)
	if err != nil {
		t.Errorf("ConvertFromIPA() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}

	// --
	input = "\u02C8brɑː.ka"
	expect = "\"brA:$ka"
	result, err = ss.ConvertFromIPA(input)
	if err != nil {
		t.Errorf("ConvertFromIPA() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}
