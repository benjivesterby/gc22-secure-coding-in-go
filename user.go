package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func (api *API) User(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		api.GetUser(w, req)
	case "POST":
		api.CreateUser(w, req)
	case "PUT":
		api.UpdateUser(w, req)
	case "DELETE":
		api.DeleteUser(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (api *API) GetUser(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	row := api.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var user User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (api *API) CreateUser(w http.ResponseWriter, req *http.Request) {
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

	if user.Name == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := NewUser(api.db, user.Name, user.Email, user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the user with the new id and clear the password
	user.ID = id
	user.Password = ""

	log.Printf("New user <%s> id: %d\n", user.Email, id)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
