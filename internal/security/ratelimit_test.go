package security

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_IsAllowed(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 2,
		WindowDuration:    time.Minute,
		CleanupInterval:   time.Minute,
		BurstSize:         5,
	}

	rl := NewRateLimiter(config)
	ip := "192.168.1.1"

	// First two requests should be allowed
	if !rl.IsAllowed(ip) {
		t.Error("First request should be allowed")
	}
	if !rl.IsAllowed(ip) {
		t.Error("Second request should be allowed")
	}

	// Third request should be denied
	if rl.IsAllowed(ip) {
		t.Error("Third request should be denied")
	}
}

func TestRateLimiter_GetRemainingRequests(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 3,
		WindowDuration:    time.Minute,
		CleanupInterval:   time.Minute,
		BurstSize:         5,
	}

	rl := NewRateLimiter(config)
	ip := "192.168.1.1"

	// Initially should have 3 remaining
	if remaining := rl.GetRemainingRequests(ip); remaining != 3 {
		t.Errorf("Expected 3 remaining requests, got %d", remaining)
	}

	// After one request, should have 2 remaining
	rl.IsAllowed(ip)
	if remaining := rl.GetRemainingRequests(ip); remaining != 2 {
		t.Errorf("Expected 2 remaining requests, got %d", remaining)
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 1,
		WindowDuration:    time.Minute,
		CleanupInterval:   time.Minute,
		BurstSize:         5,
	}

	rl := NewRateLimiter(config)
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	// Both IPs should be allowed for first request
	if !rl.IsAllowed(ip1) {
		t.Error("First IP first request should be allowed")
	}
	if !rl.IsAllowed(ip2) {
		t.Error("Second IP first request should be allowed")
	}

	// Second requests should be denied
	if rl.IsAllowed(ip1) {
		t.Error("First IP second request should be denied")
	}
	if rl.IsAllowed(ip2) {
		t.Error("Second IP second request should be denied")
	}
}

func TestRateLimitMiddleware(t *testing.T) {
	config := RateLimitConfig{
		RequestsPerMinute: 1,
		WindowDuration:    time.Minute,
		CleanupInterval:   time.Minute,
		BurstSize:         5,
	}

	rl := NewRateLimiter(config)
	middleware := RateLimitMiddleware(rl)

	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Wrap with middleware
	wrappedHandler := middleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.1.1:12345"

	// First request should succeed
	w1 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w1, req)
	if w1.Code != http.StatusOK {
		t.Errorf("First request should succeed, got status %d", w1.Code)
	}

	// Check rate limit headers
	if w1.Header().Get("X-RateLimit-Limit") == "" {
		t.Error("Missing X-RateLimit-Limit header")
	}
	if w1.Header().Get("X-RateLimit-Remaining") == "" {
		t.Error("Missing X-RateLimit-Remaining header")
	}

	// Second request should be rate limited
	w2 := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(w2, req)
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Second request should be rate limited, got status %d", w2.Code)
	}

	// Check rate limit headers
	if w2.Header().Get("X-RateLimit-Limit") == "" {
		t.Error("Missing X-RateLimit-Limit header")
	}
	if w2.Header().Get("X-RateLimit-Remaining") == "" {
		t.Error("Missing X-RateLimit-Remaining header")
	}
	if w2.Header().Get("Retry-After") == "" {
		t.Error("Missing Retry-After header")
	}
}

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name     string
		request  *http.Request
		expected string
	}{
		{
			name: "X-Forwarded-For header",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Forwarded-For", "192.168.1.1")
				return req
			}(),
			expected: "192.168.1.1",
		},
		{
			name: "X-Real-IP header",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Real-IP", "192.168.1.2")
				return req
			}(),
			expected: "192.168.1.2",
		},
		{
			name: "RemoteAddr fallback",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.RemoteAddr = "192.168.1.3:12345"
				return req
			}(),
			expected: "192.168.1.3",
		},
		{
			name: "X-Forwarded-For with multiple IPs",
			request: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Forwarded-For", "192.168.1.1, 10.0.0.1, 172.16.0.1")
				return req
			}(),
			expected: "192.168.1.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetClientIP(tt.request)
			if result != tt.expected {
				t.Errorf("GetClientIP() = %v, want %v", result, tt.expected)
			}
		})
	}
}
