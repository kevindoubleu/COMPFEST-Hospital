package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Appointment struct {
	Id int
	Doctor string
	Description string
	Capacity int
}

func init() {
	log.Println("initializing database")

	connStr := `
		user=compfestadmin
		password=compfestadmin
		dbname=hospital
		host=127.0.0.1`
	// db url and user will be provided in an env var in heroku

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("successfully connected to db")

	// init the actual db
	_, err = db.Exec(`DROP TABLE IF EXISTS appointments CASCADE`)
	ErrPanic(err)
	_, err = db.Exec(`
		CREATE TABLE appointments(
			id          SERIAL PRIMARY KEY,
			doctor      TEXT   NOT NULL,
			description TEXT   NOT NULL,
			capacity    INT    NOT NULL
		)
	`)
	ErrPanic(err)
	_, err = db.Exec(`DROP TABLE IF EXISTS patients`)
	ErrPanic(err)
	_, err = db.Exec(`
		CREATE TABLE patients(
			id             SERIAL PRIMARY KEY,
			firstname      TEXT   NOT NULL,
			lastname       TEXT   NOT NULL,
			email          TEXT   NOT NULL,
			age            INT    NOT NULL,
			username       TEXT   NOT NULL,
			password       TEXT   NOT NULL,
			appointment_id INT    references appointments(id)
		)
	`)
	ErrPanic(err)

	// insert dummy values
	_, err = db.Exec(`
		INSERT INTO appointments (doctor, description, capacity)
		VALUES
			('Dr. Some Guy', 'I will talk about some covid stuff', 3),
			('Dr. Pepper', 'some more covid stuff', 4),
			('Mr. Strange', 'hey im a doctor yknow', 5)
	`)
	ErrPanic(err)

	// read
	rows, err := db.Query("SELECT * FROM appointments;")
	ErrPanic(err)
	defer rows.Close()

	availableAppointments := make([]Appointment, 0)
	for rows.Next() {
		a := Appointment{}
		err := rows.Scan(&a.Id, &a.Doctor, &a.Description, &a.Capacity)
		ErrPanic(err)
		availableAppointments = append(availableAppointments, a)
	}

	// check for errors on retrieving rows
	ErrPanic(rows.Err())

	for _, v := range availableAppointments {
		log.Println(v)
	}

	fmt.Println("initialized database")
}