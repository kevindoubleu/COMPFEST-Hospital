package main

import (
	"net/http"
	"time"
)

// session cookie name
var scName string

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	return getJwtClaims(w, r) != nil
}

func createSession(w http.ResponseWriter, username string) {
	tokenStr := createJwtString(username)

	// create and give session cookie
	http.SetCookie(w, &http.Cookie{
		Name: scName,
		Value: tokenStr,
		Path: "/",
		MaxAge: int(sessionDuration) / int(time.Second),
		HttpOnly: true,
		// Secure: true,
	})
}

// on any authenticated action
func refreshSession(w http.ResponseWriter, r *http.Request) {
	// parse old jwt
	claims := getJwtClaims(w, r)

	if claims != nil {
		// create new session with same user
		createSession(w, claims.Username)
	} else {
		// expired session
		http.Redirect(w, r, "/login?msg="+ErrMsgSessionTimeout, http.StatusSeeOther)
	}
}

func createTemplateSessionData(w http.ResponseWriter, r *http.Request) TemplateSessionData {
	claims := getJwtClaims(w, r)
	if claims == nil {
		return TemplateSessionData{}
	} else {
		return TemplateSessionData{
				IsLoggedIn: isLoggedIn(w, r),
				Username: claims.Username,
		}
	}
}