package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all models
type BaseModel struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

// User represents a user in the system
type User struct {
	BaseModel
	Email     string     `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Username  string     `json:"username" gorm:"uniqueIndex;not null" validate:"required,min=3,max=20"`
	Password  string     `json:"-" gorm:"not null"` // Hidden from JSON
	FirstName string     `json:"first_name" validate:"max=50"`
	LastName  string     `json:"last_name" validate:"max=50"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	IsAdmin   bool       `json:"is_admin" gorm:"default:false"`
	LastLogin *time.Time `json:"last_login,omitempty"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	if u.LastName != "" {
		return u.LastName
	}
	return u.Username
}

// IsValid checks if the user data is valid
func (u *User) IsValid() bool {
	return u.Email != "" && u.Username != "" && u.Password != ""
}
