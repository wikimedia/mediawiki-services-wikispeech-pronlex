package main

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-errors/errors"
	"github.com/gorilla/mux"
	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
	"golang.org/x/net/websocket"
)

// TODO Split this file into packages + main?

func getParam(paramName string, r *http.Request) string {
	res := r.FormValue(paramName)
	if res != "" {
		return res
	}
	vars := mux.Vars(r)
	return vars[paramName]
}

func (rout *subRouter) addHandler(handler urlHandler) {
	rout.router.HandleFunc(handler.url, handler.handler)
	rout.handlers = append(rout.handlers, handler)
}

type subRouter struct {
	root     string
	router   *mux.Router
	handlers []urlHandler
}

func newSubRouter(rout *mux.Router, root string) *subRouter {
	var res = subRouter{
		router: rout.PathPrefix(root).Subrouter(),
		root:   root,
	}

	helpHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		html := "<h1>" + res.root + "</h1>"
		for _, handler := range res.handlers {
			html = html + handler.helpHtml(res.root)
		}
		fmt.Fprint(w, html)
	}

	res.router.HandleFunc("/", helpHandler)
	return &res
}

type urlHandler struct {
	name     string
	handler  func(w http.ResponseWriter, r *http.Request)
	url      string
	help     string
	examples []string
}

func (h urlHandler) helpHtml(root string) string {
	s := "<h2>" + h.name + "</h2> " + h.help
	for _, x := range h.examples {
		url := root + x
		s = s + `<pre><a href="` + url + `">` + url + `</a></pre>`
	}
	return s
}

// protect: use this call in handlers to catch 'panic' and stack traces and returning a general error to the calling client
func protect(w http.ResponseWriter) {
	if r := recover(); r != nil {
		defer http.Error(w, fmt.Sprintf("%s", "Internal server error"), http.StatusInternalServerError)
		fmt.Println(errors.Wrap(r, 2).ErrorStack())
		// TODO: log the actual error to a server log file (but do not return to client)
	}
}

// TODO remove calls to this, add error handling
func ff(f string, err error) {
	if err != nil {
		log.Fatalf(f, err)
	}
}

// TODO should go into config file
var uploadFileArea = filepath.Join(".", "upload_area")
var downloadFileArea = filepath.Join(".", "download_area")
var symbolSetFileArea = filepath.Join(".", "symbol_set_file_area")

// TODO config stuff
func init() {
	// TODO sane error handling

	// If the upload area dir doesn't exist, create it
	if _, err := os.Stat(uploadFileArea); err != nil {
		if os.IsNotExist(err) {
			err2 := os.Mkdir(uploadFileArea, 0755)
			if err2 != nil {
				fmt.Printf("lexserver.init: failed to create %s : %v", uploadFileArea, err2)
			}
		} else {
			fmt.Printf("lexserver.init: peculiar error : %v", err)
		}
	} // else: already exists, hopefullly

	// If the download area dir doesn't exist, create it
	if _, err := os.Stat(downloadFileArea); err != nil {
		if os.IsNotExist(err) {
			err2 := os.Mkdir(downloadFileArea, 0755)
			if err2 != nil {
				fmt.Printf("lexserver.init: failed to create %s : %v", downloadFileArea, err2)
			}
		} else {
			fmt.Printf("lexserver.init: peculiar error : %v", err)
		}
	} // else: already exists, hopefullly

	// If the symbol set dir doesn't exist, create it
	if _, err := os.Stat(symbolSetFileArea); err != nil {
		if os.IsNotExist(err) {
			err2 := os.Mkdir(symbolSetFileArea, 0755)
			if err2 != nil {
				fmt.Printf("lexserver.init: failed to create %s : %v", symbolSetFileArea, err2)
			}
		} else {
			fmt.Printf("lexserver.init: peculiar error : %v", err)
		}
	} // else: already exists, hopefullly

}

// TODO remove pretty-print option, since you can use the JSONView plugin to Chrome instead
// pretty print if the URL paramer 'pp' has a value
func marshal(v interface{}, r *http.Request) ([]byte, error) {

	if "" != strings.TrimSpace(getParam("pp", r)) {
		return json.MarshalIndent(v, "", "  ")
	}

	return json.Marshal(v)
}

func ipaTableHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/ipa_table.txt")
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Serving index file for '%v'\n", r.URL.Path)
	http.ServeFile(w, r, "./static/index.html")
}

func deleteLexHandler(w http.ResponseWriter, r *http.Request) {

	idS := getParam("id", r)
	if idS == "" {
		msg := "deleteLexHander expected a lexicon id defined by 'id'"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id, _ := strconv.ParseInt(idS, 10, 64)

	err := dbapi.DeleteLexicon(db, id)
	if err != nil {
		log.Printf("lexserver.deleteLexHandler got error : %v\n", err)
		http.Error(w, fmt.Sprintf("failed deleting lexicon : %v", err), http.StatusExpectationFailed)
		return
	}
}

func superDeleteLexHandler(w http.ResponseWriter, r *http.Request) {
	// Aha! Turns out that Go treats POST and GET the same way, as I understand it.
	// No need for checking whether GET or POST, as far as I understand.
	idS := getParam("id", r)
	if idS == "" {
		msg := "deleteLexHander expected a lexicon id defined by 'id'"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id, _ := strconv.ParseInt(idS, 10, 64)

	uuid := getParam("client_uuid", r)
	log.Println("lexserver.superDeleteLexHandler was called")
	messageToClientWebSock(uuid, fmt.Sprintf("Super delete was called. This may take quite a while. Lexicon id %d", id))
	err := dbapi.SuperDeleteLexicon(db, id)
	if err != nil {

		http.Error(w, fmt.Sprintf("failed super deleting lexicon : %v", err), http.StatusExpectationFailed)
		return
	}

	messageToClientWebSock(uuid, fmt.Sprintf("Done deleting lexicon with id %d", id))
}

func sqlite3AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("lexserver: running the Sqlite3 ANALYZE command...")
	_, err := db.Exec("ANALYZE")

	if err != nil {
		log.Printf("Failed to exec ANALYZE %v", err)
		http.Error(w, fmt.Sprintf("/admin/sqlite3_analyze failed : %v", err), http.StatusInternalServerError)
		return
	}
	log.Print("... done!\n")
	w.Write([]byte("OK"))
}

// TODO report unused URL parameters

// TODO Gör konstanter som kan användas istället för strängar
var knownParams = map[string]int{
	"lexicons":            1,
	"words":               1,
	"lemmas":              1,
	"wordlike":            1,
	"wordregexp":          1,
	"transcriptionlike":   1,
	"transcriptionregexp": 1,
	"partofspeechlike":    1,
	"partofspeechregexp":  1,
	"lemmalike":           1,
	"lemmaregexp":         1,
	"readinglike":         1,
	"readingregexp":       1,
	"paradigmlike":        1,
	"paradigmregexp":      1,
	"hasEntryValidation":  1,
	"page":                1,
	"pagelength":          1,
	"pp":                  1,
}

// list of values to the same param splits on comma and/or space
var splitRE = regexp.MustCompile("[, ]+")

func queryFromParams(r *http.Request) (dbapi.Query, error) {

	lexs := dbapi.RemoveEmptyStrings(
		splitRE.Split(getParam("lexicons", r), -1))
	words := dbapi.RemoveEmptyStrings(
		splitRE.Split(getParam("words", r), -1))
	lemmas := dbapi.RemoveEmptyStrings(
		splitRE.Split(getParam("lemmas", r), -1))

	wordLike := strings.TrimSpace(getParam("wordlike", r))
	wordRegexp := strings.TrimSpace(getParam("wordregexp", r))
	transcriptionLike := strings.TrimSpace(getParam("transcriptionlike", r))
	transcriptionRegexp := strings.TrimSpace(getParam("transcriptionregexp", r))
	partOfSpeechLike := strings.TrimSpace(getParam("partofspeechlike", r))
	partOfSpeechRegexp := strings.TrimSpace(getParam("partofspeechregexp", r))
	lemmaLike := strings.TrimSpace(getParam("lemmalike", r))
	lemmaRegexp := strings.TrimSpace(getParam("lemmaregexp", r))
	readingLike := strings.TrimSpace(getParam("readinglike", r))
	readingRegexp := strings.TrimSpace(getParam("readingregexp", r))
	paradigmLike := strings.TrimSpace(getParam("paradigmlike", r))
	paradigmRegexp := strings.TrimSpace(getParam("paradigmregexp", r))
	var entryStatus []string
	if "" != getParam("entrystatus", r) {
		entryStatus = splitRE.Split(getParam("entrystatus", r), -1)
	}
	// If true, returns only entries with at least one EntryValidation issue
	hasEntryValidation := false
	if strings.ToLower(getParam("hasEntryValidation", r)) == "true" {
		hasEntryValidation = true
	}
	// TODO report error if getParam("page", r) != ""?
	// Silently sets deafault if no value, or faulty value
	page, err := strconv.ParseInt(getParam("page", r), 10, 64)
	if err != nil {
		page = 0
		//log.Printf("failed to parse page parameter (using default value 0): %v", err)
	}

	// TODO report error if getParam("pagelength", r) != ""?
	// Silently sets deafault if no value, or faulty value
	pageLength, err := strconv.ParseInt(getParam("pagelength", r), 10, 64)
	if err != nil {
		pageLength = 25
	}

	dbLexs, err := dbapi.GetLexicons(db, lexs)

	q := dbapi.Query{
		Lexicons:            dbLexs,
		Words:               words,
		WordLike:            wordLike,
		WordRegexp:          wordRegexp,
		TranscriptionLike:   transcriptionLike,
		TranscriptionRegexp: transcriptionRegexp,
		PartOfSpeechLike:    partOfSpeechLike,
		PartOfSpeechRegexp:  partOfSpeechRegexp,
		Lemmas:              lemmas,
		LemmaLike:           lemmaLike,
		LemmaRegexp:         lemmaRegexp,
		ReadingLike:         readingLike,
		ReadingRegexp:       readingRegexp,
		ParadigmLike:        paradigmLike,
		ParadigmRegexp:      paradigmRegexp,
		EntryStatus:         entryStatus,
		Page:                page,
		PageLength:          pageLength,
		HasEntryValidation:  hasEntryValidation,
	}

	return q, err
}

// Remove initial and trailing " or ' from string
func delQuote(s string) string {
	res := s
	res = strings.TrimPrefix(res, `"`)
	res = strings.TrimPrefix(res, `'`)
	res = strings.TrimSuffix(res, `"`)
	res = strings.TrimSuffix(res, `'`)

	return strings.TrimSpace(res)
}

func updateEntryHandler(w http.ResponseWriter, r *http.Request) {
	entryJSON := getParam("entry", r)
	//body, err := ioutil.ReadAll(r.Body)
	var e lex.Entry
	err := json.Unmarshal([]byte(entryJSON), &e)
	if err != nil {
		log.Printf("lexserver: Failed to unmarshal json: %v", err)
		http.Error(w, fmt.Sprintf("failed to process incoming Entry json : %v", err), http.StatusInternalServerError)
		return
	}

	// Underscore below matches bool indicating if any update has taken place. Return this info?
	res, _, err2 := dbapi.UpdateEntry(db, e)
	if err2 != nil {
		log.Printf("lexserver: Failed to update entry : %v", err2)
		http.Error(w, fmt.Sprintf("failed to update Entry : %v", err2), http.StatusInternalServerError)
		return
	}

	res0, err3 := json.Marshal(res)
	if err3 != nil {
		log.Printf("lexserver: Failed to marshal entry : %v", err3)
		http.Error(w, fmt.Sprintf("failed return updated Entry : %v", err3), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprint(w, string(res0))
}

func adminAdminHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/admin/admin.html")
}

func adminHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/admin/index.html")
}

func adminLexDefinitionHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/admin_lex_definition.html")
}

func adminLexImportHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/admin/lex_import.html")
}

func lexiconValidateHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/lexicon/validate.html")
}

func adminCreateLexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/admin/create_lex.html")
}

func adminEditSymbolSetHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/admin/edit_symbolset.html")
}

var wsChan = make(chan string)

// See https://blog.golang.org/go-maps-in-action "Concurrency"
var webSocks = struct {
	sync.RWMutex
	clients map[string]*websocket.Conn
}{clients: make(map[string]*websocket.Conn)}

func deleteWebSocketClient(id string) {
	webSocks.Lock()
	//defer webSocks.Unlock()
	delete(webSocks.clients, id)
	webSocks.Unlock()
}

func messageToClientWebSock(clientUUID string, msg string) {
	if strings.TrimSpace(clientUUID) != "" {
		webSocks.Lock()
		if ws, ok := webSocks.clients[clientUUID]; ok {
			websocket.Message.Send(ws, msg)
		} else {
			log.Printf("messageToClientWebSock called with unknown UUID string '%s'", clientUUID)
		}
		webSocks.Unlock()
	} else {
		log.Printf("messageToClientWebSock called with empty UUID string and message '%s'", msg)
	}
}

func webSockRegHandler(ws *websocket.Conn) {
	var id string
	for {
		var msg string
		err := websocket.Message.Receive(ws, &msg)
		if err != nil {
			log.Printf("webSockRegHandler error : %v\n", err)
			log.Printf("webSockRegHandler removing socket with id %s", id)
			deleteWebSocketClient(id)
			id = ""
			break
		}

		log.Printf("webSockRegHandler: " + msg)

		var pref = "CLIENT_ID: "
		//var id string
		if strings.HasPrefix(msg, pref) {
			id = strings.TrimPrefix(strings.TrimSpace(msg), pref)
			webSocks.Lock()
			//defer webSocks.Unlock()
			webSocks.clients[id] = ws
			webSocks.Unlock()
			log.Printf("Processed id: %v", id)
		}

		log.Printf("id is now: %v", id)

		var reply = "HI THERE! " + time.Now().Format("Mon, 02 Jan 2006 15:04:05 PST")
		err = websocket.Message.Send(ws, reply)
		if err != nil {
			log.Printf("webSockRegHandler error : %v\n", err)
			deleteWebSocketClient(id)
			id = ""
			break
		}
	}
}

func keepClientsAlive() {
	c := time.Tick(67 * time.Second)
	for _ = range c {

		webSocks.Lock()
		//log.Printf("keepClientsAlive: pinging number of clients: %v\n", len(webSocks.clients))
		//defer webSocks.Unlock()
		for client, ws := range webSocks.clients {

			err := websocket.Message.Send(ws, "WS_KEEPALIVE") //+time.Now().Format("Mon, 02 Jan 2006 15:04:05 PST"))
			if err != nil {
				log.Printf("keepClientsAlive got error from websocket send : %v", err)
				delete(webSocks.clients, client)
				log.Printf("keepClientsAlive closed socket to client %s", client)
			}
		}
		webSocks.Unlock()
	}
}

func exportLexiconHandler(w http.ResponseWriter, r *http.Request) {
	lexiconID, err := strconv.ParseInt(getParam("id", r), 10, 64)
	if err != nil {
		msg := "exportLexiconHandler got no lexicon id"
		fmt.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	lexicon, err := dbapi.LexiconFromID(db, lexiconID)
	if err != nil {
		msg := fmt.Sprintf("exportLexiconHandler failed to get lexicon from id : %v", err)
		log.Println(msg)
		// TODO this might not be a proper error: the client could simply have asked for a lexicon id that doesn't exist.
		// Handle more gracefully (but for now, let's crash).
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// If client sends UUID, messages can be written to client socket
	clientUUID := getParam("client_uuid", r)
	messageToClientWebSock(clientUUID, fmt.Sprintf("This will take a while. Starting to export lexicon %s", lexicon.Name))

	// local output file
	fName := filepath.Join(downloadFileArea, lexicon.Name+".txt.gz")
	f, err := os.Create(fName)
	//bf := bufio.NewWriter(f)
	gz := gzip.NewWriter(f)
	//defer gz.Flush()
	defer gz.Close()
	// Query that returns all entries of lexicon
	ls := []dbapi.Lexicon{dbapi.Lexicon{ID: lexicon.ID}}
	q := dbapi.Query{Lexicons: ls}

	log.Printf("Query for exporting: %v", q)

	//nstFmt, err := line.NewNST()
	wsFmt, err := line.NewWS()
	if err != nil {
		log.Fatal(err)
		http.Error(w, "exportLexicon failed to create line writer", http.StatusInternalServerError)
		return
	}
	wsW := line.FileWriter{Parser: wsFmt, Writer: gz}
	dbapi.LookUp(db, q, wsW)
	defer gz.Close()
	gz.Flush()
	messageToClientWebSock(clientUUID, fmt.Sprintf("Done exporting lexicon %s to %s", lexicon.Name, fName))

	msg := fmt.Sprintf("Lexicon exported to '%s'", fName)
	log.Print(msg)
	fmt.Fprint(w, filepath.Base(fName))
}

func lexiconFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("lexiconFileUploadHandler: method: ", r.Method)
	if r.Method != "POST" {
		http.Error(w, fmt.Sprintf("lexiconfileupload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
		return
	}

	clientUUID := getParam("client_uuid", r)

	lexiconID, err := strconv.ParseInt(getParam("lexicon_id", r), 10, 64)
	if err != nil {
		msg := "lexiconFileUploadHandler got no lexicon id"
		fmt.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	lexiconName := getParam("lexicon_name", r)
	log.Printf("lexiconFileUploadHandler: incoming db lexicon name: %v\n", lexiconName)
	symbolSetName := getParam("symbolset_name", r)
	log.Printf("lexiconFileUploadHandler: incoming db symbol set name: %v\n", symbolSetName)

	if "" == strings.TrimSpace(clientUUID) {
		msg := "lexiconFileUploadHandler got no client uuid"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if "" == strings.TrimSpace(lexiconName) {
		msg := "lexiconFileUploadHandler got no lexicon name"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	if "" == strings.TrimSpace(symbolSetName) {
		msg := "lexiconFileUploadHandler got no symbolset name"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Lifted from https://github.com/astaxie/build-web-application-with-golang/blob/master/de/04.5.md

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("upload_file")
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("lexiconFileUploadHandler failed reading file : %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fName := filepath.Join(uploadFileArea, handler.Filename)
	f, err := os.OpenFile(fName, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("lexiconFileUploadHandler failed opening local output file : %v", err), http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		msg := fmt.Sprintf("lexiconFileUploadHandler failed copying local output file : %v", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// TODO temporarily try to directly load the uploaded file
	// into the database.  However, this should be a step of its
	// own: first upload file, validate it, etc, and then make it
	// possible to load into the db, not blindely just adding
	// stuff.  Should check if there are duplicates:
	// words+transcription in upload text file already present in
	// db.
	f.Close()
	loadLexiconFileIntoDB(clientUUID, lexiconID, lexiconName, symbolSetName, fName)

	fmt.Fprintf(w, "%v", handler.Header)
}

// TODO temporary test thingy
// TODO This guy should somehow report back what it's doing to the client.
// TODO Return some sort of result? Stats?
// TODO Set 'status' value for imported entries (now hard-wired to 'import' below)
// TODO Set 'source' value for imported entries (now hard-wired to 'unknown' below)
func loadLexiconFileIntoDB(clientUUID string, lexiconID int64, lexiconName string, symbolSetName string, uploadFileName string) error {
	fmt.Printf("clientUUID: %s\n", clientUUID)
	fmt.Printf("lexid: %v\n", lexiconID)
	fmt.Printf("lexiconName: %v\n", lexiconName)
	fmt.Printf("symbolSetName: %v\n", symbolSetName)
	fmt.Printf("uploadFile: %v\n", uploadFileName)

	fh, err := os.Open(uploadFileName)
	if err != nil {
		var msg = fmt.Sprintf("loadLexiconFileIntoDB failed to open file : %v", err)
		messageToClientWebSock(clientUUID, msg)
		return fmt.Errorf("loadLexiconFileIntoDB failed to open file : %v", err)
	}

	s := bufio.NewScanner(fh)
	// for s.Scan() {
	// 	l := s.Text()
	// 	fmt.Println(l)
	// }

	// ---------------------------->

	//nstFmt, err := line.NewNST()
	wsFmt, err := line.NewWS()
	if err != nil {
		//log.Fatal(err)
		var msg = fmt.Sprintf("lexserver failed to instantiate lexicon line parser : %v", err)
		messageToClientWebSock(clientUUID, msg)
		return fmt.Errorf("lexserver failed to instantiate lexicon line parser : %v", err)
	}
	lexicon := dbapi.Lexicon{ID: lexiconID, Name: lexiconName, SymbolSetName: symbolSetName}

	msg := fmt.Sprintf("Trying to load file: %s", uploadFileName)
	messageToClientWebSock(clientUUID, msg)
	log.Print(msg)

	n := 0
	var eBuf []lex.Entry
	for s.Scan() {
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
		l := s.Text()
		e, err := wsFmt.ParseToEntry(l)
		if err != nil {
			log.Fatal(err)
		}
		eBuf = append(eBuf, e)
		n++
		if n%10000 == 0 {
			_, err = dbapi.InsertEntries(db, lexicon, eBuf)
			if err != nil {
				log.Fatal(err)
			}
			eBuf = make([]lex.Entry, 0)
			//fmt.Printf("\rLines read: %d               \r", n)
			msg2 := fmt.Sprintf("Lines so far: %d", n)
			messageToClientWebSock(clientUUID, msg2)
			fmt.Println(msg2)
		}
	}
	dbapi.InsertEntries(db, lexicon, eBuf) // flushing the buffer

	_, err = db.Exec("ANALYZE")
	if err != nil {
		log.Fatal(err)
	}

	msg3 := fmt.Sprintf("Lines read:\t%d", n)
	messageToClientWebSock(clientUUID, msg3)
	log.Println(msg3)
	// TODO Copied from addNSTLexToDB
	// <----------------------------

	if err := s.Err(); err != nil {
		msg4 := fmt.Sprintf("lexserver failed to instantiate lexicon line parser : %v", err)
		messageToClientWebSock(clientUUID, msg4)
		return fmt.Errorf(msg4)
	}

	// TODO
	return nil
}

func downloadFileHandler(w http.ResponseWriter, r *http.Request) {
	fName := getParam("file", r)
	if fName == "" {
		msg := fmt.Sprint("downloadFileHandler got empty 'file' param")
		log.Println(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	path := filepath.Join(".", downloadFileArea, fName)
	log.Printf("downloadFileHandler: file path: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		msg := fmt.Sprintf("download: no such file '%s'", fName)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fName)
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	http.ServeFile(w, r, path)
}

var db *sql.DB

func keepAlive(wsC chan string) {
	c := time.Tick(57 * time.Second)
	for _ = range c {
		wsC <- "WS_KEEPALIVE"
	}
}

func apiChangedHandler(msg string) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, fmt.Sprintf("the API has changed: %s\n", msg), http.StatusBadRequest)
		return
	}
}

func searchDemoHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "../web/lexsearch/index.html")
}

func loadSymbolSetFile(fName string) (symbolset.SymbolSet, error) {
	return symbolset.LoadSymbolSet(fName)
}

func main() {

	dbFile := "./pronlex.db"
	port := ":8787"

	if len(os.Args) > 3 || len(os.Args) == 2 {
		log.Println("Usages:")
		log.Println("$ go run lexserver.go <SQLITE DB FILE> <PORT>")
		log.Println("$ go run lexserver.go")
		log.Println("  - defaults to db file " + dbFile + ", port " + port)
		os.Exit(1)
	} else if len(os.Args) == 3 {
		dbFile = os.Args[1] // "./pronlex.db"
		port = os.Args[2]   //":8787"
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	rout := mux.NewRouter().StrictSlash(true)

	var err error // återanvänds för alla fel

	// kolla att db-filen existerar
	_, err = os.Stat(dbFile)
	ff("lexserver: Cannot find db file. %v", err)

	dbapi.Sqlite3WithRegex()

	log.Print("lexserver: connecting to Sqlite3 db ", dbFile)
	db, err = sql.Open("sqlite3_with_regexp", dbFile)
	ff("Failed to open dbfile %v", err)
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	ff("Failed to exec PRAGMA call %v", err)
	_, err = db.Exec("PRAGMA case_sensitive_like=ON")
	ff("Failed to exec PRAGMA call %v", err)
	_, err = db.Exec("PRAGMA journal_mode=WAL")
	ff("Failed to exec PRAGMA call %v", err)
	//_, err = db.Exec("PRAGMA busy_timeout=500") // doesn't seem to do the trick
	//ff("Failed to exec PRAGMA call %v", err)
	db.SetMaxOpenConns(1) // to avoid locking errors (but it makes it slow...?) https://github.com/mattn/go-sqlite3/issues/274

	// load symbol set mappers
	err = loadSymbolSets(symbolSetFileArea)
	ff("Failed to load symbol sets from dir "+symbolSetFileArea+" : %v", err)
	log.Printf("lexserver: Loaded symbol set mappers from dir %s", symbolSetFileArea)

	err = loadValidators(symbolSetFileArea)
	ff("Failed to load validators : %v", err)
	log.Printf("lexserver: Loaded validators : %v", validatorNames())

	// static
	rout.HandleFunc("/", indexHandler)
	rout.HandleFunc("/favicon.ico", faviconHandler)
	rout.HandleFunc("/ipa_table.txt", ipaTableHandler)

	// function calls
	rout.HandleFunc("/lexicon", lexiconHelpHandler)
	rout.HandleFunc("/lexicon/list", listLexsWithEntryCountHandler)
	rout.HandleFunc("/lexicon/list_current_entry_statuses", listCurrentEntryStatuses)
	rout.HandleFunc("/lexicon/list_all_entry_statuses", listAllEntryStatuses)
	rout.HandleFunc("/lexicon/info", lexInfoHandler)
	rout.HandleFunc("/lexicon/stats", lexiconStatsHandler)
	rout.HandleFunc("/lexicon/lookup", lexLookUpHandler)
	rout.HandleFunc("/lexicon/addentry", addEntryHandler)
	rout.HandleFunc("/lexicon/updateentry", updateEntryHandler)

	// defined in file move_new_entries_handler.go.
	rout.HandleFunc("/lexicon/move_new_entries", moveNewEntriesHandler)
	rout.HandleFunc("/lexicon/validate", lexiconValidateHandler)
	rout.HandleFunc("/lexicon/do_validate", lexiconRunValidateHandler)

	rout.HandleFunc("/validation", validationHelpHandler)
	rout.HandleFunc("/validation/validateentry", validateEntryHandler)
	rout.HandleFunc("/validation/validateentries", validateEntriesHandler)
	rout.HandleFunc("/validation/list", listValidationHandler)
	rout.HandleFunc("/validation/has_validator", hasValidatorHandler)
	rout.HandleFunc("/validation/stats", validationStatsHandler)

	// admin pages/calls
	rout.HandleFunc("/admin/lex_import", adminLexImportHandler)
	rout.HandleFunc("/admin/lex_do_import", adminDoLexImportHandler)
	rout.HandleFunc("/admin", adminHelpHandler)
	rout.HandleFunc("/admin/deletelexicon", deleteLexHandler)
	rout.HandleFunc("/admin/superdeletelexicon", superDeleteLexHandler)

	// Sqlite3 ANALYZE command in some instances make search quicker,
	// but it takes a while to perform
	//rout.HandleFunc("/admin/sqlite3_analyze", sqlite3AnalyzeHandler)

	rout.Handle("/websockreg", websocket.Handler(webSockRegHandler))

	symbolset := newSubRouter(rout, "/symbolset")
	symbolset.addHandler(symbolsetList)
	symbolset.addHandler(symbolsetDelete)
	symbolset.addHandler(symbolsetContent)
	symbolset.addHandler(symbolsetReloadOne)
	symbolset.addHandler(symbolsetReloadAll)
	symbolset.addHandler(symbolsetUpload)
	rout.HandleFunc("/symbolset/do_upload", doUploadSymbolSetHandler) // hidden -- not part of API

	mapper := newSubRouter(rout, "/mapper")
	mapper.addHandler(mapperList)
	mapper.addHandler(mapperMap)
	mapper.addHandler(mapperMaptable)

	// typescript experiments
	rout.HandleFunc("/search_demo", searchDemoHandler)

	r0 := http.StripPrefix("/lexsearch/built/", http.FileServer(http.Dir("../web/lexsearch/built/")))
	rout.Handle("/lexsearch/built/", r0)

	r1 := http.StripPrefix("/lexsearch/externals/", http.FileServer(http.Dir("../web/lexsearch/externals/")))
	rout.Handle("/lexsearch/externals/", r1)

	// serve static folder
	rout.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	fmt.Println("Serving urls:")
	rout.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		fmt.Println(t)
		return nil
	})

	// Pinging connected websocket clients
	go keepClientsAlive()

	log.Print("lexserver: listening on port ", port)
	log.Fatal(http.ListenAndServe(port, rout))

}
