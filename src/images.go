package src

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

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

func addImageToAppointment(r *http.Request, id int) (success bool) {
	// process all images into bytes
	err := r.ParseMultipartForm(1000000)
	if err != nil {
		log.Println("appointment image:", err)
		return false
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
			return false
		}
	}

	return true
}