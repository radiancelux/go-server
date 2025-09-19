# Go Server API Documentation

A comprehensive guide to the Go Server API endpoints, request/response formats, and usage examples.

## Table of Contents

- [Overview](#overview)
- [Base URL](#base-url)
- [Authentication](#authentication)
- [Response Format](#response-format)
- [Error Handling](#error-handling)
- [Endpoints](#endpoints)
  - [Health Check](#health-check)
  - [Version Information](#version-information)
  - [System Metrics](#system-metrics)
  - [Server Configuration](#server-configuration)
  - [Detailed Status](#detailed-status)
  - [Main API](#main-api)
- [API Actions](#api-actions)
- [Examples](#examples)
- [Rate Limiting](#rate-limiting)
- [Troubleshooting](#troubleshooting)

## Overview

The Go Server API provides a RESTful interface for interacting with a modular Go server. The API supports both direct HTTP endpoints and a unified JSON API endpoint for all operations.

## Base URL

```
http://localhost:8080
```

## Authentication

Currently, no authentication is required. All endpoints are publicly accessible.

## Response Format

All API responses follow a consistent JSON structure:

```json
{
  "status": "success|error",
  "message": "Human-readable message",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    // Response-specific data (optional)
  }
}
```

### Response Fields

| Field       | Type   | Description                                           |
| ----------- | ------ | ----------------------------------------------------- |
| `status`    | string | Response status: `"success"` or `"error"`             |
| `message`   | string | Human-readable message describing the result          |
| `timestamp` | string | ISO 8601 timestamp of when the response was generated |
| `data`      | object | Optional response data (varies by endpoint)           |

## Error Handling

The API uses standard HTTP status codes and provides detailed error information:

| Status Code | Description                                     |
| ----------- | ----------------------------------------------- |
| 200         | Success                                         |
| 400         | Bad Request - Invalid input or validation error |
| 405         | Method Not Allowed - Wrong HTTP method          |
| 500         | Internal Server Error - Server-side error       |

### Error Response Example

```json
{
  "status": "error",
  "message": "message is required",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Endpoints

### Health Check

**Endpoint:** `GET /health`

**Description:** Basic health check to verify server is running.

**Response:**

```json
{
  "status": "healthy",
  "message": "Server is running",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

**Example:**

```bash
curl http://localhost:8080/health
```

### Version Information

**Endpoint:** `GET /version`

**Description:** Returns server version and runtime information.

**Response:**

```json
{
  "status": "success",
  "message": "Version information",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "server": "go-server",
    "version": "1.0.0",
    "go_version": "go1.21.0",
    "os": "linux",
    "arch": "amd64",
    "num_cpu": 4
  }
}
```

**Example:**

```bash
curl http://localhost:8080/version
```

### System Metrics

**Endpoint:** `GET /metrics`

**Description:** Returns system metrics including memory usage and runtime statistics.

**Response:**

```json
{
  "status": "success",
  "message": "System metrics",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "memory": {
      "alloc_mb": 2,
      "total_alloc_mb": 10,
      "sys_mb": 8,
      "num_gc": 5,
      "gc_cpu_fraction": 0.001
    },
    "runtime": {
      "goroutines": 12,
      "cpus": 4
    },
    "timestamp": 1704110400
  }
}
```

**Example:**

```bash
curl http://localhost:8080/metrics
```

### Server Configuration

**Endpoint:** `GET /config`

**Description:** Returns server configuration and environment settings.

**Response:**

```json
{
  "status": "success",
  "message": "Server configuration",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "server": {
      "port": "8080",
      "host": "localhost"
    },
    "environment": {
      "go_env": "production",
      "port": "8080",
      "log_level": "info"
    },
    "features": {
      "graceful_shutdown": true,
      "request_validation": true,
      "structured_logging": true,
      "metrics": true
    }
  }
}
```

**Example:**

```bash
curl http://localhost:8080/config
```

### Detailed Status

**Endpoint:** `GET /status`

**Description:** Returns comprehensive server status including system information and metrics.

**Response:**

```json
{
  "status": "success",
  "message": "Detailed server status",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "status": "healthy",
    "server": {
      "name": "go-server",
      "version": "1.0.0",
      "port": "8080",
      "uptime": "running"
    },
    "system": {
      "go_version": "go1.21.0",
      "os": "linux",
      "arch": "amd64",
      "goroutines": 12,
      "cpus": 4
    },
    "memory": {
      "alloc_mb": 2,
      "total_alloc_mb": 10,
      "sys_mb": 8,
      "num_gc": 5
    },
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

**Example:**

```bash
curl http://localhost:8080/status
```

### Main API

**Endpoint:** `POST /api`

**Description:** Unified API endpoint for all server operations. Accepts JSON requests with different actions.

**Request Body:**

```json
{
  "message": "string",
  "user_id": 123,
  "action": "string"
}
```

**Request Fields:**

| Field     | Type    | Required | Description              |
| --------- | ------- | -------- | ------------------------ |
| `message` | string  | Yes      | The message content      |
| `user_id` | integer | No       | Optional user identifier |
| `action`  | string  | Yes      | The action to perform    |

**Example:**

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello World",
    "user_id": 123,
    "action": "echo"
  }'
```

## API Actions

The main API endpoint supports the following actions:

### Echo Action

**Action:** `echo`

**Description:** Echoes back the provided message.

**Request:**

```json
{
  "message": "Hello World",
  "action": "echo"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Echo successful",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "echoed_message": "Hello World"
  }
}
```

### Greet Action

**Action:** `greet`

**Description:** Creates a personalized greeting message.

**Request:**

```json
{
  "message": "Nice to meet you",
  "user_id": 123,
  "action": "greet"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Greeting generated",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "greeting": "Hello User 123! You said: Nice to meet you"
  }
}
```

### Info Action

**Action:** `info`

**Description:** Returns basic server information.

**Request:**

```json
{
  "message": "Tell me about the server",
  "action": "info"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Server information",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "server": "go-server",
    "version": "1.0.0",
    "port": "8080",
    "user_input": "Tell me about the server"
  }
}
```

### Version Action

**Action:** `version`

**Description:** Returns detailed version and runtime information.

**Request:**

```json
{
  "message": "version",
  "action": "version"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Version information",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "server": "go-server",
    "version": "1.0.0",
    "go_version": "go1.21.0",
    "os": "linux",
    "arch": "amd64",
    "num_cpu": 4
  }
}
```

### Metrics Action

**Action:** `metrics`

**Description:** Returns system metrics and performance data.

**Request:**

```json
{
  "message": "metrics",
  "action": "metrics"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "System metrics",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "memory": {
      "alloc_mb": 2,
      "total_alloc_mb": 10,
      "sys_mb": 8,
      "num_gc": 5,
      "gc_cpu_fraction": 0.001
    },
    "runtime": {
      "goroutines": 12,
      "cpus": 4
    },
    "timestamp": 1704110400
  }
}
```

### Config Action

**Action:** `config`

**Description:** Returns server configuration and environment settings.

**Request:**

```json
{
  "message": "config",
  "action": "config"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Server configuration",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "server": {
      "port": "8080",
      "host": "localhost"
    },
    "environment": {
      "go_env": "production",
      "port": "8080",
      "log_level": "info"
    },
    "features": {
      "graceful_shutdown": true,
      "request_validation": true,
      "structured_logging": true,
      "metrics": true
    }
  }
}
```

### Status Action

**Action:** `status`

**Description:** Returns comprehensive server status and system information.

**Request:**

```json
{
  "message": "status",
  "action": "status"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Detailed server status",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    "status": "healthy",
    "server": {
      "name": "go-server",
      "version": "1.0.0",
      "port": "8080",
      "uptime": "running"
    },
    "system": {
      "go_version": "go1.21.0",
      "os": "linux",
      "arch": "amd64",
      "goroutines": 12,
      "cpus": 4
    },
    "memory": {
      "alloc_mb": 2,
      "total_alloc_mb": 10,
      "sys_mb": 8,
      "num_gc": 5
    },
    "timestamp": "2024-01-01T12:00:00Z"
  }
}
```

## Examples

### Basic Health Check

```bash
# Check if server is running
curl http://localhost:8080/health
```

### Get Server Version

```bash
# Get version information
curl http://localhost:8080/version
```

### Echo a Message

```bash
# Echo a message through the API
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello World",
    "action": "echo"
  }'
```

### Create a Greeting

```bash
# Create a personalized greeting
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Nice to meet you",
    "user_id": 123,
    "action": "greet"
  }'
```

### Get System Metrics

```bash
# Get system metrics
curl http://localhost:8080/metrics
```

### Get Server Configuration

```bash
# Get server configuration
curl http://localhost:8080/config
```

### Get Detailed Status

```bash
# Get comprehensive server status
curl http://localhost:8080/status
```

## Rate Limiting

Currently, no rate limiting is implemented. All endpoints are available without restrictions.

## Troubleshooting

### Common Issues

#### 1. Server Not Responding

**Problem:** `curl: (7) Failed to connect to localhost port 8080: Connection refused`

**Solution:** Ensure the server is running:

```bash
go run main.go
```

#### 2. Invalid JSON Format

**Problem:** `400 Bad Request` with message "Invalid JSON format"

**Solution:** Ensure your JSON is properly formatted:

```bash
# Correct
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{"message": "test", "action": "echo"}'

# Incorrect (missing quotes)
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{message: test, action: echo}'
```

#### 3. Missing Required Fields

**Problem:** `400 Bad Request` with message "message is required"

**Solution:** Include all required fields:

```json
{
  "message": "required field",
  "action": "echo"
}
```

#### 4. Unknown Action

**Problem:** `400 Bad Request` with message "Unknown action 'invalid'"

**Solution:** Use a valid action from the supported list:

- `echo`
- `greet`
- `info`
- `version`
- `metrics`
- `config`
- `status`

### Debug Mode

To enable debug logging, set the `LOG_LEVEL` environment variable:

```bash
LOG_LEVEL=debug go run main.go
```

### Port Configuration

To run the server on a different port, set the `PORT` environment variable:

```bash
PORT=3000 go run main.go
```

## Support

For additional support or questions about the API, please refer to the project documentation or create an issue in the repository.
