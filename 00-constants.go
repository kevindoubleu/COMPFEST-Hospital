package main

import (
	"crypto/rand"
	"database/sql"
	"log"
	"time"
)

var db *sql.DB

var toastFail string
var toastSuccess string

var ErrMsgGeneric string
var ErrMsgHasSession string
var ErrMsgNoSession string
var ErrMsgSessionTimeout string
var ErrMsgLoginFail string
var ErrMsgRegisterFail string
var ErrMsgConfirmPasswordFail string
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

	// toast types
	toastFail = "&alert=alert-danger"
	toastSuccess = "&alert=alert-success"

	// messages and err messages for toasts
	ErrMsgGeneric = "Something went wrong"+toastFail
	ErrMsgHasSession = "You are already logged in"+toastFail
	ErrMsgNoSession = "You are not logged in"+toastFail
	ErrMsgSessionTimeout = "Your session has expired, please login again"+toastFail
	ErrMsgLoginFail = "Incorrect username or password"+toastFail
	ErrMsgRegisterFail = "Sorry, that username is taken or something went wrong"+toastFail
	ErrMsgConfirmPasswordFail = "Password and password confirmation doesn't match"+toastFail
	ErrMsgChangePasswordFail = "Old password didn't match"+toastFail

	ErrMsgApplyFail = "Sorry, that appointment is fully booked or cancelled"+toastFail
	ErrMsgCancelFail = "Sorry, we can't cancel that appointment for now"+toastFail

	ErrMsgInsertFail = "Couldn't create entry in database"+toastFail
	ErrMsgDeleteFail = "Couldn't delete entry from database"+toastFail
	ErrMsgUpdateFail = "Couldn't update entry in database"+toastFail

	MsgRegistered = "You have registered successfully"+toastSuccess
	MsgLoggedIn = "You are logged in"+toastSuccess
	MsgLoggedOut = "You have logged out"+toastSuccess
	MsgChangePasswordSuccess = "Password changed"+toastSuccess

	MsgApplySuccess = "You have successfully applied for the appointment"+toastSuccess
	MsgCancelSuccess = "You have cancelled your appointment"+toastSuccess

	MsgInsertSuccess = "Successfully created entry in database"+toastSuccess
	MsgDeleteSuccess = "Successfully deleted entry from database"+toastSuccess
	MsgUpdateSuccess = "Successfully updated entry in database"+toastSuccess

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