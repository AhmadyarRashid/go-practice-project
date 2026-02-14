# Go Enterprise API

A production-ready, enterprise-level REST API built with Go, featuring best practices for authentication, database management, logging, error handling, and more.

## Table of Contents

- [Features](#features)
- [Project Structure](#project-structure)
- [Architecture](#architecture)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Database](#database)
- [Middleware](#middleware)
- [Error Handling](#error-handling)
- [Logging](#logging)
- [Testing](#testing)
- [Deployment](#deployment)
- [Best Practices](#best-practices)

## Features

- **Clean Architecture**: Separation of concerns with handlers, services, repositories
- **JWT Authentication**: Secure token-based authentication with refresh tokens
- **Role-Based Access Control**: User, Moderator, Admin roles
- **GORM ORM**: Database abstraction with migrations and relationships
- **Structured Logging**: JSON logging with Zap
- **Rate Limiting**: Protect endpoints from abuse
- **CORS Support**: Configurable cross-origin resource sharing
- **Graceful Shutdown**: Proper server shutdown handling
- **Docker Support**: Ready for containerization
- **Hot Reload**: Development with Air
- **Input Validation**: Request validation with custom validators
- **Error Handling**: Centralized error handling with custom error types
- **Health Checks**: Liveness and readiness endpoints

## Project Structure

```
go-enterprise-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection and utilities
│   ├── handlers/
│   │   ├── auth_handler.go      # Authentication handlers
│   │   ├── user_handler.go      # User CRUD handlers
│   │   ├── post_handler.go      # Post CRUD handlers
│   │   └── health_handler.go    # Health check handlers
│   ├── middleware/
│   │   ├── auth.go              # Authentication middleware
│   │   ├── cors.go              # CORS middleware
│   │   ├── logger.go            # Request logging middleware
│   │   ├── ratelimit.go         # Rate limiting middleware
│   │   └── recovery.go          # Panic recovery middleware
│   ├── models/
│   │   ├── base.go              # Base model with UUID and timestamps
│   │   ├── user.go              # User model
│   │   └── post.go              # Post and Tag models
│   ├── repository/
│   │   ├── repository.go        # Generic repository interface
│   │   ├── user_repository.go   # User repository
│   │   └── post_repository.go   # Post repository
│   ├── routes/
│   │   └── routes.go            # Route definitions
│   └── services/
│       ├── auth_service.go      # Authentication service
│       ├── user_service.go      # User service
│       └── post_service.go      # Post service
├── pkg/
│   ├── errors/
│   │   └── errors.go            # Custom error types
│   ├── logger/
│   │   └── logger.go            # Logger utilities
│   ├── response/
│   │   └── response.go          # API response helpers
│   └── validator/
│       └── validator.go         # Input validation
├── .air.toml                     # Air hot reload config
├── .env.example                  # Environment variables template
├── .gitignore                    # Git ignore rules
├── docker-compose.yml            # Docker Compose config
├── Dockerfile                    # Docker build file
├── go.mod                        # Go modules
├── Makefile                      # Build automation
└── README.md                     # This file
```

## Architecture

This project follows **Clean Architecture** principles:

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Layer                              │
│  (Routes, Middleware, Handlers)                             │
├─────────────────────────────────────────────────────────────┤
│                    Service Layer                             │
│  (Business Logic, Use Cases)                                │
├─────────────────────────────────────────────────────────────┤
│                   Repository Layer                           │
│  (Data Access, Database Operations)                         │
├─────────────────────────────────────────────────────────────┤
│                     Data Layer                               │
│  (Models, Database)                                         │
└─────────────────────────────────────────────────────────────┘
```

### Layer Responsibilities

| Layer | Responsibility |
|-------|---------------|
| **Handlers** | HTTP request/response handling, input validation |
| **Services** | Business logic, orchestration |
| **Repositories** | Data access, database queries |
| **Models** | Data structures, domain entities |

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Make (optional but recommended)
- Docker (optional, for containerization)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/go-enterprise-api.git
   cd go-enterprise-api
   ```

2. **Copy environment file**
   ```bash
   cp .env.example .env
   ```

3. **Edit the .env file** and update the values (especially `JWT_SECRET`)

4. **Install dependencies**
   ```bash
   make deps
   # or
   go mod download
   ```

5. **Run the application**
   ```bash
   make run
   # or
   go run cmd/api/main.go
   ```

6. **For development with hot reload**
   ```bash
   make dev
   ```

### Using Docker

```bash
# Build and run with Docker Compose
docker-compose up -d

# Or build manually
make docker-build
make docker-run
```

## Configuration

Configuration is managed through environment variables. See `.env.example` for all options.

### Key Configuration Options

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Environment (development/production) | development |
| `APP_PORT` | Server port | 8080 |
| `DB_DRIVER` | Database driver (postgres/sqlite) | sqlite |
| `DB_HOST` | Database host | localhost |
| `DB_PORT` | Database port | 5432 |
| `DB_NAME` | Database name | enterprise.db |
| `JWT_SECRET` | JWT signing secret (min 32 chars) | *required* |
| `JWT_EXPIRY_HOURS` | Access token expiry | 24 |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | debug |

## API Endpoints

### Health Checks
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/health` | Basic health check |
| GET | `/api/v1/health/ready` | Readiness check |
| GET | `/api/v1/health/live` | Liveness check |

### Authentication
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | Register new user | No |
| POST | `/api/v1/auth/login` | Login user | No |
| POST | `/api/v1/auth/logout` | Logout user | Yes |
| POST | `/api/v1/auth/refresh` | Refresh tokens | No |
| GET | `/api/v1/auth/me` | Get current user | Yes |
| POST | `/api/v1/auth/change-password` | Change password | Yes |

### Users
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/users` | List all users | Yes |
| GET | `/api/v1/users/:id` | Get user by ID | Yes |
| PUT | `/api/v1/users/:id` | Update user | Yes |
| DELETE | `/api/v1/users/:id` | Delete user | Admin |
| GET | `/api/v1/users/search` | Search users | Yes |
| PATCH | `/api/v1/users/:id/status` | Update status | Admin |
| PATCH | `/api/v1/users/:id/role` | Update role | Admin |

### Posts
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/posts` | List posts | No* |
| POST | `/api/v1/posts` | Create post | Yes |
| GET | `/api/v1/posts/:id` | Get post | No* |
| PUT | `/api/v1/posts/:id` | Update post | Yes |
| DELETE | `/api/v1/posts/:id` | Delete post | Yes |
| GET | `/api/v1/posts/my` | Get my posts | Yes |
| GET | `/api/v1/posts/search` | Search posts | No |
| GET | `/api/v1/posts/slug/:slug` | Get by slug | No* |

*Optional auth - authenticated users may see draft posts they own

## Authentication

### JWT Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │────▶│  Login   │────▶│  Server  │
└──────────┘     └──────────┘     └──────────┘
                       │
                       ▼
              ┌────────────────┐
              │  Access Token  │ (Short-lived: 24h)
              │ Refresh Token  │ (Long-lived: 7 days)
              └────────────────┘
                       │
                       ▼
┌──────────┐     ┌──────────┐     ┌──────────┐
│  Client  │────▶│ API Call │────▶│  Server  │
│          │◀────│ + Bearer │◀────│          │
└──────────┘     └──────────┘     └──────────┘
```

### Password Requirements

- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one number
- At least one special character

### Using the API

```bash
# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!",
    "first_name": "John",
    "last_name": "Doe"
  }'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "Password123!"
  }'

# Access protected endpoint
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

## Database

### Supported Databases

- **SQLite** (default, for development)
- **PostgreSQL** (recommended for production)

### Models

#### User
- UUID primary key
- Email (unique)
- Password (bcrypt hashed)
- Role (user/admin/moderator)
- Status (active/inactive/banned/pending)
- Profile fields (first name, last name, avatar, bio)

#### Post
- UUID primary key
- Title, Slug (unique), Content
- Status (draft/published/archived)
- Author relationship
- Tags (many-to-many)

### Migrations

Migrations run automatically on startup using GORM's AutoMigrate.

## Middleware

| Middleware | Description |
|------------|-------------|
| **Recovery** | Recovers from panics and returns 500 |
| **Logger** | Logs all requests with timing |
| **CORS** | Handles cross-origin requests |
| **RateLimit** | Limits requests per client |
| **Auth** | Validates JWT tokens |
| **RequireRole** | Checks user role permissions |

### Middleware Chain

```
Request → Recovery → Logger → CORS → RateLimit → [Auth] → Handler
```

## Error Handling

### Error Response Format

```json
{
  "success": false,
  "error": {
    "code": 2003,
    "message": "Invalid credentials",
    "details": "Password does not match"
  }
}
```

### Error Codes

| Range | Category |
|-------|----------|
| 1000-1999 | General errors |
| 2000-2999 | Authentication errors |
| 3000-3999 | User errors |
| 4000-4999 | Database errors |

## Logging

Uses **Zap** for structured logging.

### Log Format (JSON)

```json
{
  "timestamp": "2024-01-15T10:30:00Z",
  "level": "info",
  "message": "Request completed",
  "request_id": "abc-123",
  "method": "GET",
  "path": "/api/v1/users",
  "status": 200,
  "latency": "1.5ms"
}
```

### Log Levels

- `debug` - Detailed debugging information
- `info` - General information
- `warn` - Warning messages
- `error` - Error messages
- `fatal` - Fatal errors (causes exit)

## Testing

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run short tests only
make test-short
```

## Deployment

### Production Checklist

- [ ] Set `APP_ENV=production`
- [ ] Use strong `JWT_SECRET` (32+ characters)
- [ ] Use PostgreSQL instead of SQLite
- [ ] Enable SSL for database connection
- [ ] Set appropriate rate limits
- [ ] Configure CORS for your domains
- [ ] Use HTTPS (reverse proxy)
- [ ] Set up monitoring/alerting
- [ ] Configure log aggregation

### Docker Production Build

```bash
# Build production image
docker build -t go-enterprise-api:latest .

# Run with environment variables
docker run -p 8080:8080 \
  -e APP_ENV=production \
  -e DB_DRIVER=postgres \
  -e DB_HOST=your-db-host \
  -e JWT_SECRET=your-secret \
  go-enterprise-api:latest
```

### Kubernetes

Example deployment:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-enterprise-api
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-enterprise-api
  template:
    metadata:
      labels:
        app: go-enterprise-api
    spec:
      containers:
      - name: api
        image: go-enterprise-api:latest
        ports:
        - containerPort: 8080
        env:
        - name: APP_ENV
          value: "production"
        livenessProbe:
          httpGet:
            path: /api/v1/health/live
            port: 8080
        readinessProbe:
          httpGet:
            path: /api/v1/health/ready
            port: 8080
```

## Best Practices

### Code Organization

1. **Handlers** - Keep thin, delegate to services
2. **Services** - Business logic lives here
3. **Repositories** - Only database operations
4. **Models** - Data structures with validation

### Security

1. Never store plain text passwords
2. Use parameterized queries (GORM handles this)
3. Validate all user input
4. Use HTTPS in production
5. Rotate JWT secrets periodically
6. Implement rate limiting
7. Log security events

### Performance

1. Use connection pooling
2. Implement caching where appropriate
3. Use pagination for list endpoints
4. Add database indexes
5. Use Gzip compression

### Error Handling

1. Use custom error types
2. Never expose internal errors to clients
3. Log errors with context
4. Return appropriate HTTP status codes

## Makefile Commands

```bash
make help          # Show all commands
make build         # Build the application
make run           # Run the application
make dev           # Run with hot reload
make test          # Run tests
make coverage      # Run tests with coverage
make lint          # Run linter
make fmt           # Format code
make clean         # Clean build artifacts
make docker-build  # Build Docker image
make docker-run    # Run Docker container
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
