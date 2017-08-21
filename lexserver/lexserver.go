package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
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
	"github.com/stts-se/pronlex/symbolset"
	"golang.org/x/net/websocket"
)

func getParam(paramName string, r *http.Request) string {
	res := r.FormValue(paramName)
	if res != "" {
		return res
	}
	vars := mux.Vars(r)
	return vars[paramName]
}

func getLexRefParam(r *http.Request) (lex.LexRef, error) {
	lexRefS := getParam("lexicon_name", r)
	if strings.TrimSpace(lexRefS) == "" {
		msg := "input param <lexicon_name> must not be empty"
		return lex.LexRef{}, fmt.Errorf(msg)
	}
	lexRef, err := lex.ParseLexRef(lexRefS)
	if err != nil {
		return lex.LexRef{}, err
	}
	return lexRef, nil
}

func (rout *subRouter) addHandler(handler urlHandler) {
	rout.router.HandleFunc(handler.url, handler.handler)
	rout.handlers = append(rout.handlers, handler)
}

type subRouter struct {
	root     string
	router   *mux.Router
	handlers []urlHandler
	desc     string
}

var subRouters []*subRouter

type urlHandler struct {
	name     string
	handler  func(w http.ResponseWriter, r *http.Request)
	url      string
	help     string
	examples []string
}

// TODO: Neat URL encoding...
func urlEnc(url string) string {
	return strings.Replace(strings.Replace(strings.Replace(url, " ", "%20", -1), "\n", "", -1), `"`, "%22", -1)
}

func (h urlHandler) helpHTML(root string) string {
	s := "<h2>" + h.name + "</h2> " + h.help
	if strings.Contains(h.url, "{") {
		s = s + `<p>API URL: <code>` + root + h.url + `</code></p>`
	}
	if len(h.examples) > 0 {
		//s = s + `<p>Example invocation:`
		for _, x := range h.examples {
			urlPretty := root + x
			url := root + urlEnc(x)
			s = s + `<pre><a href="` + url + `">` + urlPretty + `</a></pre>`
		}
		//s = s + "</p>"
	}
	return s
}
func isHandeledPage(url string) bool {
	for _, sub := range subRouters {
		if sub.root == url || sub.root+"/" == url {
			return true
		}
		for _, handler := range sub.handlers {
			if sub.root+handler.url == url {
				return true
			}
		}
	}
	return false
}

var initialSlashRe = regexp.MustCompile("^/")

func removeInitialSlash(url string) string {
	return initialSlashRe.ReplaceAllString(url, "")
}

func (sr subRouter) handlerExamples() []string {
	res := []string{}
	for _, handler := range sr.handlers {
		for _, example := range handler.examples {
			res = append(res, sr.root+example)
		}
	}
	return res
}

func newSubRouter(rout *mux.Router, root string, description string) *subRouter {
	var res = subRouter{
		router: rout.PathPrefix(root).Subrouter(),
		root:   root,
		desc:   description,
	}

	helpHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		html := "<h1>" + removeInitialSlash(res.root) + "</h1> <em>" + res.desc + "</em>"
		for _, handler := range res.handlers {
			html = html + handler.helpHTML(res.root)
		}
		fmt.Fprint(w, html)
	}

	res.router.HandleFunc("/", helpHandler)
	subRouters = append(subRouters, &res)
	return &res
}

// protect: use this call in handlers to catch 'panic' and stack traces and returning a general error to the calling client
func protect(w http.ResponseWriter) {
	if r := recover(); r != nil {
		defer http.Error(w, fmt.Sprintf("%s", "Internal server error"), http.StatusInternalServerError)
		fmt.Println(errors.Wrap(r, 2).ErrorStack())
		// TODO: log the actual error to a server log file (but do not return to client)
	}
}

// TODO should go into config file
var uploadFileArea = filepath.Join(".", "upload_area")
var downloadFileArea = filepath.Join(".", "download_area")
var symbolSetFileArea string // = filepath.Join(".", "symbol_files")
var dbFileArea string        // = filepath.Join(".", "db_files")
var staticFolder string      // = "."

// TODO config stuff
func initFolders() {
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

	// If the db dir doesn't exist, create it
	if _, err := os.Stat(dbFileArea); err != nil {
		if os.IsNotExist(err) {
			err2 := os.Mkdir(dbFileArea, 0755)
			if err2 != nil {
				fmt.Printf("lexserver.init: failed to create %s : %v", dbFileArea, err2)
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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `<h1>Lexserver</h1>`

	for _, subRouter := range subRouters {
		html = html + `<p><a href="` + subRouter.root + `"><b>` + removeInitialSlash(subRouter.root) + `</b></a>`
		html = html + " | " + subRouter.desc + "</p>\n\n"

	}

	html = html + `<hr/><h2>Create db and start server</h2>
Instructions on how to create a lexicon database and start the server are available from the <a target="blank" href="https://github.com/stts-se/lexdata/wiki/Create-lexicon-database">Lexdata git Wiki</a>.
`
	fmt.Fprint(w, html)
}

// func sqlite3AnalyzeHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Print("lexserver: running the Sqlite3 ANALYZE command...")
// 	_, err := db.Exec("ANALYZE")

// 	if err != nil {
// 		log.Printf("Failed to exec ANALYZE %v", err)
// 		http.Error(w, fmt.Sprintf("/admin/sqlite3_analyze failed : %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	log.Print("... done!\n")
// 	w.Write([]byte("OK"))
// }

// TODO report unused URL parameters

// TODO Gör konstanter som kan användas istället för strängar
var knownParams = map[string]int{
	"lexicons":            1,
	"entryids":            1,
	"words":               1,
	"lemmas":              1,
	"wordlike":            1,
	"wordregexp":          1,
	"wordparts":           1,
	"wordpartslike":       1,
	"wordpartsregexp":     1,
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

func queryFromParams(r *http.Request) (dbapi.DBMQuery, error) {

	lexs := dbapi.RemoveEmptyStrings(
		splitRE.Split(getParam("lexicons", r), -1))
	words := dbapi.RemoveEmptyStrings(
		splitRE.Split(getParam("words", r), -1))
	wordParts := dbapi.RemoveEmptyStrings(
		splitRE.Split(getParam("wordparts", r), -1))
	lemmas := dbapi.RemoveEmptyStrings(
		splitRE.Split(getParam("lemmas", r), -1))
	entryIDStrings := dbapi.RemoveEmptyStrings(
		splitRE.Split(getParam("entryids", r), -1))
	entryIDs := []int64{}
	for _, s := range entryIDStrings {
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return dbapi.DBMQuery{}, fmt.Errorf("couldn't create int64 from input string '%s' : %v", s, err)
		}
		entryIDs = append(entryIDs, id)
	}

	wordLike := strings.TrimSpace(getParam("wordlike", r))
	wordRegexp := strings.TrimSpace(getParam("wordregexp", r))
	wordPartsLike := strings.TrimSpace(getParam("wordpartslike", r))
	wordPartsRegexp := strings.TrimSpace(getParam("wordpartsregexp", r))
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

	lexRefs := []lex.LexRef{}
	for _, l := range lexs {
		ref, err := lex.ParseLexRef(l)
		if err != nil {
			return dbapi.DBMQuery{}, fmt.Errorf("couldn't parse lexicon reference from string %s", l)
		}
		lexRefs = append(lexRefs, ref)
	}

	q := dbapi.Query{
		Words:               words,
		WordParts:           wordParts,
		WordPartsLike:       wordPartsLike,
		WordPartsRegexp:     wordPartsRegexp,
		EntryIDs:            entryIDs,
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
	dq := dbapi.DBMQuery{
		Query:   q,
		LexRefs: lexRefs,
	}
	return dq, nil
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

var dbm = dbapi.NewDBManager()

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

func isStaticPage(url string) bool {
	return url == "/" || strings.Contains(url, "externals") || strings.Contains(url, "built") || url == "/websockreg" || url == "/favicon.ico" || url == "/static/" || url == "/ipa_table.txt"
}

func main() {

	port := ":8787"
	testPort := ":8799"
	tag := "standard"

	var test = flag.Bool("test", false, "run server tests")
	var ssFiles = flag.String("ss_files", filepath.Join(".", "symbol_sets"), "location for symbol set files")
	var dbFiles = flag.String("db_files", filepath.Join(".", "db_files"), "location for db files")
	var static = flag.String("static", filepath.Join(".", "static"), "location for static html files")
	var help = flag.Bool("help", false, "print usage/help and exit")

	usage := `Usage:
     $ go run *.go <PORT>
     $ go run *.go
      - use default port

Flags:
     -test       bool    run server tests and exit (default: false)
     -ss_files   string  location for symbol set files (default: symbol_sets)
     -db_files   string  location for db files (default: db_files)
     -static     string  location for static html files (default: ./)


Default ports:
     ` + port + `  for the standard server
     ` + testPort + `  for the test server
`

	flag.Parse()

	if *help {
		fmt.Println(usage)
		os.Exit(1)
	}

	if *test {
		port = testPort
		tag = "test"
	}

	if len(flag.Args()) > 1 {
		fmt.Println(usage)
		os.Exit(1)
	} else if len(flag.Args()) == 1 {
		port = flag.Args()[0]
	}
	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	symbolSetFileArea = *ssFiles
	dbFileArea = *dbFiles
	staticFolder = *static

	initFolders()

	dbapi.Sqlite3WithRegex()

	log.Println("lexserver: started")

	err := setupDemoDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "COULDN'T INITIALISE DEMO DB : %v\n", err)
		os.Exit(1)
	}

	log.Printf("lexserver: creating %s server on port %s", tag, port)
	s, err := createServer(port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "COULDN'T CREATE SERVER : %v\n", err)
		os.Exit(1)
	}

	if *test {
		err = runInitTests(s, port)
		if err != nil {
			log.Printf("lexserver: %v", err)
			os.Exit(1)
		}
	} else { // start the standard server
		stop := make(chan os.Signal, 1)

		signal.Notify(stop, os.Interrupt)
		go func() {
			if err := s.ListenAndServe(); err != nil {
				log.Fatal(fmt.Errorf("lexserver: couldn't start server on port %s : %v", port, err))
			}
		}()
		log.Printf("lexserver: server up and running using port " + port)

		<-stop

		fmt.Fprintf(os.Stderr, "\n")
		log.Println("lexserver: shutting down...")

		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		defer s.Shutdown(ctx)

		// shut down databases nicely
		dbNames, err := dbm.ListDBNames()
		if err != nil {
			log.Printf("couldn't close databases properly, will exit anyway : %v", err)
		}
		for _, dbName := range dbNames {
			err = dbm.CloseDB(dbName)
			if err != nil {
				log.Printf("couldn't close database %s properly, will exit anyway : %v", string(dbName), err)
			}
			log.Printf("lexserver: closed database %s", string(dbName))
		}
	}
	log.Println("lexserver: BYE!\n")
}

func createServer(port string) (*http.Server, error) {

	var s *http.Server

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}
	rout := mux.NewRouter().StrictSlash(true)

	var err error // återanvänds för alla fel

	log.Print("lexserver: loading dbs from folder ", dbFileArea)
	files, err := ioutil.ReadDir(dbFileArea)
	if err != nil {
		return s, fmt.Errorf("couldn't open db file area: %v", err)
	}

	nDbs := 0
	for _, f := range files {
		dbPath := filepath.Join(dbFileArea, f.Name())
		if !strings.HasSuffix(dbPath, ".db") {
			log.Printf("lexserver: skipping file: '%s'\n", dbPath)
			continue
		}
		nDbs = nDbs + 1
		log.Print("lexserver: connecting to Sqlite3 db ", dbPath)
		// kolla att db-filen existerar
		_, err = os.Stat(dbPath)
		if err != nil {
			return s, fmt.Errorf("lexserver: cannot find db file. %v", err)
		}
		var db *sql.DB
		db, err = sql.Open("sqlite3_with_regexp", dbPath)
		if err != nil {
			return s, fmt.Errorf("Failed to open dbfile %v", err)
		}
		_, err = db.Exec("PRAGMA foreign_keys = ON")
		if err != nil {
			return s, fmt.Errorf("Failed to exec PRAGMA call %v", err)
		}
		_, err = db.Exec("PRAGMA case_sensitive_like=ON")
		if err != nil {
			return s, fmt.Errorf("Failed to exec PRAGMA call %v", err)
		}
		_, err = db.Exec("PRAGMA journal_mode=WAL")
		if err != nil {
			return s, fmt.Errorf("Failed to exec PRAGMA call %v", err)
		}
		//_, err = db.Exec("PRAGMA busy_timeout=500") // doesn't seem to do the trick
		//		if err != nil {
		//return s, fmt.Errorf("Failed to exec PRAGMA call %v", err)
		//}
		db.SetMaxOpenConns(1) // to avoid locking errors (but it makes it slow...?) https://github.com/mattn/go-sqlite3/issues/274

		dbName := filepath.Base(dbPath)
		var extension = filepath.Ext(dbName)
		dbName = dbName[0 : len(dbName)-len(extension)]
		dbRef := lex.DBRef(dbName)
		err = dbm.AddDB(dbRef, db)
		if err != nil {
			return s, fmt.Errorf("Failed to add db: %v", err)
		}

	}

	log.Printf("lexserver: loaded %v db(s)", nDbs)

	// load symbol set mappers
	err = loadSymbolSets(symbolSetFileArea)
	if err != nil {
		return s, fmt.Errorf("Failed to load symbol sets from dir "+symbolSetFileArea+" : %v", err)
	}
	log.Printf("lexserver: loaded symbol sets from dir %s", symbolSetFileArea)

	err = loadConverters(symbolSetFileArea)
	if err != nil {
		return s, fmt.Errorf("Failed to load converters from dir "+symbolSetFileArea+" : %v", err)
	}
	log.Printf("lexserver: loaded converters from dir %s", symbolSetFileArea)

	err = loadValidators(symbolSetFileArea)
	if err != nil {
		return s, fmt.Errorf("Failed to load validators : %v", err)
	}
	log.Printf("lexserver: loaded validators : %v", validatorNames())

	rout.HandleFunc("/", indexHandler)

	lexicon := newSubRouter(rout, "/lexicon", "Lexicon management/admin, including full validation")
	lexicon.addHandler(lexiconList)
	lexicon.addHandler(lexiconLookup) // has its own index page in static/
	lexicon.addHandler(lexiconInfo)
	lexicon.addHandler(lexiconStats)
	lexicon.addHandler(lexiconListCurrentEntryStatuses)
	lexicon.addHandler(lexiconListAllEntryStatuses)
	lexicon.addHandler(lexiconUpdateEntry)
	lexicon.addHandler(lexiconValidationPage)
	lexicon.addHandler(lexiconValidation)
	lexicon.addHandler(lexiconAddEntry)

	validation := newSubRouter(rout, "/validation", "Transcription/entry validation")
	validation.addHandler(validationValidateEntry)
	validation.addHandler(validationValidateEntries)
	validation.addHandler(validationListValidators)
	validation.addHandler(validationStats)
	validation.addHandler(validationHasValidator)

	symbolset := newSubRouter(rout, "/symbolset", "Handle transcription symbol sets")
	symbolset.addHandler(symbolsetList)
	symbolset.addHandler(symbolsetDelete)
	symbolset.addHandler(symbolsetContent)
	symbolset.addHandler(symbolsetReloadOne)
	symbolset.addHandler(symbolsetReloadAll)
	symbolset.addHandler(symbolsetUploadPage)
	symbolset.addHandler(symbolsetUpload)

	mapper := newSubRouter(rout, "/mapper", "Map transcriptions between different symbol sets")
	mapper.addHandler(mapperList)
	mapper.addHandler(mapperMap)
	mapper.addHandler(mapperMaptable)

	converter := newSubRouter(rout, "/converter", "Convert transcriptions between languages")
	converter.addHandler(converterConvert)
	converter.addHandler(converterList)
	converter.addHandler(converterTable)

	admin := newSubRouter(rout, "/admin", "Misc admin tools")
	admin.addHandler(adminLexImportPage)
	admin.addHandler(adminLexImport)
	admin.addHandler(adminListDBs)
	admin.addHandler(adminCreateDB)
	admin.addHandler(adminMoveNewEntries)
	admin.addHandler(adminDeleteLex)
	admin.addHandler(adminSuperDeleteLex)

	// Sqlite3 ANALYZE command in some instances make search quicker,
	// but it takes a while to perform. TODO: Re-add this call?
	//rout.HandleFunc("/admin/sqlite3_analyze", sqlite3AnalyzeHandler)

	rout.Handle("/websockreg", websocket.Handler(webSockRegHandler))

	// typescript experiments
	demo := newSubRouter(rout, "/demo", "Search demo for testing (using typesscript)")
	var searchDemoHandler = urlHandler{
		name:     "search demo page",
		url:      "/search",
		help:     "Typescript search demo.",
		examples: []string{"/search"},
		handler: func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, "../web/lexsearch/index.html")
		},
	}

	// search demo / typescript test
	demo.addHandler(searchDemoHandler)
	rout.PathPrefix("/lexsearch/externals/").Handler(http.StripPrefix("/lexsearch/externals/", http.FileServer(http.Dir("../web/lexsearch/externals"))))
	rout.PathPrefix("/lexsearch/built/").Handler(http.StripPrefix("/lexsearch/built/", http.FileServer(http.Dir("../web/lexsearch/built"))))

	// static
	rout.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticFolder, "favicon.ico"))
	})
	rout.HandleFunc("/ipa_table.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticFolder, "ipa_table.txt"))
	})
	rout.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(staticFolder)))))

	var urls = []string{}
	rout.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		url, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		if !isStaticPage(url) && !isHandeledPage(url) {
			log.Print("Unhandeled url: ", url)
			urls = append(urls, url+" (UNHANDELED)")
		} else {
			urls = append(urls, url)
		}
		return nil
	})

	meta := newSubRouter(rout, "/meta", "Meta API calls (list served URLs, etc)")
	meta.addHandler(metaURLsHandler(urls))
	meta.addHandler(metaExamplesHandler)

	// Pinging connected websocket clients
	go keepClientsAlive()

	log.Print("lexserver: server created but not started for port ", port)

	s = &http.Server{
		Addr:           port,
		Handler:        rout,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s, nil
}
