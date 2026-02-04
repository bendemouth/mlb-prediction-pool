package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/bendemouth/mlb-prediction-pool/internal/database"
	"github.com/bendemouth/mlb-prediction-pool/internal/services"
)

// Define Handler struct
type Handler struct {
	db                 *database.DB
	healthcheckService *services.HealthcheckService
}

// Create new Handler
func NewHandler(db *database.DB) *Handler {
	return &Handler{
		db:                 db,
		healthcheckService: services.NewHealthcheckService(db),
	}
}

// Encode response as JSON
func (h *Handler) respondJson(writer http.ResponseWriter, status int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)

	if data != nil {
		json.NewEncoder(writer).Encode(data)
	}
}

// Encode error response as JSON
func (h *Handler) respondError(writer http.ResponseWriter, status int, message string) {
	h.respondJson(writer, status, map[string]string{"error": message})
}

// Decode JSON request body
func (h *Handler) decodeJsonBody(request *http.Request, dst interface{}) error {
	return json.NewDecoder(request.Body).Decode(dst)
}
