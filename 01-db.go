package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Appointment struct {
	Id int
	Doctor string
	Description string
	Capacity int
}

type Patient struct {
	Id int
	Firstname string
	Lastname string
	Age int
	Email string
	Username string
	Password string
	Appointment_id int
}

func initDB() *sql.DB {
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

	// reserve "admin" username
	hash, _ := bcrypt.GenerateFromPassword(
		// []byte("compfesthospitaladmin"),
		[]byte("admin"),
		bcrypt.DefaultCost)
	_, err = db.Exec(`
		INSERT INTO patients (firstname, lastname, age, email, username, password)
		VALUES
			('admin', 'istrator', 0, 'admin@compfest.local', 'admin', $1)`,
		string(hash))
	ErrPanic(err)

	// insert dummy values
	_, err = db.Exec(`
		INSERT INTO appointments (doctor, description, capacity)
		VALUES
			('Dr. Some Guy', 'I will talk about some covid stuff', 3),
			('Dr. Pepper', 'some more covid stuff', 5),
			('Mr. Strange', 'hey im a doctor yknow', 20)
	`)
	ErrPanic(err)
	tmphash, _ := bcrypt.GenerateFromPassword(
		// []byte("compfesthospitaladmin"),
		[]byte("andi"),
		bcrypt.DefaultCost)
	_, err = db.Exec(`
		INSERT INTO patients (firstname, lastname, age, email, username, password, appointment_id)
		VALUES
			('Andi', 'boots', 18, 'andi@email.com', 'aboots', $1, 1),
			('Budi', 'man', 19, 'budi@email.com', 'budiman', '', 1),
			('Cindy', 'gulla', 20, 'cindy@email.com', 'gulamanis', '', 1),
			('Deni', 'korbusir', 21, 'budi@email.com', 'corbusir', '', 2),
			('Eddy', 'gordo', 22, 'budi@email.com', 'tekken7', '', 2)`,
		string(tmphash))
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

	log.Println("initialized database", db)
	return db
}