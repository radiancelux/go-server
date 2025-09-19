package repositories

import (
	"context"

	"go-server/internal/database/models"
	"gorm.io/gorm"
)

// UserRepository handles user-related database operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser creates a new user
func (ur *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	return ur.db.WithContext(ctx).Create(user).Error
}

// GetUserByID retrieves a user by ID
func (ur *UserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := ur.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email
func (ur *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := ur.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by username
func (ur *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := ur.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates a user
func (ur *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return ur.db.WithContext(ctx).Save(user).Error
}

// DeleteUser soft deletes a user
func (ur *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	return ur.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

// ListUsers retrieves users with pagination
func (ur *UserRepository) ListUsers(ctx context.Context, offset, limit int) ([]models.User, error) {
	var users []models.User
	err := ur.db.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	return users, err
}

// CountUsers returns the total number of users
func (ur *UserRepository) CountUsers(ctx context.Context) (int64, error) {
	var count int64
	err := ur.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return count, err
}

// GetActiveUsers retrieves only active users
func (ur *UserRepository) GetActiveUsers(ctx context.Context, offset, limit int) ([]models.User, error) {
	var users []models.User
	err := ur.db.WithContext(ctx).
		Where("is_active = ?", true).
		Offset(offset).
		Limit(limit).
		Find(&users).Error
	return users, err
}
