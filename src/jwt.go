package src

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func destroyJwtCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie(scName)
	if err == http.ErrNoCookie {
		return
	}
	c.MaxAge = -1
	http.SetCookie(w, c)
}

func getJwtClaims(w http.ResponseWriter, r *http.Request) *Claims {
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
		destroyJwtCookie(w, r)
		return nil
	}

	if !token.Valid {
		destroyJwtCookie(w, r)
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