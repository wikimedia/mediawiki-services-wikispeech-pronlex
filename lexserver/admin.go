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

var adminLexImportPage = urlHandler{
	name:     "lex_import (page)",
	url:      "/lex_import_page",
	help:     "Import lexicon file (GUI).",
	examples: []string{"/lex_import_page"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/admin/lex_import_page.html")
	},
}

var adminLexImport = urlHandler{
	name:     "lex_import (api)",
	url:      "/lex_import",
	help:     "Import lexicon file (API). Requires POST request. Mainly for server internal use.<p/>Available params: lexicon_name, symbolset_name, validate, file",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {

		defer protect(w) // use this call in handlers to catch 'panic' and stack traces and returning a general error to the calling client

		if r.Method != "POST" {
			http.Error(w, fmt.Sprintf("lexiconfileupload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		clientUUID := getParam("client_uuid", r)

		if "" == strings.TrimSpace(clientUUID) {
			msg := "adminLexImport got no client uuid"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}
		conn, ok := webSocks.clients[clientUUID]
		if !ok {
			msg := fmt.Sprintf("adminLexImport couldn't find connection for uuid %v", clientUUID)
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
			http.Error(w, fmt.Sprintf("adminLexImport failed parsing boolean argument %s : %v", vString, err), http.StatusInternalServerError)
			return
		}
		// (partially) lifted from https://github.com/astaxie/build-web-application-with-golang/blob/master/de/04.5.md

		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("adminLexImport failed reading file : %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		serverPath := filepath.Join(uploadFileArea, handler.Filename)

		f, err := os.OpenFile(serverPath, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("adminLexImport failed opening local output file : %v", err), http.StatusInternalServerError)
			return
		}
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			msg := fmt.Sprintf("adminLexImport failed copying local output file : %v", err)
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
		lexicon, err = dbapi.DefineLexicon(db, lexicon)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			deleteUploadedFile(serverPath)
			return
		}
		log.Println("Created lexicon with name:", lexicon.Name)

		var validator *validation.Validator = &validation.Validator{}
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
	},
}

var adminDeleteLex = urlHandler{
	name:     "deletelexicon",
	url:      "/deletelexicon/{lexicon_name}",
	help:     "Delete a lexicon reference from the database without removing associated entries.",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexName := getParam("lexicon_name", r)
		if lexName == "" {
			msg := "adminDeleteLex expected a lexicon name defined by 'lexicon_name'"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		err := dbapi.DeleteLexicon(db, lexName)
		if err != nil {
			log.Printf("adminDeleteLex got error : %v\n", err)
			http.Error(w, fmt.Sprintf("failed deleting lexicon : %v", err), http.StatusExpectationFailed)
			return
		}
	},
}

var adminSuperDeleteLex = urlHandler{
	name:     "superdeletelexicon",
	url:      "/superdeletelexicon/{lexicon_name}",
	help:     "Delete a complete lexicon from the database, including associated entries. This make take some time.",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexName := getParam("lexicon_name", r)
		if lexName == "" {
			msg := "adminDeleteLex expected a lexicon name defined by 'lexicon_name'"
			log.Println(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		uuid := getParam("client_uuid", r)
		log.Println("adminSuperDeleteLex was called")
		messageToClientWebSock(uuid, fmt.Sprintf("Super delete was called. This may take quite a while. Lexicon name %s", lexName))
		err := dbapi.SuperDeleteLexicon(db, lexName)
		if err != nil {

			http.Error(w, fmt.Sprintf("failed super deleting lexicon : %v", err), http.StatusExpectationFailed)
			return
		}

		messageToClientWebSock(uuid, fmt.Sprintf("Done deleting lexicon with name %s", lexName))
	},
}
