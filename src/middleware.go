package src

import "net/http"

func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAdmin(w, r) {
			http.Redirect(w, r, "/login", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func LoggedInOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isLoggedIn(w, r) {
			http.Redirect(w, r, "/login?msg="+ErrMsgNoSession, http.StatusUnauthorized)
			return
		}
		refreshSession(w, r)
		next.ServeHTTP(w, r)
	})
}

func LoggedOutOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isLoggedIn(w, r) {
			http.Redirect(w, r, "/?msg="+ErrMsgHasSession, http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func GetOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func PostOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Redirect(w, r, "/", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isLoggedIn(w http.ResponseWriter, r *http.Request) bool {
	return getJwtClaims(w, r) != nil
}

func isAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims := getJwtClaims(w, r)
	if claims != nil {
		// check if still an active admin
		row := db.QueryRow(`
			SELECT admin FROM users WHERE username = $1`,
			claims.Username)
		var active bool
		row.Scan(&active)

		if active {
			refreshSession(w, r)
			return true
		}
	}
	return false
}
