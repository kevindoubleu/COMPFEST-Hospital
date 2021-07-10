package src

import "net/http"

func homepage(w http.ResponseWriter, r *http.Request) {
	claims := getJwtClaims(w, r)
	if claims != nil {
		// check if user still in db
		row := db.QueryRow(`
			SELECT username FROM users WHERE username = $1`,
			claims.Username)
		var uname string
		row.Scan(&uname)

		if uname == "" {
			destroyJwtCookie(w, r)
		} else {
			refreshSession(w, r)
		}
	}

	data := struct{
		TemplateSessionData TemplateSessionData
	}{
		createTemplateSessionData(w, r),
	}
	tpl.ExecuteTemplate(w, "index.gohtml", data)
}