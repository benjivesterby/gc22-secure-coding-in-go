package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(Handle))

	http.ListenAndServe(":8080", router)
}

const prefix = "html"

func Handle(rw http.ResponseWriter, req *http.Request) {
	file := fmt.Sprintf("%s/index.html", prefix)
	if req.URL.Path != "/" {
		file = fmt.Sprintf("%s/%s", prefix, string(req.URL.Path[1:]))
	}

	body, err := os.ReadFile(file)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.Write(body)
}
