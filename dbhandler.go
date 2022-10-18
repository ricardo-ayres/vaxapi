package main

import (
	"database/sql"
	"log"
	"os"
)

type User struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Birth string `json:"birth"`
	Email string `json:"email"`
}

type Vaccin struct {
	Id       int64  `json: "id"`
	Name     string `json: "name"`
	NumDoses string `json: "num_doses"`
}

type Dose struct {
	Id        int64  `json: "id"`
	UserId    int64  `json: "user_id"`
	VacId     int64  `json: "vac_id"`
	DateTaken string `json: "date_taken"`
}

func SetupDatabase(filename string) *sql.DB {
	dbMissing := false

	if _, err := os.Stat(filename); err != nil {
		log.Println("Database does not exist. Will create")
		dbMissing = true
	}

	db, err := sql.Open("sqlite", filename)
	if err != nil {
		log.Fatal(err)
	}

	if dbMissing {
		createTables(db)
	}

	return db
}

func createTables(db *sql.DB) {
	var res sql.Result
	var err error

	createUsers := `
		create table users (
		id integer primary key,
		name text,
		birth text,
		email text unique
		);`

	createVaccines := `
		create table vaccines (
		id integer primary key,
		name text,
		num_doses integer
		);`

	createDoses := `
		create table doses (
		id integer primary key,
		user_id integer foreign key,
		vac_id integer foreign key,
		date_taken text
		);`

	res, err = db.Exec(createUsers)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	res, err = db.Exec(createVaccines)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	res, err = db.Exec(createDoses)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}

func GetUserById(db *sql.DB, user_id int64) (User, error) {
	var user User
	queryString := "select id, name, birth, email from users where id = ?;"
	row := db.QueryRow(queryString, user_id)
	err := row.Scan(&user.Id, &user.Name, &user.Birth, &user.Email)
	return user, err
}

func CreateNewUser(db *sql.DB, user User) (User, error) {
	var newuser User

	statement := "insert into users (name, birth, email) values (?, ?, ?);"
	result, err := db.Exec(statement, user.Name, user.Birth, user.Email)
	if err != nil {
		return newuser, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return newuser, err
	}
	newuser, err = GetUserById(db, id)
	return newuser, err
}
