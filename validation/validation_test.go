package validation

import (
	"strings"
	"testing"

	"github.com/stts-se/pronlex/dbapi"
)

var fsExp = "Expected: '%v' got: '%v'"

type mustHaveTrans struct {
}

func (r mustHaveTrans) Validate(e dbapi.Entry) []Result {
	var result = make([]Result, 0)
	if len(e.Transcriptions) == 0 {
		result = append(result, Result{"MustHaveTrans", "Format", "At least one transcription is required"})
	}
	return result
}

type noEmptyTrans struct {
}

func (r noEmptyTrans) Validate(e dbapi.Entry) []Result {
	var result = make([]Result, 0)
	for _, t := range e.Transcriptions {
		if len(strings.TrimSpace(t.Strn)) == 0 {
			result = append(result, Result{"NoEmptyTrans", "Format", "Empty transcriptions are not allowed"})
		}
	}
	return result
}

func Test1(t *testing.T) {
	var vali = Validator{
		[]Rule{mustHaveTrans{}, noEmptyTrans{}}}

	var e = dbapi.Entry{
		Strn:         "anka",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "anka",
		Transcriptions: []dbapi.Transcription{
			dbapi.Transcription{
				Strn:     "\"\" a N . k a",
				Language: "swe",
			},
		},
	}

	var result = vali.Validate([]dbapi.Entry{e})

	if len(result) != 0 {
		t.Errorf(fsExp, make([]Result, 0), result)
	}
}

func Test2(t *testing.T) {
	var vali = Validator{
		[]Rule{mustHaveTrans{}, noEmptyTrans{}}}

	var e = dbapi.Entry{
		Strn:           "anka",
		Language:       "swe",
		PartOfSpeech:   "NN",
		WordParts:      "anka",
		Transcriptions: []dbapi.Transcription{},
	}

	var result = vali.Validate([]dbapi.Entry{e})

	var expect = []Result{
		Result{
			RuleName: "MustHaveTrans",
			Level:    "Format",
			Message:  "[...]",
		},
	}

	if len(result) != 1 {
		t.Errorf(fsExp, expect, result)
	} else {
		if result[0].RuleName != "MustHaveTrans" {
			t.Errorf(fsExp, expect, result)
		}
	}
}

func Test3(t *testing.T) {
	var vali = Validator{
		[]Rule{mustHaveTrans{}, noEmptyTrans{}}}

	var e = dbapi.Entry{
		Strn:         "anka",
		Language:     "swe",
		PartOfSpeech: "NN",
		WordParts:    "anka",
		Transcriptions: []dbapi.Transcription{
			dbapi.Transcription{
				Strn:     "   ",
				Language: "swe",
			},
		},
	}

	var result = vali.Validate([]dbapi.Entry{e})

	var expect = []Result{
		Result{
			RuleName: "NoEmptyTrans",
			Level:    "Format",
			Message:  "[...]",
		},
	}

	if len(result) != 1 {
		t.Errorf(fsExp, expect, result)
	} else {
		if result[0].RuleName != "NoEmptyTrans" {
			t.Errorf(fsExp, expect, result)
		}
	}
}
