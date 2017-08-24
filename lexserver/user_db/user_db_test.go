package lexserver

import (
	"os"
	"testing"
	//"fmt"
)

var fs = "Expected '%v', got '%v'"

func Test_UserDB(t *testing.T) {
	dbFile := "user_db_test.db"
	err := CreateEmptyUserDB(dbFile)
	defer os.Remove(dbFile)
	if err != nil {
		t.Errorf("failed to create empty db file : %v", err)
	}

	udb, err := InitUserDB(dbFile)
	if err != nil {
		t.Errorf("failed to initialise empty db file : %v", err)
	}

	//_ = udb

	u := User{Name: "KalleA", Roles: "admin:cleaner", DBs: "ankeborg"}
	//_ = u

	s1, err0 := udb.GetPasswordHash("KalleA")
	if w, g := "", s1; w != g {
		t.Errorf(fs, w, g)
	}
	if err0 == nil {
		t.Error("expected error, got nil")
	}

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

	s2, err2 := udb.GetPasswordHash("KalleA")
	if s2 == "" {
		t.Errorf("expected password hash, got empty string")
	}
	if err2 != nil {
		t.Errorf("expected nil, got %v", err)
	}

	u1.Roles = "neeeewrole"
	u1.DBs = "neeewdbz"

	erru := udb.Update(u1)
	if erru != nil {
		t.Errorf("Expected nil, got %v", erru)
	}

	u2, erru2 := udb.GetUserByName("KalleA")
	if erru2 != nil {
		t.Errorf("expected nil, got %v", erru2)
	}
	if w, g := u1.Roles, u2.Roles; w != g {
		t.Errorf(fs, w, g)
	}
	// if u1.PasswordHash == "" {
	// 	t.Errorf("Expected non zero value hash: %#v", u1)
	// }

	// ==================================

	ok, err := udb.Authorized(u.Name, "sekret")
	if err != nil {
		t.Errorf("Expected nil, got %v", err)
	}

	//fmt.Printf("User: %#v\n", user)

	if w, g := true, ok; w != g {
		t.Errorf(fs, w, g)
	}

	ok, err = udb.Authorized(u.Name, "wrongily")
	if err == nil {
		t.Errorf("expected error here")
	}
	if w, g := false, ok; w != g {
		t.Errorf(fs, w, g)
	}

}
