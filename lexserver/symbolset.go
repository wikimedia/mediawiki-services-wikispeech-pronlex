package main

// The calls prefixed with '/symbolset/'

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	//"os"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/stts-se/pronlex/symbolset"
)

type JsonSymbolSet struct {
	Name    string
	Symbols []JsonSymbol
}

type JsonSymbol struct {
	Symbol string
	IPA    string
	Desc   string
	Cat    string
}

func symbolSetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if len(strings.TrimSpace(name)) == 0 {
		msg := fmt.Sprintf("symbol set should be specified by variable 'name'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	mMut.Lock()
	symbolset0, ok := mMut.service.SymbolSets[name]
	mMut.Unlock()
	if !ok {
		msg := fmt.Sprintf("failed getting symbol set : %v", name)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	symbolset := JsonSymbolSet{Name: symbolset0.Name}
	symbolset.Symbols = make([]JsonSymbol, 0)
	for _, sym := range symbolset0.Symbols {
		symbolset.Symbols = append(symbolset.Symbols, JsonSymbol{Symbol: sym.Sym1.String, IPA: sym.Sym2.String, Desc: sym.Sym1.Desc, Cat: sym.Sym1.Cat.String()})
	}

	j, err := json.Marshal(symbolset)
	if err != nil {
		msg := fmt.Sprintf("json marshalling error : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(j))
}

func loadSymbolSets(dirName string) error {
	mMut.Lock()
	mMut.service.Clear()
	mMut.Unlock()

	symbolSets, err := loadSymbolSetsFromDir(dirName)
	if err != nil {
		return err
	}
	mMut.Lock()
	mMut.service.SymbolSets = symbolSets
	mMut.Unlock()
	return nil
}

func reloadAllSymbolSetsHandler(w http.ResponseWriter, r *http.Request) {
	err := loadSymbolSets(symbolSetFileArea)
	if err != nil {
		msg := err.Error()
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	mMut.Lock()
	j, err := json.Marshal(symbolSetNames(mMut.service.SymbolSets))
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("json marshalling error : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(j))

}

func reloadOneSymbolSetHandler(w http.ResponseWriter, r *http.Request, name string) {
	mMut.Lock()
	err := mMut.service.Delete(name)
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("couldn't delete symbolset : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	serverPath := filepath.Join(symbolSetFileArea, name+symbolSetSuffix)
	mMut.Lock()
	err = mMut.service.Load(serverPath)
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("couldn't load symbolset : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	msg := fmt.Sprintf("Reloaded symbol set %s", name)
	fmt.Fprint(w, msg)

}

func reloadSymbolSetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if len(strings.TrimSpace(name)) == 0 {
		reloadAllSymbolSetsHandler(w, r)
	} else {
		reloadOneSymbolSetHandler(w, r, name)
	}

}

func deleteSymbolSetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if len(strings.TrimSpace(name)) == 0 {
		msg := fmt.Sprintf("symbol set should be specified by variable 'name'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	mMut.Lock()
	err := mMut.service.Delete(name)
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("couldn't delete symbolset : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	serverPath := filepath.Join(symbolSetFileArea, name+symbolSetSuffix)
	if _, err := os.Stat(serverPath); err != nil {
		if os.IsNotExist(err) {
			msg := fmt.Sprintf("couldn't locate server file for symbol set %s", name)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}

	err = os.Remove(serverPath)
	if err != nil {
		msg := fmt.Sprintf("couldn't delete file from server : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	msg := fmt.Sprintf("Deleted symbol set %s", name)
	fmt.Fprint(w, msg)
}

func listSymbolSetsHandler(w http.ResponseWriter, r *http.Request) {
	mMut.Lock()
	ss := symbolSetNames(mMut.service.SymbolSets)
	mMut.Unlock()
	j, err := json.Marshal(ss)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal struct : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(j))
}

func symbolSetHelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<h1>SymbolSet</h1>
<h2>reload</h2> Reloads symbol set(s) in the pre-defined folder. All affected mappers will also be removed from cache. Example invocation:
<pre><a href="/symbolset/reload">/mapper/reload</a></pre>
<pre><a href="/symbolset/reload?name=sv-se_nst-xsampa">/mapper/reload?name=sv-se_nst-xsampa</a></pre>

<h2>list</h2> Lists available symbol sets. Example invocation:
<pre><a href="/symbolset/list">/symbolset/list</a></pre>

<h2>delete</h2> Deletes a named symbol set. Example invocation:
<pre><a href="/symbolset/delete?name=sv-se_nst-xsampa">/symbolset/delete?name=sv-se_nst-xsampa</a></pre>

<h2>symbolset</h2> Lists content of a named symbolset. Example invocation:
<pre><a href="/symbolset/symbolset?name=sv-se_ws-sampa">/symbolset/symbolset?name=sv-se_ws-sampa</a></pre>

<h2>symbolset_upload</h2> Upload symbol set file
<pre><a href="/symbolset_upload">/symbolset_upload</a></pre>		
		`

	fmt.Fprint(w, html)
}

func uploadSymbolSetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/symbolset_upload.html")
}

func doUploadSymbolSetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("symbol set upload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
		return
	}

	clientUUID := r.FormValue("client_uuid")

	if "" == strings.TrimSpace(clientUUID) {
		msg := "doUploadSymbolSetHandler got no client uuid"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// (partially) lifted from https://github.com/astaxie/build-web-application-with-golang/blob/master/de/04.5.md

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("doUploadSymbolSetHandler failed reading file : %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	serverPath := filepath.Join(symbolSetFileArea, handler.Filename)
	if _, err := os.Stat(serverPath); err == nil {
		msg := fmt.Sprintf("symbol set already exists on server in file: %s", handler.Filename)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	f, err := os.OpenFile(serverPath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("doUploadSymbolSetHandler failed opening local output file : %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		msg := fmt.Sprintf("doUploadSymbolSetHandler failed copying local output file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	ss, err := loadSymbolSetFile(serverPath)
	if err != nil {
		msg := fmt.Sprintf("couldn't load symbol set file : %v", err)
		err = os.Remove(serverPath)
		if err != nil {
			msg = fmt.Sprintf("%v (couldn't delete file from server)", msg)
		} else {
			msg = fmt.Sprintf("%v (the uploaded file has been deleted from server)", msg)
		}
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	f.Close()

	mMut.Lock()
	mMut.service.SymbolSets[ss.Name] = ss
	mMut.Unlock()

	fmt.Fprintf(w, "%v", handler.Header)
}

type SymbolSetNames struct {
	SymbolSetNames []string `json:symbol_set_names`
}

func symbolSetNames(sss map[string]symbolset.SymbolSet) SymbolSetNames {
	var ssNames []string
	for ss, _ := range sss {
		ssNames = append(ssNames, ss)
	}
	sort.Strings(ssNames)
	return SymbolSetNames{SymbolSetNames: ssNames}
}
