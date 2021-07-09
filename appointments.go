package main

import (
	"database/sql"
	"log"
	"net/http"
)

type AppointmentSummary struct {
	Appointment
	RegistrantsCount int
}

type TemplatePatientData struct {
	MyAppointment Appointment
	AppointmentSummaries []AppointmentSummary
}

func init() {
	dbPing()
}

func appointments(w http.ResponseWriter, r *http.Request) {
	refreshSession(w, r)

	if r.Method == http.MethodGet {

		// get patient username
		patientUsername := getJwtClaims(w, r).Username
	
		// get patient's registered appointment if exist
		row := db.QueryRow(`
			SELECT doctor, description
			FROM users
			JOIN appointments
				ON appointment_id = appointments.id
			WHERE users.username = $1`,
			patientUsername)
		patientAppointment := Appointment{}
		row.Scan(&patientAppointment.Doctor, &patientAppointment.Description)
	
		// get appointments
		rows, err := db.Query(`
			SELECT appointments.id, doctor, description, count(users.*) as registrant_count, capacity
			FROM appointments
			LEFT JOIN users
				ON appointment_id = appointments.id
			GROUP BY doctor, description, capacity, appointments.id
			ORDER BY appointments.id
		`)
		ErrPanic(err)
		defer rows.Close()
	
		summaries := make([]AppointmentSummary, 0)
		for rows.Next() {
			s := AppointmentSummary{}
			err := rows.Scan(
				&s.Appointment.Id,
				&s.Appointment.Doctor,
				&s.Appointment.Description,
				&s.RegistrantsCount,
				&s.Appointment.Capacity)
			ErrPanic(err)
			summaries = append(summaries, s)
		}
		ErrPanic(rows.Err())
	
		data := struct{
			TemplateSessionData
			TemplatePatientData
		}{
			createTemplateSessionData(w, r),
			TemplatePatientData{
				patientAppointment,
				summaries,
			},
		}
		tpl.ExecuteTemplate(w, "appointments.gohtml", data)
	}
}

func appointmentsApply(w http.ResponseWriter, r *http.Request) {
	refreshSession(w, r)

	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/appointments", http.StatusMethodNotAllowed)
		return
	}

	if r.Method == http.MethodPost {
		// check if user is not already booked to any appointment
		uname := getJwtClaims(w, r).Username
		row := db.QueryRow(`
			SELECT
			(SELECT appointment_id FROM users WHERE username = $1) = null`,
			uname)
		var id sql.NullInt64
		row.Scan(&id)
		if id.Valid {
			http.Redirect(w, r, "/appointments?msg="+ErrMsgApplyFail, http.StatusBadRequest)
			return
		}

		// check if full
		row = db.QueryRow(`
			SELECT
			(
				SELECT count(*)
				FROM users
				JOIN appointments
					ON appointment_id = appointments.id
				WHERE appointment_id = $1
			) < (
				SELECT capacity
				FROM appointments
				WHERE id = $1
			)`,
			r.PostFormValue("id"))
		var available bool
		err := row.Scan(&available)
		if err != nil || !available {
			if err != nil {
				log.Println("select db error", err)
			}
			http.Redirect(w, r, "/appointments?msg="+ErrMsgApplyFail, http.StatusBadRequest)
			return
		}

		// assign the patient to the appointment
		_, err = db.Exec(`
			UPDATE users
			SET appointment_id = $1
			WHERE username = $2`,
			r.PostFormValue("id"),
			uname)
		if err != nil {
			http.Redirect(w, r, "/appointments?msg="+ErrMsgApplyFail, http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/appointments?msg="+MsgApplySuccess, http.StatusSeeOther)
		return
	}
}

func appointmentsCancel(w http.ResponseWriter, r *http.Request) {
	refreshSession(w, r)

	// just use get bcs we can get unique id (username) from cookie
	if r.Method == http.MethodGet {
		// get patient username
		patientUsername := getJwtClaims(w, r).Username
	
		// update in db
		_, err := db.Exec(`
			UPDATE users
			SET appointment_id = null
			WHERE username = $1`,
			patientUsername)
		if err != nil {
			log.Println("update db error:", err)
			http.Redirect(w, r, "/appointments?msg="+ErrMsgCancelFail, http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/appointments?msg="+MsgCancelSuccess, http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/appointments", http.StatusMethodNotAllowed)
		return
	}
}