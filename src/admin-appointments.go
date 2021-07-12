package src

import (
	"encoding/json"
	"io/ioutil"
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

	// process all images into bytes
	err = r.ParseMultipartForm(1000000)
	if err != nil {
		log.Println("appointment image:", err)
		http.Redirect(w, r, "/administration?msg="+ErrMsgInsertFail, http.StatusSeeOther)
		return
	}
	var imagesBytes [][]byte
	fileHeaders := r.MultipartForm.File["images"]
	for _, fh := range fileHeaders {
		f, err := fh.Open()
		if err != nil {
			break
		}
		defer f.Close()

		imgBytes, err := ioutil.ReadAll(f)
		if err != nil {
			break
		}
		imagesBytes = append(imagesBytes, imgBytes)
	}

	// get newest appointment id
	row := db.QueryRow(`SELECT id FROM appointments ORDER BY id DESC LIMIT 1;`)
	var id int
	row.Scan(&id)
	// insert images
	for _, img := range imagesBytes {
		_, err = db.Exec(`
			INSERT INTO images (appointment_id, img)
			VALUES
				($1, $2)`,
			id,
			img)
		if err != nil {
			log.Println("appointment image:", err)
			http.Redirect(w, r, "/administration?msg="+ErrMsgInsertFail, http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/administration?msg="+MsgInsertSuccess, http.StatusSeeOther)
}

// READ
func administration(w http.ResponseWriter, r *http.Request) {
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