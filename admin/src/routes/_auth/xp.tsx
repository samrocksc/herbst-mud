import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useEffect } from 'react'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { FormField, NumberField, FormError } from '../../components/FormFields'
import { showToast } from '../../components/Toast'
import { apiGet, apiPost, apiPut } from '../../utils/apiFetch'
import { XPProgressCell, parseThresholds } from './XPProgressCell'
import type { Character, ThresholdEntry } from './XPProgressCell'

export const Route = createFileRoute('/_auth/xp')({ component: XPManagement })

type RawChar = Readonly<{ ID?: number; Name?: string; Xp?: number; Level?: number; Class?: string; Race?: string; IsNPC?: boolean; [key: string]: unknown }>
type GameConfig = Readonly<{ id: number; key: string; value: string }>

function normalizeChar(c: RawChar): Character {
  return {
    id: Number(c.id ?? c.ID ?? 0), name: String(c.name ?? c.Name ?? ''),
    xp: Number(c.xp ?? c.Xp ?? 0), level: Number(c.level ?? c.Level ?? 1),
    class: String(c.class ?? c.Class ?? 'unknown'), race: String(c.race ?? c.Race ?? 'unknown'),
    isNPC: Boolean(c.isNPC ?? c.IsNPC ?? false),
  }
}

async function saveConfig(configs: GameConfig[], key: string, value: string) {
  const existing = configs.find(c => c.key === key)
  if (existing) return apiPut(`/api/game-configs/${existing.key}`, { key, value })
  return apiPost('/api/game-configs', { key, value })
}

function XPManagement() {
  const [characters, setCharacters] = useState<Character[]>([])
  const [configs, setConfigs] = useState<GameConfig[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [formError, setFormError] = useState('')
  const [xpThresholds, setXpThresholds] = useState('')
  const [deathPenalty, setDeathPenalty] = useState('')

  const loadData = useCallback(async () => {
    setLoading(true); setError(null)
    try {
      const [chars, cfg] = await Promise.all([
        apiGet<RawChar[]>('/characters'), apiGet<GameConfig[]>('/api/game-configs'),
      ])
      setCharacters(chars.map(normalizeChar).filter((c: Character) => !c.isNPC))
      setXpThresholds(cfg.find(c => c.key === 'xp_thresholds')?.value ?? '')
      setDeathPenalty(cfg.find(c => c.key === 'death_penalty_percent')?.value ?? '')
      setConfigs(cfg)
    } catch (e: unknown) { setError(e instanceof Error ? e.message : String(e)) }
    finally { setLoading(false) }
  }, [])

  useEffect(() => { loadData() }, [loadData])

  const handleSave = async (key: string, value: string, label: string) => {
    setSaving(true); setFormError('')
    try {
      await saveConfig(configs, key, value)
      showToast(`${label} saved.`, 'success')
      loadData()
    } catch (e: unknown) { setFormError(e instanceof Error ? e.message : String(e)) }
    finally { setSaving(false) }
  }

  const thresholds: ThresholdEntry[] = parseThresholds(xpThresholds)

  if (loading) return <div className="xp-page"><h1 className="xp-title">XP Management</h1><p className="xp-muted">Loading...</p></div>
  if (error) return (
    <div className="xp-page"><h1 className="xp-title">XP Management</h1>
      <p className="xp-error">Error: {error}</p><Button variant="primary" onClick={loadData}>Retry</Button></div>
  )

  const columns: Column<Character>[] = [
    { header: 'Name', accessor: 'name' },
    { header: 'Class', accessor: 'class' },
    { header: 'Race', accessor: 'race' },
    { header: 'Level', accessor: 'level', render: v => <span className="xp-level">{String(v)}</span> },
    { header: 'XP', accessor: 'xp', render: v => (v as number).toLocaleString() },
    { header: 'Next Level', accessor: 'id', render: (_, row) => <XPProgressCell char={row} thresholds={thresholds} /> },
  ]

  return (
    <div className="xp-page">
      <h1 className="xp-title">XP Management</h1>
      {formError && <FormError message={formError} />}
      <div className="xp-section">
        <h2 className="xp-section-title">XP Configuration</h2>
        <div className="xp-field">
          <FormField label="XP Thresholds" value={xpThresholds} onChange={setXpThresholds}
            placeholder='100-300-600-1000 or {"1":100,"2":300}' tooltip="Dash-separated or JSON level thresholds" />
          <Button variant="primary" onClick={() => handleSave('xp_thresholds', xpThresholds, 'XP thresholds')} disabled={saving}>
            {saving ? 'Saving...' : 'Save Thresholds'}
          </Button>
        </div>
        <div className="xp-field">
          <NumberField label="Death Penalty %" value={parseInt(deathPenalty) || 0}
            onChange={v => setDeathPenalty(String(v))} min={0} max={100} tooltip="XP lost on death (%)" />
          <Button variant="primary" onClick={() => handleSave('death_penalty_percent', deathPenalty, 'Death penalty')} disabled={saving}>
            {saving ? 'Saving...' : 'Save Penalty'}
          </Button>
        </div>
        {thresholds.length > 0 && (
          <div className="xp-thresholds">
            <p className="xp-label">Level Thresholds:</p>
            <div className="xp-threshold-list">
              {thresholds.map(t => <span key={t.level} className="xp-badge">L{t.level} → {t.xp} XP</span>)}
            </div>
          </div>
        )}
      </div>
      <div className="xp-section">
        <h2 className="xp-section-title">Player Characters</h2>
        <DataTable columns={columns} data={characters} getKey={r => r.id} variant="dark" emptyMessage="No player characters found." />
      </div>
      <div className="xp-count">{characters.length} player{characters.length !== 1 ? 's' : ''} tracked</div>
    </div>
  )
}