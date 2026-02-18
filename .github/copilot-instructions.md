# GoBox - Copilot Instructions

## Architecture Overview

GoBox is a container-as-a-service platform that dynamically provisions isolated Docker environments based on user fingerprints. The system manages lifecycle of ephemeral containers with connection-aware auto-shutdown.

**Key architectural pattern:** Ports & Adapters (Hexagonal Architecture)

- `internal/domain/` - Core business entities (e.g., Box)
- `internal/box/` - Business logic service layer with port interfaces (Repo, DockerSvc)
- `internal/repo/` - Repository implementations (adapter)
- `internal/docker/` - Docker client adapter
- `internal/rest/handler/` - HTTP/WebSocket handlers (adapter)
- `cmd/` - Application entry points (server, migrations)

## Critical Workflows

### Development Setup

```bash
make dev          # Build images, start containers, tail logs
make db-shell     # Access PostgreSQL CLI
make app-shell    # Access app container shell
```

Migrations run automatically on server startup via [main.go](main.go#L16) → `cmd.RunMigrate()`.

### Database Schema Management

- **Migrations:** Use Goose in [migrations/](migrations/) directory. Create with `make migration` (prompts for name).
- **Query Generation:** SQLC generates type-safe Go code from [queries/\*.sql](queries/) → [internal/infra/db/sqlc/](internal/infra/db/sqlc/). After modifying queries, regenerate with `sqlc generate`.
- **Custom types:** [sqlc.yml](sqlc.yml#L16-L20) defines overrides for UUIDs and enums (e.g., `box_status` → `string`).

### Container Lifecycle Pattern

The box service implements intelligent connection counting:

- [internal/box/svc.go](internal/box/svc.go#L41-L98) - `manageConnections()` goroutine tracks active WebSocket connections per fingerprint
- When last connection closes, starts 5s shutdown timer for container
- New connections cancel pending shutdowns
- Container state persists in PostgreSQL with `box_status` enum

## Code Conventions

### Dependency Injection

Services receive dependencies via constructor, defined as interface ports:

```go
// internal/box/port.go - Define dependencies as interfaces
type Repo interface { ... }
type DockerSvc interface { ... }

// internal/box/svc.go - Inject via constructor
func NewSvc(repo Repo, dockerSvc DockerSvc) *Svc
```

Wire dependencies in [cmd/server.go](cmd/server.go#L39-L42).

### Configuration

- Uses Viper to load `.env` files + environment variables ([internal/config/config.go](internal/config/config.go))
- Docker Compose overrides `POSTGRES_HOST` to `postgres` for container networking
- Dev env runs on port `8010` (configurable via `SERVER_PORT`)

### Error Handling

**AppError Pattern:** All layers use structured `AppError` from [domain/error.go](internal/domain/error.go) for consistent error handling across the application.

- **Error Types:** `NOT_FOUND`, `VALIDATION`, `INTERNAL`, `DATABASE`, `DOCKER`, `CONFLICT`, `UNAUTHORIZED`
- **Repository layer:** Maps database errors (pgx) to AppError. `ErrNoRows` → `NotFoundError`, constraint violations → `ValidationError`/`ConflictError` ([box_repo.go](internal/repo/box_repo.go#L100-L130))
- **Service layer:** Only works with AppError. Validates inputs and propagates errors from repo/docker layers ([box/svc.go](internal/box/svc.go))
- **Handler layer:** Extracts HTTP status code and error message from AppError using `domain.GetStatusCode()` and `domain.GetErrorMessage()` ([box/handler.go](internal/rest/handler/box/handler.go#L50-L65))

**Error Construction:**

```go
// Use constructor functions
domain.NewNotFoundError("box", fingerprint)
domain.NewValidationError("fingerprint cannot be empty")
domain.NewDatabaseError("create box", err)
domain.NewDockerError("start container", err)
```

**Error Checking:**

```go
// Check if error is AppError and its type
if appErr, ok := domain.IsAppError(err); ok && appErr.IsType(domain.ErrorTypeNotFound) {
    // Handle not found case
}
```

### WebSocket Communication

- Handler upgrades HTTP to WebSocket in [box/handler.go](internal/rest/handler/box/handler.go#L27-L40)
- Fingerprint passed as query parameter: `/connect?fingerprint=abc123`
- Buffer sizes: 8KB read/write ([handler.go](internal/rest/handler/box/handler.go#L10-L11))

## Docker Integration

### Host Docker Access

App container mounts `/var/run/docker.sock` ([docker-compose.yml](docker-compose.yml#L14)) to control host's Docker daemon. This enables spawning sibling containers, not nested containers.

### Base Image Strategy

[cmd/server.go](cmd/server.go#L52-L64) ensures:

1. Docker network `gobox-c-network` (172.25.0.0/16) exists
2. Base image `gobox-base:latest` built from [base-image/Dockerfile](base-image/Dockerfile)
3. All user containers attach to this network for isolation

### Network Architecture

- Main app: `gobox-network` (Compose default)
- User containers: `gobox-c-network` (custom bridge with static subnet)

## Testing & Tools

**Development tools** (install with `make install-tools`):

- `air` - Hot reload (not currently configured)
- `goose` - Migrations
- `sqlc` - Query code generation
- `golangci-lint` - Linting

**Test commands:**

```bash
make test              # Unit tests
make test-coverage     # Generate coverage.html
make lint              # Run golangci-lint
make fmt               # Format code
```

## Important Files

- [main.go](main.go) - Runs migrations then starts server
- [cmd/server.go](cmd/server.go) - Wire dependencies, setup middleware (chi router)
- [internal/box/svc.go](internal/box/svc.go) - Core business logic with connection management
- [sqlc.yml](sqlc.yml) - SQLC configuration for query generation
- [Makefile](Makefile) - All operational commands

## Common Gotchas

- When adding domain types, update SQLC overrides in [sqlc.yml](sqlc.yml#L16)
- Don't run `go run main.go` directly - migrations expect `./migrations` relative path from project root
- PostgreSQL custom enum types require manual type creation in migrations ([example](migrations/20260211192039_create_box_table.sql#L4-L9))
- SQLC queries use `pgx/v5` type bindings, not `database/sql` ([sqlc.yml](sqlc.yml#L9))
