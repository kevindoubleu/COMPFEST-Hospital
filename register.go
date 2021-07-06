package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func init() {
	err := db.Ping()
	if err == nil {
		log.Println("register connected to db")
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	// GET -> give form
	if r.Method == http.MethodGet {
		if isLoggedIn(w, r) {
			http.Redirect(w, r, "/?msg="+ErrMsgHasSession, http.StatusSeeOther)
			return
		}

		tpl.ExecuteTemplate(w, "register.gohtml", nil)
		return
	}

	// POST -> process form
	if r.Method == http.MethodPost {
		// check session cookie if logged in
		if isLoggedIn(w, r) {
			http.Redirect(w, r, "/?msg=Already logged in", http.StatusSeeOther)
			return
		}
	
		// check if duplicate username
		row := db.QueryRow(
			"SELECT * FROM patients WHERE username = $1",
			r.PostFormValue("username"))
		if err := row.Scan(); err != sql.ErrNoRows {
			http.Redirect(w, r, "/register?msg="+ErrMsgRegisterFail, http.StatusSeeOther)
			return
		}
		
		// create user in user db
		hash, _ := bcrypt.GenerateFromPassword(
			[]byte(r.PostFormValue("password")),
			bcrypt.DefaultCost)
		age, err := strconv.Atoi(r.PostFormValue("age"))
		if err != nil {
			age = 0
		}
		_, err = db.Exec(`
			INSERT INTO patients (firstname, lastname, age, email, username, password)
			VALUES
				($1, $2, $3, $4, $5, $6)`,
			r.PostFormValue("firstname"),
			r.PostFormValue("lastname"),
			age,
			r.PostFormValue("email"),
			r.PostFormValue("username"),
			string(hash))
		ErrPanic(err)

		// create session so no need to login
		createSession(w, r.PostFormValue("username"))
	
		// redirect to appointments
		http.Redirect(w, r, "/appointments?msg="+MsgRegistered, http.StatusSeeOther)
		return
	}
}