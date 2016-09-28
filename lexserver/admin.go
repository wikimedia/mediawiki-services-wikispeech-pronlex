package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
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
	lexName := r.PostFormValue("lexicon_name")

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

	err = dbapi.ImportLexiconFile(db, logger, lexName, serverPath)

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
}
