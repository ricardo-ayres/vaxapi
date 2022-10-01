package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	//"github.com/jinzhu/gorm"
)

type User struct {
	Username string `json:"username"`
	PwdHash  string `json:"pwdHash"`
	Uid      int    `json:"uid"`
}

// db stand-in
var UserMap = map[int]User{
	0: User{"ric", "asdf", 0},
	1: User{"john", "qwer", 1},
	2: User{"doe", "zxcv", 2},
}

func Users(w http.ResponseWriter, r *http.Request) {
	var err error
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")

	users := map[int]User{}

	if _, path, found := strings.Cut(r.URL.Path, "/users/"); found {
		if path == "" {
			users = UserMap
		} else {
			id, err := strconv.Atoi(path)
			if err != nil {
				log.Printf("Can't convert Atoi(\"%s\")\n", path)
				return
			}

			if user, exists := UserMap[id]; exists {
				users[id] = user
			} else {
				errmsg := fmt.Sprintf("Error: user \"%s\" not found!\n", path)
				strings.NewReader(errmsg).WriteTo(w)
				log.Printf(errmsg)
				return
			}
		}
	}

	result := make([]User, 0, len(users))
	for _, u := range users {
		result = append(result, u)
	}

	err = enc.Encode(result)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	// static pages
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/users/", Users)

	// start the server
	log.Fatal(http.ListenAndServe(":5050", nil))
}
