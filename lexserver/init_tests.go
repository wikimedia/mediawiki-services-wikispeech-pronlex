package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

func serverInitTests() error {
	port := "8799"
	s, err := createServer(port)
	if err != nil {
		return err
	}

	_, err = net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	s.ListenAndServe()

	err = runTests(port)
	if err != nil {
		return err
	}

	s.Close()

	log.Println("init_tests: done")

	return nil
}

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
		WordParts:      "kexet",
		Language:       "sv",
		Preferred:      true,
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
		Preferred:      true,
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
		Preferred:      true,
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
		Preferred:      true,
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
		Preferred:      true,
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
		Preferred:      true,
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
		Preferred:      true,
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
		Preferred:      true,
		Transcriptions: []lex.Transcription{t1},
		EntryStatus:    lex.EntryStatus{Name: "demo", Source: "auto"},
		Lemma:          lex.Lemma{Strn: "dom"},
	}
	entries = append(entries, e)

	return entries
}

func setupDemoDB() error {
	log.Println("demo_setup: creating test database ...")

	dbName := "demodb"
	lexRef := lex.NewLexRef(dbName, "demolex")
	dbPath := filepath.Join(dbFileArea, dbName+".db")

	var dbmx = dbapi.NewDBManager()
	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		log.Printf("demo_setup: deleting demo db: %v", dbPath)
		err := os.Remove(dbPath)
		if err != nil {
			return fmt.Errorf("failed to remove %s : %v", dbPath, err)
		}
	}
	if dbmx.ContainsDB(lexRef.DBRef) {
		err := dbmx.RemoveDB(lexRef.DBRef)
		if err != nil {
			return fmt.Errorf("failed to remove db: %v", err)
		}
	}

	err := dbmx.DefineSqliteDB(lexRef.DBRef, dbPath)
	if err != nil {
		return fmt.Errorf("failed to define db: %v", err)
	}

	err = dbmx.DefineLexicon(lexRef, "sv-se_ws-sampa")
	if err != nil {
		return fmt.Errorf("failed to create lexicon: %v", err)
	}

	_, err = dbmx.InsertEntries(lexRef, demoEntries())
	if err != nil {
		return fmt.Errorf("Failed to insert entries to db: %v", err)
	}

	log.Println("demo_setup: test database completed")
	return nil
}

func runTests(port string) error {
	log.Println("init_tests: running tests ... (work in progress)")
	for _, subRouter := range subRouters {
		for _, handler := range subRouter.handlers {
			for _, example := range handler.examples {
				url := "http://localhost:" + port + subRouter.root + example
				template := subRouter.root + handler.url
				log.Println("init_tests: " + template + " => " + url)
				// resp, err := http.Get(url)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// fmt.Println(resp)
				// defer resp.Body.Close()
				// _, err = io.Copy(os.Stdout, resp.Body)
				// body, err := ioutil.ReadAll(resp.Body)
				// if err != nil {
				// 	log.Fatal(err)
				// }
				// fmt.Println(body)
			}
		}
	}
	log.Println("init_tests: test completed")
	return nil
}
