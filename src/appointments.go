package src

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

// URL path is appointment_id
func appointmentsApply(w http.ResponseWriter, r *http.Request) {
	type appointmentApplyResponse struct {
		Appointment_id int
		Registrants int
		Doctor string
		Description string
		Ok bool
		ErrMsg string
	}
	resp := appointmentApplyResponse{}

	apId, err := strconv.Atoi(r.URL.Path)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp.Appointment_id = apId

	// check if user is not already booked to any appointment
	uname := getJwtClaims(w, r).Username
	row := db.QueryRow(`
		SELECT appointment_id FROM users WHERE username = $1`,
		uname)
	var id sql.NullInt64
	row.Scan(&id)
	if id.Valid {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// check if full
	row = db.QueryRow(`
		SELECT count(users.*), capacity
		FROM appointments
		LEFT JOIN users
			ON appointments.id = appointment_id
		WHERE id = $1
		GROUP BY capacity`,
		apId)
	var current, max int
	err = row.Scan(&current, &max)
	if err != nil || current >= max {
		if err != nil {
			log.Println("select db error", err)
		}
		resp.ErrMsg = ErrMsgApplyFail[:len(ErrMsgApplyFail)-len(toastFail)]
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp.Registrants = current+1

	// assign the patient to the appointment
	_, err = db.Exec(`
		UPDATE users
		SET appointment_id = $1
		WHERE username = $2`,
		apId,
		uname)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	// get the doctor and description to show
	row = db.QueryRow(`
		SELECT doctor, description
		FROM appointments
		WHERE id = $1`,
		apId)
	row.Scan(&resp.Doctor, &resp.Description)

	// set the response ok and send
	resp.Ok = true
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("appointment apply:", err)
	}
}

func appointmentsCancel(w http.ResponseWriter, r *http.Request) {
	type appointmentCancelResponse struct {
		Appointment_id int
		Registrants int
		Ok bool
	}
	resp := appointmentCancelResponse{}

	// get patient username
	patientUsername := getJwtClaims(w, r).Username

	// get old appointment id
	row := db.QueryRow(`
		SELECT appointment_id FROM users WHERE username = $1`,
		patientUsername)
	row.Scan(&resp.Appointment_id)

	// count new registrant
	row = db.QueryRow(`
		SELECT count(*)
		FROM users
		JOIN appointments
			ON appointment_id = appointments.id
		WHERE appointment_id = $1`,
		resp.Appointment_id)
	row.Scan(&resp.Registrants)
	resp.Registrants -= 1

	// update in db
	_, err := db.Exec(`
		UPDATE users
		SET appointment_id = null
		WHERE username = $1`,
		patientUsername)
	if err != nil {
		log.Println("update db error:", err)
		json.NewEncoder(w).Encode(resp)
		return
	}
	
	// set the response ok and send
	resp.Ok = true
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println("appointment cancel:", err)
	}
}

// URL path is appointment_id
func appointmentImages(w http.ResponseWriter, r *http.Request) {
	type image struct{
		Id int
		Base64 string
	}
	type response struct{
		Ok bool
		Images []image
	}
	resp := response{}

	apId := r.URL.Path
	rows, err := db.Query(`
		SELECT id, img
		FROM images
		WHERE appointment_id = $1`,
		apId)
	if err != nil {
		json.NewEncoder(w).Encode(resp)
		return
	}

	for i := 0; rows.Next(); i++ {
		resp.Images = append(resp.Images, image{})
		var imgBytes []byte
		err := rows.Scan(&resp.Images[i].Id, &imgBytes)
		resp.Images[i].Base64 = base64.StdEncoding.EncodeToString(imgBytes)
		if err != nil {
			log.Println("appointment images:", err)
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	if rows.Err() != nil {
		log.Println("appointment images:", err)
		json.NewEncoder(w).Encode(resp)
		return
	}

	// set ok and send
	resp.Ok = true
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err)
	}
}