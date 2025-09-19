package auth

import (
	"context"
	"fmt"

	"go-server/internal/database/models"
	"go-server/internal/database/repositories"

	"golang.org/x/crypto/bcrypt"
)

// RegistrationService handles user registration operations
type RegistrationService struct {
	userRepo    *repositories.UserRepository
	cacheRepo   *repositories.CacheRepository
	jwtManager  *JWTManager
}

// NewRegistrationService creates a new registration service
func NewRegistrationService(
	userRepo *repositories.UserRepository,
	cacheRepo *repositories.CacheRepository,
	jwtManager *JWTManager,
) *RegistrationService {
	return &RegistrationService{
		userRepo:   userRepo,
		cacheRepo:  cacheRepo,
		jwtManager: jwtManager,
	}
}

// Register creates a new user account
func (rs *RegistrationService) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
	// Check if email already exists
	existingUser, _ := rs.userRepo.GetUserByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("email already registered")
	}

	// Check if username already exists
	existingUser, _ = rs.userRepo.GetUserByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, fmt.Errorf("username already taken")
	}

	// Hash password
	hashedPassword, err := rs.hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &models.User{
		Email:     req.Email,
		Username:  req.Username,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		IsAdmin:   false,
	}

	if err := rs.userRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := rs.jwtManager.GenerateToken(user.ID, user.Username, user.Email, user.IsAdmin)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Get token expiration
	claims, _ := rs.jwtManager.ValidateToken(token)

	return &AuthResponse{
		Token:     token,
		User:      user,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// hashPassword hashes a password using bcrypt
func (rs *RegistrationService) hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
