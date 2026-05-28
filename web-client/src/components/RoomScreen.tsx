import { useState } from "react";
import { type RoomScreenPayload } from "../lib/types";
import { ActionExpand, type EntityAction } from "../ui";

type Props = {
  room: RoomScreenPayload;
  onTapExit: (exit: { direction: string; label: string }) => void;
  onCommand: (cmd: string) => void;
  expandedId: string | null;
  onToggleExpand: (id: string) => void;
  pendingTargets: Set<number>;
  onTogglePending: (id: number) => void;
  onConfirmAttack: (char: RoomScreenPayload["characters"][number]) => void;
};

export default function RoomScreen({
  room,
  onTapExit,
  onCommand,
  expandedId,
  onToggleExpand,
  pendingTargets,
  onTogglePending,
  onConfirmAttack,
}: Props) {
  const [descHidden, setDescHidden] = useState(false);

  const characterActions = (char: RoomScreenPayload["characters"][number]): EntityAction[] => {
    const actions: EntityAction[] = [];
    if (char.hostile) {
      if (pendingTargets.has(char.id)) {
        actions.push({ label: "Confirm", variant: "danger", onClick: () => onConfirmAttack(char) });
      } else {
        actions.push({ label: "Attack", variant: "danger", onClick: () => onTogglePending(char.id) });
      }
    } else {
      if (pendingTargets.has(char.id)) {
        actions.push({ label: "Confirm", variant: "danger", onClick: () => onConfirmAttack(char) });
      } else {
        actions.push({ label: "Attack", variant: "danger", onClick: () => onTogglePending(char.id) });
      }
      actions.push({ label: "Talk", variant: "success", onClick: () => onCommand(`talk ${char.name}`) });
    }
    actions.push({ label: "Examine", variant: "secondary", onClick: () => onCommand(`examine ${char.name}`) });
    return actions;
  };

  const itemActions = (item: RoomScreenPayload["items"][number]): EntityAction[] => {
    const actions: EntityAction[] = [
      { label: "Examine", variant: "secondary", onClick: () => onCommand(`examine ${item.name}`) },
    ];
    if (item.takeable) {
      actions.push({ label: "Take", variant: "primary", onClick: () => onCommand(`take ${item.name}`) });
    }
    return actions;
  };

  return (
    <div className="shrink-0 bg-surface border-t border-border">
      {/* Room title + description — full width */}
      <div className="px-3 py-2 border-b border-border">
        <div className="flex items-center justify-between">
          <h2 className="font-bold text-sm text-accent">{room.title}</h2>
          <button
            type="button"
            onClick={() => setDescHidden((v) => !v)}
            className="text-[10px] text-muted hover:text-foreground px-1"
            title={descHidden ? "Show description" : "Hide description"}
          >
            {descHidden ? "[+]" : "[-]"}
          </button>
        </div>
        {!descHidden && <p className="text-[11px] text-muted leading-relaxed">{room.description}</p>}
      </div>

      {/* 2x2 Quadrant Grid */}
      <div className="grid grid-cols-2 grid-rows-2 gap-px bg-border">
        {/* Top Left: Characters */}
        <div className="bg-surface p-2 overflow-y-auto" style={{ maxHeight: "154px" }}>
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Characters</div>
          {room.characters.length > 0 ? (
            <div className="flex gap-2 flex-wrap">
              {room.characters.map((char) => (
                <ActionExpand
                  key={char.id}
                  isOpen={expandedId === `char-${char.id}`}
                  onToggle={() => onToggleExpand(`char-${char.id}`)}
                  label={char.name}
                  triggerClassName={char.hostile ? "text-danger" : "text-foreground"}
                  actions={characterActions(char)}
                />
              ))}
            </div>
          ) : (
            <span className="text-[11px] text-muted italic">None</span>
          )}
        </div>

        {/* Top Right: Exits */}
        <div className="bg-surface p-2 overflow-y-auto" style={{ maxHeight: "154px" }}>
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Exits</div>
          {room.exits.length > 0 ? (
            <div className="flex gap-2 flex-wrap">
              {room.exits.map((exit) => (
                <button
                  key={exit.direction}
                  onClick={() => onTapExit(exit)}
                  className="inline-flex items-center gap-1 px-2 py-1 rounded border border-border text-[11px] font-mono hover:bg-surface-alt active:bg-accent active:text-background transition-colors"
                >
                  <span className="font-bold text-accent">{exit.direction.charAt(0).toUpperCase()}</span>
                  <span className="text-muted">{exit.label}</span>
                </button>
              ))}
            </div>
          ) : (
            <span className="text-[11px] text-muted italic">None</span>
          )}
        </div>

        {/* Bottom Left: Items */}
        <div className="bg-surface p-2 overflow-y-auto" style={{ maxHeight: "200px" }}>
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Items</div>
          {room.items.length > 0 ? (
            <div className="flex gap-2 flex-wrap">
              {room.items.map((item) => (
                <ActionExpand
                  key={item.id}
                  isOpen={expandedId === `item-${item.id}`}
                  onToggle={() => onToggleExpand(`item-${item.id}`)}
                  label={item.name}
                  triggerClassName={item.takeable ? "text-success" : "text-muted"}
                  actions={itemActions(item)}
                />
              ))}
            </div>
          ) : (
            <span className="text-[11px] text-muted italic">None</span>
          )}
        </div>

        {/* Bottom Right: Functions */}
        <div className="bg-surface p-2 overflow-y-auto" style={{ maxHeight: "200px" }}>
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Functions</div>
          <div className="flex gap-2 flex-wrap">
            <button
              onClick={() => onCommand("look")}
              className="inline-flex items-center gap-1 px-2 py-1 rounded border border-border text-[11px] font-mono hover:bg-surface-alt active:bg-accent active:text-background transition-colors"
            >
              <span className="font-bold text-accent">L</span>
              <span className="text-muted">look</span>
            </button>
            <button
              onClick={() => onCommand("examine")}
              className="inline-flex items-center gap-1 px-2 py-1 rounded border border-border text-[11px] font-mono hover:bg-surface-alt active:bg-accent active:text-background transition-colors"
            >
              <span className="font-bold text-accent">E</span>
              <span className="text-muted">examine</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
