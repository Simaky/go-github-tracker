"use server";

// Server Actions: the only mutation entry points the client can call. Each one
// proxies to the backend via the server-side API client, revalidates the list
// so the page re-renders with fresh data, and returns an ActionResult instead
// of throwing — the client turns that into a toast.

import { revalidatePath } from "next/cache";

import * as api from "@/lib/api";
import { ApiError } from "@/lib/api";
import type { ActionResult, CreateRepoInput } from "@/lib/types";

function toMessage(err: unknown): string {
  if (err instanceof ApiError) return err.message;
  return "Something went wrong. Please try again.";
}

export async function trackRepoAction(
  input: CreateRepoInput,
): Promise<ActionResult> {
  try {
    await api.createRepo(input);
    revalidatePath("/");
    return { ok: true, data: undefined };
  } catch (err) {
    return { ok: false, error: toMessage(err) };
  }
}

export async function refreshRepoAction(id: number): Promise<ActionResult> {
  try {
    await api.refreshRepo(id);
    revalidatePath("/");
    return { ok: true, data: undefined };
  } catch (err) {
    return { ok: false, error: toMessage(err) };
  }
}

export async function updateNotesAction(
  id: number,
  notes: string,
): Promise<ActionResult> {
  try {
    await api.updateNotes(id, notes);
    revalidatePath("/");
    return { ok: true, data: undefined };
  } catch (err) {
    return { ok: false, error: toMessage(err) };
  }
}

export async function deleteRepoAction(id: number): Promise<ActionResult> {
  try {
    await api.deleteRepo(id);
    revalidatePath("/");
    return { ok: true, data: undefined };
  } catch (err) {
    return { ok: false, error: toMessage(err) };
  }
}
