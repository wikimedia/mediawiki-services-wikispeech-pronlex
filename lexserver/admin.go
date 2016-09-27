package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

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

	ssName := r.PostFormValue("symbolset_name")
	lexName := r.PostFormValue("lexicon_name")
	fmt.Println(ssName)
	fmt.Println(lexName)

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
	// if _, err := os.Stat(serverPath); err == nil {
	// 	msg := fmt.Sprintf("lexicon file already exists on server in file: %s", handler.Filename)
	// 	log.Println(msg)
	// 	http.Error(w, msg, http.StatusInternalServerError)
	// 	return
	// }

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
	// _, err = importXXXX(serverPath)
	// if err != nil {
	// 	msg := fmt.Sprintf("couldn't import lexicon file : %v", err)
	// 	err = os.Remove(serverPath)
	// 	if err != nil {
	// 		msg = fmt.Sprintf("%v (couldn't delete file from server)", msg)
	// 	} else {
	// 		msg = fmt.Sprintf("%v (the uploaded file has been deleted from server)", msg)
	// 	}
	// 	log.Println(msg)
	// 	http.Error(w, msg, http.StatusInternalServerError)
	// 	return
	// }

	f.Close()

	fmt.Fprintf(w, "%v", handler.Header)
}
