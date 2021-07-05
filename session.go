package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

// session cookie name
var scName string

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	// get session cookie
	c, err := r.Cookie(scName)
	if err != nil {
		return false
	}

	// check in sessions db
	_, exists := dbSessions[c.Value]

	return exists
}

func createSession(w http.ResponseWriter, username string) {
	// generate uuid
	sid := uuid.New().String()

	// write entry in sessions db
	dbSessions[sid] = Session{
		Sid: sid,
		Username: username,
		LastModified: time.Now(),
	}

	// create and give session cookie
	http.SetCookie(w, &http.Cookie{
		Name: scName,
		Value: sid,
		Path: "/",
		MaxAge: 60 * 60 * 48, // 48 hour session per request
		HttpOnly: true,
		// Secure: true,
	})
}

// on logout
// func destroySession()

// on any action
// func extendSession()

// on a timer
// func cleanSessionDB()
