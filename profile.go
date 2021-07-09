package main

import (
	"net/http"
	"strconv"

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
		age, _ := strconv.Atoi(r.PostFormValue("age"))
		newUser := Patient{
			Firstname: r.PostFormValue("firstname"),
			Lastname: r.PostFormValue("lastname"),
			Age: age,
			Email: r.PostFormValue("email"),
			Username: uname,
		}
		_, url, code := doProfileUpdate(newUser, "/profile")

		http.Redirect(w, r, url, code)
		return
	}
}

func doProfileUpdate(newUser Patient, prev string) (success bool, url string, code int) {
	_, err := db.Exec(`
		UPDATE users
		SET
			firstname = $1,
			lastname = $2,
			age = $3,
			email = $4
		WHERE
			username = $5`,
		newUser.Firstname,
		newUser.Lastname,
		newUser.Age,
		newUser.Email,
		newUser.Username)
		
	if err != nil {
		return false, prev+"?msg="+ErrMsgUpdateFail, http.StatusInternalServerError
	} else {
		return true, prev+"?msg="+MsgUpdateSuccess, http.StatusSeeOther
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
		if correctPassword(uname, r.PostFormValue("oldpassword")) {
			http.Redirect(w, r, "/profile?msg="+ErrMsgVerifyPasswordFail, http.StatusSeeOther)
			return
		}

		// update with newpassword
		_, url, code := doProfilePasswordUpdate(uname, r.PostFormValue("newpassword"), "/profile")
		http.Redirect(w, r, url, code)
		return
	}
}

// handles the password plaintext hashing
func doProfilePasswordUpdate(username, password string, prev string) (success bool, url string, code int) {
	newHash, _ := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)
	_, err := db.Exec(`
		UPDATE users SET password = $1 WHERE username = $2`,
		string(newHash),
		username)

	if err != nil {
		return false, prev+"?msg="+ErrMsgGeneric, http.StatusInternalServerError
	} else {
		return true, prev+"?msg="+MsgChangePasswordSuccess, http.StatusSeeOther
	}
}

func profileDelete(w http.ResponseWriter, r *http.Request) {
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
		// compare password
		uname := getJwtClaims(w, r).Username
		if !correctPassword(uname, r.PostFormValue("password")) {
			http.Redirect(w, r, "/profile?msg="+ErrMsgVerifyPasswordFail, http.StatusSeeOther)
			return
		}

		// delete from db
		_, err := db.Exec(`
			DELETE FROM users WHERE username = $1`,
			uname)
		if err != nil {
			http.Redirect(w, r, "/profile?msg="+ErrMsgDeleteFail, http.StatusSeeOther)
			return
		}

		// destroy cookie
		destroyJwtCookie(w, r)

		// redirect to homepage
		http.Redirect(w, r, "/?msg="+MsgDeleteSuccess, http.StatusSeeOther)
	}
}