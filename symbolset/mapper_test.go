package symbolset

import "testing"

func testMapTranscription(t *testing.T, ms Mapper, input string, expect string) {
	result, err := ms.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here; input=%s, expect=%s : %v", input, expect, err)
		return
	} else if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

func Test_MapTranscription_EmptyDelimiterInInput1(t *testing.T) {
	symbols1 := []Symbol{
		Symbol{"a", Syllabic, "", IPASymbol{"a", "U+0061"}},
		Symbol{"r", NonSyllabic, "", IPASymbol{"r", "U+0072"}},
		Symbol{"t", NonSyllabic, "", IPASymbol{"t", "U+0074"}},
		Symbol{"r*t", NonSyllabic, "", IPASymbol{"R", "U+0052"}},
		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
	}
	symbols2 := []Symbol{
		Symbol{"A", Syllabic, "", IPASymbol{"a", "U+0061"}},
		Symbol{"R", NonSyllabic, "", IPASymbol{"r", "U+0072"}},
		Symbol{"T", NonSyllabic, "", IPASymbol{"t", "U+0074"}},
		Symbol{"RT", NonSyllabic, "", IPASymbol{"R", "U+0052"}},
		Symbol{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
	}
	ss1, err := NewSymbolSet("sampa1", symbols1)
	ss2, err := NewSymbolSet("sampa2", symbols2)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
		return
	}
	ssm, err := LoadMapper(ss1, ss2)

	// --
	input := "ar*ttr"
	expect := "A RT T R"
	result, err := ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}

	// --
	input = "ar*trt"
	expect = "A RT R T"
	result, err = ssm.MapTranscription(input)
	if err != nil {
		t.Errorf("MapTranscription() didn't expect error here : %v", err)
	}
	if result != expect {
		t.Errorf(fsExpTrans, expect, result)
	}
}

// func Test_MapTranscription_EmptyDelimiterInInput2(t *testing.T) {
// 	fromName := "ssLC"
// 	toName := "ssIPA"
// 	symbols := []Symbol{
// 		Symbol{"a", Syllabic, "", IPASymbol{"A", ""}},
// 		Symbol{"r", NonSyllabic, "", IPASymbol{"R", ""}},
// 		Symbol{"t", NonSyllabic, "", IPASymbol{"T", ""}},
// 		Symbol{"r*t", NonSyllabic, "", IPASymbol{"RT", ""}},
// 		Symbol{"", PhonemeDelimiter, "", IPASymbol{" ", ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "ar*ttrrt"
// 	expect := "A RT T R R T"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_EmptyDelimiterInOutput(t *testing.T) {
// 	fromName := "ssLC"
// 	toName := "ssIPA"
// 	symbols := []Symbol{
// 		Symbol{"a", Syllabic, "", IPASymbol{"A", ""}},
// 		Symbol{"r", NonSyllabic, "", IPASymbol{"R", ""}},
// 		Symbol{"t", NonSyllabic, "", IPASymbol{"T", ""}},
// 		Symbol{"rt", NonSyllabic, "", IPASymbol{"R*T", ""}},
// 		Symbol{" ", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "a rt r t"
// 	expect := "AR*TRT"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_FailWithUnknownSymbols_EmptyDelim(t *testing.T) {
// 	fromName := "sampa1"
// 	toName := "ipa2"
// 	symbols := []Symbol{
// 		Symbol{"a", Syllabic, "", IPASymbol{"A", ""}},
// 		Symbol{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		Symbol{"ŋ", NonSyllabic, "", IPASymbol{"N", ""}},
// 		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 		Symbol{".", SyllableDelimiter, "", IPASymbol{"$", SyllableDelimiter, ""}},
// 		Symbol{"\"", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		Symbol{"\"\"", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\"\"baŋ.ka"
// 	result, err := ssm.MapTranscription(input)
// 	if err == nil {
// 		t.Errorf("NewSymbolSet() expected error here, but got %s", result)
// 	}
// }

// func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_1(t *testing.T) {
// 	fromName := "ipa"
// 	toName := "sampa"
// 	symbols := []Symbol{
// 		Symbol{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		Symbol{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		Symbol{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		Symbol{"ŋ", NonSyllabic, "", IPASymbol{"N", ""}},
// 		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 		Symbol{".", SyllableDelimiter, "", IPASymbol{"$", SyllableDelimiter, ""}},
// 		Symbol{"\u02C8", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		Symbol{"\u02C8\u0300", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\u02C8ba\u0300ŋ.ka" // => ˈ`baŋ.ka before mapping
// 	expect := "\"\"baN$ka"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_2(t *testing.T) {
// 	fromName := "ipa"
// 	toName := "sampa"
// 	symbols := []Symbol{
// 		Symbol{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		Symbol{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		Symbol{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		Symbol{"ŋ", NonSyllabic, "", IPASymbol{"N", ""}},
// 		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 		Symbol{".", SyllableDelimiter, "", IPASymbol{"$", SyllableDelimiter, ""}},
// 		Symbol{"\u02C8", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		Symbol{"\u02C8\u0300", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\u02C8a\u0300ŋ.ka"
// 	expect := "\"\"aN$ka"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_Ipa2Sampa_WithSwedishStress_3(t *testing.T) {
// 	fromName := "ipa"
// 	toName := "sampa"
// 	symbols := []Symbol{
// 		Symbol{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		Symbol{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		Symbol{"r", NonSyllabic, "", IPASymbol{"r", ""}},
// 		Symbol{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		Symbol{"ŋ", NonSyllabic, "", IPASymbol{"N", ""}},
// 		Symbol{"", PhonemeDelimiter, "", IPASymbol{"", ""}},
// 		Symbol{".", SyllableDelimiter, "", IPASymbol{"$", SyllableDelimiter, ""}},
// 		Symbol{"\u02C8", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		Symbol{"\u02C8\u0300", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\u02C8bra\u0300ŋ.ka"
// 	expect := "\"\"braN$ka"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_NstXSAMPA_To_WsSAMPA_1(t *testing.T) {
// 	fromName := "NST-XSAMPA"
// 	toName := "WS-SAMPA_IPADUMMY"
// 	symbols := []Symbol{
// 		Symbol{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		Symbol{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		Symbol{"r", NonSyllabic, "", IPASymbol{"r", ""}},
// 		Symbol{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		Symbol{"N", NonSyllabic, "", IPASymbol{"N", ""}},
// 		Symbol{" ", PhonemeDelimiter, "", IPASymbol{" ", ""}},
// 		Symbol{"$", SyllableDelimiter, "", IPASymbol{".", SyllableDelimiter, ""}},
// 		Symbol{"\"", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		Symbol{"\"\"", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\"\" b r a N $ k a"
// 	expect := "\"\" b r a N . k a"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_MapTranscription_NstXSAMPA_To_WsSAMPA_2(t *testing.T) {
// 	fromName := "NST-XSAMPA"
// 	toName := "WS-SAMPA_IPADUMMY"
// 	symbols := []Symbol{
// 		Symbol{"a", Syllabic, "", IPASymbol{"a", ""}},
// 		Symbol{"b", NonSyllabic, "", IPASymbol{"b", ""}},
// 		Symbol{"r", NonSyllabic, "", IPASymbol{"r", ""}},
// 		Symbol{"rs", NonSyllabic, "", IPASymbol{"rs", ""}},
// 		Symbol{"s", NonSyllabic, "", IPASymbol{"s", ""}},
// 		Symbol{"k", NonSyllabic, "", IPASymbol{"k", ""}},
// 		Symbol{"N", NonSyllabic, "", IPASymbol{"N", ""}},
// 		Symbol{" ", PhonemeDelimiter, "", IPASymbol{" ", ""}},
// 		Symbol{"$", SyllableDelimiter, "", IPASymbol{".", SyllableDelimiter, ""}},
// 		Symbol{"\"", Stress, "", IPASymbol{"\"", Stress, ""}},
// 		Symbol{"\"\"", Stress, "", IPASymbol{"\"\"", Stress, ""}},
// 	}
// 	ssm, err := NewSymbolSet("test", symbols)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	input := "\"\" b r a $ rs a r s"
// 	expect := "\"\" b r a . rs a r s"
// 	result, err := ssm.MapTranscription(input)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}
// 	if result != expect {
// 		t.Errorf(fsExpTrans, expect, result)
// 	}
// }

// func Test_loadSymbolSet_NST2WS(t *testing.T) {
// 	name := "NST-XSAMPA"
// 	fromColumn := "SAMPA"
// 	toColumn := "IPA"
// 	fName := "static/sv-se_nst-xsampa.tab"
// 	ssmNST, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}

// 	name = "WS-SAMPA"
// 	fromColumn = "IPA"
// 	toColumn = "SYMBOL"
// 	fName = "static/sv-se_ws-sampa.tab"
// 	ssmWS, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 	}

// 	mappers := []SymbolSet{ssmNST, ssmWS}

// 	testMapTranscriptionX(t, mappers, "\"bOt`", "\" b O rt")
// 	testMapTranscriptionX(t, mappers, "\"ku0rd", "\" k u0 r d")
// }

// func Test_NewSymbolSet_FailIfInputContainsDuplicates(t *testing.T) {
// 	fromName := "ssLC"
// 	toName := "ssUC"
// 	symbols := []Symbol{
// 		Symbol{"A", NonSyllabic, "", IPASymbol{"a", ""}},
// 		Symbol{"A", Syllabic, "", IPASymbol{"A", ""}},
// 		Symbol{"p", NonSyllabic, "", IPASymbol{"P", ""}},
// 		Symbol{" ", PhonemeDelimiter, "", IPASymbol{" ", ""}},
// 	}
// 	_, err := NewSymbolSet("test", symbols)
// 	if err == nil {
// 		t.Errorf("NewSymbolSet() expected error when input contains duplicates")
// 	}
// }
// func Test_loadSymbolSet_CMU2MARY(t *testing.T) {
// 	name := "CMU2IPA"
// 	fromColumn := "CMU"
// 	toColumn := "IPA"
// 	fName := "static/en-us_cmu.tab"
// 	ssmCMU, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2MARY"
// 	fromColumn = "IPA"
// 	toColumn = "SYMBOL"
// 	fName = "static/en-us_sampa_mary.tab"
// 	ssmMARY, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	mappers := []SymbolSet{ssmCMU, ssmMARY}

// 	testMapTranscriptionX(t, mappers, "AX $ B AW1 T", "@ - \" b aU t")
// }

// func Test_loadSymbolSet_SAMPA2MARY(t *testing.T) {
// 	name := "SAMPA2IPA"
// 	fromColumn := "SYMBOL"
// 	toColumn := "IPA"
// 	fName := "static/sv-se_ws-sampa.tab"
// 	ssm1, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2MARY"
// 	fromColumn = "IPA"
// 	toColumn = "SAMPA"
// 	fName = "static/sv-se_sampa_mary.tab"
// 	ssm2, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}
// 	mappers := []SymbolSet{ssm1, ssm2}
// 	testMapTranscriptionX(t, mappers, "eu . r \" u: p a", "E*U - r ' u: p a")
// 	testMapTranscriptionX(t, mappers, "@ s . \"\" e", "e s - \" e")
// }

// func Test_loadSymbolSet_MARY2SAMPA(t *testing.T) {
// 	name := "MARY2IPA"
// 	fromColumn := "SAMPA"
// 	toColumn := "IPA"
// 	fName := "static/sv-se_sampa_mary.tab"
// 	ssm1, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2SAMPA"
// 	fromColumn = "IPA"
// 	toColumn = "SYMBOL"
// 	fName = "static/sv-se_ws-sampa.tab"
// 	ssm2, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}
// 	mappers := []SymbolSet{ssm1, ssm2}
// 	testMapTranscriptionX(t, mappers, "E*U - r ' u: p a", "eu . r \" u: p a")
// 	testMapTranscriptionX(t, mappers, "e s - \" e", "e s . \"\" e")
// 	testMapTranscriptionX(t, mappers, "\" e: - p a", "\"\" e: . p a")
// 	testMapTranscriptionX(t, mappers, "\" A: - p a", "\"\" A: . p a")

// 	mapper, err := LoadMapperFromFile("SAMPA", "SYMBOL", "static/sv-se_sampa_mary.tab", "static/sv-se_ws-sampa.tab")
// 	if err != nil {
// 		t.Errorf("Test_LoadMapperFromFile() didn't expect error here : %v", err)
// 		return
// 	}

// 	testMapTranscriptionY(t, mapper, "\" e: - p a", "\"\" e: . p a")
// 	testMapTranscriptionY(t, mapper, "\" A: - p a", "\"\" A: . p a")
// }

// func Test_loadSymbolSet_NST2MARY(t *testing.T) {
// 	name := "NST2IPA"
// 	fromColumn := "SAMPA"
// 	toColumn := "IPA"
// 	fName := "static/sv-se_nst-xsampa.tab"
// 	ssm1, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2MARY"
// 	fromColumn = "IPA"
// 	toColumn = "SAMPA"
// 	fName = "static/sv-se_sampa_mary.tab"
// 	ssm2, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}
// 	mappers := []SymbolSet{ssm1, ssm2}
// 	testMapTranscriptionX(t, mappers, "E*U$r\"u:t`a", "E*U - r ' u: rt a")
// }

// func Test_loadSymbolSet_NST2SAMPA(t *testing.T) {
// 	name := "NST2IPA"
// 	fromColumn := "SAMPA"
// 	toColumn := "IPA"
// 	fName := "static/sv-se_nst-xsampa.tab"
// 	ssm1, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2SAMPA"
// 	fromColumn = "IPA"
// 	toColumn = "SYMBOL"
// 	fName = "static/sv-se_ws-sampa.tab"
// 	ssm2, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	mappers := []SymbolSet{ssm1, ssm2}
// 	testMapTranscriptionX(t, mappers, "\"kaj$rU", "\" k a j . r U")
// 	testMapTranscriptionX(t, mappers, "E*U$r\"u:t`a", "eu . r \" u: rt a")
// }

// func Test_loadSymbolSet_MARY2CMU(t *testing.T) {
// 	name := "MARY2IPA"
// 	fromColumn := "SYMBOL"
// 	toColumn := "IPA"
// 	fName := "static/en-us_sampa_mary.tab"
// 	ssmMARY, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	name = "IPA2CMU"
// 	fromColumn = "IPA"
// 	toColumn = "CMU"
// 	fName = "static/en-us_cmu.tab"
// 	ssmCMU, err := loadSymbolSet_(name, fName, fromColumn, toColumn)
// 	if err != nil {
// 		t.Errorf("MapTranscription() didn't expect error here : %v", err)
// 		return
// 	}

// 	mappers := []SymbolSet{ssmMARY, ssmCMU}

// 	testMapTranscriptionX(t, mappers, "@ - \" b aU t", "AX $ B AW1 T")
// 	testMapTranscriptionX(t, mappers, "V - \" b aU t", "AH $ B AW1 T")
// }

// func Test_LoadMapperFromFile_MARY2CMU(t *testing.T) {
// 	mappers, err := LoadMapperFromFile("SYMBOL", "CMU", "static/en-us_sampa_mary.tab", "static/en-us_cmu.tab")
// 	if err != nil {
// 		t.Errorf("Test_LoadMapperFromFile() didn't expect error here : %v", err)
// 		return
// 	}

// 	testMapTranscriptionY(t, mappers, "@ - \" b aU t", "AX $ B AW1 T")
// 	testMapTranscriptionY(t, mappers, "V - \" b aU t", "AH $ B AW1 T")
// }

func Test_LoadMapperFromFile_NST2WS(t *testing.T) {
	mapper, err := LoadMapperFromFile("SAMPA", "SYMBOL", "static/nb-no_nst-xsampa.tab", "static/nb-no_ws-sampa.tab")
	if err != nil {
		t.Errorf("Test_LoadMapperFromFile() didn't expect error here : %v", err)
		return
	}

	testMapTranscription(t, mapper, "\"A:$bl@s", "\" A: . b l @ s")
	testMapTranscription(t, mapper, "\"tSE$kIsk", "\" t S e . k i s k")
	testMapTranscription(t, mapper, "\"\"b9$n@r", "\"\" b 2 . n @ r")
	testMapTranscription(t, mapper, "\"b9$n@r", "\" b 2 . n @ r")
}

func Test_LoadMapperFromFile_FailIfBothHaveTheSameName(t *testing.T) {
	_, err := LoadMapperFromFile("SAMPA", "SAMPA", "static/nb-no_nst-xsampa.tab", "static/nb-no_ws-sampa.tab")
	if err == nil {
		t.Errorf("LoadMapperFromFile() expected error here")
	}
}

func Test_LoadMapperFromFile_FailIfBothHaveTheSameFile(t *testing.T) {
	_, err := LoadMapperFromFile("XSAMPA", "SAMPA", "static/nb-no_nst-xsampa.tab", "static/nb-no_nst-xsampa.tab")
	if err == nil {
		t.Errorf("LoadMapperFromFile() expected error here")
	}
}
