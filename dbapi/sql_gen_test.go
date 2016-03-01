package dbapi

import (
	"testing"
)

var fs = "Xpctd:\t'%v' got:\t'%v'"

func TestSql_RemoveEmptyStrings(t *testing.T) {
	ss := []string{"", "a", ""}
	ss = RemoveEmptyStrings(ss)
	x := 1
	if len(ss) != x {
		t.Errorf(fs, x, len(ss))
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

func TestSql_idiotSql(t *testing.T) {
	q := NewQuery()
	s, _ := idiotSql(q)
	x := "select entry.id from entry limit 25 offset 0"
	if s != x {
		t.Errorf(fs, x, s)
	}

	q = NewQuery() // Query{Lexicons : []Lexicon{Lexicon{}, Lexicon{}}}
	q.Lexicons = []Lexicon{Lexicon{}, Lexicon{}}
	s, _ = idiotSql(q)
	x = "select entry.id from lexicon, entry where lexicon.id in (?,?) limit 25 offset 0"
	if s != x {
		t.Errorf(fs, x, s)
	}

}
