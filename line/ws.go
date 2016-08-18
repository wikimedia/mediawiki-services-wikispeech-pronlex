package line

// NB: Copy-paste + edit of line/nst.go
// TODO Review by original author of line/nst.go
// TODO Better definition of general line format (Reading field? Any number of transcriptions?)

import (
	"fmt"
	"github.com/stts-se/pronlex/lex"
	"io"
	"strings"
)

type WS struct {
	format Format
}

func (ws WS) Format() Format {
	return ws.format
}

func (ws WS) Parse(line string) (map[Field]string, error) {
	return ws.format.Parse(line)
}

func (ws WS) ParseToEntry(line string) (lex.Entry, error) {
	res := lex.Entry{}

	fs, err := ws.Parse(line)
	if err != nil {
		return res, fmt.Errorf("Parse to entry failed : %v", err)
	}

	res.Strn = fs[Orth]
	res.Language = fs[Lang]
	res.PartOfSpeech = fs[Pos]
	res.WordParts = fs[WordParts]
	if strings.TrimSpace(fs[Lemma]) != "" {
		res.Lemma = lex.Lemma{Strn: fs[Lemma], Paradigm: fs[Paradigm]} // TODO Reading : fs[Reading]?
	}
	res.Transcriptions = getTranses(fs) // <-- func getTranses declared in nst.go

	return res, nil
}

func (ws WS) String(fields map[Field]string) (string, error) {
	return ws.format.String(fields)
}

func (ws WS) Entry2String(e lex.Entry) (string, error) {
	fs, err := ws.fields(e)
	if err != nil {
		return "", err
	}
	s, err := ws.format.String(fs)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (ws WS) fields(e lex.Entry) (map[Field]string, error) {

	// Fields ID and LexiconID are database internal  and not processed here

	var fs = make(map[Field]string)
	fs[Orth] = e.Strn
	fs[Lang] = e.Language
	fs[WordParts] = e.WordParts

	fs[Pos] = e.PartOfSpeech

	//TODO Missing field for Reading
	// Lemma
	// if e.Lemma.Reading != "" {
	// 	fs[Lemma] = e.Lemma.Strn + "|" + e.Lemma.Reading
	// } else {
	// 	fs[Lemma] = e.Lemma.Strn
	// }
	// if e.Lemma.Reading != "" {
	// 	fs[Lemma] = e.Lemma.Strn + "|" + e.Lemma.Reading
	// } else {
	fs[Lemma] = e.Lemma.Strn
	//}
	fs[Paradigm] = e.Lemma.Paradigm

	for i, t := range e.Transcriptions {
		switch i {
		case 0:
			fs[Trans1] = t.Strn
			fs[Translang1] = t.Language
		case 1:
			fs[Trans2] = t.Strn
			fs[Translang2] = t.Language
		case 2:
			fs[Trans3] = t.Strn
			fs[Translang3] = t.Language
		case 3:
			fs[Trans4] = t.Strn
			fs[Translang4] = t.Language
		default:
			return map[Field]string{}, fmt.Errorf("ws line format can contain max 4 transcriptions, but found %v in: %v", len(e.Transcriptions), e)
		}
	}
	return fs, nil
}

func NewWS() (WS, error) {
	tests := []FormatTest{
		FormatTest{"storstaden	NN SIN|DEF|NOM|UTR	stor+staden	storstad|95522	s111n, a->ä, stad	SWE	\"\"stu:$%s`t`A:$den	SWE						",
			map[Field]string{
				Orth:       "storstaden",
				Pos:        "NN SIN|DEF|NOM|UTR",
				WordParts:  "stor+staden",
				Lemma:      "storstad|95522",
				Paradigm:   "s111n, a->ä, stad",
				Lang:       "SWE",
				Trans1:     "\"\"stu:$%s`t`A:$den",
				Translang1: "SWE",
				Trans2:     "",
				Translang2: "",
				Trans3:     "",
				Translang3: "",
				Trans4:     "",
				Translang4: "",
			},
			"storstaden	NN SIN|DEF|NOM|UTR	stor+staden	storstad|95522	s111n, a->ä, stad	SWE	\"\"stu:$%s`t`A:$den	SWE						",
		},
	}
	f, err := NewFormat(
		"WS",
		"	",
		map[Field]int{
			Orth: 0,
			Pos:  1,
			//Morph:      2,
			WordParts:  2,
			Lemma:      3,
			Paradigm:   4,
			Lang:       5,
			Trans1:     6,
			Translang1: 7,
			Trans2:     8,
			Translang2: 9,
			Trans3:     10,
			Translang3: 11,
			Trans4:     12,
			Translang4: 13,
		},
		14,
		tests,
	)
	if err != nil {
		return WS{}, err
	}
	return WS{f}, nil
}

type WSFileWriter struct {
	WS     WS
	Writer io.Writer
}

func (w WSFileWriter) Write(e lex.Entry) error {
	s, err := w.WS.Entry2String(e)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w.Writer, "%s\n", s)
	return err
}