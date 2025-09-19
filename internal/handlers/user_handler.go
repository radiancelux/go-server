package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go-server/internal/database/repositories"
	"go-server/internal/errors"
	"go-server/internal/logger"
	"go-server/internal/middleware"
)

// UserHandler handles user-related endpoints
type UserHandler struct {
	userRepo *repositories.UserRepository
	logger   logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *repositories.UserRepository, logger logger.Logger) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
		logger:   logger,
	}
}

// GetProfile returns the current user's profile
func (uh *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		errors.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", "NOT_AUTHENTICATED")
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// GetUserByID returns a user by ID (admin only)
func (uh *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Get user ID from URL path
	userIDStr := r.URL.Path[len("/api/users/"):]
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid user ID", "INVALID_USER_ID")
		return
	}

	// Get user from database
	user, err := uh.userRepo.GetUserByID(r.Context(), uint(userID))
	if err != nil {
		uh.logger.Error("Failed to get user", "user_id", userID, "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusNotFound, "User not found", "USER_NOT_FOUND")
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// ListUsers returns a list of users (admin only)
func (uh *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	// Set default values
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	// Get users from database
	users, err := uh.userRepo.ListUsers(r.Context(), offset, limit)
	if err != nil {
		uh.logger.Error("Failed to list users", "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to retrieve users", "DATABASE_ERROR")
		return
	}

	// Get total count
	total, err := uh.userRepo.CountUsers(r.Context())
	if err != nil {
		uh.logger.Error("Failed to count users", "error", err.Error())
		// Don't fail the request, just log the error
	}

	// Create response
	response := map[string]interface{}{
		"users": users,
		"pagination": map[string]interface{}{
			"offset": offset,
			"limit":  limit,
			"total":  total,
		},
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateProfile updates the current user's profile
func (uh *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get current user from context
	currentUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		errors.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", "NOT_AUTHENTICATED")
		return
	}

	// Parse request body
	var updateData struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Update user fields
	if updateData.FirstName != "" {
		currentUser.FirstName = updateData.FirstName
	}
	if updateData.LastName != "" {
		currentUser.LastName = updateData.LastName
	}
	if updateData.Email != "" && updateData.Email != currentUser.Email {
		// Check if email is already taken
		existingUser, err := uh.userRepo.GetUserByEmail(r.Context(), updateData.Email)
		if err == nil && existingUser.ID != currentUser.ID {
			errors.WriteErrorResponse(w, http.StatusConflict, "Email already taken", "EMAIL_TAKEN")
			return
		}
		currentUser.Email = updateData.Email
	}

	// Update user in database
	if err := uh.userRepo.UpdateUser(r.Context(), currentUser); err != nil {
		uh.logger.Error("Failed to update user profile", "user_id", currentUser.ID, "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusInternalServerError, "Failed to update profile", "DATABASE_ERROR")
		return
	}

	uh.logger.Info("User profile updated", "user_id", currentUser.ID)

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(currentUser)
}
