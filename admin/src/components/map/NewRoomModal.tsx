/* eslint-disable functional/no-mixed-types */
import { useState, useEffect } from "react";
import { Modal } from "../Modal";
import { Button } from "../Button";
import { FormField } from "../fields/FormField";
import { TextareaField } from "../fields/TextareaField";
import type { Room } from "./types";

type NewRoomModalProps = Readonly<{
  open: boolean;
  parentRoom: Room | null;
  direction: string | null;
  onConfirm: (input: { name: string; description: string }) => void;
  onCancel: () => void;
  isLoading?: boolean;
}>;

export function NewRoomModal({
  open,
  parentRoom,
  direction,
  onConfirm,
  onCancel,
  isLoading,
}: NewRoomModalProps) {
  const [name, setName] = useState("New Room");
  const [description, setDescription] = useState("");

  useEffect(() => {
    if (open && parentRoom) {
      setName("New Room");
      setDescription(parentRoom.description ?? "");
    }
  }, [open, parentRoom]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    onConfirm({ name: name.trim(), description: description.trim() });
  };

  return (
    <Modal
      isOpen={open}
      onClose={onCancel}
      title={direction ? `Add Room to the ${direction}` : "Add Room"}
    >
      <form onSubmit={handleSubmit} className="space-y-3">
        <FormField label="Name" value={name} onChange={setName} required autoFocus />
        <TextareaField
          label="Description"
          value={description}
          onChange={setDescription}
          rows={3}
        />
        <div className="flex gap-2 justify-end">
          <Button
            variant="secondary"
            onClick={onCancel}
            type="button"
            disabled={isLoading}
          >
            Cancel
          </Button>
          <Button
            variant="primary"
            type="submit"
            disabled={isLoading || !name.trim()}
          >
            {isLoading ? "Creating..." : "Create"}
          </Button>
        </div>
      </form>
    </Modal>
  );
}
