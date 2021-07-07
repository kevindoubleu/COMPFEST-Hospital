package main

import (
	"crypto/rand"
	"database/sql"
	"log"
	"time"
)

var db *sql.DB

var ErrMsgHasSession string
var ErrMsgNoSession string
var ErrMsgSessionTimeout string
var ErrMsgLoginFail string
var ErrMsgRegisterFail string

var ErrMsgInsertFail string

var MsgRegistered string
var MsgLoggedIn string
var MsgLoggedOut string

var MsgInsertSuccess string

var sessionDuration time.Duration

var secretJWT []byte
var secretJWTsize int

func init() {
	db = initDB()

	log.Println("initializing constants")

	// session cookie name
	scName = "gosessid-v0.1"

	// messages and err messages for toasts
	ErrMsgHasSession = "You are already logged in"
	ErrMsgNoSession = "You are not logged in"
	ErrMsgSessionTimeout = "Your session has expired, please login again"
	ErrMsgLoginFail = "Incorrect username or password"
	ErrMsgRegisterFail = "Sorry, that username is taken"

	ErrMsgInsertFail = "Couldn't create entry in database"

	MsgRegistered = "You have registered successfully"
	MsgLoggedIn = "You are logged in"
	MsgLoggedOut = "You have logged out"

	MsgInsertSuccess = "Successfully created entry in database"

	// config
	sessionDuration = 15 * time.Minute
	// sessionDuration = 30 * time.Second // testing
	// small session duration due to difficulty ininvalidating jwt

	// secrets
	secretJWTsize = 2048
	secretJWT := make([]byte, secretJWTsize)
	n, err := rand.Read(secretJWT)
	if err != nil || n != secretJWTsize {
		log.Fatalln("failed to initialize JWT secret:", err)
	}

	log.Println("initialized constants")
}