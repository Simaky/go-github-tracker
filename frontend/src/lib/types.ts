// Shared domain types for the frontend. These mirror the JSON shapes the Go
// backend returns (see backend/app/types.go) — kept deliberately small.

/** A tracked GitHub repository, as returned by the API. */
export interface Repo {
  id: number;
  owner: string;
  name: string;
  full_name: string;
  description: string;
  stars: number;
  language: string;
  html_url: string;
  notes: string;
  forks_count: number;
  fetched_at: string;
  created_at: string;
  updated_at: string;
}

/** Body of POST /api/repos. */
export interface CreateRepoInput {
  owner: string;
  name: string;
}

/**
 * Discriminated result returned by every Server Action. Actions never throw
 * across the network boundary; they resolve to one of these so the client can
 * render a toast without try/catch.
 */
export type ActionResult<T = undefined> =
  | { ok: true; data: T }
  | { ok: false; error: string };
