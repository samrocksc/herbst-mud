import { createFileRoute } from '@tanstack/react-router'
import { useEffect, useState, useCallback } from 'react'
import { apiGet, apiPost, apiPut } from '../utils/apiFetch'
import { Button } from '../components/Button'
import { Modal } from '../components/Modal'
import { PageHeader } from '../components/PageHeader'

type NPCTemplate = Readonly<{
  id: string
  name: string
  level: number
  xp_value: number
  respawn_rooms: string[]
  respawn_cooldown: number
}>

export const Route = createFileRoute('/npc-templates')({
  component: NPCTemplatePage,
})

function NPCTemplatePage() {
  const [templates, setTemplates] = useState<NPCTemplate[]>([])
  const [loading, setLoading] = useState(true)
  const [editingID, setEditingID] = useState<string | null>(null)
  const [saving, setSaving] = useState(false)
  const [roomsInput, setRoomsInput] = useState('')
  const [cooldownInput, setCooldownInput] = useState('60')

  // Create modal state
  const [showCreate, setShowCreate] = useState(false)
  const [createForm, setCreateForm] = useState({
    id: '',
    name: '',
    description: '',
    race: '',
    disposition: 'neutral',
    level: '1',
    xp_value: '0',
    greeting: '',
    respawn_cooldown: '60',
    respawn_rooms: '',
  })
  const [createError, setCreateError] = useState('')

  const load = useCallback(async () => {
    try {
      const data = await apiGet<NPCTemplate[]>(`${window.location.origin}/api/npc-templates`)
      setTemplates(data)
    } catch (e) {
      console.error('Failed to load NPC templates', e)
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => {
    load()
  }, [load])

  async function handleSave(id: string) {
    setSaving(true)
    try {
      const rooms = roomsInput
        .split(',')
        .map((s) => s.trim().replace(/^r/i, ''))
        .filter((s) => s !== '')
      const cooldown = parseInt(cooldownInput, 10) || 0
      const current = templates.find((t) => t.id === id)

      await apiPut(`${window.location.origin}/api/npc-templates/${id}`, {
        xp_value: current?.xp_value ?? 0,
        respawn_rooms: rooms,
        respawn_cooldown: cooldown,
      })
      await load()
      setEditingID(null)
    } catch (e) {
      alert('Save failed: ' + (e as Error).message)
    } finally {
      setSaving(false)
    }
  }

  async function handleCreate() {
    setCreateError('')
    try {
      const rooms = createForm.respawn_rooms
        .split(',')
        .map((s) => s.trim().replace(/^r/i, ''))
        .filter((s) => s !== '')

      await apiPost(`${window.location.origin}/api/npc-templates`, {
        id: createForm.id,
        name: createForm.name,
        description: createForm.description,
        race: createForm.race,
        disposition: createForm.disposition,
        level: parseInt(createForm.level, 10) || 1,
        xp_value: parseInt(createForm.xp_value, 10) || 0,
        greeting: createForm.greeting,
        respawn_cooldown: parseInt(createForm.respawn_cooldown, 10) || 60,
        respawn_rooms: rooms,
        skills: {},
        trades_with: [],
      })
      await load()
      setShowCreate(false)
      setCreateForm({
        id: '',
        name: '',
        description: '',
        race: '',
        disposition: 'neutral',
        level: '1',
        xp_value: '0',
        greeting: '',
        respawn_cooldown: '60',
        respawn_rooms: '',
      })
    } catch (e) {
      setCreateError((e as Error).message)
    }
  }

  function startEdit(t: NPCTemplate) {
    setEditingID(t.id)
    setRoomsInput(t.respawn_rooms?.join(', ') ?? '')
    setCooldownInput(String(t.respawn_cooldown ?? 60))
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-surface p-6">
        <PageHeader title="NPC Templates" backTo="/dashboard" />
        <p className="text-text-muted">Loading templates...</p>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-surface p-6">
      <PageHeader
        title="NPC Templates"
        backTo="/dashboard"
        actions={
          <Button variant="primary" onClick={() => setShowCreate(true)}>
            + New NPC
          </Button>
        }
      />

      <div className="space-y-4 max-w-[900px]">
        {templates.map((t) => (
          <div
            key={t.id}
            className="bg-surface-muted border border-border rounded p-4 space-y-3"
          >
            <div className="flex justify-between items-start">
              <div>
                <div className="text-text font-medium">{t.name}</div>
                <div className="text-text-muted text-xs">
                  Level {t.level} · XP {t.xp_value} · ID: {t.id}
                </div>
              </div>
              {editingID !== t.id ? (
                <Button variant="accent" size="sm" onClick={() => startEdit(t)}>
                  Edit
                </Button>
              ) : (
                <div className="flex gap-2">
                  <Button
                    variant="primary"
                    size="sm"
                    onClick={() => handleSave(t.id)}
                    disabled={saving}
                  >
                    {saving ? 'Saving...' : 'Save'}
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => setEditingID(null)}>
                    Cancel
                  </Button>
                </div>
              )}
            </div>

            {editingID === t.id && (
              <div className="space-y-3">
                <div>
                  <label className="text-text-muted text-xs block mb-1">
                    Respawn Rooms <span className="text-text-muted">(comma-separated room IDs)</span>
                  </label>
                  <input
                    type="text"
                    value={roomsInput}
                    onChange={(e) => setRoomsInput(e.target.value)}
                    placeholder="1, 2, 3"
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">
                    Respawn Cooldown <span className="text-text-muted">(seconds, 0 = no respawn)</span>
                  </label>
                  <input
                    type="number"
                    value={cooldownInput}
                    onChange={(e) => setCooldownInput(e.target.value)}
                    min={0}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
              </div>
            )}
          </div>
        ))}

        {templates.length === 0 && (
          <p className="text-text-muted">No NPC templates found.</p>
        )}
      </div>

      {/* Create Modal */}
      <Modal isOpen={showCreate} onClose={() => setShowCreate(false)} title="Create NPC Instance">
        <div className="space-y-3">
          {createError && (
            <div className="p-2 bg-red-100 text-red-800 rounded text-sm">{createError}</div>
          )}

          <div>
            <label className="text-text-muted text-xs block mb-1">ID *</label>
            <input
              type="text"
              value={createForm.id}
              onChange={(e) => setCreateForm((f) => ({ ...f, id: e.target.value }))}
              placeholder="e.g. goblin_guard_01"
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Name *</label>
            <input
              type="text"
              value={createForm.name}
              onChange={(e) => setCreateForm((f) => ({ ...f, name: e.target.value }))}
              placeholder="Display name"
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Description</label>
            <textarea
              value={createForm.description}
              onChange={(e) => setCreateForm((f) => ({ ...f, description: e.target.value }))}
              placeholder="Flavor text..."
              rows={2}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div className="flex gap-3">
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Race</label>
              <input
                type="text"
                value={createForm.race}
                onChange={(e) => setCreateForm((f) => ({ ...f, race: e.target.value }))}
                placeholder="e.g. goblin"
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Disposition</label>
              <select
                value={createForm.disposition}
                onChange={(e) => setCreateForm((f) => ({ ...f, disposition: e.target.value }))}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              >
                <option value="neutral">Neutral</option>
                <option value="friendly">Friendly</option>
                <option value="hostile">Hostile</option>
              </select>
            </div>
          </div>

          <div className="flex gap-3">
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Level</label>
              <input
                type="number"
                value={createForm.level}
                onChange={(e) => setCreateForm((f) => ({ ...f, level: e.target.value }))}
                min={1}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">XP Value</label>
              <input
                type="number"
                value={createForm.xp_value}
                onChange={(e) => setCreateForm((f) => ({ ...f, xp_value: e.target.value }))}
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Greeting</label>
            <textarea
              value={createForm.greeting}
              onChange={(e) => setCreateForm((f) => ({ ...f, greeting: e.target.value }))}
              placeholder="NPC greeting message..."
              rows={2}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div className="flex gap-3">
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Respawn Cooldown (s)</label>
              <input
                type="number"
                value={createForm.respawn_cooldown}
                onChange={(e) => setCreateForm((f) => ({ ...f, respawn_cooldown: e.target.value }))}
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Respawn Rooms</label>
              <input
                type="text"
                value={createForm.respawn_rooms}
                onChange={(e) => setCreateForm((f) => ({ ...f, respawn_rooms: e.target.value }))}
                placeholder="1, 2, 3"
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          <div className="flex gap-2 pt-2">
            <Button variant="primary" onClick={handleCreate}>
              Create
            </Button>
            <Button variant="secondary" onClick={() => setShowCreate(false)}>
              Cancel
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}
