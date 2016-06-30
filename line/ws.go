package line

import (
	"fmt"

	"github.com/stts-se/pronlex/lex"
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
	res.Lemma = lex.Lemma{Strn: fs[Lemma], Paradigm: fs[Paradigm]} // TODO Reading : fs[Reading]?
	res.Transcriptions = getTranses(fs)                            // <-- func getTranses declared in nst.go

	return res, nil
}
