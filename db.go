package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const db = "supersecret.db"

const users = `
  CREATE TABLE IF NOT EXISTS users (
  ID INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  Name VARCHAR(255) NOT NULL,
  Email VARCHAR(255) NOT NULL,
  Password VARCHAR(255) NOT NULL
  );`

const friends = `
  CREATE TABLE IF NOT EXISTS friends (
  UserId INTEGER NOT NULL,
  FriendId INTEGER NOT NULL
  );`

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", db)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(users); err != nil {
		return nil, err
	}

	if _, err := db.Exec(friends); err != nil {
		return nil, err
	}

	return db, nil
}

func NewUser(db *sql.DB, name, email, password string) (int, error) {
	hash := sha256.Sum256([]byte(password))
	encoded := base64.StdEncoding.EncodeToString(hash[:])

	createUser := fmt.Sprintf(
		`INSERT INTO users VALUES(NULL,'%s','%s','%s');`,
		name, email, encoded,
	)

	log.Println(createUser)

	res, err := db.Exec(createUser)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func AddFriend(db *sql.DB, userId, friendId string) error {
	log.Printf(
		"Adding friend [%s] to database for user [%s]",
		friendId,
		userId,
	)

	_, err := db.Exec(`INSERT INTO friends VALUES(?,?);`, userId, friendId)
	return err
}

// GetFriends uses StructScan which exposes too much user information
// Unsafe third party library usage
// VULN: SQL Injection w/ DoS
// localhost:8081/friends?userId=1%27%20union%20select%20%2A%20from%20users%3B--
func GetFriends(userId string) ([]*User, error) {
	dbx, err := sqlx.Open(
		"sqlite3",
		db,
	)
	if err != nil {
		return nil, err
	}

	log.Printf("Getting friends for user [%s]", userId)

	q := fmt.Sprintf(`
		SELECT users.*
		FROM users
		JOIN friends ON users.ID = friends.FriendId
		WHERE friends.UserId = '%s';
		`,
		userId,
	)

	users := []*User{}
	err = dbx.Select(&users, q)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return users, nil
}
