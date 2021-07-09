package main

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	dbPing()
}

func login(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(w, r) {
		http.Redirect(w, r, "/?msg="+ErrMsgHasSession, http.StatusSeeOther)
		return
	}
	
	// GET -> give form
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "login.gohtml", nil)
		return
	}

	// POST -> process form
	if r.Method == http.MethodPost {
		// check in user db
		row := db.QueryRow(
			"SELECT * FROM users WHERE username = $1",
			r.PostFormValue("username"))
		if err := row.Scan(); err == sql.ErrNoRows {
			http.Redirect(w, r, "/login?msg="+ErrMsgLoginFail, http.StatusSeeOther)
			return
		}
	
		// check password
		if !correctPassword(r.PostFormValue("username"), r.PostFormValue("password")) {
			http.Redirect(w, r, "/login?msg="+ErrMsgLoginFail, http.StatusSeeOther)
			return
		}
	
		// create session
		createSession(w, r.PostFormValue("username"))
	
		// redirect back to homepage
		http.Redirect(w, r, "/?msg="+MsgLoggedIn, http.StatusSeeOther)
		return
	}
}

func correctPassword(username, password string) bool {
	row := db.QueryRow(`
		SELECT password FROM users WHERE username = $1`,
		username)
	var hash string
	row.Scan(&hash)

	err := bcrypt.CompareHashAndPassword(
		[]byte(hash), []byte(password))
	return err == nil
}