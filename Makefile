# Go Server Makefile
# Provides easy commands for testing, building, and running the server

.PHONY: help test test-unit test-integration test-e2e test-performance test-coverage test-all lint build run clean

# Default target
help:
	@echo "Go Server - Available Commands"
	@echo "=============================="
	@echo ""
	@echo "Testing:"
	@echo "  test              - Run all tests"
	@echo "  test-unit         - Run unit tests only"
	@echo "  test-integration  - Run integration tests only"
	@echo "  test-e2e          - Run end-to-end tests only"
	@echo "  test-performance  - Run performance tests only"
	@echo "  test-coverage     - Run tests with coverage report"
	@echo "  test-all          - Run comprehensive test suite"
	@echo ""
	@echo "Development:"
	@echo "  lint              - Run linting checks"
	@echo "  build             - Build the server binary"
	@echo "  run               - Run the server"
	@echo "  clean             - Clean build artifacts"
	@echo ""
	@echo "Examples:"
	@echo "  make test-unit"
	@echo "  make test-coverage"
	@echo "  make run"

# Run all tests
test:
	@echo "🧪 Running all tests..."
	go test ./... -v

# Run unit tests only
test-unit:
	@echo "🧪 Running unit tests..."
	go test ./internal/... -v

# Run integration tests only
test-integration:
	@echo "🔗 Running integration tests..."
	go test ./test -v -run TestServer

# Run end-to-end tests only
test-e2e:
	@echo "🌐 Running end-to-end tests..."
	go test ./test -v -run TestHealthEndpoint -run TestAPIEndpoint -run TestVersionEndpoint -run TestMetricsEndpoint -run TestCORSEndpoint -run TestRequestSizeLimit

# Run performance tests only
test-performance:
	@echo "⚡ Running performance tests..."
	go test ./test -v -run TestLoadTest -run TestMemoryUsage

# Run tests with coverage
test-coverage:
	@echo "📈 Running tests with coverage..."
	go test ./... -cover -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

# Run comprehensive test suite
test-all:
	@echo "🚀 Running comprehensive test suite..."
	@chmod +x scripts/test.sh
	@./scripts/test.sh

# Run linting
lint:
	@echo "🔍 Running linting checks..."
	go vet ./...
	go fmt ./...

# Build the server
build:
	@echo "🔨 Building server..."
	go build -o bin/go-server main.go

# Run the server
run:
	@echo "🚀 Starting server..."
	go run main.go

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf test-results/

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod tidy
	go mod download

# Run security scan
security:
	@echo "🔒 Running security scan..."
	go list -json -deps ./... | grep -E '"(ImportPath|Imports)"' | grep -v 'go-server' | sort | uniq

# Run benchmarks
bench:
	@echo "📊 Running benchmarks..."
	go test ./test -bench=. -benchmem

# Generate test report
report:
	@echo "📋 Generating test report..."
	@chmod +x scripts/test.sh
	@./scripts/test.sh

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker build -t go-server .

docker-run:
	@echo "🐳 Running Docker container..."
	docker run -p 8080:8080 go-server

# Development server with hot reload
dev:
	@echo "🔥 Starting development server with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/cosmtrek/air@latest"; \
		go run main.go; \
	fi

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

# Check for updates
update:
	@echo "🔄 Checking for updates..."
	go list -u -m all

# Show test coverage in terminal
coverage-term:
	@echo "📊 Showing test coverage in terminal..."
	go test ./... -cover

# Run specific test
test-specific:
	@echo "🎯 Running specific test..."
	@read -p "Enter test name: " testname; \
	go test ./... -v -run $$testname

# Show help for specific command
help-test:
	@echo "Test Commands Help"
	@echo "=================="
	@echo ""
	@echo "test-unit:        Run unit tests for internal packages"
	@echo "test-integration: Run integration tests"
	@echo "test-e2e:         Run end-to-end tests"
	@echo "test-performance: Run performance and load tests"
	@echo "test-coverage:    Run tests with HTML coverage report"
	@echo "test-all:         Run comprehensive test suite with reporting"
	@echo ""
	@echo "Examples:"
	@echo "  make test-unit"
	@echo "  make test-coverage"
	@echo "  make test-specific"
