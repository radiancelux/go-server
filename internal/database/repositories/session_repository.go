package repositories

import (
	"context"
	"time"

	"go-server/internal/database/models"
	"gorm.io/gorm"
)

// SessionRepository handles session-related database operations
type SessionRepository struct {
	db *gorm.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// CreateSession creates a new session
func (sr *SessionRepository) CreateSession(ctx context.Context, session *models.Session) error {
	return sr.db.WithContext(ctx).Create(session).Error
}

// GetSessionByToken retrieves a session by token
func (sr *SessionRepository) GetSessionByToken(ctx context.Context, token string) (*models.Session, error) {
	var session models.Session
	err := sr.db.WithContext(ctx).
		Where("token = ? AND is_active = ? AND expires_at > ?", token, true, time.Now()).
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSessionsByUser retrieves all sessions for a user
func (sr *SessionRepository) GetSessionsByUser(ctx context.Context, userID uint) ([]models.Session, error) {
	var sessions []models.Session
	err := sr.db.WithContext(ctx).
		Where("user_id = ? AND is_active = ?", userID, true).
		Order("created_at DESC").
		Find(&sessions).Error
	return sessions, err
}

// DeleteSession deletes a session
func (sr *SessionRepository) DeleteSession(ctx context.Context, userID uint, sessionID string) error {
	return sr.db.WithContext(ctx).
		Where("user_id = ? AND token = ?", userID, sessionID).
		Delete(&models.Session{}).Error
}

// DeleteUserSessions deletes all sessions for a user
func (sr *SessionRepository) DeleteUserSessions(ctx context.Context, userID uint) error {
	return sr.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Delete(&models.Session{}).Error
}

// CleanupExpiredSessions removes expired sessions
func (sr *SessionRepository) CleanupExpiredSessions(ctx context.Context) error {
	return sr.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.Session{}).Error
}

// UpdateSessionLastActivity updates the last activity time for a session
func (sr *SessionRepository) UpdateSessionLastActivity(ctx context.Context, sessionID string) error {
	return sr.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("token = ?", sessionID).
		Update("updated_at", time.Now()).Error
}

// CountActiveSessions returns the number of active sessions for a user
func (sr *SessionRepository) CountActiveSessions(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := sr.db.WithContext(ctx).
		Model(&models.Session{}).
		Where("user_id = ? AND is_active = ? AND expires_at > ?", userID, true, time.Now()).
		Count(&count).Error
	return count, err
}
