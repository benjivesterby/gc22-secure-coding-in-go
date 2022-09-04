package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"go.benjiv.com/gc22/ui"
)

func main() {
	if !strings.Contains(os.Args[0], "go") {
		fmt.Print(grumpy)
		fmt.Println(`DO NOT BUILD OR INSTALL THIS!`)
		os.Exit(1)
	}

	api := &API{}

	var err error
	api.db, err = InitDB()
	if err != nil {
		log.Fatal(err)
	}

	app, err := api.App()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.Handle("/", http.FileServer(app))
	router.Handle("/imgs/", http.HandlerFunc(api.Image))
	router.Handle("/upload", http.HandlerFunc(api.Upload))
	router.Handle("/user", http.HandlerFunc(api.User))
	router.Handle("/friend", http.HandlerFunc(api.Friend))
	router.Handle("/friends", http.HandlerFunc(api.Friends))
	router.Handle("/users", http.HandlerFunc(api.Users))
	router.Handle("/images", http.HandlerFunc(api.Pictures))
	router.Handle("/search", http.HandlerFunc(api.Search))

	http.ListenAndServe(":8081", router)
}

const prefix = "html"

type API struct {
	db       *sql.DB
	sessions map[int]string
}

func (api *API) App() (http.FileSystem, error) {
	fsys, err := fs.Sub(ui.FS, "build")
	if err != nil {
		return nil, err
	}

	return http.FS(fsys), nil
}

// http://localhost:8081/imgs/1/rick.jpg
func (api *API) Image(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		rw.Header().Set("Content-Type", "image/jpeg")
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		path := strings.TrimPrefix(req.URL.Path, "/imgs/")
		file := filepath.Join(wd, "images", path)

		http.ServeFile(rw, req, file)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (api *API) Search(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		api.GetSearch(rw, req)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type Results struct {
	Query string
	Users []User
}

func (api *API) GetSearch(rw http.ResponseWriter, req *http.Request) {
	query := req.URL.Query().Get("query")

	log.Printf("Getting search results for [%s]", query)

	// TODO: Implement search
	results := Results{
		Query: query,
		Users: []User{},
	}

	data, err := json.Marshal(results)
	if err != nil {
		log.Fatal(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Write(data)
}
