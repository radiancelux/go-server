package auth

import (
	"context"

	"go-server/internal/database/models"
	"go-server/internal/database/repositories"
)

// AuthService handles authentication operations
type AuthService struct {
	loginService      *LoginService
	registrationService *RegistrationService
	sessionService    *SessionService
}

// NewAuthService creates a new authentication service
func NewAuthService(
	userRepo *repositories.UserRepository,
	cacheRepo *repositories.CacheRepository,
	sessionRepo *repositories.SessionRepository,
	jwtManager *JWTManager,
) *AuthService {
	return &AuthService{
		loginService: NewLoginService(userRepo, cacheRepo, sessionRepo, jwtManager),
		registrationService: NewRegistrationService(userRepo, cacheRepo, jwtManager),
		sessionService: NewSessionService(userRepo, cacheRepo, sessionRepo, jwtManager),
	}
}

// Login authenticates a user and returns an auth response
func (as *AuthService) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	return as.loginService.Login(ctx, req, ipAddress, userAgent)
}

// Register creates a new user account
func (as *AuthService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	return as.registrationService.Register(ctx, req)
}

// Logout invalidates a user session
func (as *AuthService) Logout(ctx context.Context, userID uint, sessionID string) error {
	return as.sessionService.Logout(ctx, userID, sessionID)
}

// ValidateToken validates a JWT token and returns the user
func (as *AuthService) ValidateToken(ctx context.Context, tokenString string) (*models.User, error) {
	return as.sessionService.ValidateToken(ctx, tokenString)
}

// RefreshToken refreshes a JWT token
func (as *AuthService) RefreshToken(ctx context.Context, tokenString string) (*AuthResponse, error) {
	return as.sessionService.RefreshToken(ctx, tokenString)
}

// CleanupExpiredSessions removes expired sessions
func (as *AuthService) CleanupExpiredSessions(ctx context.Context) error {
	return as.sessionService.CleanupExpiredSessions(ctx)
}

// GetUserSessions retrieves all active sessions for a user
func (as *AuthService) GetUserSessions(ctx context.Context, userID uint) ([]models.Session, error) {
	return as.sessionService.GetUserSessions(ctx, userID)
}

// DeleteAllUserSessions deletes all sessions for a user
func (as *AuthService) DeleteAllUserSessions(ctx context.Context, userID uint) error {
	return as.sessionService.DeleteAllUserSessions(ctx, userID)
}