import { useState, useEffect } from "react";
import { Modal } from "../Modal";
import { Button } from "../Button";
import type { Room } from "./types";

type RoomDeleteModalProps = Readonly<{
  open: boolean;
  room: Room | null;
  affectedCharacterCount?: number;
  orphanExitCount?: number;
  onConfirm: () => void;
  onCancel: () => void;
  isLoading?: boolean;
}>;

export function RoomDeleteModal({
  open,
  room,
  affectedCharacterCount = 0,
  orphanExitCount = 0,
  onConfirm,
  onCancel,
  isLoading,
}: RoomDeleteModalProps) {
  const [confirmText, setConfirmText] = useState("");

  useEffect(() => {
    if (open) setConfirmText("");
  }, [open]);

  if (!room) return null;

  const canDelete = confirmText === room.name && !isLoading;

  return (
    <Modal isOpen={open} onClose={onCancel} title={`Delete Room: ${room.name}`}>
      <div className="space-y-3">
        <p className="text-text text-sm">
          You are about to permanently delete <strong>{room.name}</strong> (ID: #{room.id}).
        </p>
        <div className="bg-surface-muted p-2 rounded text-text-muted text-xs">
          <div>This will affect:</div>
          <ul className="list-disc list-inside ml-2">
            <li>{affectedCharacterCount} character(s) currently in this room will be relocated</li>
            <li>{orphanExitCount} exit(s) in other rooms will be cleaned up</li>
          </ul>
        </div>
        <p className="text-text-muted text-xs">
          This action cannot be undone. Type <strong className="text-text">{room.name}</strong> to confirm:
        </p>
        <input
          type="text"
          value={confirmText}
          onChange={(e) => setConfirmText(e.target.value)}
          className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
          autoFocus
          disabled={isLoading}
          aria-label="Type the room name to confirm"
        />
        <div className="flex gap-2 justify-end">
          <Button variant="secondary" onClick={onCancel} type="button" disabled={isLoading}>
            Cancel
          </Button>
          <Button variant="danger" onClick={onConfirm} type="button" disabled={!canDelete}>
            {isLoading ? "Deleting…" : "Delete"}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
