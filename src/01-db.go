package src

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"

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
	Firstname string
	Lastname string
	Age int
	Email string
	Username string
	Password string
	Appointment_id sql.NullInt64
}

func initDB() *sql.DB {
	log.Println("initializing database")

	// heroku
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// local
		connStr = `
			user=compfestadmin
			password=compfestadmin
			dbname=hospital
			host=127.0.0.1`
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("successfully connected to db")

	initTables(db)

	initRecords(db)

	testQuery(db)

	log.Println("initialized database", db)
	return db
}

// tables: appointments, users
func initTables(db *sql.DB) {
	_, err := db.Exec(`DROP TABLE IF EXISTS appointments CASCADE`)
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

	// we can add more admins in the future
	// and we can revoke admin rights by setting admin to false
	_, err = db.Exec(`DROP TABLE IF EXISTS users`)
	ErrPanic(err)
	_, err = db.Exec(`
		CREATE TABLE users(
			username       TEXT   PRIMARY KEY NOT NULL,
			appointment_id INT    references appointments(id),
			firstname      TEXT   NOT NULL,
			lastname       TEXT   NOT NULL,
			age            INT    NOT NULL,
			email          TEXT   NOT NULL,
			password       TEXT   NOT NULL,
			admin          BOOL   DEFAULT FALSE
		)
	`)
	ErrPanic(err)

	// appointment images
	_, err = db.Exec(`DROP TABLE IF EXISTS images`)
	ErrPanic(err)
	_, err = db.Exec(`
		CREATE TABLE images(
			id             SERIAL PRIMARY KEY,
			appointment_id INT references appointments(id) NOT NULL,
			img            BYTEA NOT NULL
		)
	`)
	ErrPanic(err)
}

// creates 1 admin, and some dummy records
func initRecords(db *sql.DB) {
	// default admin superuser
	hash, _ := bcrypt.GenerateFromPassword(
		[]byte("compfesthospitaladmin"),
		// []byte("admin"),
		bcrypt.DefaultCost)
	_, err := db.Exec(`
		INSERT INTO users (firstname, lastname, age, email, username, password, admin)
		VALUES
			('admin', 'istrator', 0, 'admin@compfest.local', 'admin', $1, true)`,
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
		[]byte("patient"),
		bcrypt.DefaultCost)
	_, err = db.Exec(`
		INSERT INTO users (firstname, lastname, age, email, username, password, appointment_id)
		VALUES
			('Andi', 'boots', 18, 'andi@email.com', 'aboots', $1, 1),
			('Budi', 'man', 19, 'budi@email.com', 'budiman', $1, 1),
			('Cindy', 'gulla', 20, 'cindy@email.com', 'gulamanis', $1, 1),
			('Deni', 'korbusir', 21, 'deni@email.com', 'corbusir', $1, 2),
			('Eddy', 'gordo', 22, 'eddy@email.com', 'tekken7', $1, 2)`,
		string(tmphash))
	ErrPanic(err)

	var images []interface{}
	for i := 0; i < 4; i++ {
		imgFileName := fmt.Sprintf("assets/img/departments-%d.jpg", i+1)
		imgFile, err := os.Open(imgFileName)
		ErrPanic(err)
		ImgBytes, err := ioutil.ReadAll(imgFile)
		ErrPanic(err)
		images = append(images, ImgBytes)
	}
	_, err = db.Exec(`
		INSERT INTO images (appointment_id, img)
		VALUES
			(1, $1),
			(1, $2),
			(2, $3),
			(1, $4)`,
		images...)
	ErrPanic(err)
}

// run a test query on dummy records
func testQuery(db *sql.DB) {
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
}