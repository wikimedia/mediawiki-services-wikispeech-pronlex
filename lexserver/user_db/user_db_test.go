package lexserver

import (
	"os"
	"testing"
)

var fs = "Expected '%v', got '%v'"

func Test_UserDB(t *testing.T) {
	dbFile := "user_db_test.db"
	err := CreateEmptyUserDB(dbFile)
	defer os.Remove(dbFile)
	if err != nil {
		t.Errorf("failed to create empty db file : %v", err)
	}

	sql, err := InitUserDB(dbFile)
	if err != nil {
		t.Errorf("failed to initialise empty db file : %v", err)
	}

	udb := UserDB{sql}

	//_ = udb

	u := User{Name: "KalleA", Roles: "admin:cleaner", DBs: "ankeborg"}
	//_ = u

	err = udb.InsertUser(u, "sekret")
	if err != nil {
		t.Errorf("Fail: %v", err)
	}

	u1, err := udb.GetUserByName("KalleA")
	if err != nil {
		t.Errorf("oh no : %v", err)
	}
	if w, g := "kallea", u1.Name; w != g {
		t.Errorf(fs, w, g)
	}
	if u1.ID == 0 {
		t.Errorf("expected id > 0, got %d", u1.ID)
	}
	if w, g := "admin:cleaner", u1.Roles; w != g {
		t.Errorf(fs, w, g)
	}
	if w, g := "ankeborg", u1.DBs; w != g {
		t.Errorf(fs, w, g)
	}
	if u1.PasswordHash == "" {
		t.Errorf("Expected non zero value hash: %#v", u1)
	}
}
