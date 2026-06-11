# Roadmap

Incremental plan, built in small commits (see git history). Each milestone is a
self-contained, reviewable slice. Check items off as they land.

## M0 — Scaffold ✅

- [x] Monorepo: root compose, Makefile, `.env.example`, READMEs
- [x] Backend skeleton on the `go-template-ms` layout (Gin, Ent driver over pgx, config, graceful shutdown) — builds & vets clean
- [x] Frontend skeleton (Next.js + React + TS + Tailwind) — builds clean
- [x] Config moved to `.env` (godotenv + env tags)

## M1 — Data model

- [ ] `Repo` Ent schema in `backend/ent/schema/repo.go`: `owner, name, full_name (unique),
      description, stars, language, html_url, notes, fetched_at` + Ent mixins for
      `created_at` / `updated_at`
- [ ] `make ent-generate` → typed client committed
- [ ] Enable auto-migrate in `app/storage`: wrap the driver in `*ent.Client`, call
      `client.Schema.Create(ctx)` on startup

## M2 — GitHub client

- [ ] `services/github/client.go`: `GET https://api.github.com/repos/{owner}/{name}`,
      map to a typed result, handle 404 / rate-limit, optional `GITHUB_TOKEN`
- [ ] Map upstream errors to domain errors (not-found, upstream)

## M3 — Domain + storage

- [ ] `app/storage` Repo methods: create, list (with `language` filter), get, update
      notes, delete, upsert-on-refresh; map "not found" to a sentinel
- [ ] `app/interfaces.go`: `Storager` + `GitHubClient` interfaces the app needs
- [ ] `app/types.go`: `CreateRepoRequest`, `UpdateNotesRequest` + `Validate()`
- [ ] `app/application.go`: `TrackRepo`, `ListRepos`, `GetRepo`, `UpdateNotes`,
      `DeleteRepo`, `RefreshRepo` (reject duplicates on create)

## M4 — HTTP API

- [ ] `server/handlers/repos.go`: handlers for the six endpoints
- [ ] Wire routes under `/api/repos` in `server/http_server.go`
- [ ] Manual smoke test against a running Postgres

## M5 — Frontend

- [ ] API client (typed fetch wrappers) pointed at `API_BASE_URL`
- [ ] Repo table: name, stars, language, description, GitHub link
- [ ] Add-repo form (owner / name)
- [ ] Edit notes (inline or modal)
- [ ] Delete + per-row refresh actions

## M6 — Tests & polish

- [ ] GitHub client test against an `httptest` mock server
- [ ] One app/handler-level test (e.g. duplicate rejection or the language filter)
- [ ] Finalise README: architecture write-up + AI-tools specifics
- [ ] Verify full stack via `make up`
