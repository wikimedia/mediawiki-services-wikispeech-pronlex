package lexserver // TODO Restructure lexserver into sub-directories

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	PasswordHash string `json:"password_hash"` // TODO Should not be part of User?
	Roles        string `json:"roles"`
	DBs          string `json:"dbs"`
}

type UserDB struct {
	*sql.DB
}

// TODO test me
func (udb UserDB) GetUsers() ([]User, error) {
	var res []User

	tx, err := udb.Begin()
	if err != nil {
		return res, fmt.Errorf("GetUsers failed to create transaction : %v", err)
	}
	defer tx.Commit()

	rows, err := tx.Query("SELECT id, name, password_hash, roles, dbs FROM user")
	if err != nil {
		tx.Rollback()
		return res, fmt.Errorf("user db query failed : %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		rows.Scan(&u.ID, &u.Name, &u.PasswordHash, &u.Roles, &u.DBs)
		res = append(res, u)
	}

	return res, nil
}

func (udb UserDB) GetUserByName(name string) (User, error) {
	res := User{}
	tx, err := udb.Begin()
	if err != nil {
		return res, fmt.Errorf("GetUserByName failed to start transaction : %v", err)
	}
	defer tx.Commit()

	err = tx.QueryRow("SELECT id, name, password_hash, roles, dbs FROM user WHERE name = ?", strings.ToLower(name)).Scan(&res.ID, &res.Name, &res.PasswordHash, &res.Roles, &res.DBs)
	if err != nil {
		return res, fmt.Errorf("GetUserByName failed to get user '%s' : %v", name, err)
		tx.Rollback()
	}

	return res, nil
}

func (udb UserDB) InsertUser(u User, password string) error {
	tx, err := udb.Begin()
	if err != nil {
		return fmt.Errorf("InsertUser failed to start transaction : %v", err)
	}
	defer tx.Commit()

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	name := strings.ToLower(u.Name)

	//fmt.Printf("insertUser: %s %s %s", name, password, passwordHash)

	if err != nil {
		return fmt.Errorf("failed to generate hash: %v", err)
	}
	_, err = tx.Exec("INSERT INTO user (name, password_hash, roles, dbs) VALUES (?, ?, ?, ?)", name, string(passwordHash), u.Roles, u.DBs)

	if err != nil {
		return fmt.Errorf("failed to insert user into db: %v", err)
		tx.Rollback()
	}

	return nil
}

func (udb UserDB) Authorized(name, password string) (bool, User, error) {
	ok := false
	res := User{}

	res, err := udb.GetUserByName(name)
	if err != nil {
		return ok, res, fmt.Errorf("failed to get user '%s' from user db : %v", name, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(res.PasswordHash), []byte(password)); err != nil {
		return ok, res, fmt.Errorf("password doesn't match")
	}

	// password matches hash in db
	ok = true

	return ok, res, nil
}

//=================================================================================

var userDBSchema = `CREATE TABLE IF NOT EXISTS user (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  password_hash varchar[128] NOT NULL,
  roles TEXT,
  dbs TEXT);`

func CreateEmptyUserDB(fName string) error {
	if _, err := os.Stat(fName); !os.IsNotExist(err) {
		return fmt.Errorf("Cannot create file that already exists: '%s'", fName)
	}

	db, err := sql.Open("sqlite3", fName)
	if err != nil {
		return fmt.Errorf("failed to open '%s': %v", fName, err)
	}

	_, err = db.Exec(userDBSchema)
	if err != nil {
		return fmt.Errorf("failed to create user database tabl: %v", err)
	}

	return nil
}

func InitUserDB(fName string) (UserDB, error) {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return UserDB{}, fmt.Errorf("db file doesn't exist: '%s'", fName)
	}

	//var err error
	db, err := sql.Open("sqlite3", fName)
	if err != nil {
		return UserDB{}, fmt.Errorf("failed to open db file: '%s': %v", fName, err)
	}

	return UserDB{db}, nil
}
