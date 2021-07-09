package main

import (
	"database/sql"
	"net/http"
	"net/url"

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
		// parse the form for db insertion
		r.ParseForm()
		success, url, code := doRegister(r.PostForm, "/register")
		if !success {
			http.Redirect(w, r, url, code)
			return
		}

		// create session so no need to login
		createSession(w, r.PostFormValue("username"))
	
		// redirect to appointments
		http.Redirect(w, r, "/appointments?msg="+MsgRegistered, http.StatusSeeOther)
		return
	}
}

func doRegister(postFormData url.Values, prevUrl string) (success bool, redirect string, statusCode int) {
	// confirm password
	if postFormData.Get("password") != postFormData.Get("confirmpassword") {
		return false, prevUrl+"?msg="+ErrMsgConfirmPasswordFail, http.StatusSeeOther
	}

	// check if duplicate username
	row := db.QueryRow(
		"SELECT * FROM users WHERE username = $1",
		postFormData.Get("username"))
	if err := row.Scan(); err != sql.ErrNoRows {
		return false, prevUrl+"?msg="+ErrMsgConfirmPasswordFail, http.StatusSeeOther
	}

	// create user in user db
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte(postFormData.Get("password")),
		bcrypt.DefaultCost)
	_, err := db.Exec(`
		INSERT INTO users (firstname, lastname, age, email, username, password)
		VALUES
			($1, $2, $3, $4, $5, $6)`,
		postFormData.Get("firstname"),
		postFormData.Get("lastname"),
		postFormData.Get("age"),
		postFormData.Get("email"),
		postFormData.Get("username"),
		string(hash))

	return err == nil, prevUrl+"?msg="+ErrMsgConfirmPasswordFail, http.StatusSeeOther
}