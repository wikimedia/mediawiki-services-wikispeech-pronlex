package main

// The handlers of calls prefixed with '/lexicon/':

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
)

var lexiconValidationPage = urlHandler{
	name:     "validation (page)",
	url:      "/validation_page",
	help:     "Validate lexicon (GUI).",
	examples: []string{"/validation_page"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/lexicon/validation_page.html")
	},
}
var lexiconUpdateEntryURL = "/lexicon/updateentry?entry={...}"

var lexiconUpdateEntry = urlHandler{
	name:     "updateentry",
	url:      "/updateentry",
	help:     " Updates an entry in the database.",
	examples: []string{lexiconUpdateEntryURL},
	handler: func(w http.ResponseWriter, r *http.Request) {
		entryJSON := getParam("entry", r)
		//body, err := ioutil.ReadAll(r.Body)
		var e lex.Entry
		err := json.Unmarshal([]byte(entryJSON), &e)
		if err != nil {
			log.Printf("lexserver: Failed to unmarshal json: %v", err)
			http.Error(w, fmt.Sprintf("failed to process incoming Entry json : %v", err), http.StatusInternalServerError)
			return
		}

		// Underscore below matches bool indicating if any update has taken place. Return this info?
		res, _, err2 := dbm.UpdateEntry(e)
		if err2 != nil {
			log.Printf("lexserver: Failed to update entry : %v", err2)
			http.Error(w, fmt.Sprintf("failed to update Entry : %v", err2), http.StatusInternalServerError)
			return
		}

		res0, err3 := json.Marshal(res)
		if err3 != nil {
			log.Printf("lexserver: Failed to marshal entry : %v", err3)
			http.Error(w, fmt.Sprintf("failed return updated Entry : %v", err3), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, string(res0))
	},
}

// LexWithEntryCount is a struct for collecting lexicon info for json result
type LexWithEntryCount struct {
	Name          string `json:"name"`
	SymbolSetName string `json:"symbolSetName"`
	EntryCount    int64  `json:"entryCount"`
}

var lexiconList = urlHandler{
	name:     "list",
	url:      "/list",
	help:     "Lists available lexicons along with some basic info.",
	examples: []string{"/list"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexs0, err := dbm.ListLexicons() // TODO error handling
		if err != nil {
			http.Error(w, fmt.Sprintf("list lexicons failed : %v", err), http.StatusInternalServerError)
			return
		}
		var lexs []LexWithEntryCount
		for _, lex := range lexs0 {
			entryCount, err := dbm.EntryCount(lex.LexRef)
			if err != nil {
				http.Error(w, fmt.Sprintf("lexicon stats failed : %v", err), http.StatusInternalServerError)
				return
			}
			lexs = append(lexs, LexWithEntryCount{Name: lex.LexRef.String(), SymbolSetName: lex.SymbolSetName, EntryCount: entryCount})
		}
		jsn, err := marshal(lexs, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, string(jsn))
	},
}

var lexiconListCurrentEntryStatuses = urlHandler{
	name:     "list_current_entry_statuses",
	url:      "/list_current_entry_statuses/{lexicon_name}",
	help:     "List current entry statuses.",
	examples: []string{"/list_current_entry_statuses/sv-se.nst"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		statuses, err := dbm.ListCurrentEntryStatuses(lexRef)
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
	},
}

// TODO cut-n-paste from above
var lexiconListAllEntryStatuses = urlHandler{
	name:     "list_all_entry_statuses",
	url:      "/list_all_entry_statuses/{lexicon_name}",
	help:     "List all entry statuses.",
	examples: []string{"/list_all_entry_statuses/sv-se.nst"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		statuses, err := dbm.ListAllEntryStatuses(lexRef)
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
	},
}

// LexInfo is a struct for collecting lexicon info for json result
type LexInfo struct {
	Name          string `json:"name"`
	SymbolSetName string `json:"symbolSetName"`
}

var lexiconInfo = urlHandler{
	name:     "info",
	url:      "/info/{lexicon_name}",
	help:     "Get some basic lexicon info.",
	examples: []string{"/info/sv-se.nst"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		lex, err := dbm.GetLexicon(lexRef)
		if err != nil {
			http.Error(w, fmt.Sprintf("get lexicon failed : %v", err), http.StatusInternalServerError)
			return
		}
		li := LexInfo{Name: lexRef.String(), SymbolSetName: lex.SymbolSetName}
		jsn, err := marshal(li, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, string(jsn))
	},
}

var lexiconStats = urlHandler{
	name:     "stats",
	url:      "/stats/{lexicon_name}",
	help:     "Lists lexicon stats.",
	examples: []string{"/stats/sv-se.nst"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		stats, err := dbm.LexiconStats(lexRef)
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
	},
}

var lexiconLookup = urlHandler{
	name:     "lookup",
	url:      "/lookup",
	help:     "Lookup in lexicon.",
	examples: []string{"/lookup"},
	handler: func(w http.ResponseWriter, r *http.Request) {

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

		res, err := dbm.LookUpIntoSlice(q)
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
	},
}

var lexiconAddEntryURL = `/addentry?lexicon_name=sv-se.nst&entry={
    "strn": "flesk",
    "language": "sv-se",
    "partOfSpeech": "NN",
    "morphology": "SIN-PLU|IND|NOM|NEU",
    "wordParts": "flesk",
    "lemma": {
	"strn": "flesk",
	"reading": "",
	"paradigm": "s7n-övriga ex träd"
    },
    "transcriptions": [
	{
	    
	    "strn": "\" f l E s k",
	    "language": "sv-se"
	}
    ]
}
`

// TODO add tests
var lexiconAddEntry = urlHandler{
	name:     "addentry",
	url:      "/addentry",
	help:     "Add an entry to the database.",
	examples: []string{lexiconAddEntryURL},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		// lexicon, err := dbm.GetLexicon(lexRef)
		// if err != nil {
		// 	msg := fmt.Sprintf("failed to find lexicon %s in database : %v", lexRef, err)
		// 	log.Println(msg)
		// 	http.Error(w, msg, http.StatusInternalServerError)
		// 	return
		// }
		entryJSON := getParam("entry", r)
		var e lex.Entry
		err = json.Unmarshal([]byte(entryJSON), &e)
		if err != nil {
			log.Printf("lexserver: Failed to unmarshal json: %v", err)
			http.Error(w, fmt.Sprintf("failed to process incoming Entry json : %v", err), http.StatusInternalServerError)
			return
		}

		ids, err := dbm.InsertEntries(lexRef, []lex.Entry{e})
		if err != nil {
			msg := fmt.Sprintf("lexserver failed to update entry : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, ids)
	},
}

var lexiconMoveNewEntries = urlHandler{
	name:     "move_new_entries",
	url:      "/move_new_entries/{from_lexicon}/{to_lexicon}/{new_source}/{new_status}",
	help:     "Move entries from one lexicon to another.",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {
		fromLexName := delQuote(getParam("from_lexicon", r))
		if fromLexName == "" {
			http.Error(w, "no value for parameter 'from_lexicon'", http.StatusBadRequest)
			return
		}
		toLexName := delQuote(getParam("to_lexicon", r))
		if toLexName == "" {
			http.Error(w, "no value for parameter 'to_lexicon'", http.StatusBadRequest)
			return
		}

		sourceName := delQuote(getParam("new_source", r))
		if sourceName == "" {
			http.Error(w, "no value for parameter 'source'", http.StatusBadRequest)
			return
		}
		statusName := delQuote(getParam("new_status", r))
		if statusName == "" {
			http.Error(w, "no value for parameter 'status'", http.StatusBadRequest)
			return
		}

		fromLex, err := lex.ParseLexRef(fromLexName)
		if err != nil {
			http.Error(w, fmt.Sprintf("failure when trying to move entries from '%s' to '%s' : %v", fromLexName, toLexName, err), http.StatusInternalServerError)
			return
		}
		toLex, err := lex.ParseLexRef(toLexName)
		if err != nil {
			http.Error(w, fmt.Sprintf("failure when trying to move entries from '%s' to '%s' : %v", fromLexName, toLexName, err), http.StatusInternalServerError)
			return
		}

		if fromLex.DBRef != toLex.DBRef {
			http.Error(w, fmt.Sprintf("can only move entries with in the same database, from %s to %s", fromLex.String(), toLex.String()), http.StatusInternalServerError)
		}

		moveRes, err := dbm.MoveNewEntries(fromLex.DBRef, fromLex.LexName, toLex.LexName, sourceName, statusName)
		if err != nil {
			http.Error(w, fmt.Sprintf("failure when trying to move entries from '%s' to '%s' : %v", fromLexName, toLexName, err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "number of entries moved from '%s' to '%s': %d", fromLexName, toLexName, moveRes.N)
	},
}

var lexiconValidation = urlHandler{
	name:     "validation (api)",
	url:      "/validation/{lexicon_name}",
	help:     "Validate lexicon (API). Requires POST request.",
	examples: []string{"/validation/en-us.cmu"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, fmt.Sprintf("lexiconfileupload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		start := time.Now()

		clientUUID := getParam("client_uuid", r)

		if "" == strings.TrimSpace(clientUUID) {
			msg := "lexiconValidation got no client uuid"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		conn, ok := webSocks.clients[clientUUID]
		if !ok {
			msg := fmt.Sprintf("lexiconValidation couldn't find connection for uuid %v", clientUUID)
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		logger := dbapi.NewWebSockLogger(conn)

		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		lexicon, err := dbm.GetLexicon(lexRef)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't retrive lexicon : %v", err), http.StatusInternalServerError)
			return
		}
		vMut.Lock()
		v, err := vMut.service.ValidatorForName(lexicon.SymbolSetName)
		vMut.Unlock()
		if err != nil {
			msg := fmt.Sprintf("lexiconValidation failed to get validator for symbol set %v : %v", lexicon.SymbolSetName, err)
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		q := dbapi.Query{}
		stats, err := dbm.Validate(lexRef, logger, *v, q)
		if err != nil {
			msg := fmt.Sprintf("lexiconValidation failed validate : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		dur := round(time.Since(start), time.Second)
		fmt.Fprintf(w, "\nDuration %v\n", dur)
		fmt.Fprint(w, stats)
	},
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
