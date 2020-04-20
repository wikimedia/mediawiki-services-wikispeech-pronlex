package main

import (
	//"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/stts-se/pronlex/dbapi"
	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/validation"
)

func deleteUploadedFile(serverPath string) {
	// when done, delete from server!
	err := os.Remove(serverPath)
	if err != nil {
		msg := fmt.Sprintf("couldn't delete temp file from server : %v", err)
		log.Println(msg)
	} else {
		msg := "the uploaded temp file has been deleted from server"
		log.Println(msg)
	}
}

var adminListIDs = urlHandler{
	name:     "list_ids",
	url:      "/list_ids/{lexicon_name}",
	help:     "List all IDs for the entries in one lexicon.",
	examples: []string{"/list_ids/lexserver_testdb:sv"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		ids, err := dbm.ListIDs(lexRef)

		if err != nil {
			log.Printf("lexserver: Failed to get ids: %v", err)
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		jsids := IDs{ids}
		res, err := json.Marshal(jsids)
		if err != nil {
			msg := fmt.Sprintf("lexserver: Failed to marshal ids : %v", err)
			log.Print(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, string(res))

	},
}

var adminLexImportPage = urlHandler{
	name:     "lex_import (page)",
	url:      "/lex_import_page",
	help:     "Import lexicon file (GUI).",
	examples: []string{"/lex_import_page"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(staticFolder, "admin/lex_import_page.html"))
	},
}

var adminLexImport = urlHandler{
	name: "lex_import (api)",
	url:  "/lex_import",
	//help:     "Import lexicon file (API). Requires POST request. Mainly for server internal use.<p/>Available params: lexicon_name, symbolset_name, validate, file",
	help:     "Import lexicon file (API). Requires POST request. Mainly for server internal use.<p/>Available params: lexicon_name, symbolset_name, file",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {

		defer protect(w) // use this call in handlers to catch 'panic' and stack traces and returning a general error to the calling client

		if r.Method != "POST" {
			http.Error(w, fmt.Sprintf("lexiconfileupload only accepts POST request, got %s", r.Method), http.StatusBadRequest)
			return
		}

		clientUUID := getParam("client_uuid", r)

		if strings.TrimSpace(clientUUID) == "" {
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
		locale := r.PostFormValue("locale")
		if strings.TrimSpace(locale) == "" {
			msg := "input param <locale> must not be empty"
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}

		// vString := r.PostFormValue("validate")
		// if strings.TrimSpace(vString) == "" {
		// 	msg := "input param <validate> must not be empty (should be 'true' or 'false')"
		// 	log.Println(msg)
		// 	http.Error(w, msg, http.StatusInternalServerError)
		// 	return
		// }
		// validate, err := strconv.ParseBool(vString)
		// if err != nil {
		// 	log.Println(err)
		// 	http.Error(w, fmt.Sprintf("adminLexImport failed parsing boolean argument %s : %v", vString, err), http.StatusInternalServerError)
		// 	return
		// }

		// (partially) lifted from https://github.com/astaxie/build-web-application-with-golang/blob/master/de/04.5.md
		err = r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("adminLexImport failed parse multipart form : %v", err), http.StatusInternalServerError)
			return
		}

		file, handler, err := r.FormFile("file")
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("adminLexImport failed reading file : %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		serverPath := filepath.Join(uploadFileArea, handler.Filename)

		f, err := os.OpenFile(serverPath, os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("adminLexImport failed opening local output file : %v", err), http.StatusInternalServerError)
			return
		}
		/* #nosec G307 */
		defer f.Close()
		_, err = io.Copy(f, file)
		if err != nil {
			msg := fmt.Sprintf("adminLexImport failed copying local output file : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		exists, err := dbm.LexiconExists(lexRef)
		if err != nil {
			msg := fmt.Sprintf("Couldn't lookup lexicon reference: %s", lexRef.String())
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			deleteUploadedFile(serverPath)
			return
		}
		if exists {
			msg := fmt.Sprintf("Nothing will be added. Lexicon already exists: %s", lexRef.String())
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			deleteUploadedFile(serverPath)
			return
		}
		// _, err = dbapi.GetLexicon(db, lexName)
		// if err == nil {
		// 	msg := fmt.Sprintf("Nothing will be added. Lexicon already exists in database: %s", lexName)
		// 	log.Println(msg)
		// 	http.Error(w, msg, http.StatusInternalServerError)
		// 	deleteUploadedFile(serverPath)
		// 	return
		// }

		err = dbm.DefineLexicon(lexRef, symbolSetName, locale)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			deleteUploadedFile(serverPath)
			return
		}
		log.Println("Created lexicon: ", lexRef.String())

		var validator *validation.Validator = &validation.Validator{}
		// if validate {
		// 	vMut.Lock()
		// 	validator, err = vMut.service.ValidatorForName(symbolSetName)
		// 	vMut.Unlock()
		// 	if err != nil {
		// 		msg := fmt.Sprintf("adminLexImport failed to get validator for symbol set %v : %v", symbolSetName, err)
		// 		log.Println(msg)
		// 		http.Error(w, msg, http.StatusBadRequest)
		// 		return
		// 	}
		// }

		err = dbm.ImportLexiconFile(lexRef, logger, serverPath, validator)

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

		//f.Close()
		deleteUploadedFile(serverPath)

		entryCount, err := dbm.EntryCount(lexRef)
		if err != nil {
			msg := fmt.Sprintf("lexicon imported, but couldn't retrieve lexicon info from server : %v", err)
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		info := LexWithEntryCount{
			Name:          lexRef.String(),
			SymbolSetName: symbolSetName,
			EntryCount:    entryCount,
		}
		fmt.Fprintf(w, "imported %v entries into lexicon '%v'", info.EntryCount, info.Name)
	},
}

var adminDefineLex = urlHandler{
	name:     "define_lex",
	url:      "/define_lex/{lexicon_name}/{locale}/{symbolset_name}",
	help:     "Define (create) a new (empty) lexicon inside a database.",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {

		defer protect(w) // use this call in handlers to catch 'panic' and stack traces and returning a general error to the calling client

		symbolsetName := delQuote(getParam("symbolset_name", r))
		if symbolsetName == "" {
			http.Error(w, "no value for parameter 'symbolset_name'", http.StatusBadRequest)
			return
		}

		locale := delQuote(getParam("locale", r))
		if locale == "" {
			http.Error(w, "no value for parameter 'locale'", http.StatusBadRequest)
			return
		}

		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}
		exists, err := dbm.LexiconExists(lexRef)
		if err != nil {
			msg := fmt.Sprintf("Couldn't lookup lexicon reference: %s", lexRef.String())
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
		if exists {
			msg := fmt.Sprintf("Nothing will be added. Lexicon already exists: %s", lexRef.String())
			log.Println(msg)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		err = dbm.DefineLexicon(lexRef, symbolsetName, locale)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("%v", err), http.StatusInternalServerError)
			return
		}
		log.Println("Created lexicon: ", lexRef.String())
		fmt.Fprint(w, "Created lexicon "+lexRef.String())
	},
}

var adminDeleteLex = urlHandler{
	name:     "deletelexicon",
	url:      "/deletelexicon/{lexicon_name}",
	help:     "Delete a lexicon reference from the database without removing associated entries.",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {
		lexRef, err := getLexRefParam(r)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
			return
		}
		err = dbm.DeleteLexicon(lexRef)
		if err != nil {
			log.Printf("adminDeleteLex got error : %v\n", err)
			http.Error(w, fmt.Sprintf("failed deleting lexicon : %v", err), http.StatusExpectationFailed)
			return
		}
	},
}

// var adminSuperDeleteLex = urlHandler{
// 	name:     "superdeletelexicon",
// 	url:      "/superdeletelexicon/{lexicon_name}",
// 	help:     "Delete a complete lexicon from the database, including associated entries. This make take some time.",
// 	examples: []string{},
// 	handler: func(w http.ResponseWriter, r *http.Request) {
// 		lexRef, err := getLexRefParam(r)
// 		if err != nil {
// 			log.Println(err)
// 			http.Error(w, fmt.Sprintf("couldn't parse lexicon ref %v : %v", lexRef, err), http.StatusInternalServerError)
// 			return
// 		}

// 		uuid := getParam("client_uuid", r)
// 		log.Println("adminSuperDeleteLex was called")
// 		messageToClientWebSock(uuid, fmt.Sprintf("Super delete was called. This may take quite a while. Lexicon %s", lexRef.String()))
// 		err = dbm.SuperDeleteLexicon(lexRef)
// 		if err != nil {

// 			http.Error(w, fmt.Sprintf("failed super deleting lexicon : %v", err), http.StatusExpectationFailed)
// 			return
// 		}

// 		messageToClientWebSock(uuid, fmt.Sprintf("Done deleting lexicon %s", lexRef))
// 	},
// }

var adminListDBs = urlHandler{
	name:     "list_dbs",
	url:      "/list_dbs",
	help:     "Lists available lexicon databases.",
	examples: []string{"/list_dbs"},
	handler: func(w http.ResponseWriter, r *http.Request) {
		dbs, err := dbm.ListDBNames()
		if err != nil {
			http.Error(w, fmt.Sprintf("list dbs failed : %v", err), http.StatusInternalServerError)
			return
		}
		jsn, err := marshal(dbs, r)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed marshalling : %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, string(jsn))
	},
}

var adminCreateDB = urlHandler{
	name:     "create_db",
	url:      "/create_db/{db_name}",
	help:     "Create a new (empty) lexicon database.",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {
		dbName := delQuote(getParam("db_name", r))
		if dbName == "" {
			http.Error(w, "no value for parameter 'db_name'", http.StatusBadRequest)
			return
		}

		//func (dbm *DBManager) DefineDB(dbRef lex.DBRef, dbPath string) error {
		//dbPath := filepath.Join(*dbClusterLocation, dbName+".db")
		err := dbm.DefineDB(lex.DBRef(dbName), *dbClusterLocation)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't define db : %v", err), http.StatusInternalServerError)
			return
		}

		// if _, err := os.Stat(dbFile); !os.IsNotExist(err) {
		// 	http.Error(w, "Cannot create a db that already exists: "+dbName, http.StatusBadRequest)
		// 	return
		// }

		// db, err := sql.Open("sqlite3", dbFile)
		// if err != nil {
		// 	//db.Close()
		// 	http.Error(w, fmt.Sprintf("sql error : %v", err), http.StatusInternalServerError)
		// 	return
		// }
		// _, err = db.Exec("PRAGMA foreign_keys = ON")
		// if err != nil {
		// 	msg := fmt.Sprintf("sql error : %v", err)
		// 	err2 := db.Close()
		// 	if err2 != nil {
		// 		msg = fmt.Sprintf("%s : failed to close db : %v", msg, err2)
		// 	}
		// 	http.Error(w, msg, http.StatusInternalServerError)
		// 	return
		// }

		// _, err = db.Exec(dbapi.Schema)
		// if err != nil {
		// 	msg := fmt.Sprintf("sql error : %v", err)
		// 	err2 := db.Close()
		// 	if err2 != nil {
		// 		msg = fmt.Sprintf("%s : failed to close db : %v", msg, err2)
		// 	}
		// 	http.Error(w, msg, http.StatusInternalServerError)
		// 	return
		// }
		// err = dbm.AddDB(lex.DBRef(dbName), db)
		// if err != nil {
		// 	msg := fmt.Sprintf("failed to add db : %v", err)
		// 	http.Error(w, msg, http.StatusInternalServerError)
		// 	return
		// }

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		fmt.Fprint(w, "Created database "+dbName)
	},
}

var adminMoveNewEntries = urlHandler{
	name:     "move_new_entries",
	url:      "/move_new_entries/{db_name}/{from_lexicon_name}/{to_lexicon_name}/{new_source}/{new_status}",
	help:     "Move entries from one lexicon to another. N.B! Only entries that do not already exist in the right hand will be moved.",
	examples: []string{},
	handler: func(w http.ResponseWriter, r *http.Request) {
		dbName := delQuote(getParam("db_name", r))
		if dbName == "" {
			http.Error(w, "no value for parameter 'db_name'", http.StatusBadRequest)
			return
		}
		fromLexName := delQuote(getParam("from_lexicon", r))
		if fromLexName == "" {
			http.Error(w, "no value for parameter 'from_lexicon'", http.StatusBadRequest)
			return
		}
		toLexName := delQuote(getParam("to_lexicon", r))
		if toLexName == "" {
			http.Error(w, "no value for parameter 'to_lexicon'", http.StatusBadRequest)
			return
		}

		sourceName := delQuote(getParam("new_source", r))
		if sourceName == "" {
			http.Error(w, "no value for parameter 'source'", http.StatusBadRequest)
			return
		}
		statusName := delQuote(getParam("new_status", r))
		if statusName == "" {
			http.Error(w, "no value for parameter 'status'", http.StatusBadRequest)
			return
		}

		moveRes, err := dbm.MoveNewEntries(lex.DBRef(dbName), lex.LexName(fromLexName), lex.LexName(toLexName), sourceName, statusName)
		if err != nil {
			http.Error(w, fmt.Sprintf("failure when trying to move entries from '%s' to '%s' : %v", fromLexName, toLexName, err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "number of entries moved from '%s' to '%s': %d", fromLexName, toLexName, moveRes.N)
	},
}
