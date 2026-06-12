"use client";

import { useTransition } from "react";

import { refreshRepoAction } from "@/app/actions";
import type { Repo } from "@/lib/types";
import { formatCount, formatRelativeTime, formatStars } from "@/lib/format";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/ui/toast";
import {
  ExternalLinkIcon,
  ForkIcon,
  PencilIcon,
  RefreshIcon,
  StarIcon,
  TrashIcon,
} from "@/components/ui/icons";

interface RepoRowProps {
  repo: Repo;
  onEdit: (repo: Repo) => void;
  onDelete: (repo: Repo) => void;
}

export function RepoRow({ repo, onEdit, onDelete }: RepoRowProps) {
  const toast = useToast();
  const [isRefreshing, startRefresh] = useTransition();

  const handleRefresh = () => {
    startRefresh(async () => {
      const result = await refreshRepoAction(repo.id);
      if (result.ok) toast.success(`Refreshed ${repo.full_name}`);
      else toast.error(result.error);
    });
  };

  return (
    <tr className="group border-t border-slate-100 transition-colors hover:bg-slate-50/70">
      <td className="px-4 py-3 align-top">
        <div className="flex flex-col gap-1">
          <a
            href={repo.html_url}
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex w-fit items-center gap-1.5 font-semibold text-slate-900 hover:text-indigo-600"
          >
            {repo.full_name}
            <ExternalLinkIcon className="size-3.5 text-slate-400" />
          </a>
          {repo.description && (
            <p className="max-w-xl text-sm text-slate-500">{repo.description}</p>
          )}
          {repo.notes && (
            <p className="mt-1 max-w-xl rounded-md bg-amber-50 px-2 py-1 text-xs text-amber-800 ring-1 ring-amber-100">
              {repo.notes}
            </p>
          )}
        </div>
      </td>

      <td className="px-4 py-3 align-top">
        {repo.language ? (
          <Badge variant="language">{repo.language}</Badge>
        ) : (
          <span className="text-sm text-slate-400">—</span>
        )}
      </td>

      <td className="px-4 py-3 align-top">
        <span className="inline-flex items-center gap-1.5 text-sm font-medium text-slate-700">
          <StarIcon className="size-4 text-amber-400" />
          {formatStars(repo.stars)}
        </span>
      </td>

      <td className="px-4 py-3 align-top">
        <span className="inline-flex items-center gap-1.5 text-sm font-medium text-slate-700">
          <ForkIcon className="size-4 text-slate-400" />
          {formatCount(repo.forks_count)}
        </span>
      </td>

      <td
        className="px-4 py-3 align-top whitespace-nowrap text-sm text-slate-500"
        suppressHydrationWarning
      >
        {formatRelativeTime(repo.fetched_at)}
      </td>

      <td className="px-4 py-3 align-top">
        <div className="flex items-center justify-end gap-1">
          <Button
            variant="ghost"
            size="icon"
            onClick={handleRefresh}
            loading={isRefreshing}
            aria-label={`Refresh ${repo.full_name}`}
            title="Refresh from GitHub"
          >
            {!isRefreshing && <RefreshIcon className="size-4" />}
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onEdit(repo)}
            aria-label={`Edit notes for ${repo.full_name}`}
            title="Edit notes"
          >
            <PencilIcon className="size-4" />
          </Button>
          <Button
            variant="ghost"
            size="icon"
            onClick={() => onDelete(repo)}
            aria-label={`Delete ${repo.full_name}`}
            title="Delete"
            className="text-slate-500 hover:bg-red-50 hover:text-red-600"
          >
            <TrashIcon className="size-4" />
          </Button>
        </div>
      </td>
    </tr>
  );
}
