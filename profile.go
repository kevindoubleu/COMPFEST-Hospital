package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func profile(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(w, r) {
		http.Redirect(w, r, "/?msg="+ErrMsgNoSession, http.StatusSeeOther)
		return
	}

	// GET -> give form
	if r.Method == http.MethodGet {
		// get user data from db
		uname := getJwtClaims(w, r).Username
		row := db.QueryRow(`
			SELECT username, firstname, lastname, email, age
			FROM users WHERE username = $1`,
			uname)
		me := Patient{}
		row.Scan(&me.Username, &me.Firstname, &me.Lastname, &me.Email, &me.Age)

		data := struct{
			TemplateSessionData
			Patient
		}{
			createTemplateSessionData(w, r),
			me,
		}
		tpl.ExecuteTemplate(w, "profile.gohtml", data)
	}

	// POST -> process form
	if r.Method == http.MethodPost {
		// update user data in db
		uname := getJwtClaims(w, r).Username
		_, err := db.Exec(`
			UPDATE users
			SET
				firstname = $1,
				lastname = $2,
				age = $3,
				email = $4
			WHERE
				username = $5`,
			r.PostFormValue("firstname"),
			r.PostFormValue("lastname"),
			r.PostFormValue("age"),
			r.PostFormValue("email"),
			uname)
		if err != nil {
			http.Redirect(w, r, "/profile?msg="+ErrMsgUpdateFail, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/profile?msg="+MsgUpdateSuccess, http.StatusSeeOther)
		return
	}
}

func profilePassword(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(w, r) {
		http.Redirect(w, r, "/?msg="+ErrMsgNoSession, http.StatusSeeOther)
		return
	}

	// GET -> return to profile
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/profile", http.StatusSeeOther)
		return
	}

	// POST -> process
	if r.Method == http.MethodPost {
		uname := getJwtClaims(w, r).Username

		// compare oldpassword with old hash
		row := db.QueryRow(`
			SELECT password FROM users WHERE username = $1`,
			uname)
		var hash string
		row.Scan(&hash)

		err := bcrypt.CompareHashAndPassword(
			[]byte(hash), []byte(r.PostFormValue("oldpassword")))
		if err != nil {
			http.Redirect(w, r, "/profile?msg="+ErrMsgChangePasswordFail, http.StatusSeeOther)
			return
		}

		// update with newpassword
		newHash, _ := bcrypt.GenerateFromPassword(
			[]byte(r.PostFormValue("newpassword")), bcrypt.DefaultCost)
		_, err = db.Exec(`
			UPDATE users SET password = $1 WHERE username = $2`,
			string(newHash),
			uname)
		if err != nil {
			http.Redirect(w, r, "/profile?msg="+ErrMsgGeneric, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/profile?msg="+MsgChangePasswordSuccess, http.StatusSeeOther)
		return
	}
}