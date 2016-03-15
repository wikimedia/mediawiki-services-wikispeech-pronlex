package line

import (
	"fmt"
	"testing"
)

var fsExpField = "For field %v, expected: '%v' got: '%v'"
var fsExp = "Expected: '%v' got: '%v'"

func checkResult(t *testing.T, expect map[Field]string, result map[Field]string) {
	if len(expect) != len(result) {
		t.Errorf(fsExp, expect, result)
	} else {
		for f, ex := range expect {
			re := result[f]
			if re != ex {
				t.Errorf(fsExpField, f.String(), ex, re)
			}
		}
	}
}

func Test_Parse_01(t *testing.T) {
	var fs = map[Field]int{
		Orth:       0,
		Pos:        1,
		Morph:      2,
		Decomp:     3,
		WordLang:   6,
		Trans1:     11,
		Translang1: 14,
		Lemma:      15,
	}
	fmt := Format{"\t", fs, 16}

	input := "hannas	PM	GEN	hannas	-	-	swe	-	-	-	-	\"\" h a . n a s	-	-	swe	hanna_01"

	var expect = map[Field]string{
		Orth:       "hannas",
		Pos:        "PM",
		Morph:      "GEN",
		Decomp:     "hannas",
		WordLang:   "swe",
		Trans1:     "\"\" h a . n a s",
		Translang1: "swe",
		Lemma:      "hanna_01",
	}
	var result, err = fmt.Parse(input)
	if err != nil {
		t.Errorf("didn't expect error here : %v", err)
	} else {
		checkResult(t, expect, result)
	}

}

func Test_Parse_02(t *testing.T) {
	var fs = map[Field]int{
		Orth:       0,
		Pos:        1,
		Morph:      2,
		Decomp:     3,
		WordLang:   6,
		Trans1:     11,
		Translang1: 14,
		Lemma:      15,
	}
	fmt := Format{"\t", fs, 16}

	input := "hannas	PM	GEN	hannas	-	-	swe	-	-	-	-	\"\" h a . n a s	-	-	swe	hanna_01	-	-"

	var _, err = fmt.Parse(input)
	if err == nil {
		t.Errorf("Expected error here")
	}

}

func Test_String_01(t *testing.T) {
	var fs = map[Field]int{
		Orth:       0,
		Pos:        1,
		Morph:      2,
		Decomp:     3,
		WordLang:   6,
		Trans1:     11,
		Translang1: 14,
		Lemma:      15,
	}
	fmt := Format{";", fs, 16}

	expect := "hannas;PM;GEN;hannas;;;eng;;;;;\"\" h a . n a s;;;swe;hanna_01"

	var input = map[Field]string{
		Orth:       "hannas",
		Pos:        "PM",
		Morph:      "GEN",
		Decomp:     "hannas",
		WordLang:   "eng",
		Trans1:     "\"\" h a . n a s",
		Translang1: "swe",
		Lemma:      "hanna_01",
	}
	var result, err = fmt.String(input)
	if err != nil {
		t.Errorf("didn't expect error here : %v", err)
	} else if result != expect {
		t.Errorf(fsExp, expect, result)
	}

}

func Test_FieldName(t *testing.T) {
	var result, expect string

	result = Orth.String()
	expect = "Orth"
	if result != expect {
		t.Errorf(fsExp, "Orth", result)
	}
	fmt.Println(result)

}
