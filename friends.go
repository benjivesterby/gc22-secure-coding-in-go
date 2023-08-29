package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
)

func (api *API) GetFriends(
	w http.ResponseWriter,
	req *http.Request,
) {
	userID := req.URL.Query().Get("userId")

	friends, err := GetFriends(userID)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(
			http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(friends)
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(
			http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func GetFriends(userId string) ([]*User, error) {
	q := fmt.Sprintf(`
		SELECT users.*
		FROM users
		JOIN friends ON users.ID = friends.FriendId
		WHERE friends.UserId = '%s';
		`,
		userId,
	)

	return Query[*User](q)
}

func Query[T any](query string) ([]T, error) {
	dbx, err := sqlx.Open("sqlite3", db)
	if err != nil {
		return nil, err
	}

	var results []T
	err = dbx.Select(&results, query)
	if err != nil {
		return nil, err
	}

	return results, nil
}

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
