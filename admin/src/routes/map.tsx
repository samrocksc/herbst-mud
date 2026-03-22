import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useEffect, useState, useCallback, useMemo } from 'react'

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

// Move constant outside component to avoid recreation
const OPPOSITE_DIR: Record<string, string> = {
  north: 'south',
  south: 'north',
  east: 'west',
  west: 'east',
  up: 'down',
  down: 'up'
}

const DIRECTION_OFFSETS: Record<string, { dx: number; dy: number }> = {
  north: { dx: 0, dy: -120 },
  south: { dx: 0, dy: 120 },
  east: { dx: 150, dy: 0 },
  west: { dx: -150, dy: 0 }
}

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
  const [editForm, setEditForm] = useState({
    name: '',
    description: '',
    exits: {} as Record<string, string>
  })

  // Auth check - runs once on mount
  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate({ to: '/login' })
    }
  }, [navigate])

  // Load rooms and NPCs - runs once on mount
  useEffect(() => {
    const controller = new AbortController()

    Promise.all([
      fetch('http://localhost:8080/rooms', { signal: controller.signal }).then(res => res.json()),
      fetch('http://localhost:8080/npcs', { signal: controller.signal }).then(res => res.json())
    ])
      .then(([roomsData, npcsData]) => {
        setRooms(roomsData)
        setNpcs(npcsData.npcs || [])
        setLoading(false)
      })
      .catch(err => {
        if (err.name !== 'AbortError') {
          setError(err.message)
          setLoading(false)
        }
      })

    return () => controller.abort()
  }, [])

  // Load equipment for selected room - only when selectedRoom changes
  useEffect(() => {
    if (!selectedRoom) return

    const roomId = selectedRoom.id
    const controller = new AbortController()

    fetch(`http://localhost:8080/rooms/${roomId}/equipment`, { signal: controller.signal })
      .then(res => res.json())
      .then(data => {
        setRoomEquipment(prev => ({ ...prev, [roomId]: data.equipment || [] }))
      })
      .catch(() => {})

    return () => controller.abort()
  }, [selectedRoom?.id]) // Only depend on the room ID, not the room object

  // Memoize z-level calculation
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
    if (startRoom) {
      assignZLevel(startRoom.id, 0)
    }

    for (const room of rooms) {
      if (!zLevelsMap.has(room.id)) {
        zLevelsMap.set(room.id, 0)
      }
    }

    return zLevelsMap
  }, [rooms])

  // Memoize z-level range
  const zLevelRange = useMemo(() => {
    const values = Array.from(zLevels.values())
    const minZ = Math.min(...values, 0)
    const maxZ = Math.max(...values, 0)
    return Array.from({ length: maxZ - minZ + 1 }, (_, i) => minZ + i)
  }, [zLevels])

  // Memoize helper functions
  const getNPCsInRoom = useCallback((roomId: number): NPC[] =>
    npcs.filter(npc => npc.currentRoomId === roomId)
  , [npcs])

  const getEquipmentInRoom = useCallback((roomId: number): Equipment[] =>
    roomEquipment[roomId] || []
  , [roomEquipment])

  // Memoize node positions for current z-level
  const nodePositions = useMemo(() => {
    const positions = new Map<number, { x: number; y: number }>()
    const visited = new Set<number>()

    const positionRoom = (roomId: number, x: number, y: number) => {
      if (visited.has(roomId)) return
      const roomZ = zLevels.get(roomId) || 0
      if (roomZ !== currentZLevel) return

      visited.add(roomId)
      positions.set(roomId, { x, y })

      const room = rooms.find(r => r.id === roomId)
      if (!room) return

      for (const [dir, targetId] of Object.entries(room.exits || {})) {
        if (targetId && dir !== 'up' && dir !== 'down' && !visited.has(targetId)) {
          const offset = DIRECTION_OFFSETS[dir] || { dx: 150, dy: 0 }
          positionRoom(targetId, x + offset.dx, y + offset.dy)
        }
      }
    }

    const centerRoom = rooms.find(r => r.isStartingRoom) || rooms[0]
    if (centerRoom) {
      positionRoom(centerRoom.id, 400, 300)
    }

    // Position orphan rooms
    let orphanX = 800
    for (const room of rooms) {
      const roomZ = zLevels.get(room.id) || 0
      if (roomZ === currentZLevel && !visited.has(room.id)) {
        positions.set(room.id, { x: orphanX, y: 300 })
        orphanX += 150
      }
    }

    return positions
  }, [rooms, zLevels, currentZLevel])

  // Memoize handlers
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
        if (!isNaN(numId)) {
          exits[dir] = numId
        }
      }

      const response = await fetch(`http://localhost:8080/rooms/${editingRoom.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ name: editForm.name, description: editForm.description, exits })
      })

      if (!response.ok) {
        throw new Error('Failed to save room')
      }

      setRooms(prev => prev.map(r =>
        r.id === editingRoom.id
          ? { ...r, name: editForm.name, description: editForm.description, exits }
          : r
      ))
      setEditingRoom(null)
      setSelectedRoom(null)
    } catch (err) {
      console.error('Save error:', err)
      alert('Failed to save room')
    } finally {
      setSaving(false)
    }
  }, [editingRoom, editForm])

  const handleCreateRoom = useCallback(async (fromRoom: Room, direction: string) => {
    setCreating(true)
    try {
      const token = localStorage.getItem('token')

      const newRoomName = `New Room ${rooms.length + 1}`
      const newRoomDesc = 'A newly created room.'

      const createResponse = await fetch('http://localhost:8080/rooms', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          name: newRoomName,
          description: newRoomDesc,
          exits: {}
        })
      })

      if (!createResponse.ok) {
        throw new Error('Failed to create room')
      }

      const newRoom = await createResponse.json()

      const updatedExits = { ...fromRoom.exits, [direction]: newRoom.id }

      await fetch(`http://localhost:8080/rooms/${fromRoom.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ ...fromRoom, exits: updatedExits })
      })

      const reverseDirection = OPPOSITE_DIR[direction]
      await fetch(`http://localhost:8080/rooms/${newRoom.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          name: newRoomName,
          description: newRoomDesc,
          exits: { [reverseDirection]: fromRoom.id }
        })
      })

      const roomsResponse = await fetch('http://localhost:8080/rooms')
      const roomsData = await roomsResponse.json()
      setRooms(roomsData)

      setSelectedRoom(roomsData.find((r: Room) => r.id === newRoom.id) || null)
    } catch (err) {
      console.error('Create room error:', err)
      alert('Failed to create room')
    } finally {
      setCreating(false)
    }
  }, [rooms.length])

  const handleDeleteRoom = useCallback(async (roomId: number) => {
    try {
      const token = localStorage.getItem('token')

      const response = await fetch(`http://localhost:8080/rooms/${roomId}`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })

      if (!response.ok) {
        throw new Error('Failed to delete room')
      }

      // Refresh rooms
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

  if (loading) {
    return <div className="p-8 text-white">Loading map...</div>
  }

  if (error) {
    return <div className="p-8 text-red-400">Error: {error}</div>
  }

  return (
    <div className="flex h-screen bg-[#0a0a0f]">
      {/* Left Sidebar */}
      <div className="w-[220px] bg-[#1a1a2e] border-r border-[#333] flex flex-col">
        <div className="p-4 border-b border-[#333]">
          <Link
            to="/dashboard"
            className="block text-[#61dafb] no-underline p-2 rounded bg-[#16213e] text-center mb-2 hover:bg-[#1e2a4a]"
          >
            ← Dashboard
          </Link>
        </div>

        {/* Z-Level Selector */}
        <div className="p-3 border-b border-[#333]">
          <label className="text-[#888] text-xs block mb-2">Floor (Z-Level)</label>
          <div className="flex gap-1 flex-wrap">
            {zLevelRange.map(z => (
              <button
                key={z}
                onClick={() => setCurrentZLevel(z)}
                className={`px-2 py-1 rounded text-xs ${
                  currentZLevel === z
                    ? 'bg-[#27ae60] border-[#2ecc71] border'
                    : 'bg-[#16213e] border-[#333] border'
                } text-white cursor-pointer`}
              >
                {z === 0 ? 'G' : z > 0 ? `+${z}` : `${z}`}
              </button>
            ))}
          </div>
        </div>

        {/* Stats */}
        <div className="p-3 text-[#666] text-xs border-b border-[#333]">
          <div>Total: {rooms.length} rooms</div>
          <div>Floor {currentZLevel}: {Array.from(zLevels.values()).filter(z => z === currentZLevel).length}</div>
          <div>NPCs: {npcs.length}</div>
        </div>

        {/* Room List */}
        <div className="flex-1 overflow-y-auto p-3">
          <h4 className="m-0 mb-2 text-[#888] text-xs">Rooms on Floor {currentZLevel}</h4>
          <div className="flex flex-col gap-1">
            {rooms.filter(r => (zLevels.get(r.id) || 0) === currentZLevel).map(room => (
              <div
                key={room.id}
                onClick={() => setSelectedRoom(room)}
                className={`p-2 cursor-pointer rounded text-xs flex justify-between items-center ${
                  selectedRoom?.id === room.id ? 'text-[#6c5ce7] bg-[#16213e]' : 'text-white'
                }`}
              >
                <span className="overflow-hidden text-ellipsis whitespace-nowrap">{room.name}</span>
                {room.isStartingRoom && <span title="Starting room">⭐</span>}
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Map Area */}
      <div className="flex-1 overflow-hidden relative">
        {/* Header */}
        <div className="absolute top-0 left-0 right-0 p-3 bg-[rgba(26,26,46,0.95)] border-b border-[#333] flex justify-between items-center z-10">
          <h1 className="m-0 text-white text-lg">Map Builder — Floor {currentZLevel}</h1>
          <div className="flex gap-2 items-center">
            <button onClick={() => setZoom(z => Math.max(z - 0.25, 0.5))} className="px-2 py-1 bg-[#e74c3c] border-none rounded text-white cursor-pointer">−</button>
            <span className="text-[#888] text-xs w-12 text-center">{Math.round(zoom * 100)}%</span>
            <button onClick={() => setZoom(z => Math.min(z + 0.25, 2))} className="px-2 py-1 bg-[#27ae60] border-none rounded text-white cursor-pointer">+</button>
          </div>
        </div>

        {/* Map Canvas */}
        <div className="mt-[50px] h-[calc(100%-50px)] overflow-auto p-6">
          <div className="relative w-[3000px] h-[3000px]" style={{ transform: `scale(${zoom})`, transformOrigin: 'top left' }}>
            {/* Edges */}
            <svg className="absolute top-0 left-0 w-full h-full pointer-events-none">
              {rooms.map(room => {
                const pos = nodePositions.get(room.id)
                if (!pos) return null

                return Object.entries(room.exits || {}).map(([dir, targetId]) => {
                  const targetPos = nodePositions.get(targetId)
                  if (!targetPos || dir === 'up' || dir === 'down') return null

                  return (
                    <line
                      key={`${room.id}-${targetId}-${dir}`}
                      x1={pos.x + 60}
                      y1={pos.y + 35}
                      x2={targetPos.x + 60}
                      y2={targetPos.y + 35}
                      stroke="#444"
                      strokeWidth={2}
                    />
                  )
                })
              })}
            </svg>

            {/* Room Nodes */}
            {rooms.map(room => {
              const pos = nodePositions.get(room.id)
              if (!pos) return null

              const isSelected = selectedRoom?.id === room.id
              const roomNpcs = getNPCsInRoom(room.id)
              const roomItems = getEquipmentInRoom(room.id)
              const hasUp = room.exits?.up
              const hasDown = room.exits?.down

              return (
                <div
                  key={room.id}
                  onClick={() => setSelectedRoom(room)}
                  className={`w-[120px] min-h-[65px] p-2 rounded-lg cursor-pointer transition-all ${
                    room.isStartingRoom
                      ? 'bg-[#27ae60]'
                      : isSelected
                      ? 'bg-[#6c5ce7] shadow-[0_0_15px_rgba(108,92,231,0.5)] z-10'
                      : 'bg-[#2d5a27]'
                  } ${
                    isSelected ? 'border-2 border-[#a29bfe]' : 'border-2 border-[#1a3a1a]'
                  }`}
                  style={{
                    position: 'absolute',
                    left: pos.x,
                    top: pos.y,
                    boxShadow: isSelected ? '0 0 15px rgba(108, 92, 231, 0.5)' : '0 2px 8px rgba(0,0,0,0.3)'
                  }}
                >
                  <div className="text-white font-bold text-xs text-center truncate">
                    {room.name}
                    {room.isStartingRoom && ' ⭐'}
                  </div>
                  <div className="text-[#888] text-[10px] text-center">#{room.id}</div>
                  <div className="flex justify-center gap-1 mt-1">
                    {roomNpcs.length > 0 && (
                      <span className="text-[10px] text-[#f39c12]" title={`${roomNpcs.length} NPCs`}>👥{roomNpcs.length}</span>
                    )}
                    {roomItems.length > 0 && (
                      <span className="text-[10px] text-[#3498db]" title={`${roomItems.length} items`}>📦{roomItems.length}</span>
                    )}
                  </div>
                  <div className="flex justify-center gap-0.5 mt-0.5">
                    {hasUp && <span className="text-[8px] text-[#e17055]">▲{hasUp}</span>}
                    {hasDown && <span className="text-[8px] text-[#74b9ff]">▼{hasDown}</span>}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      </div>

      {/* Right Panel */}
      <div className="w-[300px] bg-[#1a1a2e] border-l border-[#333] flex flex-col">
        {editingRoom ? (
          <>
            <div className="p-3 border-b border-[#333] flex justify-between items-center">
              <h3 className="m-0 text-white text-base">Edit Room</h3>
              <button onClick={() => setEditingRoom(null)} className="bg-transparent border-none text-[#888] cursor-pointer text-xl">×</button>
            </div>
            <div className="p-3 flex-1 overflow-y-auto">
              <div className="mb-3">
                <label className="text-[#888] text-xs block mb-1">Name</label>
                <input
                  type="text"
                  value={editForm.name}
                  onChange={(e) => setEditForm({ ...editForm, name: e.target.value })}
                  className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm"
                />
              </div>
              <div className="mb-3">
                <label className="text-[#888] text-xs block mb-1">Description</label>
                <textarea
                  value={editForm.description}
                  onChange={(e) => setEditForm({ ...editForm, description: e.target.value })}
                  rows={4}
                  className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-white text-sm resize-y"
                />
              </div>
              <div className="mb-3">
                <label className="text-[#888] text-xs block mb-2">Exits</label>
                {['north', 'south', 'east', 'west', 'up', 'down'].map(dir => (
                  <div key={dir} className="flex items-center gap-2 mb-1">
                    <span className="w-[50px] text-[#666] text-xs">{dir}:</span>
                    <input
                      type="text"
                      value={editForm.exits[dir] || ''}
                      onChange={(e) => setEditForm({ ...editForm, exits: { ...editForm.exits, [dir]: e.target.value } })}
                      placeholder="room id"
                      className="flex-1 p-1 bg-[#16213e] border border-[#333] rounded text-white text-xs"
                    />
                  </div>
                ))}
              </div>
            </div>
            <div className="p-3 border-t border-[#333]">
              <button
                onClick={handleSaveRoom}
                disabled={saving}
                className="w-full p-2 bg-[#27ae60] border-none rounded text-white cursor-pointer mb-2 disabled:opacity-70"
              >
                {saving ? 'Saving...' : 'Save Changes'}
              </button>
              <button onClick={() => setEditingRoom(null)} className="w-full p-2 bg-[#16213e] border border-[#333] rounded text-[#888] cursor-pointer">
                Cancel
              </button>
            </div>
          </>
        ) : selectedRoom ? (
          <>
            <div className="p-3 border-b border-[#333] flex justify-between items-center">
              <h3 className="m-0 text-white text-base">
                {selectedRoom.name}
                {selectedRoom.isStartingRoom && <span className="text-[#f39c12]"> ⭐</span>}
              </h3>
              <button onClick={() => setSelectedRoom(null)} className="bg-transparent border-none text-[#888] cursor-pointer text-xl">×</button>
            </div>
            <div className="p-3 flex-1 overflow-y-auto">
              <div className="text-[#666] text-[10px] mb-2">
                Room ID: {selectedRoom.id}
                {selectedRoom.atmosphere && ` • ${selectedRoom.atmosphere}`}
              </div>
              <div className="text-[#aaa] mb-3 text-sm">{selectedRoom.description}</div>

              {getNPCsInRoom(selectedRoom.id).length > 0 && (
                <div className="mb-3">
                  <strong className="text-[#f39c12] text-xs">NPCs:</strong>
                  <div className="mt-1">
                    {getNPCsInRoom(selectedRoom.id).map(npc => (
                      <div key={npc.id} className="p-1 bg-[#16213e] rounded mb-1 text-xs text-white">
                        {npc.name} <span className="text-[#666]">({npc.race} {npc.class} lv.{npc.level})</span>
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {getEquipmentInRoom(selectedRoom.id).length > 0 && (
                <div className="mb-3">
                  <strong className="text-[#3498db] text-xs">Items:</strong>
                  <div className="mt-1">
                    {getEquipmentInRoom(selectedRoom.id).map(item => (
                      <div key={item.id} className="p-1 bg-[#16213e] rounded mb-1 text-xs text-white">{item.name}</div>
                    ))}
                  </div>
                </div>
              )}

              <div className="mb-3">
                <strong className="text-[#6c5ce7] text-xs">Exits:</strong>
                <div className="mt-1">
                  {['north', 'south', 'east', 'west', 'up', 'down'].map(dir => {
                    const targetId = selectedRoom.exits?.[dir]
                    const targetRoom = rooms.find(r => r.id === targetId)
                    const isZExit = dir === 'up' || dir === 'down'

                    if (targetId && targetRoom) {
                      return (
                        <div
                          key={dir}
                          onClick={() => {
                            if (isZExit) setCurrentZLevel(zLevels.get(targetId) || 0)
                            setSelectedRoom(targetRoom)
                          }}
                          className={`p-1 my-1 rounded cursor-pointer text-xs ${
                            isZExit
                              ? dir === 'up'
                                ? 'bg-[rgba(225,112,85,0.2)] border border-[#e17055]'
                                : 'bg-[rgba(116,185,255,0.2)] border border-[#74b9ff]'
                              : 'bg-[#16213e] border-none'
                          }`}
                        >
                          <strong>{dir}</strong> → {targetRoom.name}
                          {isZExit && <span className="text-[#666] ml-1 text-[10px]">(z={zLevels.get(targetId) || 0})</span>}
                        </div>
                      )
                    } else if (targetId) {
                      return (
                        <div key={dir} className="p-1 my-1 rounded text-xs bg-[#16213e] border-none">
                          <strong>{dir}</strong> → <span className="text-[#888]">Room #{targetId}</span>
                        </div>
                      )
                    } else {
                      return (
                        <div key={dir} className="flex items-center gap-2 my-1">
                          <div className="flex-1 p-1 rounded text-xs bg-[#0d0d15] border border-[#333] text-[#555]">
                            <strong>{dir}</strong> → none
                          </div>
                          <button
                            onClick={() => handleCreateRoom(selectedRoom, dir)}
                            disabled={creating}
                            className="px-2 py-1 bg-[#27ae60] border-none rounded text-white text-xs cursor-pointer hover:bg-[#2ecc71] disabled:opacity-50"
                            title={`Create room to the ${dir}`}
                          >
                            +
                          </button>
                        </div>
                      )
                    }
                  })}
                </div>
              </div>
            </div>
            <div className="p-3 border-t border-[#333] flex gap-2">
              <button
                onClick={() => handleEditRoom(selectedRoom)}
                className="flex-1 p-2 bg-[#6c5ce7] border-none rounded text-white cursor-pointer"
              >
                Edit Room
              </button>
              <button
                onClick={() => {
                  if (confirmDelete === selectedRoom.id) {
                    handleDeleteRoom(selectedRoom.id)
                  } else {
                    setConfirmDelete(selectedRoom.id)
                  }
                }}
                className={`flex-1 p-2 border-none rounded text-white cursor-pointer ${
                  confirmDelete === selectedRoom.id
                    ? 'bg-[#e67e22] hover:bg-[#d35400]'
                    : 'bg-[#c0392b] hover:bg-[#e74c3c]'
                }`}
              >
                {confirmDelete === selectedRoom.id ? 'Confirm Delete?' : 'Delete Room'}
              </button>
            </div>
          </>
        ) : (
          <div className="flex-1 flex flex-col justify-center items-center text-[#666] text-center p-4">
            <p className="m-0 mb-2">Click a room to see details</p>
            <p className="text-xs m-0 text-[#555]">
              👥 NPCs in room<br/>
              📦 Items on ground<br/>
              ▲/▼ Stairs
            </p>
          </div>
        )}
      </div>
    </div>
  )
}