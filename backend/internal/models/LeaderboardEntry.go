package models

type LeaderboardEntry struct {
	UserId              int     `json:"user_id"`
	Username            string  `json:"username"`
	TotalWinnersCorrect int     `json:"total_winners_correct"`
	WinnerAccuracy      float32 `json:"winner_accuracy"`
	TotalScoreError     float32 `json:"total_score_error"`
	TotalRunsError      float32 `json:"total_runs_error"`
	Rank                int     `json:"rank"`
}
