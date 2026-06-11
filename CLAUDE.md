# CLAUDE.md

Context for AI assistants working in this repo. Keep it short and current.

## What this is

A take-home full-stack app: **track GitHub repositories**. Add by `owner/name`,
the backend fetches metadata from the GitHub public API, stores it, and lets you
list/filter/annotate/refresh/delete. Monorepo: Go backend + Next.js frontend.

The active plan lives in **[docs/ROADMAP.md](docs/ROADMAP.md)** — read it to see
what's done and what's next. Product docs are in [README.md](README.md).

## Stack

- **Backend:** Go 1.26, **Gin** (HTTP), **Ent** (ORM) over **PostgreSQL** (`pgx`),
  `caarlos0/env` + `godotenv` (config).
- **Frontend:** Next.js 15 (App Router) + React 19 + TypeScript + Tailwind v4.
- **Infra:** Docker Compose + PostgreSQL 16.

## Layout (follow the `go-template-ms` plugin standard)

Place new backend code by concern — do not invent new top-level dirs:

| Concern                              | Goes in                              |
|--------------------------------------|--------------------------------------|
| Ent entity schemas                   | `backend/ent/schema/`                |
| DB access (Ent queries)             | `backend/app/storage/`               |
| Domain logic / orchestration         | `backend/app/` (the `App` type)      |
| Request/response types + `Validate()`| `backend/app/types.go`               |
| Interfaces the app depends on        | `backend/app/interfaces.go`          |
| Outbound clients (GitHub)            | `backend/services/github/`           |
| HTTP handlers (one file per resource)| `backend/server/handlers/`           |
| Routes / middleware                  | `backend/server/`                    |

- The `App` holds collaborators as **interfaces it declares itself**; concretes are
  wired only in `go-github-tracker-ms.go` (config → storage → app → server → run).
- HTTP routes are mounted under **`/api`** (this assignment), not `/v1`.
- Handlers stay thin: decode → validate → call app → translate error → write JSON.

## Build & verify

Backend (this machine needs the toolchain env — Go 1.23 installed, `GOSUMDB=off`):

```sh
cd backend && GOTOOLCHAIN=go1.26.0 GOSUMDB=sum.golang.org go build ./...
# the Makefile already exports these for `make run|test|lint`
```

After changing `ent/schema/`, regenerate: `cd backend && make ent-generate`.
Frontend: `cd frontend && npm run build`.

## Working rules

- **Commit history matters — do NOT squash.** Work in small, meaningful steps; one
  milestone item ≈ one (or a few) commits.
- **Never run `git commit` or `git push`.** Stage with `git add` and stop — the user
  commits. (User preference.)
- Build + verify before staging. Keep changes aligned with the plugin layout above.
