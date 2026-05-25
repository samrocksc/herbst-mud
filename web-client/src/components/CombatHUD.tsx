import { useCallback, useEffect, useRef, useState } from "react";
import type { CharacterSkill } from "../lib/types";

type Cooldown = {
  slot: number;
  expiresAt: number;
  durationMs: number;
};

type Props = {
  readonly skills: readonly CharacterSkill[];
  readonly potionCount: number;
  readonly onSkill: (slot: number) => void;
  readonly onPotion: () => void;
};

function CooldownOverlay({ expiresAt, durationMs }: { expiresAt: number; durationMs: number }) {
  const [pct, setPct] = useState(100);

  useEffect(() => {
    let raf: number;
    const tick = () => {
      const remaining = Math.max(0, expiresAt - Date.now());
      const p = durationMs > 0 ? (remaining / durationMs) * 100 : 0;
      setPct(p);
      if (remaining > 0) {
        raf = requestAnimationFrame(tick);
      }
    };
    raf = requestAnimationFrame(tick);
    return () => cancelAnimationFrame(raf);
  }, [expiresAt, durationMs]);

  if (pct <= 0) return null;

  return (
    <div className="absolute inset-0 flex items-center justify-center bg-black/60 rounded">
      <span className="text-xs font-bold text-white">{Math.ceil(pct / 100 * (durationMs / 1000))}s</span>
      <svg className="absolute inset-0 w-full h-full -rotate-90" viewBox="0 0 100 100">
        <circle cx="50" cy="50" r="46" fill="none" stroke="rgba(255,255,255,0.2)" strokeWidth="8" />
        <circle
          cx="50" cy="50" r="46" fill="none" stroke="var(--mud-accent)"
          strokeWidth="8"
          strokeDasharray={`${289 * (1 - pct / 100)} 289`}
          strokeLinecap="round"
        />
      </svg>
    </div>
  );
}

const DEFAULT_COOLDOWNS: Record<string, number> = {
  attack: 1500,
  "use potion": 3000,
};

export default function CombatHUD({ skills, potionCount, onSkill, onPotion }: Readonly<Props>) {
  const [cooldowns, setCooldowns] = useState<Map<number, Cooldown>>(new Map());
  const cooldownsRef = useRef(cooldowns);
  cooldownsRef.current = cooldowns;

  const startCooldown = useCallback((slot: number, abilityName: string | null) => {
    const duration = (abilityName && DEFAULT_COOLDOWNS[abilityName]) || 1500;
    setCooldowns((prev) => {
      const next = new Map(prev);
      next.set(slot, { slot, expiresAt: Date.now() + duration, durationMs: duration });
      return next;
    });
  }, []);

  const handleSkill = useCallback((slot: number) => {
    const cd = cooldownsRef.current.get(slot);
    if (cd && cd.expiresAt > Date.now()) return;
    const sk = skills.find((s) => s.slot === slot);
    startCooldown(slot, sk?.name ?? null);
    onSkill(slot);
  }, [skills, onSkill, startCooldown]);

  const handlePotion = useCallback(() => {
    const cd = cooldownsRef.current.get(5);
    if (cd && cd.expiresAt > Date.now()) return;
    startCooldown(5, "use potion");
    onPotion();
  }, [onPotion, startCooldown]);

  // Keyboard 1-5
  useEffect(() => {
    const handler = (e: KeyboardEvent) => {
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) return;
      const key = e.key;
      if (key >= "1" && key <= "4") {
        e.preventDefault();
        handleSkill(Number(key));
      } else if (key === "5") {
        e.preventDefault();
        handlePotion();
      }
    };
    window.addEventListener("keydown", handler);
    return () => window.removeEventListener("keydown", handler);
  }, [handleSkill, handlePotion]);

  // Clean expired cooldowns every second
  useEffect(() => {
    const id = setInterval(() => {
      const now = Date.now();
      setCooldowns((prev) => {
        const next = new Map(prev);
        for (const [slot, cd] of next) {
          if (cd.expiresAt <= now) next.delete(slot);
        }
        return next;
      });
    }, 1000);
    return () => clearInterval(id);
  }, []);

  const skillSlots = [1, 2, 3, 4];

  return (
    <div className="shrink-0 bg-surface border-t border-border px-3 py-3">
      <div className="flex gap-3 justify-center">
        {skillSlots.map((slot) => {
          const sk = skills.find((s) => s.slot === slot);
          const cd = cooldowns.get(slot);
          const onCd = !!cd && cd.expiresAt > Date.now();
          return (
            <button
              key={slot}
              onClick={() => handleSkill(slot)}
              disabled={onCd}
              className="relative flex flex-col items-center justify-center w-16 h-16 rounded-lg border border-border bg-surface-muted hover:bg-surface transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span className="text-lg font-bold text-accent">{slot}</span>
              <span className="text-[10px] text-text-muted truncate w-full px-1 text-center">
                {sk?.name ?? "—"}
              </span>
              {cd && <CooldownOverlay expiresAt={cd.expiresAt} durationMs={cd.durationMs} />}
            </button>
          );
        })}

        {/* Potion slot */}
        <button
          onClick={handlePotion}
          disabled={(cooldowns.get(5)?.expiresAt ?? 0) > Date.now()}
          className="relative flex flex-col items-center justify-center w-16 h-16 rounded-lg border border-border bg-success/10 hover:bg-success/20 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span className="text-lg font-bold text-success">5</span>
          <span className="text-[10px] text-success truncate w-full px-1 text-center">
            Potion{potionCount > 0 ? ` (${potionCount})` : ""}
          </span>
          {cooldowns.get(5) && (
            <CooldownOverlay expiresAt={cooldowns.get(5)!.expiresAt} durationMs={cooldowns.get(5)!.durationMs} />
          )}
        </button>
      </div>
    </div>
  );
}
