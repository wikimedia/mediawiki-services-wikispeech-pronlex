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
	"time"

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

func listCurrentEntryStatuses(w http.ResponseWriter, r *http.Request) {

	lexiconName := r.FormValue("lexicon_name")
	if "" == lexiconName {
		http.Error(w, "missing value for lexicon_name param", http.StatusBadRequest)
		return
	}

	statuses, err := dbapi.ListCurrentEntryStatuses(db, lexiconName)
	if err != nil {
		http.Error(w, fmt.Sprintf("listCurrentEntryStatuses : %v", err), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(statuses)
	if err != nil {
		http.Error(w, fmt.Sprintf("listCurrentEntryStatuses : %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(j))
}

// TODO cut-n-paste from above
func listAllEntryStatuses(w http.ResponseWriter, r *http.Request) {

	lexiconName := r.FormValue("lexicon_name")
	if "" == lexiconName {
		http.Error(w, "missing value for lexicon_name param", http.StatusBadRequest)
		return
	}

	statuses, err := dbapi.ListAllEntryStatuses(db, lexiconName)
	if err != nil {
		http.Error(w, fmt.Sprintf("listAllEntryStatuses : %v", err), http.StatusInternalServerError)
		return
	}
	j, err := json.Marshal(statuses)
	if err != nil {
		http.Error(w, fmt.Sprintf("listAllEntryStatuses : %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(j))
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

func lexInfoHandler(w http.ResponseWriter, r *http.Request) {
	lexName := r.FormValue("name")
	if len(strings.TrimSpace(lexName)) == 0 {
		msg := fmt.Sprintf("lexicon name should be specified by variable 'name'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	lex, err := dbapi.GetLexicon(db, lexName) // TODO error handling
	if err != nil {
		http.Error(w, fmt.Sprintf("get lexicon failed : %v", err), http.StatusInternalServerError)
		return
	}
	jsn, err := marshal(lex, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(jsn))
}

func lexiconStatsHandler(w http.ResponseWriter, r *http.Request) {
	lexiconID, err := strconv.ParseInt(r.FormValue("lexiconId"), 10, 64)
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
		msg := fmt.Sprintf("failed to find lexicon %s in database : %v", lexiconName, err)
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
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("lexiconfileupload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
		return
	}

	start := time.Now()

	clientUUID := r.FormValue("client_uuid")

	if "" == strings.TrimSpace(clientUUID) {
		msg := "lexiconRunValidateHandler got no client uuid"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	conn, ok := webSocks.clients[clientUUID]
	if !ok {
		msg := fmt.Sprintf("lexiconRunValidateHandler couldn't find connection for uuid %v", clientUUID)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	logger := dbapi.NewWebSockLogger(conn)

	lexName := r.FormValue("lexicon_name")
	lexicon, err := dbapi.GetLexicon(db, lexName)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't retrive lexicon : %v", err), http.StatusInternalServerError)
		return
	}
	vMut.Lock()
	v, err := vMut.service.ValidatorForName(lexicon.SymbolSetName)
	vMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("lexiconRunValidateHandler failed to get validator for symbol set %v : %v", lexicon.SymbolSetName, err)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	q := dbapi.Query{Lexicons: []dbapi.Lexicon{lexicon}}
	stats, err := dbapi.Validate(db, logger, *v, q)
	if err != nil {
		msg := fmt.Sprintf("lexiconRunValidateHandler failed validate : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	dur := round(time.Since(start), time.Second)
	fmt.Fprintf(w, "\nDuration %v\n", dur)
	fmt.Fprint(w, stats)
}

func round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}

func lexiconHelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	addEntryURL := "/lexicon/addentry?lexicon=&entry={...}"

	updateEntryURL := "/lexicon/updateentry?entry={...}"

	html := `<h1>Lexicon</h1>
<h2>addentry</h2>Add an entry to the database. Example invocation:
<pre><a href="` + addEntryURL + `">` + addEntryURL + `</a></pre>

<h2>updateentry</h2> Updates an entry in the database. Example invocation:
<pre><a href="` + updateEntryURL + `">` + updateEntryURL + `</a></pre>

<h2>validate</h2> Validates a list of entries.
<pre><a href="/lexicon/validate">/lexicon/validate</a></pre>

<h2>list</h2> Lists available lexicons.
<pre><a href="/lexicon/list">/lexicon/list</a></pre>

<h2>list</h2> Display lexicon info.
<pre><a href="/lexicon/info?name=sv-se.nst">/lexicon/info?name=sv-se.nst</a></pre>

<h2>stats</h2> Lists lexicon stats. Example invocation:
<pre><a href="/lexicon/stats?lexiconId=1">/lexicon/stats?lexiconId=1</a></pre>
		`

	fmt.Fprint(w, html)
}
