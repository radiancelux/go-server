package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/radiancelux/go-server/internal/config"
	"github.com/radiancelux/go-server/internal/errors"
	"github.com/radiancelux/go-server/internal/interfaces"
)

// RequestIDKey is the context key for request ID
type RequestIDKey struct{}

// Middleware represents a middleware function
type Middleware func(http.Handler) http.Handler

// Chain chains multiple middleware functions
func Chain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				// Generate a new request ID
				bytes := make([]byte, 16)
				rand.Read(bytes)
				requestID = hex.EncodeToString(bytes)
			}

			// Add request ID to context
			ctx := context.WithValue(r.Context(), RequestIDKey{}, requestID)
			r = r.WithContext(ctx)

			// Add request ID to response headers
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(logger interfaces.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			requestID := GetRequestID(r.Context())

			logger.Info("Request started: %s %s (ID: %s)", r.Method, r.URL.Path, requestID)

			// Create a response writer wrapper to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			logger.Info("Request completed: %s %s %d %v (ID: %s)",
				r.Method, r.URL.Path, wrapped.statusCode, duration, requestID)
		})
	}
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if cfg.Security.EnableCORS {
				origin := r.Header.Get("Origin")

				// Check if origin is allowed
				if isOriginAllowed(origin, cfg.Security.CORSOrigins) {
					w.Header().Set("Access-Control-Allow-Origin", origin)
				} else if contains(cfg.Security.CORSOrigins, "*") {
					w.Header().Set("Access-Control-Allow-Origin", "*")
				}

				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID")
				w.Header().Set("Access-Control-Max-Age", "86400")
			}

			// Handle preflight requests
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// SecurityHeadersMiddleware adds security headers
func SecurityHeadersMiddleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Content-Type-Options", "nosniff")
			w.Header().Set("X-Frame-Options", "DENY")
			w.Header().Set("X-XSS-Protection", "1; mode=block")
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.Header().Set("Content-Security-Policy", "default-src 'self'")

			next.ServeHTTP(w, r)
		})
	}
}

// RequestSizeMiddleware limits request body size
func RequestSizeMiddleware(cfg *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > cfg.Security.MaxRequestSize {
				requestID := GetRequestID(r.Context())
				err := errors.ErrInvalidRequest.WithDetails(
					fmt.Sprintf("Request too large: %d bytes (max: %d)",
						r.ContentLength, cfg.Security.MaxRequestSize)).
					WithRequestID(requestID)

				writeErrorResponse(w, err)
				return
			}

			// Limit the request body reader
			r.Body = http.MaxBytesReader(w, r.Body, cfg.Security.MaxRequestSize)

			next.ServeHTTP(w, r)
		})
	}
}

// RecoveryMiddleware recovers from panics
func RecoveryMiddleware(logger interfaces.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					requestID := GetRequestID(r.Context())
					logger.Error("Panic recovered: %v (ID: %s)", err, requestID)

					apiErr := errors.ErrInternal.WithRequestID(requestID)
					writeErrorResponse(w, apiErr)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey{}).(string); ok {
		return requestID
	}
	return ""
}

// isOriginAllowed checks if an origin is in the allowed list
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if origin == allowed {
			return true
		}
	}
	return false
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// writeErrorResponse writes an error response
func writeErrorResponse(w http.ResponseWriter, err *errors.APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.StatusCode)

	// In a real implementation, you'd use json.Marshal
	// For now, we'll write a simple error response
	response := fmt.Sprintf(`{"error": {"type": "%s", "message": "%s"}}`,
		err.Type, err.Message)
	w.Write([]byte(response))
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
