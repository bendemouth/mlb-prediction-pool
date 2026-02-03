package models

import "time"

type LeaderboardEntry struct {
	UserId              string    `json:"user_id" dynamodbav:"userId"`
	Username            string    `json:"username" dynamodbav:"username"`
	TotalWinnersCorrect int       `json:"total_winners_correct" dynamodbav:"totalWinnersCorrect"`
	WinnerAccuracy      float32   `json:"winner_accuracy" dynamodbav:"winnerAccuracy"`
	TotalScoreError     float32   `json:"total_score_error" dynamodbav:"totalScoreError"`
	TotalRunsError      float32   `json:"total_runs_error" dynamodbav:"totalRunsError"`
	Rank                int       `json:"rank" dynamodbav:"rank"`
	UpdatedAt           time.Time `json:"updated_at" dynamodbav:"updatedAt"`
	TTL                 int64     `json:"ttl" dynamodbav:"ttl"` // Unix timestamp for auto-deletion
}
