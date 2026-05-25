type Props = {
  readonly hp: number;
  readonly maxHp: number;
  readonly stamina: number;
  readonly maxStamina: number;
  readonly mana: number;
  readonly maxMana: number;
};

function Bar({ label, current, max, colorClass }: { label: string; current: number; max: number; colorClass: string }) {
  const pct = max > 0 ? Math.max(0, Math.min(100, (current / max) * 100)) : 0;
  return (
    <div className="space-y-1">
      <div className="flex items-center justify-between text-[10px]">
        <span className="font-bold">{label}</span>
        <span className="text-muted">
          {current}/{max}
        </span>
      </div>
      <div className="w-full h-2 bg-surface-muted rounded overflow-hidden">
        <div
          className={`h-full transition-all duration-300 ${colorClass}`}
          style={{ width: `${pct}%` }}
        />
      </div>
    </div>
  );
}

export default function CombatVitals({ hp, maxHp, stamina, maxStamina, mana, maxMana }: Props) {
  return (
    <div className="space-y-2">
      <div className="text-[10px] text-muted uppercase tracking-wider">Vitals</div>
      <Bar label="❤ HP" current={hp} max={maxHp} colorClass="hp-bar-fill" />
      <Bar label="⚡ STA" current={stamina} max={maxStamina} colorClass="sta-bar-fill" />
      <Bar label="💧 MANA" current={mana} max={maxMana} colorClass="mana-bar-fill" />
    </div>
  );
}
