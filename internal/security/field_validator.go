package security

import (
	"regexp"
	"strconv"
	"strings"
)

// FieldValidator handles field-level validation
type FieldValidator struct {
	sanitizer *Sanitizer
}

// NewFieldValidator creates a new field validator
func NewFieldValidator() *FieldValidator {
	return &FieldValidator{
		sanitizer: NewSanitizer(),
	}
}

// ValidateString validates a string field
func (v *FieldValidator) ValidateString(value, fieldName string, required bool, maxLength int) []ValidationError {
	var errors []ValidationError

	// Check if required field is empty
	if required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Field is required",
			Value:   value,
		})
		return errors
	}

	// Skip validation if field is empty and not required
	if !required && strings.TrimSpace(value) == "" {
		return errors
	}

	// Check length
	if maxLength > 0 && len(value) > maxLength {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Field too long",
			Value:   value,
		})
	}

	// Check for dangerous characters using sanitizer validation
	if !v.sanitizer.ValidateString(value) {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Field contains dangerous characters",
			Value:   value,
		})
	}

	return errors
}

// ValidateEmail validates an email field
func (v *FieldValidator) ValidateEmail(value, fieldName string, required bool) []ValidationError {
	var errors []ValidationError

	// Check if required field is empty
	if required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Email is required",
			Value:   value,
		})
		return errors
	}

	// Skip validation if field is empty and not required
	if !required && strings.TrimSpace(value) == "" {
		return errors
	}

	// Basic email validation regex
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(value) {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Invalid email format",
			Value:   value,
		})
	}

	// Check length
	if len(value) > 254 { // RFC 5321 limit
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Email too long",
			Value:   value,
		})
	}

	return errors
}

// ValidateInteger validates an integer field
func (v *FieldValidator) ValidateInteger(value, fieldName string, required bool, min, max int) []ValidationError {
	var errors []ValidationError

	// Check if required field is empty
	if required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Field is required",
			Value:   value,
		})
		return errors
	}

	// Skip validation if field is empty and not required
	if !required && strings.TrimSpace(value) == "" {
		return errors
	}

	// Parse integer
	intValue, err := strconv.Atoi(value)
	if err != nil {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Invalid integer format",
			Value:   value,
		})
		return errors
	}

	// Check range
	if min != max { // Only check range if min != max (indicating range validation is enabled)
		if intValue < min {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: "Value too small",
				Value:   value,
			})
		}
		if intValue > max {
			errors = append(errors, ValidationError{
				Field:   fieldName,
				Message: "Value too large",
				Value:   value,
			})
		}
	}

	return errors
}

// ValidateUsername validates a username field
func (v *FieldValidator) ValidateUsername(value, fieldName string, required bool) []ValidationError {
	var errors []ValidationError

	// Check if required field is empty
	if required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Username is required",
			Value:   value,
		})
		return errors
	}

	// Skip validation if field is empty and not required
	if !required && strings.TrimSpace(value) == "" {
		return errors
	}

	// Check length
	if len(value) < 3 {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Username too short (minimum 3 characters)",
			Value:   value,
		})
	}

	if len(value) > 20 {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Username too long (maximum 20 characters)",
			Value:   value,
		})
	}

	// Check for valid characters (alphanumeric and underscore only)
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(value) {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Username contains invalid characters (only letters, numbers, and underscores allowed)",
			Value:   value,
		})
	}

	return errors
}

// ValidatePassword validates a password field
func (v *FieldValidator) ValidatePassword(value, fieldName string, required bool) []ValidationError {
	var errors []ValidationError

	// Check if required field is empty
	if required && strings.TrimSpace(value) == "" {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Password is required",
			Value:   "",
		})
		return errors
	}

	// Skip validation if field is empty and not required
	if !required && strings.TrimSpace(value) == "" {
		return errors
	}

	// Check length
	if len(value) < 8 {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Password too short (minimum 8 characters)",
			Value:   "",
		})
	}

	if len(value) > 128 {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Password too long (maximum 128 characters)",
			Value:   "",
		})
	}

	// Check for at least one letter and one number
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(value)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(value)

	if !hasLetter {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Password must contain at least one letter",
			Value:   "",
		})
	}

	if !hasNumber {
		errors = append(errors, ValidationError{
			Field:   fieldName,
			Message: "Password must contain at least one number",
			Value:   "",
		})
	}

	return errors
}
