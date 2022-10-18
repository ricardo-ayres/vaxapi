package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func internalServerError(w http.ResponseWriter) {
	s := "Oops, we're having trouble!"
	e := http.StatusInternalServerError
	http.Error(w, s, e)
}

func badRequest(w http.ResponseWriter) {
	s := "Bad request!"
	e := http.StatusBadRequest
	http.Error(w, s, e)
}

func notFound(w http.ResponseWriter) {
	s := "Oops, resource not found!"
	e := http.StatusInternalServerError
	http.Error(w, s, e)
}

func Users(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var user User
		var user_id int64

		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")

		dec := json.NewDecoder(r.Body)

		_, path, pathfound := strings.Cut(r.URL.Path, "/users/")
		if !pathfound {
			http.Error(w, "Path not specified!", http.StatusBadRequest)
			return
		}
		user_id, err = strconv.ParseInt(path, 10, 64)

		switch r.Method {
		case "GET":
			if err != nil {
				internalServerError(w)
				return
			}

			user, err := GetUserById(db, user_id)
			if err == sql.ErrNoRows {
				notFound(w)
				return
			} else if err != nil {
				internalServerError(w)
				log.Printf("Bad query: %s\n", err)
				return
			}

			err = enc.Encode(user)
			if err != nil {
				internalServerError(w)
				log.Printf("Bad encode, %s\n", err)
			}

		case "POST":
			if path != "" {
				badRequest(w)
				return
			}

			err = dec.Decode(&user)
			if err != nil {
				internalServerError(w)
				return
			}

			newuser, err := CreateNewUser(db, user)
			if err != nil {
				internalServerError(w)
				return
			}
			enc.Encode(newuser)

		case "PUT":
		case "DELETE":

		default:
			internalServerError(w)
		}
	}
}

func main() {
	// setup database
	db := SetupDatabase("./sqlite.db")
	defer db.Close()

	// static pages
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/users/", Users(db))

	// start the server
	port := ":5050"
	log.Printf("Running on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
