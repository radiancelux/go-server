package security

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationResult holds validation results
type ValidationResult struct {
	Valid   bool              `json:"valid"`
	Errors  []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
}

// Validator provides request validation functions
type Validator struct {
	sanitizer *Sanitizer
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		sanitizer: NewSanitizer(),
	}
}

// ValidateRequest validates an HTTP request
func (v *Validator) ValidateRequest(r *http.Request) ValidationResult {
	var errors []ValidationError
	var warnings []ValidationError
	
	// Validate HTTP method
	if r.Method != http.MethodGet && r.Method != http.MethodPost && 
	   r.Method != http.MethodPut && r.Method != http.MethodDelete && 
	   r.Method != http.MethodOptions {
		errors = append(errors, ValidationError{
			Field:   "method",
			Message: "Invalid HTTP method",
			Value:   r.Method,
		})
	}
	
	// Validate Content-Type for POST/PUT requests
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		contentType := r.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			errors = append(errors, ValidationError{
				Field:   "content-type",
				Message: "Content-Type must be application/json",
				Value:   contentType,
			})
		}
	}
	
	// Validate Content-Length
	contentLength := r.Header.Get("Content-Length")
	if contentLength != "" {
		if length, err := strconv.Atoi(contentLength); err == nil {
			if length > 1024*1024 { // 1MB limit
				errors = append(errors, ValidationError{
					Field:   "content-length",
					Message: "Request body too large",
					Value:   contentLength,
				})
			}
		}
	}
	
	// Validate User-Agent
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		warnings = append(warnings, ValidationError{
			Field:   "user-agent",
			Message: "Missing User-Agent header",
		})
	} else if len(userAgent) > 500 {
		warnings = append(warnings, ValidationError{
			Field:   "user-agent",
			Message: "User-Agent header too long",
			Value:   userAgent[:50] + "...",
		})
	}
	
	// Validate URL path
	if !v.isValidPath(r.URL.Path) {
		errors = append(errors, ValidationError{
			Field:   "path",
			Message: "Invalid URL path",
			Value:   r.URL.Path,
		})
	}
	
	// Validate query parameters
	for key, values := range r.URL.Query() {
		if len(key) > 100 {
			errors = append(errors, ValidationError{
				Field:   "query." + key,
				Message: "Query parameter name too long",
				Value:   key,
			})
		}
		
		for _, value := range values {
			if len(value) > 1000 {
				errors = append(errors, ValidationError{
					Field:   "query." + key,
					Message: "Query parameter value too long",
					Value:   value[:50] + "...",
				})
			}
			
			if !v.sanitizer.ValidateString(value) {
				errors = append(errors, ValidationError{
					Field:   "query." + key,
					Message: "Query parameter contains invalid characters",
					Value:   value,
				})
			}
		}
	}
	
	return ValidationResult{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}
}

// ValidateJSONRequest validates a JSON request body
func (v *Validator) ValidateJSONRequest(r *http.Request, target interface{}) ValidationResult {
	var errors []ValidationError
	var warnings []ValidationError
	
	// First validate the basic request
	basicResult := v.ValidateRequest(r)
	if !basicResult.Valid {
		return basicResult
	}
	
	// Validate JSON body
	if r.Body == nil {
		errors = append(errors, ValidationError{
			Field:   "body",
			Message: "Request body is required",
		})
		return ValidationResult{Valid: false, Errors: errors}
	}
	
	// Decode JSON
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	
	if err := decoder.Decode(target); err != nil {
		errors = append(errors, ValidationError{
			Field:   "body",
			Message: "Invalid JSON format: " + err.Error(),
		})
		return ValidationResult{Valid: false, Errors: errors}
	}
	
	// Validate specific fields based on target type
	fieldErrors := v.validateFields(target)
	errors = append(errors, fieldErrors...)
	
	return ValidationResult{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}
}

// validateFields validates specific fields in the target struct
func (v *Validator) validateFields(target interface{}) []ValidationError {
	var errors []ValidationError
	
	// This is a simplified version - in a real implementation,
	// you would use reflection or a validation library like go-playground/validator
	
	// For now, we'll add basic validation for common fields
	// In a real implementation, you would use struct tags and reflection
	
	return errors
}

// ValidateString validates a string field
func (v *Validator) ValidateString(value, fieldName string, required bool, maxLength int) []ValidationError {
	var errors []ValidationError
	
	if required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Field is required",
		})
		return errors
	}
	
	if value != "" {
		if len(value) > maxLength {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: fmt.Sprintf("Field too long (max %d characters)", maxLength),
				Value:   value,
			})
		}
		
		if !v.sanitizer.ValidateString(value) {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: "Field contains invalid characters",
				Value:   value,
			})
		}
	}
	
	return errors
}

// ValidateEmail validates an email field
func (v *Validator) ValidateEmail(value, fieldName string, required bool) []ValidationError {
	var errors []ValidationError
	
	if required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Field is required",
		})
		return errors
	}
	
	if value != "" {
		if !v.sanitizer.ValidateEmail(value) {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: "Invalid email format",
				Value:   value,
			})
		}
	}
	
	return errors
}

// ValidateInteger validates an integer field
func (v *Validator) ValidateInteger(value, fieldName string, required bool, min, max int) []ValidationError {
	var errors []ValidationError
	
	if required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Field is required",
		})
		return errors
	}
	
	if value != "" {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: "Invalid integer format",
				Value:   value,
			})
		} else {
			if intValue < min || intValue > max {
				errors = append(errors, ValidationError{
					Field:   fieldName,
					Message: fmt.Sprintf("Value must be between %d and %d", min, max),
					Value:   value,
				})
			}
		}
	}
	
	return errors
}

// isValidPath validates URL path
func (v *Validator) isValidPath(path string) bool {
	// Basic path validation
	if len(path) > 1000 {
		return false
	}
	
	// Check for path traversal attempts
	if strings.Contains(path, "..") || strings.Contains(path, "//") {
		return false
	}
	
	// Check for null bytes
	if strings.Contains(path, "\x00") {
		return false
	}
	
	return true
}

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, result ValidationResult) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	
	response := map[string]interface{}{
		"status": "error",
		"type":   "validation_error",
		"message": "Validation failed",
		"errors": result.Errors,
	}
	
	if len(result.Warnings) > 0 {
		response["warnings"] = result.Warnings
	}
	
	json.NewEncoder(w).Encode(response)
}
