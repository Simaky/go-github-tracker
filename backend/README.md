# go-github-tracker-ms

The backend service for **go-github-tracker**, built to the `go-template-ms`
standard. It exposes a REST API to track GitHub repositories: the `Repo` entity
(Ent + PostgreSQL), the GitHub API client, the domain orchestrator, and the
`/api/repos` handlers sit on top of the standard layered layout, composition
root, config-from-env, error envelope, and graceful shutdown.

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

Targets Go **1.26**, set by the `go` directive in `go.mod`. With the default
`GOTOOLCHAIN=auto`, an older local Go auto-downloads the right toolchain. (Edge
case: if you've globally set `GOSUMDB=off`, that download fails its checksum
verification — install a recent Go, or prefix the command with
`GOSUMDB=sum.golang.org`.)

## Layout

```
go-github-tracker-ms.go   # package main — entry point, wiring only
app/                      # domain layer (App orchestrator, interfaces, types, errors)
  config/                 # Config + LoadConfig (env / .env)
  storage/                # Ent-backed Repo persistence over pgx (storage.go, repo.go)
ent/                      # Ent ORM: schema/repo.go + generated typed client
server/                   # Gin engine, handlers (repos.go), middleware, graceful shutdown
services/                 # outbound clients to external systems
  github/                 # GitHub REST API client
consts/                   # ServiceName
```

## Configuration

Config is read from the environment (in `app/config`): a `.env` file is loaded
if present, then env vars are decoded into the typed `Config`, then defaults are
applied and validated. **Real exported env vars always win** over `.env`, and a
missing `.env` is fine (in Docker, compose injects the variables directly).

| Variable                    | Meaning                              | Example                                                   |
|-----------------------------|--------------------------------------|-----------------------------------------------------------|
| `GOGITHUBTRACKER_DB_DSN`    | PostgreSQL DSN (**required**)        | `postgres://app:app@localhost:5432/app?sslmode=disable`   |
| `GOGITHUBTRACKER_LISTEN`    | bind address (default `:12010`)      | `0.0.0.0:12010`                                           |
| `GITHUB_TOKEN`              | GitHub API token (optional)          | `ghp_…` (raises the API rate limit)                       |

For local (non-Docker) runs, copy the example env file and start the service —
`.env` is loaded automatically:

```sh
cp .env.example .env   # adjust the DSN if needed
make run
```

`.env` is gitignored (it holds local DSNs/secrets); `.env.example` is committed.

## Endpoints

Resource API (mounted under `/api`):

| Method | Path                      | Description                                          |
|--------|---------------------------|------------------------------------------------------|
| POST   | `/api/repos`              | Fetch from GitHub, persist, return record (409 dup). |
| GET    | `/api/repos`              | List tracked repos (`?language=` filter).            |
| GET    | `/api/repos/:id`          | Return a single repo.                                |
| PATCH  | `/api/repos/:id`          | Update the user-editable notes.                      |
| POST   | `/api/repos/:id/refresh`  | Re-fetch from GitHub and update stored fields.       |
| DELETE | `/api/repos/:id`          | Remove from the watchlist.                           |

Operational: `GET /uptime` (liveness), `GET /version`, `GET /debug/pprof/*`.

## Schema changes

The `Repo` entity lives in `ent/schema/repo.go`; the typed client under `ent/` is
generated from it and committed. After editing the schema, regenerate with
`make ent-generate`. Auto-migration (`client.Schema.Create`) runs on startup, so
the table is created/upgraded automatically.

## Make targets

`make run | build | test | lint | lint-install | ent-generate | tidy`
