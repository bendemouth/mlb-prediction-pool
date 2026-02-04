package handlers

import (
	"fmt"
	"net/http"

	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

func (h *Handler) HandlePredictions(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		h.getPredictions(writer, request)
	case http.MethodPost:
		h.createPrediction(writer, request)
	default:
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Handle GET /predictions
// Eg: /predictions?userId=123
func (h *Handler) getPredictions(writer http.ResponseWriter, request *http.Request) {
	userId := request.URL.Query().Get("userId")
	if userId == "" {
		h.respondError(writer, http.StatusBadRequest, "User id is required")
		return
	}

	predictions, err := h.predictionService.GetUserPredictions(userId)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, fmt.Sprint("Failed to get predictions: ", err))
		return
	}

	h.respondJson(writer, http.StatusOK, predictions)
}

// Handle POST /predictions
func (h *Handler) SubmitPredictions(writer http.ResponseWriter, request *http.Request) {
	// TODO: Make a better way to submit predictions with user authentication token
	var req models.SubmitPredictionRequest

	if err := h.decodeJsonBody(request, &req); err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.GameId == "" || req.PredictedWinnerId == "" {
		h.respondError(writer, http.StatusBadRequest, "Game Id and predicted winner are required")
	}

	// TODO: JWT Token?
	userId := request.Header.Get("X-User-Id")

	prediction, err := h.predictionService.SubmitPrediction(request.Context(), userId, *request)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, fmt.Sprint("Failed to submit prediction: ", err))
		return
	}

	h.respondJson(writer, http.StatusCreated, prediction)
}

// Handle bulk predictions submission
// POST /predictions/bulk

func (h *Handler) HandleBulkPredictions(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Predictions []models.SubmitPredictionRequest `json:"predictions"`
	}

	if err := h.decodeJsonBody(request, &req); err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	userId := request.Header.Get("X-User-Id")

	results, err := h.predictionService.SubmitBulkPredictions(request.Context(), userId, req.Predictions)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, fmt.Sprint("Failed to submit bulk predictions: ", err))
		return
	}

	h.respondJson(writer, http.StatusCreated, results)
}
