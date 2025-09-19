# Go Server

A modular Go server that responds to JSON requests using only the standard library.

## Project Structure

```
go-server/
├── main.go                 # Minimal main entry point
├── go.mod                  # Go module definition
├── api.yaml               # OpenAPI 3.0 specification
├── example_requests.md    # Usage examples
└── internal/              # Internal packages (not importable from outside)
    ├── interfaces/        # Interface definitions
    │   └── interfaces.go
    ├── models/            # Request/Response models
    │   ├── request.go
    │   └── response.go
    ├── handlers/          # Business logic handlers
    │   ├── echo.go
    │   ├── greet.go
    │   ├── info.go
    │   └── registry.go
    ├── logger/            # Logging implementation
    │   └── logger.go
    └── server/            # HTTP server implementation
        └── server.go
```

## Architecture

### Interfaces (API Surface)

- **`APIRequest`** - Contract for incoming requests
- **`APIResponse`** - Contract for outgoing responses
- **`Handler`** - Contract for API handlers
- **`Logger`** - Contract for logging

### Modular Design

- **Handlers** - Each handler implements the `Handler` interface
- **Registry** - Dynamic handler registration and retrieval
- **Models** - Clean request/response structures with validation
- **Server** - HTTP server with graceful shutdown

## Usage

### Start the Server

```bash
go run main.go
```

### Environment Variables

- `PORT` - Server port (default: 8080)

### API Endpoints

#### Health Check

```bash
curl http://localhost:8080/health
```

#### Main API

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello World", "action": "echo"}'
```

### Supported Actions

- `echo` - Echoes back the message
- `greet` - Creates a personalized greeting
- `info` - Returns server information
- `version` - Server version and runtime info
- `metrics` - System metrics and memory usage
- `config` - Server configuration
- `status` - Detailed server status

## Documentation

- **[Complete API Documentation](docs/API.md)** - Comprehensive API reference with examples
- **[Quick Reference](docs/QuickReference.md)** - Quick lookup for endpoints and examples
- **[Postman Collection](docs/Go-Server-API.postman_collection.json)** - Import into Postman for testing
- **[OpenAPI Specification](api.yaml)** - Machine-readable API specification

## Development

### Adding New Handlers

1. Create a new handler file in `internal/handlers/`
2. Implement the `Handler` interface
3. Register the handler in `internal/server/server.go`

### Key Design Principles

- **Interfaces as Contracts** - Define clear API boundaries
- **Modular Structure** - Separate concerns into focused packages
- **Dependency Injection** - Handlers receive dependencies through constructors
- **Standard Library Only** - No external dependencies
- **Graceful Shutdown** - Production-ready server lifecycle management

## API Testing

### Using curl

```bash
# Health check
curl http://localhost:8080/health

# Echo message
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello", "action": "echo"}'
```

### Using Postman

1. Import the [Postman collection](docs/Go-Server-API.postman_collection.json)
2. Set the `base_url` variable to `http://localhost:8080`
3. Run the requests to test the API
