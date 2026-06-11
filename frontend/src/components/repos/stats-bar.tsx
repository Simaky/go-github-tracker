import type { ReactNode } from "react";

import type { Repo } from "@/lib/types";
import { formatStars } from "@/lib/format";
import { GithubIcon, StarIcon } from "@/components/ui/icons";

function Stat({ label, value, icon }: { label: string; value: string; icon: ReactNode }) {
  return (
    <div className="flex items-center gap-4 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex size-11 shrink-0 items-center justify-center rounded-lg bg-indigo-50 text-indigo-600">
        {icon}
      </div>
      <div>
        <div className="text-2xl font-semibold tracking-tight text-slate-900">{value}</div>
        <div className="text-sm text-slate-500">{label}</div>
      </div>
    </div>
  );
}

/** Three summary tiles: repositories tracked, total stars, distinct languages. */
export function StatsBar({ repos }: { repos: Repo[] }) {
  const totalStars = repos.reduce((sum, repo) => sum + repo.stars, 0);
  const languages = new Set(repos.map((repo) => repo.language).filter(Boolean));

  return (
    <div className="grid grid-cols-1 gap-4 sm:grid-cols-3">
      <Stat
        label="Repositories tracked"
        value={String(repos.length)}
        icon={<GithubIcon className="size-5" />}
      />
      <Stat
        label="Total stars"
        value={formatStars(totalStars)}
        icon={<StarIcon className="size-5" />}
      />
      <Stat
        label="Languages"
        value={String(languages.size)}
        icon={<span className="text-base font-semibold">{"{ }"}</span>}
      />
    </div>
  );
}
