package main

// The calls prefixed with '/mapper/'

import (
	"fmt"
	"io"
	"io/ioutil"
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

// TODO Mutex this variable
var mapperService = symbolset.MapperService{}

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
	result0, err := mapperService.Map(fromName, toName, trans)
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

func symbolSetMapperHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if len(strings.TrimSpace(name)) == 0 {
		msg := fmt.Sprintf("symbol set should be specified by variable 'name'")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	symbolset0, ok := mapperService.SymbolSets[name]
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
	mapper0, err := mapperService.GetMapTable(fromName, toName)
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

func loadMappersFromDir(dirName string) error {
	mapperService.Clear()

	// list files in symbol set dir
	fileInfos, err := ioutil.ReadDir(symbolSetFileArea)
	if err != nil {
		return fmt.Errorf("failed reading symbol set dir : %v", err)
	}

	var fErrs error
	var symSets []symbolset.SymbolSet
	for _, fi := range fileInfos {
		if strings.HasSuffix(fi.Name(), ".tab") {
			symset, err := symbolset.LoadSymbolSet(filepath.Join(symbolSetFileArea, fi.Name()))
			if err != nil {
				if fErrs != nil {
					fErrs = fmt.Errorf("%v : %v", fErrs, err)
				} else {
					fErrs = err
				}
			} else {
				symSets = append(symSets, symset)
			}
		}
	}

	if fErrs != nil {
		return fmt.Errorf("failed to load symbol set : %v", fErrs)
	}

	var symbolSetsMap = make(map[string]symbolset.SymbolSet)
	for _, z := range symSets {
		// TODO check that x.Name doesn't already exist
		symbolSetsMap[z.Name] = z
	}
	mapperService.SymbolSets = symbolSetsMap

	return nil
}

func loadMapperHandler(w http.ResponseWriter, r *http.Request) {
	err := loadMappersFromDir(symbolSetFileArea)
	if err != nil {
		msg := err.Error()
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(symbolSetNames(mapperService.SymbolSets))
	if err != nil {
		msg := fmt.Sprintf("json marshalling error : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(j))

}

func listMapperHandler(w http.ResponseWriter, r *http.Request) {
	ss := symbolSetNames(mapperService.SymbolSets)
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

func mapperHelpHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<h1>Mapper</h1>
<h2>load</h2> Loads mappers from a pre-defined folder on the server: <code>` + symbolSetFileArea + `</code>. Example invocation:
<pre><a href="/mapper/load">/mapper/load</a></pre>

<h2>list</h2> Lists available mappers. Example invocation:
<pre><a href="/mapper/list">/mapper/list</a></pre>

<h2>map</h2> Maps a transcription from one symbolset to another. Example invocation:
<pre><a href="/mapper/map?from=sv-se_ws-sampa&to=sv-se_sampa_mary&trans=%22%22%20p%20O%20j%20.%20k%20@">/mapper/map?from=sv-se_ws-sampa&to=sv-se_sampa_mary&trans=%22%22%20p%20O%20j%20.%20k%20@</a></pre>

<h2>symbolset</h2> Lists content of a named symbolset. Example invocation:
<pre><a href="/mapper/symbolset?name=sv-se_ws-sampa">/mapper/symbolset?name=sv-se_ws-sampa</a></pre>

<h2>maptable</h2> Lists content of a maptable given two symbolset names. Example invocation:
<pre><a href="/mapper/maptable?from=sv-se_ws-sampa&to=sv-se_sampa_mary">/mapper/maptable?from=sv-se_ws-sampa&to=sv-se_sampa_mary</a></pre>

<h2>upload</h2> Upload symbol set file
<pre><a href="/mapper_upload">/mapper_upload</a></pre>		
		`

	fmt.Fprint(w, html)
}

func uploadMapperHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/mapper_upload.html")
}

func doUploadMapperHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("lexiconfileupload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
		return
	}

	clientUUID := r.FormValue("client_uuid")

	if "" == strings.TrimSpace(clientUUID) {
		msg := "doUploadMapperHandler got no client uuid"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Lifted from https://github.com/astaxie/build-web-application-with-golang/blob/master/de/04.5.md

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("doUploadMapperHandler failed reading file : %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fName := filepath.Join(symbolSetFileArea, handler.Filename)
	f, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE, 0755)
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

	f.Close()

	fmt.Fprintf(w, "%v", handler.Header)
}

type SymbolSetNames struct {
	SymbolSetNames []string `json:symbol_set_names`
}

type ByString []string

func (a ByString) Len() int           { return len(a) }
func (a ByString) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByString) Less(i, j int) bool { return a[i] < a[j] }

func symbolSetNames(sss map[string]symbolset.SymbolSet) SymbolSetNames {
	var ssNames []string
	for ss, _ := range sss {
		ssNames = append(ssNames, ss)
	}
	sort.Sort(ByString(ssNames))
	return SymbolSetNames{SymbolSetNames: ssNames}
}