// Server-side API client for the Go backend.
//
// This module is imported only by Server Components and Server Actions, so it
// runs on the Next.js server and reaches the backend over the internal network
// (API_BASE_URL → http://backend:12010 in compose). The browser never calls the
// Go service directly, which keeps everything same-origin and avoids CORS.

import type { CreateRepoInput, Repo } from "./types";

const API_BASE_URL = process.env.API_BASE_URL ?? "http://localhost:12010";

/** The backend's error envelope: { "error": { code, message, details } }. */
interface ErrorEnvelope {
  error?: { code?: string; message?: string };
}

/** Carries the backend's human-readable message and stable error code. */
export class ApiError extends Error {
  constructor(
    message: string,
    readonly code: string,
    readonly status: number,
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API_BASE_URL}${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...init?.headers },
    // Always hit the backend; freshness is driven by revalidatePath in actions.
    cache: "no-store",
  });

  // 204 No Content (DELETE) has an empty body.
  if (res.status === 204) return undefined as T;

  const text = await res.text();
  const payload = text ? (JSON.parse(text) as unknown) : null;

  if (!res.ok) {
    const envelope = payload as ErrorEnvelope | null;
    throw new ApiError(
      envelope?.error?.message ?? `Request failed (${res.status})`,
      envelope?.error?.code ?? "UNKNOWN",
      res.status,
    );
  }

  return payload as T;
}

export function listRepos(language?: string): Promise<Repo[]> {
  const query = language ? `?language=${encodeURIComponent(language)}` : "";
  return request<Repo[]>(`/api/repos${query}`);
}

export function createRepo(input: CreateRepoInput): Promise<Repo> {
  return request<Repo>("/api/repos", {
    method: "POST",
    body: JSON.stringify(input),
  });
}

export function updateNotes(id: number, notes: string): Promise<Repo> {
  return request<Repo>(`/api/repos/${id}`, {
    method: "PATCH",
    body: JSON.stringify({ notes }),
  });
}

export function refreshRepo(id: number): Promise<Repo> {
  return request<Repo>(`/api/repos/${id}/refresh`, { method: "POST" });
}

export function deleteRepo(id: number): Promise<void> {
  return request<void>(`/api/repos/${id}`, { method: "DELETE" });
}
