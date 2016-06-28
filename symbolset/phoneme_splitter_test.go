package symbolset

import (
	"sort"
	"testing"
)

func Test_LenSort(t *testing.T) {

	s0 := "sr"
	s1 := "shrt"
	s2 := "looong"
	ss := []string{s0, s1, s2}

	sort.Sort(ByLength(ss))
	if got, want := ss[0], s2; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}
	if got, want := ss[1], s1; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}
	if got, want := ss[2], s0; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}
}

func Test_SplitIntoPhonemes(t *testing.T) {
	phs := []string{"aa", "a", "bb", "b", "ddddd", "f33"}
	s1 := "c"
	res, unk := SplitIntoPhonemes(phs, s1)
	//fmt.Printf("res: '%v' unk: '%v'\n", res, unk)
	if got, want := res[0], s1; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}
	if got, want := unk[0], s1; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}

	s2 := "a"
	res, unk = SplitIntoPhonemes(phs, s2)
	if got, want := res[0], s2; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}
	if got, want := len(unk), 0; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}

	s3 := "_azbbax"
	res, unk = SplitIntoPhonemes(phs, s3)
	if got, want := res[0], "_"; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}
	if got, want := len(unk), 3; got != want {
		t.Errorf("Got '%v' Wanted '%v'", got, want)
	}

}