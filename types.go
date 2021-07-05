package main

type Patient struct {
	fName string
	lName string
	Age int
	Email string
	Username string
	Password []byte
}

type Appointment struct {
	Doctor string
	MaxRegistrant int
	Registrants []*Patient
}
