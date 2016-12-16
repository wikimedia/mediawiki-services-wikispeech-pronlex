package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"encoding/json"

	"github.com/stts-se/pronlex/lexserver/user_db"

	"github.com/gorilla/mux"
)

var dbFile = "user_db.db"

func xtractNameAndPassword(r *http.Request) (string, string, error) {
	name := r.FormValue("name")
	password := r.FormValue("password")

	if name == "" { // Other tests on valid name?
		msg := "missing value for parameter 'name'"
		err := fmt.Errorf(msg)
		return "", "", err
	}
	if password == "" { // Other tests on valid password? See e.g. https://github.com/nbutton23/zxcvbn-go
		msg := "missing value for parameter 'password'"
		err := fmt.Errorf(msg)
		return "", "", err
	}

	return name, password, nil
}

func createUser(w http.ResponseWriter, r *http.Request) {

	name, password, err := xtractNameAndPassword(r)
	if err != nil {
		msg := fmt.Sprintf("missing param value : %v", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	userDB.InsertUser(lexserver.User{Name: name}, password)
}

func listUsers(w http.ResponseWriter, r *http.Request) {

	users, err := userDB.GetUsers()
	if err != nil {
		msg := fmt.Sprintf("user db query failed : %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(users)
	if err != nil {
		msg := fmt.Sprintf("failed to marshal db result : %v", err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(json))
}

//func deleteUser

// TODO mutexify
var userDB lexserver.UserDB

func main() {

	_, err := os.Stat(dbFile)
	if err != nil {
		log.Printf("creating empty user db file: '%s'", dbFile)

		err := lexserver.CreateEmptyUserDB(dbFile)
		if err != nil {

			msg := fmt.Sprintf("No-no-no! Failed to create empty user db file : %v", err)
			log.Println(msg)
			return
		}
	}

	userDB, err = lexserver.InitUserDB(dbFile)

	r := mux.NewRouter()
	r.HandleFunc("/admin/user_db/add_user", createUser)

}
