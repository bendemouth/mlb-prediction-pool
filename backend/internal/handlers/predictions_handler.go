package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bendemouth/mlb-prediction-pool/internal/models"
	"github.com/bendemouth/mlb-prediction-pool/internal/requests"
)

// Handle GET /predictions
// Eg: /predictions?userId=123
func (h *Handler) GetPredictionsByUser(writer http.ResponseWriter, request *http.Request) {
	userId := request.URL.Query().Get("userId")
	if userId == "" {
		h.respondError(writer, http.StatusBadRequest, "User id is required")
		return
	}

	predictions, err := h.db.GetUserPredictions(request.Context(), userId)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, fmt.Sprint("Failed to get predictions: ", err))
		return
	}

	h.respondJson(writer, http.StatusOK, predictions)
}

// POST /predictions
func (h *Handler) CreatePrediction(writer http.ResponseWriter, request *http.Request) {
	var req requests.SubmitPredictionRequest

	if err := h.decodeJsonBody(request, &req); err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.GameId == "" || req.PredictedWinnerId == "" {
		h.respondError(writer, http.StatusBadRequest, "Game Id and predicted winner are required")
		return
	}

	userId := request.URL.Query().Get("userId")
	if userId == "" {
		h.respondError(writer, http.StatusUnauthorized, "User ID header required")
		return
	}

	// Validate game exists and is upcoming
	game, err := h.db.GetGame(request.Context(), req.GameId)
	if err != nil {
		h.respondError(writer, http.StatusNotFound, "Game not found")
		return
	}

	if game.Status != "upcoming" {
		h.respondError(writer, http.StatusBadRequest, "Cannot predict for games that have started")
		return
	}

	if time.Now().After(game.Date) {
		h.respondError(writer, http.StatusBadRequest, "Game has already started")
		return
	}

	// Validate predicted winner is one of the teams
	if req.PredictedWinnerId != game.HomeTeamId && req.PredictedWinnerId != game.AwayTeamId {
		h.respondError(writer, http.StatusBadRequest, "Invalid predicted winner")
		return
	}

	// Create prediction
	prediction := &models.Prediction{
		UserId:              userId,
		GameId:              req.GameId,
		HomeScorePredicted:  req.HomeScorePredicted,
		AwayScorePredicted:  req.AwayScorePredicted,
		TotalScorePredicted: req.TotalScorePredicted,
		Confidence:          req.Confidence,
		PredictedWinnerId:   req.PredictedWinnerId,
	}

	if err := h.db.CreatePrediction(request.Context(), prediction); err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to create prediction")
		return
	}

	h.respondJson(writer, http.StatusCreated, prediction)
}

// Handle bulk predictions submission
// POST /predictions/bulk
func (h *Handler) CreateBulkPredictions(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var req struct {
		Predictions []requests.SubmitPredictionRequest `json:"predictions"`
	}

	if err := h.decodeJsonBody(request, &req); err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid request body")
		return
	}

	userId := request.URL.Query().Get("userId")
	if userId == "" {
		h.respondError(writer, http.StatusUnauthorized, "User ID header required")
		return
	}

	predictions := make([]models.Prediction, 0, len(req.Predictions))

	for _, prediction := range req.Predictions {
		game, err := h.db.GetGame(request.Context(), prediction.GameId)
		if err != nil {
			h.respondError(writer, http.StatusNotFound, fmt.Sprintf("Game not found: %s", prediction.GameId))
			return
		}

		if game.Status != "upcoming" || time.Now().After(game.Date) {
			h.respondError(writer, http.StatusBadRequest, fmt.Sprintf("Cannot predict for games that have started: %s", prediction.GameId))
			return
		}

		if prediction.PredictedWinnerId != game.HomeTeamId && prediction.PredictedWinnerId != game.AwayTeamId {
			h.respondError(writer, http.StatusBadRequest, fmt.Sprintf("Invalid predicted winner for game: %s", prediction.GameId))
			return
		}

		predictions = append(predictions, models.Prediction{
			UserId:              userId,
			GameId:              prediction.GameId,
			HomeScorePredicted:  prediction.HomeScorePredicted,
			AwayScorePredicted:  prediction.AwayScorePredicted,
			TotalScorePredicted: prediction.TotalScorePredicted,
			Confidence:          prediction.Confidence,
			PredictedWinnerId:   prediction.PredictedWinnerId,
		})
	}

	if err := h.db.BatchCreatePredictions(request.Context(), predictions); err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to create predictions")
		return
	}

	h.respondJson(writer, http.StatusCreated, predictions)
}

// Handle GET /predictions/game?gameId=123
func (h *Handler) GetPredictionsByGame(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	gameId := request.URL.Query().Get("gameId")
	if gameId == "" {
		h.respondError(writer, http.StatusBadRequest, "Game id is required")
		return
	}

	predictions, err := h.db.GetPredictionsByGame(request.Context(), gameId)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, fmt.Sprint("Failed to get predictions: ", err))
		return
	}
	h.respondJson(writer, http.StatusOK, predictions)
}
