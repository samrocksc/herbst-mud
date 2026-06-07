import { useState, useCallback, useEffect, useRef } from "react";
import { ActionChip, Button } from "./Button";

export type EntityAction = {
  readonly label: string;
  readonly variant?: "primary" | "secondary" | "ghost" | "danger" | "success" | "warning" | "info";
  readonly onClick: () => void;
  readonly keepOpen?: boolean;
};

type ActionExpandProps = {
  readonly label: string;
  readonly triggerClassName?: string;
  readonly actions: readonly EntityAction[];
  readonly isOpen?: boolean;
  readonly onToggle?: () => void;
  readonly defaultOpen?: boolean;
};

/** Inline-expandable entity action chip. */
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
  const setOpen = useCallback((v: boolean) => {
    if (!onToggle) setInternalOpen(v);
  }, [onToggle]);

  const handleToggle = useCallback(() => {
    if (onToggle) onToggle();
    else setInternalOpen((v) => !v);
  }, [onToggle]);

  const containerRef = useRef<HTMLDivElement>(null);
  const menuRef = useRef<HTMLDivElement>(null);

  // Close on outside click
  useEffect(() => {
    if (!open) return;

    const handleClick = (e: MouseEvent) => {
      const container = containerRef.current;
      const menu = menuRef.current;
      if (!container || !menu) return;

      const containerRect = container.getBoundingClientRect();
      const menuRect = menu.getBoundingClientRect();
      const x = e.clientX;
      const y = e.clientY;

      const isTriggerClick =
        x >= containerRect.left &&
        x <= containerRect.right &&
        y >= containerRect.top &&
        y <= containerRect.bottom;

      const isMenuClick =
        x >= menuRect.left &&
        x <= menuRect.right &&
        y >= menuRect.top &&
        y <= menuRect.bottom;

      if (!isTriggerClick && !isMenuClick) {
        if (onToggle) onToggle();
        else setOpen(false);
      }
    };

    document.addEventListener("click", handleClick);
    return () => {
      document.removeEventListener("click", handleClick);
    };
  }, [open, setOpen]);

  return (
    <div ref={containerRef} className="relative inline-flex flex-col min-w-0">
      <ActionChip
        onClick={handleToggle}
        className={`${triggerClassName} ${open ? "border-accent bg-surface-alt" : ""} flex-1 min-w-0`}
      >
        <span className={open ? "text-accent font-bold" : ""} title={label}>
          {label}
        </span>
      </ActionChip>
      {open && (
        <div
          ref={menuRef}
          className="fixed z-50 bg-surface border border-border rounded shadow-lg p-1 flex flex-col gap-1 min-w-[100px] animate-in fade-in slide-in-from-top-2 duration-200"
        >
          {actions.map((a) => (
            <Button
              key={a.label}
              variant={a.variant || "secondary"}
              size="sm"
              fullWidth
              onClick={(e) => {
                e.stopPropagation();
                a.onClick();
                if (!a.keepOpen) {
                  if (onToggle) onToggle();
                  else setOpen(false);
                }
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
