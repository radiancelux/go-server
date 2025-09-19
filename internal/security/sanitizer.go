package security

import (
	"html"
	"regexp"
	"strings"
	"unicode"
)

// Sanitizer provides input sanitization functions
type Sanitizer struct {
	// Regular expressions for validation
	emailRegex        *regexp.Regexp
	alphanumericRegex *regexp.Regexp
	safeStringRegex   *regexp.Regexp
}

// NewSanitizer creates a new sanitizer instance
func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		emailRegex:        regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		alphanumericRegex: regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,!?]+$`),
		safeStringRegex:   regexp.MustCompile(`^[a-zA-Z0-9\s\-_.,!?@#$%^&*()+={}[\]|\\:";'<>?/~` + "`" + `]+$`),
	}
}

// SanitizeString sanitizes a string input
func (s *Sanitizer) SanitizeString(input string) string {
	if input == "" {
		return ""
	}

	// HTML escape to prevent XSS
	sanitized := html.EscapeString(input)

	// Remove null bytes and control characters
	sanitized = strings.ReplaceAll(sanitized, "\x00", "")
	sanitized = strings.ReplaceAll(sanitized, "\r", "")
	sanitized = strings.ReplaceAll(sanitized, "\n", " ")
	sanitized = strings.ReplaceAll(sanitized, "\t", " ")

	// Trim whitespace
	sanitized = strings.TrimSpace(sanitized)

	// Limit length (prevent extremely long inputs)
	if len(sanitized) > 1000 {
		sanitized = sanitized[:1000]
	}

	return sanitized
}

// SanitizeEmail sanitizes an email address
func (s *Sanitizer) SanitizeEmail(email string) string {
	if email == "" {
		return ""
	}

	// Convert to lowercase and trim
	email = strings.ToLower(strings.TrimSpace(email))

	// Basic validation
	if !s.emailRegex.MatchString(email) {
		return ""
	}

	// HTML escape
	email = html.EscapeString(email)

	// Limit length
	if len(email) > 254 {
		email = email[:254]
	}

	return email
}

// SanitizeAlphanumeric sanitizes alphanumeric input
func (s *Sanitizer) SanitizeAlphanumeric(input string) string {
	if input == "" {
		return ""
	}

	// Remove non-alphanumeric characters except spaces, hyphens, underscores, dots, commas
	sanitized := ""
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) ||
			r == ' ' || r == '-' || r == '_' || r == '.' || r == ',' {
			sanitized += string(r)
		}
	}

	// Trim and limit length
	sanitized = strings.TrimSpace(sanitized)
	if len(sanitized) > 500 {
		sanitized = sanitized[:500]
	}

	return sanitized
}

// SanitizeSafeString sanitizes input for safe display
func (s *Sanitizer) SanitizeSafeString(input string) string {
	if input == "" {
		return ""
	}

	// HTML escape
	sanitized := html.EscapeString(input)

	// Remove potentially dangerous characters
	sanitized = strings.ReplaceAll(sanitized, "<script", "&lt;script")
	sanitized = strings.ReplaceAll(sanitized, "</script", "&lt;/script")
	sanitized = strings.ReplaceAll(sanitized, "javascript:", "")
	sanitized = strings.ReplaceAll(sanitized, "vbscript:", "")
	sanitized = strings.ReplaceAll(sanitized, "onload=", "")
	sanitized = strings.ReplaceAll(sanitized, "onerror=", "")

	// Trim and limit length
	sanitized = strings.TrimSpace(sanitized)
	if len(sanitized) > 1000 {
		sanitized = sanitized[:1000]
	}

	return sanitized
}

// ValidateString validates if a string is safe
func (s *Sanitizer) ValidateString(input string) bool {
	if input == "" {
		return true
	}

	// Check for dangerous patterns
	dangerousPatterns := []string{
		"<script",
		"</script",
		"javascript:",
		"vbscript:",
		"onload=",
		"onerror=",
		"<iframe",
		"<object",
		"<embed",
		"<form",
		"<input",
		"<textarea",
		"<select",
		"<button",
	}

	lowerInput := strings.ToLower(input)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(lowerInput, pattern) {
			return false
		}
	}

	// Check length
	if len(input) > 1000 {
		return false
	}

	return true
}

// ValidateEmail validates an email address
func (s *Sanitizer) ValidateEmail(email string) bool {
	if email == "" {
		return false
	}

	return s.emailRegex.MatchString(email)
}

// ValidateAlphanumeric validates alphanumeric input
func (s *Sanitizer) ValidateAlphanumeric(input string) bool {
	if input == "" {
		return true
	}

	return s.alphanumericRegex.MatchString(input) && len(input) <= 500
}

// SanitizeUserInput sanitizes user input based on type
func (s *Sanitizer) SanitizeUserInput(input string, inputType string) string {
	switch inputType {
	case "email":
		return s.SanitizeEmail(input)
	case "alphanumeric":
		return s.SanitizeAlphanumeric(input)
	case "safe":
		return s.SanitizeSafeString(input)
	default:
		return s.SanitizeString(input)
	}
}

// ValidateUserInput validates user input based on type
func (s *Sanitizer) ValidateUserInput(input string, inputType string) bool {
	switch inputType {
	case "email":
		return s.ValidateEmail(input)
	case "alphanumeric":
		return s.ValidateAlphanumeric(input)
	case "safe":
		return s.ValidateString(input)
	default:
		return s.ValidateString(input)
	}
}
