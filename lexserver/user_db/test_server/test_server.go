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

	// These values may be empty
	roles := r.FormValue("roles")
	dbs := r.FormValue("dbs")

	err = userDB.InsertUser(lexserver.User{Name: name, Roles: roles, DBs: dbs}, password)
	if err != nil {
		msg := fmt.Sprintf("failed to insert user '%s' : %v", name, err)
		log.Printf("%s\n", msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// TODO How do you return an 'OK'?
	fmt.Fprintf(w, "Added user '%s'", name)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	err := userDB.DeleteUser(name)
	if err != nil {
		msg := fmt.Sprintf("failed to delete user '%s' : %v", name, err)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Deleted user '%s'", name)
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
	r.HandleFunc("/admin/user_db/list_users", listUsers)
	r.HandleFunc("/admin/user_db/delete_user", deleteUser)

	port := ":8788"
	log.Println("Starting user db test_server on port %s", port)
	err = http.ListenAndServe(port, r)
	if err != nil {
		log.Fatalf("things are not working : %v", err)
	}

}
