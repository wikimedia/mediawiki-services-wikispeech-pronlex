package lexserver // TODO Restructure lexserver into sub-directories

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type UserDB struct {
	*sql.DB
}

type User struct {
	ID           int64
	Name         string
	PasswordHash string
	Roles        string
	DBs          string
}

func (udb UserDB) GetUserByName(name string) (User, error) {
	res := User{}
	tx, err := udb.Begin()
	if err != nil {
		return res, fmt.Errorf("GetUserByName failed to start transaction : %v", err)
	}
	defer tx.Commit()

	//err = tx.QueryRow("SELECT id, name, password_hash, roles, dbs FROM user WHERE name = ?", strings.ToLower(name)).Scan(&qid, &qname, &qpasswordHash, &qroles, &qdbs)
	err = tx.QueryRow("SELECT id, name, password_hash, roles, dbs FROM user WHERE name = ?", strings.ToLower(name)).Scan(&res.ID, &res.Name, &res.PasswordHash, &res.Roles, &res.DBs)
	if err != nil {
		return res, fmt.Errorf("GetUserByName failed to get user '%s' : %v", name, err)
	}

	return res, nil
}

func (udb UserDB) InsertUser(u User, password string) error {
	// TODO add transactions

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	name := strings.ToLower(u.Name)

	fmt.Printf("insertUser: %s %s %s", name, password, passwordHash)

	if err != nil {
		return fmt.Errorf("failed to generate hash: %v", err)
	}
	_, err = udb.Exec("INSERT INTO user (name, password_hash, roles, dbs) VALUES (?, ?, ?, ?)", name, string(passwordHash), u.Roles, u.DBs)

	if err != nil {
		return fmt.Errorf("failed to insert user into db: %v", err)
	}

	return nil
}

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

func InitUserDB(fName string) (*sql.DB, error) {
	if _, err := os.Stat(fName); os.IsNotExist(err) {
		return nil, fmt.Errorf("db file doesn't exist: '%s'", fName)
	}

	//var err error
	userDB, err := sql.Open("sqlite3", fName)
	if err != nil {
		return nil, fmt.Errorf("failed to open db file: '%s': %v", fName, err)
	}

	return userDB, nil
}

// func InsertUser(db *sql.DB, u User, password string) error {
// 	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

// 	fmt.Printf("insertUser: %s %s %s", u.Name, password, passwordHash)

// 	if err != nil {
// 		return fmt.Errorf("failed to generate hash: %v", err)
// 	}
// 	_, err = db.Exec("INSERT INTO user (name, password_hash, roles, dbs) VALUES (?, ?, ?, ?)", u.Name, string(passwordHash), u.Roles, u.DBs)

// 	if err != nil {
// 		return fmt.Errorf("failed to insert user into db: %v", err)
// 	}

// 	return nil
// }
