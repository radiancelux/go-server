package security

import (
	"net/http"
	"strconv"
	"strings"
)

// HTTPValidator handles HTTP request validation
type HTTPValidator struct {
	sanitizer *Sanitizer
}

// NewHTTPValidator creates a new HTTP validator
func NewHTTPValidator() *HTTPValidator {
	return &HTTPValidator{
		sanitizer: NewSanitizer(),
	}
}

// ValidateRequest validates an HTTP request
func (v *HTTPValidator) ValidateRequest(r *http.Request) ValidationResult {
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
	if (r.Method == http.MethodPost || r.Method == http.MethodPut) &&
		r.Header.Get("Content-Type") != "application/json" {
		warnings = append(warnings, ValidationError{
			Field:   "content-type",
			Message: "Expected application/json",
			Value:   r.Header.Get("Content-Type"),
		})
	}

	// Validate Content-Length
	if contentLength := r.Header.Get("Content-Length"); contentLength != "" {
		if length, err := strconv.ParseInt(contentLength, 10, 64); err == nil {
			if length > 10*1024*1024 { // 10MB limit
				errors = append(errors, ValidationError{
					Field:   "content-length",
					Message: "Request too large",
					Value:   contentLength,
				})
			}
		}
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
		for _, value := range values {
			if len(value) > 1000 { // Reasonable limit for query params
				errors = append(errors, ValidationError{
					Field:   "query." + key,
					Message: "Query parameter too long",
					Value:   value[:100] + "...",
				})
			}
		}
	}

	// Validate headers
	for key, values := range r.Header {
		if len(key) > 100 || len(values) > 10 {
			errors = append(errors, ValidationError{
				Field:   "header." + key,
				Message: "Invalid header",
				Value:   key,
			})
		}
	}

	return ValidationResult{
		Valid:    len(errors) == 0,
		Errors:   errors,
		Warnings: warnings,
	}
}

// ValidateJSONRequest validates a JSON request
func (v *HTTPValidator) ValidateJSONRequest(r *http.Request, target interface{}) ValidationResult {
	result := v.ValidateRequest(r)
	if !result.Valid {
		return result
	}

	// Additional JSON-specific validation
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		if r.Header.Get("Content-Type") != "application/json" {
			result.Errors = append(result.Errors, ValidationError{
				Field:   "content-type",
				Message: "Content-Type must be application/json",
				Value:   r.Header.Get("Content-Type"),
			})
			result.Valid = false
		}
	}

	// Validate JSON structure if target is provided
	if target != nil {
		fieldErrors := v.validateFields(target)
		result.Errors = append(result.Errors, fieldErrors...)
		result.Valid = len(result.Errors) == 0
	}

	return result
}

// isValidPath checks if a URL path is valid
func (v *HTTPValidator) isValidPath(path string) bool {
	// Basic path validation
	if path == "" || path[0] != '/' {
		return false
	}

	// Check for dangerous patterns
	dangerousPatterns := []string{
		"../",
		"..\\",
		"//",
		"\\",
		"<",
		">",
		"\"",
		"'",
		"&",
		"#",
		"%",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(path, pattern) {
			return false
		}
	}

	return true
}

// validateFields validates struct fields using reflection
func (v *HTTPValidator) validateFields(target interface{}) []ValidationError {
	// This would use reflection to validate struct fields
	// For now, return empty slice as placeholder
	return []ValidationError{}
}
