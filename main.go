package main

import (
	"log"
	"net/http"
	"text/template"

	"github.com/urfave/negroni"
)

var tpl *template.Template
func init() {
	log.Println("initializing templates")

	tpl = template.New("")
	tpl.ParseGlob("templates/components/*.gohtml")
	tpl.ParseGlob("templates/sections/*.gohtml")
	tpl.ParseGlob("templates/pages/*.gohtml")
	tpl.ParseGlob("templates/layouts/*.gohtml")

	log.Println("initialized templates")
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/test", test)

	http.HandleFunc("/", homepage)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/appointments", appointments)

	n := negroni.Classic()
	n.Use(negroni.NewLogger())
	n.UseHandler(http.DefaultServeMux)
	
	log.Printf("starting server")
	log.Fatal(http.ListenAndServe(":8080", n))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie(scName)
	if err != http.ErrNoCookie {
		refreshSession(w, r)
	}

	data := struct{
		TemplateSessionData TemplateSessionData
	}{
		createTemplateSessionData(w, r),
	}
	tpl.ExecuteTemplate(w, "index.gohtml", data)
}

func appointments(w http.ResponseWriter, r *http.Request) {
	refreshSession(w, r)

	data := struct{
		TemplateSessionData TemplateSessionData
	}{
		createTemplateSessionData(w, r),
	}
	tpl.ExecuteTemplate(w, "appointments.gohtml", data)
}

func test(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "inner-page.gohtml", nil)
}

func createTemplateSessionData(w http.ResponseWriter, r *http.Request) TemplateSessionData {
	claims := getJwtClaims(w, r)
	if claims == nil {
		return TemplateSessionData{}
	} else {
		return TemplateSessionData{
				IsLoggedIn: isLoggedIn(w, r),
				Username: claims.Username,
		}
	}
}