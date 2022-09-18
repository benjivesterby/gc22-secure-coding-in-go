package main

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (api *API) Login(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		api.LoginUser(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

/*
// Example body to detect injection

	{
	    "email": "' or 1=1 --",
	    "password": "doesn't matter"
	}

// Example body to exploit injection

	{
	    "email": "admin@friends.com' --",
	    "password": "doesn't matter"
	}
*/
func (api *API) LoginUser(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer req.Body.Close()

	user := &User{}
	if err := json.Unmarshal(body, user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("Logging in user [%s]", user.Email)
	user, err = GetUser(api.db, user.Email, Hash(user.Password))
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))

	return base64.StdEncoding.EncodeToString(hash[:])
}
