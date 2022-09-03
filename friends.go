package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (api *API) Friends(w http.ResponseWriter, req *http.Request) {
	log.Print("Users")
	switch req.Method {
	case "GET":
		api.GetFriends(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (api *API) Friend(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "PUT":
		api.AddFriend(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (api *API) AddFriend(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("userId")
	friendID := req.URL.Query().Get("friendId")

	log.Printf("Adding friend [%s] to database for user [%s]", friendID, userID)
	err := AddFriend(api.db, userID, friendID)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (api *API) GetFriends(w http.ResponseWriter, req *http.Request) {
	userID := req.URL.Query().Get("userId")

	log.Printf("Getting friends for user [%s]", userID)
	friends, err := GetFriends(userID)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(friends)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
