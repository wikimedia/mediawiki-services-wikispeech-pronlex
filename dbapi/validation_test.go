package dbapi

import (
	"database/sql"

	"log"
	"os"
	"testing"

	"reflect"

	"github.com/dlclark/regexp2"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/pronlex/vrules"
)

var vfs = "Wanted: '%v' got: '%v'"
var vDbPath = "./vtestlex.db"

func createValidator() validation.Validator {
	name := "sampa"
	symbols := []symbolset.Symbol{
		symbolset.Symbol{"a", symbolset.Syllabic, ""},
		symbolset.Symbol{"A:", symbolset.Syllabic, ""},
		symbolset.Symbol{"b", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"p", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"N", symbolset.NonSyllabic, ""},
		symbolset.Symbol{"n", symbolset.NonSyllabic, ""},
		symbolset.Symbol{" ", symbolset.PhonemeDelimiter, ""},
		symbolset.Symbol{".", symbolset.SyllableDelimiter, ""},
		symbolset.Symbol{"\"", symbolset.Stress, ""},
		symbolset.Symbol{"\"\"", symbolset.Stress, ""},
	}
	ss, err := symbolset.NewSymbols(name, symbols)
	ff("failed to init symbols : %v", err)

	primaryStressRe, err := vrules.ProcessTransRe(ss, "\"")
	ff("%v", err)

	syllabicRe, err := vrules.ProcessTransRe(ss, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	ff("%v", err)

	reFrom, err := regexp2.Compile("(.)\\1[+]\\1", regexp2.None)
	ff("%v", err)

	decomp2Orth := vrules.Decomp2Orth{CompDelim: "+",
		AcceptEmptyDecomp: true,
		PreFilterWordPartString: func(s string) (string, error) {
			res, err := reFrom.Replace(s, "$1+$1", 0, -1)
			if err != nil {
				return s, err
			}
			return res, nil
		}}

	repeatedPhnRe, err := vrules.ProcessTransRe(ss, "symbol( +[.~])? +\\1")
	ff("%v", err)

	var v = validation.Validator{
		Name: ss.Name,
		Rules: []validation.Rule{
			vrules.MustHaveTrans{},
			vrules.NoEmptyTrans{},
			vrules.RequiredTransRe{
				Name:    "primary_stress",
				Level:   "Fatal",
				Message: "Primary stress required",
				Re:      primaryStressRe,
			},
			vrules.RequiredTransRe{
				Name:    "syllabic",
				Level:   "Format",
				Message: "Each syllable needs a syllabic phoneme",
				Re:      syllabicRe,
			},
			vrules.IllegalTransRe{
				Name:    "repeated_phonemes",
				Level:   "Fatal",
				Message: "Repeated phonemes cannot be used within the same morpheme",
				Re:      repeatedPhnRe,
			},
			decomp2Orth,
			vrules.SymbolSetRule{
				SymbolSet: ss,
			},
		}}
	return v
}

func vInsertEntries(t *testing.T, lexName string) (*sql.DB, int64) {

	if _, err := os.Stat(vDbPath); !os.IsNotExist(err) {
		err := os.Remove(vDbPath)
		ff("failed to remove "+vDbPath+" : %v", err)
	}

	db, err := sql.Open("sqlite3_with_regexp", vDbPath)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
	ff("Failed to exec PRAGMA call %v", err)

	//defer db.Close()

	_, err = db.Exec(Schema) // Creates new lexicon database
	ff("Failed to create lexicon db: %v", err)

	l := Lexicon{Name: lexName, SymbolSetName: "ZZ"}

	l, err = InsertLexicon(db, l)
	if err != nil {
		t.Errorf(vfs, nil, err)
	}

	lxs, err := ListLexicons(db)
	if err != nil {
		t.Errorf(vfs, nil, err)
	}
	if len(lxs) != 1 {
		t.Errorf(vfs, 1, len(lxs))
	}

	t1a := lex.Transcription{Strn: "\" A: p a", Language: "sv-se"}
	t1b := lex.Transcription{Strn: "\" a p a", Language: "sv-se"}

	e1 := lex.Entry{Strn: "apa",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apa",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t1a, t1b},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	t2a := lex.Transcription{Strn: "\" A: p a n", Language: "sv-se"}
	t2b := lex.Transcription{Strn: "A: p a n", Language: "sv-se"}

	e2 := lex.Entry{Strn: "apan",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "apan",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t2a, t2b},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	t3a := lex.Transcription{Strn: "\" A . p a n", Language: "sv-se"}
	e3 := lex.Entry{Strn: "appan",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "appan",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t3a},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	t4a := lex.Transcription{Strn: "\" A . p a n", Language: "sv-se"}
	e4 := lex.Entry{Strn: "appan",
		PartOfSpeech:   "NN",
		Morphology:     "NEU UTR",
		WordParts:      "appan",
		Language:       "XYZZ",
		Transcriptions: []lex.Transcription{t4a},
		EntryStatus:    lex.EntryStatus{Name: "old", Source: "tst"}}

	_, errx := InsertEntries(db, l, []lex.Entry{e1, e2, e3, e4})
	if errx != nil {
		t.Errorf(vfs, "nil", errx)
	}
	return db, l.ID
}

func Test_Validation1(t *testing.T) {
	db, lexId := vInsertEntries(t, "test1")
	v := createValidator()

	q := Query{}

	stats, err := Validate(db, SilentLogger{}, v, q)
	ff("validation failed : %v", err)

	expect := ValStats{
		Values: map[string]int{
			"Total entries":                4,
			"Total validations":            5,
			"Invalid entries":              3,
			"Level: fatal":                 3,
			"Level: format":                2,
			"Rule: symbolset (fatal)":      2,
			"Rule: primary_stress (fatal)": 1,
			"Rule: syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}

	// check stats saved in db
	lexStats, err := ValidationStats(db, lexId)
	ff("validation stats failed : %v", err)

	expectFull := ValStats{
		Values: map[string]int{
			"Total entries":                4,
			"Total validations":            5,
			"Invalid entries":              3,
			"Level: fatal":                 3,
			"Level: format":                2,
			"Rule: symbolset (fatal)":      2,
			"Rule: primary_stress (fatal)": 1,
			"Rule: syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expectFull, lexStats) {
		t.Errorf(vfs, expectFull, lexStats)
	}
}

func Test_Validation2(t *testing.T) {
	db, lexId := vInsertEntries(t, "test2")
	v := createValidator()

	q := Query{WordRegexp: "a$"}

	// test 1
	stats, err := Validate(db, SilentLogger{}, v, q)
	ff("validation failed : %v", err)

	expect := ValStats{
		Values: map[string]int{
			"Total entries":     1,
			"Total validations": 0,
			"Invalid entries":   0,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}

	// check stats saved in db
	lexStats, err := ValidationStats(db, lexId)
	ff("validation stats failed : %v", err)

	expectFull := ValStats{
		Values: map[string]int{
			"Total entries":     4,
			"Total validations": 0,
			"Invalid entries":   0,
		},
	}

	if !reflect.DeepEqual(expectFull, lexStats) {
		t.Errorf(vfs, expectFull, lexStats)
	}

	// test 2
	q = Query{}

	stats, err = Validate(db, SilentLogger{}, v, q)
	ff("validation failed : %v", err)

	expect = ValStats{
		Values: map[string]int{
			"Total entries":                4,
			"Total validations":            5,
			"Invalid entries":              3,
			"Level: fatal":                 3,
			"Level: format":                2,
			"Rule: symbolset (fatal)":      2,
			"Rule: primary_stress (fatal)": 1,
			"Rule: syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}

	// check stats saved in db
	lexStats, err = ValidationStats(db, lexId)
	ff("validation stats failed : %v", err)

	expectFull = ValStats{
		Values: map[string]int{
			"Total entries":                4,
			"Total validations":            5,
			"Invalid entries":              3,
			"Level: fatal":                 3,
			"Level: format":                2,
			"Rule: symbolset (fatal)":      2,
			"Rule: primary_stress (fatal)": 1,
			"Rule: syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expectFull, lexStats) {
		t.Errorf(vfs, expectFull, lexStats)
	}
}

func Test_ValidationUpdate1(t *testing.T) {
	db, lexId := vInsertEntries(t, "test3")
	v := createValidator()
	ew := lex.EntrySliceWriter{}
	err := LookUp(db, Query{}, &ew)
	ff("lookup failed : %v", err)

	for _, e := range ew.Entries {
		e, _ = v.ValidateEntry(e)
		err = UpdateValidation(db, []lex.Entry{e})
		ff("update validation failed : %v", err)
	}

	stats, err := ValidationStats(db, lexId)
	ff("validation stats failed : %v", err)

	expect := ValStats{
		Values: map[string]int{
			"Total entries":                4,
			"Total validations":            5,
			"Invalid entries":              3,
			"Level: fatal":                 3,
			"Level: format":                2,
			"Rule: symbolset (fatal)":      2,
			"Rule: primary_stress (fatal)": 1,
			"Rule: syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}
}

func Test_ValidationUpdate2(t *testing.T) {
	db, lexId := vInsertEntries(t, "test4")
	v := createValidator()
	ew := lex.EntrySliceWriter{}
	err := LookUp(db, Query{}, &ew)
	ff("lookup failed : %v", err)

	var es []lex.Entry
	for _, e := range ew.Entries {
		e, _ = v.ValidateEntry(e)
		es = append(es, e)
	}
	err = UpdateValidation(db, es)
	ff("update validation failed : %v", err)

	stats, err := ValidationStats(db, lexId)
	ff("validation stats failed : %v", err)

	expect := ValStats{
		Values: map[string]int{
			"Total entries":                4,
			"Total validations":            5,
			"Invalid entries":              3,
			"Level: fatal":                 3,
			"Level: format":                2,
			"Rule: symbolset (fatal)":      2,
			"Rule: primary_stress (fatal)": 1,
			"Rule: syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}
}

func Test_ValidationUpdate3(t *testing.T) {
	db, lexId := vInsertEntries(t, "test5")
	v := createValidator()
	ew := lex.EntrySliceWriter{}
	err := LookUp(db, Query{}, &ew)
	ff("lookup failed : %v", err)

	es, _ := v.ValidateEntries(ew.Entries)
	err = UpdateValidation(db, es)
	ff("update validation failed : %v", err)

	stats, err := ValidationStats(db, lexId)
	ff("validation stats failed : %v", err)

	expect := ValStats{
		Values: map[string]int{
			"Total entries":                4,
			"Total validations":            5,
			"Invalid entries":              3,
			"Level: fatal":                 3,
			"Level: format":                2,
			"Rule: symbolset (fatal)":      2,
			"Rule: primary_stress (fatal)": 1,
			"Rule: syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}
}
