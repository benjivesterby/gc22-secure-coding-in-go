package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	api := &API{}

	var err error
	api.db, err = InitDB()
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(api.Route))
	router.Handle("/upload", http.HandlerFunc(api.Upload))
	router.Handle("/user", http.HandlerFunc(api.User))

	http.ListenAndServe(":8081", router)
}

const prefix = "html"

type API struct {
	db *sql.DB
}

func (api *API) Route(rw http.ResponseWriter, req *http.Request) {
	file := fmt.Sprintf("%s/index.html", prefix)
	if req.URL.Path != "/" {
		file = fmt.Sprintf("%s/%s", prefix, string(req.URL.Path[1:]))
	}

	fmt.Printf("Routing to %s\n", file)

	body, err := os.ReadFile(file)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Write(body)
}

func (api *API) Upload(w http.ResponseWriter, req *http.Request) {
	media, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	if strings.HasPrefix(media, "multipart/") {
		mr := multipart.NewReader(req.Body, params["boundary"])
		start := time.Now()

		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatal(err)
			}

			body, err := io.ReadAll(p)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(
				fmt.Sprintf("uploads/%s", p.FileName()),
				body,
				0666,
			)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("File: %s Took: %s\n", p.FileName(), time.Since(start))
		}
	}
}
