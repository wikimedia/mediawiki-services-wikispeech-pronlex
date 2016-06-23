package vrules

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
)

var fsExp = "Expected: '%v' got: '%v'"

type testMustHaveTrans struct {
}

func (r testMustHaveTrans) Validate(e lex.Entry) []validation.Result {
	name := "MustHaveTrans"
	level := "Format"
	var result = make([]validation.Result, 0)
	if len(e.Transcriptions) == 0 {
		result = append(result, validation.Result{
			RuleName: name,
			Level:    level,
			Message:  "At least one transcription is required"})
	}
	return result
}

type testNoEmptyTrans struct {
}

func (r testNoEmptyTrans) Validate(e lex.Entry) []validation.Result {
	name := "NoEmptyTrans"
	level := "Format"
	var result = make([]validation.Result, 0)
	for _, t := range e.Transcriptions {
		if len(strings.TrimSpace(t.Strn)) == 0 {
			result = append(result, validation.Result{
				RuleName: name,
				Level:    level,
				Message:  "Empty transcriptions are not allowed"})
		}
	}
	return result
}

type testDecomp2Orth struct {
}

func (r testDecomp2Orth) Validate(e lex.Entry) []validation.Result {
	name := "Decomp2Orth"
	level := "Fatal"
	var result = make([]validation.Result, 0)
	expectOrth := strings.Replace(e.WordParts, "+", "", -1)
	if expectOrth != e.Strn {
		result = append(result, validation.Result{
			RuleName: name,
			Level:    level,
			Message:  fmt.Sprintf("decomp/orth mismatch: %s/%s", e.WordParts, e.Strn)})
	}
	return result
}

func Test1(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testMustHaveTrans{}, testNoEmptyTrans{}}}

	var e = &lex.Entry{
		Strn:         "anka",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "anka",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" a N . k a",
				Language: "swe",
			},
		},
	}

	var result = vali.Validate([]*lex.Entry{e})

	if result != true {
		t.Errorf(fsExp, make([]validation.Result, 0), result)
	}
}

func Test2(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testMustHaveTrans{}, testNoEmptyTrans{}}}

	var e = &lex.Entry{
		Strn:           "anka",
		Language:       "swe",
		PartOfSpeech:   "NN",
		WordParts:      "anka",
		Transcriptions: []lex.Transcription{},
	}

	vali.Validate([]*lex.Entry{e})
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "MustHaveTrans",
			Level:    "Format",
			Message:  "[...]",
		},
	}

	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	} else {
		if result[0].RuleName != "MustHaveTrans" {
			t.Errorf(fsExp, expect, result)
		}
	}
}

func Test3(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testMustHaveTrans{}, testNoEmptyTrans{}, testDecomp2Orth{}}}

	var e = &lex.Entry{
		Strn:         "ankstjärt",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "ank+sjärt",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" a N k + % x { rt",
				Language: "swe",
			},
		},
	}

	vali.Validate([]*lex.Entry{e})
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "Decomp2Orth",
			Level:    "Fatal",
			Message:  "[...]",
		},
	}
	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	} else {
		if result[0].RuleName != "Decomp2Orth" {
			t.Errorf(fsExp, expect, result)
		}
	}
}

func Test4(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testMustHaveTrans{}, testNoEmptyTrans{}, testDecomp2Orth{}}}

	var e = &lex.Entry{
		Strn:         "ankstjärtsbad",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "ank+stjärts+bad",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" a N k + x { rt rs + % b A: d",
				Language: "swe",
			},
		},
	}

	vali.Validate([]*lex.Entry{e})
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}
}

func TestNst1(t *testing.T) {
	var vali, err = NewNSTDemoValidator()
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	var e = &lex.Entry{
		Strn:         "banen",
		Language:     "nor",
		PartOfSpeech: "NN",
		WordParts:    "banen",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\" b A: $ n @ n",
				Language: "nor",
			},
		},
	}

	vali.Validate([]*lex.Entry{e})
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{}
	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = &lex.Entry{
		Strn:         "banen",
		Language:     "nor",
		PartOfSpeech: "NN",
		WordParts:    "banen",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "b a $ \" n e: n",
				Language: "nor",
			},
		},
	}

	vali.Validate([]*lex.Entry{e})
	result = e.EntryValidations

	expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "final_nostress_nolong",
			Level:    "Warning",
			Message:  "[...]"},
	}

	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = &lex.Entry{
		Strn:         "bantorget",
		Language:     "nor",
		PartOfSpeech: "NN",
		WordParts:    "ban+torget",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" b A: n - % t O r $ g @ t",
				Language: "nor",
			},
		},
	}

	vali.Validate([]*lex.Entry{e})
	result = e.EntryValidations

	expect = []lex.EntryValidation{}

	if len(result) != len(expect) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = &lex.Entry{
		Strn:         "battorget",
		Language:     "nor",
		PartOfSpeech: "NN",
		WordParts:    "bat+torget",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" b A: t - % t O r $ g @ t",
				Language: "nor",
			},
		},
	}

	vali.Validate([]*lex.Entry{e})
	result = e.EntryValidations

	expect = []lex.EntryValidation{}

	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = &lex.Entry{
		Strn:         "battorget",
		Language:     "nor",
		PartOfSpeech: "NN",
		WordParts:    "batt+torget",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" b a t - % t O r $ g @ t",
				Language: "nor",
			},
		},
	}

	vali.Validate([]*lex.Entry{e})
	result = e.EntryValidations

	expect = []lex.EntryValidation{}

	if len(result) != len(expect) || (len(expect) > 0 && result[0].RuleName != expect[0].RuleName) {
		t.Errorf(fsExp, expect, result)
	}

	//

	e = &lex.Entry{
		Strn:         "batttorget",
		Language:     "nor",
		PartOfSpeech: "NN",
		WordParts:    "batt+torget",
		Transcriptions: []lex.Transcription{
			lex.Transcription{
				Strn:     "\"\" b a t - % t O r $ g @ t",
				Language: "nor",
			},
		},
	}

	vali.Validate([]*lex.Entry{e})
	result = e.EntryValidations

	expect = []lex.EntryValidation{
		lex.EntryValidation{
			RuleName: "Decomp2Orth",
			Level:    "Fatal",
			Message:  "[...]"},
	}

	if len(result) != len(expect) || result[0].RuleName != expect[0].RuleName {
		t.Errorf(fsExp, expect, result)
	}
}

func TestRegexp2Backrefs(t *testing.T) {
	reFrom, err := regexp2.Compile("(.)\\1[+]\\1", regexp2.None)
	if err != nil {
		t.Errorf("%v", err)
	}
	//

	reTo := "$1+$1"
	input := "hatt+torget"
	expect := "hat+torget"
	result, err := reFrom.Replace(input, reTo, 0, -1)
	if err != nil {
		t.Errorf("%v", err)
	}
	if result != expect {
		t.Errorf(fsExp, expect, result)
	}

	//

	input = "hats+torget"
	expect = "hats+torget"
	result, err = reFrom.Replace(input, reTo, 0, -1)
	if err != nil {
		t.Errorf("%v", err)
	}

	if result != expect {
		t.Errorf(fsExp, expect, result)
	}

}
