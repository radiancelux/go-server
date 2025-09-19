// Package interfaces defines the core contracts for the Go server.
package interfaces

import "time"

// APIRequest defines the contract for incoming API requests.
type APIRequest interface {
	GetMessage() string
	GetUserID() int
	GetAction() string
	Validate() error
}

// APIResponse defines the contract for outgoing API responses.
type APIResponse interface {
	GetStatus() string
	GetMessage() string
	GetTimestamp() time.Time
	GetData() any
	ToJSON() ([]byte, error)
}

// Handler defines the contract for API request handlers.
type Handler interface {
	Handle(req APIRequest) (APIResponse, error)
	GetAction() string
}

// Logger defines the contract for logging throughout the application.
type Logger interface {
	Info(msg string, args ...any)
	Error(msg string, args ...any)
	Debug(msg string, args ...any)
}
