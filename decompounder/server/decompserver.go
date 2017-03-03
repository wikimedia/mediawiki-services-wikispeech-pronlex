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
	"strings"
	"sync"

	"github.com/gorilla/mux"

	"github.com/stts-se/pronlex/decompounder"
)

// TODO mutex this:
type decomperMutex struct {
	decompounder.Decompounder
	*sync.RWMutex
}

var decomper decomperMutex //= decomperMutex{decompounder.NewDecompounder(), &sync.RWMutex{}}

func addPrefix(w http.ResponseWriter, r *http.Request) {

	prefix := r.FormValue("prefix")
	if "" == prefix {
		msg := "no value for the expected 'prefix' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	decomper.Lock()
	defer decomper.Unlock()
	decomper.AddPrefix(prefix)
	fmt.Fprintf(w, "added '%s'", prefix)
}

func addSuffix(w http.ResponseWriter, r *http.Request) {

	suffix := r.FormValue("suffix")
	if "" == suffix {
		msg := "no value for the expected 'suffix' parameter"
		log.Println()
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	decomper.Lock()
	defer decomper.Unlock()
	decomper.AddSuffix(suffix)
	fmt.Fprintf(w, "added '%s'", suffix)
}

type Decomp struct {
	Parts []string `json:"parts"`
}

func decompWord(w http.ResponseWriter, r *http.Request) {

	word := r.FormValue("word")

	// REMOVE ME:
	if word == "ERROR" {
		http.Error(w, "ERROR! TERROR", http.StatusInternalServerError)
		return
	}

	word = strings.ToLower(word)
	if "" == word {
		msg := "no value for the expected 'word' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	var res []Decomp
	//res := fmt.Sprintf("%#v", decomper.Decomp(word))
	decomper.RLock()
	defer decomper.RUnlock()

	for _, d := range decomper.Decomp(word) {
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

func decompMain(w http.ResponseWriter, r *http.Request) {
	// TODO error if file not found
	http.ServeFile(w, r, "./src/decomp_demo.html")
}

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "decompserver <DECOMPFILES DIR>\n")
		os.Exit(0)
	}
	var dn = os.Args[1] //"decomps.txt"

	files, err := ioutil.ReadDir(dn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(0)
	}

	// TODO create map of Decomponders and create one instance for each file in decomp file dir

	var fn string
	for _, f := range files {
		fn = filepath.Join(dn, f.Name())
		//fmt.Println(file.Name())
	}

	dc, err := decompounder.NewDecompounderFromFile(fn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	decomper = decomperMutex{dc, &sync.RWMutex{}}

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/decomp", decompMain).Methods("get")
	r.HandleFunc("/decomp/decomp", decompWord).Methods("get", "post")
	r.HandleFunc("/decomp/add_prefix", addPrefix).Methods("get", "post")
	r.HandleFunc("/decomp/add_suffix", addSuffix).Methods("get", "post")

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
