package main

import "net/http"

func logout(w http.ResponseWriter, r *http.Request) {
	if !isLoggedIn(w, r) {
		http.Redirect(w, r, "/?msg="+ErrMsgNoSession, http.StatusUnauthorized)
		return
	}

	// destroy cookie
	destroyJwtCookie(w, r)
	// SECURITY NOTICE: if cookie is not deleted on client side,
	// it can still be used to authenticate until the session expires

	// redirect back to homepage
	http.Redirect(w, r, "/?msg="+MsgLoggedOut, http.StatusSeeOther)
}