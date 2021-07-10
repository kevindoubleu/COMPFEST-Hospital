package src

import (
	"log"
	"net/http"
	"os"

	"github.com/urfave/negroni"
)

func Start() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/", homepage)
	http.Handle("/register",
		LoggedOutOnly(http.HandlerFunc(register)))
	http.Handle("/login",
		LoggedOutOnly(http.HandlerFunc(login)))
	http.Handle("/logout",
		LoggedInOnly(http.HandlerFunc(logout)))
	http.Handle("/profile",
		LoggedInOnly(http.HandlerFunc(profile)))
	http.Handle("/profile/password",
		LoggedInOnly(PostOnly(http.HandlerFunc(profilePassword))))
	http.Handle("/profile/delete",
		LoggedInOnly(PostOnly(http.HandlerFunc(profileDelete))))

	http.Handle("/appointments",
		LoggedInOnly(GetOnly(http.HandlerFunc(appointments))))
	http.Handle("/appointments/apply/",
		LoggedInOnly(GetOnly(http.StripPrefix("/appointments/apply/",
		http.HandlerFunc(appointmentsApply)))))
	http.Handle("/appointments/cancel",
		LoggedInOnly(GetOnly(http.HandlerFunc(appointmentsCancel))))

	http.Handle("/administration",
		AdminOnly(GetOnly(http.HandlerFunc(administration))))
	http.Handle("/administration/create",
		AdminOnly(PostOnly(http.HandlerFunc(adminCreate))))
	http.Handle("/administration/update",
		AdminOnly(PostOnly(http.HandlerFunc(adminUpdate))))
	http.Handle("/administration/delete",
		AdminOnly(PostOnly(http.HandlerFunc(adminDelete))))
	http.Handle("/administration/kick",
		AdminOnly(PostOnly(http.HandlerFunc(adminKick))))

	http.Handle("/administration/patients",
		AdminOnly(GetOnly(http.HandlerFunc(patients))))
	http.Handle("/administration/patients/create",
		AdminOnly(PostOnly(http.HandlerFunc(patientsCreate))))
	http.Handle("/administration/patients/update",
		AdminOnly(PostOnly(http.HandlerFunc(patientsUpdate))))

	n := negroni.Classic()
	n.Use(negroni.NewLogger())
	n.UseHandler(http.DefaultServeMux)
	
	log.Printf("starting server")
	port := os.Getenv("PORT")
	if port == "" {
		// local
		log.Fatal(http.ListenAndServe(":8080", n))
	} else {
		// heroku
		log.Fatal(http.ListenAndServe(":"+port, n))
	}
}