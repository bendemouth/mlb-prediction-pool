package handlers

import (
	"io"
	"net/http"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	minGbSize := 1
	maxGbSize := 10
	formName := "model-upload-form"

	r.ParseMultipartForm(int64(minGbSize) << int64(maxGbSize))

	file, handler, err := r.FormFile(formName)
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}
}
