package security

import (
	"strings"
	"testing"
)

func TestSanitizer_SanitizeString(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Normal string",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "String with HTML",
			input:    "<script>alert('xss')</script>",
			expected: "&lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;",
		},
		{
			name:     "String with null bytes",
			input:    "Hello\x00World",
			expected: "HelloWorld",
		},
		{
			name:     "String with newlines",
			input:    "Hello\nWorld\r\nTest",
			expected: "Hello World Test",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Very long string",
			input:    strings.Repeat("a", 2000),
			expected: strings.Repeat("a", 1000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_SanitizeEmail(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid email",
			input:    "test@example.com",
			expected: "test@example.com",
		},
		{
			name:     "Email with uppercase",
			input:    "TEST@EXAMPLE.COM",
			expected: "test@example.com",
		},
		{
			name:     "Email with spaces",
			input:    " test@example.com ",
			expected: "test@example.com",
		},
		{
			name:     "Invalid email",
			input:    "not-an-email",
			expected: "",
		},
		{
			name:     "Empty email",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.SanitizeEmail(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeEmail() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_ValidateString(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Safe string",
			input:    "Hello World",
			expected: true,
		},
		{
			name:     "String with script tag",
			input:    "<script>alert('xss')</script>",
			expected: false,
		},
		{
			name:     "String with javascript",
			input:    "javascript:alert('xss')",
			expected: false,
		},
		{
			name:     "String with onload",
			input:    "onload=alert('xss')",
			expected: false,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: true,
		},
		{
			name:     "Very long string",
			input:    string(make([]byte, 2000)),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.ValidateString(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateString() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSanitizer_ValidateEmail(t *testing.T) {
	sanitizer := NewSanitizer()

	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid email",
			input:    "test@example.com",
			expected: true,
		},
		{
			name:     "Invalid email",
			input:    "not-an-email",
			expected: false,
		},
		{
			name:     "Empty email",
			input:    "",
			expected: false,
		},
		{
			name:     "Email without domain",
			input:    "test@",
			expected: false,
		},
		{
			name:     "Email without @",
			input:    "testexample.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizer.ValidateEmail(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateEmail() = %v, want %v", result, tt.expected)
			}
		})
	}
}
