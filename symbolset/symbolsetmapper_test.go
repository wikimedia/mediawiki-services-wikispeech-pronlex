package symbolset

import "testing"

var fsExpTrans = "Expected: /%v/ got: /%v/"

func testMapTranscription1(t *testing.T, ssm Mapper, input string, expect string) {
	result, err := ssm.MapTranscription(input)
	//fmt.Println("ssm_test DEBUG", ssm.Name, input, result)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here; input=%s, expect=%s : %v", input, expect, err)
		return
	} else if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func testMapTranscriptionY(t *testing.T, ms Mappers, input string, expect string) {
	result, err := ms.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here; input=%s, expect=%s : %v", input, expect, err)
		return
	} else if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func testMapTranscriptionX(t *testing.T, ssms []Mapper, input string, expect string) {
	result := input
	for _, m := range ssms {
		r, err := m.MapTranscription(result)
		result = r
		//fmt.Println("ssm_test DEBUG", m.Name, input, result)
		if err != nil {
			t.Errorf("MapTranscription() didn't expect error here : %v", err)
			return
		}
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_NewMapper_WithCorrectInput1(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"P", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	_, err := NewMapper("test", fromName, toName, symbols)
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
	_, err := NewMapper("test", fromName, toName, symbols)
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
	_, err := NewMapper("test", fromName, toName, symbols)
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
	_, err := NewMapper("test", fromName, toName, symbols)
	if err == nil {
		t.Errorf("NewMapper() expected error here")
	}
}

func Test_NewMapper_FailWithAmbiguousPhonemes(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"rt", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
	}
	_, err := NewMapper("test", fromName, toName, symbols)
	if err == nil {
		t.Errorf("NewMapper() expected error here")
	}
}

func Test_MapTranscription_WithAmbiguousSymbols(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"rt", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "a rt r t"
	expect := "A RT R T"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_WithNonEmptyDelimiters(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"s", NonSyllabic, ""}, Symbol{"S", NonSyllabic, ""}},
		SymbolPair{Symbol{"t_s", NonSyllabic, ""}, Symbol{"T_S", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "a t s t_s"
	expect := "A T S T_S"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_EmptyDelimiterInInput1(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"r*t", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "ar*ttr"
	expect := "A RT T R"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_EmptyDelimiterInInput2(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"r*t", NonSyllabic, ""}, Symbol{"RT", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "ar*ttrrt"
	expect := "A RT T R R T"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_EmptyDelimiterInOutput(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"R", NonSyllabic, ""}},
		SymbolPair{Symbol{"t", NonSyllabic, ""}, Symbol{"T", NonSyllabic, ""}},
		SymbolPair{Symbol{"rt", NonSyllabic, ""}, Symbol{"R*T", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
	}
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "a rt r t"
	expect := "AR*TRT"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_Sampa2Ipa_Simple(t *testing.T) {
	fromName := "sampa"
	toName := "ipa"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"p", NonSyllabic, ""}},
		SymbolPair{Symbol{"", PhonemeDelimiter, ""}, Symbol{"", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"$", SyllableDelimiter, ""}, Symbol{".", SyllableDelimiter, ""}},
	}
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "pa$pa"
	expect := "pa.pa"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_Sampa2Ipa_WithSwedishStress_1(t *testing.T) {
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
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"pa$pa"
	expect := "\u02C8pa\u0300.pa"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_Sampa2Ipa_WithSwedishStress_2(t *testing.T) {
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
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"brA:$ka"
	expect := "\u02C8brɑː\u0300.ka"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_FailWithUnknownSymbols_EmptyDelim(t *testing.T) {
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
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\"baŋ.ka"
	result, err := ssm.MapTranscription(input)
	if err == nil {
		t.Errorf("NewMapper() expected error here, but got %s", result)
	}
}

func Test_MapTranscription_FailWithUnknownSymbols_NonEmptyDelim(t *testing.T) {
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
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\" b a ŋ . k a"
	result, err := ssm.MapTranscription(input)
	if err == nil {
		t.Errorf("NewMapper() expected error here, but got %s", result)
	}
}

func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_1(t *testing.T) {
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
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\u02C8ba\u0300ŋ.ka" // => ˈ`baŋ.ka before mapping
	expect := "\"\"baN$ka"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_2(t *testing.T) {
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
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\u02C8a\u0300ŋ.ka"
	expect := "\"\"aN$ka"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_3(t *testing.T) {
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
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\u02C8bra\u0300ŋ.ka"
	expect := "\"\"braN$ka"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_NstXSAMPA_To_WsSAMPA_1(t *testing.T) {
	fromName := "NST-XSAMPA"
	toName := "WS-SAMPA"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", Syllabic, ""}, Symbol{"a", Syllabic, ""}},
		SymbolPair{Symbol{"b", NonSyllabic, ""}, Symbol{"b", NonSyllabic, ""}},
		SymbolPair{Symbol{"r", NonSyllabic, ""}, Symbol{"r", NonSyllabic, ""}},
		SymbolPair{Symbol{"k", NonSyllabic, ""}, Symbol{"k", NonSyllabic, ""}},
		SymbolPair{Symbol{"N", NonSyllabic, ""}, Symbol{"N", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"$", SyllableDelimiter, ""}, Symbol{".", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\"", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\"\"", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\" b r a N $ k a"
	expect := "\"\" b r a N . k a"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_NstXSAMPA_To_WsSAMPA_2(t *testing.T) {
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
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
		SymbolPair{Symbol{"$", SyllableDelimiter, ""}, Symbol{".", SyllableDelimiter, ""}},
		SymbolPair{Symbol{"\"", Stress, ""}, Symbol{"\"", Stress, ""}},
		SymbolPair{Symbol{"\"\"", Stress, ""}, Symbol{"\"\"", Stress, ""}},
	}
	ssm, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	input := "\"\" b r a $ rs a r s"
	expect := "\"\" b r a . rs a r s"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_LoadMapper_NST2IPA(t *testing.T) {
	name := "NST-XSAMPA"
	fromColumn := "SAMPA"
	toColumn := "IPA"
	fName := "static/sv-se_nst-xsampa.csv"
	ssm, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	testMapTranscription1(t, ssm, "\"bOt`", "\u02C8bɔʈ")
	testMapTranscription1(t, ssm, "\"ku0rds", "\u02C8kɵrds")
}

func Test_LoadMapper_WS2IPA(t *testing.T) {
	name := "WS-SAMPA"
	fromColumn := "SAMPA"
	toColumn := "IPA"
	fName := "static/sv-se_ws-sampa.csv"
	ssm, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	testMapTranscription1(t, ssm, "\" b O rt", "\u02C8bɔʈ")
	testMapTranscription1(t, ssm, "\" k u0 r d s", "\u02C8kɵrds")
}

func Test_LoadMapper_IPA2WS(t *testing.T) {
	name := "WS-SAMPA"
	fromColumn := "IPA"
	toColumn := "SAMPA"
	fName := "static/sv-se_ws-sampa.csv"
	ssm, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	testMapTranscription1(t, ssm, "\u02C8bɔʈ", "\" b O rt")
	testMapTranscription1(t, ssm, "\u02C8kɵrds", "\" k u0 r d s")
}

func Test_LoadMapper_NST2WS(t *testing.T) {
	name := "NST-XSAMPA"
	fromColumn := "SAMPA"
	toColumn := "IPA"
	fName := "static/sv-se_nst-xsampa.csv"
	ssmNST, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}

	name = "WS-SAMPA"
	fromColumn = "IPA"
	toColumn = "SAMPA"
	fName = "static/sv-se_ws-sampa.csv"
	ssmWS, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}

	mappers := []Mapper{ssmNST, ssmWS}

	testMapTranscriptionX(t, mappers, "\"bOt`", "\" b O rt")
	testMapTranscriptionX(t, mappers, "\"ku0rd", "\" k u0 r d")
}

func Test_NewMapper_FailIfInputContainsDuplicates(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"A", NonSyllabic, ""}, Symbol{"a", NonSyllabic, ""}},
		SymbolPair{Symbol{"A", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"P", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	_, err := NewMapper("test", fromName, toName, symbols)
	if err == nil {
		t.Errorf("NewMapper() expected error when input contains duplicates")
	}
}

func Test_NewMapper_DontFailIfInputContainsDuplicates(t *testing.T) {
	fromName := "ssLC"
	toName := "ssUC"
	symbols := []SymbolPair{
		SymbolPair{Symbol{"a", NonSyllabic, ""}, Symbol{"A", NonSyllabic, ""}},
		SymbolPair{Symbol{"A", Syllabic, ""}, Symbol{"A", Syllabic, ""}},
		SymbolPair{Symbol{"p", NonSyllabic, ""}, Symbol{"P", NonSyllabic, ""}},
		SymbolPair{Symbol{" ", PhonemeDelimiter, ""}, Symbol{" ", PhonemeDelimiter, ""}},
	}
	_, err := NewMapper("test", fromName, toName, symbols)
	if err != nil {
		t.Errorf("NewMapper() didn't expect error when output phoneme set contains duplicates")
	}
}

func Test_LoadMapper_CMU2IPA(t *testing.T) {
	name := "CMU"
	fromColumn := "CMU"
	toColumn := "IPA"
	fName := "static/en-us_cmu.csv"
	ssm, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}

	testMapTranscription1(t, ssm, "AX $ B AW1 T", "ə.ˈba⁀ʊt")
}
func Test_LoadMapper_MARY2IPA(t *testing.T) {
	name := "MARY2IPA"
	fromColumn := "SAMPA"
	toColumn := "IPA"
	fName := "static/en-us_sampa_mary.csv"
	ssm, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}

	testMapTranscription1(t, ssm, "@ - \" b aU t", "ə.ˈba⁀ʊt")
}

func Test_LoadMapper_IPA2MARY(t *testing.T) {
	name := "IPA2MARY"
	fromColumn := "IPA"
	toColumn := "SAMPA"
	fName := "static/en-us_sampa_mary.csv"
	ssm, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	testMapTranscription1(t, ssm, "ə.ˈba⁀ʊt", "@ - \" b aU t")
}

func Test_LoadMapper_CMU2MARY(t *testing.T) {
	name := "CMU2IPA"
	fromColumn := "CMU"
	toColumn := "IPA"
	fName := "static/en-us_cmu.csv"
	ssmCMU, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	name = "IPA2MARY"
	fromColumn = "IPA"
	toColumn = "SAMPA"
	fName = "static/en-us_sampa_mary.csv"
	ssmMARY, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	mappers := []Mapper{ssmCMU, ssmMARY}

	testMapTranscriptionX(t, mappers, "AX $ B AW1 T", "@ - \" b aU t")
}

func Test_LoadMapper_SAMPA2MARY(t *testing.T) {
	name := "SAMPA2IPA"
	fromColumn := "SAMPA"
	toColumn := "IPA"
	fName := "static/sv-se_ws-sampa.csv"
	ssm1, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	name = "IPA2MARY"
	fromColumn = "IPA"
	toColumn = "SAMPA"
	fName = "static/sv-se_sampa_mary.csv"
	ssm2, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}
	mappers := []Mapper{ssm1, ssm2}
	testMapTranscriptionX(t, mappers, "eu . r \" u: p a", "E*U - r ' u: p a")
}

func Test_LoadMapper_NST2MARY(t *testing.T) {
	name := "NST2IPA"
	fromColumn := "SAMPA"
	toColumn := "IPA"
	fName := "static/sv-se_nst-xsampa.csv"
	ssm1, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	name = "IPA2MARY"
	fromColumn = "IPA"
	toColumn = "SAMPA"
	fName = "static/sv-se_sampa_mary.csv"
	ssm2, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}
	mappers := []Mapper{ssm1, ssm2}
	testMapTranscriptionX(t, mappers, "E*U$r\"u:t`a", "E*U - r ' u: rt a")
}

func Test_LoadMapper_IPA2SAMPA(t *testing.T) {
	name := "IPA2SAMPA"
	fromColumn := "IPA"
	toColumn := "SAMPA"
	fName := "static/sv-se_ws-sampa.csv"
	ssm, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	testMapTranscription1(t, ssm, "ˈkaj.rʊ", "\" k a j . r U")
	testMapTranscription1(t, ssm, "be.ˈliːn", "b e . \" l i: n")
}

func Test_LoadMapper_NST2SAMPA(t *testing.T) {
	name := "NST2IPA"
	fromColumn := "SAMPA"
	toColumn := "IPA"
	fName := "static/sv-se_nst-xsampa.csv"
	ssm1, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	name = "IPA2SAMPA"
	fromColumn = "IPA"
	toColumn = "SAMPA"
	fName = "static/sv-se_ws-sampa.csv"
	ssm2, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	mappers := []Mapper{ssm1, ssm2}
	testMapTranscriptionX(t, mappers, "\"kaj$rU", "\" k a j . r U")
	testMapTranscriptionX(t, mappers, "E*U$r\"u:t`a", "eu . r \" u: rt a")
}

func Test_LoadMapper_IPA2CMU(t *testing.T) {
	name := "IPA2CMU"
	fromColumn := "IPA"
	toColumn := "CMU"
	fName := "static/en-us_cmu.csv"
	ssm, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	testMapTranscription1(t, ssm, "ə.ˈba⁀ʊt", "AX $ B AW1 T")
	testMapTranscription1(t, ssm, "ʌ.ˈba⁀ʊt", "AH $ B AW1 T")
}

func Test_LoadMapper_MARY2CMU(t *testing.T) {
	name := "MARY2IPA"
	fromColumn := "SAMPA"
	toColumn := "IPA"
	fName := "static/en-us_sampa_mary.csv"
	ssmMARY, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	name = "IPA2CMU"
	fromColumn = "IPA"
	toColumn = "CMU"
	fName = "static/en-us_cmu.csv"
	ssmCMU, err := LoadMapper(name, fName, fromColumn, toColumn)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}

	mappers := []Mapper{ssmMARY, ssmCMU}

	testMapTranscriptionX(t, mappers, "@ - \" b aU t", "AX $ B AW1 T")
	testMapTranscriptionX(t, mappers, "V - \" b aU t", "AH $ B AW1 T")
}

func Test_LoadMappers_MARY2CMU(t *testing.T) {
	mappers, err := LoadMappers("SAMPA", "CMU", "static/en-us_sampa_mary.csv", "static/en-us_cmu.csv")
	if err != nil {
		t.Errorf("Test_LoadMappers() didn't expect error here : %v", err)
		return
	}

	testMapTranscriptionY(t, mappers, "@ - \" b aU t", "AX $ B AW1 T")
	testMapTranscriptionY(t, mappers, "V - \" b aU t", "AH $ B AW1 T")
}
