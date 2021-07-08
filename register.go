package main

import (
	"database/sql"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	dbPing()
}

func register(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(w, r) {
		http.Redirect(w, r, "/?msg="+ErrMsgHasSession, http.StatusSeeOther)
		return
	}

	// GET -> give form
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "register.gohtml", nil)
		return
	}

	// POST -> process form
	if r.Method == http.MethodPost {
		// check if duplicate username
		row := db.QueryRow(
			"SELECT * FROM users WHERE username = $1",
			r.PostFormValue("username"))
		if err := row.Scan(); err != sql.ErrNoRows {
			http.Redirect(w, r, "/register?msg="+ErrMsgRegisterFail, http.StatusSeeOther)
			return
		}
		
		// create user in user db
		hash, _ := bcrypt.GenerateFromPassword(
			[]byte(r.PostFormValue("password")),
			bcrypt.DefaultCost)
		_, err := db.Exec(`
			INSERT INTO users (firstname, lastname, age, email, username, password)
			VALUES
				($1, $2, $3, $4, $5, $6)`,
			r.PostFormValue("firstname"),
			r.PostFormValue("lastname"),
			r.PostFormValue("age"),
			r.PostFormValue("email"),
			r.PostFormValue("username"),
			string(hash))
		if err != nil {
			http.Redirect(w, r, "/register?msg="+ErrMsgRegisterFail, http.StatusSeeOther)
			return
		}

		// create session so no need to login
		createSession(w, r.PostFormValue("username"))
	
		// redirect to appointments
		http.Redirect(w, r, "/appointments?msg="+MsgRegistered, http.StatusSeeOther)
		return
	}
}