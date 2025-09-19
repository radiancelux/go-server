package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"go-server/internal/config"
	"go-server/internal/server"
)

// BenchmarkServer represents a benchmark server instance
type BenchmarkServer struct {
	server *server.Server
	client *http.Client
	baseURL string
}

// NewBenchmarkServer creates a new benchmark server
func NewBenchmarkServer(t *testing.B) *BenchmarkServer {
	// Use a random port to avoid conflicts (use 7xxx range for benchmarks)
	port := fmt.Sprintf("7%03d", (time.Now().UnixNano()/1000)%1000)
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: port,
		},
		Security: config.SecurityConfig{
			MaxRequestSize: 1024 * 1024,
			RateLimitRPS:   10000, // Very high limit for benchmarks
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
	
	return &BenchmarkServer{
		server:  srv,
		client:  &http.Client{Timeout: 30 * time.Second},
		baseURL: fmt.Sprintf("http://localhost:%s", port),
	}
}

// BenchmarkHealthEndpoint benchmarks the health endpoint
func BenchmarkHealthEndpoint(b *testing.B) {
	bs := NewBenchmarkServer(b)
	defer bs.cleanup()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		resp, err := bs.client.Get(bs.baseURL + "/health")
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkAPIEndpoint benchmarks the API endpoint
func BenchmarkAPIEndpoint(b *testing.B) {
	bs := NewBenchmarkServer(b)
	defer bs.cleanup()
	
	request := map[string]interface{}{
		"message": "Hello World",
		"action":  "echo",
		"user_id": 123,
	}
	
	jsonBody, _ := json.Marshal(request)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("POST", bs.baseURL+"/api", bytes.NewBuffer(jsonBody))
		if err != nil {
			b.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := bs.client.Do(req)
		if err != nil {
			b.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
}

// BenchmarkConcurrentRequests benchmarks concurrent requests
func BenchmarkConcurrentRequests(b *testing.B) {
	bs := NewBenchmarkServer(b)
	defer bs.cleanup()
	
	request := map[string]interface{}{
		"message": "Hello World",
		"action":  "echo",
		"user_id": 123,
	}
	
	jsonBody, _ := json.Marshal(request)
	
	b.ResetTimer()
	
	var wg sync.WaitGroup
	concurrency := 10
	
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			req, err := http.NewRequest("POST", bs.baseURL+"/api", bytes.NewBuffer(jsonBody))
			if err != nil {
				b.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/json")
			
			resp, err := bs.client.Do(req)
			if err != nil {
				b.Fatalf("Failed to make request: %v", err)
			}
			resp.Body.Close()
		}()
		
		if i%concurrency == 0 {
			wg.Wait()
		}
	}
	
	wg.Wait()
}

// TestLoadTest performs a load test
func TestLoadTest(t *testing.T) {
	bs := NewBenchmarkServer(&testing.B{})
	defer bs.cleanup()
	
	request := map[string]interface{}{
		"message": "Hello World",
		"action":  "echo",
		"user_id": 123,
	}
	
	jsonBody, _ := json.Marshal(request)
	
	// Test parameters
	concurrency := 10
	duration := 5 * time.Second
	requestsPerSecond := 100
	
	// Create a ticker for rate limiting
	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
	defer ticker.Stop()
	
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	
	var wg sync.WaitGroup
	requestCount := 0
	errorCount := 0
	successCount := 0
	
	// Start workers
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					req, err := http.NewRequest("POST", bs.baseURL+"/api", bytes.NewBuffer(jsonBody))
					if err != nil {
						errorCount++
						continue
					}
					req.Header.Set("Content-Type", "application/json")
					
					resp, err := bs.client.Do(req)
					if err != nil {
						errorCount++
						continue
					}
					
					if resp.StatusCode == http.StatusOK {
						successCount++
					} else {
						errorCount++
					}
					
					resp.Body.Close()
					requestCount++
				}
			}
		}()
	}
	
	// Wait for all workers to complete
	wg.Wait()
	
	// Calculate metrics
	actualRPS := float64(requestCount) / duration.Seconds()
	successRate := float64(successCount) / float64(requestCount) * 100
	errorRate := float64(errorCount) / float64(requestCount) * 100
	
	t.Logf("Load Test Results:")
	t.Logf("  Duration: %v", duration)
	t.Logf("  Concurrency: %d", concurrency)
	t.Logf("  Total Requests: %d", requestCount)
	t.Logf("  Successful Requests: %d", successCount)
	t.Logf("  Failed Requests: %d", errorCount)
	t.Logf("  Actual RPS: %.2f", actualRPS)
	t.Logf("  Success Rate: %.2f%%", successRate)
	t.Logf("  Error Rate: %.2f%%", errorRate)
	
	// Assertions
	if successRate < 95.0 {
		t.Errorf("Success rate too low: %.2f%% (expected >= 95%%)", successRate)
	}
	
	if actualRPS < float64(requestsPerSecond)*0.8 {
		t.Errorf("RPS too low: %.2f (expected >= %.2f)", actualRPS, float64(requestsPerSecond)*0.8)
	}
}

// TestMemoryUsage tests memory usage under load
func TestMemoryUsage(t *testing.T) {
	bs := NewBenchmarkServer(&testing.B{})
	defer bs.cleanup()
	
	request := map[string]interface{}{
		"message": "Hello World",
		"action":  "echo",
		"user_id": 123,
	}
	
	jsonBody, _ := json.Marshal(request)
	
	// Make many requests to test memory usage
	for i := 0; i < 1000; i++ {
		req, err := http.NewRequest("POST", bs.baseURL+"/api", bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := bs.client.Do(req)
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		resp.Body.Close()
	}
	
	// Check memory usage via metrics endpoint
	resp, err := bs.client.Get(bs.baseURL + "/metrics")
	if err != nil {
		t.Fatalf("Failed to get metrics: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read metrics: %v", err)
	}
	
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal metrics: %v", err)
	}
	
	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected data field to be an object")
	}
	
	memory, ok := data["memory"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected memory field to be an object")
	}
	
	allocMB, ok := memory["alloc_mb"].(float64)
	if !ok {
		t.Fatal("Expected alloc_mb to be a number")
	}
	
	t.Logf("Memory usage after 1000 requests: %.2f MB", allocMB)
	
	// Assert memory usage is reasonable (less than 100MB)
	if allocMB > 100 {
		t.Errorf("Memory usage too high: %.2f MB (expected < 100 MB)", allocMB)
	}
}

// cleanup stops the benchmark server
func (bs *BenchmarkServer) cleanup() {
	// In a real implementation, you'd stop the server
	// For now, we'll just log that cleanup was called
	fmt.Println("Benchmark server cleanup called")
}
