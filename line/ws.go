package line

// TODO Better definition of general line format (Reading field? Any number of transcriptions?)

import (
	"fmt"
	"strings"

	"github.com/stts-se/pronlex/lex"
)

// WS implements the line.Parser interface
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
	res.Morphology = fs[Morph]
	res.WordParts = fs[WordParts]
	if strings.TrimSpace(fs[Lemma]) != "" {
		res.Lemma = lex.Lemma{Strn: fs[Lemma], Paradigm: fs[Paradigm]} // TODO Reading : fs[Reading]?
	}
	res.Transcriptions = getTranses(fs) // <-- func getTranses declared in nst.go

	res.EntryStatus.Name = fs[StatusName]
	res.EntryStatus.Source = fs[StatusSource]

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
	fs[Morph] = e.Morphology

	fs[StatusName] = e.EntryStatus.Name
	fs[StatusSource] = e.EntryStatus.Source

	//TODO Missing field for Reading?
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
		FormatTest{"storstaden	NN	SIN|DEF|NOM|UTR	stor+staden	storstad|95522	s111n, a->ä, stad	SWE	\"\"stu:$%s`t`A:$den	SWE							imported	nst",
			map[Field]string{
				Orth:         "storstaden",
				Pos:          "NN",
				Morph:        "SIN|DEF|NOM|UTR",
				WordParts:    "stor+staden",
				Lemma:        "storstad|95522",
				Paradigm:     "s111n, a->ä, stad",
				Lang:         "SWE",
				Trans1:       "\"\"stu:$%s`t`A:$den",
				Translang1:   "SWE",
				Trans2:       "",
				Translang2:   "",
				Trans3:       "",
				Translang3:   "",
				Trans4:       "",
				Translang4:   "",
				StatusName:   "imported",
				StatusSource: "nst",
			},
			"storstaden	NN	SIN|DEF|NOM|UTR	stor+staden	storstad|95522	s111n, a->ä, stad	SWE	\"\"stu:$%s`t`A:$den	SWE							imported	nst",
		},
	}
	f, err := NewFormat(
		"WS",
		"	",
		map[Field]int{
			Orth:         0,
			Pos:          1,
			Morph:        2,
			WordParts:    3,
			Lemma:        4,
			Paradigm:     5,
			Lang:         6,
			Trans1:       7,
			Translang1:   8,
			Trans2:       9,
			Translang2:   10,
			Trans3:       11,
			Translang3:   12,
			Trans4:       13,
			Translang4:   14,
			StatusName:   15,
			StatusSource: 16,
		},
		17,
		tests,
	)
	if err != nil {
		return WS{}, err
	}
	return WS{f}, nil
}
