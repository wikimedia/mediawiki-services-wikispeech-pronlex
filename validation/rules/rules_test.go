package rules

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

func (r testMustHaveTrans) Validate(e lex.Entry) (validation.Result, error) {
	name := "MustHaveTrans"
	level := "Format"
	var messages = make([]string, 0)
	if len(e.Transcriptions) == 0 {
		messages = append(messages, "At least one transcription is required")
	}
	return validation.Result{RuleName: name, Level: level, Messages: messages}, nil

}
func (r testMustHaveTrans) ShouldAccept() []lex.Entry {
	return make([]lex.Entry, 0)
}
func (r testMustHaveTrans) ShouldReject() []lex.Entry {
	return make([]lex.Entry, 0)
}
func (r testMustHaveTrans) Name() string {
	return ""
}
func (r testMustHaveTrans) Level() string {
	return ""
}

func (r testMustHaveTrans) AddAccept(entry lex.Entry) {
}

func (r testMustHaveTrans) AddReject(entry lex.Entry) {
}

type testNoEmptyTrans struct {
}

func (r testNoEmptyTrans) Validate(e lex.Entry) (validation.Result, error) {
	name := "NoEmptyTrans"
	level := "Format"
	var messages = make([]string, 0)
	for _, t := range e.Transcriptions {
		if len(strings.TrimSpace(t.Strn)) == 0 {
			messages = append(messages, "Empty transcriptions are not allowed")
		}
	}
	return validation.Result{RuleName: name, Level: level, Messages: messages}, nil

}
func (r testNoEmptyTrans) ShouldAccept() []lex.Entry {
	return make([]lex.Entry, 0)
}
func (r testNoEmptyTrans) ShouldReject() []lex.Entry {
	return make([]lex.Entry, 0)
}
func (r testNoEmptyTrans) Name() string {
	return "NoEmptyTrans"
}
func (r testNoEmptyTrans) Level() string {
	return "Format"
}
func (r testNoEmptyTrans) AddAccept(entry lex.Entry) {
}

func (r testNoEmptyTrans) AddReject(entry lex.Entry) {
}

type testDecomp2Orth struct {
}

func (r testDecomp2Orth) Validate(e lex.Entry) (validation.Result, error) {
	name := "Decomp2Orth"
	level := "Fatal"
	var messages = make([]string, 0)
	expectOrth := strings.Replace(e.WordParts, "+", "", -1)
	if expectOrth != e.Strn {
		messages = append(messages, fmt.Sprintf("decomp/orth mismatch: %s/%s", e.WordParts, e.Strn))
	}
	return validation.Result{RuleName: name, Level: level, Messages: messages}, nil

}
func (r testDecomp2Orth) ShouldAccept() []lex.Entry {
	return make([]lex.Entry, 0)
}
func (r testDecomp2Orth) ShouldReject() []lex.Entry {
	return make([]lex.Entry, 0)
}
func (r testDecomp2Orth) Name() string {
	return ""
}
func (r testDecomp2Orth) Level() string {
	return ""
}

func (r testDecomp2Orth) AddAccept(entry lex.Entry) {
}

func (r testDecomp2Orth) AddReject(entry lex.Entry) {
}

func Test1(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testNoEmptyTrans{}, testMustHaveTrans{}}}

	var e = lex.Entry{
		Strn:         "anka",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "anka",
		Transcriptions: []lex.Transcription{
			{
				Strn:     "\"\" a N . k a",
				Language: "swe",
			},
		},
	}

	var _, result = vali.ValidateEntries([]lex.Entry{e})

	if result != true {
		t.Errorf(fsExp, make([]validation.Result, 0), result)
	}
}

func Test2(t *testing.T) {
	var vali = validation.Validator{
		Rules: []validation.Rule{testNoEmptyTrans{}, testMustHaveTrans{}}}

	var e = lex.Entry{
		Strn:           "anka",
		Language:       "swe",
		PartOfSpeech:   "NN",
		WordParts:      "anka",
		Transcriptions: []lex.Transcription{},
	}

	vali.ValidateEntry(&e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{
		{
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
		Rules: []validation.Rule{testNoEmptyTrans{}, testNoEmptyTrans{}, testDecomp2Orth{}}}

	var e = lex.Entry{
		Strn:         "ankstjärt",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "ank+sjärt",
		Transcriptions: []lex.Transcription{
			{
				Strn:     "\"\" a N k + % x { rt",
				Language: "swe",
			},
		},
	}

	vali.ValidateEntry(&e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{
		{
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
		Rules: []validation.Rule{testNoEmptyTrans{}, testNoEmptyTrans{}, testDecomp2Orth{}}}

	var e = lex.Entry{
		Strn:         "ankstjärtsbad",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "ank+stjärts+bad",
		Transcriptions: []lex.Transcription{
			{
				Strn:     "\"\" a N k + x { rt rs + % b A: d",
				Language: "swe",
			},
		},
	}

	vali.ValidateEntry(&e)
	var result = e.EntryValidations

	var expect = []lex.EntryValidation{}
	if len(result) != len(expect) {
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
