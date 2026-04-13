import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useEffect, useState, useCallback, useMemo } from 'react'
import { Modal } from '../components/Modal'

export const Route = createFileRoute('/map')({
  component: MapBuilder,
})

interface Room {
  id: number
  name: string
  description: string
  isStartingRoom?: boolean
  exits: Record<string, number>
  atmosphere?: string
}

interface NPC {
  id: number
  name: string
  class: string
  race: string
  level: number
  currentRoomId: number
}

interface Equipment {
  id: number
  name: string
  description?: string
  roomId?: number
}

const OPPOSITE_DIR: Record<string, string> = {
  north: 'south', south: 'north',
  east: 'west', west: 'east',
  northeast: 'southwest', southwest: 'northeast',
  northwest: 'southeast', southeast: 'northwest',
  up: 'down', down: 'up'
}

const DIRECTION_OFFSETS: Record<string, { dx: number; dy: number }> = {
  north: { dx: 0, dy: -120 },
  south: { dx: 0, dy: 120 },
  east: { dx: 150, dy: 0 },
  west: { dx: -150, dy: 0 },
  northeast: { dx: 106, dy: -106 },
  northwest: { dx: -106, dy: -106 },
  southeast: { dx: 106, dy: 106 },
  southwest: { dx: -106, dy: 106 }
}

const ALL_DIRECTIONS = [
  'north', 'northeast', 'east', 'southeast', 'south',
  'southwest', 'west', 'northwest', 'up', 'down'
]

function MapBuilder() {
  const navigate = useNavigate()
  const [rooms, setRooms] = useState<Room[]>([])
  const [npcs, setNpcs] = useState<NPC[]>([])
  const [roomEquipment, setRoomEquipment] = useState<Record<number, Equipment[]>>({})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null)
  const [editingRoom, setEditingRoom] = useState<Room | null>(null)
  const [zoom, setZoom] = useState(1)
  const [currentZLevel, setCurrentZLevel] = useState(0)
  const [saving, setSaving] = useState(false)
  const [creating, setCreating] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState<number | null>(null)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [editForm, setEditForm] = useState({ name: '', description: '', exits: {} as Record<string, string> })
  const [newRoomForm, setNewRoomForm] = useState({ name: '', description: '' })

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) navigate({ to: '/login' })
  }, [navigate])

  useEffect(() => {
    const controller = new AbortController()
    Promise.all([
      fetch('http://localhost:8080/rooms', { signal: controller.signal }).then(res => res.json()),
      fetch('http://localhost:8080/npcs', { signal: controller.signal }).then(res => res.json())
    ]).then(([roomsData, npcsData]) => {
      setRooms(roomsData)
      setNpcs(npcsData.npcs || [])
      setLoading(false)
    }).catch(err => {
      if (err.name !== 'AbortError') {
        setError(err.message)
        setLoading(false)
      }
    })
    return () => controller.abort()
  }, [])

  useEffect(() => {
    if (!selectedRoom) return
    const roomId = selectedRoom.id
    const controller = new AbortController()
    fetch(`http://localhost:8080/rooms/${roomId}/equipment`, { signal: controller.signal })
      .then(res => res.json())
      .then(data => setRoomEquipment(prev => ({ ...prev, [roomId]: data.equipment || [] })))
      .catch(() => {})
    return () => controller.abort()
  }, [selectedRoom?.id])

  const zLevels = useMemo(() => {
    const zLevelsMap = new Map<number, number>()
    const visited = new Set<number>()
    const assignZLevel = (roomId: number, z: number) => {
      if (visited.has(roomId)) return
      visited.add(roomId)
      zLevelsMap.set(roomId, z)
      const room = rooms.find(r => r.id === roomId)
      if (!room) return
      for (const [dir, targetId] of Object.entries(room.exits || {})) {
        if (targetId) {
          const targetZ = dir === 'up' ? z + 1 : dir === 'down' ? z - 1 : z
          assignZLevel(targetId, targetZ)
        }
      }
    }
    const startRoom = rooms.find(r => r.isStartingRoom) || rooms[0]
    if (startRoom) assignZLevel(startRoom.id, 0)
    for (const room of rooms) {
      if (!zLevelsMap.has(room.id)) zLevelsMap.set(room.id, 0)
    }
    return zLevelsMap
  }, [rooms])

  const zLevelRange = useMemo(() => {
    const values = Array.from(zLevels.values())
    const minZ = Math.min(...values, 0)
    const maxZ = Math.max(...values, 0)
    return Array.from({ length: maxZ - minZ + 1 }, (_, i) => minZ + i)
  }, [zLevels])

  const getNPCsInRoom = useCallback((roomId: number) => npcs.filter(npc => npc.currentRoomId === roomId), [npcs])
  const getEquipmentInRoom = useCallback((roomId: number) => roomEquipment[roomId] || [], [roomEquipment])

  const nodePositions = useMemo(() => {
    const positions = new Map<number, { x: number; y: number }>()
    const positioned = new Set<number>()

    // Position rooms recursively from the starting room
    const positionRoom = (roomId: number, x: number, y: number) => {
      const roomZ = zLevels.get(roomId) || 0
      if (roomZ !== currentZLevel) return
      if (positioned.has(roomId)) return // Skip if already positioned on this pass
      positioned.add(roomId)
      positions.set(roomId, { x, y })
      const room = rooms.find(r => r.id === roomId)
      if (!room) return
      for (const [dir, targetId] of Object.entries(room.exits || {})) {
        if (targetId && dir !== 'up' && dir !== 'down') {
          const offset = DIRECTION_OFFSETS[dir] || { dx: 150, dy: 0 }
          positionRoom(targetId, x + offset.dx, y + offset.dy)
        }
      }
    }

    // Start from the starting room
    const startRoom = rooms.find(r => r.isStartingRoom) || rooms[0]
    if (startRoom) positionRoom(startRoom.id, 400, 300)

    // Position any orphaned rooms (not connected to starting room)
    let orphanX = 800
    for (const room of rooms) {
      const roomZ = zLevels.get(room.id) || 0
      if (roomZ === currentZLevel && !positioned.has(room.id)) {
        positions.set(room.id, { x: orphanX, y: 300 })
        orphanX += 150
      }
    }
    return positions
  }, [rooms, zLevels, currentZLevel])

  const handleEditRoom = useCallback((room: Room) => {
    setEditingRoom(room)
    setEditForm({
      name: room.name,
      description: room.description,
      exits: Object.fromEntries(Object.entries(room.exits || {}).map(([dir, id]) => [dir, String(id)]))
    })
  }, [])

  const handleSaveRoom = useCallback(async () => {
    if (!editingRoom) return
    setSaving(true)
    try {
      const token = localStorage.getItem('token')
      const exits: Record<string, number> = {}
      for (const [dir, id] of Object.entries(editForm.exits)) {
        const numId = parseInt(id)
        if (!isNaN(numId)) exits[dir] = numId
      }

      // Save the room being edited
      const response = await fetch(`http://localhost:8080/rooms/${editingRoom.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ name: editForm.name, description: editForm.description, exits })
      })
      if (!response.ok) throw new Error('Failed to save room')

      // Update reverse exits for bidirectional linking
      const oldExits = editingRoom.exits || {}
      const dirs = ALL_DIRECTIONS

      for (const dir of dirs) {
        const newTargetId = exits[dir]
        const oldTargetId = oldExits[dir]
        const reverseDir = OPPOSITE_DIR[dir]

        // If exit changed or is new
        if (newTargetId && newTargetId !== oldTargetId) {
          // Add reverse exit to the new target room
          const targetRoom = rooms.find(r => r.id === newTargetId)
          if (targetRoom && reverseDir) {
            const updatedTargetExits = { ...targetRoom.exits, [reverseDir]: editingRoom.id }
            await fetch(`http://localhost:8080/rooms/${newTargetId}`, {
              method: 'PUT',
              headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
              body: JSON.stringify({ ...targetRoom, exits: updatedTargetExits })
            })
          }
        }

        // If exit was removed or changed to different room, remove reverse from old target
        if (oldTargetId && oldTargetId !== newTargetId) {
          const oldTargetRoom = rooms.find(r => r.id === oldTargetId)
          if (oldTargetRoom && reverseDir && oldTargetRoom.exits?.[reverseDir] === editingRoom.id) {
            const { [reverseDir]: _, ...remainingExits } = oldTargetRoom.exits
            await fetch(`http://localhost:8080/rooms/${oldTargetId}`, {
              method: 'PUT',
              headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
              body: JSON.stringify({ ...oldTargetRoom, exits: remainingExits || {} })
            })
          }
        }
      }

      // Refresh rooms from server
      const roomsResponse = await fetch('http://localhost:8080/rooms')
      const roomsData = await roomsResponse.json()
      setRooms(roomsData)
      setEditingRoom(null)
      setSelectedRoom(null)
    } catch (err) {
      console.error('Save error:', err)
      alert('Failed to save room')
    } finally {
      setSaving(false)
    }
  }, [editingRoom, editForm, rooms])

  const handleCreateRoom = useCallback(async (fromRoom: Room | null, direction: string) => {
    setCreating(true)
    try {
      const token = localStorage.getItem('token')
      const newRoomName = newRoomForm.name || `New Room ${rooms.length + 1}`
      const newRoomDesc = newRoomForm.description || 'A newly created room.'

      const createResponse = await fetch('http://localhost:8080/rooms', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
        body: JSON.stringify({ name: newRoomName, description: newRoomDesc, exits: {} })
      })
      if (!createResponse.ok) throw new Error('Failed to create room')
      const newRoom = await createResponse.json()

      if (fromRoom && direction) {
        const updatedExits = { ...fromRoom.exits, [direction]: newRoom.id }
        await fetch(`http://localhost:8080/rooms/${fromRoom.id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
          body: JSON.stringify({ ...fromRoom, exits: updatedExits })
        })
        const reverseDirection = OPPOSITE_DIR[direction]
        await fetch(`http://localhost:8080/rooms/${newRoom.id}`, {
          method: 'PUT',
          headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` },
          body: JSON.stringify({ name: newRoomName, description: newRoomDesc, exits: { [reverseDirection]: fromRoom.id } })
        })
      }

      const roomsResponse = await fetch('http://localhost:8080/rooms')
      const roomsData = await roomsResponse.json()
      setRooms(roomsData)
      setSelectedRoom(roomsData.find((r: Room) => r.id === newRoom.id) || null)
      setShowCreateModal(false)
      setNewRoomForm({ name: '', description: '' })
    } catch (err) {
      console.error('Create room error:', err)
      alert('Failed to create room')
    } finally {
      setCreating(false)
    }
  }, [rooms.length, newRoomForm])

  const handleCreateStandaloneRoom = useCallback(async () => {
    await handleCreateRoom(null, '')
  }, [handleCreateRoom])

  const handleDeleteRoom = useCallback(async (roomId: number) => {
    try {
      const token = localStorage.getItem('token')
      const response = await fetch(`http://localhost:8080/rooms/${roomId}`, {
        method: 'DELETE',
        headers: { 'Authorization': `Bearer ${token}` }
      })
      if (!response.ok) throw new Error('Failed to delete room')
      const roomsResponse = await fetch('http://localhost:8080/rooms')
      const roomsData = await roomsResponse.json()
      setRooms(roomsData)
      setSelectedRoom(null)
      setConfirmDelete(null)
    } catch (err) {
      console.error('Delete room error:', err)
      alert('Failed to delete room')
    }
  }, [])

  if (loading) return <div className="p-8 text-text">Loading map...</div>
  if (error) return <div className="p-8 text-danger">Error: {error}</div>

  return (
    <div className="flex h-screen bg-surface">
      {/* Left Sidebar */}
      <div className="w-[220px] bg-surface-muted border-r border-border flex flex-col">
        <div className="p-4 border-b border-border">
          <Link to="/dashboard" className="block text-primary no-underline p-2 rounded bg-surface-dark text-center mb-2 hover:bg-surface-darker">
            ← Dashboard
          </Link>
          <button onClick={() => setShowCreateModal(true)} className="w-full p-2 bg-primary border-none rounded text-white cursor-pointer hover:bg-primary-hover">
            + Add Room
          </button>
        </div>

        <div className="p-3 border-b border-border">
          <label className="text-text-muted text-xs block mb-2">Floor (Z-Level)</label>
          <div className="flex gap-1 flex-wrap">
            {zLevelRange.map(z => (
              <button key={z} onClick={() => setCurrentZLevel(z)}
                className={`px-2 py-1 rounded text-xs ${currentZLevel === z ? 'bg-primary border-primary-hover border' : 'bg-surface-dark border-border border'} text-white cursor-pointer`}>
                {z === 0 ? 'G' : z > 0 ? `+${z}` : `${z}`}
              </button>
            ))}
          </div>
        </div>

        <div className="p-3 text-text-muted text-xs border-b border-border">
          <div>Total: {rooms.length} rooms</div>
          <div>Floor {currentZLevel}: {Array.from(zLevels.values()).filter(z => z === currentZLevel).length}</div>
          <div>NPCs: {npcs.length}</div>
        </div>

        <div className="flex-1 overflow-y-auto p-3">
          <h4 className="m-0 mb-2 text-text-muted text-xs">Rooms on Floor {currentZLevel}</h4>
          <div className="flex flex-col gap-1">
            {rooms.filter(r => (zLevels.get(r.id) || 0) === currentZLevel).map(room => (
              <div key={room.id} onClick={() => setSelectedRoom(room)}
                className={`p-2 cursor-pointer rounded text-xs ${selectedRoom?.id === room.id ? 'text-text bg-surface-dark' : 'text-text'}`}>
                <span className="truncate">{room.name}</span>
                {room.isStartingRoom && <span> ⭐</span>}
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Map Area */}
      <div className="flex-1 overflow-hidden relative">
        <div className="absolute top-0 left-0 right-0 p-3 bg-surface-muted border-b border-border flex justify-between items-center z-10">
          <h1 className="m-0 text-text text-lg">Map Builder — Floor {currentZLevel}</h1>
          <div className="flex gap-2 items-center">
            <button onClick={() => setZoom(z => Math.max(z - 0.25, 0.5))} className="px-2 py-1 bg-danger border-none rounded text-white cursor-pointer">−</button>
            <span className="text-text-muted text-xs w-12 text-center">{Math.round(zoom * 100)}%</span>
            <button onClick={() => setZoom(z => Math.min(z + 0.25, 2))} className="px-2 py-1 bg-primary border-none rounded text-white cursor-pointer">+</button>
          </div>
        </div>

        <div className="mt-[50px] h-[calc(100%-50px)] overflow-auto p-6">
          <div className="relative w-[3000px] h-[3000px]" style={{ transform: `scale(${zoom})`, transformOrigin: 'top left' }}>
            <svg className="absolute top-0 left-0 w-full h-full pointer-events-none">
              {rooms.map(room => {
                const pos = nodePositions.get(room.id)
                if (!pos) return null
                return Object.entries(room.exits || {}).map(([dir, targetId]) => {
                  const targetPos = nodePositions.get(targetId)
                  if (!targetPos || dir === 'up' || dir === 'down') return null
                  return <line key={`${room.id}-${targetId}-${dir}`} x1={pos.x + 60} y1={pos.y + 35} x2={targetPos.x + 60} y2={targetPos.y + 35} stroke="#9ca3af" strokeWidth={2} />
                })
              })}
            </svg>

            {rooms.map(room => {
              const pos = nodePositions.get(room.id)
              if (!pos) return null
              const isSelected = selectedRoom?.id === room.id
              const roomNpcs = getNPCsInRoom(room.id)
              const roomItems = getEquipmentInRoom(room.id)
              const isColored = room.isStartingRoom || isSelected
              return (
                <div key={room.id} onClick={() => setSelectedRoom(room)}
                  className={`w-[120px] min-h-[65px] p-2 rounded-lg cursor-pointer transition-all ${room.isStartingRoom ? 'bg-primary' : isSelected ? 'bg-primary-hover shadow-lg' : 'bg-surface-dark'} ${isSelected ? 'border-2 border-accent' : 'border-2 border-border'}`}
                  style={{ position: 'absolute', left: pos.x, top: pos.y, zIndex: 1 }}>
                  <div className={`font-bold text-xs text-center truncate ${isColored ? 'text-white' : 'text-text'}`}>{room.name}{room.isStartingRoom && ' ⭐'}</div>
                  <div className="text-text-muted text-[10px] text-center">#{room.id}</div>
                  <div className="flex justify-center gap-1 mt-1">
                    {roomNpcs.length > 0 && <span className={`text-[10px] ${isColored ? 'text-warning' : 'text-warning'}`} title={`${roomNpcs.length} NPCs`}>👥{roomNpcs.length}</span>}
                    {roomItems.length > 0 && <span className={`text-[10px] ${isColored ? 'text-white/80' : 'text-success'}`} title={`${roomItems.length} items`}>📦{roomItems.length}</span>}
                  </div>
                  <div className="flex justify-center gap-0.5 mt-0.5">
                    {room.exits?.up && <span className={`text-[8px] ${isColored ? 'text-white/90' : 'text-warning'}`}>▲{room.exits.up}</span>}
                    {room.exits?.down && <span className={`text-[8px] ${isColored ? 'text-white/90' : 'text-success'}`}>▼{room.exits.down}</span>}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      </div>

      {/* Right Panel */}
      <div className="w-[300px] bg-surface-muted border-l border-border flex flex-col">
        {editingRoom ? (
          <>
            <div className="p-3 border-b border-border flex justify-between items-center">
              <h3 className="m-0 text-text text-base">Edit Room</h3>
              <button onClick={() => setEditingRoom(null)} className="bg-transparent border-none text-text-muted cursor-pointer text-xl">×</button>
            </div>
            <div className="p-3 flex-1 overflow-y-auto">
              <div className="mb-3">
                <label className="text-text-muted text-xs block mb-1">Name</label>
                <input type="text" value={editForm.name} onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
              </div>
              <div className="mb-3">
                <label className="text-text-muted text-xs block mb-1">Description</label>
                <textarea value={editForm.description} onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                  rows={4} className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y" />
              </div>
              <div className="mb-3">
                <label className="text-text-muted text-xs block mb-2">Exits</label>
                {ALL_DIRECTIONS.map(dir => (
                  <div key={dir} className="flex items-center gap-2 mb-1">
                    <span className="w-[50px] text-text-muted text-xs">{dir}:</span>
                    <input type="text" value={editForm.exits[dir] || ''} onChange={(e) => setEditForm({ ...editForm, exits: { ...editForm.exits, [dir]: e.target.value } })}
                      placeholder="room id" className="flex-1 p-1 bg-surface border border-border rounded text-text text-xs" />
                  </div>
                ))}
              </div>
            </div>
            <div className="p-3 border-t border-border">
              <button onClick={handleSaveRoom} disabled={saving}
                className="w-full p-2 bg-primary border-none rounded text-white cursor-pointer mb-2 disabled:opacity-70">
                {saving ? 'Saving...' : 'Save Changes'}
              </button>
              <button onClick={() => setEditingRoom(null)} className="w-full p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer">
                Cancel
              </button>
            </div>
          </>
        ) : selectedRoom ? (
          <>
            <div className="p-3 border-b border-border flex justify-between items-center">
              <h3 className="m-0 text-text text-base">{selectedRoom.name}{selectedRoom.isStartingRoom && <span className="text-warning"> ⭐</span>}</h3>
              <button onClick={() => setSelectedRoom(null)} className="bg-transparent border-none text-text-muted cursor-pointer text-xl">×</button>
            </div>
            <div className="p-3 flex-1 overflow-y-auto">
              <div className="text-text-muted text-[10px] mb-2">Room ID: {selectedRoom.id}{selectedRoom.atmosphere && ` • ${selectedRoom.atmosphere}`}</div>
              <div className="text-text mb-3 text-sm">{selectedRoom.description}</div>

              {getNPCsInRoom(selectedRoom.id).length > 0 && (
                <div className="mb-3">
                  <strong className="text-warning text-xs">NPCs:</strong>
                  <div className="mt-1">
                    {getNPCsInRoom(selectedRoom.id).map(npc => (
                      <div key={npc.id} className="p-1 bg-surface-dark rounded mb-1 text-xs text-text">
                        {npc.name} <span className="text-text-muted">({npc.race} {npc.class} lv.{npc.level})</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {getEquipmentInRoom(selectedRoom.id).length > 0 && (
                <div className="mb-3">
                  <strong className="text-success text-xs">Items:</strong>
                  <div className="mt-1">
                    {getEquipmentInRoom(selectedRoom.id).map(item => (
                      <div key={item.id} className="p-1 bg-surface-dark rounded mb-1 text-xs text-text">{item.name}</div>
                    ))}
                  </div>
                </div>
              )}

              <div className="mb-3">
                <strong className="text-accent text-xs">Exits:</strong>
                <div className="mt-1">
                  {ALL_DIRECTIONS.map(dir => {
                    const targetId = selectedRoom.exits?.[dir]
                    const targetRoom = rooms.find(r => r.id === targetId)
                    const isZExit = dir === 'up' || dir === 'down'

                    if (targetId && targetRoom) {
                      return (
                        <div key={dir} onClick={() => { if (isZExit) setCurrentZLevel(zLevels.get(targetId) || 0); setSelectedRoom(targetRoom) }}
                          className={`p-1 my-1 rounded cursor-pointer text-xs ${isZExit ? (dir === 'up' ? 'bg-warning/20 border border-warning' : 'bg-success/20 border border-success') : 'bg-surface-dark border-none'}`}>
                          <strong>{dir}</strong> → {targetRoom.name}
                          {isZExit && <span className="text-text-muted ml-1 text-[10px]">(z={zLevels.get(targetId) || 0})</span>}
                        </div>
                      )
                    } else if (targetId) {
                      return <div key={dir} className="p-1 my-1 rounded text-xs bg-surface-dark border-none"><strong>{dir}</strong> → <span className="text-text-muted">Room #{targetId}</span></div>
                    } else {
                      return (
                        <div key={dir} className="flex items-center gap-2 my-1">
                          <div className="flex-1 p-1 rounded text-xs bg-surface-muted border border-border text-text-muted"><strong>{dir}</strong> → none</div>
                          <button onClick={() => { setNewRoomForm({ name: '', description: '' }); setShowCreateModal(true); }}
                            disabled={creating} className="px-2 py-1 bg-primary border-none rounded text-white text-xs cursor-pointer hover:bg-primary-hover disabled:opacity-50" title={`Create room to the ${dir}`}>
                            +
                          </button>
                        </div>
                      )
                    }
                  })}
                </div>
              </div>
            </div>
            <div className="p-3 border-t border-border flex gap-2">
              <button onClick={() => handleEditRoom(selectedRoom)} className="flex-1 p-2 bg-accent border-none rounded text-white cursor-pointer">Edit Room</button>
              <button onClick={() => { if (confirmDelete === selectedRoom.id) { handleDeleteRoom(selectedRoom.id) } else { setConfirmDelete(selectedRoom.id) } }}
                className={`flex-1 p-2 border-none rounded text-white cursor-pointer ${confirmDelete === selectedRoom.id ? 'bg-warning hover:bg-warning/80' : 'bg-danger hover:bg-danger-hover'}`}>
                {confirmDelete === selectedRoom.id ? 'Confirm Delete?' : 'Delete Room'}
              </button>
            </div>
          </>
        ) : (
          <div className="flex-1 flex flex-col justify-center items-center text-text-muted text-center p-4">
            <p className="m-0 mb-2">Click a room to see details</p>
            <p className="text-xs m-0 text-text-muted">👥 NPCs in room<br/>📦 Items on ground<br/>▲/▼ Stairs</p>
          </div>
        )}
      </div>

      {/* Create Room Modal */}
      <Modal isOpen={showCreateModal} onClose={() => setShowCreateModal(false)} title="Create New Room">
        <div className="mb-4">
          <label className="text-text-muted text-xs block mb-1">Room Name</label>
          <input type="text" value={newRoomForm.name} onChange={(e) => setNewRoomForm({ ...newRoomForm, name: e.target.value })}
            placeholder="Enter room name" className="w-full p-2 bg-surface border border-border rounded text-text text-sm" />
        </div>
        <div className="mb-4">
          <label className="text-text-muted text-xs block mb-1">Description</label>
          <textarea value={newRoomForm.description} onChange={(e) => setNewRoomForm({ ...newRoomForm, description: e.target.value })}
            placeholder="Enter room description" rows={4} className="w-full p-2 bg-surface border border-border rounded text-text text-sm resize-y" />
        </div>
        <div className="flex gap-2">
          <button onClick={handleCreateStandaloneRoom} disabled={creating}
            className="flex-1 p-2 bg-primary border-none rounded text-white cursor-pointer disabled:opacity-70">
            {creating ? 'Creating...' : 'Create Room'}
          </button>
          <button onClick={() => setShowCreateModal(false)} className="flex-1 p-2 bg-surface-dark border border-border rounded text-text-muted cursor-pointer">
            Cancel
          </button>
        </div>
      </Modal>
    </div>
  )
}