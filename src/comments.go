package src

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// URL path is appointment_id
func appointmentComments(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Ok bool
		Comments []Comment
	}
	apId := r.URL.Path
	resp := response{}

	// get comments from db
	rows, err := db.Query(`
		SELECT author, comment, posted_at
		FROM comments WHERE appointment_id = $1 ORDER BY posted_at`,
		apId)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	for rows.Next() {
		c := Comment{}
		var tmptime string
		err := rows.Scan(&c.Author, &c.Comment, &tmptime)
		if err != nil {
			json.NewEncoder(w).Encode(resp)
			return
		}
		timestamp, _ := time.Parse(time.RFC3339, tmptime)
		c.Posted_at = timestamp.Format("2 Jan 2006, 15:04")
		resp.Comments = append(resp.Comments, c)
	}
	if rows.Err() != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// set ok and send response
	resp.Ok = true
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("appointment comments:", err)
	}
}

// body is JSON
// {"appointment_id":"1", "comment":"what"}
// response is JSON string of Comment struct
func addCommentToAppointment(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Ok bool
		Comment
	}
	resp := response{}

	c := Comment{}
	err := json.NewDecoder(r.Body).Decode(&c)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// insert to db
	_, err = db.Exec(`
		INSERT INTO comments (author, appointment_id, comment)
		VALUES ($1, $2, $3)`,
		getJwtClaims(w, r).Username,
		c.Appointment_id,
		c.Comment)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// get the value for front end
	row := db.QueryRow(`
		SELECT author, comment, posted_at
		FROM comments ORDER BY id DESC LIMIT 1`)
	var tmptime string
	row.Scan(&c.Author, &c.Comment, &tmptime)
	timestamp, _ := time.Parse(time.RFC3339, tmptime)
	c.Posted_at = timestamp.Format("2 Jan 2006, 15:04")

	// set ok and send
	resp.Ok = true
	resp.Comment = c
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("appointment comments add:", err)
	}
}