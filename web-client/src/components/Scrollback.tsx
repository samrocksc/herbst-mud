import type { WSLine } from "../hooks/useMUDSocket";

function styleClass(kind: WSLine["kind"]): string {
  switch (kind) {
    case "system":    return "text-accent font-bold";
    case "error":     return "text-danger font-bold";
    case "input":     return "text-accent";
    case "ping":      return "text-warning";
    case "notification": return "text-warning font-bold";  // Highlight quest/completion notifications
    default:          return "";
  }
}

type Props = {
  lines: WSLine[]
};

// Combat messages are filtered out - they appear in the CombatScreen instead
function isCombatMessage(text: string): boolean {
  // Combat messages start with these emoji prefixes or contain combat keywords
  const combatPrefixes = ["⚔", "🎲", "✨", "☠", "🛡️"];
  const combatKeywords = ["hit", "miss", "crit", "fumble", "damage", "Combat", "combat"];

  for (const prefix of combatPrefixes) {
    if (text.startsWith(prefix)) {
      return true;
    }
  }

  const lower = text.toLowerCase();
  for (const keyword of combatKeywords) {
    if (lower.includes(keyword)) {
      return true;
    }
  }

  return false;
}

export default function Scrollback({ lines }: Props) {
  return (
    <div className="flex flex-col min-h-full">
      {lines.filter((line) => !isCombatMessage(line.text)).map((line) => (
        line.text === "" ? (
          <div key={line.id} className="h-1" />
        ) : (
          <div key={line.id} className={`whitespace-pre-wrap break-words text-xs leading-relaxed ${styleClass(line.kind)}`}>
            {line.text}
          </div>
        )
      ))}
    </div>
  );
}