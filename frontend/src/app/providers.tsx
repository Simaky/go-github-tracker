"use client";

import type { ReactNode } from "react";

import { ToastProvider } from "@/components/ui/toast";

/** Client-side providers mounted once at the app root. */
export function Providers({ children }: { children: ReactNode }) {
  return <ToastProvider>{children}</ToastProvider>;
}
