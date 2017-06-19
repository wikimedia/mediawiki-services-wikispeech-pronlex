package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"

	"github.com/stts-se/pronlex/symbolset"
	//"os"
	"encoding/json"
	"strings"
)

// The calls prefixed with '/mapper/'

var mMut = struct {
	sync.RWMutex
	service symbolset.MapperService
}{
	service: symbolset.MapperService{
		SymbolSets: make(map[string]symbolset.SymbolSet),
		Mappers:    make(map[string]symbolset.Mapper),
	},
}

type JSONMapped struct {
	From   string
	To     string
	Input  string
	Result string
}

func trimTrans(trans string) string {
	re := "  +"
	repl := regexp.MustCompile(re)
	trans = repl.ReplaceAllString(trans, " ")
	return trans
}

var mapperMap = urlHandler{
	name:     "map",
	url:      "/map/{from}/{to}/{trans}",
	help:     "Maps a transcription from one symbolset to another.",
	examples: []string{"/map/sv-se_ws-sampa/sv-se_sampa_mary/%22%22%20p%20O%20j%20.%20k%20@"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		fromName := getParam("from", r)
		toName := getParam("to", r)
		trans := trimTrans(getParam("trans", r))
		if len(strings.TrimSpace(fromName)) == 0 {
			msg := fmt.Sprintf("input symbol set should be specified by variable 'from'")
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if len(strings.TrimSpace(toName)) == 0 {
			msg := fmt.Sprintf("output symbol set should be specified by variable 'to'")
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if len(strings.TrimSpace(trans)) == 0 {
			msg := fmt.Sprintf("input trans should be specified by variable 'trans'")
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		mMut.Lock()
		result0, err := mMut.service.Map(fromName, toName, trans)
		mMut.Unlock()
		if err != nil {
			msg := fmt.Sprintf("failed mapping transcription : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		result := JSONMapped{Input: trans, Result: result0, From: fromName, To: toName}
		j, err := json.Marshal(result)
		if err != nil {
			msg := fmt.Sprintf("json marshalling error : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(j))
	},
}

type JSONMapper struct {
	From    string
	To      string
	Symbols []JSONMSymbol
}

type JSONMSymbol struct {
	From string
	To   string
	IPA  JSONIPA
	Desc string
	Cat  string
}

var mapperMaptable = urlHandler{
	name:     "maptable",
	url:      "/maptable/{from}/{to}",
	help:     "Lists content of a maptable given two symbolset names.",
	examples: []string{"/maptable/sv-se_ws-sampa/sv-se_sampa_mary"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		fromName := getParam("from", r)
		toName := getParam("to", r)
		if len(strings.TrimSpace(fromName)) == 0 {
			msg := fmt.Sprintf("input symbol set should be specified by variable 'from'")
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if len(strings.TrimSpace(toName)) == 0 {
			msg := fmt.Sprintf("output symbol set should be specified by variable 'to'")
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		mMut.Lock()
		mapper0, err := mMut.service.GetMapTable(fromName, toName)
		mMut.Unlock()
		if err != nil {
			msg := fmt.Sprintf("failed getting map table: %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		mapper := JSONMapper{From: mapper0.SymbolSet1.Name, To: mapper0.SymbolSet2.Name}
		mapper.Symbols = make([]JSONMSymbol, 0)
		for _, from := range mapper0.SymbolSet1.Symbols {
			to, err := mapper0.MapSymbol(from)
			if err != nil {
				msg := fmt.Sprintf("failed getting map table: %v", err)
				log.Println(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
			mapper.Symbols = append(mapper.Symbols, JSONMSymbol{From: from.String, To: to.String, IPA: JSONIPA{String: from.IPA.String, Unicode: from.IPA.Unicode}, Desc: from.Desc, Cat: from.Cat.String()})
		}

		j, err := json.Marshal(mapper)
		if err != nil {
			msg := fmt.Sprintf("json marshalling error : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(j))
	},
}

var mapperList = urlHandler{
	name:     "list",
	url:      "/list",
	help:     "Lists cached mappers.",
	examples: []string{"/list"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		mMut.Lock()
		ms := mMut.service.MapperNames()
		mMut.Unlock()
		j, err := json.Marshal(ms)
		if err != nil {
			msg := fmt.Sprintf("failed to marshal struct : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, string(j))
	},
}

/// MAPPER INIT TESTS

type mapperTests struct {
	fromName string
	toName   string
	tests    []mapperTest
}

type mapperTest struct {
	from string
	to   string
}

func parseMapperTestLine(l string) (mapperTest, error) {
	fs := strings.Split(l, "\t")
	if fs[0] != "TEST" {
		return mapperTest{}, fmt.Errorf("mapper test line must start with TEST; found %s", l)
	}
	if len(fs) != 3 {
		return mapperTest{}, fmt.Errorf("mapper test line must have 3 fields, found %s", l)
	}
	from := fs[1]
	to := fs[2]
	return mapperTest{
		from: from,
		to:   to}, nil
}

func parseMapperTest(l string) (mapperTest, error) {
	fs := strings.Split(l, "\t")
	if fs[0] != "TEST" {
		return mapperTest{}, fmt.Errorf("mapper test line must start with TEST; found %s", l)
	}
	if len(fs) != 3 {
		return mapperTest{}, fmt.Errorf("mapper test line must have 3 fields, found %s", l)
	}
	from := fs[1]
	to := fs[2]
	return mapperTest{
		from: from,
		to:   to}, nil

}

func parseMapperTests(mapperLine string, testLines []string) (mapperTests, error) {
	fs := strings.Split(mapperLine, "\t")
	if fs[0] != "MAPPER" {
		return mapperTests{}, fmt.Errorf("mapper definition line must start with MAPPER; found %s", mapperLine)
	}
	if len(fs) != 3 {
		return mapperTests{}, fmt.Errorf("mapper definition line must have 3 fields, found %s", mapperLine)
	}
	from := fs[1]
	to := fs[2]
	tests := []mapperTest{}
	for _, l := range testLines {
		t, err := parseMapperTest(l)
		if err != nil {
			return mapperTests{}, err
		}
		tests = append(tests, t)
	}
	return mapperTests{
		fromName: from,
		toName:   to,
		tests:    tests}, nil
}

func loadMapperTestsFromFile(fName string) ([]mapperTests, error) {
	var res []mapperTests
	fh, err := os.Open(fName)
	defer fh.Close()
	if err != nil {
		return []mapperTests{}, err
	}
	s := bufio.NewScanner(fh)
	n := 0
	prevMapper := ""
	cachedTests := []string{}
	for s.Scan() {
		if err := s.Err(); err != nil {
			return []mapperTests{}, err
		}
		n++
		l := s.Text()
		if len(strings.TrimSpace(l)) == 0 {
			// empty line
		} else if strings.HasPrefix(l, "#") {
			// comment line
		} else if strings.HasPrefix(l, "TEST\t") {
			cachedTests = append(cachedTests, l)
		} else if strings.HasPrefix(l, "MAPPER\t") {
			if prevMapper != "" {
				mTests, err := parseMapperTests(prevMapper, cachedTests)
				if err != nil {
					return []mapperTests{}, err
				}
				res = append(res, mTests)
			}
			cachedTests = []string{}
			prevMapper = l
		}
	}
	if prevMapper != "" {
		mTests, err := parseMapperTests(prevMapper, cachedTests)
		if err != nil {
			return []mapperTests{}, err
		}
		res = append(res, mTests)
	}
	return res, nil
}

func testMappers(mDefFile string) error {
	errs := []string{}
	if _, err := os.Stat(mDefFile); !os.IsNotExist(err) {
		log.Println("lexserver: loading mapper definitions from file", mDefFile)
		mTests, err := loadMapperTestsFromFile(mDefFile)
		if err != nil {
			return nil
		}
		for _, mt := range mTests {
			log.Println("lexserver: initializing mapper", mt)
			mMut.Lock()
			mtab, err := mMut.service.GetMapTable(mt.fromName, mt.toName)
			mMut.Unlock()
			if err != nil {
				msg := fmt.Sprintf("failed getting map table : %v", err)
				log.Println(msg)
				return err
			}
			for _, from := range mtab.SymbolSet1.Symbols {
				_, err := mtab.MapSymbol(from)
				if err != nil {
					mMut.service.DeleteMapper(mt.fromName, mt.toName)
					msg := fmt.Sprintf("failed getting map table : %v", err)
					log.Println(msg)
					return err
				}
			}

			for _, t := range mt.tests {
				mMut.Lock()
				mapped, err := mMut.service.Map(mt.fromName, mt.toName, t.from)
				mMut.Unlock()
				if err != nil {
					return err
				}
				if mapped != t.to {
					msg := fmt.Sprintf("from /%s/ expected /%s/, found /%s/", t.from, t.to, mapped)
					log.Println(msg)
					errs = append(errs, msg)
				}
			}
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "; "))
	}
	return nil
}
