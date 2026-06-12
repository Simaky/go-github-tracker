# go-github-tracker

<img width="1381" height="902" alt="image" src="https://github.com/user-attachments/assets/225ac500-b7d6-485f-85fa-486d083ac882" />

Track GitHub repositories: add one by `owner/name`, the backend fetches its
metadata from the GitHub public API and stores it; then list, filter, annotate
with notes, refresh, and remove repos.

Monorepo — a **Go** backend API and a **Next.js** frontend, run together with
Docker Compose (+ PostgreSQL).

## How to run

Needs Docker. From the repo root:

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

The defaults in `.env.example` work as-is. To run a part on its own instead of
the whole stack: `make backend-run` (needs a local Postgres + Go 1.26+, configured
via `backend/.env`) or `make frontend-dev` (needs Node 22+).

### API

| Method | Path                     | Body / Query                  |
|--------|--------------------------|-------------------------------|
| POST   | `/api/repos`             | `{ "owner": "", "name": "" }` |
| GET    | `/api/repos`             | `?language=Go` (optional)     |
| GET    | `/api/repos/:id`         | —                             |
| PATCH  | `/api/repos/:id`         | `{ "notes": "" }`             |
| POST   | `/api/repos/:id/refresh` | —                             |
| DELETE | `/api/repos/:id`         | —                             |

## Architectural choices

```
backend/   Go API — Gin (HTTP) + Ent (ORM) over PostgreSQL (pgx)
frontend/  Next.js (App Router) + TypeScript + Tailwind
```

**Backend** is layered; imports only ever point down:

```
main → server (Gin / handlers) → app (domain) → app/storage + services
```

- `app.App` orchestrates the domain and holds its collaborators as **interfaces
  it declares itself** — so it is unit-testable with stubs and never forms an
  import cycle. The concrete types (Ent-backed storage, GitHub client) are
  constructed once in `main` and injected.
- **Persistence:** Ent over PostgreSQL via the `pgx` driver. The schema lives in
  `backend/ent/schema/`; Ent generates the typed client and auto-migrates on
  startup.
- **HTTP:** Gin. Handlers are thin translators (decode → validate → call the app
  → map the error → write JSON); clients only ever see a sanitised error
  envelope. Routes are mounted under `/api`.
- **GitHub:** a small client in `services/github/` over `*http.Client`. An
  optional `GITHUB_TOKEN` raises the rate limit.

> The backend layout comes from a personal `go-template-ms` standard. Two
> defaults were swapped for this assignment's required stack: **Gin** (vs chi)
> and **Ent** (vs sqlx).

**Frontend** — the browser talks only to Next.js. The initial list is fetched in
a Server Component; every mutation goes through a **Server Action** that proxies
to the backend over the internal network and calls `revalidatePath` so the page
re-renders with fresh data. There are no direct browser → Go calls, so there is
no CORS to manage. The UI is built on plain Tailwind (no component library).

**Tests & CI:** the backend has unit tests for the GitHub client (against an
`httptest` server) and the app orchestration (with stubbed collaborators).
GitHub Actions run on each push: `backend-ci` (go test + golangci-lint) and
`frontend-ci` (eslint + typecheck + build).

## AI tools used

Built with **Claude Code** (Anthropic's CLI) as a pair-programming assistant:

- Scaffolded the monorepo and the Go service from my personal `go-template-ms`
  standard (packaged as a Claude Code plugin), keeping the layout consistent.
- Generated boilerplate (config, error envelope, middleware, Dockerfiles,
  compose, Ent schema, GitHub client, handlers) and the whole frontend (Tailwind
  UI, Server Actions, dialogs/toasts).
- Reviewed structure, wiring, and the diff as the work progressed — and found a
  real bug along the way (Ent was opening the `postgres` sql driver instead of
  the `pgx` one the binary registers).

Every change was reviewed before staging; commits are made by hand in small,
meaningful steps (see the commit history).
