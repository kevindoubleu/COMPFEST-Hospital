package main

var ErrMsgHasSession string
var ErrMsgNoSession string
var ErrMsgLoginFail string

var MsgRegistered string
var MsgLoggedIn string
var MsgLoggedOut string

func init() {
	// session cookie name
	scName = "gosessid-v0.1"

	// messages and err messages for toasts
	ErrMsgHasSession = "You are already logged in"
	ErrMsgNoSession = "You are not logged in"
	ErrMsgLoginFail = "Incorrect username or password"

	MsgRegistered = "You have registered successfully"
	MsgLoggedIn = "You are logged in"
	MsgLoggedOut = "You have logged out"
}