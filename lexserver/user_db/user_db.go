package lexserver // TODO Restructure lexserver into sub-directories

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	//PasswordHash string `json:"password_hash"` // TODO Should not be part of User?
	Roles string `json:"roles"`
	DBs   string `json:"dbs"`
}

type UserDB struct {
	mutex *sync.Mutex
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

	//rows, err := tx.Query("SELECT id, name, password_hash, roles, dbs FROM user")
	rows, err := tx.Query("SELECT id, name, roles, dbs FROM user")
	if err != nil {
		tx.Rollback()
		return res, fmt.Errorf("user db query failed : %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		rows.Scan(&u.ID, &u.Name /*&u.PasswordHash,*/, &u.Roles, &u.DBs)
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

	//err = tx.QueryRow("SELECT id, name, password_hash, roles, dbs FROM user WHERE name = ?", strings.ToLower(name)).Scan(&res.ID, &res.Name, &res.PasswordHash, &res.Roles, &res.DBs)
	err = tx.QueryRow("SELECT id, name, roles, dbs FROM user WHERE name = ?", strings.ToLower(name)).Scan(&res.ID, &res.Name /*&res.PasswordHash,*/, &res.Roles, &res.DBs)
	if err != nil {
		return res, fmt.Errorf("GetUserByName failed to get user '%s' : %v", name, err)
		tx.Rollback()
	}

	return res, nil
}

// GetPasswordHash returns the password_hash value for userName. If no
// such value is found, the empty string is returned (along with a
// non-nil error value)
func (udb UserDB) GetPasswordHash(userName string) (string, error) {
	name := strings.ToLower(userName)
	var res string
	tx, err := udb.Begin()
	if err != nil {
		msg := fmt.Sprintf("failed starting transaction : %v", err)
		return res, fmt.Errorf(msg)
	}
	defer tx.Commit()

	err = tx.QueryRow("SELECT password_hash FROM user WHERE name = ?", name).Scan(&res)
	if err != nil || err == sql.ErrNoRows || res == "" {
		return res, fmt.Errorf("password hash not found for user '%s'", userName)
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

func (udb UserDB) DeleteUser(userName string) error {
	name := strings.ToLower(userName)
	tx, err := udb.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction : %v", err)
	}
	defer tx.Commit()

	rez, err := tx.Exec("DELETE FROM user WHERE name = ?", name)
	if err != nil {
		return fmt.Errorf("failed to delete user '%s' : ", userName, err)
	}

	ra, _ := rez.RowsAffected()
	if ra < 1 {
		return fmt.Errorf("could not delete user '%v' (nu such user?)", userName)
	}

	//if rez.RowsAffected

	return nil
}

// Update updates the fields of User except for User.ID and User.Name.
// Zero valued fields (empty string) will be treated as acceptable
// values, and updated to the empty string in the DB.

func (udb UserDB) Update(user User) error {
	name := strings.ToLower(user.Name)
	tx, err := udb.Begin()
	if err != nil {
		return fmt.Errorf("failed to sdtart transaction : %v", err)
	}
	defer tx.Commit()

	rez, err := tx.Exec("UPDATE user SET roles = ?, dbs = ? WHERE name = ?", user.Roles, user.DBs, name)
	if err != nil {
		return fmt.Errorf("failed to update user '%s' : %v", name, err)
	}

	ra, _ := rez.RowsAffected()
	if ra < 1 {
		return fmt.Errorf("failed to update user '%s' (does the user exist?)", user.Name)
	}

	return nil
}

//TODO UpdatePassword

func (udb UserDB) Authorized(userName, password string) (bool, error) {
	ok := false
	//res := ""
	name := strings.ToLower(userName)

	res, err := udb.GetPasswordHash(name)
	if err != nil {
		return ok, fmt.Errorf("failed to get user '%s' from user db : %v", name, err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(res), []byte(password)); err != nil {
		return ok, fmt.Errorf("password doesn't match")
	}

	// password matches hash in db
	ok = true

	return ok, nil
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

	//var m sync.Mutex
	return UserDB{&sync.Mutex{}, db}, nil
}
