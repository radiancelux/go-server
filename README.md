# 🚀 Go Server - Production-Ready Backend Foundation

A **complete, modular, production-ready Go server** with database integration, authentication, security features, and comprehensive testing. Perfect as a foundation for web applications and APIs.

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org/)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED.svg)](https://www.docker.com/)
[![Security](https://img.shields.io/badge/Security-Scanned-green.svg)](https://github.com/golang/vuln)
[![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen.svg)](https://github.com/radiancelux/go-server)

## ✨ Features

### 🏗️ **Modular Architecture**
- **96% reduction** in large files through complete modularization
- **Clean separation** of concerns across all layers
- **Repository pattern** for data access
- **Service layer** for business logic
- **Interface-driven design** for testability

### 🗄️ **Database Integration**
- **PostgreSQL** - Primary production database with connection pooling
- **Redis** - Caching and session storage
- **SQLite** - Development and testing
- **GORM ORM** - Clean data access layer
- **Migration system** - Database schema management
- **Graceful fallback** - Server continues without databases

### 🔐 **Authentication & Security**
- **JWT tokens** - Secure authentication
- **User management** - Registration, login, profile management
- **Session management** - Secure session handling
- **Password security** - bcrypt hashing
- **Role-based access** - Admin and user roles
- **Input validation** - Comprehensive data sanitization
- **Rate limiting** - Per-IP request limiting
- **CORS support** - Cross-origin resource sharing

### 🧪 **Testing & Quality**
- **Unit tests** - All packages tested
- **Integration tests** - End-to-end testing
- **Performance tests** - Load testing capabilities
- **Postman integration** - API testing and documentation
- **Security scanning** - govulncheck integration
- **Go test runner** - Custom testing orchestration

### 🐳 **Production Ready**
- **Docker support** - Containerized deployment
- **CI/CD pipeline** - GitHub Actions automation
- **Structured logging** - Request tracking and debugging
- **Graceful shutdown** - Clean server termination
- **Health checks** - Container and endpoint monitoring
- **Interactive docs** - Auto-generated API documentation

## 🚀 Quick Start

### Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/radiancelux/go-server.git
cd go-server

# Start the server
docker-compose up -d

# Check status
docker ps

# View logs
docker logs go-server-go-server-1

# Stop server
docker-compose down
```

### Local Development

```bash
# Clone the repository
git clone https://github.com/radiancelux/go-server.git
cd go-server

# Install dependencies
go mod tidy

# Run the server
go run ./main.go

# Run tests
go test ./...

# Run security scan
go run ./cmd/security

# Test database integration
go run ./cmd/test-db
```

## 📊 Project Structure

```
go-server/
├── cmd/                    # Command-line applications
│   ├── main.go            # Main server entry point
│   ├── security/          # Security scanning tool
│   └── test-db/           # Database testing tool
├── internal/              # Internal packages
│   ├── auth/              # Authentication system
│   │   ├── jwt.go         # JWT token management
│   │   ├── service.go     # Auth service coordinator
│   │   ├── login_service.go
│   │   ├── registration_service.go
│   │   ├── session_service.go
│   │   └── types.go       # Auth data structures
│   ├── config/            # Configuration management
│   ├── database/          # Database layer
│   │   ├── config.go      # Database configuration
│   │   ├── connection.go  # Connection management
│   │   ├── migrate.go     # Database migrations
│   │   ├── models/        # Data models
│   │   └── repositories/  # Data access layer
│   ├── docs/              # API documentation
│   │   ├── generator.go   # Documentation generation
│   │   ├── parser.go      # Postman collection parsing
│   │   ├── templates.go   # HTML templates
│   │   └── types.go       # Documentation types
│   ├── errors/            # Error handling
│   ├── handlers/          # HTTP handlers
│   ├── interfaces/        # Interface definitions
│   ├── logger/            # Logging system
│   ├── middleware/        # HTTP middleware
│   ├── models/            # Data models
│   ├── security/          # Security utilities
│   │   ├── validator.go   # Main validator
│   │   ├── field_validator.go
│   │   ├── http_validator.go
│   │   └── validation_utils.go
│   ├── server/            # Server implementation
│   │   ├── server.go      # Main server
│   │   ├── types.go       # Server types
│   │   ├── handlers.go    # HTTP handlers
│   │   ├── routes.go      # Route setup
│   │   ├── lifecycle.go   # Server lifecycle
│   │   └── utils.go       # Server utilities
│   ├── services/          # Business logic
│   └── testrunner/        # Testing framework
├── migrations/            # Database migrations
├── postman/              # Postman collections
├── test/                 # Test files
└── docs/                 # Documentation
```

## 🌐 API Endpoints

### Core Endpoints
- `GET /health` - Health check
- `GET /version` - Version information
- `GET /docs` - Interactive API documentation
- `GET /metrics` - Server metrics
- `GET /config` - Configuration info
- `GET /status` - Server status

### API Actions
- `POST /api` - Main API endpoint with actions:
  - `echo` - Echoes back the message
  - `greet` - Creates a personalized greeting
  - `info` - Returns server information

### Authentication Endpoints
- `POST /auth/register` - User registration
- `POST /auth/login` - User login
- `POST /auth/logout` - User logout
- `POST /auth/refresh` - Token refresh

## 🗄️ Database Configuration

### Environment Variables

```bash
# PostgreSQL
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=password
POSTGRES_DB=go_server
POSTGRES_SSLMODE=disable

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Connection Settings
MAX_CONNECTIONS=100
MAX_IDLE_CONNS=10
CONN_MAX_LIFETIME=1h
CONN_MAX_IDLE_TIME=30m
```

### Database Support
- **PostgreSQL** - Primary production database
- **Redis** - Caching and session storage
- **SQLite** - Development and testing
- **Graceful fallback** - Server works without databases

## 🧪 Testing

### Run All Tests
```bash
# Run unit tests
go test ./...

# Run integration tests
go test ./test/...

# Run performance tests
go test -bench=. ./test/...

# Run security scan
go run ./cmd/security

# Test database integration
go run ./cmd/test-db
```

### Test Coverage
- **Unit Tests** - All packages tested
- **Integration Tests** - End-to-end testing
- **Performance Tests** - Load testing
- **Security Tests** - Vulnerability scanning
- **Postman Tests** - API testing

## 📚 Documentation

### Interactive Documentation
- **[Live API Docs](http://localhost:8080/docs)** - Interactive HTML documentation
- **Postman Collection** - Import into Postman for testing
- **OpenAPI Specification** - Machine-readable API spec

### Documentation Features
- **Live API Testing** - Test endpoints directly from browser
- **Copy/Download Collection** - One-click Postman export
- **Always Up-to-Date** - Generated from actual test collection
- **Comprehensive Coverage** - All endpoints and scenarios

## 🔧 Development

### Adding New Features

1. **New Handlers** - Add to `internal/handlers/`
2. **New Services** - Add to `internal/services/`
3. **New Models** - Add to `internal/models/`
4. **New Repositories** - Add to `internal/database/repositories/`
5. **New Tests** - Add to appropriate test files

### Key Design Principles

- **Interfaces as Contracts** - Define clear API boundaries
- **Modular Structure** - Separate concerns into focused packages
- **Dependency Injection** - Clean dependency management
- **Repository Pattern** - Clean data access layer
- **Service Layer** - Business logic separation
- **Graceful Shutdown** - Production-ready lifecycle management

## 🚀 Deployment

### Docker Deployment
```bash
# Build and run
docker-compose up -d

# Check status
docker ps

# View logs
docker logs go-server-go-server-1

# Stop
docker-compose down
```

### Production Considerations
- Set up PostgreSQL and Redis databases
- Configure environment variables
- Set up monitoring and logging
- Configure load balancing if needed
- Set up SSL/TLS certificates

## 📊 Performance

### Benchmarks
- **Response Time** - Sub-millisecond for simple endpoints
- **Throughput** - High concurrent request handling
- **Memory Usage** - Efficient memory management
- **Database** - Connection pooling and optimization

### Monitoring
- **Health Checks** - Container and endpoint monitoring
- **Metrics** - Server and application metrics
- **Logging** - Structured logging with request tracking
- **Security** - Vulnerability scanning and monitoring

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run the test suite
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🎯 Roadmap

- [ ] WebSocket support
- [ ] GraphQL API
- [ ] Microservices architecture
- [ ] Kubernetes deployment
- [ ] Advanced monitoring
- [ ] API versioning

## 🙏 Acknowledgments

- Built with Go standard library and community packages
- Inspired by clean architecture principles
- Designed for production use and scalability

---

**Repository**: https://github.com/radiancelux/go-server  
**Status**: ✅ Production Ready  
**Last Updated**: September 2024

*Ready to power your next web application!* 🚀