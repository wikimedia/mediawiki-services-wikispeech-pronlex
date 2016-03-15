package line

import (
	"testing"
)

var fsExpField = "For field %v, expected: '%v' got: '%v'"
var fsExpFields = "Expected: '%v' got: '%v'"

func checkResult(t *testing.T, expect map[Field]string, result map[Field]string) {
	if len(expect) != len(result) {
		t.Errorf(fsExpFields, expect, result)
	} else {
		for f, ex := range expect {
			re := result[f]
			if re != ex {
				fn, err := FieldName(f)
				if err != nil {
					t.Errorf("didn't expect error here : %v", err)
				}
				t.Errorf(fsExpField, fn, ex, re)
			}
		}
	}
}

func Test_1(t *testing.T) {
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
	fmt := Format{"\t", fs}

	input := "hannas	PM	GEN	hannas	-	-	swe	-	-	-	-	\"\" h a . n a s	-	-	swe	hanna_01	-	-"

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
