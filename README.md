# LogLine

**A centralized, self-hosted logging platform written in Go.**
LogLine gives your services a single HTTP endpoint to ship structured logs to, a queryable
storage layer backed by PostgreSQL, and a server-rendered web dashboard for triage and
observability — with no external SaaS dependency and a small enough footprint to run on a
single container.

[![Go Version](https://img.shields.io/badge/go-1.27-00ADD8?logo=go)](https://go.dev)
[![Go Build](https://github.com/MdRasB/LogLine/actions/workflows/go.yml/badge.svg)](https://github.com/MdRasB/LogLine/actions/workflows/go.yml)
[![PostgreSQL](https://img.shields.io/badge/postgres-18-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Docker](https://img.shields.io/badge/docker-ready-2496ED?logo=docker&logoColor=white)](Dockerfile)
[![Status](https://img.shields.io/badge/status-active--development-yellow)](#roadmap--known-limitations)

---

## Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Quick Start (Docker Compose)](#quick-start-docker-compose)
  - [Local Development (without Docker)](#local-development-without-docker)
- [Configuration](#configuration)
- [Database Schema](#database-schema)
- [API Reference](#api-reference)
  - [Authentication](#authentication)
  - [Ingest](#post-ingest)
  - [Query Logs](#get-logs)
  - [Dashboard](#get-dashboard)
  - [Health](#get-health)
- [Security Model](#security-model)
- [Request Lifecycle & Middleware Pipeline](#request-lifecycle--middleware-pipeline)
- [Testing](#testing)
- [Makefile Reference](#makefile-reference)
- [Deployment](#deployment)
- [Roadmap / Known Limitations](#roadmap--known-limitations)
- [Contributing](#contributing)

---

## Overview

LogLine is a small, dependency-light logging backend built around three primitives:

1. **Ingest** — services `POST` structured log entries over HTTP.
2. **Store** — entries land in PostgreSQL with indexed columns for service, level, and time.
3. **Query / Visualize** — logs are queryable via a JSON API and browsable through a
   server-rendered HTML dashboard with aggregate charts (volume over time, level breakdown,
   per-service breakdown).

The project is intentionally built on the Go standard library (`net/http`, `html/template`)
rather than a web framework — routing, middleware chaining, and the dashboard rendering
pipeline are all hand-rolled, which makes the whole request path easy to trace end to end.

## Architecture

```mermaid
flowchart TB
    subgraph Clients
        SVC[Application services]
        USR[Operator / Browser]
    end

    subgraph LogLine["LogLine API Server (net/http)"]
        MW["Middleware Chain\nRequestID → Recovery → RateLimit → (Auth) → Logging"]
        H_ING[Ingest Handler]
        H_LOG[Logs Handler]
        H_DASH[Dashboard Handler]
        H_AUTH[Auth Handler]
        H_HEALTH[Health Handler]
        TPL[HTML Template Engine]
    end

    subgraph Domain["Domain Services"]
        AUTHSVC[auth.Service\nsession + password logic]
    end

    subgraph Storage
        PG[(PostgreSQL)]
        LOGSTORE[DBStore]
        USERSTORE[UserStore]
        SESSSTORE[SessionStore]
    end

    SVC -->|"POST /ingest (API-authenticated)"| MW
    USR -->|"GET /dashboard, /logs, /auth/*"| MW

    MW --> H_ING
    MW --> H_LOG
    MW --> H_DASH
    MW --> H_AUTH
    MW --> H_HEALTH

    H_ING --> LOGSTORE
    H_LOG --> LOGSTORE
    H_DASH --> LOGSTORE
    H_DASH --> TPL
    H_AUTH --> AUTHSVC
    H_HEALTH -->|ping| PG

    AUTHSVC --> USERSTORE
    AUTHSVC --> SESSSTORE

    LOGSTORE --> PG
    USERSTORE --> PG
    SESSSTORE --> PG
```

**Design notes:**

- The server is a single Go binary (`cmd/api`) that wires config → database pool → stores →
  auth service → middleware → routes in `internal/server`. There is no dependency-injection
  framework; wiring is explicit and visible in `NewServer`.
- Two middleware chains exist — `publicChain` and `protectedChain` — composed from the same
  middleware primitives in different orders (see [Request Lifecycle](#request-lifecycle--middleware-pipeline)).
- The dashboard is server-rendered (`html/template`) and hydrated client-side only for charts
  (Chart.js reads a JSON payload embedded in the page and draws volume/level/service graphs).
- Data access is a thin, explicit SQL layer over `pgx/v5` — no ORM. Filter-to-SQL translation
  lives in `internal/db/query.go` and is shared between the counting query and the paginated
  query.

## Features

- **HTTP log ingestion** — `POST /ingest` accepts structured JSON log entries (`level`,
  `message`, `service`, `timestamp`, optional `metadata`) with server-side validation.
- **Filterable query API** — `GET /logs` supports filtering by `service`, `level`, free-text
  `search` (case-insensitive substring match on message), a `from`/`to` RFC3339 time range,
  and cursor-free `page`/`limit` pagination (hard-capped at 100 rows/page).
- **Web dashboard** — aggregate stats (total logs, total services, error count), an hourly log
  volume chart, a level breakdown, and a service breakdown, plus a filterable/paginated log
  table — all rendered server-side and charted with Chart.js.
- **Session-based authentication** — email/password registration and login, bcrypt password
  hashing, cryptographically random session tokens (SHA-256 hashed at rest, never stored in
  plaintext), and bearer-token protected routes.
- **API key primitives** — key generation and constant-time verification utilities
  (`ll_live_*` prefix) plus a dedicated `api_keys` table, intended for service-to-service
  ingest authentication (see [Roadmap](#roadmap--known-limitations)).
- **Per-IP rate limiting** — token-bucket limiter (`golang.org/x/time/rate`) keyed by client
  IP, with automatic cleanup of stale clients.
- **Structured request logging** — every request is logged via `log/slog` with a generated
  request ID, method, path, status code, duration, and (where available) user/session ID.
- **Panic recovery** — a dedicated recovery middleware catches panics per-request, logs the
  stack trace, and returns a clean `500` instead of crashing the server.
- **Health checks** — `GET /health` reports process uptime, build version, and a live
  database ping, suitable for container/orchestrator liveness or readiness probes.
- **Schema-managed migrations** — versioned SQL migrations via [goose](https://github.com/pressly/goose).
- **Graceful shutdown** — `SIGINT`/`SIGTERM` triggers a bounded (15s) graceful HTTP shutdown
  followed by a clean database pool close.
- **Container-first** — multi-stage Dockerfile producing a small Alpine runtime image running
  as a non-root user, plus a ready-to-run `docker compose` stack (app + PostgreSQL with a
  health-gated startup order).

## Tech Stack

| Layer              | Choice                                                            |
|---------------------|--------------------------------------------------------------------|
| Language             | Go 1.27                                                            |
| HTTP routing         | Standard library `net/http.ServeMux` (no third-party framework)   |
| Database             | PostgreSQL 18 (Alpine image in Compose)                           |
| DB driver            | [`jackc/pgx/v5`](https://github.com/jackc/pgx) (`pgxpool` connection pool) |
| Migrations           | [`pressly/goose`](https://github.com/pressly/goose)                |
| Templating           | Standard library `html/template`                                  |
| Frontend charts      | [Chart.js](https://www.chartjs.org/) (vanilla JS, no bundler/build step) |
| Password hashing     | `bcrypt` (cost factor 12) via `golang.org/x/crypto`                |
| Structured logging   | Standard library `log/slog`                                       |
| Rate limiting        | `golang.org/x/time/rate` (token bucket)                            |
| Config               | Environment variables (`.env` support via `joho/godotenv`)         |
| Containerization     | Docker (multi-stage build) + Docker Compose                        |
| CI                   | GitHub Actions (build on push/PR to `main`)                        |

> **Note on `go.sum`:** the module graph includes a long tail of indirect dependencies for
> ClickHouse, MySQL, SQLite, YDB, Vertica, etc. These are **transitive dependencies of goose**
> (which supports many database backends), not drivers used by LogLine itself. The application
> only ever talks to PostgreSQL via `pgx`.

## Project Structure

```
LogLine/
├── cmd/
│   └── api/                # main.go — process entrypoint, signal handling, graceful shutdown
├── internal/
│   ├── auth/                # Session/password/API-key crypto primitives + auth.Service (business logic)
│   ├── config/               # Environment-variable configuration loader
│   ├── contextutil/           # Typed context helpers (request ID, user ID, session ID)
│   ├── dashboard/             # View-model structs for the HTML dashboard
│   ├── db/                    # pgx connection pool + Store implementations (logs, users, sessions) + SQL filter builder
│   ├── handler/                # HTTP handlers (ingest, logs, dashboard, auth, health) + request validation
│   ├── middleware/             # RequestID, Recovery, Logging, RateLimiter, Auth, chaining
│   ├── model/                  # Domain structs (Logs, LogEntry, LogFilter, User, Session, PaginatedLogs)
│   ├── server/                  # Server struct, route registration, middleware chain composition
│   └── web/                     # html/template manager + template helper funcs
├── migrations/                  # goose SQL migrations (logs, users, sessions, api_keys)
├── web/
│   ├── static/                   # CSS + vanilla JS (Chart.js dashboard hydration)
│   └── templates/                 # base.html, dashboard.html, partials/ (navbar, filters, charts, log table, pagination…)
├── Dockerfile                     # Multi-stage build → minimal Alpine runtime image
├── compose.yaml                   # app + postgres services for local/dev orchestration
├── Makefile                       # Dev, DB, migration, test, and Docker workflow shortcuts
└── example.env                    # Template for local environment configuration
```

## Getting Started

### Prerequisites

- [Go 1.27+](https://go.dev/dl/)
- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)
  (recommended path — spins up PostgreSQL for you)
- [PostgreSQL 15+](https://www.postgresql.org/) if you prefer to run the database yourself
- `goose` (installed on demand via `go run` — no separate install required) for migrations

### Quick Start (Docker Compose)

This builds the app image and starts it alongside a health-checked PostgreSQL instance.

```bash
git clone https://github.com/MdRasB/LogLine.git
cd LogLine

# Build and start both services
make compose-up
# equivalent to: docker compose up -d --build

# Apply database migrations against the containerized Postgres
make migrate

# Tail logs
make compose-logs
```

The API is now listening on `http://localhost:8080`. Verify it with:

```bash
curl http://localhost:8080/health
```

Bring the stack down with `make compose-down`.

### Local Development (without Docker)

```bash
# 1. Start only the database container
make db-up

# 2. Copy and adjust environment variables
cp example.env .env

# 3. Apply migrations
make migrate

# 4. Run the server directly with `go run`
make run
# equivalent to: go run cmd/api/main.go
```

## Configuration

LogLine is configured entirely via environment variables (loaded from a `.env` file if
present, via `godotenv`). See `example.env` for a working template.

| Variable                | Default    | Description                                                                 |
|--------------------------|------------|-------------------------------------------------------------------------------|
| `PORT`                    | `:8079`    | Listen address/port. Must be in `:PORT` form.                                |
| `DB_URL`                  | *(none)*   | **Required.** PostgreSQL connection string, e.g. `postgres://user:pass@host:5432/db`. |
| `REQLIMIT`                | `5`        | Allowed requests per second, per client IP (token bucket refill rate).       |
| `BURST`                   | `10`       | Token bucket burst capacity per client IP.                                   |
| `VERSION`                 | `v1.0.0`   | Build/version string surfaced by `GET /health`.                              |
| `READ_TIMEOUT`            | `10s`      | `http.Server.ReadTimeout`.                                                    |
| `READ_HEADER_TIMEOUT`     | `5s`       | `http.Server.ReadHeaderTimeout`.                                              |
| `WRITE_TIMEOUT`           | `15s`      | `http.Server.WriteTimeout`.                                                   |
| `IDLE_TIMEOUT`            | `60s`      | `http.Server.IdleTimeout`.                                                    |

Startup is **intended** to fail fast (non-zero exit) if `DB_URL` is unset, or if any
timeout/rate-limit value fails to parse or is non-positive. **This currently doesn't work
correctly** — see the `config.Load()` bug under
[Roadmap / Known Limitations](#roadmap--known-limitations): misconfiguration today causes a
nil-pointer panic on startup rather than the clean error message this table implies.

## Database Schema

Managed via four sequential goose migrations under `migrations/`:

| Migration                  | Table       | Purpose                                                                 |
|------------------------------|-------------|---------------------------------------------------------------------------|
| `001_create_logs.sql`         | `logs`       | Core log storage: `level`, `message`, `service`, `timestamp`, `metadata` (JSONB), `created_at`. Indexed on `service`, `level`, and `timestamp DESC`. |
| `002_create_users.sql`        | `users`      | Dashboard/API account records: `email` (unique), `password_hash`, `created_at`. |
| `003_create_sessions.sql`     | `sessions`   | Active login sessions: SHA-256 `token_hash` (unique), `expires_at`, cascade-deletes with the owning user. |
| `004_create_api_keys.sql`     | `api_keys`   | Service-scoped API keys: `key_prefix`, SHA-256 `key_hash` (unique), `revoked_at` for soft revocation. |

```bash
# Apply all pending migrations
make migrate

# Check migration status
make migrate-status
```

## API Reference

All request/response bodies are JSON unless otherwise noted. Errors are returned as
`{"error": "<message>"}` with an appropriate HTTP status code.

### Authentication

#### `POST /auth/register`

Creates a new user account.

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "password": "a-strong-password"}'
```

```json
{ "message": "registration successful" }
```

#### `POST /auth/login`

Authenticates and issues a session token.

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "you@example.com", "password": "a-strong-password"}'
```

```json
{ "session_token": "ll_sess_..." }
```

#### `POST /auth/logout` 🔒

Invalidates the current session. Requires `Authorization: Bearer <session_token>`.

```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer ll_sess_..."
```

### `POST /ingest` 🔒

Accepts a single structured log entry.

```bash
curl -X POST http://localhost:8080/ingest \
  -H "Authorization: Bearer ll_sess_..." \
  -H "Content-Type: application/json" \
  -d '{
        "level": "error",
        "message": "payment webhook timed out",
        "service": "billing-api",
        "timestamp": "2025-01-15T10:00:00Z",
        "metadata": { "order_id": "ord_123", "retry": 2 }
      }'
```

| Field       | Type    | Required | Notes                                                            |
|-------------|---------|----------|-------------------------------------------------------------------|
| `level`      | string   | yes       | One of `error`, `warn`, `info`, `debug`, `fatal`.                  |
| `message`    | string   | yes       | Free-text log message.                                             |
| `service`    | string   | yes       | Name of the emitting service.                                      |
| `timestamp`  | string   | yes       | RFC3339 timestamp.                                                  |
| `metadata`   | object   | no        | Arbitrary JSON, stored as `JSONB`.                                  |

```json
{ "message": "log accepted" }
```

### `GET /logs` 🔒

Query stored logs with filters and pagination.

```bash
curl "http://localhost:8080/logs?service=billing-api&level=error&search=timeout&page=1&limit=20" \
  -H "Authorization: Bearer ll_sess_..."
```

| Query param | Type   | Notes                                                          |
|-------------|--------|-------------------------------------------------------------------|
| `service`    | string  | Exact match.                                                      |
| `level`      | string  | One of `error`, `warn`, `info`, `debug`, `fatal`.                  |
| `search`     | string  | Case-insensitive substring match against `message`.               |
| `from` / `to`| string  | RFC3339 time range bounds (inclusive).                             |
| `page`       | int     | Default `1`.                                                       |
| `limit`      | int     | Default `20`, hard-capped at `100`.                                 |

```json
{
  "Logs": [ { "ID": "...", "Level": "error", "Message": "...", "Service": "billing-api", "Timestamp": "...", "Metadata": {}, "CreatedAt": "..." } ],
  "Total": 42,
  "Page": 1,
  "Limit": 20,
  "HasMore": true
}
```

### `GET /dashboard`

Renders the server-side HTML dashboard (accepts the same `service`, `level`, `search`,
`page`, `limit` query parameters as `/logs`). Currently unauthenticated, and — unlike
`GET /logs` — its `limit` parameter is **not** capped at 100, so it's an easy target for
resource-exhaustion abuse until both are fixed (see
[Security Model](#security-model) and [Known Limitations](#roadmap--known-limitations)).

### `GET /health`

Liveness/readiness probe.

```json
{
  "status": "healthy",
  "database": "up",
  "uptime": "2h13m5s",
  "started_at": "2025-01-15T08:00:00Z",
  "time": "2025-01-15T10:13:05Z",
  "version": "v0.8.10"
}
```

Returns `503 Service Unavailable` with `"status": "unhealthy"` if the database is unreachable.

## Security Model

- **Passwords** are hashed with `bcrypt` at cost factor `12` and never logged or returned.
- **Session tokens** (`ll_sess_<64 hex chars>`) are generated from `crypto/rand`, returned to
  the client exactly once at login, and stored server-side only as a SHA-256 hash — the
  plaintext token is unrecoverable from the database. Comparison during validation would
  ideally use `crypto/subtle.ConstantTimeCompare` end-to-end (see [Roadmap](#roadmap--known-limitations)
  regarding session lookup).
- **API keys** (`ll_live_<64 hex chars>`) follow the same generate-once / hash-at-rest /
  constant-time-verify pattern as sessions, via `auth.GenerateAPIKey` / `auth.VerifyAPIKey`.
- **Rate limiting** is applied per client IP on every route (public and protected) before any
  handler logic runs, mitigating brute-force and abuse traffic at the edge of the middleware
  chain. Note it currently keys off `r.RemoteAddr` directly, so behind a reverse proxy or load
  balancer every client will appear to share the proxy's IP unless a trusted
  `X-Forwarded-For`/`X-Real-IP` header is read instead.
- **Body size limits** — ingest and auth handlers wrap the request body in
  `http.MaxBytesReader` (1 MiB) to bound memory usage from oversized payloads.
- **Strict JSON decoding** — request bodies are decoded with `DisallowUnknownFields`, rejecting
  unexpected fields rather than silently ignoring them.
- **Panic isolation** — the recovery middleware ensures an unhandled panic in any single
  handler returns a `500` instead of taking down the process or leaking a stack trace to the
  client.
- **`model.User.PasswordHash` has no `json:"-"` tag.** No current handler serializes a `User`
  struct directly back to a client, so this isn't exploited today, but it's a latent leak: the
  first future endpoint that returns a `User` (e.g. a "get profile" route) by mistake would
  ship the bcrypt hash to the client. Worth tagging defensively now.

## Request Lifecycle & Middleware Pipeline

Every route is built by composing the same middleware primitives in a different order for
public vs. protected endpoints (`internal/server/routes.go`):

```mermaid
sequenceDiagram
    participant C as Client
    participant RID as RequestID
    participant REC as Recovery
    participant RL as RateLimiter
    participant AUTH as AuthMiddleware
    participant LOG as Logging
    participant H as Handler

    Note over C,H: Public route (e.g. /auth/login, /health, /dashboard)
    C->>RID: HTTP request
    RID->>REC: inject request_id, set X-Request-ID header
    REC->>LOG: recover() wrapper installed
    LOG->>RL: (start timer)
    RL->>H: token available? forward
    H-->>C: response
    LOG-->>C: log line: request_id, method, path, status, duration

    Note over C,H: Protected route (e.g. /ingest, /logs, /auth/logout)
    C->>RID: HTTP request
    RID->>REC: inject request_id
    REC->>RL: recover() wrapper installed
    RL->>AUTH: token available? forward
    AUTH->>LOG: session validated → inject user_id, session_id
    LOG->>H: (start timer)
    H-->>C: response
    LOG-->>C: log line: request_id, user_id, session_id, method, path, status, duration
```

Note the deliberate ordering difference: on the protected chain, **rate limiting and recovery
run before authentication**, so unauthenticated/abusive traffic is throttled and can't crash
the process before a single DB lookup for session validation occurs. **Logging runs last** in
both chains so that the final response status and duration are captured accurately.

## Testing

```bash
make test
# equivalent to: go test ./...
```

- `internal/auth` — pure unit tests covering API key generation/hashing/verification and
  bcrypt password hashing/verification (no external dependencies).
- `internal/db` — integration-style tests (`user_store_test.go`) that exercise `UserStore`
  against a real PostgreSQL instance. They expect a reachable database at
  `postgres://postgres:postgres@localhost:5432/logline_test?sslmode=disable` — point this at a
  disposable test database before running the full suite.

Static analysis / linting:

```bash
make check
# equivalent to: golangci-lint run -v
```

> CI (`.github/workflows/go.yml`) currently runs `go build ./...` on every push/PR to `main`;
> the `go test` step is present but commented out pending a CI-provisioned test database.

## Makefile Reference

| Command              | Description                                                   |
|-----------------------|--------------------------------------------------------------------|
| `make run`             | Run the API server directly with `go run`.                        |
| `make build`           | Build the `logline` binary from `cmd/api`.                        |
| `make compile`         | Compile all packages (`go build ./...`) without producing a binary.|
| `make db-up`           | Start only the PostgreSQL container.                               |
| `make db-down`         | Stop the PostgreSQL container.                                     |
| `make db-logs`         | Tail PostgreSQL container logs.                                    |
| `make db-run`          | Open a `psql` shell inside the running database container.         |
| `make migrate`         | Apply all pending goose migrations.                                 |
| `make migrate-status`  | Show goose migration status.                                        |
| `make test`            | Run the Go test suite.                                              |
| `make check`           | Run `golangci-lint`.                                                |
| `make docker-build`    | Build the `logline:dev` image standalone.                           |
| `make docker-up`       | `docker compose up -d --build`.                                     |
| `make docker-ps`       | Show running compose services.                                      |
| `make docker-shell`    | Shell into the running `app` container.                             |
| `make docker-rebuild`  | Rebuild images with `--no-cache`.                                   |
| `make docker-clean`    | Tear down compose stack and remove local images.                    |
| `make compose-build`   | Build compose services.                                             |
| `make compose-up`      | Start the full stack detached.                                      |
| `make compose-down`    | Stop and remove the full stack.                                     |
| `make compose-logs`    | Tail logs for all compose services.                                  |

## Deployment

The `Dockerfile` uses a two-stage build:

1. **Builder stage** (`golang:1.27-alpine`) — downloads modules, compiles a stripped, trimmed,
   statically-linked (`CGO_ENABLED=0`) binary.
2. **Runtime stage** (`alpine:3.24`) — copies only the compiled binary, `migrations/`, and
   `web/` (templates + static assets) into a minimal image, creates a dedicated non-root
   `logline` user/group, and runs the process as that user on port `8080`.

`compose.yaml` orchestrates the app alongside a `postgres:18-alpine` service with a
`pg_isready` healthcheck gate, so the app container only starts once the database is actually
ready to accept connections — not just once the container process has started.

For production use beyond local Compose, typical next steps are: externalizing `DB_URL` and
secrets via your platform's secret manager, fronting the service with a reverse proxy/ingress
for TLS termination, and running migrations as a separate release step (`make migrate`) ahead
of a rolling deploy.

## Roadmap / Known Limitations

This project is under active development. Notable gaps an architect/reviewer should be aware
of before relying on it in production:

- **`config.Load()` doesn't actually fail fast — it can nil-pointer panic instead.** Three
  validation branches in `internal/config/config.go` (the empty-`DB_URL` check, the
  negative-`Burst`/`ReqPerSec` check, and the non-positive-timeout check) `return nil, err`
  where `err` is left over from an earlier, *successful* parse and is therefore `nil`. So on
  misconfiguration, `config.Load()` returns `(nil, nil)`; `main.go`'s `if err != nil` guard
  never fires; the nil `*Config` gets passed into `server.NewServer`; and the process panics
  on the first field access instead of printing the clean, actionable error message the log
  line right before it implies. Fix: replace `err` with a fresh `errors.New(...)` /
  `fmt.Errorf(...)` on each of those three `return` statements.
- **`api_keys` table exists but isn't wired up yet** — key generation/verification primitives
  are implemented in `internal/auth`, but there is no `APIKeyStore`, handler, or middleware
  consuming them yet. Today, `/ingest` is protected by the same session-token auth as the
  dashboard, not by service-scoped API keys.
- **`GET /dashboard` is currently unauthenticated** (registered under `publicChain`), while
  `GET /logs` (the equivalent JSON data) requires a session. Worth aligning before exposing the
  dashboard beyond a trusted network. Compounding this, its `limit` query parameter has no
  upper bound (unlike `/logs`, which is capped at 100), so it's also a resource-exhaustion
  vector against Postgres while it stays public.
- **`DashboardHandler.Stats`** is an unimplemented stub (`internal/handler/dashboard.go`) and,
  separately, is never registered as a route in `internal/server/routes.go` — it's currently
  dead code either way.
- **`model.User.PasswordHash` lacks a `json:"-"` tag** — not exploited by any handler today,
  but a defensive fix worth making before any "get current user" endpoint is added.
- **Rate limiter keys off `r.RemoteAddr`**, not a proxy-aware header, so behind a reverse
  proxy/load balancer every client currently looks like the same IP unless
  `X-Forwarded-For`/`X-Real-IP` is read and trusted appropriately.
- **No session-store lookup index note**: `sessions.token_hash` is unique-indexed, but there is
  no background job yet to purge expired sessions beyond the ad-hoc
  `SessionStore.DeleteExpiredSessions` method — it isn't currently scheduled anywhere.
- **CI does not run tests** — `go test` is commented out in the GitHub Actions workflow,
  pending a provisioned test database in CI.
- **No `LICENSE` file** is currently present in the repository.
- **No rotation/pruning strategy for the `logs` table** — there's no TTL, partitioning, or
  archival job yet; at high ingest volume this table will grow unbounded.

## Contributing

1. Fork the repository and create a feature branch.
2. Run `make check` and `make test` before opening a PR.
3. Keep migrations additive and reversible (`-- +goose Up` / `-- +goose Down`) — never edit a
   migration that has already shipped.
4. Open a pull request against `main`; the Go build workflow will run automatically.
