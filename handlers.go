package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	_ "github.com/glebarez/go-sqlite"
	"log"
	"net/http"
	"pi3/vaxapi/model"
	"strings"
)

// Helper functions
func internalServerError(w http.ResponseWriter, err error) {
	s := err.Error()
	e := http.StatusInternalServerError
	http.Error(w, s, e)
	log.Println(s)
}

func badRequest(w http.ResponseWriter, err error) {
	s := err.Error()
	e := http.StatusBadRequest
	http.Error(w, s, e)
	log.Println(s)
}

func notFound(w http.ResponseWriter, err error) {
	s := err.Error()
	e := http.StatusInternalServerError
	http.Error(w, s, e)
	log.Println(s)
}

func parsePath(r *http.Request, pattern string) string {
	_, path, pathfound := strings.Cut(r.URL.Path, pattern)
	if !pathfound {
		return ""
	}
	return path
}

func credsFromHeader(r *http.Request) (string, string, error) {
	var err error
	var user string
	var password string

	Authorization, ok := r.Header["Authorization"]
	if !ok {
		err = errors.New("Unauthorized")
		return user, password, err
	}

	auth, creds64, ok := strings.Cut(Authorization[0], " ")
	if !ok {
		err = errors.New("Malformed header")
		return user, password, err
	}

	if auth != "Basic" {
		err = errors.New("Bad authorization mode. Only Basic supported.")
		return user, password, err
	}

	tmp := make([]byte, base64.StdEncoding.DecodedLen(len(creds64)))
	n, err := base64.StdEncoding.Decode(tmp, []byte(creds64))
	if err != nil {
		return user, password, err
	}

	creds := string(tmp[:n])
	user, password, ok = strings.Cut(creds, ":")
	if !ok {
		err = errors.New("Malformed Basic authorization token.")
	}

	return user, password, err
}

func sendJson(w http.ResponseWriter, v any) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	return enc.Encode(v)
}

func parseJson(r *http.Request, v any) error {
	dec := json.NewDecoder(r.Body)
	return dec.Decode(v)
}

// Generic Handler
type VaxHandler struct {
	db             *sql.DB
	pattern        string
	methodHandlers map[string]func(VaxCtx)
}

type VaxCtx struct {
	w http.ResponseWriter
	r *http.Request
	h VaxHandler
}

// Implement the Handler interface with ServeHTTP
func (h VaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var ctx = VaxCtx{w: w, r: r, h: h}
	h.methodHandlers[r.Method](ctx)
}

// UsersHandler

// queryUser() handles GET requests for retrieving user data if authenticated
func queryUser(ctx VaxCtx) {
	var err error
	var username string
	var password string
	var userdata model.User

	username, password, err = credsFromHeader(ctx.r)
	if err != nil {
		badRequest(ctx.w, err)
		return
	}

	userdata, err = model.GetUser(ctx.h.db, username, password)
	if err != nil {
		internalServerError(ctx.w, err)
		return
	}

	sendJson(ctx.w, userdata)
}

// requestNewUser() handles POST requests for new users
func requestNewUser(ctx VaxCtx) {
	var newuser model.User
	var err error

	err = parseJson(ctx.r, &newuser)
	if err != nil {
		badRequest(ctx.w, err)
		return
	}

	newuser, err = model.CreateNewUser(ctx.h.db, newuser)
	if err != nil {
		internalServerError(ctx.w, err)
		return
	}

	err = sendJson(ctx.w, newuser)
	if err != nil {
		internalServerError(ctx.w, err)
	}
}

// updateUser() handles PUT requests for updating existing users
func updateUser(ctx VaxCtx) {
	var err error
	var username string
	var password string
	var newdata model.User

	username, password, err = credsFromHeader(ctx.r)
	if err != nil {
		badRequest(ctx.w, err)
		return
	}

	err = parseJson(ctx.r, &newdata)
	if err != nil {
		badRequest(ctx.w, err)
		return
	}

	newdata, err = model.UpdateUser(ctx.h.db, newdata, username, password)
	if err != nil {
		internalServerError(ctx.w, err)
		return
	}

	err = sendJson(ctx.w, newdata)
	if err != nil {
		internalServerError(ctx.w, err)
	}
	return
}

// deleteUser() handles DELETE requests for removing users
func deleteUser(ctx VaxCtx) {
	var err error
	var username string
	var password string

	username, password, err = credsFromHeader(ctx.r)
	if err != nil {
		badRequest(ctx.w, err)
		return
	}

	err = model.DelUser(ctx.h.db, username, password)
	if err != nil {
		internalServerError(ctx.w, err)
	}
	return
}

func NewUsersHandler(db *sql.DB, pattern string) VaxHandler {
	var h VaxHandler
	h.db = db
	h.pattern = pattern
	h.methodHandlers = make(map[string]func(VaxCtx))
	h.methodHandlers["GET"] = queryUser
	h.methodHandlers["POST"] = requestNewUser
	h.methodHandlers["PUT"] = updateUser
	h.methodHandlers["DELETE"] = deleteUser
	return h
}

/* a simple closure to allow access to the db variable */
func VacHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vax, err := model.GetVac(db)
		if err != nil {
			internalServerError(w, err)
		}
		sendJson(w, vax)
		return
	}
}
