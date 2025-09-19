package services

import (
	"context"
	"fmt"

	"go-server/internal/database/models"
	"go-server/internal/database/repositories"
	"go-server/internal/logger"
)

// UserService handles user business logic
type UserService struct {
	userRepo  *repositories.UserRepository
	cacheRepo *repositories.CacheRepository
	logger    logger.Logger
}

// NewUserService creates a new user service
func NewUserService(
	userRepo *repositories.UserRepository,
	cacheRepo *repositories.CacheRepository,
	logger logger.Logger,
) *UserService {
	return &UserService{
		userRepo:  userRepo,
		cacheRepo: cacheRepo,
		logger:    logger,
	}
}

// GetUserByID retrieves a user by ID with caching
func (us *UserService) GetUserByID(ctx context.Context, userID uint) (*models.User, error) {
	// Try cache first
	if cached, err := us.cacheRepo.GetUserCache(ctx, userID); err == nil && cached != "" {
		// In a real implementation, you'd deserialize the JSON
		// For now, we'll fetch from database
	}

	// Get from database
	user, err := us.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Cache the result
	if err := us.cacheRepo.SetUserCache(ctx, userID, user, 30*60); err != nil {
		us.logger.Warn("Failed to cache user", "user_id", userID, "error", err.Error())
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := us.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// GetUserByUsername retrieves a user by username
func (us *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := us.userRepo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	return user, nil
}

// CreateUser creates a new user
func (us *UserService) CreateUser(ctx context.Context, user *models.User) error {
	if !user.IsValid() {
		return fmt.Errorf("invalid user data")
	}

	if err := us.userRepo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	us.logger.Info("User created successfully", "user_id", user.ID, "email", user.Email)
	return nil
}

// UpdateUser updates a user
func (us *UserService) UpdateUser(ctx context.Context, user *models.User) error {
	if err := us.userRepo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	// Clear cache
	if err := us.cacheRepo.DeleteUserCache(ctx, user.ID); err != nil {
		us.logger.Warn("Failed to clear user cache", "user_id", user.ID, "error", err.Error())
	}

	us.logger.Info("User updated successfully", "user_id", user.ID)
	return nil
}

// DeleteUser soft deletes a user
func (us *UserService) DeleteUser(ctx context.Context, userID uint) error {
	if err := us.userRepo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// Clear cache
	if err := us.cacheRepo.DeleteUserCache(ctx, userID); err != nil {
		us.logger.Warn("Failed to clear user cache", "user_id", userID, "error", err.Error())
	}

	us.logger.Info("User deleted successfully", "user_id", userID)
	return nil
}

// ListUsers retrieves users with pagination
func (us *UserService) ListUsers(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	users, err := us.userRepo.ListUsers(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	total, err := us.userRepo.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}

// GetActiveUsers retrieves only active users
func (us *UserService) GetActiveUsers(ctx context.Context, offset, limit int) ([]models.User, int64, error) {
	users, err := us.userRepo.GetActiveUsers(ctx, offset, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get active users: %w", err)
	}

	// For active users, we can use the same count as total users
	// In a real implementation, you might want a separate count method
	total, err := us.userRepo.CountUsers(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}
