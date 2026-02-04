package handlers

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/bendemouth/mlb-prediction-pool/internal/models"
	"github.com/google/uuid"
)

// HandleCreateUser creates a new user
// POST /users/create
func (h *Handler) HandleCreateUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	newUuid, err := uuid.NewUUID()
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to generate user ID")
		return
	}

	newUserGuid := newUuid.String()

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := h.decodeJsonBody(request, &req); err != nil {
		h.respondError(writer, http.StatusBadRequest, "Invalid request body")
		return
	}
	if !isValidEmail(req.Email) {
		h.respondError(writer, http.StatusBadRequest, "Invalid email address")
		return
	}
	if req.Username == "" {
		h.respondError(writer, http.StatusBadRequest, "Username is required")
		return
	}

	newUserRequest := &models.User{
		Id:        newUserGuid,
		Username:  req.Username,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	err = h.db.CreateUser(request.Context(), newUserRequest)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to create user")
		return
	}

	h.respondJson(writer, http.StatusCreated, newUserRequest)
}

// HandleGetUser retrieves a user by userId
// GET /users?user_id=username
func (h *Handler) HandleGetUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userId := request.URL.Query().Get("user_id")
	if userId == "" {
		h.respondError(writer, http.StatusBadRequest, "Missing user_id parameter")
		return
	}

	user, err := h.db.GetUser(request.Context(), userId)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to get user")
		return
	}
	h.respondJson(writer, http.StatusOK, user)
}

// HandleListUsers retrieves all users
// GET /users/listUsers
func (h *Handler) HandleListUsers(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	users, err := h.db.ListUsers(request.Context())
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to list users")
		return
	}
	h.respondJson(writer, http.StatusOK, users)
}

// HandleGetUserStats retrieves statistics for a specific user
// GET /users/stats?user_id=username
func (h *Handler) HandleGetUserStats(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userId := request.URL.Query().Get("user_id")
	if userId == "" {
		h.respondError(writer, http.StatusBadRequest, "Missing user_id parameter")
		return
	}

	stats, err := h.db.GetUserStats(request.Context(), userId)
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to get user stats")
		return
	}

	h.respondJson(writer, http.StatusOK, stats)
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
