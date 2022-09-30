package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

const prefix = "html"

type API struct {
	db       *sql.DB
	sessions map[int]string
}

func main() {
	if !strings.Contains(os.Args[0], "go") {
		fmt.Print(grumpy)
		fmt.Println(`DO NOT BUILD OR INSTALL THIS!`)
		os.Exit(1)
	}

	api := &API{}

	api.db, err = InitDB()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	// Auth
	router.Handle("/login", http.HandlerFunc(api.Login))

	// Users
	router.Handle("/user", http.HandlerFunc(api.User))
	router.Handle("/users", http.HandlerFunc(api.Users))
	router.Handle("/friend", http.HandlerFunc(api.Friend))
	router.Handle("/friends", http.HandlerFunc(api.Friends))

	// Images
	router.Handle("/images", http.HandlerFunc(api.Pictures))
	router.Handle("/imgs/", http.HandlerFunc(api.Image))
	router.Handle("/upload", http.HandlerFunc(api.Upload))

	log.Printf("Listening on port 8081")
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", router))
}
