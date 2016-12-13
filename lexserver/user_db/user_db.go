package lexserver // TODO Restructure lexserver into sub-directories

import (
	"database/sql"
	"fmt"
	"os"
)

type UserDB struct {
	userDB *sql.DB
}

type User struct {
	ID    int64
	Name  string
	Roles string
	DBs   string
}

func (udb UserDB) getUserByName(name string) User {

	return User{}
}

func (udb UserDB) insertUser(u User) error {

	return nil
}

var userDBSchema = `CREATE TABLE IF NOT EXISTS user (
  id INTEGER PRIMARY KEY,
  name TEXT NOT NULL UNIQUE,
  password_hash varchar[128] NOT NULL,
  roles TEXT,
  dbs TEXT);`

func createEmptyUserDB(fName string) error {
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

func initUserDB(fName string) (*sql.DB, error) {
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
