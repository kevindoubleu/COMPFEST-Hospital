package main

import (
	"crypto/rand"
	"database/sql"
	"log"
	"time"
)

var db *sql.DB

var ErrMsgGeneric string
var ErrMsgHasSession string
var ErrMsgNoSession string
var ErrMsgSessionTimeout string
var ErrMsgLoginFail string
var ErrMsgRegisterFail string
var ErrMsgChangePasswordFail string

var ErrMsgApplyFail string
var ErrMsgCancelFail string

var ErrMsgInsertFail string
var ErrMsgDeleteFail string
var ErrMsgUpdateFail string

var MsgRegistered string
var MsgLoggedIn string
var MsgLoggedOut string
var MsgChangePasswordSuccess string

var MsgApplySuccess string
var MsgCancelSuccess string

var MsgInsertSuccess string
var MsgUpdateSuccess string
var MsgDeleteSuccess string

var sessionDuration time.Duration

var secretJWT []byte
var secretJWTsize int

func init() {
	db = initDB()

	log.Println("initializing constants")

	// session cookie name
	scName = "gosessid-v0.1"

	// messages and err messages for toasts
	ErrMsgGeneric = "Something went wrong"
	ErrMsgHasSession = "You are already logged in"
	ErrMsgNoSession = "You are not logged in"
	ErrMsgSessionTimeout = "Your session has expired, please login again"
	ErrMsgLoginFail = "Incorrect username or password"
	ErrMsgRegisterFail = "Sorry, that username is taken or something went wrong"
	ErrMsgChangePasswordFail = "Old password didn't match"

	ErrMsgApplyFail = "Sorry, that appointment is fully booked or cancelled"
	ErrMsgCancelFail = "Sorry, we can't cancel that appointment for now"

	ErrMsgInsertFail = "Couldn't create entry in database"
	ErrMsgDeleteFail = "Couldn't delete entry from database"
	ErrMsgUpdateFail = "Couldn't update entry in database"

	MsgRegistered = "You have registered successfully"
	MsgLoggedIn = "You are logged in"
	MsgLoggedOut = "You have logged out"
	MsgChangePasswordSuccess = "Password changed"

	MsgApplySuccess = "You have successfully applied for the appointment"
	MsgCancelSuccess = "You have cancelled your appointment"

	MsgInsertSuccess = "Successfully created entry in database"
	MsgDeleteSuccess = "Successfully deleted entry from database"
	MsgUpdateSuccess = "Successfully updated entry in database"

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