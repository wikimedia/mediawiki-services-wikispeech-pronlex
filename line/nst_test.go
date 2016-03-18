package line

import (
	"fmt"
	"testing"

	"github.com/stts-se/pronlex/dbapi"
)

//var fsExpField = "For field %v, expected: '%v' got: '%v'"
//var fsExp = "Expected: '%v' got: '%v'"

func checkNSTResultField(t *testing.T, field string, x string, r string) {
	if x != r {
		t.Errorf(fsExpField, field, x, r)
	}
}

func checkNSTResult(t *testing.T, x dbapi.Entry, r dbapi.Entry) {
	checkNSTResultField(t, Orth.String(), x.Strn, r.Strn)
	checkNSTResultField(t, Pos.String(), x.PartOfSpeech, r.PartOfSpeech)
	checkNSTResultField(t, WordParts.String(), x.WordParts, r.WordParts)
	checkNSTResultField(t, Lang.String(), x.Language, r.Language)
	checkNSTResultField(t, "Lemma.Reading", x.Lemma.Reading, r.Lemma.Reading)
	checkNSTResultField(t, "Lemma.Strn", x.Lemma.Strn, r.Lemma.Strn)
	checkNSTResultField(t, "Lemma.Paradigm", x.Lemma.Paradigm, r.Lemma.Paradigm)

	if len(x.Transcriptions) != len(r.Transcriptions) {
		t.Errorf("Expected %v transcriptions, got %v", len(x.Transcriptions), len(r.Transcriptions))
	} else {
		for i, trans := range x.Transcriptions {
			transID := fmt.Sprintf("Trans%d", (i + 1))
			translangID := fmt.Sprintf("Translang%d", (i + 1))
			checkNSTResultField(t, transID, trans.Strn, r.Transcriptions[i].Strn)
			checkNSTResultField(t, translangID, trans.Language, r.Transcriptions[i].Language)
		}
	}
}

func Test_NewNST(t *testing.T) {
	_, err := NewNST()
	if err != nil {
		t.Errorf("didn't expect error here : %s", err)
	}
}

func Test_NSTParse_01(t *testing.T) {
	nst, err := NewNST()
	if err != nil {
		t.Errorf("didn't expect error here : %s", err)
		return
	}

	input := "storstaden;NN;SIN|DEF|NOM|UTR;stor+staden;JJ+NN;LEX|INFL;SWE;;;;;\"\"stu:$%s`t`A:$den;1;STD;SWE;\"\"stu:$%s`t`A:n;;;SWE;;;;;;;;;;18174;enter_se|inflector;;INFLECTED;storstad|95522;s111n, a->ä, stad;s111;;;;;;;;;;;;;storstaden;;;88748"

	expect := dbapi.Entry{
		Strn:         "storstaden",
		PartOfSpeech: "NN SIN|DEF|NOM|UTR",
		WordParts:    "stor+staden",
		Language:     "SWE",
		Lemma: dbapi.Lemma{
			Strn:     "storstad",
			Reading:  "95522",
			Paradigm: "s111n, a->ä, stad",
		},
		Transcriptions: []dbapi.Transcription{
			dbapi.Transcription{
				Strn:     "\"\"stu:$%s`t`A:$den",
				Language: "SWE",
			},
			dbapi.Transcription{
				Strn:     "\"\"stu:$%s`t`A:n",
				Language: "SWE",
			},
		},
	}
	result, err := nst.Parse(input)
	if err != nil {
		t.Errorf("didn't expect error here : %v", err)
	} else {
		checkNSTResult(t, expect, result)
	}

}
