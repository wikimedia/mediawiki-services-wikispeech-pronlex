package dbapi

import (
	"testing"
)

var fs = "Wanted: '%v' got: '%v'"

func TestSql_RemoveEmptyStrings(t *testing.T) {
	ss := []string{"", "a", ""}
	ss = RemoveEmptyStrings(ss)
	x := 1
	if len(ss) != x {
		t.Errorf(fs, x, len(ss))
	}
}

func TestSql_ToLower(t *testing.T) {
	ss := []string{"AaaA", "aB", "cCC"}
	ss = ToLower(ss)
	x := 3
	if len(ss) != x {
		t.Errorf(fs, x, len(ss))
	}

	if ss[0] != "aaaa" {
		t.Errorf(fs, "aaaa", ss[0])
	}
	if ss[1] != "ab" {
		t.Errorf(fs, "ab", ss[1])
	}

	if ss[2] != "ccc" {
		t.Errorf(fs, "ccc", ss[2])
	}

}

func TestSql_words(t *testing.T) {
	w, wv := words(Query{})
	if "" != w {
		t.Error("Gah!")
	}
	if len(wv) != 0 {
		t.Error("Gah2!")
	}

	w, wv = words(Query{Words: []string{"fimbul"}})
	x := "entry.strn in (?)"
	if w != x {
		t.Errorf(fs, x, w)
	}
	if len(wv) != 1 {
		t.Errorf(fs, 1, len(wv))
	}

	w, _ = words(Query{Lexicons: []Lexicon{Lexicon{}}, Words: []string{"fimbul", "vinter"}})
	x = "entry.strn in (?,?) and entry.lexiconid = lexicon.id"
	if w != x {
		t.Errorf(fs, x, w)
	}

}

func TestSql_nQs(t *testing.T) {
	q1 := nQs(0)
	if q1 != "" {
		t.Errorf("Expected empty string, got %q", q1)
	}

	q1 = nQs(2)
	if q1 != "(?,?)" {
		t.Error("X (?,?) got ", q1)
	}
}

// func TestSql_idiotSQL(t *testing.T) {
// 	q := NewQuery()
// 	s, _ := idiotSQL(q)
// 	x := "select entry.id from entry limit 25 offset 0"
// 	if s != x {
// 		t.Errorf(fs, x, s)
// 	}

// 	q = NewQuery() // Query{Lexicons : []Lexicon{Lexicon{}, Lexicon{}}}
// 	q.Lexicons = []Lexicon{Lexicon{}, Lexicon{}}
// 	s, _ = idiotSQL(q)
// 	x = "select entry.id from lexicon, entry where lexicon.id in (?,?) limit 25 offset 0"
// 	if s != x {
// 		t.Errorf(fs, x, s)
// 	}

// }

func TestSql_SelectEntriesSQL(t *testing.T) {
	q := Query{LemmaLike: "%gal_", ReadingLike: "%grus_"}
	sq := selectEntriesSQL(q)
	if sq.sql == "" {
		t.Error(fs, "non empty", sq.sql)
	}

	//fmt.Printf("%s\n\n%v\n", sqlS, vs)

	if len(sq.values) != 2 {
		t.Error(fs, 2, len(sq.values))
	}
}
