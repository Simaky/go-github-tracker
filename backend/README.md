# go-github-tracker-ms

The backend service for **go-github-tracker**, built to the `go-template-ms`
standard. It is a bare skeleton: the full layered layout, the composition root,
config-from-env, the error envelope, and graceful shutdown are all in
place, with no domain entities or endpoints yet.

## Stack

This service swaps three of the template's default libraries. The template's
*structure* is unchanged; only the libraries differ, with reasons:

| Concern     | Template default      | Here                          | Reason                |
|-------------|-----------------------|-------------------------------|-----------------------|
| HTTP router | `go-chi/chi`          | **`gin-gonic/gin`**           | Chosen stack          |
| Persistence | `jmoiron/sqlx` + MySQL| **`ent` ORM** over PostgreSQL | Chosen stack          |
| DB driver   | `go-sql-driver/mysql` | **`jackc/pgx/v5`** (stdlib)   | PostgreSQL            |
| Migrations  | `goose` + `migrate.sh`| **Ent auto-migrate** at start | Chosen approach       |

Everything else follows the template: `app/` (domain) → `app/storage` + `services/`,
`server/` (Gin + handlers + middleware), `consts/`, a single
`package main` entry (`go-github-tracker-ms.go`) wiring `config → storage → app →
server → run`.

## Go version

Targets the latest Go (**1.26**), pinned via the `toolchain` directive in
`go.mod`. If your machine has an older Go, the toolchain auto-downloads — but that
download verifies against the checksum database, so if you have `GOSUMDB=off` set
it will refuse. The `Makefile` re-enables `GOSUMDB=sum.golang.org` for its targets;
for raw `go` commands either set that env var once or install Go 1.26+.

## Layout

```
go-github-tracker-ms.go   # package main — entry point, wiring only
app/                      # domain layer (App orchestrator, interfaces, types, errors)
  config/                 # Config + LoadConfig (env / JSON file)
  storage/                # DB connection via Ent's SQL driver over pgx
ent/                      # Ent ORM: generate.go + schema/ (add entities here)
server/                   # Gin engine, handlers, middleware, graceful shutdown
services/                 # outbound clients to external systems (empty)
consts/                   # ServiceName
```

## Configuration

Config is read from the environment (in `app/config`): a `.env` file is loaded
if present, then env vars are decoded into the typed `Config`, then defaults are
applied and validated. **Real exported env vars always win** over `.env`, and a
missing `.env` is fine (in Docker, compose injects the variables directly).

| Variable                    | Meaning                          | Example                                                   |
|-----------------------------|----------------------------------|-----------------------------------------------------------|
| `GOGITHUBTRACKER_DB_DSN`    | PostgreSQL DSN (**required**)    | `postgres://app:app@localhost:5432/app?sslmode=disable`   |
| `GOGITHUBTRACKER_LISTEN`    | bind address (default `:12010`)  | `0.0.0.0:12010`                                           |

For local (non-Docker) runs, copy the example env file and start the service —
`.env` is loaded automatically:

```sh
cp .env.example .env   # adjust the DSN if needed
make run
```

`.env` is gitignored (it holds local DSNs/secrets); `.env.example` is committed.

## Endpoints

- `GET /uptime` — plain-text liveness.
- `GET /version` — build version string.
- `GET /debug/pprof/*` — pprof.
- `/v1/...` — versioned API; resource routes are added here as the domain grows.

## Adding your first entity (and enabling auto-migrate)

Ent generates its typed client from schemas, so there is nothing to migrate until
you add one:

1. Create `ent/schema/<entity>.go` declaring an `ent.Schema` type.
2. `make ent-generate` — regenerates the typed client into `ent/`.
3. In `app/storage/storage.go`, wrap the driver in the generated client and call
   `client.Schema.Create(ctx)` in `New` (the exact snippet is in that file's doc
   comment). This runs auto-migration on startup.
4. Add `Storager` methods in `app/interfaces.go`, implement them in
   `app/storage`, and expose handlers under `/v1`.

## Make targets

`make run | build | test | lint | lint-install | ent-generate | tidy`
