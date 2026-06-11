import type { ReactNode } from "react";

import { cn } from "./cn";

type Variant = "neutral" | "language";

const VARIANTS: Record<Variant, string> = {
  neutral: "bg-slate-100 text-slate-700 ring-slate-200",
  language: "bg-indigo-50 text-indigo-700 ring-indigo-100",
};

interface BadgeProps {
  variant?: Variant;
  className?: string;
  children: ReactNode;
}

/** Small non-interactive status/label pill. */
export function Badge({ variant = "neutral", className, children }: BadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium ring-1 ring-inset",
        VARIANTS[variant],
        className,
      )}
    >
      {children}
    </span>
  );
}
