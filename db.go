package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log"
)

const db = "supersecret.db"

const users = `
  CREATE TABLE IF NOT EXISTS users (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL
  );`

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", db)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(users); err != nil {
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
