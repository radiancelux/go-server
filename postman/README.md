# Postman Integration for Go Server

This directory contains Postman collections and environment files for testing the Go Server API.

## Files

- `Go-Server-API.postman_collection.json` - Complete API collection with all endpoints and test scenarios
- `Go-Server-Environment.postman_environment.json` - Environment variables for testing
- `README.md` - This documentation file

## Collection Structure

### Health & Status

- **Health Check** (`GET /health`) - Basic health monitoring
- **Server Status** (`GET /status`) - Detailed server information

### API Information

- **API Version** (`GET /version`) - Version information
- **API Configuration** (`GET /config`) - Server configuration
- **API Documentation** (`GET /docs`) - HTML documentation

### Metrics & Monitoring

- **Server Metrics** (`GET /metrics`) - Performance metrics

### API Actions

- **Echo Action** (`POST /api`) - Echo back messages
- **Greet Action** (`POST /api`) - Greet users
- **Info Action** (`POST /api`) - Get server info

### Security Tests

- **XSS Protection Test** - Test input sanitization
- **Rate Limiting Test** - Test rate limiting functionality
- **CORS Preflight Test** - Test CORS headers
- **Invalid JSON Test** - Test error handling
- **Missing Required Fields Test** - Test validation

### Performance Tests

- **Load Test - Health Endpoint** - Performance testing
- **Load Test - API Endpoint** - API performance testing

## Prerequisites

### Install Newman (Command Line Runner)

```bash
npm install -g newman
```

### Install Newman HTML Reporter (Optional)

```bash
npm install -g newman-reporter-html
```

## Usage

### 1. Manual Testing with Postman App

1. Import `Go-Server-API.postman_collection.json` into Postman
2. Import `Go-Server-Environment.postman_environment.json` as environment
3. Start the Go server: `go run main.go`
4. Run the collection or individual requests

### 2. Command Line Testing with Newman

#### Basic Test Run

```bash
# Run all tests
newman run postman/Go-Server-API.postman_collection.json -e postman/Go-Server-Environment.postman_environment.json

# Run with HTML report
newman run postman/Go-Server-API.postman_collection.json \
  -e postman/Go-Server-Environment.postman_environment.json \
  --reporters cli,html \
  --reporter-html-export test-results/postman-report.html
```

#### Using Go Test Runner

```bash
# Run Postman tests only
go run ./cmd/test -type postman -v

# Run all tests including Postman
go run ./cmd/test -type all -v

# Using Makefile
make test-postman
make test-all
```

### 3. Docker Testing

```bash
# Run tests in Docker (includes Postman)
make test-docker
```

## Environment Variables

The environment file contains these variables:

- `base_url` - Server base URL (default: http://localhost:8080)
- `api_version` - API version (default: v1)
- `test_user_id` - Test user ID (default: 123)
- `test_message` - Test message (default: "Hello from Postman!")
- `xss_test_message` - XSS test payload
- `rate_limit_test_message` - Rate limit test message

## Test Scenarios

### Security Testing

- **XSS Protection**: Tests that malicious scripts are properly sanitized
- **Rate Limiting**: Verifies rate limiting headers and behavior
- **CORS**: Tests cross-origin request handling
- **Input Validation**: Tests various invalid input scenarios

### Performance Testing

- **Response Time**: Ensures responses are within acceptable limits
- **Load Testing**: Tests server under multiple concurrent requests

### Functional Testing

- **API Endpoints**: Tests all available endpoints
- **Error Handling**: Tests various error scenarios
- **Data Validation**: Tests request/response data structures

## Reports

Test runs generate several types of reports:

- **CLI Output**: Real-time test results in terminal
- **JSON Report**: Machine-readable test results
- **HTML Report**: Human-readable test results with detailed information

## Customization

### Adding New Tests

1. Open the collection in Postman
2. Add new requests to appropriate folders
3. Add test scripts using JavaScript
4. Export the updated collection
5. Update this README if needed

### Modifying Environment

1. Edit `Go-Server-Environment.postman_environment.json`
2. Add or modify variables as needed
3. Test with the updated environment

### Integration with CI/CD

The Postman tests are integrated into the Go test runner and can be run as part of CI/CD pipelines:

```yaml
# GitHub Actions example
- name: Run Postman Tests
  run: |
    npm install -g newman
    go run ./cmd/test -type postman -v
```

## Troubleshooting

### Newman Not Found

```bash
# Install Newman globally
npm install -g newman

# Verify installation
newman --version
```

### Server Not Running

```bash
# Start the server
go run main.go

# Or with Docker
docker-compose up go-server
```

### Permission Issues (Linux/Mac)

```bash
# Make scripts executable
chmod +x scripts/run-postman-tests.sh
```

### Windows Execution Policy

```powershell
# Allow script execution
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
```

## Best Practices

1. **Always test locally first** before running in CI/CD
2. **Keep tests independent** - each test should work standalone
3. **Use environment variables** for configuration
4. **Add meaningful test names** and descriptions
5. **Include both positive and negative test cases**
6. **Test security scenarios** thoroughly
7. **Monitor performance** with response time assertions
8. **Keep collections updated** with API changes
