# MySQL Migrations

Migration files use paired `up` and `down` SQL files:

```text
YYYYMMDDHHMMSS_name.up.sql
YYYYMMDDHHMMSS_name.down.sql
```

Create a new pair from the `server` directory:

```bash
go run ./cmd/migrate create -name create_users
```

Validate naming and pair completeness:

```bash
go run ./cmd/migrate validate
```

Rules:

- Production schema changes should be reviewed SQL migrations, not implicit `AutoMigrate`.
- Every `up` file must have a matching `down` file.
- Keep required seed data separate from demo seed data.
- Do not put secrets or environment-specific values in migration SQL.
