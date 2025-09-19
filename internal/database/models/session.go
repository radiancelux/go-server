package models

import (
	"time"
)

// Session represents a user session
type Session struct {
	BaseModel
	UserID    uint      `json:"user_id" gorm:"not null"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null" validate:"required"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
	IsActive  bool      `json:"is_active" gorm:"default:true"`
}

// TableName returns the table name for Session
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the session is valid (active and not expired)
func (s *Session) IsValid() bool {
	return s.IsActive && !s.IsExpired()
}

// Extend extends the session expiration time
func (s *Session) Extend(duration time.Duration) {
	s.ExpiresAt = time.Now().Add(duration)
}

// Deactivate marks the session as inactive
func (s *Session) Deactivate() {
	s.IsActive = false
}
