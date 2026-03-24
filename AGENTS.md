# Agent Coding Guidelines

This document provides guidelines for agents working on this codebase.

## Project Overview

Go microservices project with:

- **Language**: Go 1.25.8
- **API Framework**: Gin (REST), gRPC
- **Databases**: PostgreSQL, MySQL, MongoDB, Redis
- **Messaging**: Kafka

## Services

| Service | Port | Protocol |
| --- | --- | --- |
| auth-service | 8081 | REST + gRPC |
| user-service | 8082 | REST + gRPC |
| order-service | - | gRPC |
| notification-service | - | Consumer |
| analytics-service | - | Consumer |
| common | - | Shared packages |

## Build/Lint/Test Commands

### Root Level

```bash
make build-all  # Build all modules
```

### Per-Service

```bash
cd auth-service

make install-tools  # Install swag, golangci-lint, goose
make update         # go mod tidy + go get -u
make fmt            # go fmt
make vet            # go vet
make lint           # golangci-lint run
make check          # fmt + vet + lint
make swagger        # Generate swagger docs
make build          # Build application
make run            # Run application
make clean          # Clean build artifacts
```

### Running a Single Test

```bash
go test -v -run TestAuthService ./src/internal/service/auth/...
go test -v -cover ./...
```

## Code Style

### Naming Conventions

- **Interfaces**: `XxxItf` (e.g., `AuthServiceItf`)
- **Private structs**: lowercase (e.g., `authService`)
- **Init functions**: `InitXxxService`, `InitXxxRepository`
- **Config structs**: `Options` with YAML tags
- **Constants**: `PascalCase`

### Project Structure

```text
service-name/
├── src/
│   ├── cmd/           # main.go, app.go
│   ├── internal/
│   │   ├── config/
│   │   ├── handler/
│   │   ├── model/
│   │   ├── repository/
│   │   └── service/
├── etc/migrations/, sql/
├── config.yaml
├── go.mod
└── Makefile
```

### Imports

Group with blank lines:

1. Standard library
2. External packages (github.com)
3. Local packages

```go
import (
    "context"
    "fmt"

    "github.com/linggaaskaedo/go-kill/auth-service/src/internal/repository/auth"
    authpb "github.com/linggaaskaedo/go-kill/common/pkg/proto/auth"

    "golang.org/x/crypto/bcrypt"
)
```

### Error Handling

- Use `palantir/stacktrace` for error wrapping
- Return `(value, error)` pattern
- Use error codes from `common/pkg/errors/error_code.go`

```go
return nil, stacktrace.Propagate(err, "failed to create user")
```

### Logging

- Use `rs/zerolog`
- Initialize via `logger.Init(config)` in `common/pkg/logger/logger.go`

```go
log.Info().Str("user_id", userID).Msg("created")
log.Error().Err(err).Msg("failed")
```

### Configuration

```go
type Options struct {
    JwtSecret string `yaml:"jwt_secret"`
    Topic     string `yaml:"topic"`
}
```

### Database Patterns

- `sqlx` for SQL databases
- MongoDB driver v2
- `go-redis/v9` for Redis
- Repository pattern with interfaces

### HTTP Handlers

- Gin framework for REST APIs
- Return JSON responses with proper status codes

### Messaging

- Kafka for event-driven communication
- Consume in notification-service and analytics-service
