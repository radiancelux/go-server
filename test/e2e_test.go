package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"go-server/internal/config"
	"go-server/internal/server"
)

// TestServer represents a test server instance
type TestServer struct {
	server  *server.Server
	client  *http.Client
	baseURL string
}

// NewTestServer creates a new test server
func NewTestServer(t *testing.T) *TestServer {
	// Use a random port to avoid conflicts (use 6xxx range for E2E tests)
	port := fmt.Sprintf("6%03d", (time.Now().UnixNano()/1000)%1000)
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: port,
		},
		Security: config.SecurityConfig{
			MaxRequestSize: 1024 * 1024,
			RateLimitRPS:   10000, // Very high limit for tests
			RateLimitBurst: 20000,
		},
	}

	srv := server.NewServer(cfg)

	// Start server in background
	go func() {
		if err := srv.Start(); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	return &TestServer{
		server:  srv,
		client:  &http.Client{Timeout: 5 * time.Second},
		baseURL: fmt.Sprintf("http://localhost:%s", port),
	}
}

// TestHealthEndpoint tests the health endpoint
func TestHealthEndpoint(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.cleanup()

	resp, err := ts.client.Get(ts.baseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", response["status"])
	}

	// Check for request ID header
	if resp.Header.Get("X-Request-ID") == "" {
		t.Error("Expected X-Request-ID header to be present")
	}

	// Check for security headers
	securityHeaders := []string{
		"X-Content-Type-Options",
		"X-Frame-Options",
		"X-XSS-Protection",
		"Referrer-Policy",
	}

	for _, header := range securityHeaders {
		if resp.Header.Get(header) == "" {
			t.Errorf("Expected security header %s to be present", header)
		}
	}
}

// TestAPIEndpoint tests the main API endpoint
func TestAPIEndpoint(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.cleanup()

	tests := []struct {
		name                string
		request             map[string]interface{}
		expectedStatus      int
		expectedStatusField string
	}{
		{
			name: "Valid echo request",
			request: map[string]interface{}{
				"message": "Hello World",
				"action":  "echo",
				"user_id": 123,
			},
			expectedStatus:      http.StatusOK,
			expectedStatusField: "success",
		},
		{
			name: "Valid greet request",
			request: map[string]interface{}{
				"message": "John",
				"action":  "greet",
				"user_id": 456,
			},
			expectedStatus:      http.StatusOK,
			expectedStatusField: "success",
		},
		{
			name: "Invalid action",
			request: map[string]interface{}{
				"message": "Hello",
				"action":  "invalid_action",
				"user_id": 123,
			},
			expectedStatus:      http.StatusNotFound,
			expectedStatusField: "error",
		},
		{
			name: "Missing required fields",
			request: map[string]interface{}{
				"action": "echo",
			},
			expectedStatus:      http.StatusBadRequest,
			expectedStatusField: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.request)
			req, err := http.NewRequest("POST", ts.baseURL+"/api", bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")

			resp, err := ts.client.Do(req)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response: %v", err)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(body, &response); err != nil {
				t.Fatalf("Failed to unmarshal response: %v", err)
			}

			// For error cases, check if we have a status field or if it's an error response
			if tt.expectedStatusField == "error" {
				if response["status"] != nil && response["status"] != "error" {
					t.Errorf("Expected status field 'error' or nil, got '%s'", response["status"])
				}
			} else {
				if response["status"] != tt.expectedStatusField {
					t.Errorf("Expected status field '%s', got '%s'", tt.expectedStatusField, response["status"])
				}
			}
		})
	}
}

// TestVersionEndpoint tests the version endpoint
func TestVersionEndpoint(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.cleanup()

	resp, err := ts.client.Get(ts.baseURL + "/version")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}

	// Check for required fields in data
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data field to be an object")
	}

	requiredFields := []string{"server", "version", "go_version", "os", "arch"}
	for _, field := range requiredFields {
		if data[field] == nil {
			t.Errorf("Expected data field '%s' to be present", field)
		}
	}
}

// TestMetricsEndpoint tests the metrics endpoint
func TestMetricsEndpoint(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.cleanup()

	resp, err := ts.client.Get(ts.baseURL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["status"] != "success" {
		t.Errorf("Expected status 'success', got '%s'", response["status"])
	}

	// Check for required fields in data
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data field to be an object")
	}

	requiredFields := []string{"memory", "runtime"}
	for _, field := range requiredFields {
		if data[field] == nil {
			t.Errorf("Expected data field '%s' to be present", field)
		}
	}
}

// TestCORSEndpoint tests CORS functionality
func TestCORSEndpoint(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.cleanup()

	// Test OPTIONS request
	req, err := http.NewRequest("OPTIONS", ts.baseURL+"/api", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Origin", "https://example.com")

	resp, err := ts.client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check CORS headers
	corsHeaders := []string{
		"Access-Control-Allow-Origin",
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
	}

	for _, header := range corsHeaders {
		if resp.Header.Get(header) == "" {
			t.Errorf("Expected CORS header %s to be present", header)
		}
	}
}

// TestRequestSizeLimit tests request size limiting
func TestRequestSizeLimit(t *testing.T) {
	ts := NewTestServer(t)
	defer ts.cleanup()

	// Create a large request body (2MB)
	largeBody := make([]byte, 2*1024*1024)
	for i := range largeBody {
		largeBody[i] = 'x'
	}

	request := map[string]interface{}{
		"message": string(largeBody),
		"action":  "echo",
	}

	jsonBody, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", ts.baseURL+"/api", bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := ts.client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Should be rejected due to size limit
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for large request, got %d", resp.StatusCode)
	}
}

// cleanup stops the test server
func (ts *TestServer) cleanup() {
	// In a real implementation, you'd stop the server
	// For now, we'll just log that cleanup was called
	fmt.Println("Test server cleanup called")
}
