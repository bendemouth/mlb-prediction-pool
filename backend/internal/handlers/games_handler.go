package handlers

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/bendemouth/mlb-prediction-pool/internal/models"
)

// GamePredictionSummary combines a game with aggregated prediction stats
type GamePredictionSummary struct {
	models.Game
	PredictionCount        int     `json:"prediction_count"`
	AvgHomeScorePredicted  float64 `json:"avg_home_score_predicted"`
	AvgAwayScorePredicted  float64 `json:"avg_away_score_predicted"`
	AvgTotalScorePredicted float64 `json:"avg_total_score_predicted"`
	AvgConfidence          float64 `json:"avg_confidence"`
}

// GET /games/upcoming
// Returns upcoming games with aggregated community prediction stats
func (h *Handler) GetUpcomingGamesSummary(writer http.ResponseWriter, request *http.Request) {
	games, err := h.db.GetUpcomingGames(request.Context())
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, fmt.Sprint("Failed to get upcoming games: ", err))
		return
	}

	summaries := make([]GamePredictionSummary, 0, len(games))

	for _, game := range games {
		predictions, err := h.db.GetPredictionsByGame(request.Context(), game.GameId)
		if err != nil {
			h.respondError(writer, http.StatusInternalServerError, fmt.Sprintf("Failed to get predictions for game %s: %s", game.GameId, err))
			return
		}

		summary := GamePredictionSummary{
			Game:            game,
			PredictionCount: len(predictions),
		}

		if len(predictions) > 0 {
			var totalHome, totalAway, totalRuns, totalConf float64
			for _, p := range predictions {
				totalHome += float64(p.HomeScorePredicted)
				totalAway += float64(p.AwayScorePredicted)
				totalRuns += float64(p.TotalScorePredicted)
				totalConf += float64(p.Confidence)
			}
			n := float64(len(predictions))
			summary.AvgHomeScorePredicted = totalHome / n
			summary.AvgAwayScorePredicted = totalAway / n
			summary.AvgTotalScorePredicted = totalRuns / n
			summary.AvgConfidence = totalConf / n
		}

		summaries = append(summaries, summary)
	}

	// Sort by game date ascending
	sort.Slice(summaries, func(i, j int) bool {
		return summaries[i].Date.Before(summaries[j].Date)
	})

	h.respondJson(writer, http.StatusOK, summaries)
}
