package main

import "log"

type TemplateSessionData struct {
	IsLoggedIn bool
	Username string
	IsAdmin bool
}

func ErrFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func ErrPanic(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func dbPing() {
	err := db.Ping()
	if err != nil {
		log.Panic("can't connect to db")
	}
}