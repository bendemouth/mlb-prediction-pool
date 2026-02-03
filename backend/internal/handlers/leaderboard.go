package handlers

import (
	"net/http"
)

// HandleLeaderboard returns current standings
// GET /leaderboard
func (h *Handler) HandleLeaderboard(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	leaderboard, err := h.leaderboardService.GetLeaderboard(request.Context())

	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to get leaderboard")
		return
	}

	h.respondJson(writer, http.StatusOK, leaderboard)
}

func (h *Handler) HandleUserStats(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userId := request.URL.Query().Get("user_id")
	if userId == "" {
		h.respondError(writer, http.StatusBadRequest, "Missing user_id parameter")
		return
	}

	stats, err := h.leaderboardService.GetUserStats(request.Context(), userId)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to get user stats")
		return
	}

	h.respondJson(writer, http.StatusOK, stats)
}
