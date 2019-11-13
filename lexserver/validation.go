package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/validation"
	"github.com/stts-se/pronlex/validation/validators"
)

// The calls prefixed with '/validation/'

var vMut = struct {
	sync.RWMutex
	service validators.ValidatorService
}{
	service: validators.ValidatorService{Validators: make(map[string]*validation.Validator)},
}

func loadValidators(symsetDirName string) error {
	symbolSets, err := symbolset.LoadSymbolSetsFromDir(symsetDirName)
	if err != nil {
		return err
	}
	vMut.Lock()
	err = vMut.service.Load(symbolSets, symsetDirName)
	vMut.Unlock()
	return err
}

// TODO code duplication between validateEntriesHandler and validateEntryHandler
var validationValidateEntry = urlHandler{
	name:     "validateentry",
	url:      "/validateentry",
	help:     "Validates one entry. Input in JSON format. For examples, see <a href=\"https://godoc.org/github.com/stts-se/pronlex/lex\">package documentation</a>",
	examples: []string{`/validateentry?symbolsetname=sv-se_ws-sampa-DEMO&entry={%22id%22:371546,%22lexiconId%22:1,%22strn%22:%22h%C3%A4st%22,%22language%22:%22SWE%22,%22partOfSpeech%22:%22NN%20SIN|IND|NOM|UTR%22,%22wordParts%22:%22h%C3%A4st%22,%22lemma%22:{%22id%22:42815,%22strn%22:%22h%C3%A4st%22,%22reading%22:%22%22,%22paradigm%22:%22s2q-lapp%22},%22transcriptions%22:[{%22id%22:377191,%22entryId%22:371546,%22strn%22:%22\%22%20h%20E%20s%20t%22,%22language%22:%22SWE%22,%22sources%22:[]}],%22status%22:{%22id%22:371546,%22name%22:%22imported%22,%22source%22:%22nst%22,%22timestamp%22:%222016-09-06T12:54:12Z%22,%22current%22:true}}`},
	handler: func(w http.ResponseWriter, r *http.Request) {
		entryJSON := getParam("entry", r)
		if entryJSON == "" {
			msg := "validateentry expected param entry"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		symbolSetName := getParam("symbolsetname", r)

		var e lex.Entry
		err := json.Unmarshal([]byte(entryJSON), &e)
		if err != nil {
			msg := fmt.Sprintf("lexserver: Failed to unmarshal json: %v : %v", entryJSON, err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if symbolSetName == "" {
			msg := "validateentry expected a symbol set name"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		// TODO Hardwired stuff below!!!!
		vMut.Lock()
		vdator, err := vMut.service.ValidatorForName(symbolSetName)
		vMut.Unlock()
		if err != nil {
			msg := fmt.Sprintf("validateentry failed to get validator for symbol set %v : %v", symbolSetName, err)
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		e = trimEntry(e)
		vdator.ValidateEntry(&e)

		res0, err3 := json.Marshal(e)
		if err3 != nil {
			msg := fmt.Sprintf("lexserver: Failed to marshal entry : %v", err3)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, string(res0))
	},
}

var validationStats = urlHandler{
	name:     "stats",
	url:      "/stats/{lexicon_name}",
	help:     "Lists validation stats.",
	examples: []string{"/stats/lexserver_testdb:sv"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		stats, err := dbm.ValidationStats(lexRef)
		if err != nil {
			msg := fmt.Sprintf("validationStatsHandler failed to retrieve validation stats : %v", err)
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
	},
}

// TODO code duplication between validateEntries and validateEntry
var validationValidateEntries = urlHandler{
	name:     "validateentries",
	url:      "/validateentries",
	help:     "Validates a list of entries.  Input as a list in JSON format. For examples, see <a href=\"https://godoc.org/github.com/stts-se/pronlex/lex\">package documentation</a>",
	examples: []string{`/validateentries?symbolsetname=sv-se_ws-sampa-DEMO&entries=[{%22id%22:371546,%22lexiconId%22:1,%22strn%22:%22h%C3%A4st%22,%22language%22:%22SWE%22,%22partOfSpeech%22:%22NN%20SIN|IND|NOM|UTR%22,%22wordParts%22:%22h%C3%A4st%22,%22lemma%22:{%22id%22:42815,%22strn%22:%22h%C3%A4st%22,%22reading%22:%22%22,%22paradigm%22:%22s2q-lapp%22},%22transcriptions%22:[{%22id%22:377191,%22entryId%22:371546,%22strn%22:%22\%22%20h%20E%20s%20t%22,%22language%22:%22SWE%22,%22sources%22:[]}],%22status%22:{%22id%22:371546,%22name%22:%22imported%22,%22source%22:%22nst%22,%22timestamp%22:%222016-09-06T12:54:12Z%22,%22current%22:true}},{%22id%22:371546,%22lexiconId%22:1,%22strn%22:%22host%22,%22language%22:%22SWE%22,%22partOfSpeech%22:%22NN%20SIN|IND|NOM|UTR%22,%22wordParts%22:%22host%22,%22lemma%22:{%22id%22:42815,%22strn%22:%22h%C3%A4st%22,%22reading%22:%22%22,%22paradigm%22:%22s2q-lapp%22},%22transcriptions%22:[{%22id%22:377191,%22entryId%22:371546,%22strn%22:%22\%22%20h%20U%20s%20t%22,%22language%22:%22SWE%22,%22sources%22:[]}],%22status%22:{%22id%22:371546,%22name%22:%22imported%22,%22source%22:%22nst%22,%22timestamp%22:%222016-09-06T12:54:12Z%22,%22current%22:true}}]`},
	handler: func(w http.ResponseWriter, r *http.Request) {
		entriesJSON := getParam("entries", r)
		if entriesJSON == "" {
			msg := "validateentry expected param entries"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		symbolSetName := getParam("symbolsetname", r)

		var es []lex.Entry
		err := json.Unmarshal([]byte(entriesJSON), &es) //TODO check if OK. NL 20161019 es -> &es
		if err != nil {
			msg := fmt.Sprintf("lexserver: Failed to unmarshal json: %v : %v", entriesJSON, err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		if symbolSetName == "" {
			msg := "validateentries expected a symbol set name"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		vMut.Lock()
		vdator, err := vMut.service.ValidatorForName(symbolSetName)
		vMut.Unlock()
		if err != nil {
			msg := fmt.Sprintf("validateentries failed to get validator for symbol set %v : %v", symbolSetName, err)
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
	},
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

func validatorNames() []string {
	var vNames []string
	vMut.Lock()
	for vName := range vMut.service.Validators {
		vNames = append(vNames, vName)
	}
	vMut.Unlock()
	sort.Strings(vNames)
	return vNames
}

func hasValidator(symbolSet string) bool {
	vMut.Lock()
	res := vMut.service.HasValidator(symbolSet)
	vMut.Unlock()
	return res
}

var validationListValidators = urlHandler{
	name:     "list",
	url:      "/list",
	help:     "Lists available validators.",
	examples: []string{"/list"},
	handler: func(w http.ResponseWriter, r *http.Request) {
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
	},
}

var validationHasValidator = urlHandler{
	name:     "has_validator",
	url:      "/has_validator/{symbolset}",
	help:     "Checks if a symbol set has an associated validator.",
	examples: []string{"/has_validator/sv-se_ws-sampa", "/has_validator/ar_ws-sampa"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		symbolSet := getParam("symbolset", r)
		if len(strings.TrimSpace(symbolSet)) == 0 {
			msg := fmt.Sprintf("symbol set should be specified by variable 'symbolset'")
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		res := hasValidator(symbolSet)

		j, err := json.Marshal(res)
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

/*
func validationHelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<h1>Validation</h1>`

	fmt.Fprint(w, html)
}
*/
