package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Logging  LoggingConfig
	Security SecurityConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level  string
	Format string
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	MaxRequestSize int64
	RateLimitRPS   int
	RateLimitBurst int
	EnableCORS     bool
	CORSOrigins    []string

	// Input validation
	EnableInputValidation bool
	MaxStringLength       int
	MaxEmailLength        int

	// Security headers
	EnableSecurityHeaders bool
	ContentSecurityPolicy string
}

// Load loads configuration from environment variables with defaults
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     getDurationEnv("READ_TIMEOUT", 30*time.Second),
			WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 30*time.Second),
			IdleTimeout:     getDurationEnv("IDLE_TIMEOUT", 120*time.Second),
			ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
		Security: SecurityConfig{
			MaxRequestSize: getInt64Env("MAX_REQUEST_SIZE", 1024*1024), // 1MB
			RateLimitRPS:   getIntEnv("RATE_LIMIT_RPS", 100),
			RateLimitBurst: getIntEnv("RATE_LIMIT_BURST", 200),
			EnableCORS:     getBoolEnv("ENABLE_CORS", true),
			CORSOrigins:    getStringSliceEnv("CORS_ORIGINS", []string{"*"}),

			// Input validation
			EnableInputValidation: getBoolEnv("ENABLE_INPUT_VALIDATION", true),
			MaxStringLength:       getIntEnv("MAX_STRING_LENGTH", 1000),
			MaxEmailLength:        getIntEnv("MAX_EMAIL_LENGTH", 254),

			// Security headers
			EnableSecurityHeaders: getBoolEnv("ENABLE_SECURITY_HEADERS", true),
			ContentSecurityPolicy: getEnv("CONTENT_SECURITY_POLICY", "default-src 'self'"),
		},
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("server port cannot be empty")
	}

	if c.Server.ReadTimeout <= 0 {
		return fmt.Errorf("read timeout must be positive")
	}

	if c.Server.WriteTimeout <= 0 {
		return fmt.Errorf("write timeout must be positive")
	}

	if c.Server.IdleTimeout <= 0 {
		return fmt.Errorf("idle timeout must be positive")
	}

	if c.Server.ShutdownTimeout <= 0 {
		return fmt.Errorf("shutdown timeout must be positive")
	}

	if c.Security.MaxRequestSize <= 0 {
		return fmt.Errorf("max request size must be positive")
	}

	if c.Security.RateLimitRPS <= 0 {
		return fmt.Errorf("rate limit RPS must be positive")
	}

	if c.Security.RateLimitBurst <= 0 {
		return fmt.Errorf("rate limit burst must be positive")
	}

	return nil
}

// GetServerAddress returns the full server address
func (c *Config) GetServerAddress() string {
	return ":" + c.Server.Port
}

// Helper functions for environment variable parsing

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getInt64Env(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Simple comma-separated values parsing
		// In production, you might want more sophisticated parsing
		return []string{value}
	}
	return defaultValue
}
