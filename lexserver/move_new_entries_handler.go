package main

import (
	"fmt"
	"net/http"

	"github.com/stts-se/pronlex/dbapi"
)

func moveNewEntriesHandler(w http.ResponseWriter, r *http.Request) {

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
	sourceName := delQuote(getParam("source", r))
	if sourceName == "" {
		http.Error(w, "no value for parameter 'source'", http.StatusBadRequest)
		return
	}
	statusName := delQuote(getParam("status", r))
	if statusName == "" {
		http.Error(w, "no value for parameter 'status'", http.StatusBadRequest)
		return
	}

	moveRes, err := dbapi.MoveNewEntries(db, fromLexName, toLexName, sourceName, statusName)
	if err != nil {
		http.Error(w, fmt.Sprintf("failure when trying to move entries from '%s' to '%s' : %v", fromLexName, toLexName, err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "number of entries moved from '%s' to '%s': %d", fromLexName, toLexName, moveRes.N)
}
