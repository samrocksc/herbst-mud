import type { CombatTarget } from "../lib/types";

function HPBar({ current, max }: { current: number; max: number }) {
  const pct = max > 0 ? Math.max(0, Math.min(100, (current / max) * 100)) : 0;
  return (
    <div className="w-full h-2 bg-surface-muted rounded overflow-hidden">
      <div
        className="h-full transition-all duration-300 hp-bar-fill"
        style={{ width: `${pct}%` }}
      />
    </div>
  );
}

type Props = {
  readonly targets: readonly CombatTarget[];
};

export default function CombatTargetList({ targets }: Props) {
  if (targets.length === 0) return null;

  return (
    <div className="space-y-2">
      <div className="text-[10px] text-muted uppercase tracking-wider">Targets</div>
      {targets.map((t) => (
        <div key={t.id} className="space-y-1">
          <div className="flex items-center justify-between text-xs">
            <span className="font-bold text-danger">🎯 {t.name}</span>
            <span className="text-muted text-[10px]">
              {t.hp}/{t.maxHp} HP
            </span>
          </div>
          <HPBar current={t.hp} max={t.maxHp} />
        </div>
      ))}
    </div>
  );
}
