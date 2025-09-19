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

// APIRequest represents the incoming JSON request structure for API endpoints
type APIRequest struct {
	Action string                 `json:"action"`
	Data   map[string]interface{} `json:"data,omitempty"`
}

// NewAPIRequest creates a new APIRequest instance
func NewAPIRequest(action string, data map[string]interface{}) *APIRequest {
	return &APIRequest{
		Action: action,
		Data:   data,
	}
}

// GetMessage returns the message (implements APIRequest interface)
func (r APIRequest) GetMessage() string {
	if msg, ok := r.Data["message"].(string); ok {
		return msg
	}
	return ""
}

// GetUserID returns the user ID (implements APIRequest interface)
func (r APIRequest) GetUserID() int {
	if userID, ok := r.Data["user_id"].(float64); ok {
		return int(userID)
	}
	return 0
}

// GetAction returns the action (implements APIRequest interface)
func (r APIRequest) GetAction() string {
	return r.Action
}

// Validate validates the API request (implements APIRequest interface)
func (r APIRequest) Validate() error {
	if r.Action == "" {
		return fmt.Errorf("action is required")
	}
	return nil
}
