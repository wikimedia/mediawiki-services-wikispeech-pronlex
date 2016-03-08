package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stts-se/pronlex/dbapi"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func f(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func ff(f string, err error) {
	if err != nil {
		log.Fatalf(f, err)
	}
}

// pretty print if the URL paramer 'pp' has a value
func marshal(v interface{}, r *http.Request) ([]byte, error) {

	if "" != strings.TrimSpace(r.FormValue("pp")) {
		return json.MarshalIndent(v, "", "  ")
	}

	return json.Marshal(v)
}

func ipaTableHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/ipa_table.txt")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

func listLexsHandler(w http.ResponseWriter, r *http.Request) {

	lexs, err := dbapi.ListLexicons(db) // TODO error handling
	jsn, err := marshal(lexs, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/javascript") // TODO Behövs denna?
	fmt.Fprint(w, string(jsn))
}

func insertOrUpdateLexHandler(w http.ResponseWriter, r *http.Request) {
	// if no id or not an int, simply set id to 0:
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	name := strings.TrimSpace(r.FormValue("name"))
	symbolSetName := strings.TrimSpace(r.FormValue("symbolsetname"))
	if name == "" || symbolSetName == "" {
		http.Error(w, fmt.Sprintf("missing parameter value, expecting value for 'name' and 'symbolsetname'"), http.StatusExpectationFailed)
		return
	}

	res, err := dbapi.InsertOrUpdateLexicon(db, dbapi.Lexicon{Id: id, Name: name, SymbolSetName: symbolSetName})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed database call : %v", err), http.StatusInternalServerError)
		return
	}

	jsn, err := marshal(res, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(jsn))
}

func deleteLexHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	//name := strings.TrimSpace(r.FormValue("name"))
	//symbolSetName := strings.TrimSpace(r.FormValue("symbolsetname"))

	err := dbapi.DeleteLexicon(db, id) //dbapi.Lexicon{Id: id, Name: name, SymbolSetName: symbolSetName})
	if err != nil {

		http.Error(w, fmt.Sprintf("failed deleting lexicon : %v", err), http.StatusInternalServerError)
		return
	}

}

// TODO report unused URL parameters

// TODO Gör konstanter som kan användas istället för strängar
var knownParams = map[string]int{
	"lexicons":          1,
	"words":             1,
	"lemmas":            1,
	"wordlike":          1,
	"transcriptionlike": 1,
	"partofspeechlike":  1,
	"lemmalike":         1,
	"readinglike":       1,
	"paradigmlike":      1,
	"page":              1,
	"pagelike":          1,
	"pp":                1,
}

// TODO return error
func queryFromParams(r *http.Request) dbapi.Query {

	lexs := dbapi.RemoveEmptyStrings(
		regexp.MustCompile("[, ]").Split(r.FormValue("lexicons"), -1))
	words := dbapi.RemoveEmptyStrings(
		regexp.MustCompile("[, ]").Split(r.FormValue("words"), -1))
	lemmas := dbapi.RemoveEmptyStrings(
		regexp.MustCompile("[, ]").Split(r.FormValue("lemmas"), -1))

	wordLike := strings.TrimSpace(r.FormValue("wordlike"))
	transcriptionLike := strings.TrimSpace(r.FormValue("transcriptionlike"))
	partOfSpeechLike := strings.TrimSpace(r.FormValue("partofspeechlike"))
	lemmaLike := strings.TrimSpace(r.FormValue("lemmalike"))
	readingLike := strings.TrimSpace(r.FormValue("readinglike"))
	paradigmLike := strings.TrimSpace(r.FormValue("paradigmlike"))

	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		page = 0
	}
	pageLength, err := strconv.ParseInt(r.FormValue("pagelength"), 10, 64)
	if err != nil {
		pageLength = 25
	}

	dbLexs, err := dbapi.GetLexicons(db, lexs)
	// TODO return error
	if err != nil {
		log.Printf("GetLexicons failed : %v", err)
	}
	q := dbapi.Query{
		Lexicons:          dbLexs,
		Words:             words,
		WordLike:          wordLike,
		TranscriptionLike: transcriptionLike,
		PartOfSpeechLike:  partOfSpeechLike,
		Lemmas:            lemmas,
		LemmaLike:         lemmaLike,
		ReadingLike:       readingLike,
		ParadigmLike:      paradigmLike,
		Page:              page,
		PageLength:        pageLength}

	return q
}

func lexLookUpHandler(w http.ResponseWriter, r *http.Request) {

	var err error
	// TODO Felhantering?

	// TODO report unknown params to client
	u, err := url.Parse(r.URL.String())
	ff("lexLookUpHandler failed to get params: %v", err)
	params := u.Query()
	if len(params) == 0 {
		log.Print("lexLookUpHandler: zero params, serving lexlookup.html")
		http.ServeFile(w, r, "./static/lexlookup.html")
	}
	for k, v := range params {
		if _, ok := knownParams[k]; !ok {
			log.Printf("lexLookUpHandler: unknown URL parameter: '%s': '%s'", k, v)
		}
	}

	q := queryFromParams(r)
	res, err := dbapi.GetEntries(db, q)
	if err != nil {
		log.Printf("lexserver: Failed to get entries: %v", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}

	jsn, err := marshal(res, r)
	if err != nil {
		log.Printf("lexserver: Failed to marshal json: %v", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/javascript") // TODO Behövs denna?
	fmt.Fprint(w, string(jsn))                               // TODO should not be called if error occurs?
}

// func updateEntryHandler(w http.ResponseWriter, r *http.Request) {

// }

func adminCreateLexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/admin/create_lex.html")
}

var db *sql.DB

func main() {

	port := ":8787"

	dbFile := "./pronlex.db"

	var err error // återanvänds för alla fel

	// kolla att db-filen existerar
	_, err = os.Stat(dbFile)
	ff("lexserver: Cannot find db file. %v", err)

	db, err = sql.Open("sqlite3", dbFile)
	ff("Failed to open dbfile %v", err)
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	ff("Failed to exec PRAGMA call %v", err)
	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
	ff("Failed to exec PRAGMA call %v", err)

	// static
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/ipa_table.txt", ipaTableHandler)

	// function calls
	http.HandleFunc("/listlexicons", listLexsHandler)
	http.HandleFunc("/lexlookup", lexLookUpHandler)
	//http.HandleFunc("/updateentry", updateEntryHandler)

	// admin page
	http.HandleFunc("/admin/createlex", adminCreateLexHandler)
	http.HandleFunc("/admin/insertorupdatelexicon", insertOrUpdateLexHandler)
	http.HandleFunc("/admin/deletelexicon", deleteLexHandler)

	//            (Why this http.StripPrefix?)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Print("Lexicon server listening on port ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
