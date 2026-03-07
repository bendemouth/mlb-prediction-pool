package database

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

// CalculateLeaderboard recalculates the leaderboard based on user scores.
func (db *DB) CalculateLeaderboard(ctx context.Context) ([]models.LeaderboardEntry, error) {
	users, err := db.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	leaderboard := make([]models.LeaderboardEntry, 0, len(users))

	for _, user := range users {
		predictions, err := db.GetUserPredictions(ctx, user.Id)
		if err != nil {
			return nil, fmt.Errorf("failed to get predictions for user %s: %w", user.Id, err)
		}

		winnerAccuracy, totalWinnersCorrect := calculateWinnerAccuracyAndTotalCorrectWinners(predictions)
		totalScoreMse := calculateTotalScoreMse(predictions)
		teamScoreMse := calculateTeamScoreMse(predictions)
		leaderboardScore := getLeaderboardScore(predictions)

		leaderboard = append(leaderboard, models.LeaderboardEntry{
			UserId:              user.Id,
			Username:            user.Username,
			TotalWinnersCorrect: totalWinnersCorrect,
			WinnerAccuracy:      winnerAccuracy,
			TeamScoreMse:        totalScoreMse,
			TotalRunsMse:        teamScoreMse,
			LeaderboardScore:    leaderboardScore,
			Rank:                0, // Rank will be assigned later
		})
	}

	sortLeaderboardEntries(leaderboard)

	for i := range leaderboard {
		leaderboard[i].Rank = i + 1
	}

	return leaderboard, nil
}

// GetUserStats retrieves statistics for a specific user
func (db *DB) GetUserStats(ctx context.Context, userId string) (*models.LeaderboardEntry, error) {
	user, err := db.GetUser(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	predictions, err := db.GetUserPredictions(ctx, user.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user predictions: %w", err)
	}

	winnerAccuracy, totalWinnersCorrect := calculateWinnerAccuracyAndTotalCorrectWinners(predictions)
	totalScoreError := calculateTotalScoreMse(predictions)
	totalRunsError := calculateTeamScoreMse(predictions)

	leaderboard, err := db.CalculateLeaderboard(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate leaderboard: %w", err)
	}

	rank := 0
	for _, entry := range leaderboard {
		if entry.UserId == user.Id {
			rank = entry.Rank
			break
		}
	}

	return &models.LeaderboardEntry{
		UserId:              user.Id,
		Username:            user.Username,
		TotalWinnersCorrect: totalWinnersCorrect,
		WinnerAccuracy:      winnerAccuracy,
		TeamScoreMse:        totalScoreError,
		TotalRunsMse:        totalRunsError,
		Rank:                rank,
	}, nil
}

func calculateWinnerAccuracyAndTotalCorrectWinners(predictions []models.Prediction) (winnerAccuracy float32, totalWinnersCorrect int) {
	var totalPredictions int
	for _, pred := range predictions {
		if pred.WinnerCorrect != nil {
			totalPredictions++
			if *pred.WinnerCorrect {
				totalWinnersCorrect++
			}
		}
	}
	if totalPredictions > 0 {
		winnerAccuracy = float32(totalWinnersCorrect) / float32(totalPredictions)
	}
	return winnerAccuracy, totalWinnersCorrect
}

func calculateTeamScoreMse(predictions []models.Prediction) float32 {
	var totalPredictions int
	var sumSquaredErrors float32

	for _, pred := range predictions {
		if pred.WinnerCorrect == nil {
			continue
		}
		sumSquaredErrors += (pred.HomeScoreError * pred.HomeScoreError) + (pred.AwayScoreError * pred.AwayScoreError)
		totalPredictions++
	}
	if totalPredictions == 0 {
		return 0
	}
	mse := float64(sumSquaredErrors / float32(2*totalPredictions))
	return float32(math.Sqrt(mse))
}

func calculateTotalScoreMse(predictions []models.Prediction) float32 {
	var totalPredictions int
	var sumSquaredErrors float32

	for _, pred := range predictions {
		if pred.WinnerCorrect == nil {
			continue
		}
		sumSquaredErrors += pred.TotalScoreError * pred.TotalScoreError
		totalPredictions++
	}
	if totalPredictions == 0 {
		return 0
	}
	mse := float64(sumSquaredErrors / float32(totalPredictions))
	return float32(math.Sqrt(mse))
}

func getLeaderboardScore(predictions []models.Prediction) (leaderboardScore float32) {
	winnerAccuracyWeight := float32(0.6)
	teamScoreMseWeight := float32(0.2)
	totalScoreMseWeight := float32(0.2)

	winnerAccuracy, _ := calculateWinnerAccuracyAndTotalCorrectWinners(predictions)
	teamScoreRmse := calculateTeamScoreMse(predictions)
	totalScoreRmse := calculateTotalScoreMse(predictions)

	// Normalize RMSE to a 0-1 score using exponential decay.
	// A lower RMSE yields a score closer to 1, higher RMSE closer to 0.
	// The decay constant controls sensitivity — tune as needed.
	teamScoreComponent := float32(math.Exp(float64(-teamScoreRmse) / 3.0))
	totalScoreComponent := float32(math.Exp(float64(-totalScoreRmse) / 3.0))

	leaderboardScore = (winnerAccuracy * winnerAccuracyWeight) +
		(teamScoreComponent * teamScoreMseWeight) +
		(totalScoreComponent * totalScoreMseWeight)

	return leaderboardScore
}

func sortLeaderboardEntries(entries []models.LeaderboardEntry) {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].LeaderboardScore > entries[j].LeaderboardScore {
			return true
		}
		if entries[i].LeaderboardScore == entries[j].LeaderboardScore {
			if entries[i].TeamScoreMse < entries[j].TeamScoreMse {
				return true
			}
			if entries[i].TeamScoreMse == entries[j].TeamScoreMse {
				return entries[i].TotalRunsMse < entries[j].TotalRunsMse
			}
		}
		return false
	})
}
