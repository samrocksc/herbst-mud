import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useEffect } from 'react'
import type { ReactNode } from 'react'
import { Button } from '../../components/Button'
import { DataTable } from '../../components/DataTable'

export const Route = createFileRoute('/_auth/config')({
  component: ConfigManagement,
})

type GameConfig = Readonly<{
  id: number
  key: string
  value: string
}>

type ConfigForm = Readonly<{
  key: string
  value: string
}>

/** Human-readable labels for known config keys */
const KEY_LABELS: Readonly<Record<string, string>> = {
  xp_thresholds: 'XP Thresholds',
  death_penalty_percent: 'Death Penalty (%)',
  starting_room_id: 'Starting Room ID',
  max_level: 'Max Level',
  campaign_type: 'Campaign Type',
  corpse_rot_minutes: 'Corpse Rot Minutes',
  xp_per_kill: 'XP Per Kill',
  respawn_room_id: 'Respawn Room ID',
  max_inventory_size: 'Max Inventory Size',
  pvp_enabled: 'PvP Enabled',
  death_penalty_type: 'Death Penalty Type',
  starting_hp: 'Starting HP',
  starting_mana: 'Starting Mana',
  starting_gold: 'Starting Gold',
  regen_tick_seconds: 'Regen Tick (seconds)',
}

/**
 * Convert a raw config key to a human-readable label.
 * If an explicit mapping exists in KEY_LABELS, use it.
 * Otherwise, split on underscores and title-case each word.
 */
function humanizeKey(raw: string): string {
  if (KEY_LABELS[raw]) return KEY_LABELS[raw]
  return raw
    .split('_')
    .map(segment => segment.charAt(0).toUpperCase() + segment.slice(1))
    .join(' ')
}

/**
 * Attempt to parse a string as JSON. Returns the parsed object/array
 * if successful, or null if it's not valid JSON (i.e. it's a plain string/number).
 */
function tryParseJSON(value: string): unknown | null {
  try {
    const parsed = JSON.parse(value)
    if (typeof parsed === 'object' && parsed !== null) return parsed
    return null
  } catch {
    return null
  }
}

/** Format a parsed JSON object/array with 2-space indentation */
function formatJSON(value: unknown): string {
  return JSON.stringify(value, null, 2)
}

const PRESETS = [
  { label: 'XP Thresholds', key: 'xp_thresholds', value: '{"1":100,"2":300,"3":600,"4":1000,"5":1500}' },
  { label: 'Death Penalty %', key: 'death_penalty_percent', value: '10' },
  { label: 'Corpse Rot Minutes', key: 'corpse_rot_minutes', value: '5' },
  { label: 'XP Per Kill (default)', key: 'xp_per_kill', value: '50' },
  { label: 'Max Level', key: 'max_level', value: '100' },
  { label: 'Starting Room ID', key: 'starting_room_id', value: '1' },
]

/** Renders a config value cell: pretty-prints JSON objects/arrays, truncates long plain strings */
function ConfigValueCell({ value }: { value: string }) {
  const parsed = tryParseJSON(value)

  if (parsed !== null) {
    const formatted = formatJSON(parsed)
    const isLong = formatted.split('\n').length > 4
    const [expanded, setExpanded] = useState(false)

    return (
      <div className="text-xs">
        <pre className={`font-mono text-text-secondary whitespace-pre-wrap m-0 ${!expanded ? 'max-h-16 overflow-hidden' : ''}`}>
          {formatted}
        </pre>
        {isLong && (
          <button
            type="button"
            className="text-primary text-xs mt-1 hover:underline cursor-pointer"
            onClick={() => setExpanded(e => !e)}
          >
            {expanded ? 'Show less' : 'Show more'}
          </button>
        )}
      </div>
    )
  }

  return (
    <span className="inline-block max-w-md overflow-hidden text-ellipsis whitespace-nowrap text-text-secondary text-xs">
      {value.length > 60 ? value.slice(0, 60) + '…' : value}
    </span>
  )
}

/** Collapsible JSON preview shown above the textarea in the edit modal */
function CollapsibleJSONPreview({ value }: { value: string }) {
  const parsed = tryParseJSON(value)
  const [expanded, setExpanded] = useState(false)

  if (parsed === null) return null

  const formatted = formatJSON(parsed)

  return (
    <div className="mb-3">
      <button
        type="button"
        className="text-xs text-primary hover:underline cursor-pointer flex items-center gap-1 mb-1"
        onClick={() => setExpanded(e => !e)}
      >
        <span className={`inline-block transition-transform ${expanded ? 'rotate-90' : ''}`}>▶</span>
        {expanded ? 'Collapse JSON preview' : 'Expand JSON preview'}
      </button>
      {expanded && (
        <pre className="bg-surface-muted border-2 border-border rounded p-3 text-xs font-mono whitespace-pre-wrap overflow-auto max-h-64">
          {formatted}
        </pre>
      )}
    </div>
  )
}

function ConfigManagement() {
  const [configs, setConfigs] = useState<GameConfig[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<GameConfig | null>(null)
  const [form, setForm] = useState<ConfigForm>({ key: '', value: '' })
  const [saving, setSaving] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<GameConfig | null>(null)

  const fetchConfigs = useCallback(async () => {
    setLoading(true)
    setError(null)
    const token = localStorage.getItem('token')
    try {
      const res = await fetch(`${window.location.origin}/api/game-configs`, {
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const data = await res.json()
      setConfigs(data)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }, [])

  // eslint-disable-next-line react-hooks/exhaustive-deps
  useEffect(() => { fetchConfigs() }, [fetchConfigs])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    const token = localStorage.getItem('token')
    try {
      const res = await fetch(`${window.location.origin}/api/game-configs`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify(form),
      })
      if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error || `HTTP ${res.status}`)
      }
      setShowForm(false)
      setForm({ key: '', value: '' })
      fetchConfigs()
    } catch (e: any) {
      alert(`Failed to create config: ${e.message}`)
    } finally {
      setSaving(false)
    }
  }

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!editing) return
    setSaving(true)
    const token = localStorage.getItem('token')
    try {
      const res = await fetch(`${window.location.origin}/api/game-configs/${editing.key}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({ value: form.value }),
      })
      if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error || `HTTP ${res.status}`)
      }
      setEditing(null)
      setForm({ key: '', value: '' })
      fetchConfigs()
    } catch (e: any) {
      alert(`Failed to update config: ${e.message}`)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget) return
    const token = localStorage.getItem('token')
    try {
      const res = await fetch(`${window.location.origin}/api/game-configs/${deleteTarget.key}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      setDeleteTarget(null)
      fetchConfigs()
    } catch (e: any) {
      alert(`Failed to delete: ${e.message}`)
    }
  }

  const startEdit = (cfg: GameConfig) => {
    setEditing(cfg)
    setForm({ key: cfg.key, value: cfg.value })
    setShowForm(false)
  }

  const startCreate = () => {
    setShowForm(true)
    setEditing(null)
    setForm({ key: '', value: '' })
  }

  const applyPreset = (preset: typeof PRESETS[0]) => {
    if (showForm) {
      setForm({ key: preset.key, value: preset.value })
    }
  }

  const filtered = configs.filter(c =>
    c.key.toLowerCase().includes(search.toLowerCase()) ||
    c.value.toLowerCase().includes(search.toLowerCase())
  )

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Game Configs</h2>
        <Button variant="primary" onClick={startCreate}>+ New Config</Button>
      </div>

      {error && <div className="error-banner">{error}</div>}

      <div className="flex items-center gap-3 mb-4">
        <input
          type="text"
          placeholder="Search keys or values..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          className="flex-1 px-3 py-2 bg-surface-muted border-2 border-border color-text rounded"
        />
        <Button variant="secondary" onClick={fetchConfigs}>Refresh</Button>
      </div>

      {loading ? (
        <div className="loading">Loading configs...</div>
      ) : (
        <DataTable
          columns={[
            {
              header: 'Key',
              accessor: 'key',
              render: (_, row): ReactNode => (
                <div>
                  <span className="text-text text-sm font-medium">{humanizeKey(row.key)}</span>
                  <br />
                  <code className="text-primary text-xs">{row.key}</code>
                </div>
              ),
            },
            {
              header: 'Value',
              accessor: 'value',
              render: (val) => <ConfigValueCell value={val as string} />,
            },
            {
              header: 'Actions',
              accessor: '_actions',
              render: (_, row): ReactNode => (
                <div className="flex gap-2">
                  <Button variant="accent" size="sm" onClick={() => startEdit(row)}>Edit</Button>
                  <Button variant="danger" size="sm" onClick={() => setDeleteTarget(row)}>Delete</Button>
                </div>
              ),
            },
          ]}
          data={filtered}
          getKey={(row) => row.id}
          emptyMessage={configs.length === 0 ? 'No configs found. Create one below.' : 'No configs match your search.'}
        />
      )}

      {/* Create / Edit Form Modal */}
      {(showForm || editing) && (
        <div className="modal-overlay" onClick={() => { setShowForm(false); setEditing(null) }}>
          <div className="modal-content max-w-2xl" onClick={e => e.stopPropagation()}>
            <h3>{editing ? `Edit: ${humanizeKey(editing.key)}` : 'New Game Config'}</h3>
            <form onSubmit={editing ? handleUpdate : handleCreate}>
              <div className="form-group">
                <label>Key</label>
                {editing ? (
                  <code className="block p-2 bg-surface-muted rounded">{editing.key}</code>
                ) : (
                  <input
                    type="text"
                    value={form.key}
                    onChange={e => setForm(f => ({ ...f, key: e.target.value }))}
                    placeholder="e.g. xp_thresholds"
                    required
                    className="w-full p-2 bg-surface-muted border-2 border-border color-text rounded"
                  />
                )}
              </div>
              <div className="form-group">
                <label>Value (JSON or plain)</label>
                <CollapsibleJSONPreview value={form.value} />
                <textarea
                  value={form.value}
                  onChange={e => setForm(f => ({ ...f, value: e.target.value }))}
                  className="w-full p-2 bg-surface-muted border-2 border-border color-text rounded font-mono text-sm"
                  rows={6}
                  placeholder='{"key": "value"} or plain text'
                  required
                />
              </div>
              {!editing && (
                <div className="form-group">
                  <label>Presets</label>
                  <div className="flex flex-wrap gap-2">
                    {PRESETS.map(p => (
                      <Button type="button" key={p.key} variant="ghost" size="sm" onClick={() => applyPreset(p)}>
                        {p.label}
                      </Button>
                    ))}
                  </div>
                </div>
              )}
              <div className="flex gap-3 justify-end mt-4">
                <Button type="button" variant="secondary" onClick={() => { setShowForm(false); setEditing(null) }}>Cancel</Button>
                <Button type="submit" variant="primary" disabled={saving}>{saving ? 'Saving…' : (editing ? 'Update' : 'Create')}</Button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Delete Confirmation */}
      {deleteTarget && (
        <div className="modal-overlay" onClick={() => setDeleteTarget(null)}>
          <div className="modal-content max-w-md" onClick={e => e.stopPropagation()}>
            <h3>Delete Config?</h3>
            <p>Are you sure you want to delete <code>{deleteTarget.key}</code>? This cannot be undone.</p>
            <div className="flex gap-3 justify-end mt-4">
              <Button variant="secondary" onClick={() => setDeleteTarget(null)}>Cancel</Button>
              <Button variant="danger" onClick={handleDelete}>Delete</Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}