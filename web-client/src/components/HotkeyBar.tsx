import { ActionChip } from "../ui";
import type { CharacterSkill } from "../lib/types";

type Props = {
  readonly onActivate: (slot: string) => void;
  readonly skills?: readonly CharacterSkill[];
};

const STATIC_SLOTS = [
  { key: "5", label: "concentrate", color: "text-accent" },
  { key: "L", label: "look", color: "text-muted" },
  { key: "E", label: "examine", color: "text-muted" },
  { key: "Q", label: "quit", color: "text-danger" },
  { key: "R", label: "use potion", color: "text-success" },
];

const SKILL_SLOT_KEYS = ["1", "2", "3", "4"];

function buildSlots(skills?: readonly CharacterSkill[]) {
  const skillMap = new Map(
    skills?.map((sk) => [sk.slot, sk] as const) ?? [],
  );

  const skillSlots = SKILL_SLOT_KEYS.map((key, idx) => {
    const slot = idx + 1;
    const sk = skillMap.get(slot);
    return {
      key,
      label: sk?.name ?? "\u2014",
      color: sk?.name ? "text-accent" : "text-muted",
    };
  });

  const staticSlots = STATIC_SLOTS.map((s) => {
    const numKey = Number(s.key);
    const sk = Number.isNaN(numKey) ? undefined : skillMap.get(numKey);
    if (sk?.name) {
      return { key: s.key, label: sk.name, color: "text-accent" };
    }
    return s;
  });

  return [...skillSlots, ...staticSlots];
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