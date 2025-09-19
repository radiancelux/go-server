package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-server/internal/auth"
	"go-server/internal/database/models"
	"go-server/internal/errors"
	"go-server/internal/logger"
)

// AuthMiddleware handles JWT authentication
type AuthMiddleware struct {
	authService *auth.AuthService
	logger      logger.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService *auth.AuthService, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// RequireAuth middleware that requires authentication
func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		token := am.extractToken(r)
		if token == "" {
			am.logger.Error("No token provided")
			errors.WriteErrorResponse(w, http.StatusUnauthorized, "Authentication required", "NO_TOKEN")
			return
		}

		// Validate token and get user
		user, err := am.authService.ValidateToken(r.Context(), token)
		if err != nil {
			am.logger.Error("Invalid token", "error", err.Error())
			errors.WriteErrorResponse(w, http.StatusUnauthorized, "Invalid token", "INVALID_TOKEN")
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), "user", user)
		ctx = context.WithValue(ctx, "user_id", user.ID)
		ctx = context.WithValue(ctx, "is_admin", user.IsAdmin)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAdmin middleware that requires admin privileges
func (am *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return am.RequireAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user is admin
		isAdmin, ok := r.Context().Value("is_admin").(bool)
		if !ok || !isAdmin {
			am.logger.Error("Admin access required", "user_id", r.Context().Value("user_id"))
			errors.WriteErrorResponse(w, http.StatusForbidden, "Admin access required", "ADMIN_REQUIRED")
			return
		}

		next.ServeHTTP(w, r)
	}))
}

// OptionalAuth middleware that adds user info if token is present
func (am *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from Authorization header
		token := am.extractToken(r)
		if token != "" {
			// Validate token and get user
			user, err := am.authService.ValidateToken(r.Context(), token)
			if err == nil {
				// Add user to request context
				ctx := context.WithValue(r.Context(), "user", user)
				ctx = context.WithValue(ctx, "user_id", user.ID)
				ctx = context.WithValue(ctx, "is_admin", user.IsAdmin)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// extractToken extracts JWT token from Authorization header
func (am *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check for "Bearer " prefix
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return ""
	}

	return strings.TrimPrefix(authHeader, bearerPrefix)
}

// GetUserFromContext extracts user from request context
func GetUserFromContext(ctx context.Context) (*models.User, bool) {
	user, ok := ctx.Value("user").(*models.User)
	return user, ok
}

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (uint, bool) {
	userID, ok := ctx.Value("user_id").(uint)
	return userID, ok
}

// IsAdminFromContext checks if user is admin from request context
func IsAdminFromContext(ctx context.Context) bool {
	isAdmin, ok := ctx.Value("is_admin").(bool)
	return ok && isAdmin
}
