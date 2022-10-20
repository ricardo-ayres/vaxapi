package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func internalServerError(w http.ResponseWriter, err) {
	s := err.Error()
	e := http.StatusInternalServerError
	http.Error(w, s, e)
}

func badRequest(w http.ResponseWriter, err) {
	s := err.Error()
	e := http.StatusBadRequest
	http.Error(w, s, e)
}

func notFound(w http.ResponseWriter, err) {
	s := err.Error()
	e := http.StatusInternalServerError
	http.Error(w, s, e)
}

func parsePath(r *http.Request, pattern string) string {
		_, path, pathfound := strings.Cut(r.URL.Path, pattern)
		if !pathfound {
			http.Error(w, "Path not specified!", http.StatusBadRequest)
			return nil
		}
}

func sendJson(w http.ResponseWriter, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(v)
}

func parseJson(r *http.Request, v *any) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(&v)
}

func Users(db *sql.DB, pattern string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var user User

		path := parsePath(r, pattern)
		// WIP
	}
}

func main() {
	// setup database
	db := SetupDatabase("./sqlite.db")
	defer db.Close()

	// static pages
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/users/", Users(db, "/users/"))

	// decide what port to use, run on :80 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	port = ":" + port

	// start the server
	log.Printf("Running on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
