package main

import (
	"crypto/rand"
	"log"
	"time"
)

var ErrMsgHasSession string
var ErrMsgNoSession string
var ErrMsgLoginFail string

var MsgRegistered string
var MsgLoggedIn string
var MsgLoggedOut string

var sessionDuration time.Duration

var secretJWT []byte
var secretJWTsize int

func init() {
	log.Println("initializing constants")

	// session cookie name
	scName = "gosessid-v0.1"

	// messages and err messages for toasts
	ErrMsgHasSession = "You are already logged in"
	ErrMsgNoSession = "You are not logged in"
	ErrMsgLoginFail = "Incorrect username or password"

	MsgRegistered = "You have registered successfully"
	MsgLoggedIn = "You are logged in"
	MsgLoggedOut = "You have logged out"

	// config
	sessionDuration = 15 * time.Minute
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