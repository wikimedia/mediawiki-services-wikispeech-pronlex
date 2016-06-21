package dbapi

import (
	"sort"
	"testing"

	"github.com/stts-se/pronlex/lex"
)

// Test that TranscriptionSlice sorts on .Id
func TestStruct_TranscriptionSliceSort(t *testing.T) {
	t1 := lex.Transcription{ID: 47}
	t2 := lex.Transcription{ID: 1047}
	t3 := lex.Transcription{ID: 11}

	ts1 := []lex.Transcription{t1, t2, t3}
	sort.Sort(lex.TranscriptionSlice(ts1))
	if ts1[0].ID != 11 {
		t.Errorf(fs, 11, ts1[0].ID)
	}
	if ts1[1].ID != 47 {
		t.Errorf(fs, 47, ts1[1].ID)
	}
	if ts1[2].ID != 1047 {
		t.Errorf(fs, 1047, ts1[2].ID)
	}
}
