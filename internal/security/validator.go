package security

import (
	"net/http"
)

// Validator provides comprehensive validation functionality
type Validator struct {
	httpValidator  *HTTPValidator
	fieldValidator *FieldValidator
	sanitizer      *Sanitizer
}

// NewValidator creates a new validator
func NewValidator() *Validator {
	return &Validator{
		httpValidator:  NewHTTPValidator(),
		fieldValidator: NewFieldValidator(),
		sanitizer:      NewSanitizer(),
	}
}

// ValidateRequest validates an HTTP request
func (v *Validator) ValidateRequest(r *http.Request) ValidationResult {
	return v.httpValidator.ValidateRequest(r)
}

// ValidateJSONRequest validates a JSON request
func (v *Validator) ValidateJSONRequest(r *http.Request, target interface{}) ValidationResult {
	return v.httpValidator.ValidateJSONRequest(r, target)
}

// ValidateString validates a string field
func (v *Validator) ValidateString(value, fieldName string, required bool, maxLength int) []ValidationError {
	return v.fieldValidator.ValidateString(value, fieldName, required, maxLength)
}

// ValidateEmail validates an email field
func (v *Validator) ValidateEmail(value, fieldName string, required bool) []ValidationError {
	return v.fieldValidator.ValidateEmail(value, fieldName, required)
}

// ValidateInteger validates an integer field
func (v *Validator) ValidateInteger(value, fieldName string, required bool, min, max int) []ValidationError {
	return v.fieldValidator.ValidateInteger(value, fieldName, required, min, max)
}

// ValidateUsername validates a username field
func (v *Validator) ValidateUsername(value, fieldName string, required bool) []ValidationError {
	return v.fieldValidator.ValidateUsername(value, fieldName, required)
}

// ValidatePassword validates a password field
func (v *Validator) ValidatePassword(value, fieldName string, required bool) []ValidationError {
	return v.fieldValidator.ValidatePassword(value, fieldName, required)
}

// Note: WriteValidationError and WriteValidationSuccess are available from validation_utils.go