import type { CombatLogEntry } from "../lib/types";

type Props = {
  readonly entries: readonly CombatLogEntry[];
};

function kindColor(kind: CombatLogEntry["kind"]): string {
  switch (kind) {
    case "hit": return "text-foreground";
    case "crit": return "text-warning";
    case "miss": return "text-muted";
    case "heal": return "text-success";
    case "flee": return "text-info";
    case "defeat": return "text-danger";
    case "queue": return "text-accent";
    default: return "text-muted";
  }
}

export default function CombatLog({ entries }: Props) {
  return (
    <div className="flex-1 min-h-0 overflow-y-auto bg-black/20 rounded border border-border/50 px-3 py-2 space-y-1">
      {entries.length === 0 && (
        <p className="text-[10px] text-muted text-center italic">Combat log empty...</p>
      )}
      {entries.map((entry, idx) => (
        <div key={idx} className={`text-[11px] font-mono leading-snug ${kindColor(entry.kind)}`}>
          <span className="opacity-50 mr-1">
            {new Date(entry.timestamp).toLocaleTimeString("en-US", {
              hour12: false,
              hour: "2-digit",
              minute: "2-digit",
              second: "2-digit",
            })}
          </span>
          {entry.text}
        </div>
      ))}
    </div>
  );
}
