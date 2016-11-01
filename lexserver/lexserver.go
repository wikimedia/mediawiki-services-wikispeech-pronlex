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

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
	"golang.org/x/net/websocket"
)

// TODO Split this file into packages + main?

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

	if "" != strings.TrimSpace(r.FormValue("pp")) {
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

	idS := r.FormValue("id")
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
	idS := r.FormValue("id")
	if idS == "" {
		msg := "deleteLexHander expected a lexicon id defined by 'id'"
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	id, _ := strconv.ParseInt(idS, 10, 64)

	uuid := r.FormValue("client_uuid")
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
	"page":                1,
	"pagelength":          1,
	"pp":                  1,
	"hasEntryValidation":  1,
}

var splitRE = regexp.MustCompile("[, ]")

func queryFromParams(r *http.Request) (dbapi.Query, error) {

	lexs := dbapi.RemoveEmptyStrings(
		splitRE.Split(r.FormValue("lexicons"), -1))
	words := dbapi.RemoveEmptyStrings(
		splitRE.Split(r.FormValue("words"), -1))
	lemmas := dbapi.RemoveEmptyStrings(
		splitRE.Split(r.FormValue("lemmas"), -1))

	wordLike := strings.TrimSpace(r.FormValue("wordlike"))
	wordRegexp := strings.TrimSpace(r.FormValue("wordregexp"))
	transcriptionLike := strings.TrimSpace(r.FormValue("transcriptionlike"))
	transcriptionRegexp := strings.TrimSpace(r.FormValue("transcriptionregexp"))
	partOfSpeechLike := strings.TrimSpace(r.FormValue("partofspeechlike"))
	partOfSpeechRegexp := strings.TrimSpace(r.FormValue("partofspeechregexp"))
	lemmaLike := strings.TrimSpace(r.FormValue("lemmalike"))
	lemmaRegexp := strings.TrimSpace(r.FormValue("lemmaregexp"))
	readingLike := strings.TrimSpace(r.FormValue("readinglike"))
	readingRegexp := strings.TrimSpace(r.FormValue("readingregexp"))
	paradigmLike := strings.TrimSpace(r.FormValue("paradigmlike"))
	paradigmRegexp := strings.TrimSpace(r.FormValue("paradigmregexp"))

	// If true, returns only entries with at least one EntryValidation issue
	hasEntryValidation := false
	if strings.ToLower(r.FormValue("hasEntryValidation")) == "true" {
		hasEntryValidation = true
	}
	// TODO report error if r.FormValue("page") != ""?
	// Silently sets deafault if no value, or faulty value
	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil {
		page = 0
		//log.Printf("failed to parse page parameter (using default value 0): %v", err)
	}

	// TODO report error if r.FormValue("pagelength") != ""?
	// Silently sets deafault if no value, or faulty value
	pageLength, err := strconv.ParseInt(r.FormValue("pagelength"), 10, 64)
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
		Page:                page,
		PageLength:          pageLength,
		HasEntryValidation:  hasEntryValidation,
	}

	return q, err
}

func updateEntryHandler(w http.ResponseWriter, r *http.Request) {
	entryJSON := r.FormValue("entry")
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
				log.Print("keepClientsAlive closed socket to client %s", client)
			}
		}
		webSocks.Unlock()
	}
}

func exportLexiconHandler(w http.ResponseWriter, r *http.Request) {
	lexiconID, err := strconv.ParseInt(r.FormValue("id"), 10, 64)
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
	clientUUID := r.FormValue("client_uuid")
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
	wsW := line.FileWriter{wsFmt, gz}
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

	clientUUID := r.FormValue("client_uuid")

	lexiconID, err := strconv.ParseInt(r.FormValue("lexicon_id"), 10, 64)
	if err != nil {
		msg := "lexiconFileUploadHandler got no lexicon id"
		fmt.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}
	lexiconName := r.FormValue("lexicon_name")
	log.Printf("lexiconFileUploadHandler: incoming db lexicon name: %v\n", lexiconName)
	symbolSetName := r.FormValue("symbolset_name")
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
	fName := r.FormValue("file")
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

	// load symbol set mappers
	err = loadSymbolSets(symbolSetFileArea)
	ff("Failed to load symbol sets : %v", err)
	log.Printf("Loaded symbol set mappers from dir %s", symbolSetFileArea)

	err = loadValidators(symbolSetFileArea)
	ff("Failed to load validators : %v", err)
	log.Printf("Loaded validators : %v", validatorNames())

	// static
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/ipa_table.txt", ipaTableHandler)

	// Old junk
	// Temp, testing
	//http.HandleFunc("/lexicon/listknown", listKnownHandler)

	// function calls
	http.HandleFunc("/lexicon", lexiconHelpHandler)
	http.HandleFunc("/lexicon/list", listLexsWithEntryCountHandler)
	http.HandleFunc("/listlexicons", apiChangedHandler("use /lexicon/list instead"))
	http.HandleFunc("/lexicon/listlexicons", apiChangedHandler("use /lexicon/list instead"))

	http.HandleFunc("/lexicon/stats", lexiconStatsHandler)
	http.HandleFunc("/lexiconstats", apiChangedHandler("use /lexicon/stats instead"))

	http.HandleFunc("/lexicon/lookup", lexLookUpHandler)
	http.HandleFunc("/lexlookup", apiChangedHandler("use /lexicon/lookup instead"))

	http.HandleFunc("/lexicon/addentry", addEntryHandler)
	http.HandleFunc("/addentry", apiChangedHandler("use /lexicon/addentry instead"))

	http.HandleFunc("/lexicon/updateentry", updateEntryHandler)
	http.HandleFunc("/updateentry", apiChangedHandler("/lexicon/updateentry instead"))

	http.HandleFunc("/lexicon/validate", lexiconValidateHandler)
	http.HandleFunc("/lexicon/do_validate", lexiconRunValidateHandler)

	http.HandleFunc("/validation/validateentry", validateEntryHandler)
	http.HandleFunc("/validateentry", apiChangedHandler("/validation/validateentry"))
	http.HandleFunc("/validation/validateentries", validateEntriesHandler)
	http.HandleFunc("/validateentries", apiChangedHandler("/validation/validateentries"))
	http.HandleFunc("/validation", validationHelpHandler)
	http.HandleFunc("/validation/list", listValidationHandler)
	http.HandleFunc("/validation/stats", validationStatsHandler)

	//http.HandleFunc("/download", downloadFileHandler)

	// admin pages/calls

	//http.HandleFunc("/admin_lex_definition.html", adminLexDefinitionHandler)
	http.HandleFunc("/admin/lex_import", adminLexImportHandler)
	http.HandleFunc("/admin/lex_do_import", adminDoLexImportHandler)
	//http.HandleFunc("/admin/admin.html", adminAdminHandler)
	http.HandleFunc("/admin", adminHelpHandler)
	//http.HandleFunc("/admin", adminHandler)
	//http.HandleFunc("/admin/createlex", adminCreateLexHandler)
	//http.HandleFunc("/admin/editsymbolset", adminEditSymbolSetHandler)
	//http.HandleFunc("/admin/listsymbolset", listSymbolSetHandler)
	//http.HandleFunc("/admin/savesymbolset", saveSymbolSetHandler)
	//http.HandleFunc("/admin/insertorupdatelexicon", insertOrUpdateLexHandler)
	http.HandleFunc("/admin/deletelexicon", deleteLexHandler)
	http.HandleFunc("/admin/superdeletelexicon", superDeleteLexHandler)

	// Sqlite3 ANALYZE command in some instances make search quicker,
	// but it takes a while to perform
	//http.HandleFunc("/admin/sqlite3_analyze", sqlite3AnalyzeHandler)

	//http.HandleFunc("/websocktest", webSockTestHandler)

	http.Handle("/websockreg", websocket.Handler(webSockRegHandler))

	//http.HandleFunc("/admin/lexiconfileupload", lexiconFileUploadHandler)
	//http.HandleFunc("/admin/exportlexicon", exportLexiconHandler)

	// TODO Split this main func into several files

	http.HandleFunc("/symbolset", symbolSetHelpHandler)
	http.HandleFunc("/symbolset/list", listSymbolSetsHandler)
	http.HandleFunc("/symbolset/delete", deleteSymbolSetHandler)
	http.HandleFunc("/symbolset/reload", reloadSymbolSetHandler)
	http.HandleFunc("/symbolset/symbolset", symbolSetHandler)
	http.HandleFunc("/symbolset/upload", uploadSymbolSetHandler)
	http.HandleFunc("/symbolset/do_upload", doUploadSymbolSetHandler)

	http.HandleFunc("/mapper", mapperHelpHandler)
	http.HandleFunc("/mapper/list", listMappersHandler)
	http.HandleFunc("/mapper/map", mapMapperHandler)
	http.HandleFunc("/mapper/maptable", mapTableMapperHandler)

	// TODO Why this http.StripPrefix? Looks odd.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	//http.Handle("/download_area/", http.StripPrefix("/download_area/", http.FileServer(http.Dir("./download_area"))))

	// Pinging connected websocket clients
	go keepClientsAlive()

	log.Print("lexserver: listening on port ", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
