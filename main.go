package main

import (
	"net/http"
	"log"
)

func main() {
	// setup database
	db := SetupDatabase("./sqlite.db")
	defer db.Close()

	// static pages
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/users/", Users(db))

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