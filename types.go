package main

type Patient struct {
	fName string
	lName string
	Age int
	Email string
	Username string
	Password []byte
}

type Doctor struct {
	Name string
}

type Appointment struct {
	Doctor Doctor
	Patients [3]*Patient
}
