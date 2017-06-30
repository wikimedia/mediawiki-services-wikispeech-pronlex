package lex

import "testing"

func Test_ParseLexRef(t *testing.T) {

	n1 := "sv_se-nst:full_words:v1.0"
	lr1, err := ParseLexRef(n1)
	db1 := string(lr1.DBRef)
	l1 := string(lr1.LexName)
	if w, g := "sv_se-nst", db1; w != g {
		t.Errorf("wanted %s got '%s'", w, g)
	}
	if w, g := "full_words:v1.0", l1; w != g {
		t.Errorf("wanted %s got '%s'", w, g)
	}
	if err != nil {
		t.Errorf("Auch! %v", err)
	}

	// Invalid db name
	n2 := ":full_words:v1.0"
	lr2, err := ParseLexRef(n2)
	db2 := string(lr2.DBRef)
	l2 := string(lr2.LexName)
	if err == nil {
		t.Errorf("wanted error, got nil")
	}
	if w, g := "", db2; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}
	if w, g := "", l2; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

	// Invalid db name
	n3 := "full_wordsv1.0:"
	lr3, err := ParseLexRef(n3)
	db3 := string(lr3.DBRef)
	l3 := string(lr3.LexName)
	if err == nil {
		t.Errorf("wanted error, got nil")
	}
	if w, g := "", db3; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}
	if w, g := "", l3; w != g {
		t.Errorf("wanted '%s' got '%s'", w, g)
	}

}
