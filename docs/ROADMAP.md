# Roadmap

Incremental plan, built in small commits (see git history). Each milestone is a
self-contained, reviewable slice. Check items off as they land.

## M0 — Scaffold ✅

- [x] Monorepo: root compose, Makefile, `.env.example`, READMEs
- [x] Backend skeleton on the `go-template-ms` layout (Gin, Ent driver over pgx, config, graceful shutdown) — builds & vets clean
- [x] Frontend skeleton (Next.js + React + TS + Tailwind) — builds clean
- [x] Config moved to `.env` (godotenv + env tags)

## M1 — Data model ✅

- [x] `Repo` Ent schema in `backend/ent/schema/repo.go`: `owner, name, full_name (unique),
      description, stars, language, html_url, notes, fetched_at` + `TimeMixin` for
      `created_at` / `updated_at`
- [x] `make ent-generate` → typed client committed
- [x] Enable auto-migrate in `app/storage`: wrap the driver in `*ent.Client`, call
      `client.Schema.Create(ctx)` on startup

## M2 — GitHub client ✅

- [x] `services/github/client.go`: `GET https://api.github.com/repos/{owner}/{name}`,
      map to a typed result, handle 404 / rate-limit, optional `GITHUB_TOKEN`
- [x] Map upstream errors to domain errors (not-found, upstream)

## M3 — Domain + storage ✅

- [x] `app/storage` Repo methods: create, list (with `language` filter), get, update
      notes, delete, refresh; translate "not found" / conflict to domain errors
- [x] `app/interfaces.go`: `Storager` + `GitHubClient` interfaces the app needs
- [x] `app/types.go`: `CreateRepoRequest`, `UpdateNotesRequest` + `Validate()`
- [x] `app/application.go`: `TrackRepo`, `ListRepos`, `GetRepo`, `UpdateNotes`,
      `DeleteRepo`, `RefreshRepo` (reject duplicates on create)

## M4 — HTTP API ✅

- [x] `server/handlers/repos.go`: handlers for the six endpoints
- [x] Wire routes under `/api/repos` in `server/http_server.go`
- [x] Manual smoke test against a running Postgres (fixed Ent+pgx driver wiring:
      open `sql.Open("pgx", …)` + `entsql.OpenDB`, not `entsql.Open("postgres", …)`)

## M5 — Frontend ✅

Plain Tailwind v4 (no component library). The browser talks only to Next.js;
Server Actions proxy to the backend over `API_BASE_URL`, so there is no CORS.

- [x] API client (typed fetch wrappers) pointed at `API_BASE_URL` (`src/lib/api.ts`,
      server-only) + Server Actions in `src/app/actions.ts` with `revalidatePath`
- [x] Repo table: name, stars, language, description, GitHub link
- [x] Add-repo form (owner / name) — modal dialog, parses `owner/name`
- [x] Edit notes — modal dialog
- [x] Delete (confirm dialog) + per-row refresh actions
- [x] Toasts, empty state, search + language filter, stats summary

## M6 — Tests & polish

- [x] GitHub client test against an `httptest` mock server
- [x] App-level tests (track/refresh orchestration, duplicate rejection, validation)
- [x] Single root README: how to run, architecture choices, AI tools
- [x] `frontend-ci` workflow (eslint + typecheck + build) alongside `backend-ci`
- [x] Verify full stack via `make up` (postgres healthy, SSR renders live data)
