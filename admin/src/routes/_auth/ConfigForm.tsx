import { useState } from 'react'
import { Button } from '../../components/Button'
import { FormField, TextareaField, FormError } from '../../components/FormFields'
import { apiPost, apiPut } from '../../utils/apiFetch'
import { showToast } from '../../components/Toast'

export type GameConfig = Readonly<{ id: number; key: string; value: string }>

const KEY_LABELS: Readonly<Record<string, string>> = {
  xp_thresholds: 'XP Thresholds', death_penalty_percent: 'Death Penalty (%)', starting_room_id: 'Starting Room ID',
  max_level: 'Max Level', campaign_type: 'Campaign Type', corpse_rot_minutes: 'Corpse Rot Minutes',
  xp_per_kill: 'XP Per Kill', respawn_room_id: 'Respawn Room ID', max_inventory_size: 'Max Inventory Size',
  pvp_enabled: 'PvP Enabled', death_penalty_type: 'Death Penalty Type', starting_hp: 'Starting HP',
  starting_mana: 'Starting Mana', starting_gold: 'Starting Gold', regen_tick_seconds: 'Regen Tick (seconds)',
}

export function humanizeKey(raw: string): string {
  if (KEY_LABELS[raw]) return KEY_LABELS[raw]
  return raw.split('_').map(s => s.charAt(0).toUpperCase() + s.slice(1)).join(' ')
}

export function tryParseJSON(value: string): unknown | null {
  try { const p = JSON.parse(value); return typeof p === 'object' && p !== null ? p : null }
  catch { return null }
}

export const PRESETS = [
  { label: 'XP Thresholds', key: 'xp_thresholds', value: '{"1":100,"2":300,"3":600,"4":1000,"5":1500}' },
  { label: 'Death Penalty %', key: 'death_penalty_percent', value: '10' },
  { label: 'Corpse Rot Minutes', key: 'corpse_rot_minutes', value: '5' },
  { label: 'XP Per Kill', key: 'xp_per_kill', value: '50' },
  { label: 'Max Level', key: 'max_level', value: '100' },
  { label: 'Starting Room ID', key: 'starting_room_id', value: '1' },
]

function CollapsibleJSONPreview({ value }: { value: string }) {
  const parsed = tryParseJSON(value)
  const [expanded, setExpanded] = useState(false)
  if (parsed === null) return null
  const formatted = JSON.stringify(parsed, null, 2)
  return (
    <div className="mb-3">
      <button type="button" className="text-xs text-primary hover:underline cursor-pointer flex items-center gap-1 mb-1"
        onClick={() => setExpanded(e => !e)}>
        <span className={`inline-block transition-transform ${expanded ? 'rotate-90' : ''}`}>&#9654;</span>
        {expanded ? 'Collapse' : 'Expand'} JSON preview
      </button>
      {expanded && <pre className="bg-surface-muted border-2 border-border rounded p-3 text-xs font-mono whitespace-pre-wrap overflow-auto max-h-64">{formatted}</pre>}
    </div>
  )
}

export function ConfigForm({ editing, onDone }: Readonly<{ editing: GameConfig | null; onDone: () => void }>) {
  const [form, setForm] = useState(editing ? { key: editing.key, value: editing.value } : { key: '', value: '' })
  const [saving, setSaving] = useState(false)
  const [formError, setFormError] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault(); setSaving(true); setFormError('')
    try {
      if (editing) { await apiPut(`/api/game-configs/${editing.key}`, { value: form.value }); showToast('Config updated.', 'success') }
      else { await apiPost('/api/game-configs', form); showToast('Config created.', 'success') }
      onDone()
    } catch (e: unknown) { setFormError(e instanceof Error ? e.message : 'Unknown error') }
    finally { setSaving(false) }
  }

  return (
    <div className="modal-overlay" onClick={onDone}>
      <div className="modal-content max-w-2xl" onClick={e => e.stopPropagation()}>
        <h3>{editing ? `Edit: ${humanizeKey(editing.key)}` : 'New Game Config'}</h3>
        {formError && <FormError message={formError} />}
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label>Key</label>
            {editing ? <code className="block p-2 bg-surface-muted rounded">{editing.key}</code>
              : <FormField label="Key" value={form.key} onChange={v => setForm(f => ({ ...f, key: v }))} placeholder="e.g. xp_thresholds" required />}
          </div>
          <div className="form-group">
            <label>Value (JSON or plain)</label>
            <CollapsibleJSONPreview value={form.value} />
            <TextareaField label="" value={form.value} onChange={v => setForm(f => ({ ...f, value: v }))} rows={6} placeholder='{"key": "value"} or plain text' required />
          </div>
          {!editing && (
            <div className="form-group">
              <label>Presets</label>
              <div className="flex flex-wrap gap-2">
                {PRESETS.map(p => <Button type="button" key={p.key} variant="ghost" size="sm" onClick={() => setForm({ key: p.key, value: p.value })}>{p.label}</Button>)}
              </div>
            </div>
          )}
          <div className="flex gap-3 justify-end mt-4">
            <Button type="button" variant="secondary" onClick={onDone}>Cancel</Button>
            <Button type="submit" variant="primary" disabled={saving}>{saving ? 'Saving...' : (editing ? 'Update' : 'Create')}</Button>
          </div>
        </form>
      </div>
    </div>
  )
}