package main

import (
	//"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/stts-se/pronlex/decompounder"
)

type decomperMutex struct {
	// map from language name to decompounder, as read from word parts file dir.
	// lang is used as HTTP request parameter to select a decompounder
	decompers map[string]decompounder.Decompounder
	// map from language name to word parts file name path.  Used
	// for appending new word parts to the original text file, to
	// keep the word parts text file in sync with the in-memory
	// Decompounder. (Maybe there is a saner way to handle this?)
	files map[string]string
	mutex *sync.RWMutex
}

var decomper = decomperMutex{
	decompers: make(map[string]decompounder.Decompounder),
	files:     make(map[string]string),
	mutex:     &sync.RWMutex{},
}

// appendToWordPartsFile writes a line to a file.
// NB that it is not thread-safe, and should be called after locking.
func appendToWordPartsFile(fn string, line string) error {

	fh, err := os.OpenFile(fn, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer fh.Close()

	_, err = fh.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return nil
}

func addPrefix(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	if "" == lang {
		msg := "no value for the expected 'lang' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	prefix := strings.ToLower(r.FormValue("prefix"))
	if "" == prefix {
		msg := "no value for the expected 'prefix' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	decomper.mutex.Lock()
	defer decomper.mutex.Unlock()
	fn, ok := decomper.files[lang]
	if !ok {
		msg := "unknown 'lang': " + lang
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if decomper.decompers[lang].ContainsPrefix(prefix) {
		fmt.Fprintf(w, "prefix already found: '%s'", prefix)
		return
	}

	decomper.decompers[lang].AddPrefix(prefix)
	err := appendToWordPartsFile(fn, "PREFIX:"+prefix)
	if err != nil {
		msg := fmt.Sprintf("decompounder: failed to append to word parts file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "added '%s'", prefix)
}

// TODO cut-and-paste of addPrefix
func removePrefix(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	if "" == lang {
		msg := "no value for the expected 'lang' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	prefix := strings.ToLower(r.FormValue("prefix"))
	if "" == prefix {
		msg := "no value for the expected 'prefix' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	decomper.mutex.Lock()
	defer decomper.mutex.Unlock()
	fn, ok := decomper.files[lang]
	if !ok {
		msg := "unknown 'lang': " + lang
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	//writeToWordPartsFile("PREFIX")
	if !decomper.decompers[lang].ContainsPrefix(prefix) {
		fmt.Fprintf(w, "prefix not found: '%s'", prefix)
		return
	}
	decomper.decompers[lang].RemovePrefix(prefix)
	err := appendToWordPartsFile(fn, "REMOVE:PREFIX:"+prefix)
	if err != nil {
		msg := fmt.Sprintf("decompounder: failed to append to word parts file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "removed prefix: '%s'", prefix)
}

func addSuffix(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	if "" == lang {
		msg := "no value for the expected 'lang' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	suffix := r.FormValue("suffix")
	if "" == suffix {
		msg := "no value for the expected 'suffix' parameter"
		log.Println()
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	decomper.mutex.Lock()
	defer decomper.mutex.Unlock()
	fn, ok := decomper.files[lang]
	if !ok {
		msg := "unknown 'lang': " + lang
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if decomper.decompers[lang].ContainsSuffix(suffix) {
		fmt.Fprintf(w, "sufffix already found: '%s'", suffix)
		return
	}

	decomper.decompers[lang].AddSuffix(suffix)
	err := appendToWordPartsFile(fn, "SUFFIX:"+suffix)
	if err != nil {
		msg := fmt.Sprintf("decompounder: failed to append to word parts file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "added '%s'", suffix)
}

// TODO cut-and-paste of addSuffix
func removeSuffix(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	if "" == lang {
		msg := "no value for the expected 'lang' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	suffix := r.FormValue("suffix")
	if "" == suffix {
		msg := "no value for the expected 'suffix' parameter"
		log.Println()
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	decomper.mutex.Lock()
	defer decomper.mutex.Unlock()
	fn, ok := decomper.files[lang]
	if !ok {
		msg := "unknown 'lang': " + lang
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if !decomper.decompers[lang].ContainsSuffix(suffix) {
		fmt.Fprintf(w, "suffix not found: '%s'", suffix)
		return
	}
	decomper.decompers[lang].RemoveSuffix(suffix)
	err := appendToWordPartsFile(fn, "REMOVE:SUFFIX:"+suffix)
	if err != nil {
		msg := fmt.Sprintf("decompounder: failed to append to word parts file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "removed '%s'", suffix)
}

type Decomp struct {
	Parts []string `json:"parts"`
}

// langFromFilePath returns the base file name stripped from any '.txt' extension
func langFromFilePath(p string) string {
	b := filepath.Base(p)
	if strings.HasSuffix(b, ".txt") {
		b = b[0 : len(b)-4]
	}
	return b
}

func decompWord(w http.ResponseWriter, r *http.Request) {

	lang := r.FormValue("lang")
	if "" == lang {
		msg := "no value for the expected 'lang' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	word := r.FormValue("word")
	word = strings.ToLower(word)
	if "" == word {
		msg := "no value for the expected 'word' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var res []Decomp
	decomper.mutex.RLock()
	defer decomper.mutex.RUnlock()
	_, ok := decomper.files[lang]
	if !ok {
		msg := "unknown 'lang': " + lang
		var langs []string
		for l, _ := range decomper.decompers {
			langs = append(langs, l)
		}
		msg = fmt.Sprintf("%s. Known 'lang' values: %s", msg, strings.Join(langs, ", "))
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	for _, d := range decomper.decompers[lang].Decomp(word) {
		res = append(res, Decomp{Parts: d})
	}
	log.Println(res)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	j, err := json.Marshal(res)
	if err != nil {
		msg := fmt.Sprintf("failed json marshalling : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(j))
}

func listLanguages(w http.ResponseWriter, r *http.Request) {

	decomper.mutex.RLock()
	var res []string // res0 contains path to file
	for l, _ := range decomper.decompers {
		res = append(res, l)
	}
	decomper.mutex.RUnlock()

	sort.Strings(res)
	j, err := json.Marshal(res)
	if err != nil {
		msg := fmt.Sprintf("failed json marshalling : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, string(j))
}

func decompMain(w http.ResponseWriter, r *http.Request) {
	// TODO error if file not found
	http.ServeFile(w, r, "./src/decomp_demo.html")
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "decompserver <DECOMPFILES DIR>\n")
		os.Exit(0)
	}

	// word decomp file dir. Each file in dn with .txt extension
	// is treated as a word parts file
	var dn = os.Args[1]

	files, err := ioutil.ReadDir(dn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(0)
	}

	// populate map of decompounders from word parts files.
	// The base file name minus '.txt' is the language name.
	var fn string
	for _, f := range files {
		fn = filepath.Join(dn, f.Name())
		if !strings.HasSuffix(fn, ".txt") {
			fmt.Fprintf(os.Stderr, "decompserver: skipping file: '%s'\n", fn)
			continue
		}

		dc, err := decompounder.NewDecompounderFromFile(fn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			fmt.Fprintf(os.Stderr, "decompserver: skipping file: '%s'\n", fn)
			continue
		}

		lang := langFromFilePath(fn)
		decomper.mutex.Lock()
		decomper.decompers[lang] = dc
		decomper.files[lang] = fn
		decomper.mutex.Unlock()
		fmt.Fprintf(os.Stderr, "decomper: loaded file '%s'\n", fn)
	}

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/decomp", decompMain).Methods("get")
	r.HandleFunc("/decomp/decomp", decompWord).Methods("get", "post")
	r.HandleFunc("/decomp/add_prefix", addPrefix).Methods("get", "post")
	r.HandleFunc("/decomp/remove_prefix", removePrefix).Methods("get", "post")
	r.HandleFunc("/decomp/add_suffix", addSuffix).Methods("get", "post")
	r.HandleFunc("/decomp/remove_suffix", removeSuffix).Methods("get", "post")
	r.HandleFunc("/decomp/list_languages", listLanguages).Methods("get", "post")

	r0 := http.StripPrefix("/decomp/built/", http.FileServer(http.Dir("./built/")))
	r.PathPrefix("/decomp/built/").Handler(r0)

	r1 := http.StripPrefix("/decomp/externals/", http.FileServer(http.Dir("./externals/")))
	r.PathPrefix("/decomp/externals/").Handler(r1)

	port := ":6778"
	log.Printf("starting decomp server at port %s\n", port)
	err = http.ListenAndServe(port, r)
	if err != nil {

		log.Fatalf("no fun: %v\n", err)
	}

}
