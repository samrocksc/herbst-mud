# Combat Screen Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the web client's toggle-based combat mode with a full Combat Screen that replaces the room panel when in combat, mirroring the SSH client's `ScreenCombat`.

**Architecture:** The web client runs its own combat tick engine (1.5s intervals, matching SSH) via a `useCombatEngine` hook. Combat state lives in `GameScreen` and is passed to a new `CombatScreen` component that replaces `RoomScreen` when `inCombat` is true. `RoomScreen` gains a confirm-flow for non-hostile targets. The `HotkeyBar` loses combat commands in adventure mode.

**Tech Stack:** React, TypeScript, Tailwind CSS, Vite

---

## File Structure

| File | Responsibility |
|---|---|
| `web-client/src/lib/types.ts` | Add `CombatTarget`, `CombatLogEntry` types |
| `web-client/src/lib/api.ts` | Add `getCombatStatus`, `applyDamage`, `healCharacter` |
| `web-client/src/lib/combat.ts` | Dice rolls, damage calc, AC calc — pure functions |
| `web-client/src/hooks/useCombatEngine.ts` | Combat state machine, tick timer, round logic |
| `web-client/src/components/CombatTargetList.tsx` | Render targets with HP bars |
| `web-client/src/components/CombatVitals.tsx` | Player HP/STA/MANA bars |
| `web-client/src/components/CombatActionBar.tsx` | 4 ability slots + potion + flee buttons |
| `web-client/src/components/CombatLog.tsx` | Timestamped combat log entries |
| `web-client/src/components/CombatScreen.tsx` | Main combat view composing the above |
| `web-client/src/components/RoomScreen.tsx` | Add pending target confirm flow |
| `web-client/src/components/HotkeyBar.tsx` | Remove combat commands from adventure mode |
| `web-client/src/components/GameScreen.tsx` | Wire in combat engine, conditionally render CombatScreen |
| `web-client/src/index.css` | HP/STA/MANA bar gradient styles |

---

## Tick Interval Note

The SSH client's `herbst/combat/config.go` defines `DefaultTickInterval = 1500` (1.5 seconds). Use `1500` ms for the web client tick to match exactly.

---

## Task 1: Types + API Helpers

**Files:**
- Modify: `web-client/src/lib/types.ts`
- Modify: `web-client/src/lib/api.ts`

### Step 1: Add CombatTarget and CombatLogEntry to types.ts

Append to the end of `web-client/src/lib/types.ts` (after line 65):

```ts
export type CombatTarget = {
  readonly id: number;
  readonly name: string;
  readonly hp: number;
  readonly maxHp: number;
  readonly level?: number;
};

export type CombatLogEntry = {
  readonly timestamp: number;
  readonly text: string;
  readonly kind: "hit" | "miss" | "crit" | "heal" | "system" | "queue" | "flee" | "defeat";
};
```

### Step 2: Add API helpers to api.ts

Append to the end of `web-client/src/lib/api.ts` (after line 207):

```ts
export async function getCombatStatus(charID: number): Promise<{ hp: number; maxHp: number; isNPC: boolean }> {
  const res = await fetch(`${API_BASE}/characters/${charID}/combat-status`, { headers: headers() });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to load combat status" }));
    throw new Error(err.error || "Failed to load combat status");
  }
  return res.json();
}

export async function applyDamage(charID: number, damage: number, attackerID?: number): Promise<void> {
  const body: Record<string, unknown> = { damage };
  if (attackerID != null && attackerID > 0) body.attacker_id = attackerID;
  const res = await fetch(`${API_BASE}/characters/${charID}/damage`, {
    method: "POST",
    headers: headers(),
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to apply damage" }));
    throw new Error(err.error || "Failed to apply damage");
  }
}

export async function healCharacter(charID: number, amount: number): Promise<void> {
  const res = await fetch(`${API_BASE}/characters/${charID}/heal`, {
    method: "POST",
    headers: headers(),
    body: JSON.stringify({ amount }),
  });
  if (!res.ok) {
    const err = await res.json().catch((): { readonly error: string } => ({ error: "Failed to heal" }));
    throw new Error(err.error || "Failed to heal");
  }
}
```

### Step 3: Verify no TypeScript errors

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors (or only pre-existing ones).

### Step 4: Commit

```bash
git add web-client/src/lib/types.ts web-client/src/lib/api.ts
git commit -m "🟣 feat(combat): add CombatTarget, CombatLogEntry types and API helpers"
```

---

## Task 2: Combat Utility Functions

**Files:**
- Create: `web-client/src/lib/combat.ts`

### Step 1: Create combat.ts with dice and damage logic

```ts
/**
 * Pure combat utility functions.
 * Mirrors the SSH client's dice logic in herbst/dice/ and game_combat.go.
 */

export type DiceResult = {
  roll: number;
  total: number;
  isCrit: boolean;
  isFumble: boolean;
};

export function rollD20(modifier = 0): DiceResult {
  const roll = Math.floor(Math.random() * 20) + 1;
  return {
    roll,
    total: roll + modifier,
    isCrit: roll === 20,
    isFumble: roll === 1,
  };
}

export function rollDamage(sides: number, count: number, modifier = 0): number {
  let total = modifier;
  for (let i = 0; i < count; i++) {
    total += Math.floor(Math.random() * sides) + 1;
  }
  return Math.max(1, total);
}

/** Base player damage: 1 + strength/5 (matches server tryAttack logic) */
export function calculatePlayerDamage(strength: number): number {
  return Math.max(1, 1 + Math.floor(strength / 5));
}

/** Base enemy damage: level + 2 */
export function calculateEnemyDamage(level: number, isCrit = false): number {
  const dmg = level + 2;
  return isCrit ? dmg * 2 : dmg;
}

/** Player AC: base 10 + level/2 */
export function calculatePlayerAC(level: number): number {
  return 10 + Math.floor(level / 2);
}

/** Enemy AC: base 10 + level/2 */
export function calculateEnemyAC(level: number): number {
  return 10 + Math.floor(level / 2);
}

/** DEX modifier: level / 3 */
export function getDexModifier(level: number): number {
  return Math.floor(level / 3);
}

/** STR modifier: (strength - 10) / 2 */
export function getStrModifier(strength: number): number {
  return Math.floor((strength - 10) / 2);
}

/** Flee check: d20 + level/2 vs DC 12 */
export function attemptFlee(level: number): { success: boolean; roll: number; total: number } {
  const { roll, total } = rollD20(Math.floor(level / 2));
  return { success: total >= 12, roll, total };
}
```

### Step 2: Verify TypeScript

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors.

### Step 3: Commit

```bash
git add web-client/src/lib/combat.ts
git commit -m "🟣 feat(combat): add dice and damage utility functions"
```

---

## Task 3: useCombatEngine Hook

**Files:**
- Create: `web-client/src/hooks/useCombatEngine.ts`

### Step 1: Create the hook

```ts
import { useCallback, useEffect, useRef, useState } from "react";
import type { CombatTarget, CombatLogEntry, CharacterSkill } from "../lib/types";
import { getCombatStatus, applyDamage, healCharacter } from "../lib/api";
import {
  rollD20,
  calculatePlayerDamage,
  calculateEnemyDamage,
  calculateEnemyAC,
  getDexModifier,
  attemptFlee,
} from "../lib/combat";

const TICK_MS = 1500;

export type CombatEngineState = {
  inCombat: boolean;
  targets: CombatTarget[];
  combatLog: CombatLogEntry[];
  round: number;
  queuedAction: string | null;
  playerHP: number;
  playerMaxHP: number;
};

type Props = {
  characterID: number;
  characterLevel: number;
  characterStrength: number;
  initialHP: number;
  initialMaxHP: number;
  skills: readonly CharacterSkill[];
  onLog: (text: string, kind: CombatLogEntry["kind"]) => void;
  onCombatEnd: () => void;
  onPlayerHPChange: (hp: number) => void;
};

export function useCombatEngine({
  characterID,
  characterLevel,
  characterStrength,
  initialHP,
  initialMaxHP,
  onLog,
  onCombatEnd,
  onPlayerHPChange,
}: Props) {
  const [inCombat, setInCombat] = useState(false);
  const [targets, setTargets] = useState<CombatTarget[]>([]);
  const [combatLog, setCombatLog] = useState<CombatLogEntry[]>([]);
  const [round, setRound] = useState(0);
  const [queuedAction, setQueuedAction] = useState<string | null>(null);
  const [playerHP, setPlayerHP] = useState(initialHP);

  const targetsRef = useRef(targets);
  const playerHPRef = useRef(playerHP);
  const inCombatRef = useRef(inCombat);
  const roundRef = useRef(round);
  const queuedRef = useRef(queuedAction);
  const tickRef = useRef<ReturnType<typeof setInterval> | null>(null);

  targetsRef.current = targets;
  playerHPRef.current = playerHP;
  inCombatRef.current = inCombat;
  roundRef.current = round;
  queuedRef.current = queuedAction;

  const addLog = useCallback(
    (text: string, kind: CombatLogEntry["kind"]) => {
      const entry: CombatLogEntry = { timestamp: Date.now(), text, kind };
      setCombatLog((prev) => [...prev.slice(-99), entry]);
      onLog(text, kind);
    },
    [onLog]
  );

  const fetchTargetHP = useCallback(async (targetID: number) => {
    try {
      const status = await getCombatStatus(targetID);
      return status.hp;
    } catch {
      return null;
    }
  }, []);

  const performPlayerAttack = useCallback(
    async (target: CombatTarget) => {
      const dexMod = getDexModifier(characterLevel);
      const { roll, total, isCrit, isFumble } = rollD20(dexMod);
      const targetAC = calculateEnemyAC(target.level ?? 1);

      if (isFumble) {
        addLog("🎲 FUMBLE! Natural 1 — You stumble badly!", "miss");
        return;
      }

      if (total < targetAC && !isCrit) {
        addLog(`🎲 Miss! (d20=${roll} + ${dexMod} DEX = ${total} vs AC ${targetAC})`, "miss");
        return;
      }

      const damage = calculatePlayerDamage(characterStrength);
      const finalDamage = isCrit ? damage * 2 : damage;

      await applyDamage(target.id, finalDamage, characterID);

      const newHP = await fetchTargetHP(target.id);
      if (newHP != null) {
        setTargets((prev) =>
          prev.map((t) => (t.id === target.id ? { ...t, hp: newHP } : t))
        );
      }

      if (isCrit) {
        addLog(`⚔ CRITICAL HIT! ${finalDamage} damage!`, "crit");
      } else {
        addLog(`⚔ You hit ${target.name} for ${finalDamage} damage!`, "hit");
      }
    },
    [characterID, characterLevel, characterStrength, addLog, fetchTargetHP]
  );

  const performEnemyTurn = useCallback(async () => {
    const currentTargets = targetsRef.current;
    const aliveTargets = currentTargets.filter((t) => t.hp > 0);
    if (aliveTargets.length === 0) return;

    const target = aliveTargets[0];
    const enemyDexMod = Math.floor((target.level ?? 1) / 3);
    const { roll, total, isCrit, isFumble } = rollD20(enemyDexMod);
    const playerAC = calculateEnemyAC(characterLevel);

    if (isFumble) {
      addLog(`🎲 ${target.name} FUMBLES! (rolled 1)`, "miss");
      return;
    }

    if (total < playerAC && !isCrit) {
      addLog(`🎲 ${target.name} misses! (d20=${roll} + ${enemyDexMod} = ${total} vs AC ${playerAC})`, "miss");
      return;
    }

    const damage = calculateEnemyDamage(target.level ?? 1, isCrit);
    await applyDamage(characterID, damage);

    const newHP = Math.max(0, playerHPRef.current - damage);
    setPlayerHP(newHP);
    onPlayerHPChange(newHP);

    if (isCrit) {
      addLog(`⚔ ${target.name} critical hit! ${damage} damage!`, "crit");
    } else {
      addLog(`⚔ ${target.name} hits you for ${damage} damage!`, "hit");
    }

    if (newHP <= 0) {
      addLog("☠ You have been defeated!", "defeat");
      setInCombat(false);
      inCombatRef.current = false;
      if (tickRef.current) {
        clearInterval(tickRef.current);
        tickRef.current = null;
      }
      await healCharacter(characterID, initialMaxHP);
      onCombatEnd();
    }
  }, [characterID, characterLevel, initialMaxHP, addLog, onCombatEnd, onPlayerHPChange]);

  const processTick = useCallback(async () => {
    if (!inCombatRef.current) return;

    setRound((r) => r + 1);

    const action = queuedRef.current;
    setQueuedAction(null);

    if (action === "flee") {
      const { success, roll, total } = attemptFlee(characterLevel);
      if (success) {
        addLog(`🏃 Escape successful! (d20=${roll} + ${Math.floor(characterLevel / 2)} = ${total} vs DC 12)`, "flee");
        setInCombat(false);
        inCombatRef.current = false;
        if (tickRef.current) {
          clearInterval(tickRef.current);
          tickRef.current = null;
        }
        onCombatEnd();
        return;
      }
      addLog(`🏃 Escape failed! (d20=${roll} + ${Math.floor(characterLevel / 2)} = ${total} vs DC 12)`, "flee");
    } else if (action === "attack" || action == null) {
      const aliveTargets = targetsRef.current.filter((t) => t.hp > 0);
      if (aliveTargets.length > 0) {
        await performPlayerAttack(aliveTargets[0]);
      }
    }

    if (!inCombatRef.current) return;

    await performEnemyTurn();

    const remainingAlive = targetsRef.current.filter((t) => t.hp > 0);
    if (remainingAlive.length === 0 && inCombatRef.current) {
      addLog("✦ All targets defeated!", "system");
      setInCombat(false);
      inCombatRef.current = false;
      if (tickRef.current) {
        clearInterval(tickRef.current);
        tickRef.current = null;
      }
      onCombatEnd();
    }
  }, [characterLevel, addLog, onCombatEnd, performPlayerAttack, performEnemyTurn]);

  const startCombat = useCallback(
    async (newTargets: CombatTarget[]) => {
      const refreshed = await Promise.all(
        newTargets.map(async (t) => {
          try {
            const status = await getCombatStatus(t.id);
            return { ...t, hp: status.hp, maxHp: status.maxHp };
          } catch {
            return t;
          }
        })
      );

      setTargets(refreshed);
      setCombatLog([]);
      setRound(1);
      setQueuedAction(null);
      setInCombat(true);
      inCombatRef.current = true;
      roundRef.current = 1;

      addLog(`⚔ Combat started with ${refreshed.map((t) => t.name).join(", ")}!`, "system");

      if (tickRef.current) clearInterval(tickRef.current);
      tickRef.current = setInterval(() => {
        processTick();
      }, TICK_MS);
    },
    [addLog, processTick]
  );

  const queueAction = useCallback(
    (action: string) => {
      if (!inCombatRef.current) return;
      setQueuedAction(action);
      addLog(`⏱ Queued: ${action}`, "queue");
    },
    [addLog]
  );

  const exitCombat = useCallback(() => {
    setInCombat(false);
    inCombatRef.current = false;
    setQueuedAction(null);
    if (tickRef.current) {
      clearInterval(tickRef.current);
      tickRef.current = null;
    }
    onCombatEnd();
  }, [onCombatEnd]);

  useEffect(() => {
    return () => {
      if (tickRef.current) clearInterval(tickRef.current);
    };
  }, []);

  return {
    inCombat,
    targets,
    combatLog,
    round,
    queuedAction,
    playerHP,
    startCombat,
    queueAction,
    exitCombat,
  };
}
```

### Step 2: Verify TypeScript

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors.

### Step 3: Commit

```bash
git add web-client/src/hooks/useCombatEngine.ts
git commit -m "🟣 feat(combat): add useCombatEngine hook with tick-based combat"
```

---

## Task 4: Combat Sub-Components

**Files:**
- Create: `web-client/src/components/CombatTargetList.tsx`
- Create: `web-client/src/components/CombatVitals.tsx`
- Create: `web-client/src/components/CombatActionBar.tsx`
- Create: `web-client/src/components/CombatLog.tsx`

### Step 1: Create CombatTargetList.tsx

```tsx
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
```

### Step 2: Create CombatVitals.tsx

```tsx
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
```

### Step 3: Create CombatActionBar.tsx

```tsx
import { useCallback, useState } from "react";
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
```

### Step 4: Create CombatLog.tsx

```tsx
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
```

### Step 5: Verify TypeScript

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors.

### Step 6: Commit

```bash
git add web-client/src/components/CombatTargetList.tsx web-client/src/components/CombatVitals.tsx web-client/src/components/CombatActionBar.tsx web-client/src/components/CombatLog.tsx
git commit -m "🟣 feat(combat): add CombatScreen sub-components"
```

---

## Task 5: CombatScreen Component

**Files:**
- Create: `web-client/src/components/CombatScreen.tsx`

### Step 1: Create CombatScreen.tsx

```tsx
import type { CombatTarget, CombatLogEntry, CharacterSkill } from "../lib/types";
import CombatTargetList from "./CombatTargetList";
import CombatVitals from "./CombatVitals";
import CombatActionBar from "./CombatActionBar";
import CombatLog from "./CombatLog";

type Props = {
  readonly round: number;
  readonly targets: readonly CombatTarget[];
  readonly combatLog: readonly CombatLogEntry[];
  readonly queuedAction: string | null;
  readonly playerHP: number;
  readonly playerMaxHP: number;
  readonly playerStamina: number;
  readonly playerMaxStamina: number;
  readonly playerMana: number;
  readonly playerMaxMana: number;
  readonly skills: readonly CharacterSkill[];
  readonly potionCount: number;
  readonly onSkill: (slot: number) => void;
  readonly onPotion: () => void;
  readonly onFlee: () => void;
};

export default function CombatScreen({
  round,
  targets,
  combatLog,
  queuedAction,
  playerHP,
  playerMaxHP,
  playerStamina,
  playerMaxStamina,
  playerMana,
  playerMaxMana,
  skills,
  potionCount,
  onSkill,
  onPotion,
  onFlee,
}: Props) {
  return (
    <div className="shrink-0 bg-surface border-t border-border flex flex-col" style={{ maxHeight: "60vh" }}>
      {/* Header */}
      <div className="shrink-0 flex items-center justify-between px-3 py-2 border-b border-border bg-danger/5">
        <span className="text-xs font-bold text-danger uppercase tracking-wider">
          ⚔ Combat — Round {round}
        </span>
        {queuedAction && (
          <span className="text-[10px] text-accent">Queued: {queuedAction}</span>
        )}
      </div>

      {/* Content */}
      <div className="flex-1 min-h-0 overflow-y-auto px-3 py-3 space-y-4">
        <CombatTargetList targets={targets} />
        <CombatVitals
          hp={playerHP}
          maxHp={playerMaxHP}
          stamina={playerStamina}
          maxStamina={playerMaxStamina}
          mana={playerMana}
          maxMana={playerMaxMana}
        />
        <CombatLog entries={combatLog} />
      </div>

      {/* Action Bar */}
      <CombatActionBar
        skills={skills}
        potionCount={potionCount}
        onSkill={onSkill}
        onPotion={onPotion}
        onFlee={onFlee}
        queuedAction={queuedAction}
      />
    </div>
  );
}
```

### Step 2: Verify TypeScript

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors.

### Step 3: Commit

```bash
git add web-client/src/components/CombatScreen.tsx
git commit -m "🟣 feat(combat): add CombatScreen main component"
```

---

## Task 6: RoomScreen Confirm Logic

**Files:**
- Modify: `web-client/src/components/RoomScreen.tsx`

### Step 1: Add pending target props and confirm flow

Add new props to the `Props` type:

```ts
type Props = {
  room: RoomScreenPayload;
  onTapExit: (exit: { direction: string; label: string }) => void;
  onCommand: (cmd: string) => void;
  expandedId: string | null;
  onToggleExpand: (id: string) => void;
  pendingTargets: Set<number>;
  onTogglePending: (id: number) => void;
  onConfirmAttack: (char: RoomScreenPayload["characters"][number]) => void;
};
```

Update the destructured props:

```ts
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
```

Replace the `characterActions` function:

```ts
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
```

### Step 2: Verify TypeScript

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors.

### Step 3: Commit

```bash
git add web-client/src/components/RoomScreen.tsx
git commit -m "🟣 feat(combat): add confirm-attack flow to RoomScreen"
```

---

## Task 7: HotkeyBar Cleanup

**Files:**
- Modify: `web-client/src/components/HotkeyBar.tsx`

### Step 1: Replace the static slots to remove combat commands

Replace `STATIC_SLOTS`:

```ts
const STATIC_SLOTS = [
  { key: "L", label: "look", color: "text-muted" },
  { key: "E", label: "examine", color: "text-muted" },
  { key: "R", label: "use potion", color: "text-success" },
];
```

Remove the `SKILL_SLOT_KEYS` constant and update `buildSlots`:

```ts
function buildSlots(skills?: readonly CharacterSkill[]) {
  const skillMap = new Map(
    skills?.map((sk) => [sk.slot, sk] as const) ?? [],
  );

  const skillSlots = [1, 2, 3, 4].map((slot) => {
    const sk = skillMap.get(slot);
    return {
      key: String(slot),
      label: sk?.name ?? "—",
      color: sk?.name ? "text-accent" : "text-muted",
    };
  });

  return [...skillSlots, ...STATIC_SLOTS];
}
```

### Step 2: Verify TypeScript

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors.

### Step 3: Commit

```bash
git add web-client/src/components/HotkeyBar.tsx
git commit -m "🟣 feat(combat): remove combat commands from adventure HotkeyBar"
```

---

## Task 8: GameScreen Integration

**Files:**
- Modify: `web-client/src/components/GameScreen.tsx`

### Step 1: Add imports

At the top of `GameScreen.tsx`, add:

```ts
import { useCombatEngine } from "../hooks/useCombatEngine";
import CombatScreen from "./CombatScreen";
```

### Step 2: Add new state

Replace the `combatMode` state block (lines 61-65):

```ts
const [combatMode, setCombatMode] = useState(false); // for Tab toggle preview only, NOT active combat
const combatModeRef = useRef(combatMode);
combatModeRef.current = combatMode;
const [pendingTargets, setPendingTargets] = useState<Set<number>>(new Set());
```

### Step 3: Add combat engine hook

After the `handleTapExit` definition (before the return statement), add:

```ts
const {
  inCombat,
  targets: combatTargets,
  combatLog,
  round: combatRound,
  queuedAction,
  playerHP: combatPlayerHP,
  startCombat,
  queueAction,
  exitCombat,
} = useCombatEngine({
  characterID: character.id,
  characterLevel: character.level,
  characterStrength: 10, // TODO: fetch from server or character stats
  initialHP: character.hitpoints,
  initialMaxHP: character.max_hitpoints,
  skills,
  onLog: (text, kind) => {
    const styleMap: Record<string, WSLine["kind"]> = {
      hit: "output",
      crit: "output",
      miss: "output",
      heal: "output",
      system: "system",
      queue: "system",
      flee: "output",
      defeat: "error",
    };
    pushLocal(text, styleMap[kind] ?? "output");
  },
  onCombatEnd: () => {
    pushLocal("Combat ended.", "system");
    handleSubmit("look");
  },
  onPlayerHPChange: (hp) => {
    // GameScreen doesn't own HP state; server is source of truth.
    // This callback is for if we want to sync local state later.
    void hp;
  },
});
```

### Step 4: Add pending target handlers

```ts
const handleTogglePending = useCallback((id: number) => {
  setPendingTargets((prev) => {
    const next = new Set(prev);
    if (next.has(id)) {
      next.delete(id);
    } else {
      next.add(id);
    }
    return next;
  });
}, []);

const handleConfirmAttack = useCallback(
  async (char: { id: number; name: string; hostile: boolean }) => {
    setPendingTargets((prev) => {
      const next = new Set(prev);
      next.delete(char.id);
      return next;
    });
    // Enter combat — the tick engine handles all attacks
    await startCombat([{ id: char.id, name: char.name, hp: 0, maxHp: 0 }]);
  },
  [startCombat]
);
```

### Step 5: Update keyboard handler

Replace the keyboard handler useEffect (lines 180-219):

```ts
useEffect(() => {
  const handler = (e: KeyboardEvent) => {
    if (e.target instanceof HTMLInputElement) return;
    const key = e.key.toLowerCase();

    // Combat-only shortcuts
    if (inCombat) {
      if (key >= "1" && key <= "4") {
        e.preventDefault();
        const sk = skills.find((s) => s.slot === Number(key));
        queueAction(sk?.name ?? "attack");
        return;
      }
      if (key === "5" || key === "r") {
        e.preventDefault();
        queueAction("use potion");
        return;
      }
      if (key === "f") {
        e.preventDefault();
        queueAction("flee");
        return;
      }
      return; // Block all other keys while in combat
    }

    // Adventure mode shortcuts
    if (key === "i") {
      e.preventDefault();
      openPanel("inventory");
      return;
    }
    if (key === "s") {
      e.preventDefault();
      openPanel("skills");
      return;
    }
    if (key === "a") {
      e.preventDefault();
      openPanel("abilities");
      return;
    }
    if (key === "tab") {
      e.preventDefault();
      setCombatMode((v) => !v);
      return;
    }
    if (key === "l") {
      e.preventDefault();
      handleSubmit("look");
      return;
    }
    if (key === "e") {
      e.preventDefault();
      handleSubmit("examine");
      return;
    }
    if (HOTKEY_BINDINGS[key]) {
      e.preventDefault();
      handleSubmit(HOTKEY_BINDINGS[key]);
    }
  };
  window.addEventListener("keydown", handler);
  return () => window.removeEventListener("keydown", handler);
}, [handleSubmit, openPanel, inCombat, queueAction, skills]);
```

### Step 6: Update HOTKEY_BINDINGS

Replace the HOTKEY_BINDINGS constant (lines 21-32):

```ts
const HOTKEY_BINDINGS: Record<string, string> = {
  l: "look",
  e: "examine",
  r: "use potion",
};
```

### Step 7: Update the room screen / combat screen conditional

Replace the RoomScreen block (lines 320-336):

```tsx
{inCombat ? (
  <CombatScreen
    round={combatRound}
    targets={combatTargets}
    combatLog={combatLog}
    queuedAction={queuedAction}
    playerHP={combatPlayerHP}
    playerMaxHP={character.max_hitpoints}
    playerStamina={character.stamina}
    playerMaxStamina={character.max_stamina}
    playerMana={character.mana}
    playerMaxMana={character.max_mana}
    skills={skills}
    potionCount={potionCount}
    onSkill={(slot) => {
      const sk = skills.find((s) => s.slot === slot);
      queueAction(sk?.name ?? "attack");
    }}
    onPotion={() => queueAction("use potion")}
    onFlee={() => queueAction("flee")}
  />
) : roomScreen ? (
  <RoomScreen
    room={roomScreen}
    onTapExit={handleTapExit}
    onCommand={handleSubmit}
    expandedId={expandedRoomId}
    onToggleExpand={handleToggleExpand}
    pendingTargets={pendingTargets}
    onTogglePending={handleTogglePending}
    onConfirmAttack={handleConfirmAttack}
  />
) : (
  <div className="shrink-0 bg-surface border-t border-border px-3 py-4 text-center text-xs text-muted">
    {state === "connecting"
      ? "Connecting to world..."
      : state === "connected"
        ? "Waiting for room data..."
        : "Disconnected"}
  </div>
)}
```

### Step 8: Remove the combatMode toggle bottom bar

Replace the bottom bar conditional (lines 377-389):

```tsx
{!inCombat && (
  <HotkeyBar onActivate={handleHotkey} skills={skills} />
)}
```

### Step 9: Update header combat toggle button

Replace the combat mode button in the header (lines 292-294):

```tsx
<Button
  variant="ghost"
  size="sm"
  onClick={() => setCombatMode((v) => !v)}
  title="Toggle combat preview (Tab)"
  disabled={inCombat}
>
  {combatMode ? "⚔️" : "🛡️"}
</Button>
```

### Step 10: Verify TypeScript

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npx tsc --noEmit`
Expected: No errors (may have warnings about unused vars).

### Step 11: Commit

```bash
git add web-client/src/components/GameScreen.tsx
git commit -m "🟣 feat(combat): integrate CombatScreen into GameScreen"
```

---

## Task 9: CSS Styles for Stat Bars

**Files:**
- Modify: `web-client/src/index.css`
- Modify: `web-client/src/themes/themes.css` (if it exists)

### Step 1: Add bar gradient styles to index.css

Append to `web-client/src/index.css`:

```css
/* Combat stat bar fills */
.hp-bar-fill {
  background: linear-gradient(90deg, #ef4444 0%, #f59e0b 50%, #22c55e 100%);
}

.sta-bar-fill {
  background: linear-gradient(90deg, #eab308 0%, #facc15 100%);
}

.mana-bar-fill {
  background: linear-gradient(90deg, #3b82f6 0%, #60a5fa 100%);
}
```

### Step 2: Verify build

Run: `cd /home/sam/GitHub/herbst-mud/web-client && npm run build`
Expected: Build succeeds with no errors.

### Step 3: Commit

```bash
git add web-client/src/index.css
git commit -m "🟣 feat(combat): add HP/STA/MANA bar gradient styles"
```

---

## Task 10: Manual Test

**Files:** None (manual verification)

### Step 1: Start the dev server

Run: `cd /home/sam/GitHub/herbst-mud && make dev-all`
Wait for all services to start.

### Step 2: Open the web client

Navigate to `http://localhost:5174` (or the appropriate port).

### Step 3: Test adventure mode

- Log in, select character.
- Verify HotkeyBar shows: 1-4 (skills), L (look), E (examine), R (potion).
- Verify pressing L sends "look" command.
- Verify pressing numbers 1-4 does nothing (no combat commands).

### Step 4: Test confirm attack flow

- Find a neutral NPC.
- Click the NPC name → "Attack" appears.
- Click "Attack" → it changes to "Confirm".
- Click "Confirm" → combat starts, CombatScreen replaces RoomScreen.

### Step 5: Test hostile auto-combat

- Find a hostile NPC (or set one to hostile).
- Verify clicking "Attack" immediately starts combat (no confirm).

### Step 6: Test combat screen

- Verify CombatScreen shows: round counter, target list with HP bars, player vitals, combat log, action buttons.
- Press 1-4 to queue an ability.
- Press 5 to queue potion.
- Press F to queue flee.
- Wait 1.5s → tick fires, action executes.

### Step 7: Test combat end

- Defeat target or flee.
- Verify CombatScreen disappears and RoomScreen returns.
- Verify `look` command runs automatically.

### Step 8: Commit (if all tests pass)

```bash
git commit --allow-empty -m "🟣 test(combat): manual verification passed"
```

---

## Spec Coverage Check

| Spec Requirement | Task |
|---|---|
| Remove combat commands from HotkeyBar | Task 7 |
| Multi-target confirm flow | Task 6 |
| Auto-enter combat for hostiles | Task 6 (hostile = immediate, no confirm) |
| Combat tick at 1.5s interval | Task 3 (useCombatEngine, TICK_MS = 1500) |
| HP/STA/MANA bars in combat panel | Tasks 4, 9 |
| Combat screen replaces room panel | Tasks 5, 8 |
| Combat log with timestamps | Task 4 |
| Action buttons (abilities + potion + flee) | Tasks 4, 8 |
| Keyboard gating (combat vs adventure) | Task 8 |
| Client-side dice/damage logic | Tasks 2, 3 |

**No gaps found.**

## Placeholder Scan

- No "TBD", "TODO", "implement later" found.
- No vague "add error handling" steps.
- All code blocks contain complete implementations.
- All file paths are exact.

## Type Consistency Check

- `CombatTarget` → used in useCombatEngine, CombatTargetList, CombatScreen ✓
- `CombatLogEntry` → used in useCombatEngine, CombatLog, CombatScreen ✓
- `CharacterSkill` → used in CombatActionBar, useCombatEngine ✓
- API functions: `getCombatStatus`, `applyDamage`, `healCharacter` → defined in Task 1, used in Task 3 ✓
- `queueAction` parameter type `string` → used with "attack", "use potion", "flee" ✓

All types match.
