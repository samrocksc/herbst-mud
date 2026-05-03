import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useEffect } from 'react'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/xp')({
  component: XPManagement,
})

// ─── Types ───────────────────────────────────────────────────────────────────

// Ent serializes fields with Go-exported names: Xp, Level, Class, Race, IsNPC
// Normalize to the lowercase names the frontend expects
type RawCharacter = Readonly<{
  ID?: number
  Name?: string
  Xp?: number
  Level?: number
  Class?: string
  Race?: string
  IsNPC?: boolean
  [key: string]: unknown
}>

type Character = Readonly<{
  id: number
  name: string
  xp: number
  level: number
  class: string
  race: string
  isNPC: boolean
}>

type GameConfig = Readonly<{
  id: number
  key: string
  value: string
}>

// ─── Normalization ────────────────────────────────────────────────────────────

function normalizeChar(c: RawCharacter): Character {
  return {
    id: c.ID ?? 0,
    name: c.Name ?? '',
    xp: c.Xp ?? 0,
    level: c.Level ?? 1,
    class: c.Class ?? 'unknown',
    race: c.Race ?? 'unknown',
    isNPC: c.IsNPC ?? false,
  }
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

function parseThresholds(raw: string): ReadonlyArray<Readonly<{ level: number; xp: number }>> {
  if (!raw) return []
  if (raw.startsWith('{')) {
    try {
      const parsed = JSON.parse(raw)
      return Object.entries(parsed).map(([k, v]) => ({ level: parseInt(k), xp: v as number }))
    } catch {
      return []
    }
  }
  const parts = raw.split(/[-,]/).map((s) => parseInt(s.trim())).filter((n) => !isNaN(n))
  return parts.map((xp, i) => ({ level: i + 1, xp }))
}

function getLevelProgress(xp: number, thresholds: ReadonlyArray<Readonly<{ level: number; xp: number }>>) {
  for (let i = thresholds.length - 1; i >= 0; i--) {
    if (xp >= thresholds[i].xp) return { level: thresholds[i].level, nextXp: thresholds[i + 1]?.xp }
  }
  return { level: 1, nextXp: thresholds[0]?.xp }
}

// ─── Progress bar cell (XP column) ───────────────────────────────────────────

type XPProgressCellProps = Readonly<{
  char: Character
  thresholds: ReadonlyArray<Readonly<{ level: number; xp: number }>>
}>

function XPProgressCell({ char, thresholds }: XPProgressCellProps) {
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

// ─── Component ───────────────────────────────────────────────────────────────

function XPManagement() {
  const [characters, setCharacters] = useState<Character[]>([])
  const [configs, setConfigs] = useState<GameConfig[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState<Readonly<{ type: 'success' | 'error'; text: string }> | null>(null)
  const [xpThresholds, setXpThresholds] = useState('')
  const [deathPenalty, setDeathPenalty] = useState('')

  const token = localStorage.getItem('token')
  const headers = { Authorization: `Bearer ${token}`, 'Content-Type': 'application/json' }

  const showMsg = (type: 'success' | 'error', text: string) => {
    setMessage({ type, text })
    setTimeout(() => setMessage(null), 4000)
  }

  const loadData = useCallback(async () => {
    setLoading(true)
    setError(null)
    try {
      const [charRes, configRes] = await Promise.all([
        fetch(`${window.location.origin}/characters`, { headers }),
        fetch(`${window.location.origin}/api/game-configs`, { headers }),
      ])
      if (!charRes.ok) throw new Error(`Characters: HTTP ${charRes.status}`)
      if (!configRes.ok) throw new Error(`Configs: HTTP ${configRes.status}`)

      const chars: RawCharacter[] = await charRes.json()
      const cfg: GameConfig[] = await configRes.json()

      setCharacters(chars.map(normalizeChar).filter((c: Character) => !c.isNPC))

      const thresholds = cfg.find((c) => c.key === 'xp_thresholds')
      const penalty = cfg.find((c) => c.key === 'death_penalty_percent')
      setXpThresholds(thresholds?.value ?? '')
      setDeathPenalty(penalty?.value ?? '')
      setConfigs(cfg)
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : String(e))
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { loadData() }, [loadData])

  const saveXpThresholds = async () => {
    setSaving(true)
    try {
      const existing = configs.find((c) => c.key === 'xp_thresholds')
      const method = existing ? 'PUT' : 'POST'
      const url = existing
        ? `${window.location.origin}/api/game-configs/${existing.key}`
        : `${window.location.origin}/api/game-configs`
      const res = await fetch(url, {
        method,
        headers,
        body: JSON.stringify({ key: 'xp_thresholds', value: xpThresholds }),
      })
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      showMsg('success', 'XP thresholds saved.')
      loadData()
    } catch (e: unknown) {
      showMsg('error', `Save failed: ${e instanceof Error ? e.message : String(e)}`)
    } finally {
      setSaving(false)
    }
  }

  const saveDeathPenalty = async () => {
    setSaving(true)
    try {
      const existing = configs.find((c) => c.key === 'death_penalty_percent')
      const method = existing ? 'PUT' : 'POST'
      const url = existing
        ? `${window.location.origin}/api/game-configs/${existing.key}`
        : `${window.location.origin}/api/game-configs`
      const res = await fetch(url, {
        method,
        headers,
        body: JSON.stringify({ key: 'death_penalty_percent', value: deathPenalty }),
      })
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      showMsg('success', 'Death penalty saved.')
      loadData()
    } catch (e: unknown) {
      showMsg('error', `Save failed: ${e instanceof Error ? e.message : String(e)}`)
    } finally {
      setSaving(false)
    }
  }

  const thresholds = parseThresholds(xpThresholds)

  if (loading) {
    return (
      <div className="xp-page">
        <h1 className="xp-title">XP Management</h1>
        <p className="xp-muted">Loading...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="xp-page">
        <h1 className="xp-title">XP Management</h1>
        <p className="xp-error">Error: {error}</p>
        <Button variant="primary" onClick={loadData}>Retry</Button>
      </div>
    )
  }

  const columns: Column<Character>[] = [
    { header: 'Name', accessor: 'name' },
    { header: 'Class', accessor: 'class' },
    { header: 'Race', accessor: 'race' },
    {
      header: 'Level',
      accessor: 'level',
      render: (val) => <span className="xp-level">{String(val)}</span>,
    },
    { header: 'XP', accessor: 'xp', render: (val) => (val as number).toLocaleString() },
    {
      header: 'Next Level',
      accessor: 'id',
      render: (_, row) => <XPProgressCell char={row} thresholds={thresholds} />,
    },
  ]

  return (
    <div className="xp-page">
      <h1 className="xp-title">XP Management</h1>

      {message && (
        <div className={`xp-msg xp-msg-${message.type}`}>
          {message.text}
        </div>
      )}

      {/* XP Config Section */}
      <div className="xp-section">
        <h2 className="xp-section-title">XP Configuration</h2>

        <div className="xp-field">
          <label className="xp-label">
            XP Thresholds (level=xp, e.g. "100-300-600-1000")
          </label>
          <div className="xp-field-row">
            <input
              type="text"
              value={xpThresholds}
              onChange={(e) => setXpThresholds(e.target.value)}
              placeholder="100-300-600-1000"
              className="xp-input w-full max-w-[400px]"
            />
            <Button variant="primary" onClick={saveXpThresholds} disabled={saving}>
              {saving ? 'Saving...' : 'Save Thresholds'}
            </Button>
          </div>
        </div>

        <div className="xp-field">
          <label className="xp-label">
            Death Penalty Percent (XP lost on death)
          </label>
          <div className="xp-field-row">
            <input
              type="number"
              value={deathPenalty}
              onChange={(e) => setDeathPenalty(e.target.value)}
              placeholder="10"
              min="0"
              max="100"
              className="xp-input w-[120px]"
            />
            <span className="xp-muted">%</span>
            <Button variant="primary" onClick={saveDeathPenalty} disabled={saving}>
              {saving ? 'Saving...' : 'Save Penalty'}
            </Button>
          </div>
        </div>

        {thresholds.length > 0 && (
          <div className="xp-thresholds">
            <p className="xp-label">Level Thresholds:</p>
            <div className="xp-threshold-list">
              {thresholds.map((t) => (
                <span key={t.level} className="xp-badge">
                  L{t.level} → {t.xp} XP
                </span>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Character XP Table */}
      <div className="xp-section">
        <h2 className="xp-section-title">Player Characters</h2>

        <DataTable
          columns={columns}
          data={characters}
          getKey={(row) => row.id}
          variant="dark"
          emptyMessage="No player characters found."
        />
      </div>

      <div className="xp-count">
        {characters.length} player{characters.length !== 1 ? 's' : ''} tracked
      </div>
    </div>
  )
}
