package symbolset

import (
	"testing"
)

var fsExpTrans = "Expected: /%v/ got: /%v/"

func Test_NewMapper_WithCorrectInput1(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"P", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	_, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("NewMapper() didn't expect error here : %v", err)
	}
}

func Test_NewMapper_FailIfInputLacksPhonemeDelimiter(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"P", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", NonSyllabic, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	_, err := NewMapper(fromName, toName, symbols)
	if err == nil {
		t.Errorf("NewMapper() expected error here")
	}
}

func Test_NewMapper_FailIfOutputLacksPhonemeDelimiter(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"P", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", NonSyllabic, ""}},
	}
	_, err := NewMapper(fromName, toName, symbols)
	if err == nil {
		t.Errorf("NewMapper() expected error here")
	}
}

func Test_NewMapper_FailIfBothSymbolSetsHaveTheSameName(t *testing.T) {
	fromName := "ssLC"
	toName := "ssLC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"P", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	_, err := NewMapper(fromName, toName, symbols)
	if err == nil {
		t.Errorf("NewMapper() expected error here")
	}
}

func Test_NewMapper_FailWithAmbiguousPhonemesAndNoExplicitDelimiter(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"rt", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
	}
	_, err := NewMapper(fromName, toName, symbols)
	if err == nil {
		t.Errorf("NewMapper() expected error here")
	}
}

func Test_mapTranscription_WithAmbiguousSymbols(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"rt", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "a rt r t"
	expect := "A RT R T"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_WithNonEmptyDelimiters(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"s", NonSyllabic, ""}, Symbol{"S", NonSyllabic, ""}},
		SymbolPair{Symbol{"t_s", NonSyllabic, ""}, Symbol{"T_S", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "a t s t_s"
	expect := "A T S T_S"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_EmptyDelimiterInInput1(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"rt", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"-", ExplicitPhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "arttr"
	expect := "A RT T R"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_EmptyDelimiterInInput2(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"rt", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"-", ExplicitPhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "arttrr-t"
	expect := "A RT T R R T"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_EmptyDelimiterInOutput(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"rt", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{"-", ExplicitPhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "a rt r t"
	expect := "ARTR-T"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_Sampa2Ipa_Simple(t *testing.T) {
	fromName := "sampa"
	toName := "ipa"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"p", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"$", SyllableDelimiter, ""}, Symbol{".", SyllableDelimiter, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "pa$pa"
	expect := "pa.pa"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_Sampa2Ipa_WithSwedishStress_1(t *testing.T) {
	fromName := "sampa"
	toName := "ipa"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"p", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"$", SyllableDelimiter, ""}, Symbol{".", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\"", Stress, ""}, Symbol{"\u02C8", Stress, ""}},
		SymbolPair{Symbol{"\"\"", Stress, ""}, Symbol{"\u02C8\u0300", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"pa$pa"
	expect := "\u02C8pa\u0300.pa"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_Sampa2Ipa_WithSwedishStress_2(t *testing.T) {
	fromName := "sampa"
	toName := "ipa"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"r", NonSyllabic, ""}},
		SymbolPair{Symbol{"k", NonSyllabic, ""}, Symbol{"k", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"A:", Syllabic, ""}, Symbol{"ɑː", Syllabic, ""}},
		SymbolPair{Symbol{"$", SyllableDelimiter, ""}, Symbol{".", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\"", Stress, ""}, Symbol{"\u02C8", Stress, ""}},
		SymbolPair{Symbol{"\"\"", Stress, ""}, Symbol{"\u02C8\u0300", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"brA:$ka"
	expect := "\u02C8brɑː\u0300.ka"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_FailWithUnknownSymbols_EmptyDelim(t *testing.T) {
	fromName := "sampa1"
	toName := "sampa2 "
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"ŋ", NonSyllabic, ""}, Symbol{"N", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{".", SyllableDelimiter, ""}, Symbol{"$", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\"", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\"\"", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"baŋ.ka"
	result, err := ssm.mapTranscription(input)
	if err == nil {
		t.Errorf("NewMapper() expected error here, but got %s", result)
	}
}

func Test_mapTranscription_FailWithUnknownSymbols_NonEmptyDelim(t *testing.T) {
	fromName := "sampa1"
	toName := "sampa2 "
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"ŋ", NonSyllabic, ""}, Symbol{"N", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{".", SyllableDelimiter, ""}, Symbol{"$", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\"", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\"\"", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\" b a ŋ . k a"
	result, err := ssm.mapTranscription(input)
	if err == nil {
		t.Errorf("NewMapper() expected error here, but got %s", result)
	}
}

func Test_mapTranscription_Ipa2Sampa_WithSwedishStress_1(t *testing.T) {
	fromName := "ipa"
	toName := "sampa"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"k", NonSyllabic, ""}, Symbol{"k", NonSyllabic, ""}},
		SymbolPair{Symbol{"ŋ", NonSyllabic, ""}, Symbol{"N", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{".", SyllableDelimiter, ""}, Symbol{"$", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\u02C8", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\u02C8\u0300", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\u02C8ba\u0300ŋ.ka" // => ˈ`baŋ.ka before mapping
	expect := "\"\"baN$ka"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_Ipa2Sampa_WithSwedishStress_2(t *testing.T) {
	fromName := "ipa"
	toName := "sampa"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"k", NonSyllabic, ""}, Symbol{"k", NonSyllabic, ""}},
		SymbolPair{Symbol{"ŋ", NonSyllabic, ""}, Symbol{"N", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{".", SyllableDelimiter, ""}, Symbol{"$", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\u02C8", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\u02C8\u0300", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\u02C8a\u0300ŋ.ka"
	expect := "\"\"aN$ka"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_Ipa2Sampa_WithSwedishStress_3(t *testing.T) {
	fromName := "ipa"
	toName := "sampa"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"r", NonSyllabic, ""}},
		SymbolPair{Symbol{"k", NonSyllabic, ""}, Symbol{"k", NonSyllabic, ""}},
		SymbolPair{Symbol{"ŋ", NonSyllabic, ""}, Symbol{"N", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{".", SyllableDelimiter, ""}, Symbol{"$", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\u02C8", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\u02C8\u0300", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\u02C8bra\u0300ŋ.ka"
	expect := "\"\"braN$ka"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_NstXSAMPA_To_WsSAMPA_1(t *testing.T) {
	fromName := "NST-XSAMPA"
	toName := "WS-SAMPA"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"r", NonSyllabic, ""}},
		SymbolPair{Symbol{"k", NonSyllabic, ""}, Symbol{"k", NonSyllabic, ""}},
		SymbolPair{Symbol{"N", NonSyllabic, ""}, Symbol{"N", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"-", ExplicitPhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"$", SyllableDelimiter, ""}, Symbol{".", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\"", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\"\"", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"braN$ka"
	expect := "\"\" b r a N . k a"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_mapTranscription_NstXSAMPA_To_WsSAMPA_2(t *testing.T) {
	fromName := "NST-XSAMPA"
	toName := "WS-SAMPA"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"r", NonSyllabic, ""}},
		SymbolPair{Symbol{"rs", NonSyllabic, ""}, Symbol{"rs", NonSyllabic, ""}},
		SymbolPair{Symbol{"s", NonSyllabic, ""}, Symbol{"s", NonSyllabic, ""}},
		SymbolPair{Symbol{"k", NonSyllabic, ""}, Symbol{"k", NonSyllabic, ""}},
		SymbolPair{Symbol{"N", NonSyllabic, ""}, Symbol{"N", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"-", ExplicitPhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"$", SyllableDelimiter, ""}, Symbol{".", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\"", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\"\"", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper(fromName, toName, symbols)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"bra$rsar-s"
	expect := "\"\" b r a . rs a r s"
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func testMapTranscription(t *testing.T, ssm Mapper, input string, expect string) {
	result, err := ssm.mapTranscription(input)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_LoadMapper_NST2WS(t *testing.T) {
	fromName := "NST-XSAMPA"
	toName := "WS-SAMPA"
	fName := "static/sv_maptable.csv"
	ssm, err := LoadMapper(fName, fromName, toName)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	testMapTranscription(t, ssm, "\"bOt`", "\" b O rt")
	testMapTranscription(t, ssm, "\"ku0r-d", "\" k u0 r d")
}

func Test_LoadMapper_NST2IPA(t *testing.T) {
	fromName := "NST-XSAMPA"
	toName := "IPA"
	fName := "static/sv_maptable.csv"
	ssm, err := LoadMapper(fName, fromName, toName)
	if err != nil {
		t.Errorf("mapTranscription() didn't expect error here : %v", err)
	}
	testMapTranscription(t, ssm, "\"bOt`", "\u02C8bɔʈ")
	testMapTranscription(t, ssm, "\"ku0r-ds", "\u02C8kɵrds")
}
