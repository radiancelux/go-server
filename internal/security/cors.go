package security

import (
	"fmt"
	"net/http"
	"strings"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// DefaultCORSConfig returns a default CORS configuration
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD", "PATCH"},
		AllowedHeaders: []string{
			"Accept",
			"Accept-Language",
			"Content-Language",
			"Content-Type",
			"Authorization",
			"X-Requested-With",
			"X-Request-ID",
			"X-Forwarded-For",
			"X-Real-IP",
		},
		ExposedHeaders: []string{
			"X-Request-ID",
			"X-RateLimit-Limit",
			"X-RateLimit-Remaining",
			"X-RateLimit-Reset",
		},
		AllowCredentials: false,
		MaxAge:           86400, // 24 hours
	}
}

// CORSHandler handles CORS requests
type CORSHandler struct {
	config CORSConfig
}

// NewCORSHandler creates a new CORS handler
func NewCORSHandler(config CORSConfig) *CORSHandler {
	return &CORSHandler{config: config}
}

// HandleCORS handles CORS preflight and actual requests
func (c *CORSHandler) HandleCORS(w http.ResponseWriter, r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Handle preflight request first
	if r.Method == http.MethodOptions {
		// Set CORS headers for preflight
		c.setCORSHeaders(w, origin)
		w.WriteHeader(http.StatusOK)
		return true
	}

	// For non-OPTIONS requests, check if origin is allowed
	if !c.isOriginAllowed(origin) {
		return false
	}

	// Set CORS headers
	c.setCORSHeaders(w, origin)

	return false
}

// isOriginAllowed checks if the origin is allowed
func (c *CORSHandler) isOriginAllowed(origin string) bool {
	if origin == "" {
		return false
	}

	// Check for wildcard
	for _, allowedOrigin := range c.config.AllowedOrigins {
		if allowedOrigin == "*" {
			return true
		}
		if allowedOrigin == origin {
			return true
		}
	}

	return false
}

// setCORSHeaders sets the CORS headers
func (c *CORSHandler) setCORSHeaders(w http.ResponseWriter, origin string) {
	// Set Access-Control-Allow-Origin
	if len(c.config.AllowedOrigins) > 0 && c.config.AllowedOrigins[0] == "*" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	// Set Access-Control-Allow-Methods
	if len(c.config.AllowedMethods) > 0 {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(c.config.AllowedMethods, ", "))
	}

	// Set Access-Control-Allow-Headers
	if len(c.config.AllowedHeaders) > 0 {
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(c.config.AllowedHeaders, ", "))
	}

	// Set Access-Control-Expose-Headers
	if len(c.config.ExposedHeaders) > 0 {
		w.Header().Set("Access-Control-Expose-Headers", strings.Join(c.config.ExposedHeaders, ", "))
	}

	// Set Access-Control-Allow-Credentials
	if c.config.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	// Set Access-Control-Max-Age
	if c.config.MaxAge > 0 {
		w.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", c.config.MaxAge))
	}
}

// CORSMiddleware creates a CORS middleware
func CORSMiddleware(config CORSConfig) func(http.Handler) http.Handler {
	corsHandler := NewCORSHandler(config)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Handle CORS
			if corsHandler.HandleCORS(w, r) {
				return
			}

			// Continue to next handler
			next.ServeHTTP(w, r)
		})
	}
}

// ValidateCORSRequest validates a CORS request
func ValidateCORSRequest(r *http.Request, config CORSConfig) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true // No origin header, not a CORS request
	}

	corsHandler := NewCORSHandler(config)
	return corsHandler.isOriginAllowed(origin)
}
