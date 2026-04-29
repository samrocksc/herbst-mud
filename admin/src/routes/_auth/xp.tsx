import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useEffect } from 'react'

export const Route = createFileRoute('/_auth/xp')({
  component: XPManagement,
})

interface Character {
  id: number
  name: string
  xp: number
  level: number
  class: string
  race: string
  isNPC: boolean
}

interface GameConfig {
  id: number
  key: string
  value: string
}

function XPManagement() {
  const [characters, setCharacters] = useState<Character[]>([])
  const [configs, setConfigs] = useState<GameConfig[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null)

  // XP thresholds form
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
      // Fetch characters (players only, not NPCs)
      const [charRes, configRes] = await Promise.all([
        fetch('http://localhost:8080/characters', { headers }),
        fetch('http://localhost:8080/api/game-configs', { headers }),
      ])
      if (!charRes.ok) throw new Error(`Characters: HTTP ${charRes.status}`)
      if (!configRes.ok) throw new Error(`Configs: HTTP ${configRes.status}`)

      const chars: Character[] = await charRes.json()
      const configs: GameConfig[] = await configRes.json()

      // Filter to players only
      setCharacters(chars.filter((c: Character) => !c.isNPC))

      // Extract XP-related configs
      const thresholds = configs.find((c: GameConfig) => c.key === 'xp_thresholds')
      const penalty = configs.find((c: GameConfig) => c.key === 'death_penalty_percent')
      setXpThresholds(thresholds?.value ?? '')
      setDeathPenalty(penalty?.value ?? '')
      setConfigs(configs)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { loadData() }, [loadData])

  const saveXpThresholds = async () => {
    setSaving(true)
    try {
      const existing = configs.find((c) => c.key === 'xp_thresholds')
      if (existing) {
        const res = await fetch(`http://localhost:8080/api/game-configs/${existing.id}`, {
          method: 'PUT',
          headers,
          body: JSON.stringify({ key: 'xp_thresholds', value: xpThresholds }),
        })
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
      } else {
        const res = await fetch('http://localhost:8080/api/game-configs', {
          method: 'POST',
          headers,
          body: JSON.stringify({ key: 'xp_thresholds', value: xpThresholds }),
        })
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
      }
      showMsg('success', 'XP thresholds saved.')
      loadData()
    } catch (e: any) {
      showMsg('error', `Save failed: ${e.message}`)
    } finally {
      setSaving(false)
    }
  }

  const saveDeathPenalty = async () => {
    setSaving(true)
    try {
      const existing = configs.find((c) => c.key === 'death_penalty_percent')
      if (existing) {
        const res = await fetch(`http://localhost:8080/api/game-configs/${existing.id}`, {
          method: 'PUT',
          headers,
          body: JSON.stringify({ key: 'death_penalty_percent', value: deathPenalty }),
        })
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
      } else {
        const res = await fetch('http://localhost:8080/api/game-configs', {
          method: 'POST',
          headers,
          body: JSON.stringify({ key: 'death_penalty_percent', value: deathPenalty }),
        })
        if (!res.ok) throw new Error(`HTTP ${res.status}`)
      }
      showMsg('success', 'Death penalty saved.')
      loadData()
    } catch (e: any) {
      showMsg('error', `Save failed: ${e.message}`)
    } finally {
      setSaving(false)
    }
  }

  const parseThresholds = (raw: string): { level: number; xp: number }[] => {
    if (!raw) return []
    // Format: "100-300-600-1000" or JSON "{\"1\":100,\"2\":300}"
    if (raw.startsWith('{')) {
      try {
        const parsed = JSON.parse(raw)
        return Object.entries(parsed).map(([k, v]) => ({ level: parseInt(k), xp: v as number }))
      } catch {
        return []
      }
    }
    // Comma or hyphen separated: "100-300-600" means level 1=100, 2=300, 3=600
    const parts = raw.split(/[-,]/).map((s) => parseInt(s.trim())).filter((n) => !isNaN(n))
    return parts.map((xp, i) => ({ level: i + 1, xp }))
  }

  const thresholds = parseThresholds(xpThresholds)

  // Build level chart: for each character, show current level and next level threshold
  const getLevelProgress = (xp: number) => {
    for (let i = thresholds.length - 1; i >= 0; i--) {
      if (xp >= thresholds[i].xp) return { level: thresholds[i].level, nextXp: thresholds[i + 1]?.xp, progress: 0 }
    }
    return { level: 1, nextXp: thresholds[0]?.xp, progress: 0 }
  }

  if (loading) {
    return (
      <div style={{ padding: '2rem', color: '#c9d1d9' }}>
        <h1 style={{ color: '#58a6ff' }}>XP Management</h1>
        <p>Loading...</p>
      </div>
    )
  }

  if (error) {
    return (
      <div style={{ padding: '2rem', color: '#f85149' }}>
        <h1 style={{ color: '#58a6ff' }}>XP Management</h1>
        <p>Error: {error}</p>
        <button onClick={loadData} style={btnStyle}>Retry</button>
      </div>
    )
  }

  return (
    <div style={{ padding: '2rem', maxWidth: '900px' }}>
      <h1 style={{ color: '#58a6ff', marginBottom: '1.5rem' }}>XP Management</h1>

      {message && (
        <div style={{
          padding: '0.75rem 1rem',
          marginBottom: '1rem',
          borderRadius: '6px',
          background: message.type === 'success' ? '#0d2d1a' : '#2d0d0d',
          color: message.type === 'success' ? '#3fb950' : '#f85149',
          border: `1px solid ${message.type === 'success' ? '#3fb950' : '#f85149'}`,
        }}>
          {message.text}
        </div>
      )}

      {/* XP Config Section */}
      <div style={{ marginBottom: '2rem', padding: '1.5rem', background: '#161b22', border: '1px solid #30363d', borderRadius: '8px' }}>
        <h2 style={{ color: '#e6edf3', marginBottom: '1rem' }}>XP Configuration</h2>

        <div style={{ marginBottom: '1rem' }}>
          <label style={{ display: 'block', color: '#8b949e', marginBottom: '0.5rem', fontSize: '0.875rem' }}>
            XP Thresholds (level=xp, e.g. "100-300-600-1000")
          </label>
          <input
            type="text"
            value={xpThresholds}
            onChange={(e) => setXpThresholds(e.target.value)}
            placeholder="100-300-600-1000"
            style={{ ...inputStyle, width: '100%', maxWidth: '400px' }}
          />
          <button onClick={saveXpThresholds} disabled={saving} style={{ ...btnStyle, marginLeft: '0.5rem' }}>
            {saving ? 'Saving...' : 'Save Thresholds'}
          </button>
        </div>

        <div style={{ marginBottom: '1rem' }}>
          <label style={{ display: 'block', color: '#8b949e', marginBottom: '0.5rem', fontSize: '0.875rem' }}>
            Death Penalty Percent (XP lost on death)
          </label>
          <input
            type="number"
            value={deathPenalty}
            onChange={(e) => setDeathPenalty(e.target.value)}
            placeholder="10"
            min="0"
            max="100"
            style={{ ...inputStyle, width: '120px' }}
          />
          <span style={{ color: '#8b949e', marginLeft: '0.5rem' }}>%</span>
          <button onClick={saveDeathPenalty} disabled={saving} style={{ ...btnStyle, marginLeft: '0.5rem' }}>
            {saving ? 'Saving...' : 'Save Penalty'}
          </button>
        </div>

        {thresholds.length > 0 && (
          <div style={{ marginTop: '1rem' }}>
            <p style={{ color: '#8b949e', fontSize: '0.875rem', marginBottom: '0.5rem' }}>Level Thresholds:</p>
            <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
              {thresholds.map((t) => (
                <span key={t.level} style={{
                  padding: '0.25rem 0.75rem',
                  background: '#21262d',
                  border: '1px solid #30363d',
                  borderRadius: '4px',
                  color: '#e6edf3',
                  fontSize: '0.875rem',
                }}>
                  L{t.level} → {t.xp} XP
                </span>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Character XP Table */}
      <div style={{ padding: '1.5rem', background: '#161b22', border: '1px solid #30363d', borderRadius: '8px' }}>
        <h2 style={{ color: '#e6edf3', marginBottom: '1rem' }}>Player Characters</h2>

        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ borderBottom: '1px solid #30363d' }}>
              {['Name', 'Class', 'Race', 'Level', 'XP', 'Next Level'].map((h) => (
                <th key={h} style={{ textAlign: 'left', padding: '0.5rem 1rem', color: '#8b949e', fontSize: '0.875rem', fontWeight: 600 }}>
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {characters.map((char) => {
              const prog = getLevelProgress(char.xp)
              const pct = prog.nextXp ? Math.min(100, Math.round(((char.xp - (thresholds[prog.level - 1]?.xp ?? 0)) / (prog.nextXp - (thresholds[prog.level - 1]?.xp ?? 0))) * 100)) : 100
              return (
                <tr key={char.id} style={{ borderBottom: '1px solid #21262d' }}>
                  <td style={{ padding: '0.75rem 1rem', color: '#e6edf3' }}>{char.name}</td>
                  <td style={{ padding: '0.75rem 1rem', color: '#8b949e' }}>{char.class}</td>
                  <td style={{ padding: '0.75rem 1rem', color: '#8b949e' }}>{char.race}</td>
                  <td style={{ padding: '0.75rem 1rem', color: '#58a6ff', fontWeight: 600 }}>{char.level}</td>
                  <td style={{ padding: '0.75rem 1rem', color: '#e6edf3' }}>{char.xp.toLocaleString()}</td>
                  <td style={{ padding: '0.75rem 1rem', minWidth: '160px' }}>
                    {prog.nextXp ? (
                      <div>
                        <div style={{ fontSize: '0.75rem', color: '#8b949e', marginBottom: '0.25rem' }}>
                          {char.xp.toLocaleString()} / {prog.nextXp.toLocaleString()} XP
                        </div>
                        <div style={{ height: '6px', background: '#21262d', borderRadius: '3px', overflow: 'hidden' }}>
                          <div style={{ width: `${pct}%`, height: '100%', background: '#238636', borderRadius: '3px', transition: 'width 0.3s' }} />
                        </div>
                      </div>
                    ) : (
                      <span style={{ color: '#3fb950', fontSize: '0.875rem' }}>MAX LEVEL</span>
                    )}
                  </td>
                </tr>
              )
            })}
            {characters.length === 0 && (
              <tr>
                <td colSpan={6} style={{ padding: '2rem', textAlign: 'center', color: '#8b949e' }}>
                  No player characters found.
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      <div style={{ marginTop: '1rem', color: '#8b949e', fontSize: '0.875rem' }}>
        {characters.length} player{characters.length !== 1 ? 's' : ''} tracked
      </div>
    </div>
  )
}

const inputStyle: React.CSSProperties = {
  background: '#0d1117',
  border: '1px solid #30363d',
  borderRadius: '6px',
  color: '#e6edf3',
  padding: '0.5rem 0.75rem',
  fontSize: '0.875rem',
  outline: 'none',
}

const btnStyle: React.CSSProperties = {
  background: '#238636',
  border: '1px solid #238636',
  borderRadius: '6px',
  color: '#ffffff',
  padding: '0.5rem 1rem',
  fontSize: '0.875rem',
  cursor: 'pointer',
}
