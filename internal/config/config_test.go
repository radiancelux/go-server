package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Test with default values
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test default values
	if cfg.Server.Port != "8080" {
		t.Errorf("Expected port 8080, got %s", cfg.Server.Port)
	}

	if cfg.Server.ReadTimeout != 30*time.Second {
		t.Errorf("Expected read timeout 30s, got %v", cfg.Server.ReadTimeout)
	}

	if cfg.Logging.Level != "info" {
		t.Errorf("Expected log level 'info', got %s", cfg.Logging.Level)
	}

	if cfg.Security.MaxRequestSize != 1024*1024 {
		t.Errorf("Expected max request size 1MB, got %d", cfg.Security.MaxRequestSize)
	}
}

func TestLoadWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("LOG_LEVEL", "debug")
	os.Setenv("MAX_REQUEST_SIZE", "2048")
	defer func() {
		os.Unsetenv("PORT")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("MAX_REQUEST_SIZE")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Test environment variable values
	if cfg.Server.Port != "9090" {
		t.Errorf("Expected port 9090, got %s", cfg.Server.Port)
	}

	if cfg.Logging.Level != "debug" {
		t.Errorf("Expected log level 'debug', got %s", cfg.Logging.Level)
	}

	if cfg.Security.MaxRequestSize != 2048 {
		t.Errorf("Expected max request size 2048, got %d", cfg.Security.MaxRequestSize)
	}
}

func TestValidate(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            "8080",
			ReadTimeout:     30 * time.Second,
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     120 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Security: SecurityConfig{
			MaxRequestSize: 1024 * 1024,
			RateLimitRPS:   100,
			RateLimitBurst: 200,
		},
	}

	if err := cfg.Validate(); err != nil {
		t.Errorf("Valid config should not return error: %v", err)
	}
}

func TestValidateInvalidConfig(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port:            "", // Invalid: empty port
			ReadTimeout:     -1, // Invalid: negative timeout
			WriteTimeout:    30 * time.Second,
			IdleTimeout:     120 * time.Second,
			ShutdownTimeout: 10 * time.Second,
		},
		Security: SecurityConfig{
			MaxRequestSize: 1024 * 1024,
			RateLimitRPS:   100,
			RateLimitBurst: 200,
		},
	}

	if err := cfg.Validate(); err == nil {
		t.Error("Invalid config should return error")
	}
}

func TestGetServerAddress(t *testing.T) {
	cfg := &Config{
		Server: ServerConfig{
			Port: "8080",
		},
	}

	expected := ":8080"
	if cfg.GetServerAddress() != expected {
		t.Errorf("Expected %s, got %s", expected, cfg.GetServerAddress())
	}
}
