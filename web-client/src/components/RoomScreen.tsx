import { type RoomScreenPayload } from "../lib/types";
import { ActionExpand, type EntityAction } from "../ui";

type Props = {
  room: RoomScreenPayload
  onTapExit: (exit: { direction: string; label: string }) => void
  onCommand: (cmd: string) => void
  expandedId: string | null
  onToggleExpand: (id: string) => void
  pendingTargets: Set<number>
  onTogglePending: (id: number) => void
  onConfirmAttack: (char: RoomScreenPayload["characters"][number]) => void
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
  const characterActions = (char: RoomScreenPayload["characters"][number]): EntityAction[] => {
  const actions: EntityAction[] = [];

  if (char.hostile) {
    actions.push({
      label: "Attack",
      variant: "danger",
      onClick: () => onConfirmAttack(char),
    });
  } else {
    if (pendingTargets.has(char.id)) {
      actions.push({
        label: "Confirm",
        variant: "danger",
        onClick: () => onConfirmAttack(char),
      });
    } else {
      actions.push({
        label: "Attack",
        variant: "danger",
        onClick: () => onTogglePending(char.id),
      });
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
      <div className="px-3 py-2">
        <h2 className="font-bold text-sm text-accent mb-1">{room.title}</h2>
        <p className="text-[11px] text-muted leading-relaxed">{room.description}</p>
      </div>

      {/* Exits stay simple — tap = go */}
      <div className="px-3 py-2 border-t border-border">
        <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Exits</div>
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
      </div>

      {/* Characters — expandable */}
      {room.characters.length > 0 && (
        <div className="px-3 py-2 border-t border-border">
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Characters</div>
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
        </div>
      )}

      {/* Items — expandable */}
      {room.items.length > 0 && (
        <div className="px-3 py-2 border-t border-border">
          <div className="text-[10px] text-muted mb-1 uppercase tracking-wider">Items</div>
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
        </div>
      )}
    </div>
  );
}