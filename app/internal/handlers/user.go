package handlers

import (
	"encoding/json"
	"net/http"
	"net/mail"
	"strings"

	"github.com/jasonsavage/financeflow/internal/middleware"
	"github.com/jasonsavage/financeflow/internal/repository"
)

type UserHandler struct {
	Repo repository.DatabaseRepo
}

func NewUserHandler(repo repository.DatabaseRepo) *UserHandler {
	return &UserHandler{Repo: repo}
}

type userDetailsResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type updateUserDetailsRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// GET /user/details
func (h *UserHandler) GetDetails(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.Repo.GetUserDetails(r.Context(), userID)
	if err != nil {
		http.Error(w, "Failed to fetch user details", http.StatusInternalServerError)
		return
	}

	resp := userDetailsResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}

	writeJSON(w, http.StatusOK, resp)
}

// PUT /user/details
func (h *UserHandler) UpdateDetails(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req updateUserDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate email format
	if req.Email != "" {
		_, err := mail.ParseAddress(req.Email)
		if err != nil {
			http.Error(w, "Validation error: invalid email format", http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, "Validation error: email is required", http.StatusBadRequest)
		return
	}

	err := h.Repo.UpdateUserDetails(r.Context(), userID, req.FirstName, req.LastName, req.Email)
	if err != nil {
		// PostgreSQL unique violation error detection
		if strings.Contains(err.Error(), "users_email_key") || strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "23505") {
			http.Error(w, "Validation error: email already in use", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to update user details", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "User details successfully updated"})
}
