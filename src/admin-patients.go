package src

import (
	"net/http"
	"strconv"
)

// CREATE
func patientsCreate(w http.ResponseWriter, r *http.Request) {
	// use the same controller as in the normal register
	r.ParseForm()
	success, url, code := doRegister(r.PostForm, "/administration/patients")
	if !success {
		http.Redirect(w, r, url, code)
		return
	}
	
	http.Redirect(w, r, "/administration/patients?msg="+MsgInsertSuccess, http.StatusSeeOther)
}

// READ
func patients(w http.ResponseWriter, r *http.Request) {
	// get list of patients
	rows, err := db.Query(`
		SELECT username, firstname, lastname, age, email, appointment_id
		FROM users
		WHERE admin = false
		ORDER BY username
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
}

// UPDATE
func patientsUpdate(w http.ResponseWriter, r *http.Request) {
	// update record in db
	age, _ := strconv.Atoi(r.PostFormValue("age"))
	newUser := Patient{
		Firstname: r.PostFormValue("firstname"),
		Lastname: r.PostFormValue("lastname"),
		Age: age,
		Email: r.PostFormValue("email"),
		Username: r.PostFormValue("username"),
	}
	if r.PostFormValue("password") != "" {
		doProfilePasswordUpdate(
			r.PostFormValue("username"),
			r.PostFormValue("password"),
			"/administration/patients")
	}

	_, url, code := doProfileUpdate(newUser, "/administration/patients")
	http.Redirect(w, r, url, code)
}