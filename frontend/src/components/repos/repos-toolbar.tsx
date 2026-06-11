import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { PlusIcon, SearchIcon } from "@/components/ui/icons";

interface ReposToolbarProps {
  query: string;
  onQueryChange: (value: string) => void;
  language: string;
  onLanguageChange: (value: string) => void;
  languages: string[];
  onAdd: () => void;
}

/** Search box + language filter on the left, "Add repository" on the right. */
export function ReposToolbar({
  query,
  onQueryChange,
  language,
  onLanguageChange,
  languages,
  onAdd,
}: ReposToolbarProps) {
  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div className="flex flex-1 flex-col gap-3 sm:flex-row sm:items-center">
        <div className="relative sm:max-w-xs sm:flex-1">
          <SearchIcon className="pointer-events-none absolute top-1/2 left-3 size-4 -translate-y-1/2 text-slate-400" />
          <Input
            type="search"
            placeholder="Search repositories…"
            aria-label="Search repositories"
            value={query}
            onChange={(e) => onQueryChange(e.target.value)}
            className="pl-9"
          />
        </div>

        <select
          aria-label="Filter by language"
          value={language}
          onChange={(e) => onLanguageChange(e.target.value)}
          className="h-10 rounded-lg border border-slate-300 bg-white px-3 text-sm text-slate-900 shadow-sm focus-visible:outline focus-visible:outline-2 focus-visible:outline-indigo-500"
        >
          <option value="">All languages</option>
          {languages.map((lang) => (
            <option key={lang} value={lang}>
              {lang}
            </option>
          ))}
        </select>
      </div>

      <Button onClick={onAdd}>
        <PlusIcon className="size-4" />
        Add repository
      </Button>
    </div>
  );
}
