package main

import "net/http"

func logout(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(w, r) {
		http.Redirect(w, r, "/?msg="+ErrMsgNoSession, http.StatusSeeOther)
		return
	}

	// destroy cookie
	c, _ := r.Cookie(scName)
	c.MaxAge = -1
	http.SetCookie(w, c)

	// delete entry in session db
	delete(dbSessions, c.Value)

	// redirect back to homepage
	http.Redirect(w, r, "/?msg="+MsgLoggedOut, http.StatusSeeOther)
}