import type { ReactNode } from "react";

interface FieldProps {
  /** Input id the label points at. */
  htmlFor: string;
  label: string;
  /** Optional hint shown under the label. */
  description?: ReactNode;
  /** Validation message; renders in red and marks the field. */
  error?: string;
  children: ReactNode;
}

/** Label + control + description/error wrapper for consistent form layout. */
export function Field({ htmlFor, label, description, error, children }: FieldProps) {
  return (
    <div className="flex flex-col gap-1.5">
      <label htmlFor={htmlFor} className="text-sm font-medium text-slate-700">
        {label}
      </label>
      {children}
      {error ? (
        <p className="text-xs text-red-600">{error}</p>
      ) : description ? (
        <p className="text-xs text-slate-500">{description}</p>
      ) : null}
    </div>
  );
}
