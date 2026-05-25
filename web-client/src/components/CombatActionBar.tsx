import type { CharacterSkill } from "../lib/types";

type Props = {
  readonly skills: readonly CharacterSkill[];
  readonly potionCount: number;
  readonly onSkill: (slot: number) => void;
  readonly onPotion: () => void;
  readonly onFlee: () => void;
  readonly queuedAction: string | null;
};

function ActionButton({
  slot,
  label,
  onClick,
  isQueued,
  variant = "primary",
}: {
  slot: string;
  label: string;
  onClick: () => void;
  isQueued: boolean;
  variant?: "primary" | "success" | "danger";
}) {
  const variantClass =
    variant === "danger"
      ? "border-danger bg-danger/10 hover:bg-danger/20 text-danger"
      : variant === "success"
        ? "border-success bg-success/10 hover:bg-success/20 text-success"
        : "border-accent bg-accent/10 hover:bg-accent/20 text-accent";

  return (
    <button
      onClick={onClick}
      className={`relative flex flex-col items-center justify-center w-16 h-16 rounded-lg border transition-all active:scale-95 ${variantClass} ${
        isQueued ? "ring-2 ring-accent ring-offset-1 ring-offset-background" : ""
      }`}
    >
      <span className="text-lg font-bold">{slot}</span>
      <span className="text-[10px] truncate w-full px-1 text-center">{label}</span>
    </button>
  );
}

export default function CombatActionBar({
  skills,
  potionCount,
  onSkill,
  onPotion,
  onFlee,
  queuedAction,
}: Props) {
  const skillSlots = [1, 2, 3, 4];

  return (
    <div className="shrink-0 bg-surface border-t border-border px-3 py-3">
      <div className="flex gap-3 justify-center">
        {skillSlots.map((slot) => {
          const sk = skills.find((s) => s.slot === slot);
          return (
            <ActionButton
              key={slot}
              slot={String(slot)}
              label={sk?.name ?? "—"}
              onClick={() => onSkill(slot)}
              isQueued={queuedAction === (sk?.name ?? "attack")}
            />
          );
        })}
        <ActionButton
          slot="5"
          label={`Potion${potionCount > 0 ? ` (${potionCount})` : ""}`}
          onClick={onPotion}
          isQueued={queuedAction === "use potion"}
          variant="success"
        />
        <ActionButton
          slot="F"
          label="Flee"
          onClick={onFlee}
          isQueued={queuedAction === "flee"}
          variant="danger"
        />
      </div>
    </div>
  );
}
