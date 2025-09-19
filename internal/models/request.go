package models

import (
	"fmt"
)

// Request represents the incoming JSON request structure
type Request struct {
	Message string `json:"message"`
	UserID  int    `json:"user_id,omitempty"`
	Action  string `json:"action,omitempty"`
}

// NewRequest creates a new Request instance
func NewRequest(message, action string, userID int) *Request {
	return &Request{
		Message: message,
		UserID:  userID,
		Action:  action,
	}
}

// GetMessage returns the message
func (r Request) GetMessage() string { return r.Message }

// GetUserID returns the user ID
func (r Request) GetUserID() int { return r.UserID }

// GetAction returns the action
func (r Request) GetAction() string { return r.Action }

// Validate validates the request
func (r Request) Validate() error {
	if r.Message == "" {
		return fmt.Errorf("message is required")
	}
	if r.Action == "" {
		return fmt.Errorf("action is required")
	}
	return nil
}
