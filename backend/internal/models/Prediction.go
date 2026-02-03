package models

import "time"

type Prediction struct {
	UserId              string    `json:"user_id"`
	GameId              string    `json:"game_id"`
	HomeScorePredicted  float32   `json:"home_score_predicted"`
	AwayScorePredicted  float32   `json:"away_score_predicted"`
	TotalScorePredicted float32   `json:"total_score_predicted"`
	Confidence          float32   `json:"confidence"`
	PredictedWinnerId   string    `json:"predicted_winner_id"`
	ActualWinnerId      string    `json:"actual_winner_id,omitempty"`
	WinnerCorrect       *bool     `json:"winner_correct,omitempty"`
	HomeScoreError      float32   `json:"home_score_error,omitempty"`
	AwayScoreError      float32   `json:"away_score_error,omitempty"`
	TotalScoreError     float32   `json:"total_score_error,omitempty"`
	SubmittedAt         time.Time `json:"submitted_at"`
}
