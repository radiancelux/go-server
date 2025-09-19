package security

import (
	"encoding/json"
	"net/http"
)

// WriteValidationError writes a validation error response
func WriteValidationError(w http.ResponseWriter, result ValidationResult) {
	w.Header().Set("Content-Type", "application/json")
	
	if result.Valid {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
	
	json.NewEncoder(w).Encode(result)
}

// WriteValidationSuccess writes a validation success response
func WriteValidationSuccess(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"valid":   true,
		"message": message,
	}
	
	json.NewEncoder(w).Encode(response)
}

// AddError adds an error to a validation result
func (vr *ValidationResult) AddError(field, message, value string) {
	vr.Errors = append(vr.Errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
	vr.Valid = false
}

// AddWarning adds a warning to a validation result
func (vr *ValidationResult) AddWarning(field, message, value string) {
	vr.Warnings = append(vr.Warnings, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors returns true if there are validation errors
func (vr *ValidationResult) HasErrors() bool {
	return len(vr.Errors) > 0
}

// HasWarnings returns true if there are validation warnings
func (vr *ValidationResult) HasWarnings() bool {
	return len(vr.Warnings) > 0
}

// GetErrorMessages returns all error messages as a slice
func (vr *ValidationResult) GetErrorMessages() []string {
	messages := make([]string, len(vr.Errors))
	for i, err := range vr.Errors {
		messages[i] = err.Message
	}
	return messages
}

// GetWarningMessages returns all warning messages as a slice
func (vr *ValidationResult) GetWarningMessages() []string {
	messages := make([]string, len(vr.Warnings))
	for i, warning := range vr.Warnings {
		messages[i] = warning.Message
	}
	return messages
}
