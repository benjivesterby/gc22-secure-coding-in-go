package main

import (
	"database/sql"
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
	createUser := fmt.Sprintf(
		`INSERT INTO users VALUES(NULL,'%s','%s','%s');`,
		name, email, password,
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

func GetUser(db *sql.DB, email, password string) (*User, error) {
	log.Printf("Getting user [%s] from database", email)

	q := fmt.Sprintf(`
		SELECT * 
		FROM users 
		WHERE Email='%s' AND Password = '%s'
		LIMIT 1;`,
		email,
		password,
	)

	fmt.Println(q)

	row := db.QueryRow(q)
	if row == nil {
		return nil, fmt.Errorf("User not found")
	}

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, err
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

func GetFriends(userId string) ([]*User, error) {
	q := fmt.Sprintf(`
		SELECT users.*
		FROM users
		JOIN friends ON users.ID = friends.FriendId
		WHERE friends.UserId = '%s';
		`,
		userId,
	)

	return Query[*User](q)
}

func Query[T any](query string) ([]T, error) {
	dbx, err := sqlx.Open("sqlite3", db)
	if err != nil {
		return nil, err
	}

	var results []T
	err = dbx.Select(&results, query)
	if err != nil {
		return nil, err
	}

	return results, nil
}
