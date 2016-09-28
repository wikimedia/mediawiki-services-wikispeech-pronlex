package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	//"os"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/stts-se/pronlex/symbolset"
)

// The calls prefixed with '/mapper/'

var mMut = struct {
	sync.RWMutex
	service symbolset.MapperService
}{
	service: symbolset.MapperService{
		SymbolSets: make(map[string]symbolset.SymbolSet),
		Mappers:    make(map[string]symbolset.Mapper),
	},
}

type JsonMapped struct {
	From   string
	To     string
	Input  string
	Result string
}

func mapMapperHandler(w http.ResponseWriter, r *http.Request) {
	fromName := r.FormValue("from")
	toName := r.FormValue("to")
	trans := r.FormValue("trans")
	if len(strings.TrimSpace(fromName)) == 0 {
		msg := fmt.Sprintf("input symbol set should be specified by variable 'from'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	if len(strings.TrimSpace(toName)) == 0 {
		msg := fmt.Sprintf("output symbol set should be specified by variable 'to'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	if len(strings.TrimSpace(trans)) == 0 {
		msg := fmt.Sprintf("input trans should be specified by variable 'trans'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	mMut.Lock()
	result0, err := mMut.service.Map(fromName, toName, trans)
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("failed mapping transcription : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	result := JsonMapped{Input: trans, Result: result0, From: fromName, To: toName}
	j, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("json marshalling error : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(j))
}

type JsonMapper struct {
	From    string
	To      string
	Symbols []JsonMSymbol
}

type JsonMSymbol struct {
	From string
	To   string
	IPA  string
	Desc string
	Cat  string
}

func mapTableMapperHandler(w http.ResponseWriter, r *http.Request) {
	fromName := r.FormValue("from")
	toName := r.FormValue("to")
	if len(strings.TrimSpace(fromName)) == 0 {
		msg := fmt.Sprintf("input symbol set should be specified by variable 'from'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	if len(strings.TrimSpace(toName)) == 0 {
		msg := fmt.Sprintf("output symbol set should be specified by variable 'to'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	mMut.Lock()
	mapper0, err := mMut.service.GetMapTable(fromName, toName)
	mMut.Unlock()
	if err != nil {
		msg := fmt.Sprintf("failed getting map table : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	mapper := JsonMapper{From: mapper0.SymbolSet1.FromName, To: mapper0.SymbolSet2.ToName}
	mapper.Symbols = make([]JsonMSymbol, 0)
	for _, sym := range mapper0.SymbolSet1.Symbols {
		from := sym.Sym1
		to, err := mapper0.MapSymbol(from)
		if err != nil {
			msg := fmt.Sprintf("failed getting map table : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		mapper.Symbols = append(mapper.Symbols, JsonMSymbol{From: from.String, To: to.String, IPA: sym.Sym2.String, Desc: sym.Sym1.Desc, Cat: sym.Sym1.Cat.String()})
	}

	j, err := json.Marshal(mapper)
	if err != nil {
		msg := fmt.Sprintf("json marshalling error : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(j))
}

func listMappersHandler(w http.ResponseWriter, r *http.Request) {
	mMut.Lock()
	ms := mMut.service.MapperNames()
	mMut.Unlock()
	j, err := json.Marshal(ms)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal struct : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(j))
}

func mapperHelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<h1>Mapper</h1>
<h2>map</h2> Maps a transcription from one symbolset to another. Example invocation:
<pre><a href="/mapper/map?from=sv-se_ws-sampa&to=sv-se_sampa_mary&trans=%22%22%20p%20O%20j%20.%20k%20@">/mapper/map?from=sv-se_ws-sampa&to=sv-se_sampa_mary&trans=%22%22%20p%20O%20j%20.%20k%20@</a></pre>

<h2>list</h2> Lists cached mappers. Example invocation:
<pre><a href="/mapper/list">/mapper/list</a></pre>

<h2>maptable</h2> Lists content of a maptable given two symbolset names. Example invocation:
<pre><a href="/mapper/maptable?from=sv-se_ws-sampa&to=sv-se_sampa_mary">/mapper/maptable?from=sv-se_ws-sampa&to=sv-se_sampa_mary</a></pre>
		`

	fmt.Fprint(w, html)
}

func uploadMapperHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/mapper_upload_symbolset.html")
}

func doUploadMapperHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("symbol set upload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
		return
	}

	clientUUID := r.FormValue("client_uuid")

	if "" == strings.TrimSpace(clientUUID) {
		msg := "doUploadMapperHandler got no client uuid"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// (partially) lifted from https://github.com/astaxie/build-web-application-with-golang/blob/master/de/04.5.md

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("doUploadMapperHandler failed reading file : %v", err), http.StatusInternalServerError)
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
		http.Error(w, fmt.Sprintf("doUploadMapperHandler failed opening local output file : %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		msg := fmt.Sprintf("doUploadMapperHandler failed copying local output file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	_, err = loadSymbolSetFile(serverPath)
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

	fmt.Fprintf(w, "%v", handler.Header)
}
