/* eslint-disable react-refresh/only-export-components */

import { getLevelProgress, parseThresholds, normalizeChar } from "../../components/map/xpUtils";
import type { Character, ThresholdEntry, RawChar, GameConfig } from "../../components/map/xpUtils";

export type { Character, ThresholdEntry, RawChar, GameConfig };
export { getLevelProgress, parseThresholds, normalizeChar };

/** Renders an XP progress bar showing current xp / next-level xp */
export function XPProgressCell({ char, thresholds }: Readonly<{
  char: Character; thresholds: ThresholdEntry[]
}>) {
  const prog = getLevelProgress(char.xp, thresholds);
  if (!prog.nextXp) return <span className="xp-max">MAX LEVEL</span>;
  const prev = thresholds[prog.level - 1]?.xp ?? 0;
  const pct = Math.min(100, Math.round(((char.xp - prev) / (prog.nextXp - prev)) * 100));
  return (
    <div>
      <div className="xp-progress-label">{char.xp.toLocaleString()} / {prog.nextXp.toLocaleString()} XP</div>
      <div className="xp-progress-track"><div className="xp-progress-fill" style={{ width: `${pct}%` }} /></div>
    </div>
  );
}
