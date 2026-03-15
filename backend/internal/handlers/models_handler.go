package handlers

import (
	"net/http"
	"strings"

	"github.com/bendemouth/mlb-prediction-pool/internal/middleware"
)

func (h *Handler) GetUserModelsHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userId, ok := r.Context().Value(middleware.UserSubKey).(string)
	if !ok || userId == "" {
		h.respondError(w, http.StatusUnauthorized, "Unauthorized: invalid user context")
		return
	}

	// Get all models for the user
	userModels, err := h.db.GetModelsByUserId(r.Context(), userId)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to retrieve models: "+err.Error())
		return
	}

	// Ensure we return a JSON array, not null
	var response interface{}
	if len(userModels) == 0 {
		response = []interface{}{}
	} else {
		response = userModels
	}

	h.respondJson(w, http.StatusOK, response)
}

func (h *Handler) DeleteModelHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userId, ok := r.Context().Value(middleware.UserSubKey).(string)
	if !ok || userId == "" {
		h.respondError(w, http.StatusUnauthorized, "Unauthorized: invalid user context")
		return
	}

	// Get model ID from URL path (e.g., /models/delete/{modelId})
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		h.respondError(w, http.StatusBadRequest, "Invalid request path")
		return
	}

	modelId := parts[len(parts)-1]
	if modelId == "" {
		h.respondError(w, http.StatusBadRequest, "Model ID is required")
		return
	}

	// Verify the model exists and belongs to the user
	model, err := h.db.GetModelById(r.Context(), modelId, userId)
	if err != nil || model == nil {
		h.respondError(w, http.StatusNotFound, "Model not found")
		return
	}

	// Delete the model from DynamoDB
	if err := h.db.DeleteModel(r.Context(), modelId, userId); err != nil {
		h.respondError(w, http.StatusInternalServerError, "Failed to delete model: "+err.Error())
		return
	}

	// TODO: Delete the file from S3
	// success, _, err := h.S3Handler.DeleteFileFromS3(model.S3Key, r.Context())
	// if !success {
	//     // Log error but don't fail the request - model is already deleted from DB
	//     log.Printf("Failed to delete S3 file: %v", err)
	// }

	h.respondJson(w, http.StatusOK, map[string]string{
		"message": "Model deleted successfully",
	})
}

func (h *Handler) GetModelHandler(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userId, ok := r.Context().Value(middleware.UserSubKey).(string)
	if !ok || userId == "" {
		h.respondError(w, http.StatusUnauthorized, "Unauthorized: invalid user context")
		return
	}

	// Get model ID from URL path (e.g., /models/{modelId})
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		h.respondError(w, http.StatusBadRequest, "Invalid request path")
		return
	}

	modelId := parts[len(parts)-1]
	if modelId == "" {
		h.respondError(w, http.StatusBadRequest, "Model ID is required")
		return
	}

	// Get the model
	model, err := h.db.GetModelById(r.Context(), modelId, userId)
	if err != nil || model == nil {
		h.respondError(w, http.StatusNotFound, "Model not found")
		return
	}

	h.respondJson(w, http.StatusOK, model)
}
