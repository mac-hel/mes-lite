# MES Lite

A lightweight Manufacturing Execution System for small manufacturing companies.

## Stack

- **Language:** Go
- **HTTP:** Fuego
- **Database:** PostgreSQL via pgx + sqlc
- **Migrations:** goose
- **Container:** Docker

## Development

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- golangci-lint (optional, CI runs it)
- air (optional, for live reload)

### Quick Start

```bash
# Start database
make docker-up

# Run the server (with live reload if air is installed)
make run
```

### Commands

```bash
make build    # Build binary
make test     # Run tests
make lint     # Run linter
make fmt      # Format code
make tidy     # Tidy modules
```

## Project Structure

```
cmd/server/     Application entry point
internal/       Private application packages
migrations/     Database migrations
docs/adr/       Architecture Decision Records
```
