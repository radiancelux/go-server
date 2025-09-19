package security

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// RateLimiter implements per-IP rate limiting
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
	cleanup  time.Duration
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	RequestsPerMinute int
	WindowDuration    time.Duration
	CleanupInterval   time.Duration
	BurstSize         int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(config RateLimitConfig) *RateLimiter {
	rl := &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    config.RequestsPerMinute,
		window:   config.WindowDuration,
		cleanup:  config.CleanupInterval,
	}
	
	// Start cleanup goroutine
	go rl.cleanupExpired()
	
	return rl
}

// IsAllowed checks if a request from the given IP is allowed
func (rl *RateLimiter) IsAllowed(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	// Get existing requests for this IP
	requests, exists := rl.requests[ip]
	if !exists {
		requests = make([]time.Time, 0)
	}
	
	// Remove old requests outside the window
	var validRequests []time.Time
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validRequests = append(validRequests, reqTime)
		}
	}
	
	// Check if under limit
	if len(validRequests) >= rl.limit {
		return false
	}
	
	// Add current request
	validRequests = append(validRequests, now)
	rl.requests[ip] = validRequests
	
	return true
}

// GetRemainingRequests returns the number of remaining requests for an IP
func (rl *RateLimiter) GetRemainingRequests(ip string) int {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	requests, exists := rl.requests[ip]
	if !exists {
		return rl.limit
	}
	
	// Count valid requests
	validCount := 0
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			validCount++
		}
	}
	
	remaining := rl.limit - validCount
	if remaining < 0 {
		return 0
	}
	
	return remaining
}

// GetResetTime returns when the rate limit resets for an IP
func (rl *RateLimiter) GetResetTime(ip string) time.Time {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	
	now := time.Now()
	cutoff := now.Add(-rl.window)
	
	requests, exists := rl.requests[ip]
	if !exists {
		return now.Add(rl.window)
	}
	
	// Find the oldest valid request
	var oldestTime time.Time
	for _, reqTime := range requests {
		if reqTime.After(cutoff) {
			if oldestTime.IsZero() || reqTime.Before(oldestTime) {
				oldestTime = reqTime
			}
		}
	}
	
	if oldestTime.IsZero() {
		return now.Add(rl.window)
	}
	
	return oldestTime.Add(rl.window)
}

// cleanupExpired removes expired entries from the rate limiter
func (rl *RateLimiter) cleanupExpired() {
	ticker := time.NewTicker(rl.cleanup)
	defer ticker.Stop()
	
	for range ticker.C {
		rl.mutex.Lock()
		now := time.Now()
		cutoff := now.Add(-rl.window)
		
		for ip, requests := range rl.requests {
			var validRequests []time.Time
			for _, reqTime := range requests {
				if reqTime.After(cutoff) {
					validRequests = append(validRequests, reqTime)
				}
			}
			
			if len(validRequests) == 0 {
				delete(rl.requests, ip)
			} else {
				rl.requests[ip] = validRequests
			}
		}
		rl.mutex.Unlock()
	}
}

// GetClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if comma := strings.Index(xff, ","); comma != -1 {
			xff = xff[:comma]
		}
		xff = strings.TrimSpace(xff)
		if xff != "" {
			return xff
		}
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	
	// Fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	
	return ip
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rateLimiter *RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := GetClientIP(r)
			
			if !rateLimiter.IsAllowed(clientIP) {
				remaining := rateLimiter.GetRemainingRequests(clientIP)
				resetTime := rateLimiter.GetResetTime(clientIP)
				
				w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rateLimiter.limit))
				w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
				w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
				w.Header().Set("Retry-After", fmt.Sprintf("%d", int(time.Until(resetTime).Seconds())))
				
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			
			// Add rate limit headers to successful requests
			remaining := rateLimiter.GetRemainingRequests(clientIP)
			resetTime := rateLimiter.GetResetTime(clientIP)
			
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rateLimiter.limit))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprintf("%d", resetTime.Unix()))
			
			next.ServeHTTP(w, r)
		})
	}
}
