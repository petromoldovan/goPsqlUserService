package main

import (
	"fmt"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
	"encoding/json"
)

var db *sql.DB

type User struct {
	ID 					int			`json:"id"`
	FirstName 			string		`json:"first_name"`
	Surname				string		`json:"surname"`
	PhoneNumber			string		`json:"phone_number"`
	Email				string		`json:"email"`
	isActive			bool
}

func init() {
	var err error
	db, err = sql.Open("postgres", "postgres://postgres:postgres@localhost/bonddb?sslmode=disable")
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("db connected")
}

func main() {
	http.HandleFunc("/users", userShowByID)
	http.HandleFunc("/users/show", usersShow)
	http.HandleFunc("/users/create", userCreate)
	http.HandleFunc("/users/delete", userDelete)
	http.HandleFunc("/users/update", userUpdate)
	http.ListenAndServe(":8080", nil)
}

//
// get all users
//
func usersShow(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT * FROM users")

	fmt.Println(rows)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		err := rows.Scan(&user.ID, &user.FirstName, &user.Surname, &user.PhoneNumber, &user.Email, &user.isActive)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	marshaledUsers, err := json.Marshal(users)

	//set cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Server", "User Service")
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(marshaledUsers)
}

//
// get specific user by id
//
func userShowByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(511), http.StatusMethodNotAllowed)
		return
	}

	//reading query param
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	row := db.QueryRow("SELECT * FROM users WHERE id=$1", id)

	user := User{}

	err := row.Scan(&user.ID, &user.FirstName, &user.Surname, &user.PhoneNumber, &user.Email, &user.isActive)
	switch {
	case err == sql.ErrNoRows:
		http.NotFound(w, r)
		return
	case err != nil:
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func userCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	newUser := &User{}
	json.NewDecoder(r.Body).Decode(&newUser)

	if newUser.FirstName == "" || newUser.Surname == "" || newUser.PhoneNumber == "" || newUser.Email == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("INSERT INTO users (first_name, surname, phone_number, email) VALUES ($1, $2, $3, $4)", newUser.FirstName, newUser.Surname, newUser.PhoneNumber, newUser.Email)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	fmt.Println("user created", newUser)
}

func userUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	newUser := &User{}
	json.NewDecoder(r.Body).Decode(&newUser)

	if newUser.FirstName == "" || newUser.Surname == "" || newUser.PhoneNumber == "" || newUser.Email == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE users SET first_name=$1, surname=$2, phone_number=$3, email=$4 WHERE id=$5", newUser.FirstName, newUser.Surname, newUser.PhoneNumber, newUser.Email, newUser.ID)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	fmt.Println("user updated", newUser)
}

func userDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
	}

	id := r.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("DELETE FROM users WHERE id=$1", id)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}

	fmt.Println("user deleted", id)
}
