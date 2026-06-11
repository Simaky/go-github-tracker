import type { Repo } from "@/lib/types";

import { RepoRow } from "./repo-row";

interface ReposTableProps {
  repos: Repo[];
  onEdit: (repo: Repo) => void;
  onDelete: (repo: Repo) => void;
}

const HEADERS = ["Repository", "Language", "Stars", "Updated", ""];

export function ReposTable({ repos, onEdit, onDelete }: ReposTableProps) {
  return (
    <div className="overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm">
      <div className="overflow-x-auto">
        <table className="w-full border-collapse text-left">
          <thead>
            <tr className="bg-slate-50 text-xs font-semibold tracking-wide text-slate-500 uppercase">
              {HEADERS.map((header, i) => (
                <th
                  key={header || "actions"}
                  className={i === HEADERS.length - 1 ? "px-4 py-3 text-right" : "px-4 py-3"}
                  scope="col"
                >
                  {header}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {repos.length === 0 ? (
              <tr>
                <td colSpan={HEADERS.length} className="px-4 py-12 text-center text-sm text-slate-500">
                  No repositories match your filters.
                </td>
              </tr>
            ) : (
              repos.map((repo) => (
                <RepoRow key={repo.id} repo={repo} onEdit={onEdit} onDelete={onDelete} />
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
}
