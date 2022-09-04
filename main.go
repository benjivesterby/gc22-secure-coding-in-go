package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(api.Route))
	router.Handle("/upload", http.HandlerFunc(api.Upload))
	router.Handle("/user", http.HandlerFunc(api.User))
	router.Handle("/friend", http.HandlerFunc(api.Friend))
	router.Handle("/friends", http.HandlerFunc(api.Friends))
	router.Handle("/users", http.HandlerFunc(api.Users))
	router.Handle("/images", http.HandlerFunc(api.Pictures))

	out, err := Command("ls", ".")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(out))

	http.ListenAndServe(":8081", router)
}

const prefix = "html"

type API struct {
	db       *sql.DB
	sessions map[int]string
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

func Command(cmd string, args ...string) ([]byte, error) {
	log.Printf("Command: %s Args: %s\n", cmd, args)

	args = append([]string{"-c", cmd}, args...)

	return exec.Command("/bin/bash", args...).Output()
}
