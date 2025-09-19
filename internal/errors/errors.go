package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorType represents the type of error
type ErrorType string

const (
	ErrorTypeValidation   ErrorType = "validation"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	ErrorTypeForbidden    ErrorType = "forbidden"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeInternal     ErrorType = "internal"
	ErrorTypeBadRequest   ErrorType = "bad_request"
	ErrorTypeRateLimit    ErrorType = "rate_limit"
)

// APIError represents a structured API error
type APIError struct {
	Type       ErrorType `json:"type"`
	Message    string    `json:"message"`
	Code       string    `json:"code,omitempty"`
	Details    string    `json:"details,omitempty"`
	RequestID  string    `json:"request_id,omitempty"`
	StatusCode int       `json:"-"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Type, e.Code, e.Message)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

// NewAPIError creates a new API error
func NewAPIError(errorType ErrorType, message string, statusCode int) *APIError {
	return &APIError{
		Type:       errorType,
		Message:    message,
		StatusCode: statusCode,
	}
}

// NewAPIErrorWithCode creates a new API error with a specific code
func NewAPIErrorWithCode(errorType ErrorType, code, message string, statusCode int) *APIError {
	return &APIError{
		Type:       errorType,
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// WithDetails adds details to the error
func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

// WithRequestID adds request ID to the error
func (e *APIError) WithRequestID(requestID string) *APIError {
	e.RequestID = requestID
	return e
}

// Predefined common errors
var (
	// Validation errors
	ErrInvalidRequest = NewAPIError(ErrorTypeValidation, "Invalid request", http.StatusBadRequest)
	ErrMissingField   = NewAPIErrorWithCode(ErrorTypeValidation, "MISSING_FIELD", "Required field is missing", http.StatusBadRequest)
	ErrInvalidFormat  = NewAPIErrorWithCode(ErrorTypeValidation, "INVALID_FORMAT", "Invalid data format", http.StatusBadRequest)

	// Not found errors
	ErrNotFound        = NewAPIError(ErrorTypeNotFound, "Resource not found", http.StatusNotFound)
	ErrHandlerNotFound = NewAPIErrorWithCode(ErrorTypeNotFound, "HANDLER_NOT_FOUND", "Handler not found for action", http.StatusNotFound)

	// Authentication/Authorization errors
	ErrUnauthorized = NewAPIError(ErrorTypeUnauthorized, "Unauthorized", http.StatusUnauthorized)
	ErrForbidden    = NewAPIError(ErrorTypeForbidden, "Forbidden", http.StatusForbidden)

	// Conflict errors
	ErrConflict = NewAPIError(ErrorTypeConflict, "Resource conflict", http.StatusConflict)

	// Internal errors
	ErrInternal = NewAPIError(ErrorTypeInternal, "Internal server error", http.StatusInternalServerError)
	ErrDatabase = NewAPIErrorWithCode(ErrorTypeInternal, "DATABASE_ERROR", "Database operation failed", http.StatusInternalServerError)

	// Rate limiting
	ErrRateLimit = NewAPIError(ErrorTypeRateLimit, "Rate limit exceeded", http.StatusTooManyRequests)
)

// WrapError wraps an existing error with additional context
func WrapError(err error, message string) *APIError {
	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}

	return &APIError{
		Type:       ErrorTypeInternal,
		Message:    message,
		Details:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	}
}

// WrapErrorWithType wraps an existing error with a specific error type
func WrapErrorWithType(err error, errorType ErrorType, message string, statusCode int) *APIError {
	return &APIError{
		Type:       errorType,
		Message:    message,
		Details:    err.Error(),
		StatusCode: statusCode,
	}
}

// WriteErrorResponse writes an error response to the HTTP response writer
func WriteErrorResponse(w http.ResponseWriter, statusCode int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := APIError{
		Type:       ErrorTypeInternal,
		Message:    message,
		Code:       code,
		StatusCode: statusCode,
	}

	// Set appropriate error type based on status code
	switch statusCode {
	case http.StatusBadRequest:
		errorResponse.Type = ErrorTypeBadRequest
	case http.StatusUnauthorized:
		errorResponse.Type = ErrorTypeUnauthorized
	case http.StatusForbidden:
		errorResponse.Type = ErrorTypeForbidden
	case http.StatusNotFound:
		errorResponse.Type = ErrorTypeNotFound
	case http.StatusConflict:
		errorResponse.Type = ErrorTypeConflict
	case http.StatusTooManyRequests:
		errorResponse.Type = ErrorTypeRateLimit
	}

	json.NewEncoder(w).Encode(errorResponse)
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) *APIError {
	return NewAPIErrorWithCode(ErrorTypeValidation, "VALIDATION_ERROR", message, http.StatusBadRequest).WithDetails(field)
}
