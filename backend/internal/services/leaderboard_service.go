package services

import (
	"context"

	"github.com/bendemouth/mlb-prediction-pool/internal/database"
	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

type LeaderboardService struct {
	db *database.DB
}

func NewLeaderboardService(db *database.DB) *LeaderboardService {
	return &LeaderboardService{db: db}
}

func (s *LeaderboardService) GetLeaderboard(ctx context.Context) ([]models.LeaderboardEntry, error) {
	// TODO: Implement real calculation
	// Mock data for now
	return []models.LeaderboardEntry{
		{UserId: 1, Username: "bendemouth", TotalWinnersCorrect: 13, WinnerAccuracy: 0.8, TotalScoreError: 51.26, TotalRunsError: 75.68, Rank: 1},
		{UserId: 2, Username: "user2", TotalWinnersCorrect: 11, WinnerAccuracy: 0.7, TotalScoreError: 60.12, TotalRunsError: 80.45, Rank: 2},
		{UserId: 3, Username: "user3", TotalWinnersCorrect: 10, WinnerAccuracy: 0.65, TotalScoreError: 70.34, TotalRunsError: 90.23, Rank: 3},
	}, nil
}

func (s *LeaderboardService) GetUserStats(ctx context.Context, userId int) (*models.LeaderboardEntry, error) {
	// TODO: Implement real query
	return &models.LeaderboardEntry{
		UserId:              userId,
		Username:            "bendemouth",
		TotalWinnersCorrect: 13,
		WinnerAccuracy:      0.8,
		TotalScoreError:     51.26,
		TotalRunsError:      75.68,
		Rank:                1,
	}, nil
}
