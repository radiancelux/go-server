package auth

import (
	"context"
	"fmt"

	"go-server/internal/database/models"
	"go-server/internal/database/repositories"
)

// SessionService handles session management operations
type SessionService struct {
	userRepo    *repositories.UserRepository
	cacheRepo   *repositories.CacheRepository
	sessionRepo *repositories.SessionRepository
	jwtManager  *JWTManager
}

// NewSessionService creates a new session service
func NewSessionService(
	userRepo *repositories.UserRepository,
	cacheRepo *repositories.CacheRepository,
	sessionRepo *repositories.SessionRepository,
	jwtManager *JWTManager,
) *SessionService {
	return &SessionService{
		userRepo:    userRepo,
		cacheRepo:   cacheRepo,
		sessionRepo: sessionRepo,
		jwtManager:  jwtManager,
	}
}

// Logout invalidates a user session
func (ss *SessionService) Logout(ctx context.Context, userID uint, sessionID string) error {
	// Delete session from database
	if err := ss.sessionRepo.DeleteSession(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	// Delete session from cache
	if err := ss.cacheRepo.DeleteUserSession(ctx, userID, sessionID); err != nil {
		// Log error but don't fail logout
		fmt.Printf("Warning: failed to delete session from cache: %v\n", err)
	}

	return nil
}

// ValidateToken validates a JWT token and returns the user
func (ss *SessionService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Validate JWT token
	claims, err := ss.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Get user from database
	user, err := ss.userRepo.GetUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Check if user is still active
	if !user.IsActive {
		return nil, fmt.Errorf("user account is deactivated")
	}

	return user, nil
}

// RefreshToken refreshes a JWT token
func (ss *SessionService) RefreshToken(ctx context.Context, tokenString string) (*AuthResponse, error) {
	// Validate current token
	user, err := ss.ValidateToken(ctx, tokenString)
	if err != nil {
		return nil, err
	}

	// Generate new token
	newToken, err := ss.jwtManager.GenerateToken(user.ID, user.Username, user.Email, user.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate new token: %w", err)
	}

	// Get new token expiration
	claims, _ := ss.jwtManager.ValidateToken(newToken)

	return &AuthResponse{
		Token:     newToken,
		User:      user,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// CleanupExpiredSessions removes expired sessions
func (ss *SessionService) CleanupExpiredSessions(ctx context.Context) error {
	return ss.sessionRepo.CleanupExpiredSessions(ctx)
}

// GetUserSessions retrieves all active sessions for a user
func (ss *SessionService) GetUserSessions(ctx context.Context, userID uint) ([]models.Session, error) {
	return ss.sessionRepo.GetSessionsByUser(ctx, userID)
}

// DeleteAllUserSessions deletes all sessions for a user
func (ss *SessionService) DeleteAllUserSessions(ctx context.Context, userID uint) error {
	return ss.sessionRepo.DeleteUserSessions(ctx, userID)
}
