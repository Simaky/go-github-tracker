import { Button } from "@/components/ui/button";
import { GithubIcon, PlusIcon } from "@/components/ui/icons";

/** Shown when no repositories are tracked yet — friendly first-run CTA. */
export function EmptyState({ onAdd }: { onAdd: () => void }) {
  return (
    <div className="flex flex-col items-center justify-center rounded-xl border border-dashed border-slate-300 bg-white px-6 py-16 text-center shadow-sm">
      <div className="flex size-14 items-center justify-center rounded-full bg-slate-100 text-slate-500">
        <GithubIcon className="size-7" />
      </div>
      <h3 className="mt-5 text-lg font-semibold text-slate-900">No repositories yet</h3>
      <p className="mt-1 max-w-sm text-sm text-slate-500">
        Track your first GitHub repository to pull in its stars, language, and
        description — then annotate and refresh it any time.
      </p>
      <Button onClick={onAdd} className="mt-6">
        <PlusIcon className="size-4" />
        Add repository
      </Button>
    </div>
  );
}
