package model

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"pi3/vaxapi/auth"
)

type User struct {
	UserId   int    `json:"user_id,omitempty"`
	Username string `json:"username,omitempty"`
	Name     string `json:"name,omitempty"`
	Birth    string `json:"birth,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type Vaccin struct {
	VacId    int64  `json: "vac_id"`
	Name     string `json: "name"`
	NumDoses string `json: "num_doses"`
	Obs      string `json: "obs"`
}

type Dose struct {
	DoseId    int64  `json: "dose_id"`
	UserId    int64  `json: "user_id"`
	VacId     int64  `json: "vac_id"`
	DateTaken string `json: "date_taken"`
}

func SetupDatabase(filename string) *sql.DB {
	if _, err := os.Stat(filename); err != nil {
		log.Println("Database does not exist. Will create")
	}

	db, err := sql.Open("sqlite", filename)
	if err != nil {
		log.Fatal(err)
	}

	createTables(db)

	return db
}

func createTables(db *sql.DB) {
	var err error

	pragmaForeignKeys := `pragma foreign_keys = on;`

	createUsers := `
		create table if not exists users (
		user_id integer not null primary key,
		username text not null unique,
		name text not null,
		birth text not null,
		email text not null unique,
		pwd_hash binary not null,
		pwd_salt binary not null unique,
		check(length(pwd_hash) == 32),
		check(length(pwd_salt) == 32)
		);`

	createVaccines := `
		create table if not exists vaccines (
		vac_id integer not null primary key,
		name text unique not null,
		num_doses integer not null,
		obs text,
		check(num_doses > 0)
		);`

	createDoses := `
		create table if not exists doses (
		dose_id integer not null primary key,
		user_id integer,
		vac_id integer,
		date_taken text,
		foreign key(user_id) references users(user_id) on delete cascade,
		foreign key(vac_id) references vaccines(vac_id)
		);`

	_, err = db.Exec(pragmaForeignKeys)
	if err != nil {
		log.Fatal(err)
	}

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

/*
 * FUNCTIONS FOR TABLE users
 */
func CreateNewUser(db *sql.DB, user User) (User, error) {
	var newuser User
	var creds auth.Credentials
	var err error

	creds, err = auth.NewCredentials(user.Password)
	if err != nil {
		return user, err
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
		creds.Hash[:],
		creds.Salt[:])

	if err != nil {
		return newuser, err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return newuser, err
	}

	newuser, err = GetUser(db, user.Username, user.Password)
	return newuser, err
}

func GetUser(db *sql.DB, username string, password string) (User, error) {
	var userdata User
	var blank User
	var creds auth.Credentials

	hash := make([]byte, auth.Size)
	salt := make([]byte, auth.Size)

	queryString := `select * from users where username = ?;`
	row := db.QueryRow(queryString, username)
	err := row.Scan(
		&userdata.UserId,
		&userdata.Username,
		&userdata.Name,
		&userdata.Birth,
		&userdata.Email,
		&hash,
		&salt)

	if err != nil {
		return blank, err
	}

	n := copy(creds.Hash[:], hash)
	if n != auth.Size {
		return blank, errors.New("Hash copying error")
	}

	n = copy(creds.Salt[:], salt)
	if n != auth.Size {
		return blank, errors.New("Salt copying error")
	}

	if auth.CheckPassword(password, creds) {
		return userdata, err
	} else {
		return blank, errors.New("Wrong password.")
	}
}

func UpdateUser(db *sql.DB, newdata User, username string, password string) (User, error) {
	var err, newerr error
	var olddata User
	var blank User
	var creds auth.Credentials

	updateCreds := `;` // do nothing by default
	updateData := `update users set
		name = ?,
		birth = ?,
		email = ?
		where user_id = ?;`

	olddata, err = GetUser(db, username, password)
	if err != nil {
		return blank, err
	}

	/* Copy old data over blank/missing fields */
	if newdata.Name == "" {
		newdata.Name = olddata.Name
	}
	if newdata.Birth == "" {
		newdata.Birth = olddata.Birth
	}
	if newdata.Email == "" {
		newdata.Email = olddata.Email
	}

	/* If password updated: we need to generate a new hash and a new salt */
	if newdata.Password != "" {
		creds, err = auth.NewCredentials(newdata.Password)
		if err != nil {
			return olddata, err
		}

		updateCreds = `update users set
			pwd_hash = ?,
			pwd_salt = ?
			where user_id = ?;`

		newdata.Password = "UPDATED"
	}

	tx, err := db.Begin()
	if err != nil {
		return blank, err
	}

	/* update credentials */
	_, err = tx.Exec(updateCreds,
		creds.Hash[:],
		creds.Salt[:],
		olddata.UserId)
	if err != nil {
		newerr = tx.Rollback()
		if newerr != nil {
			err = newerr
		}
		return blank, err
	}

	/* update data */
	_, err = tx.Exec(updateData,
		newdata.Name,
		newdata.Birth,
		newdata.Email,
		olddata.UserId)

	if err != nil {
		newerr = tx.Rollback()
		if newerr != nil {
			err = newerr
		}
		return blank, err
	}

	tx.Commit()
	return newdata, err
}

func DelUser(db *sql.DB, username string, password string) error {
	var err error
	var userdata User
	statement := `delete from users where user_id = ?;`

	userdata, err = GetUser(db, username, password)
	if err != nil {
		return err
	}

	_, err = db.Exec(statement, userdata.UserId)
	if err != nil {
		return err
	}

	return err
}

/*
 * FUNCTIONS FOR TABLE vaccines
 */
func GetVac(db *sql.DB) ([]Vaccin, error) {
	var err error
	var vac Vaccin

	/*
	 * O SUS oferece 20 tipos diferentes de vacinas
	 * durante a vida, alocando um slice com cap de 32
	 * para evitar realocações nos appends
	 */
	vax := make([]Vaccin, 0, 32)
	statement := `select * from vaccines`

	rows, err := db.Query(statement)
	if err != nil {
		return vax, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&vac.VacId, &vac.Name, &vac.NumDoses)
		if err != nil {
			return vax, err
		}
		vax = append(vax, vac)
	}

	return vax, err
}

/* GetDoses sql statement: */
/* SELECT users.name, vaccines.name, vaccines.num_doses, doses.date_taken
FROM doses JOIN users ON doses.user_id=users.user_id
JOIN vaccines ON doses.vac_id=vaccines.vac_id
WHERE users.user_id=?
*/
