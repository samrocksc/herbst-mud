import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useEffect, useState, useCallback } from 'react'
import { Button } from '../components/Button'
import { DashboardIcon } from '../components/icons/DashboardIcon'
import { MapIcon } from '../components/icons/MapIcon'
import { ItemsIcon } from '../components/icons/ItemsIcon'
import { MenuIcon } from '../components/icons/MenuIcon'

export const Route = createFileRoute('/npcs')({
  component: NPCManager,
})

type Room = Readonly<{
  id: number
  name: string
}>

type NPC = Readonly<{
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
}>

type NPCForm = Readonly<{
  name: string
  race: string
  class: string
  level: number
  currentRoomId: number
  hitpoints: number
  max_hitpoints: number
  roaming: boolean
  roamingTime: number
}>

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
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [form, setForm] = useState<NPCForm>({
    name: '',
    race: 'human',
    class: 'adventurer',
    level: 1,
    currentRoomId: 0,
    hitpoints: 100,
    max_hitpoints: 100,
    roaming: false,
    roamingTime: 0,
  })

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate({ to: '/login' })
      return
    }

    Promise.all([
      fetch(`${window.location.origin}/npcs`).then((res) => res.json()),
      fetch(`${window.location.origin}/rooms`).then((res) => res.json()),
    ])
      .then(([npcsData, roomsData]) => {
        setNpcs(npcsData.npcs || [])
        setRooms(roomsData)
        setLoading(false)
      })
      .catch((err) => {
        console.error('Failed to load data:', err)
        setLoading(false)
      })
  }, [navigate])

  const refreshNPCs = useCallback(async () => {
    const npcsResponse = await fetch(`${window.location.origin}/npcs`)
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
      const response = await fetch(`${window.location.origin}/characters`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
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
          max_hitpoints: form.max_hitpoints,
        }),
      })

      if (!response.ok) throw new Error('Failed to create NPC')

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
        roamingTime: 0,
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
      const response = await fetch(`${window.location.origin}/characters/${editingNPC.id}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          name: form.name,
          race: form.race,
          class: form.class,
          level: form.level,
          currentRoomId: form.currentRoomId,
          hitpoints: form.hitpoints,
          max_hitpoints: form.max_hitpoints,
        }),
      })

      if (!response.ok) throw new Error('Failed to update NPC')

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

  const handleDeleteNPC = useCallback(
    async (npcId: number) => {
      try {
        const token = localStorage.getItem('token')
        const response = await fetch(`${window.location.origin}/characters/${npcId}`, {
          method: 'DELETE',
          headers: { Authorization: `Bearer ${token}` },
        })

        if (!response.ok) throw new Error('Failed to delete NPC')

        await refreshNPCs()
        setSelectedNPC(null)
        setConfirmDelete(null)
      } catch (err) {
        console.error('Delete NPC error:', err)
        alert('Failed to delete NPC')
      }
    },
    [refreshNPCs]
  )

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
      roamingTime: 0,
    })
    setShowCreateForm(false)
  }


  const getRoomName = (roomId: number): string => {
    const room = rooms.find((r) => r.id === roomId)
    return room ? room.name : `Room ${roomId}`
  }

  const getHPStatus = (npc: NPC) => {
    if (npc.max_hitpoints === 0)
      return { text: '∞ Endless', cls: 'text-success' }
    return {
      text: `${npc.hitpoints}/${npc.max_hitpoints}`,
      cls:
        npc.hitpoints < npc.max_hitpoints * 0.3
          ? 'text-danger'
          : 'text-text-muted',
    }
  }

  if (loading) {
    return (
      <div className="flex h-screen bg-surface">
        <NPcsSidebar
          npcs={npcs}
          rooms={rooms}
          selectedNPC={selectedNPC}
          onSelectNPC={(npc) => {
            setSelectedNPC(npc)
            setEditingNPC(null)
            setShowCreateForm(false)
          }}
          onCreateClick={() => {
            setShowCreateForm(true)
            setSelectedNPC(null)
            setEditingNPC(null)
          }}
        />
        <main className="flex-1 flex items-center justify-center text-text-muted">
          Loading NPCs...
        </main>
      </div>
    )
  }

  return (
    <div className="flex h-screen bg-surface">
      {/* Mobile hamburger */}
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setSidebarOpen(true)}
        aria-label="Open NPC list"
        className="fixed top-3 left-3 z-50 p-2 bg-surface border border-border text-text-muted hover:bg-surface-muted hover:text-text lg:hidden"
      >
        <MenuIcon stroke="currentColor" />
      </Button>

      {/* Sidebar overlay (mobile) / always-visible (desktop) */}
      <div className={['lg:block lg:relative lg:inset-auto lg:z-auto', sidebarOpen ? 'block' : 'hidden'].join(' ')}>
        <div className="fixed inset-y-0 left-0 z-40 lg:static">
          <NPcsSidebar
            npcs={npcs}
            rooms={rooms}
            selectedNPC={selectedNPC}
            onSelectNPC={(npc) => {
              setSelectedNPC(npc)
              setEditingNPC(null)
              setShowCreateForm(false)
              setSidebarOpen(false)
            }}
            onCreateClick={() => {
              setShowCreateForm(true)
              setSelectedNPC(null)
              setEditingNPC(null)
              setSidebarOpen(false)
            }}
          />
        </div>
      </div>

      {/* Mobile backdrop */}
      {sidebarOpen && (
        <div
          className="fixed inset-0 bg-black/30 z-30 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <main className="flex-1 overflow-y-auto p-6">
        {showCreateForm ? (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Create New NPC</h2>
            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <FormFields form={form} setForm={setForm} rooms={rooms} />
              {form.hitpoints === 0 && (
                <div className="mb-4 p-2 bg-primary/20 border border-primary rounded text-success text-xs">
                  ⚡ Endless HP: This NPC will never die from combat. Great for
                  training dummies!
                </div>
              )}
              <div className="flex gap-2">
                <Button
                  variant="primary"
                  size="md"
                  fullWidth
                  onClick={handleCreateNPC}
                  disabled={saving}
                >
                  {saving ? 'Creating...' : 'Create NPC'}
                </Button>
                <Button
                  variant="secondary"
                  size="md"
                  fullWidth
                  onClick={() => setShowCreateForm(false)}
                >
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        ) : editingNPC ? (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-text">Edit NPC</h2>
            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <FormFields form={form} setForm={setForm} rooms={rooms} />
              <div className="flex gap-2 mt-4">
                <Button
                  variant="primary"
                  size="md"
                  fullWidth
                  onClick={handleUpdateNPC}
                  disabled={saving}
                >
                  {saving ? 'Saving...' : 'Save Changes'}
                </Button>
                <Button
                  variant="secondary"
                  size="md"
                  fullWidth
                  onClick={() => setEditingNPC(null)}
                >
                  Cancel
                </Button>
              </div>
            </div>
          </div>
        ) : selectedNPC ? (
          <div className="max-w-[600px] mx-auto">
            <div className="flex justify-between items-center mb-4">
              <h2 className="m-0 text-text">{selectedNPC.name}</h2>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setSelectedNPC(null)}
                aria-label="Close"
              >
                ×
              </Button>
            </div>

            <div className="bg-surface-muted rounded-lg p-4 border border-border">
              <div className="text-text-muted text-[10px] mb-2">
                ID: #{selectedNPC.id}
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <div className="text-text-muted text-xs mb-1">Race</div>
                  <div className="text-text">{selectedNPC.race}</div>
                </div>
                <div>
                  <div className="text-text-muted text-xs mb-1">Class</div>
                  <div className="text-text">{selectedNPC.class}</div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <div className="text-text-muted text-xs mb-1">Level</div>
                  <div className="text-text">{selectedNPC.level}</div>
                </div>
                <div>
                  <div className="text-text-muted text-xs mb-1">Room</div>
                  <div className="text-text">{getRoomName(selectedNPC.currentRoomId)}</div>
                  <div className="text-text-muted text-[10px]">
                    Room ID: {selectedNPC.currentRoomId}
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <div className="text-text-muted text-xs mb-1">Hitpoints</div>
                  <div className={getHPStatus(selectedNPC).cls}>
                    {selectedNPC.max_hitpoints === 0
                      ? '∞ Endless'
                      : `${selectedNPC.hitpoints} / ${selectedNPC.max_hitpoints}`}
                  </div>
                </div>
                <div>
                  <div className="text-text-muted text-xs mb-1">Status</div>
                  <div className="text-text">
                    {selectedNPC.max_hitpoints === 0 ? 'Training Dummy' : 'Combatant'}
                  </div>
                </div>
              </div>

              <div className="flex gap-2">
                <Button
                  variant="accent"
                  size="md"
                  fullWidth
                  onClick={() => startEditing(selectedNPC)}
                >
                  Edit NPC
                </Button>
                <Button
                  variant={confirmDelete === selectedNPC.id ? 'secondary' : 'danger'}
                  size="md"
                  fullWidth
                  onClick={() => {
                    if (confirmDelete === selectedNPC.id) {
                      handleDeleteNPC(selectedNPC.id)
                    } else {
                      setConfirmDelete(selectedNPC.id)
                    }
                  }}
                >
                  {confirmDelete === selectedNPC.id ? 'Confirm Delete?' : 'Delete NPC'}
                </Button>
              </div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-full text-text-muted text-center">
            <p className="mb-2">Select an NPC from the list or create a new one</p>
            <p className="text-xs">
              💡 Set HP to 0 for endless training dummies
              <br />
              NPCs can be placed in any room
            </p>
          </div>
        )}
      </main>
    </div>
  )
}

// ─── Shared form fields (used by create + edit) ─────────────────────────────

function FormFields({
  form,
  setForm,
  rooms,
}: {
  form: NPCForm
  setForm: (f: NPCForm) => void
  rooms: Room[]
}) {
  return (
    <>
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
            {RACES.map((race) => (
              <option key={race} value={race}>
                {race}
              </option>
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
            {CLASSES.map((cls) => (
              <option key={cls} value={cls}>
                {cls}
              </option>
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
            onChange={(e) =>
              setForm({ ...form, currentRoomId: parseInt(e.target.value) })
            }
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
          >
            <option value={0}>Select room...</option>
            {rooms.map((room) => (
              <option key={room.id} value={room.id}>
                {room.name} (ID: {room.id})
              </option>
            ))}
          </select>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-4">
        <div>
          <label className="text-text-muted text-xs block mb-1">
            Hitpoints (0 = endless)
          </label>
          <input
            type="number"
            value={form.hitpoints}
            onChange={(e) =>
              setForm({ ...form, hitpoints: parseInt(e.target.value) || 0 })
            }
            min={0}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
        <div>
          <label className="text-text-muted text-xs block mb-1">Max Hitpoints</label>
          <input
            type="number"
            value={form.max_hitpoints}
            onChange={(e) =>
              setForm({ ...form, max_hitpoints: parseInt(e.target.value) || 0 })
            }
            min={0}
            className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
          />
        </div>
      </div>
    </>
  )
}

// ─── NPC sidebar ─────────────────────────────────────────────────────────────

type NPcsSidebarProps = {
  npcs: NPC[]
  rooms: Room[]
  selectedNPC: NPC | null
  onSelectNPC: (npc: NPC) => void
  onCreateClick: () => void
}

function NPcsSidebar({
  npcs,
  rooms,
  selectedNPC,
  onSelectNPC,
  onCreateClick,
}: NPcsSidebarProps) {
  const getRoomName = (roomId: number) => {
    const room = rooms.find((r) => r.id === roomId)
    return room ? room.name : `Room ${roomId}`
  }

  const getHPStatus = (npc: NPC) => {
    if (npc.max_hitpoints === 0)
      return { text: '∞ Endless', cls: 'text-success' }
    return {
      text: `${npc.hitpoints}/${npc.max_hitpoints}`,
      cls:
        npc.hitpoints < npc.max_hitpoints * 0.3
          ? 'text-danger'
          : 'text-text-muted',
    }
  }
  return (
    <div className="w-[220px] bg-surface-muted border-r border-border flex flex-col flex-shrink-0">
      {/* Dashboard + secondary nav */}
      <div className="p-3 border-b border-border flex flex-col gap-1">
        <Link
          to="/dashboard"
          activeProps={{
            className: 'bg-primary/10 text-primary border-l-4 border-primary font-semibold',
          }}
          inactiveProps={{
            className: 'text-text-muted hover:bg-surface-muted hover:text-text',
          }}
          className="flex items-center gap-3 px-3 py-2 rounded text-sm no-underline transition-colors"
        >
          <span className="flex-shrink-0">
            <DashboardIcon stroke="currentColor" />
          </span>
          <span className="whitespace-nowrap">Dashboard</span>
        </Link>
        <Link
          to="/map"
          activeProps={{
            className: 'bg-primary/10 text-primary border-l-4 border-primary font-semibold',
          }}
          inactiveProps={{
            className: 'text-text-muted hover:bg-surface-muted hover:text-text',
          }}
          className="flex items-center gap-3 px-3 py-2 rounded text-sm no-underline transition-colors"
        >
          <span className="flex-shrink-0">
            <MapIcon stroke="currentColor" />
          </span>
          <span className="whitespace-nowrap">Map Builder</span>
        </Link>
        <Link
          to="/items"
          activeProps={{
            className: 'bg-primary/10 text-primary border-l-4 border-primary font-semibold',
          }}
          inactiveProps={{
            className: 'text-text-muted hover:bg-surface-muted hover:text-text',
          }}
          className="flex items-center gap-3 px-3 py-2 rounded text-sm no-underline transition-colors"
        >
          <span className="flex-shrink-0">
            <ItemsIcon stroke="currentColor" />
          </span>
          <span className="whitespace-nowrap">Items</span>
        </Link>
      </div>

      {/* Create NPC */}
      <div className="p-3 border-b border-border">
        <Button variant="primary" size="md" fullWidth onClick={onCreateClick}>
          + Create NPC
        </Button>
      </div>

      {/* Title + count */}
      <div className="p-3 border-b border-border">
        <h2 className="m-0 text-text text-base font-semibold">NPC Manager</h2>
        <p className="text-text-muted text-xs mt-0.5 m-0">{npcs.length} NPCs</p>
      </div>

      {/* NPC list */}
      <div className="flex-1 overflow-y-auto p-3">
        <div className="flex flex-col gap-1">
          {npcs.map((npc) => {
            const hpStatus = getHPStatus(npc)
            return (
              <div
                key={npc.id}
                onClick={() => onSelectNPC(npc)}
                className={[
                  'p-2 rounded text-xs cursor-pointer transition-colors',
                  selectedNPC?.id === npc.id
                    ? 'bg-primary/10 text-text border border-primary/30'
                    : 'text-text-muted hover:bg-surface hover:text-text',
                ].join(' ')}
              >
                <div className="font-semibold truncate">{npc.name}</div>
                <div className="text-text-muted">
                  {npc.race} {npc.class} lv.{npc.level}
                </div>
                <div className="flex justify-between mt-0.5">
                  <span className="text-text-muted truncate">
                    {getRoomName(npc.currentRoomId)}
                  </span>
                  <span className={hpStatus.cls}>HP: {hpStatus.text}</span>
                </div>
              </div>
            )
          })}
        </div>
      </div>
    </div>
  )
}
