# MySQL Migrations

SQL migration files — applied by GORM AutoMigrate at startup (`bootstrap/database.go`).

## Migration files

```text
YYYYMMDDHHMMSS_name.up.sql
YYYYMMDDHHMMSS_name.down.sql
```

## Rules

- Production schema changes should be reviewed SQL migrations, not implicit `AutoMigrate`.
- Every `up` file should have a matching `down` file.
- Keep required seed data separate from demo seed data.
- Do not put secrets or environment-specific values in migration SQL.
