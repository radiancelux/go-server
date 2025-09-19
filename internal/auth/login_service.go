package auth

import (
	"context"
	"fmt"
	"time"

	"go-server/internal/database/models"
	"go-server/internal/database/repositories"

	"golang.org/x/crypto/bcrypt"
)

// LoginService handles login operations
type LoginService struct {
	userRepo    *repositories.UserRepository
	cacheRepo   *repositories.CacheRepository
	jwtManager  *JWTManager
	sessionRepo *repositories.SessionRepository
}

// NewLoginService creates a new login service
func NewLoginService(
	userRepo *repositories.UserRepository,
	cacheRepo *repositories.CacheRepository,
	sessionRepo *repositories.SessionRepository,
	jwtManager *JWTManager,
) *LoginService {
	return &LoginService{
		userRepo:    userRepo,
		cacheRepo:   cacheRepo,
		sessionRepo: sessionRepo,
		jwtManager:  jwtManager,
	}
}

// Login authenticates a user and returns an auth response
func (ls *LoginService) Login(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	// Get user by email
	user, err := ls.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	// Verify password
	if err := ls.verifyPassword(req.Password, user.Password); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Generate JWT token
	token, err := ls.jwtManager.GenerateToken(user.ID, user.Username, user.Email, user.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Generate session token
	sessionToken, err := ls.generateSessionToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Create session
	session := &models.Session{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour session
		IPAddress: ipAddress,
		UserAgent: userAgent,
		IsActive:  true,
	}

	if err := ls.sessionRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLogin = &now
	if err := ls.userRepo.UpdateUser(ctx, user); err != nil {
		// Log error but don't fail login
		fmt.Printf("Warning: failed to update last login: %v\n", err)
	}

	// Cache user session
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	if err := ls.cacheRepo.Set(ctx, cacheKey, user, 30*time.Minute); err != nil {
		// Log error but don't fail login
		fmt.Printf("Warning: failed to cache user: %v\n", err)
	}

	// Get token expiration
	claims, _ := ls.jwtManager.ValidateToken(token)

	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: claims.ExpiresAt.Time,
		SessionID: sessionToken,
	}, nil
}

// verifyPassword verifies a password against a hash
func (ls *LoginService) verifyPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// generateSessionToken generates a random session token
func (ls *LoginService) generateSessionToken() (string, error) {
	// This would generate a secure random token
	// For now, return a simple implementation
	return fmt.Sprintf("session_%d", time.Now().UnixNano()), nil
}
