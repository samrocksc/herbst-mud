import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect, useCallback } from 'react'

export const Route = createFileRoute('/_auth/races')({
  component: RacesManagement,
})

interface Race {
  id: number
  name: string
  display_name: string
  description: string
  stat_modifiers: string
  skill_grants: string
  is_playable: boolean
}

interface RaceForm {
  name: string
  display_name: string
  description: string
  stat_modifiers: string
  skill_grants: string
  is_playable: boolean
}

const emptyForm: RaceForm = {
  name: '',
  display_name: '',
  description: '',
  stat_modifiers: '{}',
  skill_grants: '[]',
  is_playable: true,
}

function RacesManagement() {
  const [races, setRaces] = useState<Race[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [search, setSearch] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [editing, setEditing] = useState<Race | null>(null)
  const [form, setForm] = useState<RaceForm>(emptyForm)
  const [saving, setSaving] = useState(false)
  const [deleteTarget, setDeleteTarget] = useState<Race | null>(null)

  const fetchRaces = useCallback(async () => {
    setLoading(true)
    setError(null)
    const token = localStorage.getItem('token')
    try {
      const res = await fetch('http://localhost:8080/api/races', {
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const data = await res.json()
      setRaces(data)
    } catch (e: any) {
      setError(e.message)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { fetchRaces() }, [fetchRaces])

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setSaving(true)
    const token = localStorage.getItem('token')
    try {
      const res = await fetch('http://localhost:8080/api/races', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify(form),
      })
      if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error || `HTTP ${res.status}`)
      }
      setShowForm(false)
      setForm(emptyForm)
      fetchRaces()
    } catch (e: any) {
      alert(`Failed to create race: ${e.message}`)
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
      const res = await fetch(`http://localhost:8080/api/races/${editing.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({
          display_name: form.display_name,
          description: form.description,
          stat_modifiers: form.stat_modifiers,
          skill_grants: form.skill_grants,
          is_playable: form.is_playable,
        }),
      })
      if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error || `HTTP ${res.status}`)
      }
      setEditing(null)
      setForm(emptyForm)
      fetchRaces()
    } catch (e: any) {
      alert(`Failed to update race: ${e.message}`)
    } finally {
      setSaving(false)
    }
  }

  const handleDelete = async () => {
    if (!deleteTarget) return
    const token = localStorage.getItem('token')
    try {
      const res = await fetch(`http://localhost:8080/api/races/${deleteTarget.id}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      setDeleteTarget(null)
      fetchRaces()
    } catch (e: any) {
      alert(`Failed to delete: ${e.message}`)
    }
  }

  const startEdit = (race: Race) => {
    setEditing(race)
    setForm({
      name: race.name,
      display_name: race.display_name,
      description: race.description,
      stat_modifiers: race.stat_modifiers || '{}',
      skill_grants: race.skill_grants || '[]',
      is_playable: race.is_playable,
    })
    setShowForm(false)
  }

  const startCreate = () => {
    setShowForm(true)
    setEditing(null)
    setForm(emptyForm)
  }

  const cancelForm = () => {
    setShowForm(false)
    setEditing(null)
    setForm(emptyForm)
  }

  const filtered = races.filter(r =>
    r.name.toLowerCase().includes(search.toLowerCase()) ||
    r.display_name.toLowerCase().includes(search.toLowerCase())
  )

  const formatJSON = (val: string) => {
    try { return JSON.stringify(JSON.parse(val), null, 2) }
    catch { return val }
  }

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Races</h2>
        <button className="btn-primary" onClick={startCreate}>+ New Race</button>
      </div>

      {error && <div className="error-banner">{error}</div>}

      <div className="toolbar-row" style={{ marginBottom: '1rem', display: 'flex', gap: '0.75rem', alignItems: 'center' }}>
        <input
          type="text"
          placeholder="Search by name or display name..."
          value={search}
          onChange={e => setSearch(e.target.value)}
          style={{ flex: 1, padding: '0.5rem 0.75rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px' }}
        />
        <button className="btn-secondary" onClick={fetchRaces}>Refresh</button>
      </div>

      {loading ? (
        <div className="loading">Loading races...</div>
      ) : (
        <table className="data-table">
          <thead>
            <tr>
              <th>Name</th>
              <th>Display Name</th>
              <th>Description</th>
              <th>Stats</th>
              <th>Skills</th>
              <th>Playable</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {filtered.length === 0 ? (
              <tr>
                <td colSpan={7} style={{ textAlign: 'center', color: 'var(--text-muted)', padding: '2rem' }}>
                  {races.length === 0 ? 'No races found. Create one below.' : 'No races match your search.'}
                </td>
              </tr>
            ) : (
              filtered.map(race => (
                <tr key={race.id}>
                  <td><code style={{ color: 'var(--primary)', fontSize: '0.85rem' }}>{race.name}</code></td>
                  <td style={{ fontWeight: 500 }}>{race.display_name}</td>
                  <td style={{ color: 'var(--text-secondary)', fontSize: '0.8rem', maxWidth: '200px', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {race.description || <span style={{ color: 'var(--text-muted)' }}>—</span>}
                  </td>
                  <td>
                    <code style={{ fontSize: '0.7rem', background: 'var(--surface-muted)', padding: '0.15rem 0.3rem', borderRadius: '3px' }}>
                      {race.stat_modifiers ? (
                        <span title={formatJSON(race.stat_modifiers)}>
                          {race.stat_modifiers.length > 20 ? race.stat_modifiers.slice(0, 20) + '…' : race.stat_modifiers}
                        </span>
                      ) : <span style={{ color: 'var(--text-muted)' }}>—</span>}
                    </code>
                  </td>
                  <td>
                    {race.skill_grants && race.skill_grants !== '[]' ? (
                      <code style={{ fontSize: '0.7rem', background: 'var(--surface-muted)', padding: '0.15rem 0.3rem', borderRadius: '3px' }}>
                        {race.skill_grants.length > 20 ? race.skill_grants.slice(0, 20) + '…' : race.skill_grants}
                      </code>
                    ) : <span style={{ color: 'var(--text-muted)' }}>—</span>}
                  </td>
                  <td>
                    {race.is_playable
                      ? <span style={{ color: 'var(--success)', fontWeight: 600 }}>Yes</span>
                      : <span style={{ color: 'var(--text-muted)' }}>No</span>}
                  </td>
                  <td>
                    <div style={{ display: 'flex', gap: '0.5rem' }}>
                      <button className="btn-small" onClick={() => startEdit(race)}>Edit</button>
                      <button className="btn-small btn-danger" onClick={() => setDeleteTarget(race)}>Delete</button>
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
        <div className="modal-overlay" onClick={cancelForm}>
          <div className="modal-content" onClick={e => e.stopPropagation()} style={{ maxWidth: '560px' }}>
            <h3>{editing ? `Edit: ${editing.display_name}` : 'New Race'}</h3>
            <form onSubmit={editing ? handleUpdate : handleCreate}>
              <div className="form-group">
                <label>Internal Name</label>
                {editing ? (
                  <code style={{ display: 'block', padding: '0.5rem', background: 'var(--surface-muted)', borderRadius: '4px', color: 'var(--primary)' }}>{editing.name}</code>
                ) : (
                  <input
                    type="text"
                    value={form.name}
                    onChange={e => setForm(f => ({ ...f, name: e.target.value }))}
                    placeholder="e.g. human, mutant"
                    required
                    pattern="[a-z_]+"
                    title="Lowercase letters and underscores only"
                    style={{ width: '100%', padding: '0.5rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px' }}
                  />
                )}
              </div>
              <div className="form-group">
                <label>Display Name</label>
                <input
                  type="text"
                  value={form.display_name}
                  onChange={e => setForm(f => ({ ...f, display_name: e.target.value }))}
                  placeholder="e.g. Human, Mutant"
                  required
                  style={{ width: '100%', padding: '0.5rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px' }}
                />
              </div>
              <div className="form-group">
                <label>Description</label>
                <textarea
                  value={form.description}
                  onChange={e => setForm(f => ({ ...f, description: e.target.value }))}
                  rows={3}
                  placeholder="Flavor text shown in character creation"
                  style={{ width: '100%', padding: '0.5rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px', resize: 'vertical' }}
                />
              </div>
              <div className="form-group">
                <label>Stat Modifiers (JSON)</label>
                <textarea
                  value={form.stat_modifiers}
                  onChange={e => setForm(f => ({ ...f, stat_modifiers: e.target.value }))}
                  rows={4}
                  placeholder='{"strength": 0, "dexterity": 0, ...}'
                  style={{ width: '100%', padding: '0.5rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px', fontFamily: 'monospace', fontSize: '0.85rem', resize: 'vertical' }}
                />
              </div>
              <div className="form-group">
                <label>Skill Grants (JSON array)</label>
                <textarea
                  value={form.skill_grants}
                  onChange={e => setForm(f => ({ ...f, skill_grants: e.target.value }))}
                  rows={2}
                  placeholder='["skill_id"]'
                  style={{ width: '100%', padding: '0.5rem', background: 'var(--surface-muted)', border: '2px solid var(--border)', color: 'var(--text)', borderRadius: '4px', fontFamily: 'monospace', fontSize: '0.85rem', resize: 'vertical' }}
                />
              </div>
              <div className="form-group">
                <label style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', cursor: 'pointer' }}>
                  <input
                    type="checkbox"
                    checked={form.is_playable}
                    onChange={e => setForm(f => ({ ...f, is_playable: e.target.checked }))}
                  />
                  Playable (available for player character creation)
                </label>
              </div>
              <div style={{ display: 'flex', gap: '0.75rem', justifyContent: 'flex-end', marginTop: '1rem' }}>
                <button type="button" className="btn-secondary" onClick={cancelForm}>Cancel</button>
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
            <h3>Delete Race?</h3>
            <p>Are you sure you want to delete <strong>{deleteTarget.display_name}</strong> (<code>{deleteTarget.name}</code>)? This cannot be undone.</p>
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
