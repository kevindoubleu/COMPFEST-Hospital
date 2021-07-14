package src

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type AppointmentDetail struct {
	Appointment Appointment
	Registrants []Patient
	RegistrantsCount int
}

func init() {
	dbPing()
}

// CREATE
func adminCreate(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec(`
		INSERT INTO appointments (doctor, description, capacity)
		VALUES
			($1, $2, $3)`,
		r.PostFormValue("doctor"),
		r.PostFormValue("description"),
		r.PostFormValue("capacity"))
	if err != nil {
		log.Println("insert db error:", err)
		http.Redirect(w, r, "/administration?msg="+ErrMsgInsertFail, http.StatusSeeOther)
		return
	}

	// get newest appointment id
	row := db.QueryRow(`SELECT id FROM appointments ORDER BY id DESC LIMIT 1;`)
	var id int
	row.Scan(&id)
	
	if success := addImageToAppointment(r, id); !success {
		http.Redirect(w, r, "/administration?msg="+ErrMsgInsertFail, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/administration?msg="+MsgInsertSuccess, http.StatusSeeOther)
}

// READ
func administration(w http.ResponseWriter, r *http.Request) {
	details := make([]AppointmentDetail, 0)

	// get appointments from db
	rows, err := db.Query(`
		SELECT id, doctor, description, capacity
		FROM appointments ORDER BY id;`)
	ErrPanic(err)
	defer rows.Close()

	for rows.Next() {
		a := Appointment{}
		err := rows.Scan(&a.Id, &a.Doctor, &a.Description, &a.Capacity)
		ErrPanic(err)
		details = append(details, AppointmentDetail{
			Appointment: a,
		})
	}
	ErrPanic(rows.Err())

	// get registrants of each appointment
	for i, a := range details {
		rows, err := db.Query(`
			SELECT username, firstname, lastname, age, email, appointment_id
			FROM users
			JOIN appointments
				ON appointment_id = appointments.id
			WHERE appointments.id = $1`,
			a.Appointment.Id)
		ErrPanic(err)
		defer rows.Close()

		for rows.Next() {
			r := Patient{}
			err := rows.Scan(&r.Username, &r.Firstname, &r.Lastname, &r.Age, &r.Email, &r.Appointment_id)
			ErrPanic(err)
			details[i].Registrants = append(details[i].Registrants, r)

			details[i].RegistrantsCount++
		}
		ErrPanic(rows.Err())
	}

	data := struct{
		TemplateSessionData
		AppointmentDetails []AppointmentDetail
	}{
		createTemplateSessionData(w, r),
		details,
	}

	tpl.ExecuteTemplate(w, "admin-appointments.gohtml", data)
}

// UPDATE
func adminUpdate(w http.ResponseWriter, r *http.Request) {
	// update all fields bcs the original field values
	// are already supplied to frontend
	_, err := db.Exec(`
		UPDATE appointments
		SET doctor = $1,
			description = $2,
			capacity = $3
		WHERE id = $4;`,
		r.PostFormValue("doctor"),
		r.PostFormValue("description"),
		r.PostFormValue("capacity"),
		r.PostFormValue("id"))
	if err != nil {
		log.Println("update db error:", err)
		http.Redirect(w, r, "/administration?msg="+ErrMsgUpdateFail, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/administration?msg="+MsgUpdateSuccess, http.StatusSeeOther)
}

// DELETE
func adminDelete(w http.ResponseWriter, r *http.Request) {
	// unbook patients, delete all images, and delete appointment in db
	_, err1 := db.Exec(`
		UPDATE users
		SET appointment_id = null
		WHERE appointment_id = $1`,
		r.PostFormValue("id"))
	_, err2 := db.Exec(`
		DELETE FROM images WHERE appointment_id = $1`,
		r.PostFormValue("id"))
	_, err3 := db.Exec(`
		DELETE FROM appointments WHERE id = $1`,
		r.PostFormValue("id"))
	if err1 != nil || err2 != nil {
		log.Println("update db error:", err1)
		log.Println("delete db error:", err2, err3)
		http.Redirect(w, r, "/administration?msg="+ErrMsgDeleteFail, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/administration?msg="+MsgDeleteSuccess, http.StatusSeeOther)
}

// URL path is username
func adminKick(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path

	type adminKickResponse struct {
		Registrants int
		Ok bool
	}
	resp := adminKickResponse{}

	// calculate new registrant count on appointment
	row := db.QueryRow(`
		SELECT count(*)
		FROM users
		WHERE appointment_id = (
			SELECT appointment_id FROM users WHERE username = $1)`,
		username)
	row.Scan(&resp.Registrants)
	resp.Registrants -= 1

	// kick the patient
	_, err := db.Exec(`
		UPDATE users
		SET appointment_id = null
		WHERE username = $1`,
		username)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp.Ok = true
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("admin kick:", err)
	}
}

// IMAGE ADD
func adminImagesAdd(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PostFormValue("id"))
	if err != nil {
		http.Redirect(w, r, "/administration?msg="+ErrMsgInsertFail, http.StatusSeeOther)
		return
	}

	if success := addImageToAppointment(r, id); !success {
		http.Redirect(w, r, "/administration?msg="+ErrMsgInsertFail, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/administration?msg="+MsgInsertSuccess, http.StatusSeeOther)
}

// IMAGE DELETE
// URL path is id
func adminImagesDelete(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Ok bool
	}
	id := r.URL.Path
	resp := response{}

	_, err := db.Exec(`DELETE FROM images WHERE appointment_id = $1`, id)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// set ok and send
	resp.Ok = true
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("appointment img delete:", err)
	}
}

// TOGGLE COMMENTS
// URL path is appointment_id
func adminToggleComments(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Ok bool
		Commentable bool
	}
	resp := response{}
	apId := r.URL.Path

	// find current state
	row := db.QueryRow(`SELECT commentable FROM appointments WHERE id = $1`,
		apId)
	var current bool
	row.Scan(&current)

	// set new state
	_, err := db.Exec(`UPDATE appointments SET commentable = $1 WHERE id = $2`,
		!current,
		apId)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// set ok and send response
	resp.Ok = true
	resp.Commentable = !current
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("admin toggle comments:", err)
	}
}