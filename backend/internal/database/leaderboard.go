package database

import (
	"context"
	"fmt"
	"sort"

	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

// CalculateLeaderboard computes leaderboard from predictions
func (db *DB) CalculateLeaderboard(ctx context.Context) ([]models.LeaderboardEntry, error) {
	// Get all users
	users, err := db.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	leaderboard := make([]models.LeaderboardEntry, 0, len(users))

	for _, user := range users {
		// Get user's predictions
		predictions, err := db.GetUserPredictions(ctx, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get predictions for user %s: %w", user.Username, err)
		}

		// Calculate stats
		var correct, total int
		for _, pred := range predictions {
			if pred.WinnerCorrect {
				correct++
				total++
			} else if !pred.WinnerCorrect {
				total++
			}
			// Skip "pending" predictions
		}

		var accuracy float32
		if total > 0 {
			accuracy = float32(correct) / float32(total)
		}

		leaderboard = append(leaderboard, models.LeaderboardEntry{
			UserId:              user.ID,
			Username:            user.Username,
			TotalWinnersCorrect: correct,
			WinnerAccuracy:      accuracy,
			// TODO: Add logic for calculating error
		})
	}

	// Sort by accuracy (descending), then by total (descending)
	sort.Slice(leaderboard, func(i, j int) bool {
		if leaderboard[i].WinnerAccuracy != leaderboard[j].WinnerAccuracy {
			return leaderboard[i].WinnerAccuracy > leaderboard[j].WinnerAccuracy
		}
		return leaderboard[i].TotalWinnersCorrect > leaderboard[j].TotalWinnersCorrect
	})

	// Assign ranks
	for i := range leaderboard {
		leaderboard[i].Rank = i + 1
	}

	return leaderboard, nil
}

// GetUserStats retrieves stats for a specific user
func (db *DB) GetUserStats(ctx context.Context, userID string) (*models.LeaderboardEntry, error) {
	user, err := db.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	predictions, err := db.GetUserPredictions(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get predictions: %w", err)
	}

	var correct, total int
	for _, pred := range predictions {
		if pred.WinnerCorrect {
			correct++
			total++
		} else if !pred.WinnerCorrect {
			total++
		}
	}

	var accuracy float32
	if total > 0 {
		accuracy = float32(correct) / float32(total)
	}

	// Get rank (requires calculating full leaderboard - optimize this later)
	leaderboard, err := db.CalculateLeaderboard(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate leaderboard: %w", err)
	}

	rank := 0
	for _, entry := range leaderboard {
		if entry.UserId == user.ID {
			rank = entry.Rank
			break
		}
	}

	return &models.LeaderboardEntry{
		UserId:              user.ID,
		Username:            user.Username,
		TotalWinnersCorrect: correct,
		WinnerAccuracy:      accuracy,
		Rank:                rank,
	}, nil
}
