package main

var MsgHasSession string
var MsgRegistered string

func init() {
	// session cookie name
	scName = "gosessid-v0.1"

	// messages for toasts
	MsgHasSession = "Already logged in"
	MsgRegistered = "Registered successfully"
}