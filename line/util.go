package line

import (
	"fmt"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/symbolset/mapper"
)

// MapTranscriptions maps the input entry's transcriptions (in-place)
func MapTranscriptions(m mapper.Mapper, e *lex.Entry) error {
	var newTs []lex.Transcription
	var errs []string
	for _, t := range e.Transcriptions {
		newT, err := m.MapTranscription(t.Strn)
		if err != nil {
			errs = append(errs, err.Error())
		}
		newTs = append(newTs, lex.Transcription{ID: t.ID, Strn: newT, EntryID: t.EntryID, Language: t.Language, Sources: t.Sources})
	}
	e.Transcriptions = newTs
	if len(errs) > 0 {
		return fmt.Errorf("%v", strings.Join(errs, "; "))
	}
	return nil
}

func equals(expect map[Field]string, result map[Field]string) bool {
	if len(expect) != len(result) {
		return false
	}
	for f, expS := range expect {
		resS := result[f]
		if resS != expS {
			//fmt.Printf("%v: %v vs %v", f, resS, expS)
			return false
		}
	}
	return true
}

// Equals compares two line.Format instances
func (f Format) Equals(other Format) bool {
	if f.Name != other.Name {
		return false
	}
	if f.FieldSep != other.FieldSep {
		return false
	}
	if f.NFields != other.NFields {
		return false
	}
	if len(f.Fields) != len(other.Fields) {
		return false
	}
	for f, expS := range f.Fields {
		resS := other.Fields[f]
		if resS != expS {
			return false
		}
	}
	return true
}

// stringSlice used for sorting if necessary
// type stringSlice []string

// func (a stringSlice) Len() int      { return len(a) }
// func (a stringSlice) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
// func (a stringSlice) Less(i, j int) bool {
// 	return a[i] < a[j]
// }
