"use client";

import { useTransition } from "react";

import { deleteRepoAction } from "@/app/actions";
import type { Repo } from "@/lib/types";
import { Button } from "@/components/ui/button";
import { Modal } from "@/components/ui/modal";
import { useToast } from "@/components/ui/toast";

interface DeleteRepoDialogProps {
  /** The repo pending deletion, or null when the dialog is closed. */
  repo: Repo | null;
  onClose: () => void;
}

export function DeleteRepoDialog({ repo, onClose }: DeleteRepoDialogProps) {
  const toast = useToast();
  const [isDeleting, startDelete] = useTransition();

  const handleDelete = () => {
    if (!repo) return;
    startDelete(async () => {
      const result = await deleteRepoAction(repo.id);
      if (result.ok) {
        toast.success(`Removed ${repo.full_name}`);
        onClose();
      } else {
        toast.error(result.error);
      }
    });
  };

  return (
    <Modal
      open={repo !== null}
      onClose={onClose}
      title="Delete repository"
      description="This stops tracking the repository. It does not affect anything on GitHub."
      footer={
        <>
          <Button variant="outline" onClick={onClose}>
            Cancel
          </Button>
          <Button variant="destructive" loading={isDeleting} onClick={handleDelete}>
            Delete
          </Button>
        </>
      }
    >
      <p className="text-sm text-slate-600">
        Are you sure you want to stop tracking{" "}
        <span className="font-semibold text-slate-900">{repo?.full_name}</span>?
      </p>
    </Modal>
  );
}
