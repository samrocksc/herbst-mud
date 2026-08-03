import { useCallback, useEffect, useState, type ReactNode } from "react";
import { Fragment } from "react";
import type { DialogNode } from "../lib/types";

export type ConversationOverlayProps = {
  npcName: string;
  nodes: Map<string, DialogNode>;
  currentNodeId: string;
  onChoice: (index: number) => void;
  onClose: () => void;
};

export default function ConversationOverlay({
  npcName,
  nodes,
  currentNodeId,
  onChoice,
  onClose,
}: Readonly<ConversationOverlayProps>) {
  const [loading, setLoading] = useState(false);
  const node = nodes.get(currentNodeId);

  // Reset loading when conversation advances to a new node
  useEffect(() => {
    setLoading(false);
  }, [currentNodeId]);

  const handleClickChoice = useCallback(
    (index: number) => {
      setLoading(true);
      onChoice(index);
    },
    [onChoice],
  );

  // Keyboard support: 1-9 select responses, 0/Escape/q closes
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;

      if (e.key === "Escape" || e.key === "0" || e.key.toLowerCase() === "q") {
        e.preventDefault();
        onClose();
        return;
      }

      const num = Number(e.key);
      if (num >= 1 && num <= 9 && node) {
        const idx = num - 1;
        if (idx < node.responses.length) {
          e.preventDefault();
          handleClickChoice(idx);
        }
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [node, onClose, handleClickChoice]);

  if (!node) {
    return (
      <div className="fixed inset-0 z-50 bg-background flex items-center justify-center">
        <p className="text-xs text-muted font-mono italic">
          Loading conversation...
        </p>
      </div>
    );
  }

  // Render response buttons. Dead-end responses (no next_node_id) are rendered as
  // a single Leave option so we never duplicate the Leave control.
  const renderResponseButtons = () => {
    const buttons: ReactNode[] = [];
    let hasLeave = false;
    node.responses.forEach((resp, idx) => {
      const isLeave = !resp.next_node_id;
      if (isLeave) {
        hasLeave = true;
        return;
      }
      buttons.push(
        <button
          key={`${currentNodeId}-${idx}`}
          type="button"
          disabled={loading}
          onClick={() => handleClickChoice(idx)}
          className="w-full flex items-center gap-2 px-3 py-2 rounded border border-border transition-colors disabled:opacity-50 text-sm text-left bg-surface-alt hover:bg-accent hover:text-background"
        >
          <span className="shrink-0 inline-flex items-center justify-center w-5 h-5 rounded text-[10px] font-bold bg-accent/20 text-accent">
            {idx + 1}
          </span>
          <span>{resp.label || "Continue"}</span>
        </button>
      );
    });
    if (hasLeave || buttons.length === 0) {
      buttons.push(
        <button
          key="leave-conversation"
          type="button"
          onClick={onClose}
          className="w-full px-3 py-2 rounded border border-border bg-muted/30 hover:bg-muted/50 transition-colors text-sm text-muted"
        >
          Leave
        </button>
      );
    }
    return <Fragment>{buttons}</Fragment>;
  };

  return (
    <div className="fixed inset-0 z-50 bg-background flex flex-col">
      {/* Header */}
      <header className="shrink-0 flex items-center justify-between px-4 py-3 border-b border-border bg-surface">
        <div className="flex items-center gap-2 text-xs font-mono">
          <span className="font-bold text-accent">{npcName}</span>
        </div>
        <button
          type="button"
          onClick={onClose}
          className="text-muted hover:text-foreground text-xl leading-none px-1"
          aria-label="Close conversation"
        >
          &times;
        </button>
      </header>

      {/* Body */}
      <div className="flex-1 min-h-0 flex flex-col">
        {/* NPC text speech bubble — centered vertically */}
        <div className="flex-1 flex items-center justify-center px-4 py-6">
          <div className="w-full max-w-lg">
            <div className="relative bg-surface border border-border rounded p-5 font-mono">
              {/* Speech bubble arrow */}
              <div
                className="absolute left-6 -bottom-2 w-3 h-3 bg-surface border-r border-b border-border rotate-45"
                aria-hidden="true"
              />
              <p className="text-sm leading-relaxed text-foreground italic">
                &ldquo;{node.npc_text}&rdquo;
              </p>
            </div>
          </div>
        </div>

        {/* Response buttons — anchored to bottom */}
        <div className="shrink-0 border-t border-border bg-surface px-4 py-4">
          <div className="w-full max-w-lg mx-auto space-y-2">
            {node.responses.length > 0 ? (
              renderResponseButtons()
            ) : (
              <button
                type="button"
                onClick={onClose}
                className="w-full px-3 py-2 rounded border border-border bg-surface-alt hover:bg-accent hover:text-background transition-colors text-sm text-muted"
              >
                End conversation
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
