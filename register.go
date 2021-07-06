package main

import (
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

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
		if _, exists := dbUsers[r.PostFormValue("username")]; exists {
			http.Redirect(w, r, "/register?msg=Duplicate username", http.StatusSeeOther)
			return
		}
	
		// create user in user db
		hash, _ := bcrypt.GenerateFromPassword(
			[]byte(r.PostFormValue(r.PostFormValue("password"))),
			bcrypt.DefaultCost)
		age, err := strconv.Atoi(r.PostFormValue("age"))
		if err != nil {
			age = 0
		}
		dbUsers[r.PostFormValue("username")] = Patient{
			fName: r.PostFormValue("firstname"),
			lName: r.PostFormValue("lastname"),
			Age: age,
			Email: r.PostFormValue("email"),
			Username: r.PostFormValue("username"),
			Password: hash,
		}

		// create session so no need to login
		createSession(w, r.PostFormValue("username"))

		// debug
		log.Printf("dbUsers: %#+v\n", dbUsers)
	
		// redirect to appointments
		http.Redirect(w, r, "/appointments?msg="+MsgRegistered, http.StatusSeeOther)
		return
	}
}