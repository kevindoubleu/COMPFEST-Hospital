package src

import (
	"net/http"
	"time"
)

// session cookie name
var scName string

func init() {
	dbPing()
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	return getJwtClaims(w, r) != nil
}

func isAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims := getJwtClaims(w, r)
	if claims != nil {
		// check if still an active admin
		row := db.QueryRow(`
			SELECT admin FROM users WHERE username = $1`,
			claims.Username)
		var active bool
		row.Scan(&active)

		if active {
			refreshSession(w, r)
			return true
		}
	}
	return false
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
		SameSite: http.SameSiteLaxMode,
		Secure: true,
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
		http.Redirect(w, r, "/login?msg="+ErrMsgSessionTimeout, http.StatusUnauthorized)
	}
}

func createTemplateSessionData(w http.ResponseWriter, r *http.Request) TemplateSessionData {
	claims := getJwtClaims(w, r)
	if claims == nil {
		return TemplateSessionData{}
	} else {
		// check in db if username is admin
		row := db.QueryRow(`
			SELECT admin FROM users WHERE username = $1`,
			claims.Username)
		var isAdmin bool
		row.Scan(&isAdmin)
		
		return TemplateSessionData{
				IsLoggedIn: isLoggedIn(w, r),
				Username: claims.Username,
				IsAdmin: isAdmin,
		}
	}
}