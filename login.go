package main

import (
	"database/sql"
	"log"
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
	}

	// POST -> process form
	if r.Method == http.MethodPost {
		// check in user db
		row := db.QueryRow(
			"SELECT * FROM patients WHERE username = $1",
			r.PostFormValue("username"))
		if err := row.Scan(); err == sql.ErrNoRows {
			http.Redirect(w, r, "/login?msg="+ErrMsgLoginFail, http.StatusSeeOther)
			return
		}
	
		// check password
		row = db.QueryRow(
			"SELECT password FROM patients WHERE username = $1",
			r.PostFormValue("username"))
		var hash string
		row.Scan(&hash)
		err := bcrypt.CompareHashAndPassword(
			[]byte(hash),
			[]byte(r.PostFormValue("password")))
		if err != nil {
			http.Redirect(w, r, "/login?msg="+ErrMsgLoginFail, http.StatusSeeOther)
			log.Printf("err: %#+v\n", err)
			return
		}
	
		// create session
		createSession(w, r.PostFormValue("username"))
	
		// redirect back to homepage
		http.Redirect(w, r, "/?msg="+MsgLoggedIn, http.StatusSeeOther)
		return
	}
}