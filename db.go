package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

// not using an actual db yet
// use in memory db for now

// ===== SESSIONS =====
// key is uuid
// val is Patient.Username
var dbSessions map[string]string

// ===== USERS =====
// key is Patient.Username
var dbUsers map[string]Patient

func init() {
	log.Println("initializing database")

	dbSessions = make(map[string]string)
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
	log.Println("admin superuser created, admin:compfesthospitaladmin")

	log.Println("database initialized")
}