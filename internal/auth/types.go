package auth

import (
	"time"

	"go-server/internal/database/models"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required,min=3,max=20"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"max=50"`
	LastName  string `json:"last_name" validate:"max=50"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token     string      `json:"token"`
	User      *models.User `json:"user"`
	ExpiresAt time.Time   `json:"expires_at"`
	SessionID string      `json:"session_id,omitempty"`
}

// TokenRefreshRequest represents a token refresh request
type TokenRefreshRequest struct {
	Token string `json:"token" validate:"required"`
}

// PasswordChangeRequest represents a password change request
type PasswordChangeRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}

// ProfileUpdateRequest represents a profile update request
type ProfileUpdateRequest struct {
	FirstName string `json:"first_name" validate:"max=50"`
	LastName  string `json:"last_name" validate:"max=50"`
	Email     string `json:"email" validate:"email"`
}
