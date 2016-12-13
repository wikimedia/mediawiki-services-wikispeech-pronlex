package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/pronlex/validation/rules"
)

// The calls prefixed with '/validation/'

var vMut = struct {
	sync.RWMutex
	service rules.ValidatorService
}{
	service: rules.ValidatorService{Validators: make(map[string]*validation.Validator)},
}

func loadValidators(symsetDirName string) error {
	symbolSets, err := symbolset.LoadSymbolSetsFromDir(symsetDirName)
	if err != nil {
		return err
	}
	vMut.Lock()
	err = vMut.service.Load(symbolSets)
	vMut.Unlock()
	return err
}

// TODO code duplication between validateEntriesHandler and validateEntryHandler

func validateEntriesHandler(w http.ResponseWriter, r *http.Request) {
	entriesJSON := r.FormValue("entries")
	symbolSetName := r.FormValue("symbolsetname")

	var es []lex.Entry
	err := json.Unmarshal([]byte(entriesJSON), &es) //TODO check if OK. NL 20161019 es -> &es
	if err != nil {
		msg := fmt.Sprintf("lexserver: Failed to unmarshal json: %v : %v", entriesJSON, err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if symbolSetName == "" {
		msg := "validateEntryHandler expected a symbol set name"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	vMut.Lock()
	vdator, err := vMut.service.ValidatorForName(symbolSetName)
	vMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("validateEntryHandler failed to get validator for symbol set %v : %v", symbolSetName, err)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	es = trimEntries(es)
	es, _ = vdator.ValidateEntries(es)

	res0, err3 := json.Marshal(es)
	if err3 != nil {
		msg := fmt.Sprintf("lexserver: Failed to marshal entry : %v", err3)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(res0))
}

// TODO code duplication between validateEntriesHandler and validateEntryHandler

func validateEntryHandler(w http.ResponseWriter, r *http.Request) {
	entryJSON := r.FormValue("entry")
	symbolSetName := r.FormValue("symbolsetname")

	var e lex.Entry
	err := json.Unmarshal([]byte(entryJSON), &e)
	if err != nil {
		msg := fmt.Sprintf("lexserver: Failed to unmarshal json: %v : %v", entryJSON, err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if symbolSetName == "" {
		msg := "validateEntryHandler expected a symbol set name"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// TODO Hardwired stuff below!!!!
	vMut.Lock()
	vdator, err := vMut.service.ValidatorForName(symbolSetName)
	vMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("validateEntryHandler failed to get validator for symbol set %v : %v", symbolSetName, err)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	e = trimEntry(e)
	e, _ = vdator.ValidateEntry(e)

	res0, err3 := json.Marshal(e)
	if err3 != nil {
		msg := fmt.Sprintf("lexserver: Failed to marshal entry : %v", err3)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(res0))
}

func validationStatsHandler(w http.ResponseWriter, r *http.Request) {
	lexiconIdS := r.FormValue("lexiconId")
	if len(strings.TrimSpace(lexiconIdS)) == 0 {
		msg := fmt.Sprintf("lexicon id should be specified by variable 'lexiconId'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	lexiconId, err := strconv.ParseInt(lexiconIdS, 10, 64)
	if err != nil {
		msg := fmt.Sprintf("lexicon id should be an integer")
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	stats, err := dbapi.ValidationStats(db, lexiconId)
	if err != nil {
		msg := fmt.Sprintf("validationStatsHandler failed to retreive validation stats : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	j, err := json.Marshal(stats)
	if err != nil {
		msg := fmt.Sprintf("lexserver: Failed to marshal stats : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(j))
}

var trimWhitespaceRe = regexp.MustCompile("[\\s]+")

func trimTranscriptions(e lex.Entry) lex.Entry {
	var newTs []lex.Transcription
	for _, t := range e.Transcriptions {
		s := trimWhitespaceRe.ReplaceAllString(strings.TrimSpace(t.Strn), " ")
		t.Strn = s
		newTs = append(newTs, t)
	}
	e.Transcriptions = newTs
	return e
}

func trimEntry(e lex.Entry) lex.Entry {
	e = trimTranscriptions(e)
	e.Strn = strings.TrimSpace(e.Strn)
	e.WordParts = strings.TrimSpace(e.WordParts)
	e.PartOfSpeech = trimWhitespaceRe.ReplaceAllString(strings.TrimSpace(e.PartOfSpeech), " ")
	return e
}

func trimEntries(entries []lex.Entry) []lex.Entry {
	var res []lex.Entry
	for _, e := range entries {
		trimEntry(e)
		res = append(res, e)
	}
	return res
}

type ValidatorNames struct {
	ValidatorNames []string `json:validator_names`
}

func validatorNames() ValidatorNames {
	var vNames []string
	vMut.Lock()
	for vName, _ := range vMut.service.Validators {
		vNames = append(vNames, vName)
	}
	vMut.Unlock()
	sort.Strings(vNames)
	return ValidatorNames{ValidatorNames: vNames}
}

func listValidationHandler(w http.ResponseWriter, r *http.Request) {
	vs := validatorNames()
	j, err := json.Marshal(vs)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal struct : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(j))
}

func validationHelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	validateEntryUrl_en_us := `/validation/validateentry?symbolsetname=en-us_sampa_mary&entry={%22id%22:1703348,%22lexiconId%22:3,%22strn%22:%22barn%22,%22language%22:%22en-us%22,%22partOfSpeech%22:%22%22,%22wordParts%22:%22%22,%22lemma%22:{%22id%22:0,%22strn%22:%22%22,%22reading%22:%22%22,%22paradigm%22:%22%22},%22transcriptions%22:[{%22id%22:1717337,%22entryId%22:1703348,%22strn%22:%22\%22%20b%20A%20r%20n%22,%22language%22:%22%22,%22sources%22:[]}],%22status%22:{%22id%22:1703348,%22name%22:%22imported%22,%22source%22:%22cmu%22,%22timestamp%22:%222016-09-06T13:16:07Z%22,%22current%22:true},%22entryValidations%22:[]}`

	validateEntryUrl := `/validation/validateentry?symbolsetname=sv-se_ws-sampa&entry={%22id%22:371546,%22lexiconId%22:1,%22strn%22:%22h%C3%A4st%22,%22language%22:%22SWE%22,%22partOfSpeech%22:%22NN%20SIN|IND|NOM|UTR%22,%22wordParts%22:%22h%C3%A4st%22,%22lemma%22:{%22id%22:42815,%22strn%22:%22h%C3%A4st%22,%22reading%22:%22%22,%22paradigm%22:%22s2q-lapp%22},%22transcriptions%22:[{%22id%22:377191,%22entryId%22:371546,%22strn%22:%22\%22%20h%20E%20s%20t%22,%22language%22:%22SWE%22,%22sources%22:[]}],%22status%22:{%22id%22:371546,%22name%22:%22imported%22,%22source%22:%22nst%22,%22timestamp%22:%222016-09-06T12:54:12Z%22,%22current%22:true}}`

	validateEntriesUrl := "/validation/validateentries?symbolsetname=sv-se_ws-sampa&entries=..."

	html := `<h1>Validation</h1>
<h2>validateentry</h2> Validates an entry. Example invocation:
<pre><a href="` + validateEntryUrl + `">` + validateEntryUrl + `</a></pre>
<pre><a href="` + validateEntryUrl_en_us + `">` + validateEntryUrl_en_us + `</a></pre>

<h2>validateentries</h2> Validates a list of entries. Example invocation:
<pre><a href="` + validateEntriesUrl + `">` + validateEntriesUrl + `</a></pre>

<h2>list</h2> Lists available validators. Example invocation:
<pre><a href="/validation/list">/validation/list</a></pre>

<h2>stats</h2> Lists validation stats. Example invocation:
<pre><a href="/validation/stats?lexiconId=1">/validation/stats?lexiconId=1</a></pre>
		`

	fmt.Fprint(w, html)
}
