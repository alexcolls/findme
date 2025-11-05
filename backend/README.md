# FindMe Backend API

Go backend service for the FindMe dating application.

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Databases**: PostgreSQL, Qdrant (vector DB), Redis
- **Authentication**: JWT
- **Architecture**: Clean Architecture / Domain-Driven Design

## Quick Start

### Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Make (optional)

### Setup

1. **Start dependencies**:
   ```bash
   docker-compose up -d
   ```

2. **Configure environment**:
   ```bash
   cp .env.sample .env
   # Edit .env with your settings
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Run the server**:
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

### Useful Commands

```bash
# Run with hot reload (install air first: go install github.com/air-verse/air@latest)
air

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Build for production
go build -o bin/api cmd/api/main.go

# Run with Docker
docker build -t findme-api .
docker run -p 8080:8080 findme-api

# Check linting
golangci-lint run
```

## Project Structure

```
backend/
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── api/              # HTTP layer (handlers, middleware, routes)
│   ├── domain/           # Business domain models
│   ├── service/          # Business logic
│   ├── repository/       # Data access layer
│   ├── infrastructure/   # External services
│   └── config/           # Configuration
├── pkg/                  # Public packages
├── migrations/           # Database migrations
├── scripts/              # Utility scripts
├── docker-compose.yml    # Development services
└── Dockerfile            # Production build
```

## API Endpoints

### Health Checks
- `GET /health` - API health status
- `GET /ready` - Readiness check

### API v1
- `GET /api/v1/` - API information

See [API Documentation](../docs/API.md) for complete endpoint reference.

## Development

### Running Dependencies

```bash
# Start all services
docker-compose up -d

# Check service status
docker-compose ps

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Database Migrations

```bash
# Run migrations
migrate -path migrations/postgres -database $DATABASE_URL up

# Rollback migration
migrate -path migrations/postgres -database $DATABASE_URL down 1
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Environment Variables

See `.env.sample` for all available configuration options.

Key variables:
- `API_PORT` - Server port (default: 8080)
- `APP_ENV` - Environment (development/production)
- `DATABASE_URL` - PostgreSQL connection string
- `QDRANT_URL` - Qdrant vector database URL
- `REDIS_URL` - Redis connection URL
- `JWT_SECRET_KEY` - JWT signing secret

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for development guidelines.

## License

See [LICENSE](../LICENSE) for license information.
