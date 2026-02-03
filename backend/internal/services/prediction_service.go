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

func (s *PredictionService) GetUserPredictions(userId int) ([]models.Prediction, error) {
	return []models.Prediction{
		{Id: 1, UserId: 5, GameId: "TestGame1", HomeScorePredicted: 4, AwayScorePredicted: 2, TotalScorePredicted: 6, Confidence: 0.7, PredictedWinnerId: 1, ActualWinnerId: 1, WinnerCorrect: true, HomeScoreError: 1, AwayScoreError: 1, TotalScoreError: 1},
		{Id: 2, UserId: 5, GameId: "TestGame2", HomeScorePredicted: 3, AwayScorePredicted: 5, TotalScorePredicted: 8, Confidence: 0.6, PredictedWinnerId: 2, ActualWinnerId: 1, WinnerCorrect: false, HomeScoreError: 2, AwayScoreError: 0, TotalScoreError: 2},
	}, nil
}

func (s *PredictionService) CreatePrediction(writer http.ResponseWriter, request *http.Request) (models.Prediction, error) {
	return models.Prediction{}, nil
}

func (s *PredictionService) SubmitPrediction(ctx context.Context, userId int, httprequest http.Request) (models.Prediction, error) {
	return models.Prediction{}, nil
}

func (s *PredictionService) SubmitBulkPredictions(ctx context.Context, userId int, requests []models.SubmitPredictionRequest) ([]models.Prediction, error) {
	return []models.Prediction{}, nil
}
