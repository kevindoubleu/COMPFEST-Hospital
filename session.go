package main

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// session cookie name
var scName string

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func isLoggedIn(r *http.Request) bool {
	return getJwtClaims(r) != nil
}

func getJwtClaims(r *http.Request) *Claims {
	// get session cookie
	c, err := r.Cookie(scName)
	if err != nil {
		return nil
	}

	// validate jwt
	tokenStr := c.Value
	claims := &Claims{}
	getJwtKey := func(t *jwt.Token) (interface{}, error) {
		return secretJWT, nil
	}
	token, err := jwt.ParseWithClaims(tokenStr, claims, getJwtKey)
	if err != nil {
		return nil
	}

	if !token.Valid {
		return nil
	}

	return claims
}

func createJwtString(username string) string {
	// prepare jwt
	expirationTime := time.Now().Add(sessionDuration)
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// create jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secretJWT)
	if err != nil {
		log.Println("JWT creation error:", err)
	}

	return tokenStr
}

func createSession(w http.ResponseWriter, username string) {
	tokenStr := createJwtString(username)

	// create and give session cookie
	http.SetCookie(w, &http.Cookie{
		Name: scName,
		Value: tokenStr,
		Path: "/",
		MaxAge: int(sessionDuration) / int(time.Second),
		HttpOnly: true,
		// Secure: true,
	})
}

// on logout
// func destroySession()

// on any authenticated action
func refreshSession(w http.ResponseWriter, r *http.Request) {
	// parse old jwt
	claims := getJwtClaims(r)

	// create new session with same user
	createSession(w, claims.Username)
}

// on a timer
// func cleanSessionDB()
