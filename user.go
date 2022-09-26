package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

func (api *API) Users(w http.ResponseWriter, req *http.Request) {
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
	ID       int    `json:"id" db:"ID"`
	Name     string `json:"name" db:"Name"`
	Email    string `json:"email" db:"Email"`
	Password string `json:"password" db:"Password"`
	Role     string `json:"role, omitempty" db:"Role"`
}

func (api *API) GetUser(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("userId")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(
		"SELECT id,name,email FROM users WHERE id = '%s' limit 1",
		id,
	)

	log.Println(query)
	rows, err := api.db.Query(query)
	if err != nil {
		log.Print("Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	users, err := api.readUsers(rows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users[0])
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

	users, err := api.readUsers(rows)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (api *API) readUsers(rows *sql.Rows) ([]User, error) {
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
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

	id, err := NewUser(api.db, user.Name, user.Email, Hash(user.Password))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update the user with the new id and clear the password
	user.ID = id
	user.Password = ""

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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (api *API) DeleteUser(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("userId")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	q := fmt.Sprintf("DELETE FROM users WHERE id = '%s'", id)
	_, err := api.db.Exec(q)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := strconv.Atoi(id); err != nil {
		fmt.Print(hacked)
		err = nil
	}

	w.WriteHeader(http.StatusOK)
}
