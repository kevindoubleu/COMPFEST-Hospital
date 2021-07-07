package main

import (
	"log"
	"net/http"
)

type AppointmentDetail struct {
	Appointment Appointment
	Registrants []Patient
	RegistrantsCount int
}

func init() {
	dbPing()
}

func isAdmin(w http.ResponseWriter, r *http.Request) bool {
	claims := getJwtClaims(w, r)
	return claims != nil && claims.Username == "admin"
}

// CREATE
func adminCreate(w http.ResponseWriter, r *http.Request) {
	// validate admin
	if !isAdmin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET -> not accepted
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/administration", http.StatusSeeOther)
		return
	}

	// POST -> insert to db
	if r.Method == http.MethodPost {
		_, err := db.Exec(`
			INSERT INTO appointments (doctor, description, capacity)
			VALUES
				($1, $2, $3)`,
			r.PostFormValue("doctor"),
			r.PostFormValue("description"),
			r.PostFormValue("capacity"))
		if err != nil {
			log.Println("insert db error:", err)
			http.Redirect(w, r, "/administration?msg="+ErrMsgInsertFail, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/administration?msg="+MsgInsertSuccess, http.StatusSeeOther)
		return
	}
}

// READ
func administration(w http.ResponseWriter, r *http.Request) {
	// validate admin
	if !isAdmin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET -> show crud funcs
	if r.Method == http.MethodGet {
		details := make([]AppointmentDetail, 0)

		// get appointments from db
		rows, err := db.Query("SELECT * FROM appointments ORDER BY id;")
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

	// POST -> not acepted
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/administration", http.StatusSeeOther)
		return
	}
}

// UPDATE
func adminUpdate(w http.ResponseWriter, r *http.Request) {
	// validate admin
	if !isAdmin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET -> not accepted
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/administration", http.StatusSeeOther)
		return
	}

	// POST -> delete appointment based on id
	if r.Method == http.MethodPost {
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
			http.Redirect(w, r, "/administration?msg="+ErrMsgUpdateFail, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/administration?msg="+MsgUpdateSuccess, http.StatusSeeOther)
		return
	}
}

// DELETE
func adminDelete(w http.ResponseWriter, r *http.Request) {
	// validate admin
	if !isAdmin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET -> not accepted
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/administration", http.StatusSeeOther)
		return
	}

	// POST -> delete appointment based on id
	if r.Method == http.MethodPost {
		// unbook patients and delete appointment in db
		_, err1 := db.Exec(`
			UPDATE patients
			SET appointment_id = null
			WHERE appointment_id = $1`,
			r.PostFormValue("id"))
		_, err2 := db.Exec(`
			DELETE FROM appointments WHERE id = $1`,
			r.PostFormValue("id"))
		if err1 != nil || err2 != nil {
			log.Println("update db error:", err1)
			log.Println("delete db error:", err2)
			http.Redirect(w, r, "/administration?msg="+ErrMsgInsertFail, http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/administration?msg="+MsgDeleteSuccess, http.StatusSeeOther)
		return
	}
}