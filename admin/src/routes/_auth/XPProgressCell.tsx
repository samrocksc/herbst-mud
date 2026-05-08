export type Character = Readonly<{
  id: number
  name: string
  xp: number
  level: number
  class: string
  race: string
  isNPC: boolean
}>

export type ThresholdEntry = Readonly<{ level: number; xp: number }>

/** Parse XP threshold strings: JSON object or dash/comma-separated numbers */
export function parseThresholds(raw: string): ThresholdEntry[] {
  if (!raw) return []
  if (raw.startsWith('{')) {
    try {
      const parsed = JSON.parse(raw)
      return Object.entries(parsed).map(([k, v]) => ({ level: parseInt(k), xp: v as number }))
    } catch { return [] }
  }
  const parts = raw.split(/[-,]/).map(s => parseInt(s.trim())).filter(n => !isNaN(n))
  return parts.map((xp, i) => ({ level: i + 1, xp }))
}

/** Determine current level and next-level XP from thresholds */
export function getLevelProgress(xp: number, thresholds: ThresholdEntry[]) {
  for (let i = thresholds.length - 1; i >= 0; i--) {
    if (xp >= thresholds[i].xp) return { level: thresholds[i].level, nextXp: thresholds[i + 1]?.xp }
  }
  return { level: 1, nextXp: thresholds[0]?.xp }
}

/** Renders an XP progress bar showing current xp / next-level xp */
export function XPProgressCell({ char, thresholds }: Readonly<{
  char: Character
  thresholds: ThresholdEntry[]
}>) {
  const prog = getLevelProgress(char.xp, thresholds)
  if (!prog.nextXp) return <span className="xp-max">MAX LEVEL</span>
  const prev = thresholds[prog.level - 1]?.xp ?? 0
  const pct = Math.min(100, Math.round(((char.xp - prev) / (prog.nextXp - prev)) * 100))
  return (
    <div>
      <div className="xp-progress-label">
        {char.xp.toLocaleString()} / {prog.nextXp.toLocaleString()} XP
      </div>
      <div className="xp-progress-track">
        <div className="xp-progress-fill" style={{ width: `${pct}%` }} />
      </div>
    </div>
  )
}