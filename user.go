package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/davecgh/go-spew/spew"
)

func (api *API) Users(w http.ResponseWriter, req *http.Request) {
	log.Print("Users")
	switch req.Method {
	case "GET":
		api.GetUsers(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (api *API) User(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		api.GetUser(w, req)
	case "POST":
		api.UpdateUser(w, req)
	case "PUT":
		api.CreateUser(w, req)
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
	Role     string `json:"role, omitempty"`
}

func (api *API) GetUser(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("SELECT id,name,email FROM users WHERE id = '%s'", id)

	log.Print("Get user id: ", id)
	row := api.db.QueryRow(query)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Print("Loaded User: ", spew.Sdump(user))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (api *API) GetUsers(w http.ResponseWriter, req *http.Request) {
	if req.URL.Query().Get("isAdmin") != "1" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	rows, err := api.db.Query("SELECT id,name,email FROM users")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
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

	log.Print("Created User: ", spew.Sdump(user))

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (api *API) UpdateUser(w http.ResponseWriter, req *http.Request) {
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

	if user.ID == 0 || user.Name == "" || user.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = api.db.Exec(`
		UPDATE users SET name = ?, email = ? WHERE id = ?
	`,
		user.Name, user.Email, user.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Print("Updated User: ", spew.Sdump(user))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (api *API) DeleteUser(w http.ResponseWriter, req *http.Request) {
	safe := req.URL.Query().Get("safe")

	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if safe == "" {
		// VULN: curl -X DELETE 'localhost:8081/user?id="7%27%20or%201%3D1--"'
		// 7' or 1=1--
		_, err := api.db.Exec(fmt.Sprintf("DELETE FROM users WHERE id = '%s'", id))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if _, err := strconv.Atoi(id); err != nil {
			fmt.Print(hacked)
			err = nil
		}
	} else {
		_, err := api.db.Exec("DELETE FROM users WHERE id = '%s'", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	// Should really be
	// _, err := api.db.Exec("DELETE FROM users WHERE id = '?'", id)

	log.Print("Deleted User: ", id)

	w.WriteHeader(http.StatusOK)
}
