// Package models contains the data structures for requests and responses.
// These models implement the interfaces defined in the interfaces package
// and provide the concrete types used throughout the application.
package models

import (
	"encoding/json"
	"time"
)

// Response represents the JSON response structure.
// It implements the APIResponse interface and provides consistent
// response formatting across all endpoints.
type Response struct {
	Status    string    `json:"status"`              // Response status (success/error)
	Message   string    `json:"message"`             // Response message
	Timestamp time.Time `json:"timestamp"`           // When the response was created
	Data      any       `json:"data,omitempty"`      // Optional response data
}

// NewResponse creates a new Response instance with the provided values.
func NewResponse(status, message string, data any) *Response {
	return &Response{
		Status:    status,
		Message:   message,
		Timestamp: time.Now(),
		Data:      data,
	}
}

// NewSuccessResponse creates a success response with the provided message and data.
func NewSuccessResponse(message string, data any) *Response {
	return NewResponse("success", message, data)
}

// NewErrorResponse creates an error response with the provided message.
func NewErrorResponse(message string) *Response {
	return NewResponse("error", message, nil)
}

// GetStatus returns the response status.
func (r Response) GetStatus() string { return r.Status }

// GetMessage returns the response message.
func (r Response) GetMessage() string { return r.Message }

// GetTimestamp returns when the response was created.
func (r Response) GetTimestamp() time.Time { return r.Timestamp }

// GetData returns the optional response data.
func (r Response) GetData() any { return r.Data }

// ToJSON serializes the response to JSON bytes.
func (r Response) ToJSON() ([]byte, error) {
	return json.Marshal(r)
}
