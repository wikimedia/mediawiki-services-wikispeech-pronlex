package line

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/stts-se/pronlex/lex"

	"testing"
)

func Test_NewWS(t *testing.T) {

	nws, err := NewWS()
	if err != nil {
		t.Errorf("My heart bleeds: %v", err)
	}
	_ = nws // Hooray!
}

func TestParseComments(t *testing.T) {
	var cmtString string
	var got, expect []lex.EntryComment
	var nws WS
	var err error

	nws, err = NewWS()
	if err != nil {
		t.Errorf("My heart bleeds: %v", err)
	}

	//
	cmtString = "[other: ] (nisse) §§§ [spelling: this is a typo] (kalle)"
	got, err = nws.parseComments(cmtString)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	expect = []lex.EntryComment{
		{Label: "other", Comment: "", Source: "nisse"},
		{Label: "spelling", Comment: "this is a typo", Source: "kalle"},
	}

	if !reflect.DeepEqual(got, expect) {
		t.Errorf("Expected %#v, found %#v", expect, got)
	}

	//
	cmtString = "[other: ] (nisse) §§§ [spelling: this is a typo] (kalle)"
	got, err = nws.parseComments(cmtString)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	expect = []lex.EntryComment{
		{Label: "other", Comment: "", Source: "nisse"},
		{Label: "spelling", Comment: "this is a typo", Source: "kalle"},
	}

	if !reflect.DeepEqual(got, expect) {
		t.Errorf("Expected %#v, found %#v", expect, got)
	}

}

func Test_WSParse_01(t *testing.T) {
	ws, err := NewWS()
	if err != nil {
		t.Errorf("didn't expect error here : %s", err)
		return
	}

	input := `anka	NN	SIN|IND|NOM|UTR	anka	anka	s1a-flicka	sv-se	"" a N . k a	sv-se							imported	hanna	false	duck	`

	expect := lex.Entry{
		Strn:         "anka",
		PartOfSpeech: "NN",
		Morphology:   "SIN|IND|NOM|UTR",
		WordParts:    "anka",
		Language:     "sv-se",
		Lemma: lex.Lemma{
			Strn:     "anka",
			Reading:  "",
			Paradigm: "s1a-flicka",
		},
		Transcriptions: []lex.Transcription{
			{
				Strn:     "\"\" a N . k a",
				Language: "sv-se",
				Sources:  []string{"hanna"},
			},
		},
		Tag: "duck",
	}
	result, err := ws.ParseToEntry(input)
	if err != nil {
		t.Errorf("didn't expect error here : %v", err)
	} else {
		checkWSResult(t, expect, result)
	}

}

func checkWSResult(t *testing.T, x lex.Entry, r lex.Entry) {
	checkWSResultField(t, Orth.String(), x.Strn, r.Strn)
	checkWSResultField(t, Pos.String(), x.PartOfSpeech, r.PartOfSpeech)
	checkWSResultField(t, Morph.String(), x.Morphology, r.Morphology)
	checkWSResultField(t, WordParts.String(), x.WordParts, r.WordParts)
	checkWSResultField(t, Lang.String(), x.Language, r.Language)
	checkWSResultField(t, "Lemma.Reading", x.Lemma.Reading, r.Lemma.Reading)
	checkWSResultField(t, "Lemma.Strn", x.Lemma.Strn, r.Lemma.Strn)
	checkWSResultField(t, "Lemma.Paradigm", x.Lemma.Paradigm, r.Lemma.Paradigm)

	if len(x.Transcriptions) != len(r.Transcriptions) {
		t.Errorf("Expected %v transcriptions, got %v", len(x.Transcriptions), len(r.Transcriptions))
	} else {
		for i, trans := range x.Transcriptions {
			transID := fmt.Sprintf("Trans%d", (i + 1))
			translangID := fmt.Sprintf("Translang%d", (i + 1))
			transsourceID := fmt.Sprintf("Transsource%d", (i + 1))
			checkWSResultField(t, transID, trans.Strn, r.Transcriptions[i].Strn)
			checkWSResultField(t, translangID, trans.Language, r.Transcriptions[i].Language)
			checkWSResultField(t, transsourceID, strings.Join(trans.Sources, ", "), strings.Join(r.Transcriptions[i].Sources, ","))
		}
	}
}

func checkWSResultField(t *testing.T, field string, x string, r string) {
	if x != r {
		t.Errorf(fsExpField, field, x, r)
	}
}
