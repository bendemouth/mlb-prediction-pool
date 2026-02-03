package handlers

import (
	"fmt"
	"net/http"
	"strconv"

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
	userIdQuery := request.URL.Query().Get("userId")
	if userIdQuery == "" {
		h.respondError(writer, http.StatusBadRequest, "User id is required")
		return
	}

	userId, err := strconv.Atoi(userIdQuery)
	if err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid user id")
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
func (h *Handler) submitPredictions(writer http.ResponseWriter, request *http.Request) {
	// TODO: Make a better way to submit predictions with user authentication token
	var req models.SubmitPredictionRequest

	if err := h.decodeJsonBody(request, &req); err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.GameId == "" || req.PredictedWinnerId == 0 {
		h.respondError(writer, http.StatusBadRequest, "Game Id and predicted winner are required")
	}

	// TODO: JWT Token?
	userId := request.Header.Get("X-User-Id")

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid user id")
		return
	}

	prediction, err := h.predictionService.SubmitPrediction(request.Context(), userIdInt, *request)
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

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid user id")
		return
	}

	results, err := h.predictionService.SubmitBulkPredictions(request.Context(), userIdInt, req.Predictions)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, fmt.Sprint("Failed to submit bulk predictions: ", err))
		return
	}

	h.respondJson(writer, http.StatusCreated, results)
}
