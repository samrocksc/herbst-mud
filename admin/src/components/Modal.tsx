/* eslint-disable functional/no-mixed-types, functional/immutable-data */
import { useEffect } from "react";
import type { ReactNode } from "react";
import { Button } from "./Button";

type ModalProps = Readonly<{
  isOpen: boolean
  onClose: () => void
  title: string
  children: ReactNode
}>

export function Modal({ isOpen, onClose, title, children }: ModalProps) {
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscape);
      document.body.style.overflow = "hidden";
      return () => {
        document.removeEventListener("keydown", handleEscape);
        document.body.style.overflow = "";
      };
    }
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="relative bg-surface-muted border-2 border-border-dark rounded-lg max-w-lg w-full mx-4 max-h-[90vh] overflow-y-auto shadow-xl">
        {/* Header */}
        <div className="sticky top-0 flex items-center justify-between p-4 border-b border-border bg-surface-muted">
          <h2 className="m-0 text-text text-lg">{title}</h2>
          <Button variant="ghost" size="sm" onClick={onClose} aria-label="Close">
            ×
          </Button>
        </div>

        {/* Content */}
        <div className="p-4">
          {children}
        </div>
      </div>
    </div>
  );
}