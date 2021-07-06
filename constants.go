package main

var MsgHasSession string
var MsgNoSession string

var MsgRegistered string

var MsgLoggedOut string

func init() {
	// session cookie name
	scName = "gosessid-v0.1"

	// messages for toasts
	MsgHasSession = "You're already logged in"
	MsgNoSession = "You're not logged in"
	MsgRegistered = "You have registered successfully"
	MsgLoggedOut = "You have logged out"
}