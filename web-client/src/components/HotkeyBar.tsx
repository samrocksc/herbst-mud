import { ActionChip } from "../ui";
import type { CharacterSkill } from "../lib/types";

type Props = {
  readonly onActivate: (slot: string) => void;
  readonly skills?: readonly CharacterSkill[];
};

const STATIC_SLOTS = [
  { key: "L", label: "look", color: "text-muted" },
  { key: "E", label: "examine", color: "text-muted" },
  { key: "R", label: "use potion", color: "text-success" },
];

function buildSlots(skills?: readonly CharacterSkill[]) {
  const skillMap = new Map(
    skills?.map((sk) => [sk.slot, sk] as const) ?? [],
  );

  const skillSlots = [1, 2, 3, 4].map((slot) => {
    const sk = skillMap.get(slot);
    return {
      key: String(slot),
      label: sk?.name ?? "\u2014",
      color: sk?.name ? "text-accent" : "text-muted",
    };
  });

  return [...skillSlots, ...STATIC_SLOTS];
}

export default function HotkeyBar({ onActivate, skills }: Readonly<Props>) {
  const slots = buildSlots(skills);

  return (
    <div className="shrink-0 bg-surface border-t border-border px-2 py-2">
      <div className="flex gap-2 flex-wrap">
        {slots.map(({ key, label, color }) => (
          <ActionChip key={key} onClick={() => onActivate(key.toLowerCase())}>
            <span className="font-bold text-accent">{key}</span>
            <span className={color}>{label}</span>
          </ActionChip>
        ))}
      </div>
    </div>
  );
}