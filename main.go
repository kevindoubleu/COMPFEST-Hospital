package main

import (
	"log"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", homepage)
	http.HandleFunc("/register", register)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)

	http.HandleFunc("/appointments", appointments)

	http.HandleFunc("/administration", administration)
	http.HandleFunc("/administration/create", adminCreate)
	http.HandleFunc("/administration/update", adminUpdate)
	http.HandleFunc("/administration/delete", adminDelete)

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