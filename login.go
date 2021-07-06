package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func login(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
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
		user, exists := dbUsers[r.PostFormValue("username")]
		if !exists {
			http.Redirect(w, r, "/login?msg="+ErrMsgLoginFail, http.StatusSeeOther)
			return
		}
	
		// check password
		err := bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(r.PostFormValue("password")))
		if err != nil {
			http.Redirect(w, r, "/login?msg="+ErrMsgLoginFail, http.StatusSeeOther)
			return
		}
	
		// create session
		createSession(w, user.Username)
	
		// redirect back to homepage
		http.Redirect(w, r, "/?msg="+MsgLoggedIn, http.StatusSeeOther)
		return
	}
}