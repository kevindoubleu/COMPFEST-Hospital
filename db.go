package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// not using an actual db yet
// use in memory db for now

// ===== USERS =====
// key is Patient.Username
// val is Patient data
type Patient struct {
	fName string
	lName string
	Age int
	Email string
	Username string
	Password []byte
}
var dbUsers map[string]Patient

func init() {
	log.Println("initializing database")

	// dbSessions = make(map[string]Session)
	dbUsers = make(map[string]Patient)

	log.Println("creating admin superuser")
	adminHash, _ := bcrypt.GenerateFromPassword(
		[]byte("compfesthospitaladmin"),
		bcrypt.DefaultCost)
	dbUsers["admin"] = Patient{
		fName: "admin",
		lName: "istrator",
		Age: 0,
		Email: "admin@compfest.local",
		Username: "admin",
		Password: adminHash,
	}
	log.Println("created admin superuser, admin:compfesthospitaladmin")

	log.Println("initialized database")
}