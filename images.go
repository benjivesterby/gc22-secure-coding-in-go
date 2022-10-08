package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

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

		fmt.Println(file)

		http.ServeFile(rw, req, file)
	default:
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

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
	// Path Traversal
	// curl -H "userId: ../../" localhost:8081/images
	userId := req.Header.Get("userId")
	log.Printf("Getting pictures for user [%s]", userId)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	wd = filepath.Join(wd, "images", userId)

	cmd := exec.Command("ls", wd)

	out, err := cmd.Output()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(strings.Split(string(out), "\n"))
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
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
