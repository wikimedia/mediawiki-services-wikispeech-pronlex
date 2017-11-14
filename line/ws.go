package line

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/stts-se/pronlex/lex"
)

// WS implements the line.Parser interface
type WS struct {
	format Format
}

// Format is the line.Format instance used for line parsing inside of this parser
func (ws WS) Format() Format {
	return ws.format
}

// Parse is used for parsing input lines (calls underlying Format.Parse)
func (ws WS) Parse(line string) (map[Field]string, error) {
	return ws.format.Parse(line)
}

// ParseToEntry is used for parsing input lines (calls underlying Format.Parse)
func (ws WS) ParseToEntry(line string) (lex.Entry, error) {
	res := lex.Entry{}

	fs, err := ws.Parse(line)
	if err != nil {
		return res, fmt.Errorf("parse to entry failed : %v", err)
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
	res.Preferred, err = strconv.ParseBool(fs[Preferred])
	res.Tag = fs[Tag]
	if err != nil {
		err := fmt.Errorf("couldn't convert string to boolean preferred field: %v", err)
		return res, fmt.Errorf("parse to entry failed : %v", err)
	}

	err = ws.sanityChecks(res)
	if err != nil {
		return res, fmt.Errorf("parse to entry failed for line %s: %v", line, err)
	}

	return res, nil
}

func (ws WS) sanityChecks(e lex.Entry) error {
	if strings.TrimSpace(e.Strn) == "" {
		return fmt.Errorf("input orthography cannot be empty")
	}
	if len(e.Transcriptions) == 0 {
		return fmt.Errorf("there must be at least one input transcription")
	}
	for _, t := range e.Transcriptions {
		if strings.TrimSpace(t.Strn) == "" {
			return fmt.Errorf("input transcription cannot be empty")
		}
	}
	return nil
}

// String is used to generate an output line from a set of fields (calls underlying Format.String)
func (ws WS) String(fields map[Field]string) (string, error) {
	return ws.format.String(fields)
}

// Entry2String is used to generate an output line from a lex.Entry (calls underlying Format.String)
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
	fs[Preferred] = strconv.FormatBool(e.Preferred)
	fs[Tag] = e.Tag

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

// NewWS is used to create a new instance of the WS parser
func NewWS() (WS, error) {
	tests := []FormatTest{
		{"storstaden	NN	SIN|DEF|NOM|UTR	stor+staden	storstad|95522	s111n, a->ä, stad	SWE	\"\"stu:$%s`t`A:$den	SWE							imported	nst	false	big_city",
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
				Preferred:    "false",
				Tag:          "big_city",
			},
			"storstaden	NN	SIN|DEF|NOM|UTR	stor+staden	storstad|95522	s111n, a->ä, stad	SWE	\"\"stu:$%s`t`A:$den	SWE							imported	nst	false	big_city",
		},
		{"storstaden	NN	SIN|DEF|NOM|UTR	stor+staden	storstad|95522	s111n, a->ä, stad	SWE	\"\"stu:$%s`t`A:$den	SWE							imported	nst	true	",
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
				Preferred:    "true",
				Tag:          "",
			},
			"storstaden	NN	SIN|DEF|NOM|UTR	stor+staden	storstad|95522	s111n, a->ä, stad	SWE	\"\"stu:$%s`t`A:$den	SWE							imported	nst	true	",
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
			Preferred:    17,
			Tag:          18,
		},
		19,
		tests,
	)
	if err != nil {
		return WS{}, err
	}
	return WS{f}, nil
}
