 
export type Character = Readonly<{
  id: number; name: string; xp: number; level: number
  class: string; race: string; isNPC: boolean
}>

export type ThresholdEntry = Readonly<{ level: number; xp: number }>

export type RawChar = Readonly<{
  ID?: number; Name?: string; Xp?: number; Level?: number
  Class?: string; Race?: string; IsNPC?: boolean; [key: string]: unknown
}>

export type GameConfig = Readonly<{ id: number; key: string; value: string }>

export function normalizeChar(c: RawChar): Character {
  return {
    id: Number(c.id ?? c.ID ?? 0), name: String(c.name ?? c.Name ?? ""),
    xp: Number(c.xp ?? c.Xp ?? 0), level: Number(c.level ?? c.Level ?? 1),
    class: String(c.class ?? c.Class ?? "unknown"), race: String(c.race ?? c.Race ?? "unknown"),
    isNPC: Boolean(c.isNPC ?? c.IsNPC ?? false),
  };
}

/** Parse XP thresholds: JSON object or dash/comma-separated numbers */
export function parseThresholds(raw: string): ThresholdEntry[] {
  if (!raw) return [];
  if (raw.startsWith("{")) {
    try {
      const parsed = JSON.parse(raw);
      return Object.entries(parsed).map(([k, v]) => ({ level: parseInt(k), xp: v as number }));
    } catch { return []; }
  }
  const parts = raw.split(/[-,]/).map(s => parseInt(s.trim())).filter(n => !isNaN(n));
  return parts.map((xp, i) => ({ level: i + 1, xp }));
}

/** Determine current level and next-level XP from thresholds */
export function getLevelProgress(xp: number, thresholds: ReadonlyArray<ThresholdEntry>) {
  const reversed = [...thresholds].reverse();
  const found = reversed.find(t => xp >= t.xp);
  const idx = found ? thresholds.indexOf(found) : -1;
  return {
    level: found?.level ?? 1,
    nextXp: idx >= 0 && idx < thresholds.length - 1 ? thresholds[idx + 1]?.xp : undefined,
  };
}
