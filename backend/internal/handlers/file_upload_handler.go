package handlers

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/bendemouth/mlb-prediction-pool/internal/middleware"
	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

// generateUUID generates a simple UUID-like string
func generateUUID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}

func (h *Handler) UploadModelHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userId, ok := r.Context().Value(middleware.UserSubKey).(string)
	if !ok || userId == "" {
		http.Error(w, "Unauthorized: invalid user context", http.StatusUnauthorized)
		return
	}

	// Parse the multipart form with a max upload size of 500MB
	if err := r.ParseMultipartForm(500 * 1024 * 1024); err != nil {
		http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the model name from form data
	modelName := r.FormValue("modelName")
	if modelName == "" {
		http.Error(w, "Model name is required", http.StatusBadRequest)
		return
	}

	// Get the file from form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to read file from request: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if path.Ext(header.Filename) != ".pkl" {
		http.Error(w, "Invalid file type. Only .pkl files are allowed.", http.StatusBadRequest)
		return
	}

	// Generate model ID
	modelId := generateUUID()

	// Generate S3 key with user ID and timestamp to ensure uniqueness
	s3Key := fmt.Sprintf("models/%s/%d/%s", userId, time.Now().Unix(), header.Filename)

	// Upload file to S3
	success, key, err := h.S3Handler.UploadFileToS3(file, s3Key, r.Context())
	if !success {
		http.Error(w, "Failed to upload file to S3: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create model metadata record in DynamoDB
	model := &models.ModelMetadata{
		ModelId:   modelId,
		ModelName: modelName,
		UserId:    userId,
		FileName:  header.Filename,
		S3Key:     key,
		Status:    "active",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.db.CreateModel(r.Context(), model); err != nil {
		// If DB write fails, we should ideally delete the S3 file, but for now just return error
		http.Error(w, "Failed to save model metadata: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	h.respondJson(w, http.StatusOK, map[string]interface{}{
		"message": "Model uploaded successfully",
		"model":   model,
	})
}
