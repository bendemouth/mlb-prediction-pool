package handlers

import (
	"net/http"
)

// GetLeaderboard returns current standings
// GET /leaderboard
func (h *Handler) GetLeaderboard(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	leaderboard, err := h.db.CalculateLeaderboard(request.Context())

	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to get leaderboard")
		return
	}

	h.respondJson(writer, http.StatusOK, leaderboard)
}
