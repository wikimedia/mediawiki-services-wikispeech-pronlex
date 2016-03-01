package dbapi

import (
	"sort"
	"testing"
)

// Test that TranscriptionSlice sorts on .Id
func TestStruct_TranscriptionSliceSort(t *testing.T) {
	t1 := Transcription{Id: 47}
	t2 := Transcription{Id: 1047}
	t3 := Transcription{Id: 11}

	ts1 := []Transcription{t1, t2, t3}
	sort.Sort(TranscriptionSlice(ts1))
	if ts1[0].Id != 11 {
		t.Errorf(fs, 11, ts1[0].Id)
	}
	if ts1[1].Id != 47 {
		t.Errorf(fs, 47, ts1[1].Id)
	}
	if ts1[2].Id != 1047 {
		t.Errorf(fs, 1047, ts1[2].Id)
	}
}
