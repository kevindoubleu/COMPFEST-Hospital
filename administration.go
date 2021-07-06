package main

import (
	"net/http"
)

type AppointmentDetail struct {
	Appointment Appointment
	Registrants []Patient
	RegistrantsCount int
}

func administration(w http.ResponseWriter, r *http.Request) {
	// validate admin
	claims := getJwtClaims(w, r)
	if claims == nil || claims.Username != "admin" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET -> show crud funcs
	if r.Method == http.MethodGet {
		details := make([]AppointmentDetail, 0)

		// get appointments from db
		rows, err := db.Query("SELECT * FROM appointments;")
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
				SELECT firstname, lastname, age, email
				FROM patients
				JOIN appointments
					ON appointment_id = appointments.id
				WHERE appointments.id = $1`,
				a.Appointment.Id)
			ErrPanic(err)
			defer rows.Close()

			for rows.Next() {
				r := Patient{}
				err := rows.Scan(&r.Firstname, &r.Lastname, &r.Age, &r.Email)
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

		tpl.ExecuteTemplate(w, "administration.gohtml", data)
	}

	// POST
}