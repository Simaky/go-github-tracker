import * as api from "@/lib/api";
import type { Repo } from "@/lib/types";
import { GithubIcon } from "@/components/ui/icons";
import { ReposClient } from "@/components/repos/repos-client";
import { StatsBar } from "@/components/repos/stats-bar";

// Always render fresh: revalidatePath in the Server Actions drives updates.
export const dynamic = "force-dynamic";

export default async function Home() {
  let repos: Repo[] = [];
  let loadError: string | null = null;

  try {
    repos = await api.listRepos();
  } catch {
    loadError = "Couldn’t reach the API. Make sure the backend is running.";
  }

  return (
    <div className="min-h-screen">
      <header className="border-b border-slate-800 bg-slate-900">
        <div className="mx-auto flex max-w-5xl items-center gap-3 px-6 py-4">
          <div className="flex size-9 items-center justify-center rounded-lg bg-white/10 text-white">
            <GithubIcon className="size-5" />
          </div>
          <span className="text-lg font-semibold tracking-tight text-white">
            GitHub Repository Tracker
          </span>
        </div>
      </header>

      <main className="mx-auto flex max-w-5xl flex-col gap-8 px-6 py-10">
        <div>
          <h1 className="text-2xl font-bold tracking-tight text-slate-900">
            Tracked repositories
          </h1>
          <p className="mt-1 text-sm text-slate-500">
            Add a repository by <code className="text-slate-700">owner/name</code>; we pull
            its metadata from GitHub so you can filter, annotate, and refresh it.
          </p>
        </div>

        {loadError ? (
          <div className="rounded-xl border border-red-200 bg-red-50 p-6 text-sm text-red-700">
            {loadError}
          </div>
        ) : (
          <>
            <StatsBar repos={repos} />
            <ReposClient repos={repos} />
          </>
        )}
      </main>
    </div>
  );
}
