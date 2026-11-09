# LungCare backend

Production-oriented Go REST API for the lung disease prediction system. It uses Gin for HTTP routing, GORM for persistence, MySQL, explicit dependency wiring, environment-based configuration, structured logs, health checks, and graceful shutdown.

## Structure

```text
cmd/api/                 executable entry point
internal/app/            dependency wiring and process lifecycle
internal/config/         validated environment configuration
internal/database/       MySQL connection and pool policy
internal/domain/         persistence/domain models and migrations
internal/repository/     data-access boundary
internal/service/        application and authentication logic
internal/handler/        HTTP routes, validation, and responses
internal/httpx/          shared HTTP error contract
```

The dependency direction is `handler -> service -> repository`; configuration and database setup are owned by the application bootstrap.

## Run locally

1. Install Go 1.25+ and MySQL 8+.
2. Copy `.env.example` to `.env`, replace the sample credentials, then export them with `set -a; source .env; set +a`. The application intentionally does not load `.env` files itself in production.
3. Create the configured database.
4. Run `make migrate` to apply the schema.
5. Run `make run`.

The service exposes:

- `GET /health/live` for process liveness.
- `GET /health/ready` for database readiness.
- API resources under `/api/v1` (the original endpoint paths are retained).

## Quality checks

```sh
make fmt
make vet
make test
```

## Production notes

- Keep `DATABASE_DSN` in a secret manager, never in source control.
- Set `APP_ENV=production` and list exact browser origins in `CORS_ALLOWED_ORIGINS`.
- Run `make migrate` as a controlled deployment step. The API process does not mutate its schema at startup.
- Authentication currently validates credentials but does not issue an access token. Add an identity/session layer before exposing protected write endpoints publicly.
