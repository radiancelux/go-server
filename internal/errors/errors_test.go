package errors

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewAPIError(t *testing.T) {
	err := NewAPIError(ErrorTypeValidation, "Test error", http.StatusBadRequest)
	
	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected type %s, got %s", ErrorTypeValidation, err.Type)
	}
	
	if err.Message != "Test error" {
		t.Errorf("Expected message 'Test error', got %s", err.Message)
	}
	
	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, err.StatusCode)
	}
}

func TestNewAPIErrorWithCode(t *testing.T) {
	err := NewAPIErrorWithCode(ErrorTypeValidation, "INVALID_FIELD", "Field is invalid", http.StatusBadRequest)
	
	if err.Type != ErrorTypeValidation {
		t.Errorf("Expected type %s, got %s", ErrorTypeValidation, err.Type)
	}
	
	if err.Code != "INVALID_FIELD" {
		t.Errorf("Expected code 'INVALID_FIELD', got %s", err.Code)
	}
	
	if err.Message != "Field is invalid" {
		t.Errorf("Expected message 'Field is invalid', got %s", err.Message)
	}
}

func TestWithDetails(t *testing.T) {
	err := NewAPIError(ErrorTypeValidation, "Test error", http.StatusBadRequest)
	err = err.WithDetails("Additional details")
	
	if err.Details != "Additional details" {
		t.Errorf("Expected details 'Additional details', got %s", err.Details)
	}
}

func TestWithRequestID(t *testing.T) {
	err := NewAPIError(ErrorTypeValidation, "Test error", http.StatusBadRequest)
	err = err.WithRequestID("req-123")
	
	if err.RequestID != "req-123" {
		t.Errorf("Expected request ID 'req-123', got %s", err.RequestID)
	}
}

func TestError(t *testing.T) {
	err := NewAPIErrorWithCode(ErrorTypeValidation, "TEST_ERROR", "Test error", http.StatusBadRequest)
	
	expected := "[validation] TEST_ERROR: Test error"
	if err.Error() != expected {
		t.Errorf("Expected error string '%s', got '%s'", expected, err.Error())
	}
}

func TestWrapError(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapError(originalErr, "wrapped message")
	
	if wrappedErr.Type != ErrorTypeInternal {
		t.Errorf("Expected type %s, got %s", ErrorTypeInternal, wrappedErr.Type)
	}
	
	if wrappedErr.Message != "wrapped message" {
		t.Errorf("Expected message 'wrapped message', got %s", wrappedErr.Message)
	}
	
	if wrappedErr.Details != "original error" {
		t.Errorf("Expected details 'original error', got %s", wrappedErr.Details)
	}
}

func TestWrapErrorWithType(t *testing.T) {
	originalErr := errors.New("original error")
	wrappedErr := WrapErrorWithType(originalErr, ErrorTypeValidation, "validation failed", http.StatusBadRequest)
	
	if wrappedErr.Type != ErrorTypeValidation {
		t.Errorf("Expected type %s, got %s", ErrorTypeValidation, wrappedErr.Type)
	}
	
	if wrappedErr.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got %s", wrappedErr.Message)
	}
	
	if wrappedErr.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, wrappedErr.StatusCode)
	}
}

func TestPredefinedErrors(t *testing.T) {
	tests := []struct {
		err      *APIError
		expected ErrorType
		status   int
	}{
		{ErrInvalidRequest, ErrorTypeValidation, http.StatusBadRequest},
		{ErrMissingField, ErrorTypeValidation, http.StatusBadRequest},
		{ErrNotFound, ErrorTypeNotFound, http.StatusNotFound},
		{ErrHandlerNotFound, ErrorTypeNotFound, http.StatusNotFound},
		{ErrUnauthorized, ErrorTypeUnauthorized, http.StatusUnauthorized},
		{ErrForbidden, ErrorTypeForbidden, http.StatusForbidden},
		{ErrConflict, ErrorTypeConflict, http.StatusConflict},
		{ErrInternal, ErrorTypeInternal, http.StatusInternalServerError},
		{ErrDatabase, ErrorTypeInternal, http.StatusInternalServerError},
		{ErrRateLimit, ErrorTypeRateLimit, http.StatusTooManyRequests},
	}

	for _, test := range tests {
		if test.err.Type != test.expected {
			t.Errorf("Expected type %s, got %s", test.expected, test.err.Type)
		}
		
		if test.err.StatusCode != test.status {
			t.Errorf("Expected status code %d, got %d", test.status, test.err.StatusCode)
		}
	}
}
