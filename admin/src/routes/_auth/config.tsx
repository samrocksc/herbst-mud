import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useEffect } from 'react'

export const Route = createFileRoute('/_auth/config')({
  component: ConfigManagement,
})

interface GameConfig {
  id: number
  key: string
  value: string
}

interface ConfigForm {
  key: string
  value: string
}

const PRESETS = [
  { label: 'XP Thresholds', key: 'xp_thresholds', value: '{"1":100,"2":300,"3":600,"4":1000,"5":1500}' },
  { label: 'Death Penalty %', key: 'death_penalty_percent', value: '10' },
  { label: 'Corpse Rot Minutes', key: 'corpse_rot_minutes', value: '5' },
  { label: 'XP Per Kill (default)', key: 'xp_per_kill', value: '50' },
  { label: 'Max Level', key: 'max_level', value: '100' },
  { label: 'Starting Room ID', key: 'starting_room_id', value: '1' },
]

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
      const res = await fetch('http://localhost:8080/api/game-configs', {
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
      const res = await fetch('http://localhost:8080/api/game-configs', {
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
      const res = await fetch(`http://localhost:8080/api/game-configs/${editing.key}`, {
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
      const res = await fetch(`http://localhost:8080/api/game-configs/${deleteTarget.key}`, {
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
        <button className="btn-primary" onClick={startCreate}>+ New Config</button>
      </div>

      {error && <div className="error-banner">{error}</div>}

      <div className="toolbar-row" style={{ marginBottom: '1rem', display: 'flex', gap: '0.75rem', alignItems: 'center' }}>
        <input
          type="text"
          placeholder="Search keys or values..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          style={{ flex: 1, padding: '0.5rem 0.75rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px' }}
        />
        <button className="btn-secondary" onClick={fetchConfigs}>Refresh</button>
      </div>

      {loading ? (
        <div className="loading">Loading configs...</div>
      ) : (
        <table className="data-table">
          <thead>
            <tr>
              <th>Key</th>
              <th>Value</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {filtered.length === 0 ? (
              <tr>
                <td colSpan={3} style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '2rem' }}>
                  {configs.length === 0 ? 'No configs found. Create one below.' : 'No configs match your search.'}
                </td>
              </tr>
            ) : (
              filtered.map(cfg => (
                <tr key={cfg.id}>
                  <td>
                    <code style={{ color: 'var(--primary)', fontSize: '0.9rem' }}>{cfg.key}</code>
                  </td>
                  <td>
                    <span style={{
                      display: 'inline-block',
                      maxWidth: '400px',
                      overflow: 'hidden',
                      textOverflow: 'ellipsis',
                      whiteSpace: 'nowrap',
                      color: 'var(--text-secondary)',
                      fontSize: '0.85rem',
                    }}>
                      {cfg.value.length > 60 ? cfg.value.slice(0, 60) + '…' : cfg.value}
                    </span>
                  </td>
                  <td>
                    <div style={{ display: 'flex', gap: '0.5rem' }}>
                      <button className="btn-small" onClick={() => startEdit(cfg)}>Edit</button>
                      <button className="btn-small btn-danger" onClick={() => setDeleteTarget(cfg)}>Delete</button>
                    </div>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      )}

      {/* Create / Edit Form Modal */}
      {(showForm || editing) && (
        <div className="modal-overlay" onClick={() => { setShowForm(false); setEditing(null) }}>
          <div className="modal-content" onClick={e => e.stopPropagation()} style={{ maxWidth: '560px' }}>
            <h3>{editing ? `Edit: ${editing.key}` : 'New Game Config'}</h3>
            <form onSubmit={editing ? handleUpdate : handleCreate}>
              <div className="form-group">
                <label>Key</label>
                {editing ? (
                  <code style={{ display: 'block', padding: '0.5rem', background: 'var(--surface-muted)', borderRadius: '4px' }}>{editing.key}</code>
                ) : (
                  <input
                    type="text"
                    value={form.key}
                    onChange={e => setForm(f => ({ ...f, key: e.target.value }))}
                    placeholder="e.g. xp_thresholds"
                    required
                    style={{ width: '100%', padding: '0.5rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px' }}
                  />
                )}
              </div>
              <div className="form-group">
                <label>Value (JSON or plain)</label>
                <textarea
                  value={form.value}
                  onChange={e => setForm(f => ({ ...f, value: e.target.value }))}
                  rows={6}
                  placeholder='{"key": "value"} or plain text'
                  required
                  style={{ width: '100%', padding: '0.5rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px', fontFamily: 'monospace', fontSize: '0.85rem' }}
                />
              </div>
              {!editing && (
                <div className="form-group">
                  <label>Presets</label>
                  <div style={{ display: 'flex', flexWrap: 'wrap', gap: '0.5rem' }}>
                    {PRESETS.map(p => (
                      <button type="button" key={p.key} className="btn-small" onClick={() => applyPreset(p)}>
                        {p.label}
                      </button>
                    ))}
                  </div>
                </div>
              )}
              <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end', marginTop: '1rem' }}>
                <button type="button" className="btn-secondary" onClick={() => { setShowForm(false); setEditing(null) }}>Cancel</button>
                <button type="submit" className="btn-primary" disabled={saving}>{saving ? 'Saving…' : (editing ? 'Update' : 'Create')}</button>
              </div>
            </form>
          </div>
        </div>
      )}

      {/* Delete Confirmation */}
      {deleteTarget && (
        <div className="modal-overlay" onClick={() => setDeleteTarget(null)}>
          <div className="modal-content" onClick={e => e.stopPropagation()} style={{ maxWidth: '400px' }}>
            <h3>Delete Config?</h3>
            <p>Are you sure you want to delete <code>{deleteTarget.key}</code>? This cannot be undone.</p>
            <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end', marginTop: '1rem' }}>
              <button className="btn-secondary" onClick={() => setDeleteTarget(null)}>Cancel</button>
              <button className="btn-danger" onClick={handleDelete}>Delete</button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
