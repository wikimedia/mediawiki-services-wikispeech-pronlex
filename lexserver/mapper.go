package main

import (
	"fmt"
	"log"
	"net/http"
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
			msg := fmt.Sprintf("failed getting map table : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		mapper := JSONMapper{From: mapper0.SymbolSet1.Name, To: mapper0.SymbolSet2.Name}
		mapper.Symbols = make([]JSONMSymbol, 0)
		for _, from := range mapper0.SymbolSet1.Symbols {
			to, err := mapper0.MapSymbol(from)
			if err != nil {
				msg := fmt.Sprintf("failed getting map table : %v", err)
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
