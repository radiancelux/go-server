# Go Server Example Requests

This document provides examples of how to interact with the Go server.

## Starting the Server

```bash
go run main.go
```

The server will start on port 8080 by default. You can change the port by setting the `PORT` environment variable:

```bash
PORT=3000 go run main.go
```

## API Endpoints

### Health Check

```bash
curl http://localhost:8080/health
```

### Main API Endpoint

#### Echo Action

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Hello World",
    "action": "echo"
  }'
```

#### Greet Action

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Nice to meet you",
    "user_id": 123,
    "action": "greet"
  }'
```

#### Info Action

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{
    "message": "Tell me about the server",
    "action": "info"
  }'
```

## Expected Responses

All responses follow this structure:

```json
{
  "status": "success|error",
  "message": "Description of the result",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    // Additional data based on the action
  }
}
```

### Example Echo Response

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

### Example Greet Response

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

### Example Info Response

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

## Error Handling

The server returns appropriate HTTP status codes and error messages:

- `400 Bad Request` - Invalid JSON or unknown action
- `405 Method Not Allowed` - Wrong HTTP method
- `500 Internal Server Error` - Server-side errors

### Example Error Response

```json
{
  "status": "error",
  "message": "Invalid JSON format",
  "timestamp": "2024-01-01T12:00:00Z"
}
```
