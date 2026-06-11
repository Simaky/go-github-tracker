"use client";

import { useState, useTransition } from "react";
import type { FormEvent } from "react";

import { trackRepoAction } from "@/app/actions";
import { Button } from "@/components/ui/button";
import { Field } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { Modal } from "@/components/ui/modal";
import { useToast } from "@/components/ui/toast";

interface AddRepoDialogProps {
  open: boolean;
  onClose: () => void;
}

/** Parse "owner/name" into its two segments, or null if malformed. */
function parseSlug(value: string): { owner: string; name: string } | null {
  const parts = value.trim().split("/");
  if (parts.length !== 2) return null;
  const [owner, name] = parts.map((p) => p.trim());
  if (!owner || !name) return null;
  return { owner, name };
}

export function AddRepoDialog({ open, onClose }: AddRepoDialogProps) {
  const toast = useToast();
  const [value, setValue] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, startSubmit] = useTransition();

  const close = () => {
    setValue("");
    setError(null);
    onClose();
  };

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    const parsed = parseSlug(value);
    if (!parsed) {
      setError('Enter the repository as "owner/name", e.g. facebook/react.');
      return;
    }
    setError(null);
    startSubmit(async () => {
      const result = await trackRepoAction(parsed);
      if (result.ok) {
        toast.success(`Now tracking ${parsed.owner}/${parsed.name}`);
        close();
      } else {
        setError(result.error);
      }
    });
  };

  return (
    <Modal
      open={open}
      onClose={close}
      title="Add repository"
      description="Paste a GitHub repository to start tracking it."
    >
      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <Field
          htmlFor="repo-slug"
          label="Repository"
          description='Format: owner/name (e.g. "vercel/next.js").'
          error={error ?? undefined}
        >
          <Input
            id="repo-slug"
            autoFocus
            placeholder="facebook/react"
            value={value}
            invalid={!!error}
            onChange={(e) => setValue(e.target.value)}
          />
        </Field>

        <div className="flex justify-end gap-3">
          <Button type="button" variant="outline" onClick={close}>
            Cancel
          </Button>
          <Button type="submit" loading={isSubmitting}>
            Track repository
          </Button>
        </div>
      </form>
    </Modal>
  );
}
