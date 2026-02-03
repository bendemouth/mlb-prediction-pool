package database

import (
	"context"
	"fmt"
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

		totalWinnersCorrect, totalScoreError, totalRunsError, winnerAccuracy := calculateWinnersAndError(predictions)

		leaderboard = append(leaderboard, models.LeaderboardEntry{
			UserId:              user.Id,
			Username:            user.Username,
			TotalWinnersCorrect: totalWinnersCorrect,
			WinnerAccuracy:      winnerAccuracy,
			TotalScoreError:     totalScoreError,
			TotalRunsError:      totalRunsError,
			Rank:                0, // Rank will be assigned later
		})
	}

	// Sort by the following criteria:
	// 1. Winners correct
	// 2. Total winners correct
	// 3. Total score error (lower is better)
	// 4. Total runs error (lower is better)
	sort.Slice(leaderboard, func(i, j int) bool {
		if leaderboard[i].WinnerAccuracy > leaderboard[j].WinnerAccuracy {
			return leaderboard[i].WinnerAccuracy > leaderboard[j].WinnerAccuracy
		}
		if leaderboard[i].TotalWinnersCorrect != leaderboard[j].TotalWinnersCorrect {
			return leaderboard[i].TotalWinnersCorrect > leaderboard[j].TotalWinnersCorrect
		}
		if leaderboard[i].TotalScoreError != leaderboard[j].TotalScoreError {
			return leaderboard[i].TotalScoreError < leaderboard[j].TotalScoreError
		}
		return leaderboard[i].TotalRunsError < leaderboard[j].TotalRunsError
	})

	// Assign ranks, handling ties
	for i := range leaderboard {
		if i == 0 {
			leaderboard[i].Rank = 1
		} else {
			if leaderboard[i].WinnerAccuracy == leaderboard[i-1].WinnerAccuracy &&
				leaderboard[i].TotalWinnersCorrect == leaderboard[i-1].TotalWinnersCorrect &&
				leaderboard[i].TotalScoreError == leaderboard[i-1].TotalScoreError &&
				leaderboard[i].TotalRunsError == leaderboard[i-1].TotalRunsError {
				leaderboard[i].Rank = leaderboard[i-1].Rank
			} else {
				leaderboard[i].Rank = i + 1
			}
		}
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

	totalWinnersCorrect, totalScoreError, totalRunsError, winnerAccuracy := calculateWinnersAndError(predictions)

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
		TotalScoreError:     totalScoreError,
		TotalRunsError:      totalRunsError,
		Rank:                rank,
	}, nil
}

func calculateWinnersAndError(predictions []models.Prediction) (totalWinnersCorrect int, totalScoreError float32, totalRunsError float32, winnerAccuracy float32) {
	var totalPredictions int

	for _, pred := range predictions {
		if pred.WinnerCorrect != nil {
			totalPredictions++

			if *pred.WinnerCorrect {
				totalWinnersCorrect++
			}

			totalScoreError += pred.HomeScoreError + pred.AwayScoreError
			totalRunsError += pred.TotalScoreError
		}
	}

	if totalPredictions > 0 {
		winnerAccuracy = float32(totalWinnersCorrect) / float32(totalPredictions)
	}
	return totalWinnersCorrect, totalScoreError, totalRunsError, winnerAccuracy
}
