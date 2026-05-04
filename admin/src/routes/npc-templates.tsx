import { createFileRoute, Link } from '@tanstack/react-router'
import { useEffect, useState, useCallback } from 'react'
import { Button } from '../components/Button'
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

  const load = useCallback(async () => {
    const token = localStorage.getItem('token')
    const res = await fetch(`${window.location.origin}/api/npc-templates`, {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (!res.ok) throw new Error('failed to load')
    const data = await res.json()
    setTemplates(data)
    setLoading(false)
  }, [])

  useEffect(() => {
    load()
  }, [load])

  async function handleSave(id: string) {
    setSaving(true)
    try {
      const token = localStorage.getItem('token')
      const rooms = roomsInput
        .split(',')
        .map((s) => s.trim().replace(/^r/i, '')) // remove optional "r" prefix
        .filter((s) => s !== '')
      const cooldown = parseInt(cooldownInput, 10) || 0

      const res = await fetch(`${window.location.origin}/api/npc-templates/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          xp_value: templates.find((t) => t.id === id)?.xp_value ?? 0,
          respawn_rooms: rooms,
          respawn_cooldown: cooldown,
        }),
      })
      if (!res.ok) throw new Error('save failed')
      await load()
      setEditingID(null)
    } catch (e) {
      alert('Save failed: ' + (e as Error).message)
    } finally {
      setSaving(false)
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
      <PageHeader title="NPC Templates" backTo="/dashboard" />

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
    </div>
  )
}
