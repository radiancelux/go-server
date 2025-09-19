package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"go-server/internal/auth"
	"go-server/internal/errors"
	"go-server/internal/logger"
	"go-server/internal/models"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *auth.AuthService
	logger      logger.Logger
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(authService *auth.AuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

// Login handles user login
func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req auth.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.logger.Error("Invalid login request", "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Validate request
	if err := validateLoginRequest(&req); err != nil {
		ah.logger.Error("Login validation failed", "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	// Get client info
	ipAddress := getClientIP(r)
	userAgent := r.Header.Get("User-Agent")

	// Attempt login
	response, err := ah.authService.Login(r.Context(), &req, ipAddress, userAgent)
	if err != nil {
		ah.logger.Error("Login failed", "email", req.Email, "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid credentials", "LOGIN_FAILED")
		return
	}

	ah.logger.Info("User logged in successfully", "user_id", response.User.ID, "email", response.User.Email)

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Register handles user registration
func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req auth.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		ah.logger.Error("Invalid registration request", "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Invalid request body", "INVALID_REQUEST")
		return
	}

	// Validate request
	if err := validateRegisterRequest(&req); err != nil {
		ah.logger.Error("Registration validation failed", "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	// Attempt registration
	response, err := ah.authService.Register(r.Context(), &req)
	if err != nil {
		ah.logger.Error("Registration failed", "email", req.Email, "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusConflict, err.Error(), "REGISTRATION_FAILED")
		return
	}

	ah.logger.Info("User registered successfully", "user_id", response.User.ID, "email", response.User.Email)

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout
func (ah *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*auth.AuthResponse)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", "NOT_AUTHENTICATED")
		return
	}

	// Get session ID from request or context
	sessionID := r.Header.Get("X-Session-ID")
	if sessionID == "" {
		// Try to get from context if available
		if session, ok := r.Context().Value("session_id").(string); ok {
			sessionID = session
		}
	}

	if sessionID != "" {
		// Logout with session
		if err := ah.authService.Logout(r.Context(), user.User.ID, sessionID); err != nil {
			ah.logger.Error("Logout failed", "user_id", user.User.ID, "error", err.Error())
			// Don't fail logout if session cleanup fails
		}
	}

	ah.logger.Info("User logged out successfully", "user_id", user.User.ID)

	// Write success response
	response := models.NewSuccessResponse("Logged out successfully", nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RefreshToken handles token refresh
func (ah *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Get current token from Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		errors.WriteErrorResponse(w, http.StatusBadRequest, "Authorization header required", "NO_AUTH_HEADER")
		return
	}

	// Extract token (assuming "Bearer " prefix)
	token := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	}

	// Refresh token
	response, err := ah.authService.RefreshToken(r.Context(), token)
	if err != nil {
		ah.logger.Error("Token refresh failed", "error", err.Error())
		errors.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token", "REFRESH_FAILED")
		return
	}

	ah.logger.Info("Token refreshed successfully", "user_id", response.User.ID)

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetProfile returns the current user's profile
func (ah *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by auth middleware)
	user, ok := r.Context().Value("user").(*auth.AuthResponse)
	if !ok {
		errors.WriteErrorResponse(w, http.StatusUnauthorized, "User not authenticated", "NOT_AUTHENTICATED")
		return
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user.User)
}

// Validation functions
func validateLoginRequest(req *auth.LoginRequest) error {
	if req.Email == "" {
		return errors.NewValidationError("email", "Email is required")
	}
	if req.Password == "" {
		return errors.NewValidationError("password", "Password is required")
	}
	if len(req.Password) < 6 {
		return errors.NewValidationError("password", "Password must be at least 6 characters")
	}
	return nil
}

func validateRegisterRequest(req *auth.RegisterRequest) error {
	if req.Email == "" {
		return errors.NewValidationError("email", "Email is required")
	}
	if req.Username == "" {
		return errors.NewValidationError("username", "Username is required")
	}
	if len(req.Username) < 3 {
		return errors.NewValidationError("username", "Username must be at least 3 characters")
	}
	if len(req.Username) > 20 {
		return errors.NewValidationError("username", "Username must be at most 20 characters")
	}
	if req.Password == "" {
		return errors.NewValidationError("password", "Password is required")
	}
	if len(req.Password) < 6 {
		return errors.NewValidationError("password", "Password must be at least 6 characters")
	}
	return nil
}

// Helper function to get client IP
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if idx := strings.Index(xff, ","); idx != -1 {
			return strings.TrimSpace(xff[:idx])
		}
		return strings.TrimSpace(xff)
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
