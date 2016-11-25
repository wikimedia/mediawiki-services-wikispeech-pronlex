package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/stts-se/pronlex/decompounder"
)

// TODO mutex this:
var decomper = decompounder.NewDecompounder()

func addPrefix(w http.ResponseWriter, r *http.Request) {

	prefix := r.FormValue("prefix")
	if "" == prefix {
		msg := "no value for the expected 'prefix' parameter"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

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

	decomper.AddSuffix(suffix)
	fmt.Fprintf(w, "added '%s'", suffix)
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

	res := fmt.Sprintf("%#v", decomper.Decomp(word))
	log.Println(res)

	w.Header().Set("Content-Type", "text/plain;charset=UTF-8")
	fmt.Fprintf(w, res)
}

func decompMain(w http.ResponseWriter, r *http.Request) {
	// TODO error if file not found
	http.ServeFile(w, r, "./src/decomp_demo.html")
}

func main() {

	decomper.AddPrefix("bil")
	decomper.AddSuffix("skrot")

	decomper.AddPrefix("skrot")
	decomper.AddSuffix("bil")

	decomper.AddPrefix("last")
	decomper.AddSuffix("båt")

	decomper.AddPrefix("båt")
	decomper.AddSuffix("last")

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/decomp", decompMain).Methods("get")
	r.HandleFunc("/decomp/decomp", decompWord).Methods("get", "post")
	r.HandleFunc("/decomp/add_prefix", addPrefix).Methods("get", "post")
	r.HandleFunc("/decomp/add_suffix", addSuffix).Methods("get", "post")

	r0 := http.StripPrefix("/decomp/built/", http.FileServer(http.Dir("./built/")))
	r.PathPrefix("/decomp/built/").Handler(r0)

	port := ":6778"
	log.Printf("starting decomp server at port %s\n", port)
	err := http.ListenAndServe(port, r)
	if err != nil {

		log.Fatalf("no fun: %v\n", err)
	}

}
