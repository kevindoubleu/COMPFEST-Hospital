package main

import (
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

func appointments(w http.ResponseWriter, r *http.Request) {
	refreshSession(w, r)

	if r.Method == http.MethodGet {

		// get patient username
		patientUsername := getJwtClaims(w, r).Username
	
		// get patient's registered appointment if exist
		row := db.QueryRow(`
			SELECT doctor, description
			FROM patients
			JOIN appointments
				ON appointment_id = appointments.id
			WHERE patients.username = $1`,
			patientUsername)
		patientAppointment := Appointment{}
		row.Scan(&patientAppointment.Doctor, &patientAppointment.Description)
	
		// get appointments
		rows, err := db.Query(`
			SELECT doctor, description, count(patients.id) as registrant_count, capacity
			FROM appointments
			LEFT JOIN patients
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

func appointmentsCancel(w http.ResponseWriter, r *http.Request) {
	refreshSession(w, r)

	if r.Method == http.MethodGet {
		// get patient username
		patientUsername := getJwtClaims(w, r).Username
	
		// update in db
		_, err := db.Exec(`
			UPDATE patients
			SET appointment_id = null
			WHERE username = $1`,
			patientUsername)
		if err != nil {
			log.Println("update db error:", err)
			http.Redirect(w, r, "/appointments?msg="+ErrMsgDeleteFail, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/appointments?msg="+MsgDeleteSuccess, http.StatusSeeOther)
		return
	}
}