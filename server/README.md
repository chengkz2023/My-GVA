# BoyKing Admin Server

Go 1.24+ backend with Gin + GORM + Casbin + JWT.

## Build & Test

```bash
go build ./cmd/admin-api       # V2 entrypoint (web server)
go build ./cmd/migrate         # migration CLI
go build ./cmd/modulegen       # code generator
go test ./...                  # all tests (28 packages)
```

## Project Layout

```
.
├── main.go                    # bootstrap.Run()
├── cmd/
│   ├── admin-api/             # future entrypoint
│   ├── migrate/               # migration CLI (up/down/status/dry-run)
│   └── modulegen/             # module code generator
├── internal/
│   ├── app/                   # bootstrap, container, migration engine, scaffold
│   ├── interfaces/http/       # V2 router + middleware
│   ├── modules/               # system + business modules (11 modules, 50 endpoints)
│   └── platform/              # shared infra (auth/JWT/claims, authz, casbin, config, db, ...)
├── config/                    # config struct definitions
├── model/                     # legacy GORM models (25 files)
├── utils/                     # backward-compat wrappers, timer, AST
└── global/                    # legacy global variables
```

## Adding a Module

```bash
go run ./cmd/modulegen -name <module-name>
```

Then register in `internal/modules/modules.go`. Architecture test enforces boundary rules automatically.

## Documentation

- `../docs/backend-v2-handoff.md` — full API list and architecture overview
- `../docs/backend-architecture-blueprint.md` — design spec
