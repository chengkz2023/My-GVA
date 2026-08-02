# Module Guide

## Quick Start

Generate a new lightweight CRUD module:

```bash
cd server
go run ./cmd/modulegen -name order
```

This creates:
```
internal/modules/business/order/
  handler.go       # HTTP handler
  service.go       # Business logic
  repository.go    # Data access
  model.go         # Domain entity
  dto.go           # Request/response DTOs
  module.go        # Module wiring
  handler_test.go  # HTTP test
```

Then register it in `internal/modules/modules.go`:

```go
import businessorder "github.com/chengkz2023/My-GVA/server/internal/modules/business/order"

func HTTPModules(c *container.Container) []apphttp.Module {
    return []apphttp.Module{
        businessorder.NewModule(c),
        // ... existing modules
    }
}
```

## Module Shapes

### Lightweight (6 files)

For simple admin CRUD modules:

```
handler.go     — HTTP binding + response mapping
service.go     — Business logic orchestration
repository.go  — Database access (GORM)
model.go       — Database entity
dto.go         — Request/response DTOs
module.go      — Wiring (repo → service → handler)
```

### Full Layered (4 directories)

For modules with identities, permissions, state transitions, or reusable domain rules:

```
domain/           — Business entity + repository interface (no framework imports)
application/      — Use case orchestration + DTOs
infrastructure/   — GORM persistence implementation
transport/http/   — Gin handler + routes
module.go         — Wiring
```

See `system/user` and `system/role` for reference implementations.

## Architecture Boundaries

Enforced by `internal/architecture_test.go`:

| Layer | Must NOT import |
|---|---|
| domain | Gin, GORM, Zap, Viper, Redis, Casbin, transport, infrastructure |
| application | Gin |
| transport/http | GORM, infrastructure |
| platform | Any `internal/modules` package |
| module A | Module B's `infrastructure/` |

## Platform Dependencies

New modules should use platform packages for shared capabilities:

```
internal/platform/
  auth/        → Actor, PasswordHasher
  authz/       → Authorizer, PolicyProvider, PolicySyncer, AuthorityChecker
  config/      → Typed config reading
  database/    → DB connection
  errors/      → Unified error types (Validation, Unauthorized, Forbidden, NotFound, Conflict, Internal)
  pagination/  → Page, Result[T]
  response/    → OK(), Error(), Fail()
  transaction/ → Manager interface
```

## PR Checklist

```bash
cd server
go test ./...                    # All tests pass
go build ./cmd/admin-api         # Main build
go build ./cmd/modulegen         # Module generator
go test ./internal -run TestArchitectureBoundaries -count=1  # Architecture rules
```

## Running the Server

```bash
cd server
go run .
# Or:
go run ./cmd/admin-api
```
