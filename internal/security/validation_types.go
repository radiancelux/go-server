package security

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationResult holds validation results
type ValidationResult struct {
	Valid    bool              `json:"valid"`
	Errors   []ValidationError `json:"errors,omitempty"`
	Warnings []ValidationError `json:"warnings,omitempty"`
}

// ValidationRule defines a validation rule
type ValidationRule struct {
	Field     string
	Required  bool
	MinLength int
	MaxLength int
	Pattern   string
	Min       int
	Max       int
}

// ValidationConfig holds validation configuration
type ValidationConfig struct {
	MaxRequestSize int64
	AllowedMethods []string
	RequiredFields []string
	OptionalFields []string
}
