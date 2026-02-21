package models

import "time"

type Prediction struct {
	UserId              string    `json:"user_id"              dynamodbav:"userId"`
	GameId              string    `json:"game_id"              dynamodbav:"gameId"`
	HomeScorePredicted  float32   `json:"home_score_predicted"  dynamodbav:"homeScorePredicted"`
	AwayScorePredicted  float32   `json:"away_score_predicted"  dynamodbav:"awayScorePredicted"`
	TotalScorePredicted float32   `json:"total_score_predicted" dynamodbav:"totalScorePredicted"`
	Confidence          float32   `json:"confidence"            dynamodbav:"confidence"`
	PredictedWinnerId   string    `json:"predicted_winner_id"   dynamodbav:"predictedWinnerId"`
	ActualWinnerId      string    `json:"actual_winner_id,omitempty" dynamodbav:"actualWinnerId,omitempty"`
	WinnerCorrect       *bool     `json:"winner_correct,omitempty"   dynamodbav:"winnerCorrect,omitempty"`
	HomeScoreError      float32   `json:"home_score_error,omitempty" dynamodbav:"homeScoreError,omitempty"`
	AwayScoreError      float32   `json:"away_score_error,omitempty" dynamodbav:"awayScoreError,omitempty"`
	TotalScoreError     float32   `json:"total_score_error,omitempty" dynamodbav:"totalScoreError,omitempty"`
	SubmittedAt         time.Time `json:"submitted_at"          dynamodbav:"submittedAt"`
}
