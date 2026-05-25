import { useState } from "react";
import { ActionChip, Button } from "./Button";

export type EntityAction = {
  readonly label: string;
  readonly variant?: "primary" | "secondary" | "ghost" | "danger" | "success" | "warning" | "info";
  readonly onClick: () => void;
};

type ActionExpandProps = {
  readonly label: string;
  readonly triggerClassName?: string;
  readonly actions: readonly EntityAction[];
  readonly isOpen?: boolean;
  readonly onToggle?: () => void;
  readonly defaultOpen?: boolean;
};

/** Inline-expandable entity action chip. Single-select when used inside RoomScreen. */
export function ActionExpand({
  label,
  triggerClassName = "",
  actions,
  isOpen,
  onToggle,
  defaultOpen = false,
}: Readonly<ActionExpandProps>) {
  const [internalOpen, setInternalOpen] = useState(defaultOpen);
  const open = isOpen !== undefined ? isOpen : internalOpen;
  const setOpen = (v: boolean) => {
    if (!onToggle) setInternalOpen(v);
  };

  const handleToggle = () => {
    if (onToggle) onToggle();
    else setInternalOpen((v) => !v);
  };

  return (
    <div className="flex flex-col gap-1">
      <ActionChip
        onClick={handleToggle}
        className={`${triggerClassName} ${open ? "border-accent bg-surface-alt" : ""}`}
      >
        <span className={open ? "text-accent font-bold" : ""}>{label}</span>
      </ActionChip>
      {open && (
        <div className="flex gap-1 flex-wrap animate-in fade-in slide-in-from-top-1 duration-150">
          {actions.map((a) => (
            <Button
              key={a.label}
              variant={a.variant || "secondary"}
              size="sm"
              onClick={(e) => {
                e.stopPropagation();
                a.onClick();
                setOpen(false);
              }}
            >
              {a.label}
            </Button>
          ))}
        </div>
      )}
    </div>
  );
}