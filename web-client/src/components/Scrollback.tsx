import type { WSLine } from "../hooks/useMUDSocket";

function styleClass(kind: WSLine["kind"]): string {
  switch (kind) {
    case "system": return "text-accent font-bold";
    case "error":  return "text-danger font-bold";
    case "input":  return "text-accent";
    case "ping":   return "text-warning";
    default:       return "";
  }
}

type Props = {
  lines: WSLine[]
};

export default function Scrollback({ lines }: Props) {
  return (
    <div className="flex flex-col min-h-full">
      {lines.map((line) => (
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