package models

type Prediction struct {
	Id                  int     `json:"id"`
	UserId              int     `json:"user_id"`
	GameId              string  `json:"game_id"`
	HomeScorePredicted  float32 `json:"home_score_predicted"`
	AwayScorePredicted  float32 `json:"away_score_predicted"`
	TotalScorePredicted float32 `json:"total_score_predicted"`
	Confidence          float32 `json:"confidence"`
	PredictedWinnerId   int     `json:"predicted_winner_id"`
	ActualWinnerId      int     `json:"actual_winner_id"`
	WinnerCorrect       bool    `json:"winner_correct"`
	HomeScoreError      float32 `json:"home_score_error"`
	AwayScoreError      float32 `json:"away_score_error"`
	TotalScoreError     float32 `json:"total_score_error"`
}
