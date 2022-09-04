package main

import (
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (api *API) Pictures(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		api.GetPictures(w, req)
	case "POST":
		api.Upload(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (api *API) GetPictures(w http.ResponseWriter, req *http.Request) {

}

// curl -H "userId: 1" -X POST -F file=@rick.jpg localhost:8081/images
func (api *API) Upload(w http.ResponseWriter, req *http.Request) {
	media, params, err := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if err != nil {
		log.Fatal(err)
	}

	// Path Traversal
	// curl -H "userId: ../../" -X POST -F file=@rick.jpg localhost:8081/images
	userId := req.Header.Get("userId")
	log.Printf("Uploading picture for user [%s]", userId)

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

			path := filepath.Join("images", userId, p.FileName())

			err = os.WriteFile(
				path,
				body,
				0666,
			)
			if err != nil {
				log.Fatal(err)
			}

			log.Printf("File: %s Took: %s\n", p.FileName(), time.Since(start))
		}
	}
}
