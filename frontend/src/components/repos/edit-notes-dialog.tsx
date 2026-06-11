"use client";

import { useEffect, useState, useTransition } from "react";
import type { FormEvent } from "react";

import { updateNotesAction } from "@/app/actions";
import type { Repo } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Field } from "@/components/ui/field";
import { Modal } from "@/components/ui/modal";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/toast";

interface EditNotesDialogProps {
  /** The repo being edited, or null when the dialog is closed. */
  repo: Repo | null;
  onClose: () => void;
}

export function EditNotesDialog({ repo, onClose }: EditNotesDialogProps) {
  const toast = useToast();
  const [notes, setNotes] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSaving, startSave] = useTransition();

  // Re-seed the textarea whenever a different repo opens the dialog.
  useEffect(() => {
    if (repo) {
      setNotes(repo.notes);
      setError(null);
    }
  }, [repo]);

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    if (!repo) return;
    startSave(async () => {
      const result = await updateNotesAction(repo.id, notes);
      if (result.ok) {
        toast.success(`Saved notes for ${repo.full_name}`);
        onClose();
      } else {
        setError(result.error);
      }
    });
  };

  return (
    <Modal
      open={repo !== null}
      onClose={onClose}
      title="Edit notes"
      description={repo ? repo.full_name : undefined}
    >
      <form onSubmit={handleSubmit} className="flex flex-col gap-5">
        <Field
          htmlFor="repo-notes"
          label="Notes"
          description="Your private notes about this repository."
          error={error ?? undefined}
        >
          <Textarea
            id="repo-notes"
            autoFocus
            rows={5}
            placeholder="Why are you tracking this repo?"
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
          />
        </Field>

        <div className="flex justify-end gap-3">
          <Button type="button" variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" loading={isSaving}>
            Save notes
          </Button>
        </div>
      </form>
    </Modal>
  );
}
