package handlers

import (
	"net/http"
	"time"

	"github.com/bendemouth/mlb-prediction-pool/internal/models"
	"github.com/google/uuid"
)

// HandleCreateUser creates a new user
// POST /users
func (h *Handler) HandleCreateUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		h.respondError(writer, http.StatusMethodNotAllowed, "Method not allowed")
	}

	newUuid, err := uuid.NewUUID()
	if err != nil {
		h.respondError(writer, http.StatusInternalServerError, "Failed to generate user ID")
	}

	newUserGuid := newUuid.String()

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
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
// GET /users
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
