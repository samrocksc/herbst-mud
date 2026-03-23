import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useEffect, useState, useCallback } from 'react'

export const Route = createFileRoute('/npcs')({
  component: NPCManager,
})

interface Room {
  id: number
  name: string
}

interface NPC {
  id: number
  name: string
  class: string
  race: string
  level: number
  currentRoomId: number
  isNPC: boolean
  hitpoints: number
  max_hitpoints: number
  stamina?: number
  max_stamina?: number
  mana?: number
  max_mana?: number
}

interface NPCForm {
  name: string
  race: string
  class: string
  level: number
  currentRoomId: number
  hitpoints: number
  max_hitpoints: number
  roaming: boolean
  roamingTime: number
}

const RACES = ['human', 'elf', 'dwarf', 'halfling', 'half-dog', 'mutant', 'robot']
const CLASSES = ['adventurer', 'warrior', 'mage', 'rogue', 'healer', 'merchant']

function NPCManager() {
  const navigate = useNavigate()
  const [npcs, setNpcs] = useState<NPC[]>([])
  const [rooms, setRooms] = useState<Room[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedNPC, setSelectedNPC] = useState<NPC | null>(null)
  const [editingNPC, setEditingNPC] = useState<NPC | null>(null)
  const [saving, setSaving] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [form, setForm] = useState<NPCForm>({
    name: '',
    race: 'human',
    class: 'adventurer',
    level: 1,
    currentRoomId: 0,
    hitpoints: 100,
    max_hitpoints: 100,
    roaming: false,
    roamingTime: 0
  })

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate({ to: '/login' })
      return
    }

    Promise.all([
      fetch('http://localhost:8080/npcs').then(res => res.json()),
      fetch('http://localhost:8080/rooms').then(res => res.json())
    ])
      .then(([npcsData, roomsData]) => {
        setNpcs(npcsData.npcs || [])
        setRooms(roomsData)
        setLoading(false)
      })
      .catch(err => {
        console.error('Failed to load data:', err)
        setLoading(false)
      })
  }, [navigate])

  const refreshNPCs = useCallback(async () => {
    const npcsResponse = await fetch('http://localhost:8080/npcs')
    const npcsData = await npcsResponse.json()
    setNpcs(npcsData.npcs || [])
  }, [])

  const handleCreateNPC = useCallback(async () => {
    if (!form.name || !form.currentRoomId) {
      alert('Please fill in all required fields')
      return
    }

    setSaving(true)
    try {
      const token = localStorage.getItem('token')

      const response = await fetch('http://localhost:8080/characters', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          name: form.name,
          isNPC: true,
          currentRoomId: form.currentRoomId,
          startingRoomId: form.currentRoomId,
          race: form.race,
          class: form.class,
          level: form.level,
          hitpoints: form.hitpoints,
          max_hitpoints: form.max_hitpoints
        })
      })

      if (!response.ok) {
        throw new Error('Failed to create NPC')
      }

      await refreshNPCs()
      setForm({
        name: '',
        race: 'human',
        class: 'adventurer',
        level: 1,
        currentRoomId: 0,
        hitpoints: 100,
        max_hitpoints: 100,
        roaming: false,
        roamingTime: 0
      })
      setShowCreateForm(false)
    } catch (err) {
      console.error('Create NPC error:', err)
      alert('Failed to create NPC')
    } finally {
      setSaving(false)
    }
  }, [form, refreshNPCs])

  const handleUpdateNPC = useCallback(async () => {
    if (!editingNPC) return
    setSaving(true)
    try {
      const token = localStorage.getItem('token')

      const response = await fetch(`http://localhost:8080/characters/${editingNPC.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          name: form.name,
          race: form.race,
          class: form.class,
          level: form.level,
          currentRoomId: form.currentRoomId,
          hitpoints: form.hitpoints,
          max_hitpoints: form.max_hitpoints
        })
      })

      if (!response.ok) {
        throw new Error('Failed to update NPC')
      }

      await refreshNPCs()
      setEditingNPC(null)
      setSelectedNPC(null)
    } catch (err) {
      console.error('Update NPC error:', err)
      alert('Failed to update NPC')
    } finally {
      setSaving(false)
    }
  }, [editingNPC, form, refreshNPCs])

  const handleDeleteNPC = useCallback(async (npcId: number) => {
    try {
      const token = localStorage.getItem('token')

      const response = await fetch(`http://localhost:8080/characters/${npcId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (!response.ok) {
        throw new Error('Failed to delete NPC')
      }

      await refreshNPCs()
      setSelectedNPC(null)
      setConfirmDelete(null)
    } catch (err) {
      console.error('Delete NPC error:', err)
      alert('Failed to delete NPC')
    }
  }, [refreshNPCs])

  const startEditing = (npc: NPC) => {
    setEditingNPC(npc)
    setForm({
      name: npc.name,
      race: npc.race || 'human',
      class: npc.class || 'adventurer',
      level: npc.level || 1,
      currentRoomId: npc.currentRoomId,
      hitpoints: npc.hitpoints || 100,
      max_hitpoints: npc.max_hitpoints || 100,
      roaming: false,
      roamingTime: 0
    })
    setShowCreateForm(false)
  }

  const getRoomName = (roomId: number): string => {
    const room = rooms.find(r => r.id === roomId)
    return room ? room.name : `Room ${roomId}`
  }

  const getHPStatus = (npc: NPC): { text: string; color: string } => {
    if (npc.max_hitpoints === 0) {
      return { text: '∞ Endless', color: '#4a7c4e' }
    }
    return { text: `${npc.hitpoints}/${npc.max_hitpoints}`, color: npc.hitpoints < npc.max_hitpoints * 0.3 ? '#8b4444' : '#e8dcc4' }
  }

  if (loading) {
    return <div className="p-8 text-text">Loading NPCs...</div>
  }

  return (
    <div className="flex h-screen bg-surface">
      {/* Left Sidebar */}
      <div className="w-[280px] bg-surface-muted border-r border-border flex flex-col">
        <div className="p-4 border-b border-border">
          <Link
            to="/dashboard"
            className="block text-primary no-underline p-2 rounded bg-surface-dark text-center mb-2 hover:bg-surface-darker"
          >
            ← Dashboard
          </Link>
          <Link
            to="/map"
            className="block text-text-muted no-underline p-2 rounded bg-surface-dark text-center mb-2 hover:bg-surface-darker"
          >
            Map Builder
          </Link>
          <Link
            to="/items"
            className="block text-text-muted no-underline p-2 rounded bg-surface-dark text-center hover:bg-surface-darker"
          >
            Item Manager
          </Link>
        </div>

        <div className="p-3 border-b border-border">
          <h2 className="m-0 text-text text-lg">NPC Manager</h2>
          <p className="text-text-muted text-xs mt-1 mb-0">{npcs.length} NPCs</p>
        </div>

        {/* NPC List */}
        <div className="flex-1 overflow-y-auto p-3">
          <div className="flex flex-col gap-1">
            {npcs.map(npc => {
              const hpStatus = getHPStatus(npc)
              return (
                <div
                  key={npc.id}
                  onClick={() => { setSelectedNPC(npc); setEditingNPC(null); setShowCreateForm(false); }}
                  className={`p-2 cursor-pointer rounded text-xs ${selectedNPC?.id === npc.id ? 'text-text bg-surface-dark' : 'text-text'}`}
                >
                  <div className="font-bold">{npc.name}</div>
                  <div className="text-text-muted">
                    {npc.race} {npc.class} lv.{npc.level}
                  </div>
                  <div className="text-[10px] flex justify-between">
                    <span className="text-text-muted">{getRoomName(npc.currentRoomId)}</span>
                    <span style={{ color: hpStatus.color }}>HP: {hpStatus.text}</span>
                  </div>
                </div>
              )
            })}
          </div>
        </div>

        {/* Create NPC Button */}
        <div className="p-3 border-t border-border">
          <button
            onClick={() => { setShowCreateForm(true); setSelectedNPC(null); setEditingNPC(null); }}
            className="w-full p-2 bg-primary border-none rounded text-white cursor-pointer hover:bg-primary-hover"
          >
            + Create NPC
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {showCreateForm ? (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Create New NPC</h2>

            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Name *</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="NPC name"
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Race</label>
                  <select
                    value={form.race}
                    onChange={(e) => setForm({ ...form, race: e.target.value })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {RACES.map(race => (
                      <option key={race} value={race}>{race}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Class</label>
                  <select
                    value={form.class}
                    onChange={(e) => setForm({ ...form, class: e.target.value })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {CLASSES.map(cls => (
                      <option key={cls} value={cls}>{cls}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Level</label>
                  <input
                    type="number"
                    value={form.level}
                    onChange={(e) => setForm({ ...form, level: parseInt(e.target.value) || 1 })}
                    min={1}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Room *</label>
                  <select
                    value={form.currentRoomId}
                    onChange={(e) => setForm({ ...form, currentRoomId: parseInt(e.target.value) })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    <option value={0}>Select room...</option>
                    {rooms.map(room => (
                      <option key={room.id} value={room.id}>{room.name} (ID: {room.id})</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Hitpoints (0 = endless)</label>
                  <input
                    type="number"
                    value={form.hitpoints}
                    onChange={(e) => setForm({ ...form, hitpoints: parseInt(e.target.value) || 0 })}
                    min={0}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Max Hitpoints</label>
                  <input
                    type="number"
                    value={form.max_hitpoints}
                    onChange={(e) => setForm({ ...form, max_hitpoints: parseInt(e.target.value) || 0 })}
                    min={0}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
              </div>

              {form.hitpoints === 0 && (
                <div className="mb-4 p-2 bg-primary/20 border border-primary rounded text-success text-xs">
                  ⚡ Endless HP: This NPC will never die from combat. Great for training dummies!
                </div>
              )}

              <div className="flex gap-2">
                <button
                  onClick={handleCreateNPC}
                  disabled={saving}
                  className="flex-1 p-2 bg-primary border-none rounded text-white cursor-pointer disabled:opacity-70"
                >
                  {saving ? 'Creating...' : 'Create NPC'}
                </button>
                <button
                  onClick={() => setShowCreateForm(false)}
                  className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        ) : editingNPC ? (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Edit NPC</h2>

            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="mb-4">
                <label className="text-text-muted text-xs block mb-1">Name *</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Race</label>
                  <select
                    value={form.race}
                    onChange={(e) => setForm({ ...form, race: e.target.value })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {RACES.map(race => (
                      <option key={race} value={race}>{race}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Class</label>
                  <select
                    value={form.class}
                    onChange={(e) => setForm({ ...form, class: e.target.value })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {CLASSES.map(cls => (
                      <option key={cls} value={cls}>{cls}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Level</label>
                  <input
                    type="number"
                    value={form.level}
                    onChange={(e) => setForm({ ...form, level: parseInt(e.target.value) || 1 })}
                    min={1}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Room</label>
                  <select
                    value={form.currentRoomId}
                    onChange={(e) => setForm({ ...form, currentRoomId: parseInt(e.target.value) })}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  >
                    {rooms.map(room => (
                      <option key={room.id} value={room.id}>{room.name} (ID: {room.id})</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Hitpoints (0 = endless)</label>
                  <input
                    type="number"
                    value={form.hitpoints}
                    onChange={(e) => setForm({ ...form, hitpoints: parseInt(e.target.value) || 0 })}
                    min={0}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Max Hitpoints</label>
                  <input
                    type="number"
                    value={form.max_hitpoints}
                    onChange={(e) => setForm({ ...form, max_hitpoints: parseInt(e.target.value) || 0 })}
                    min={0}
                    className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                  />
                </div>
              </div>

              {form.hitpoints === 0 && (
                <div className="mb-4 p-2 bg-primary/20 border border-primary rounded text-success text-xs">
                  ⚡ Endless HP: This NPC will never die from combat.
                </div>
              )}

              <div className="flex gap-2">
                <button
                  onClick={handleUpdateNPC}
                  disabled={saving}
                  className="flex-1 p-2 bg-primary border-none rounded text-white cursor-pointer disabled:opacity-70"
                >
                  {saving ? 'Saving...' : 'Save Changes'}
                </button>
                <button
                  onClick={() => setEditingNPC(null)}
                  className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        ) : selectedNPC ? (
          <div className="max-w-[600px] mx-auto">
            <div className="flex justify-between items-center mb-4">
              <h2 className="m-0 text-text">{selectedNPC.name}</h2>
              <button onClick={() => setSelectedNPC(null)} className="bg-transparent border-none text-text-muted cursor-pointer text-xl">×</button>
            </div>

            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="text-text-muted text-[10px] mb-2">ID: #{selectedNPC.id}</div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Race</label>
                  <div className="text-text">{selectedNPC.race}</div>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Class</label>
                  <div className="text-text">{selectedNPC.class}</div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Level</label>
                  <div className="text-text">{selectedNPC.level}</div>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Room</label>
                  <div className="text-text">{getRoomName(selectedNPC.currentRoomId)}</div>
                  <div className="text-text-muted text-[10px]">Room ID: {selectedNPC.currentRoomId}</div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-text-muted text-xs block mb-1">Hitpoints</label>
                  <div className="text-text">
                    {selectedNPC.max_hitpoints === 0 ? (
                      <span className="text-success">∞ Endless</span>
                    ) : (
                      `${selectedNPC.hitpoints} / ${selectedNPC.max_hitpoints}`
                    )}
                  </div>
                </div>
                <div>
                  <label className="text-text-muted text-xs block mb-1">Status</label>
                  <div className="text-text">
                    {selectedNPC.max_hitpoints === 0 ? 'Training Dummy' : 'Combatant'}
                  </div>
                </div>
              </div>

              <div className="flex gap-2">
                <button
                  onClick={() => startEditing(selectedNPC)}
                  className="flex-1 p-2 bg-accent border-none rounded text-white cursor-pointer hover:bg-accent-hover"
                >
                  Edit NPC
                </button>
                <button
                  onClick={() => {
                    if (confirmDelete === selectedNPC.id) {
                      handleDeleteNPC(selectedNPC.id)
                    } else {
                      setConfirmDelete(selectedNPC.id)
                    }
                  }}
                  className={`flex-1 p-2 border-none rounded text-white cursor-pointer ${
                    confirmDelete === selectedNPC.id
                      ? 'bg-warning hover:bg-warning/80'
                      : 'bg-danger hover:bg-danger-hover'
                  }`}
                >
                  {confirmDelete === selectedNPC.id ? 'Confirm Delete?' : 'Delete NPC'}
                </button>
              </div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-full text-text-muted">
            <p className="mb-2">Select an NPC from the list or create a new one</p>
            <p className="text-xs text-center">
              💡 Set HP to 0 for endless training dummies<br/>
              NPCs can be placed in any room
            </p>
          </div>
        )}
      </div>
    </div>
  )
}