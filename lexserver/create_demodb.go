package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

func demoEntries() []lex.Entry {
	entries := []lex.Entry{}
	var e lex.Entry
	var t1 lex.Transcription
	var t2 lex.Transcription

	//
	t1 = lex.Transcription{Strn: "\" k e k s", Language: "sv"}
	t2 = lex.Transcription{Strn: "\" C e k s", Language: "sv"}

	e = lex.Entry{Strn: "kex",
		PartOfSpeech:   "NN",
		Morphology:     "NEU IND SIN",
		WordParts:      "kex",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1, t2},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "kex"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\" k e k . s @ t", Language: "sv"}
	t2 = lex.Transcription{Strn: "\" C e k . s @ t", Language: "sv"}
	e = lex.Entry{Strn: "kexet",
		PartOfSpeech:   "NN",
		Morphology:     "NEU DEF SIN",
		WordParts:      "kexet",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1, t2},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "kex"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\"\" k e k + p a . k % e: t", Language: "sv"}
	t2 = lex.Transcription{Strn: "\"\" C e k + p a . k % e: t", Language: "sv"}
	e = lex.Entry{Strn: "kexpaket",
		PartOfSpeech:   "NN",
		Morphology:     "NEU IND SIN",
		WordParts:      "kex+paket",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1, t2},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "kexpaket"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\" h u0 n d", Language: "sv"}
	e = lex.Entry{Strn: "hund",
		PartOfSpeech:   "NN",
		Morphology:     "NEU IND SIN",
		WordParts:      "hund",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "hund"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\" h E s t", Language: "sv"}
	e = lex.Entry{Strn: "häst",
		PartOfSpeech:   "NN",
		Morphology:     "NEU IND SIN",
		WordParts:      "häst",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "häst"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\" h E . s t a r", Language: "sv"}
	e = lex.Entry{Strn: "hästar",
		PartOfSpeech:   "NN",
		Morphology:     "NEU IND PLU",
		WordParts:      "hästar",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "häst"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\" h { . s t a r", Language: "sv"}
	e = lex.Entry{Strn: "hästar",
		PartOfSpeech:   "NN",
		Morphology:     "NEU IND PLU",
		WordParts:      "hästar",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "häst"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\" d U m", Language: "sv"}
	e = lex.Entry{Strn: "dom",
		PartOfSpeech:   "NN",
		Morphology:     "UTR IND SIN",
		WordParts:      "dom",
		Language:       "sv",
		Preferred:      true,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "dom"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\" d o: m", Language: "sv"}
	e = lex.Entry{Strn: "dom",
		PartOfSpeech:   "NN",
		Morphology:     "UTR IND SIN",
		WordParts:      "dom",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "dom"},
	}
	entries = append(entries, e)

	//
	t1 = lex.Transcription{Strn: "\" d O m", Language: "sv"}
	e = lex.Entry{Strn: "dom",
		PartOfSpeech:   "PM",
		Morphology:     "",
		WordParts:      "de",
		Language:       "sv",
		Preferred:      false,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "dom"},
	}
	entries = append(entries, e)

	return entries
}

func setupDemoDB(engine dbapi.DBEngine) error {
	var err error

	log.Println("demo_setup: creating demo database ...")

	dbName := "lexserver_testdb"
	lexRef := lex.NewLexRef(dbName, "sv")
	dbPath := filepath.Join(dbFileArea, dbName+".db")

	dbRelatedPaths, err := filepath.Glob(dbPath + "*")
	if err != nil {
		return fmt.Errorf("failed to retrieve list of db files for '%s' : %v", dbPath, err)
	}
	for _, file := range dbRelatedPaths {
		if _, err = os.Stat(file); !os.IsNotExist(err) {
			log.Printf("demo_setup: deleting db file: %v", file)
			err := os.Remove(file)
			if err != nil {
				return fmt.Errorf("failed to remove %s : %v", file, err)
			}
		}
	}

	dbmx, err := dbapi.NewDBManager(engine)
	if err != nil {
		return fmt.Errorf("failed to init db manager : %v", err)
	}
	defer dbmx.CloseDB(lexRef.DBRef)
	if dbmx.ContainsDB(lexRef.DBRef) {
		err := dbmx.RemoveDB(lexRef.DBRef)
		if err != nil {
			return fmt.Errorf("failed to remove db: %v", err)
		}
	}

	err = dbmx.DefineSqliteDB(lexRef.DBRef, dbPath)
	if err != nil {
		return fmt.Errorf("failed to define db %s | %v : %v", dbPath, lexRef, err)
	}

	err = dbmx.DefineLexicon(lexRef, "sv-se_ws-sampa", "sv")
	if err != nil {
		return fmt.Errorf("failed to create lexicon %v: %v", lexRef, err)
	}

	_, err = dbmx.InsertEntries(lexRef, demoEntries())
	if err != nil {
		return fmt.Errorf("failed to insert entries to db %v: %v", lexRef, err)
	}

	log.Println("demo_setup: test database completed")
	return nil
}
