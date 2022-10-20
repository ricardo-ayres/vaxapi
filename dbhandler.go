package main

import (
	"database/sql"
	"log"
	"os"
)

type User struct {
	UserId   int64  `json:"user_id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Birth    string `json:"birth"`
	Email    string `json:"email"`
	password string
}

type Vaccin struct {
	VacId    int64  `json: "vac_id"`
	Name     string `json: "name"`
	NumDoses string `json: "num_doses"`
}

type Dose struct {
	DoseId    int64  `json: "dose_id"`
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
	var err error

	createUsers := `
		create table users (
		user_id integer primary key,
		username text not null,
		name text not null,
		birth text not null,
		email text unique,
		pwd_hash hash not null,
		pwd_salt hash not null unique,
		check(length(pwd_hash) == 32),
		check(length(pwd_salt) == 32)
		);`

	createVaccines := `
		create table vaccines (
		vac_id integer primary key,
		name text unique not null,
		num_doses integer not null,
		check(num_doses > 0)
		);`

	createDoses := `
		create table doses (
		dose_id integer primary key,
		user_id integer,
		vac_id integer,
		date_taken text,
		foreign key(user_id) references users(user_id),
		foreign key(vac_id) references vaccines(vac_id)
		);`

	_, err = db.Exec(createUsers)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createVaccines)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(createDoses)
	if err != nil {
		log.Fatal(err)
	}
}

func GetUser(db *sql.DB, user User) (User, error) {
	var userdata User
	var creds Credentials

	queryString := "select * from users where username = ?;"
	row := db.QueryRow(queryString, username)
	err := row.Scan(
		&userdata.UserId,
		&userdata.Username,
		&userdata.Name,
		&userdata.Birth,
		&userdata.Email,
		&creds.hash,
		&creds.salt)

	if CheckPassword(user.password, creds) {
		return userdata, err
	} else {
		return user, errors.New("Wrong password.")
	}
}

func CreateNewUser(db *sql.DB, user User) (User, error) {
	var newuser User
	var creds Credentials
	var err error
	
	creds, err = NewCredentials(user.password)
	if err != nil {
		return User, err
	}

	statement := `insert into users (
		username,
		name,
		birth,
		email,
		pwd_hash,
		pwd_salt
		) values (?, ?, ?, ?, ?, ?);`

	result, err := db.Exec(statement,
		user.Username,
		user.Name,
		user.Birth,
		user.Email,
		creds.hash,
		creds.salt)

	if err != nil {
		return newuser, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return newuser, err
	}
	newuser, err = GetUser(db, user)
	return newuser, err
}
