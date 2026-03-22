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
}

interface NPCForm {
  name: string
  race: string
  class: string
  level: number
  currentRoomId: number
  roaming: boolean
  roamingTime: number
  roamingRooms: number[]
}

const RACES = ['human', 'elf', 'dwarf', 'halfling', 'half-dog', 'mutant', 'robot']
const CLASSES = ['adventurer', 'warrior', 'mage', 'rogue', 'healer', 'merchant']

function NPCManager() {
  const navigate = useNavigate()
  const [npcs, setNpcs] = useState<NPC[]>([])
  const [rooms, setRooms] = useState<Room[]>([])
  const [loading, setLoading] = useState(true)
  const [selectedNPC, setSelectedNPC] = useState<NPC | null>(null)
  const [saving, setSaving] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [form, setForm] = useState<NPCForm>({
    name: '',
    race: 'human',
    class: 'adventurer',
    level: 1,
    currentRoomId: 0,
    roaming: false,
    roamingTime: 0,
    roamingRooms: []
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

  const handleCreateNPC = useCallback(async () => {
    if (!form.name || !form.currentRoomId) {
      alert('Please fill in all required fields')
      return
    }

    setSaving(true)
    try {
      const token = localStorage.getItem('token')

      // Create character as NPC
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
          startingRoomId: form.currentRoomId
        })
      })

      if (!response.ok) {
        throw new Error('Failed to create NPC')
      }

      // Refresh NPC list
      const npcsResponse = await fetch('http://localhost:8080/npcs')
      const npcsData = await npcsResponse.json()
      setNpcs(npcsData.npcs || [])

      // Reset form
      setForm({
        name: '',
        race: 'human',
        class: 'adventurer',
        level: 1,
        currentRoomId: 0,
        roaming: false,
        roamingTime: 0,
        roamingRooms: []
      })
      setShowCreateForm(false)
    } catch (err) {
      console.error('Create NPC error:', err)
      alert('Failed to create NPC')
    } finally {
      setSaving(false)
    }
  }, [form])

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

      // Refresh NPC list
      const npcsResponse = await fetch('http://localhost:8080/npcs')
      const npcsData = await npcsResponse.json()
      setNpcs(npcsData.npcs || [])
      setSelectedNPC(null)
      setConfirmDelete(null)
    } catch (err) {
      console.error('Delete NPC error:', err)
      alert('Failed to delete NPC')
    }
  }, [])

  const handleUpdateNPCRoom = useCallback(async (npcId: number, newRoomId: number) => {
    try {
      const token = localStorage.getItem('token')

      const response = await fetch(`http://localhost:8080/characters/${npcId}`, {
        method: 'PATCH',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          currentRoomId: newRoomId
        })
      })

      if (!response.ok) {
        throw new Error('Failed to update NPC')
      }

      // Refresh NPC list
      const npcsResponse = await fetch('http://localhost:8080/npcs')
      const npcsData = await npcsResponse.json()
      setNpcs(npcsData.npcs || [])
      setSelectedNPC(null)
    } catch (err) {
      console.error('Update NPC error:', err)
      alert('Failed to update NPC')
    }
  }, [])

  const getRoomName = (roomId: number): string => {
    const room = rooms.find(r => r.id === roomId)
    return room ? room.name : `Room ${roomId}`
  }

  if (loading) {
    return <div className="p-8 text-white">Loading NPCs...</div>
  }

  return (
    <div className="flex h-screen bg-[#0a0a0f]">
      {/* Left Sidebar */}
      <div className="w-[280px] bg-[#1a1a2e] border-r border-[#333] flex flex-col">
        <div className="p-4 border-b border-[#333]">
          <Link
            to="/dashboard"
            className="block text-[#61dafb] no-underline p-2 rounded bg-[#16213e] text-center mb-2 hover:bg-[#1e2a4a]"
          >
            ← Dashboard
          </Link>
          <Link
            to="/map"
            className="block text-[#888] no-underline p-2 rounded bg-[#16213e] text-center hover:bg-[#1e2a4a]"
          >
            Map Builder
          </Link>
        </div>

        <div className="p-3 border-b border-[#333]">
          <h2 className="m-0 text-white text-lg">NPC Manager</h2>
          <p className="text-[#666] text-xs mt-1 mb-0">{npcs.length} NPCs</p>
        </div>

        {/* NPC List */}
        <div className="flex-1 overflow-y-auto p-3">
          <div className="flex flex-col gap-1">
            {npcs.map(npc => (
              <div
                key={npc.id}
                onClick={() => {
                  setSelectedNPC(npc)
                  setShowCreateForm(false)
                }}
                className={`p-2 cursor-pointer rounded text-xs ${
                  selectedNPC?.id === npc.id ? 'text-[#6c5ce7] bg-[#16213e]' : 'text-white'
                }`}
              >
                <div className="font-bold">{npc.name}</div>
                <div className="text-[#666]">
                  {npc.race} {npc.class} lv.{npc.level}
                </div>
                <div className="text-[#888] text-[10px]">
                  {getRoomName(npc.currentRoomId)}
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* Create NPC Button */}
        <div className="p-3 border-t border-[#333]">
          <button
            onClick={() => {
              setShowCreateForm(true)
              setSelectedNPC(null)
            }}
            className="w-full p-2 bg-[#27ae60] border-none rounded text-white cursor-pointer hover:bg-[#2ecc71]"
          >
            + Create NPC
          </button>
        </div>
      </div>

      {/* Main Content */}
      <div className="flex-1 overflow-y-auto p-6">
        {showCreateForm ? (
          <div className="max-w-[600px] mx-auto">
            <h2 className="mt-0 mb-4 text-white">Create New NPC</h2>

            <div className="bg-[#1a1a2e] rounded-lg p-4">
              <div className="mb-4">
                <label className="text-[#888] text-xs block mb-1">Name *</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="NPC name"
                  className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm"
                />
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-[#888] text-xs block mb-1">Race</label>
                  <select
                    value={form.race}
                    onChange={(e) => setForm({ ...form, race: e.target.value })}
                    className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm"
                  >
                    {RACES.map(race => (
                      <option key={race} value={race}>{race}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="text-[#888] text-xs block mb-1">Class</label>
                  <select
                    value={form.class}
                    onChange={(e) => setForm({ ...form, class: e.target.value })}
                    className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm"
                  >
                    {CLASSES.map(cls => (
                      <option key={cls} value={cls}>{cls}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-[#888] text-xs block mb-1">Level</label>
                  <input
                    type="number"
                    value={form.level}
                    onChange={(e) => setForm({ ...form, level: parseInt(e.target.value) || 1 })}
                    min={1}
                    className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm"
                  />
                </div>
                <div>
                  <label className="text-[#888] text-xs block mb-1">Room *</label>
                  <select
                    value={form.currentRoomId}
                    onChange={(e) => setForm({ ...form, currentRoomId: parseInt(e.target.value) })}
                    className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm"
                  >
                    <option value={0}>Select room...</option>
                    {rooms.map(room => (
                      <option key={room.id} value={room.id}>{room.name}</option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="mb-4">
                <label className="flex items-center gap-2 text-[#888] text-xs cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.roaming}
                    onChange={(e) => setForm({ ...form, roaming: e.target.checked })}
                    className="cursor-pointer"
                  />
                  Enable Roaming
                </label>
              </div>

              {form.roaming && (
                <div className="mb-4">
                  <label className="text-[#888] text-xs block mb-1">Roaming Time (ms, 0 = static)</label>
                  <input
                    type="number"
                    value={form.roamingTime}
                    onChange={(e) => setForm({ ...form, roamingTime: parseInt(e.target.value) || 0 })}
                    placeholder="e.g. 30000 for 30 seconds"
                    className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm"
                  />
                </div>
              )}

              <div className="flex gap-2">
                <button
                  onClick={handleCreateNPC}
                  disabled={saving}
                  className="flex-1 p-2 bg-[#27ae60] border-none rounded text-white cursor-pointer disabled:opacity-70"
                >
                  {saving ? 'Creating...' : 'Create NPC'}
                </button>
                <button
                  onClick={() => setShowCreateForm(false)}
                  className="flex-1 p-2 bg-[#16213e] border border-[#333] rounded text-[#888] cursor-pointer"
                >
                  Cancel
                </button>
              </div>
            </div>
          </div>
        ) : selectedNPC ? (
          <div className="max-w-[600px] mx-auto">
            <div className="flex justify-between items-center mb-4">
              <h2 className="m-0 text-white">{selectedNPC.name}</h2>
              <button
                onClick={() => setSelectedNPC(null)}
                className="bg-transparent border-none text-[#888] cursor-pointer text-xl"
              >
                ×
              </button>
            </div>

            <div className="bg-[#1a1a2e] rounded-lg p-4">
              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-[#888] text-xs block mb-1">Race</label>
                  <div className="text-white">{selectedNPC.race}</div>
                </div>
                <div>
                  <label className="text-[#888] text-xs block mb-1">Class</label>
                  <div className="text-white">{selectedNPC.class}</div>
                </div>
              </div>

              <div className="grid grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="text-[#888] text-xs block mb-1">Level</label>
                  <div className="text-white">{selectedNPC.level}</div>
                </div>
                <div>
                  <label className="text-[#888] text-xs block mb-1">ID</label>
                  <div className="text-white">#{selectedNPC.id}</div>
                </div>
              </div>

              <div className="mb-4">
                <label className="text-[#888] text-xs block mb-1">Current Room</label>
                <select
                  value={selectedNPC.currentRoomId}
                  onChange={(e) => handleUpdateNPCRoom(selectedNPC.id, parseInt(e.target.value))}
                  className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm"
                >
                  {rooms.map(room => (
                    <option key={room.id} value={room.id}>{room.name}</option>
                  ))}
                </select>
              </div>

              <div className="flex gap-2">
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
                      ? 'bg-[#e67e22] hover:bg-[#d35400]'
                      : 'bg-[#c0392b] hover:bg-[#e74c3c]'
                  }`}
                >
                  {confirmDelete === selectedNPC.id ? 'Confirm Delete?' : 'Delete NPC'}
                </button>
              </div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center h-full text-[#666]">
            <p className="mb-2">Select an NPC from the list or create a new one</p>
            <p className="text-xs text-[#555]">
              NPCs are characters controlled by the game<br/>
              They can roam between rooms or stay static
            </p>
          </div>
        )}
      </div>
    </div>
  )
}