package main

import (
	"net/http"
)

// CREATE
func patientsCreate(w http.ResponseWriter, r *http.Request) {
	// validate admin
	if !isAdmin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET -> not accepted
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/administration/patients", http.StatusSeeOther)
		return
	}

	// POST -> process form
	if r.Method == http.MethodPost {
		// use the same controller as in the normal register
		r.ParseForm()
		success, url, code := doRegister(r.PostForm, "/administration/patients")
		if !success {
			http.Redirect(w, r, url, code)
			return
		}
		
		http.Redirect(w, r, "/administration/patients?msg="+MsgInsertSuccess, http.StatusSeeOther)
		return
	}
}

// READ
func patients(w http.ResponseWriter, r *http.Request) {
	// validate admin
	if !isAdmin(w, r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// GET -> read patients list
	if r.Method == http.MethodGet {
		// get list of patients
		rows, err := db.Query(`
			SELECT username, firstname, lastname, age, email, appointment_id
			FROM users
		`)
		ErrPanic(err)
		defer rows.Close()

		patients := make([]Patient, 0)
		for rows.Next() {
			u := Patient{}
			err := rows.Scan(
				&u.Username, &u.Firstname, &u.Lastname, &u.Age, &u.Email, &u.Appointment_id)
			ErrPanic(err)
			patients = append(patients, u)
		}
		ErrPanic(rows.Err())

		data := struct{
			TemplateSessionData
			Patients []Patient
		}{
			createTemplateSessionData(w, r),
			patients,
		}

		tpl.ExecuteTemplate(w, "admin-patients.gohtml", data)
		return
	}

	// POST -> not accepted
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/administration/patients", http.StatusSeeOther)
		return
	}
}