# Go Server API - Quick Reference

## Base URL

```
http://localhost:8080
```

## Direct Endpoints

| Method | Endpoint   | Description          |
| ------ | ---------- | -------------------- |
| `GET`  | `/health`  | Health check         |
| `GET`  | `/version` | Version information  |
| `GET`  | `/metrics` | System metrics       |
| `GET`  | `/config`  | Server configuration |
| `GET`  | `/status`  | Detailed status      |
| `POST` | `/api`     | Main API endpoint    |

## API Actions (via POST /api)

| Action    | Description       | Required Fields     |
| --------- | ----------------- | ------------------- |
| `echo`    | Echo back message | `message`, `action` |
| `greet`   | Create greeting   | `message`, `action` |
| `info`    | Server info       | `message`, `action` |
| `version` | Version details   | `message`, `action` |
| `metrics` | System metrics    | `message`, `action` |
| `config`  | Configuration     | `message`, `action` |
| `status`  | Detailed status   | `message`, `action` |

## Quick Examples

### Health Check

```bash
curl http://localhost:8080/health
```

### Echo Message

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{"message": "Hello", "action": "echo"}'
```

### Get Version

```bash
curl http://localhost:8080/version
```

### Get Metrics

```bash
curl http://localhost:8080/metrics
```

### Greet User

```bash
curl -X POST http://localhost:8080/api \
  -H "Content-Type: application/json" \
  -d '{"message": "Hi", "user_id": 123, "action": "greet"}'
```

## Response Format

```json
{
  "status": "success|error",
  "message": "Description",
  "timestamp": "2024-01-01T12:00:00Z",
  "data": {
    /* optional data */
  }
}
```

## Error Codes

- `200` - Success
- `400` - Bad Request
- `405` - Method Not Allowed
- `500` - Internal Server Error

## Environment Variables

- `PORT` - Server port (default: 8080)
- `LOG_LEVEL` - Log level (debug, info, error)
