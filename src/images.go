package src

import (
	"io/ioutil"
	"log"
	"net/http"
)

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