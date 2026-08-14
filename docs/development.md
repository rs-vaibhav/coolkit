# Development Guide

This guide provides instructions for setting up your development environment and contributing to CoolKit.

## Prerequisites

- **Go**: 1.22 or higher
- **PostgreSQL**: 14 or higher
- **Docker & Docker Compose**: (Optional, but recommended for easy database setup)
- **golangci-lint**: For linting
- **golang-migrate**: For managing database migrations

## Setup

### macOS
```bash
brew install go postgresql golangci-lint golang-migrate
```

### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install golang postgresql
# For golangci-lint and golang-migrate, follow official docs for Linux
```

### Windows
```powershell
scoop install go postgresql golangci-lint
```

## Database Setup (Manual)

If you are not using Docker, set up PostgreSQL manually:
1. Access `psql`.
2. Run:
   ```sql
   CREATE DATABASE coolkit;
   CREATE USER coolkit_user WITH ENCRYPTED PASSWORD 'coolkit_pass';
   GRANT ALL PRIVILEGES ON DATABASE coolkit TO coolkit_user;
   ```

## Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Edit `.env` and set the variables according to your environment.
   - `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
   - `JWT_SECRET`: Secret key for signing tokens.

## Running the Server

Using Make:
```bash
make run
```
Or directly with Go:
```bash
go run cmd/server/main.go
```

## Running Migrations

Apply up migrations:
```bash
make migrate-up
```

Revert down migrations:
```bash
make migrate-down
```

Create a new migration:
```bash
migrate create -ext sql -dir db/migrations -seq your_migration_name
```

## Seeding Data

To populate the database with dummy data for testing:
```bash
make seed
# or
go run scripts/seed/main.go
```

## Running Tests

Run all tests:
```bash
make test
```

Run tests quickly (skip long tests):
```bash
make test-short
```

Run tests with coverage report:
```bash
make coverage
```

## Linting

Format code and check for lint errors:
```bash
make fmt
make lint
```

## Docker Development

Start PostgreSQL and API in containers:
```bash
make docker-up
```

Stop containers:
```bash
make docker-down
```

View logs:
```bash
make docker-logs
```

## Adding a New Feature

Follow this general flow when adding a feature:
1. **Model**: Define struct in `internal/model/entity.go`.
2. **Repository**: Create interface and GORM implementation in `internal/repository`.
3. **Service**: Create business logic in `internal/service`.
4. **Handler**: Create HTTP parsing/responses in `internal/handler`.
5. **Route**: Wire the handler to the router in `cmd/server/main.go`.
6. **Test**: Write unit tests for your service and handler.

## Troubleshooting

- **Database Connection Error**: Ensure PostgreSQL is running and credentials in `.env` match.
- **Port Already in Use**: Kill the process using port 8080 or change the port in `.env`.
