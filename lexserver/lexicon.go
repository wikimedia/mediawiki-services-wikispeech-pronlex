package main

// The handlers of calls prefixed with '/lexicon/':

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

func listLexsHandler(w http.ResponseWriter, r *http.Request) {

	lexs, err := dbapi.ListLexicons(db) // TODO error handling
	if err != nil {
		http.Error(w, fmt.Sprintf("list lexicons failed : %v", err), http.StatusInternalServerError)
		return
	}
	jsn, err := marshal(lexs, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(jsn))
}

type LexWithEntryCount struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	SymbolSetName string `json:"symbolSetName"`
	EntryCount    int64  `json:"entryCount"`
}

func listLexsWithEntryCountHandler(w http.ResponseWriter, r *http.Request) {

	lexs0, err := dbapi.ListLexicons(db) // TODO error handling
	if err != nil {
		http.Error(w, fmt.Sprintf("list lexicons failed : %v", err), http.StatusInternalServerError)
		return
	}
	var lexs []LexWithEntryCount
	for _, lex := range lexs0 {
		entryCount, err := dbapi.EntryCount(db, lex.ID)
		if err != nil {
			http.Error(w, fmt.Sprintf("lexicon stats failed : %v", err), http.StatusInternalServerError)
			return
		}
		lexs = append(lexs, LexWithEntryCount{ID: lex.ID, Name: lex.Name, SymbolSetName: lex.SymbolSetName, EntryCount: entryCount})
	}
	jsn, err := marshal(lexs, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(jsn))
}

func lexiconStatsHandler(w http.ResponseWriter, r *http.Request) {
	lexiconID, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
	if err != nil {
		msg := "lexiconStatsHandler got no lexicon id"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	stats, err := dbapi.LexiconStats(db, lexiconID)
	if err != nil {
		http.Error(w, fmt.Sprintf("lexiconStatsHandler: call to  dbapi.LexiconStats failed : %v", err), http.StatusInternalServerError)
		return
	}
	res, err := json.Marshal(stats)

	if err != nil {
		http.Error(w, fmt.Sprintf("lexiconStatsHandler: failed to marshal struct : %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(res))
}

func lexLookUpHandler(w http.ResponseWriter, r *http.Request) {

	// TODO check r.Method?

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

	q, err := queryFromParams(r)
	if err != nil {
		log.Printf("failed to process query params: %v", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusBadRequest)
		return
	}

	//res, err := dbapi.LookUpIntoMap(db, q) // GetEntries(db, q)
	res, err := dbapi.LookUpIntoSlice(db, q) // GetEntries(db, q)
	if err != nil {
		log.Printf("lexserver: Failed to get entries: %v", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	jsn, err := marshal(res, r)
	if err != nil {
		log.Printf("lexserver: Failed to marshal json: %v", err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(jsn))
}

// TODO add tests
func addEntryHandler(w http.ResponseWriter, r *http.Request) {
	// TODO error check parameters
	lexiconName := r.FormValue("lexicon")
	lexicon, err := dbapi.GetLexicon(db, lexiconName)
	if err != nil {
		msg := fmt.Sprintf("failed to find lexicon %s in database : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	entryJSON := r.FormValue("entry")
	var e lex.Entry
	err = json.Unmarshal([]byte(entryJSON), &e)
	if err != nil {
		log.Printf("lexserver: Failed to unmarshal json: %v", err)
		http.Error(w, fmt.Sprintf("failed to process incoming Entry json : %v", err), http.StatusInternalServerError)
		return
	}

	ids, err := dbapi.InsertEntries(db, lexicon, []lex.Entry{e})
	if err != nil {
		msg := fmt.Sprintf("lexserver failed to update entry : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, ids)
}

func insertOrUpdateLexHandler(w http.ResponseWriter, r *http.Request) {
	// if no id or not an int, simply set id to 0:
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	name := strings.TrimSpace(r.FormValue("name"))
	symbolSetName := strings.TrimSpace(r.FormValue("symbolsetname"))

	if name == "" || symbolSetName == "" {
		msg := fmt.Sprint("missing parameter value, expecting value for 'name' and 'symbolsetname'")
		log.Printf("%s", msg)
		http.Error(w, msg, http.StatusExpectationFailed)
		return
	}

	res, err := dbapi.InsertOrUpdateLexicon(db, dbapi.Lexicon{ID: id, Name: name, SymbolSetName: symbolSetName})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed database call : %v", err), http.StatusInternalServerError)
		return
	}

	jsn, err := marshal(res, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(jsn))
	//fmt.Fprint(w, jsn)
}

func lexiconRunValidateHandler(w http.ResponseWriter, r *http.Request) {
}
