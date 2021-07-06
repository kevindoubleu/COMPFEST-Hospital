package main

import "log"

type TemplateSessionData struct {
	IsLoggedIn bool
	Username string
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