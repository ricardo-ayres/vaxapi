package main

import (
	"log"
	"net/http"
	"os"
	"pi3/vaxapi/model"
)

func main() {
	// setup database
	db := model.SetupDatabase("./sqlite.db")
	defer db.Close()

	// static pages
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/users/", NewUsersHandler(db, "/users/"))
	http.HandleFunc("/vacs/", VacHandler(db))

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
