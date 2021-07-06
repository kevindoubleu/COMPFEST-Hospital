package main

import (
	"log"
	"net/http"
	"text/template"
)

var tpl *template.Template
func init() {
	log.Println("initializing templates")

	tpl =  template.New("")
	tpl.ParseGlob("templates/components/*.gohtml")
	tpl.ParseGlob("templates/sections/*.gohtml")
	tpl.ParseGlob("templates/pages/*.gohtml")
	tpl.ParseGlob("templates/layouts/*.gohtml")

	log.Println("parsed templates", tpl.DefinedTemplates())
}

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/test", test)

	http.HandleFunc("/", homepage)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	
	log.Printf("starting server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(w, r) {
		c, _ := r.Cookie(scName)
		data := struct{
			TemplateSessionData TemplateSessionData
		}{
			TemplateSessionData{
				IsLoggedIn: true,
				Username: dbSessions[c.Value].Username,
			},
		}

		tpl.ExecuteTemplate(w, "index.gohtml", data)
	} else {
		tpl.ExecuteTemplate(w, "index.gohtml", nil)
	}
}

func test(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "inner-page.gohtml", nil)
}
