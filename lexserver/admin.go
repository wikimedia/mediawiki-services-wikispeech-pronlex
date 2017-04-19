package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/validation"
)

func deleteUploadedFile(serverPath string) {
	// when done, delete from server!
	err := os.Remove(serverPath)
	if err != nil {
		msg := fmt.Sprintf("couldn't delete temp file from server : %v", err)
		log.Println(msg)
	} else {
		msg := fmt.Sprint("the uploaded temp file has been deleted from server")
		log.Println(msg)
	}
}

func adminDoLexImportHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("lexiconfileupload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
		return
	}

	clientUUID := r.FormValue("client_uuid")

	if "" == strings.TrimSpace(clientUUID) {
		msg := "adminDoLexImportHandler got no client uuid"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	conn, ok := webSocks.clients[clientUUID]
	if !ok {
		msg := fmt.Sprintf("adminDoLexImportHandler couldn't find connection for uuid %v", clientUUID)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	logger := dbapi.NewWebSockLogger(conn)

	symbolSetName := r.PostFormValue("symbolset_name")
	if strings.TrimSpace(symbolSetName) == "" {
		msg := "input param <symbolset_name> must not be empty"
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	lexName := r.PostFormValue("lexicon_name")
	if strings.TrimSpace(lexName) == "" {
		msg := "input param <lexicon_name> must not be empty"
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	vString := r.PostFormValue("validate")
	if strings.TrimSpace(vString) == "" {
		msg := "input param <validate> must not be empty (should be 'true' or 'false')"
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	validate, err := strconv.ParseBool(vString)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("adminDoLexImportHandler failed parsing boolean argument %s : %v", vString, err), http.StatusInternalServerError)
		return
	}
	// (partially) lifted from https://github.com/astaxie/build-web-application-with-golang/blob/master/de/04.5.md

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("adminDoLexImportHandler failed reading file : %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	serverPath := filepath.Join(uploadFileArea, handler.Filename)

	f, err := os.OpenFile(serverPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("adminDoLexImportHandler failed opening local output file : %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		msg := fmt.Sprintf("adminDoLexImportHandler failed copying local output file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, err = dbapi.GetLexicon(db, lexName)
	if err == nil {
		msg := fmt.Sprintf("Nothing will be added. Lexicon already exists in database: %s", lexName)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		deleteUploadedFile(serverPath)
		return
	}

	lexicon := dbapi.Lexicon{Name: lexName, SymbolSetName: symbolSetName}
	lexicon, err = dbapi.InsertLexicon(db, lexicon)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
		deleteUploadedFile(serverPath)
	}
	log.Println("Created lexicon with name:", lexicon.Name)

	var validator *validation.Validator
	if validate {
		vMut.Lock()
		validator, err = vMut.service.ValidatorForName(lexicon.SymbolSetName)
		vMut.Unlock()
		if err != nil {
			msg := fmt.Sprintf("lexiconRunValidateHandler failed to get validator for symbol set %v : %v", lexicon.SymbolSetName, err)
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
	}

	err = dbapi.ImportLexiconFile(db, logger, lexName, serverPath, validator)

	if err == nil {
		msg := fmt.Sprintf("lexicon file imported successfully : %v", handler.Filename)
		log.Println(msg)
	} else {
		msg := fmt.Sprintf("couldn't import lexicon file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		deleteUploadedFile(serverPath)
		return
	}

	f.Close()
	deleteUploadedFile(serverPath)

	// if validate {
	// 	start := time.Now()
	// 	logger.Write("validating lexicon ... ")
	// 	fmt.Fprintf(w, "validating lexicon ...")
	// 	vMut.Lock()
	// 	v, err := vMut.service.ValidatorForName(lexicon.SymbolSetName)
	// 	vMut.Unlock()
	// 	if err != nil {
	// 		msg := fmt.Sprintf("lexiconRunValidateHandler failed to get validator for symbol set %v : %v", lexicon.SymbolSetName, err)
	// 		log.Println(msg)
	// 		http.Error(w, msg, http.StatusBadRequest)
	// 		return
	// 	}

	// 	q := dbapi.Query{Lexicons: []dbapi.Lexicon{lexicon}}
	// 	stats, err := dbapi.Validate(db, logger, *v, q)
	// 	if err != nil {
	// 		msg := fmt.Sprintf("lexiconRunValidateHandler failed validate : %v", err)
	// 		log.Println(msg)
	// 		http.Error(w, msg, http.StatusBadRequest)
	// 		return
	// 	}
	// 	dur := round(time.Since(start), time.Second)
	// 	fmt.Fprintf(w, "\nValidation took %v\n", dur)
	// 	fmt.Fprint(w, stats)
	// }

	entryCount, err := dbapi.EntryCount(db, lexicon.ID)
	if err != nil {
		msg := fmt.Sprintf("lexicon imported, but couldn't retrieve lexicon info from server : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	info := LexWithEntryCount{
		ID:            lexicon.ID,
		Name:          lexicon.Name,
		SymbolSetName: lexicon.SymbolSetName,
		EntryCount:    entryCount,
	}
	fmt.Fprintf(w, "imported %v entries into lexicon '%v'", info.EntryCount, info.Name)
}

func adminHelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<h1>Admin</h1>
<h2>lex_import</h2>Import lexicon file.
<pre><a href="/admin/lex_import">/admin/lex_import</a></pre>

<h2>deletelexicon</h2> Delete a lexicon reference from the database without removing associated entries.
<pre>/admin/deletelexicon?id=N</a></pre>

<h2>superdeletelexicon</h2> Delete a complete lexicon from the database, including associated entries.
<pre>/admin/superdeletelexicon?id=N</a></pre>

		`

	fmt.Fprint(w, html)
}
