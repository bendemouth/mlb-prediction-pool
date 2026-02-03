package services

import (
	"context"
	"net/http"

	"github.com/bendemouth/mlb-prediction-pool/internal/database"
	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

type PredictionService struct {
	db *database.DB
}

func NewPredictionService(db *database.DB) *PredictionService {
	return &PredictionService{db: db}
}

func (s *PredictionService) MakePrediction(request models.SubmitPredictionRequest) error {
	// TODO: Implement actual logic for making predictions

	return nil
}

func (s *PredictionService) GetUserPredictions(userId string) ([]models.Prediction, error) {
	return []models.Prediction{
		{UserId: "1", GameId: "game1", PredictedWinnerId: "teamA", HomeScorePredicted: 5, AwayScorePredicted: 3, TotalScorePredicted: 8},
		{UserId: "1", GameId: "game2", PredictedWinnerId: "teamB", HomeScorePredicted: 2, AwayScorePredicted: 4, TotalScorePredicted: 6},
	}, nil
}

func (s *PredictionService) CreatePrediction(writer http.ResponseWriter, request *http.Request) (models.Prediction, error) {
	return models.Prediction{}, nil
}

func (s *PredictionService) SubmitPrediction(ctx context.Context, userId string, httprequest http.Request) (models.Prediction, error) {
	return models.Prediction{}, nil
}

func (s *PredictionService) SubmitBulkPredictions(ctx context.Context, userId string, requests []models.SubmitPredictionRequest) ([]models.Prediction, error) {
	return []models.Prediction{}, nil
}
