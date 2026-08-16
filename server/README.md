# BoyKing Admin Server

Go 1.24+ backend with Gin + GORM + Casbin + JWT.

## Build & Test

```bash
go build ./cmd/admin-api       # entrypoint (web server)
go build ./cmd/modulegen       # code generator
go test ./...                  # all tests
```

## Project Layout

```
.
├── cmd/
│   ├── admin-api/             # web server entrypoint
│   └── modulegen/             # module code generator
├── internal/
│   ├── app/                   # bootstrap, container, scaffold
│   ├── interfaces/http/       # router + middleware
│   ├── modules/               # system + business modules (11 modules, 61 endpoints)
│   └── platform/              # shared infra (audit, auth, authz, dataauth, ratelimit, config, db, ...)
├── config/                    # config struct definitions
└── configs/                   # config.yaml (safe defaults) + config.example.yaml
```

## Adding a Module

```bash
go run ./cmd/modulegen -name <module-name>
```

Then register in `internal/modules/modules.go`. Architecture test enforces boundary rules automatically.

## Documentation

- `../docs/backend-handoff.md` — full API list and architecture overview
- `../docs/backend-architecture-blueprint.md` — design spec
