package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type User struct {
	Id      int    `json:"id"`
	Name string `json:"name"`
	Birth string `json:"birth"`
	Email string `json:"email"`
}

func Users(w http.ResponseWriter, r *http.Request) {
	var err error
	var user User
	db := Db
	queryString := "select * from users"
	users := make([]User, 0)

	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")

	if _, path, found := strings.Cut(r.URL.Path, "/users/"); found {
		id, err := strconv.Atoi(path)
		if err != nil {
			log.Printf("Can't convert Atoi(\"%s\")\n", path)
		} else {
			queryString = fmt.Sprintf("%s where id = %d", queryString, id)
		}

		rows, err := db.Query(queryString)
		if err != nil {
			log.Printf("Bad query, %s\n", err)
		}

		for rows.Next() {
			err := rows.Scan(&user.Id, &user.Name, &user.Birth, &user.Email)
			if err != nil {
				log.Printf("Bad scan, %s\n", err)
			}
			users = append(users, user)
		}
	}
	err = enc.Encode(users)
	if err != nil {
		log.Printf("Bad encode, %s\n", err)
	}
}

var Db *sql.DB
func main() {
	// setup database
	createTables := false
	if _, err := os.Stat("./sqlite.db"); err != nil {
		createTables = true
		log.Println("Database does not exist. Will create")
	}

	db, err := sql.Open("sqlite", "sqlite.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	Db = db

	if createTables {
		createUsers := `
			create table users (
			id integer primary key,
			name text,
			birth text,
			email text unique
			);`
		res, err := db.Exec(createUsers)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(res)
	}

	// static pages
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/users/", Users)

	// start the server
	port := ":5050"
	log.Printf("Running on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
