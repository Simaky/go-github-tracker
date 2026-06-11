# go-github-tracker

A small full-stack app to **track GitHub repositories you care about**. Add a repo
by `owner/name`, the backend fetches its metadata from the GitHub public API and
stores it; you can list, filter, annotate with notes, refresh, and remove repos.

Monorepo: a **Go** backend API and a **Next.js** frontend, run together with
docker-compose (+ PostgreSQL).

---

## Status

Built incrementally (see commit history). The **backend API is complete** — the
`Repo` model, GitHub client, domain layer, and all six `/api/repos` endpoints,
with unit tests for the GitHub client and the app orchestration. The frontend is
next — see **[docs/ROADMAP.md](docs/ROADMAP.md)** for the plan and progress.

---

## Architecture

```
go-github-tracker/
├── docker-compose.yml      # backend + frontend + postgres
├── backend/                # Go API — go-template-ms layout
│   ├── go-github-tracker-ms.go   # composition root: config → storage → app → server → run
│   ├── app/                # domain layer (orchestrator, interfaces, types, errors)
│   │   ├── config/         # env-based config (.env supported)
│   │   └── storage/        # DB access via Ent over pgx
│   ├── ent/                # Ent ORM: schema/ + generated client
│   ├── services/           # outbound clients to external systems (GitHub) — one sub-pkg each
│   └── server/             # Gin engine, handlers, middleware
└── frontend/               # Next.js + React + TypeScript + Tailwind
```

**Backend** follows a deliberate layering (top imports down, never the reverse):

```
main → server (Gin/handlers) → app (domain) → app/storage + services (interfaces)
```

- The `app.App` orchestrator holds its collaborators as **interfaces it declares
  itself**, so it is testable with stubs and never forms an import cycle. Concrete
  types (the Ent-backed storage, the GitHub client) are constructed once in `main`.
- **Persistence** is **Ent** over **PostgreSQL** (via the `pgx` driver). The schema
  lives in `backend/ent/schema/`; Ent generates the typed client and runs
  **auto-migration** on startup (`client.Schema.Create`).
- **HTTP** is **Gin**. Handlers are thin translators (decode → validate → call the
  app → translate error → write JSON); the client only ever sees a sanitised error
  envelope. Routes are mounted under **`/api`** (per the assignment).
- The **GitHub integration** lives in `services/github/` as a client over
  `*http.Client`, calling `https://api.github.com/repos/{owner}/{name}`. An optional
  `GITHUB_TOKEN` raises the rate limit when present.

> The layout and conventions come from a personal `go-template-ms` standard
> (consistent service shape across an estate). A couple of defaults were swapped
> for this assignment's required stack — **Gin** (vs chi) and **Ent** (vs sqlx);
> see [`backend/README.md`](backend/README.md) for the rationale.

**Frontend** is a Next.js (App Router) app in TypeScript with Tailwind. It talks to
the backend over HTTP only.

---

## API

| Method | Path                      | Body / Query                | Description                                        |
|--------|---------------------------|-----------------------------|----------------------------------------------------|
| POST   | `/api/repos`              | `{ "owner": "", "name": "" }` | Fetch from GitHub, persist, return record. 409 on duplicate. |
| GET    | `/api/repos`              | `?language=Go` (optional)   | List tracked repos, optionally filtered by language. |
| GET    | `/api/repos/:id`          | —                           | Return a single repo.                              |
| PATCH  | `/api/repos/:id`          | `{ "notes": "" }`           | Update the user-editable notes.                    |
| DELETE | `/api/repos/:id`          | —                           | Remove from the watchlist.                         |
| POST   | `/api/repos/:id/refresh`  | —                           | Re-fetch from GitHub and update stored fields.     |

**`Repo` fields:** `id, owner, name, full_name (unique), description, stars,
language, html_url, notes, fetched_at, created_at, updated_at`.

Operational endpoints: `GET /uptime`, `GET /version`.

---

## Configuration

All config is via environment variables. Copy the examples and adjust:

```sh
cp .env.example .env                 # stack (compose): Postgres + ports
cp backend/.env.example backend/.env # backend (local run): DSN, listen, GitHub token
```

| Variable                  | Used by | Meaning                                              |
|---------------------------|---------|------------------------------------------------------|
| `GOGITHUBTRACKER_DB_DSN`  | backend | PostgreSQL DSN (**required**)                        |
| `GOGITHUBTRACKER_LISTEN`  | backend | bind address (default `0.0.0.0:12010`)               |
| `GITHUB_TOKEN`            | backend | optional GitHub token for higher rate limits         |
| `POSTGRES_USER/PASSWORD/DB`, `*_PORT` | compose | database + host port mapping              |

---

## How to run

### Docker Compose (whole stack)

```sh
cp .env.example .env
make up        # build + start backend, frontend, postgres (detached)
make logs      # follow logs
make down      # stop
```

| Service     | URL                    |
|-------------|------------------------|
| Frontend    | http://localhost:3000  |
| Backend API | http://localhost:12010 |
| PostgreSQL  | localhost:5432         |

### Locally (without Docker)

Backend (needs a running Postgres and Go 1.26+ — see the toolchain note in
[`backend/README.md`](backend/README.md)):

```sh
cd backend && cp .env.example .env && make run
```

Frontend (Node 22+):

```sh
cd frontend && npm install && npm run dev
```

---

## Tech stack

- **Backend:** Go 1.26 (satisfies the 1.22+ requirement), Gin, Ent, PostgreSQL (pgx),
  `caarlos0/env` + `godotenv` for config.
- **Frontend:** Next.js 15, React 19, TypeScript, Tailwind CSS v4.
- **Infra:** Docker Compose, PostgreSQL 16.

---

## AI tools used

This project was built with **Claude Code** (Anthropic's CLI) as a pair-programming
assistant, used to:

- Scaffold the monorepo and the Go service against my personal `go-template-ms`
  project-structure standard (a Claude Code plugin), keeping the layout consistent.
- Generate boilerplate (config loading, error envelope, middleware, Dockerfiles,
  compose) and the Ent schema / GitHub client / handlers.
- Review structure and wiring as the work progressed.

Every change was reviewed before staging; commits are made by hand in small,
meaningful steps (see the commit history). _This section will be expanded with
specifics as the implementation lands._

---

## Tests

Backend (`cd backend && make test`):

- **GitHub client** (`services/github`) against an `httptest` mock server — success
  mapping, 404 → not-found, rate-limit → upstream, auth-header handling.
- **App orchestration** (`app`) with stubbed storage/GitHub — track, duplicate
  rejection, validation short-circuits, and the refresh flow.
