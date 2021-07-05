package main

type Appointment struct {
	Doctor string
	MaxRegistrant int
	Registrants []*Patient
}

type TemplateSessionData struct {
	IsLoggedIn bool
	Username string
}