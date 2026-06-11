import { forwardRef } from "react";
import type { InputHTMLAttributes } from "react";

import { cn } from "./cn";

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  invalid?: boolean;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className, invalid, ...props }, ref) => (
    <input
      ref={ref}
      aria-invalid={invalid || undefined}
      className={cn(
        "h-10 w-full rounded-lg border bg-white px-3 text-sm text-slate-900 shadow-sm transition-colors",
        "placeholder:text-slate-400",
        "focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-0",
        "disabled:cursor-not-allowed disabled:opacity-50",
        invalid
          ? "border-red-400 focus-visible:outline-red-500"
          : "border-slate-300 focus-visible:outline-indigo-500",
        className,
      )}
      {...props}
    />
  ),
);

Input.displayName = "Input";
