"use client";

import { useMemo, useState } from "react";

import type { Repo } from "@/lib/types";

import { AddRepoDialog } from "./add-repo-dialog";
import { DeleteRepoDialog } from "./delete-repo-dialog";
import { EditNotesDialog } from "./edit-notes-dialog";
import { EmptyState } from "./empty-state";
import { ReposTable } from "./repos-table";
import { ReposToolbar } from "./repos-toolbar";

/**
 * Client shell for the repositories view. `repos` arrives from the server
 * component and is refreshed automatically whenever a Server Action calls
 * revalidatePath — so this component holds only UI state (search, filter,
 * which dialog is open), never a copy of the data.
 */
export function ReposClient({ repos }: { repos: Repo[] }) {
  const [query, setQuery] = useState("");
  const [language, setLanguage] = useState("");
  const [addOpen, setAddOpen] = useState(false);
  const [editing, setEditing] = useState<Repo | null>(null);
  const [deleting, setDeleting] = useState<Repo | null>(null);

  const languages = useMemo(
    () =>
      Array.from(new Set(repos.map((r) => r.language).filter(Boolean))).sort((a, b) =>
        a.localeCompare(b),
      ),
    [repos],
  );

  const filtered = useMemo(() => {
    const needle = query.trim().toLowerCase();
    return repos.filter((repo) => {
      if (language && repo.language !== language) return false;
      if (!needle) return true;
      return (
        repo.full_name.toLowerCase().includes(needle) ||
        repo.description.toLowerCase().includes(needle) ||
        repo.notes.toLowerCase().includes(needle)
      );
    });
  }, [repos, query, language]);

  return (
    <div className="flex flex-col gap-4">
      {repos.length === 0 ? (
        <EmptyState onAdd={() => setAddOpen(true)} />
      ) : (
        <>
          <ReposToolbar
            query={query}
            onQueryChange={setQuery}
            language={language}
            onLanguageChange={setLanguage}
            languages={languages}
            onAdd={() => setAddOpen(true)}
          />
          <ReposTable repos={filtered} onEdit={setEditing} onDelete={setDeleting} />
        </>
      )}

      <AddRepoDialog open={addOpen} onClose={() => setAddOpen(false)} />
      <EditNotesDialog repo={editing} onClose={() => setEditing(null)} />
      <DeleteRepoDialog repo={deleting} onClose={() => setDeleting(null)} />
    </div>
  );
}
