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
	"github.com/stts-se/pronlex/validation/rules"
)

var vfs = "Wanted: '%v' got: '%v'"
var vDbPath = "./vtestlex.db"

func createValidator() validation.Validator {
	name := "sampa"
	symbols := []symbolset.Symbol{
		symbolset.Symbol{String: "a", Cat: symbolset.Syllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "A:", Cat: symbolset.Syllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "b", Cat: symbolset.NonSyllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "p", Cat: symbolset.NonSyllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "N", Cat: symbolset.NonSyllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "n", Cat: symbolset.NonSyllabic, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: " ", Cat: symbolset.PhonemeDelimiter, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: ".", Cat: symbolset.SyllableDelimiter, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "\"", Cat: symbolset.Stress, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
		symbolset.Symbol{String: "\"\"", Cat: symbolset.Stress, Desc: "", IPA: symbolset.IPASymbol{String: "", Unicode: ""}},
	}
	ss, err := symbolset.NewSymbolSet(name, symbols)
	ff("failed to init symbols : %v", err)

	primaryStressRe, err := rules.ProcessTransRe(ss, "\"")
	ff("%v", err)

	syllabicRe, err := rules.ProcessTransRe(ss, "^(\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*( (.|-) (\"\"|\"|%)? *(nonsyllabic +)*syllabic( +nonsyllabic)*)*$")
	ff("%v", err)

	reFrom, err := regexp2.Compile("(.)\\1[+]\\1", regexp2.None)
	ff("%v", err)

	decomp2Orth := rules.Decomp2Orth{CompDelim: "+",
		AcceptEmptyDecomp: true,
		PreFilterWordPartString: func(s string) (string, error) {
			res, err := reFrom.Replace(s, "$1+$1", 0, -1)
			if err != nil {
				return s, err
			}
			return res, nil
		}}

	repeatedPhnRe, err := rules.ProcessTransRe(ss, "symbol( +[.~])? +\\1")
	ff("%v", err)

	var v = validation.Validator{
		Name: ss.Name,
		Rules: []validation.Rule{
			rules.MustHaveTrans{},
			rules.NoEmptyTrans{},
			rules.RequiredTransRe{
				NameStr:  "primary_stress",
				LevelStr: "Fatal",
				Message:  "Primary stress required",
				Re:       primaryStressRe,
			},
			rules.RequiredTransRe{
				NameStr:  "syllabic",
				LevelStr: "Format",
				Message:  "Each syllable needs a syllabic phoneme",
				Re:       syllabicRe,
			},
			rules.IllegalTransRe{
				NameStr:  "repeated_phonemes",
				LevelStr: "Fatal",
				Message:  "Repeated phonemes cannot be used within the same morpheme",
				Re:       repeatedPhnRe,
			},
			decomp2Orth,
			rules.SymbolSetRule{
				SymbolSet: ss,
			},
		}}
	return v
}

func vInsertEntries(t *testing.T, lexName string) (*sql.DB, string) {

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

	_, err = execSchema(db) // Creates new lexicon database
	ff("Failed to create lexicon db: %v", err)

	l := lexicon{name: lexName, symbolSetName: "ZZ"}

	l, err = defineLexicon(db, l)
	if err != nil {
		t.Errorf(vfs, nil, err)
	}

	lxs, err := listLexicons(db)
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

	_, errx := insertEntries(db, l, []lex.Entry{e1, e2, e3, e4})
	if errx != nil {
		t.Errorf(vfs, "nil", errx)
	}
	return db, l.name
}

func Test_Validation1(t *testing.T) {
	db, lexName := vInsertEntries(t, "test1")
	v := createValidator()

	q := Query{}

	stats, err := Validate(db, []lex.LexName{lex.LexName("test1")}, SilentLogger{}, v, q)
	ff("validation failed : %v", err)

	expect := ValStats{
		TotalEntries:     4,
		ValidatedEntries: 4,
		TotalValidations: 5,
		InvalidEntries:   3,
		Levels: map[string]int{
			"fatal":  3,
			"format": 2,
		},
		Rules: map[string]int{
			"symbolset (fatal)":      2,
			"primary_stress (fatal)": 1,
			"syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}

	// check stats saved in db
	lexStats, err := validationStats(db, lexName)
	ff("validation stats failed : %v", err)

	expectFull := ValStats{
		TotalEntries:     4,
		ValidatedEntries: 4,
		TotalValidations: 5,
		InvalidEntries:   3,
		Levels: map[string]int{
			"fatal":  3,
			"format": 2,
		},
		Rules: map[string]int{
			"symbolset (fatal)":      2,
			"primary_stress (fatal)": 1,
			"syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expectFull, lexStats) {
		t.Errorf(vfs, expectFull, lexStats)
	}
}

func Test_Validation2(t *testing.T) {
	db, lexName := vInsertEntries(t, "test2")
	v := createValidator()

	q := Query{WordRegexp: "a$"}

	// test 1
	stats, err := Validate(db, []lex.LexName{lex.LexName("test2")}, SilentLogger{}, v, q)
	ff("validation failed : %v", err)

	expect := ValStats{
		TotalEntries:     1,
		ValidatedEntries: 1,
		TotalValidations: 0,
		InvalidEntries:   0,
		Levels:           make(map[string]int),
		Rules:            make(map[string]int),
	}
	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}

	// check stats saved in db
	lexStats, err := validationStats(db, lexName)
	ff("validation stats failed : %v", err)

	expectFull := ValStats{
		TotalEntries:     4,
		ValidatedEntries: 4,
		TotalValidations: 0,
		InvalidEntries:   0,
		Levels:           make(map[string]int),
		Rules:            make(map[string]int),
	}

	if !reflect.DeepEqual(expectFull, lexStats) {
		t.Errorf(vfs, expectFull, lexStats)
	}

	// test 2
	q = Query{}

	stats, err = Validate(db, []lex.LexName{lex.LexName("test2")}, SilentLogger{}, v, q)
	ff("validation failed : %v", err)

	expect = ValStats{
		TotalEntries:     4,
		ValidatedEntries: 4,
		TotalValidations: 5,
		InvalidEntries:   3,
		Levels: map[string]int{
			"fatal":  3,
			"format": 2,
		},
		Rules: map[string]int{
			"symbolset (fatal)":      2,
			"primary_stress (fatal)": 1,
			"syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}

	// check stats saved in db
	lexStats, err = validationStats(db, lexName)
	ff("validation stats failed : %v", err)

	expectFull = ValStats{
		TotalEntries:     4,
		ValidatedEntries: 4,
		TotalValidations: 5,
		InvalidEntries:   3,
		Levels: map[string]int{
			"fatal":  3,
			"format": 2,
		},
		Rules: map[string]int{
			"symbolset (fatal)":      2,
			"primary_stress (fatal)": 1,
			"syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expectFull, lexStats) {
		t.Errorf(vfs, expectFull, lexStats)
	}
}

func Test_ValidationUpdate1(t *testing.T) {
	db, lexName := vInsertEntries(t, "test3")
	v := createValidator()
	ew := lex.EntrySliceWriter{}
	err := lookUp(db, []lex.LexName{lex.LexName("test3")}, Query{}, &ew)
	ff("lookup failed : %v", err)

	for _, e := range ew.Entries {
		v.ValidateEntry(&e)
		err = updateValidation(db, []lex.Entry{e})
		ff("update validation failed : %v", err)
	}

	stats, err := validationStats(db, lexName)
	ff("validation stats failed : %v", err)

	expect := ValStats{
		TotalEntries:     4,
		ValidatedEntries: 4,
		TotalValidations: 5,
		InvalidEntries:   3,
		Levels: map[string]int{
			"fatal":  3,
			"format": 2,
		},
		Rules: map[string]int{
			"symbolset (fatal)":      2,
			"primary_stress (fatal)": 1,
			"syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}
}

func Test_ValidationUpdate2(t *testing.T) {
	db, lexName := vInsertEntries(t, "test4")
	v := createValidator()
	ew := lex.EntrySliceWriter{}
	err := lookUp(db, []lex.LexName{lex.LexName("test4")}, Query{}, &ew)
	ff("lookup failed : %v", err)

	var es []lex.Entry
	for _, e := range ew.Entries {
		v.ValidateEntry(&e)
		es = append(es, e)
	}
	err = updateValidation(db, es)
	ff("update validation failed : %v", err)

	stats, err := validationStats(db, lexName)
	ff("validation stats failed : %v", err)

	expect := ValStats{
		TotalEntries:     4,
		ValidatedEntries: 4,
		TotalValidations: 5,
		InvalidEntries:   3,
		Levels: map[string]int{
			"fatal":  3,
			"format": 2,
		},
		Rules: map[string]int{
			"symbolset (fatal)":      2,
			"primary_stress (fatal)": 1,
			"syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}
}

func Test_ValidationUpdate3(t *testing.T) {
	db, lexName := vInsertEntries(t, "test5")
	v := createValidator()
	ew := lex.EntrySliceWriter{}
	err := lookUp(db, []lex.LexName{lex.LexName("test5")}, Query{}, &ew)
	ff("lookup failed : %v", err)

	es, _ := v.ValidateEntries(ew.Entries)
	err = updateValidation(db, es)
	ff("update validation failed : %v", err)

	stats, err := validationStats(db, lexName)
	ff("validation stats failed : %v", err)

	expect := ValStats{
		TotalEntries:     4,
		ValidatedEntries: 4,
		TotalValidations: 5,
		InvalidEntries:   3,
		Levels: map[string]int{
			"fatal":  3,
			"format": 2,
		},
		Rules: map[string]int{
			"symbolset (fatal)":      2,
			"primary_stress (fatal)": 1,
			"syllabic (format)":      2,
		},
	}

	if !reflect.DeepEqual(expect, stats) {
		t.Errorf(vfs, expect, stats)
	}
}
