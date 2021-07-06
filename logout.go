package main

import "net/http"

func logout(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(r) {
		http.Redirect(w, r, "/?msg="+ErrMsgNoSession, http.StatusSeeOther)
		return
	}

	// destroy cookie
	c, _ := r.Cookie(scName)
	c.MaxAge = -1
	http.SetCookie(w, c)
	// SECURITY NOTICE: if cookie is not deleted, it can still be used to auth

	// redirect back to homepage
	http.Redirect(w, r, "/?msg="+MsgLoggedOut, http.StatusSeeOther)
}