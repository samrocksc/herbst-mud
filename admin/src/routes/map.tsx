import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useEffect, useState, useCallback, useMemo } from 'react'
import { MapSidebar } from '../components/map/MapSidebar'
import { MapToolbar } from '../components/map/MapToolbar'
import { RoomNode } from '../components/map/RoomNode'
import { ExitLines } from '../components/map/ExitLines'
import { RoomDetailPanel } from '../components/map/RoomDetailPanel'
import { RoomEditor } from '../components/map/RoomEditor'
import { CreateRoomModal } from '../components/map/CreateRoomModal'
import { DIRECTION_OFFSETS, OPPOSITE_DIR, ALL_DIRECTIONS } from '../components/map/DirectionUtils'
import type { Room, NPC, Equipment } from '../components/map/types'

export const Route = createFileRoute('/map')({
  component: MapBuilder,
})

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
      fetch(`${window.location.origin}/rooms`, { signal: controller.signal }).then(res => res.json()),
      fetch(`${window.location.origin}/npcs`, { signal: controller.signal }).then(res => res.json()),
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

  useEffect(() => {
    if (!selectedRoom) return
    const roomId = selectedRoom.id
    const controller = new AbortController()
    fetch(`${window.location.origin}/rooms/${roomId}/equipment`, { signal: controller.signal })
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

  const getNPCsInRoom = useCallback((roomId: number) => npcs.filter(npc => npc.currentRoomId === roomId), [npcs])
  const getEquipmentInRoom = useCallback((roomId: number) => roomEquipment[roomId] || [], [roomEquipment])

  const nodePositions = useMemo(() => {
    const positions = new Map<number, { x: number; y: number }>()
    const positioned = new Set<number>()

    const positionRoom = (roomId: number, x: number, y: number) => {
      const roomZ = zLevels.get(roomId) || 0
      if (roomZ !== currentZLevel) return
      if (positioned.has(roomId)) return
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

    const startRoom = rooms.find(r => r.isStartingRoom) || rooms[0]
    if (startRoom) positionRoom(startRoom.id, 400, 300)

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

  const handleEditRoom = useCallback(
    (room: Room) => {
      setEditingRoom(room)
      setEditForm({
        name: room.name,
        description: room.description,
        exits: Object.fromEntries(Object.entries(room.exits || {}).map(([dir, id]) => [dir, String(id)])),
      })
    },
    []
  )

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

      const response = await fetch(`${window.location.origin}/rooms/${editingRoom.id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({ name: editForm.name, description: editForm.description, exits }),
      })
      if (!response.ok) throw new Error('Failed to save room')

      const oldExits = editingRoom.exits || {}

      for (const dir of ALL_DIRECTIONS) {
        const newTargetId = exits[dir]
        const oldTargetId = oldExits[dir]
        const reverseDir = OPPOSITE_DIR[dir]

        if (newTargetId && newTargetId !== oldTargetId) {
          const targetRoom = rooms.find(r => r.id === newTargetId)
          if (targetRoom && reverseDir) {
            const updatedTargetExits = { ...targetRoom.exits, [reverseDir]: editingRoom.id }
            await fetch(`${window.location.origin}/rooms/${newTargetId}`, {
              method: 'PUT',
              headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
              body: JSON.stringify({ ...targetRoom, exits: updatedTargetExits }),
            })
          }
        }

        if (oldTargetId && oldTargetId !== newTargetId) {
          const oldTargetRoom = rooms.find(r => r.id === oldTargetId)
          if (oldTargetRoom && reverseDir && oldTargetRoom.exits?.[reverseDir] === editingRoom.id) {
            const { [reverseDir]: _, ...remainingExits } = oldTargetRoom.exits
            await fetch(`${window.location.origin}/rooms/${oldTargetId}`, {
              method: 'PUT',
              headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
              body: JSON.stringify({ ...oldTargetRoom, exits: remainingExits || {} }),
            })
          }
        }
      }

      const roomsResponse = await fetch(`${window.location.origin}/rooms`)
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

  const handleCreateRoom = useCallback(
    async (fromRoom: Room | null, direction: string) => {
      setCreating(true)
      try {
        const token = localStorage.getItem('token')
        const newRoomName = newRoomForm.name || `New Room ${rooms.length + 1}`
        const newRoomDesc = newRoomForm.description || 'A newly created room.'

        const createResponse = await fetch(`${window.location.origin}/rooms`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
          body: JSON.stringify({ name: newRoomName, description: newRoomDesc, exits: {} }),
        })
        if (!createResponse.ok) throw new Error('Failed to create room')
        const newRoom = await createResponse.json()

        if (fromRoom && direction) {
          const updatedExits = { ...fromRoom.exits, [direction]: newRoom.id }
          await fetch(`${window.location.origin}/rooms/${fromRoom.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
            body: JSON.stringify({ ...fromRoom, exits: updatedExits }),
          })
          const reverseDirection = OPPOSITE_DIR[direction]
          await fetch(`${window.location.origin}/rooms/${newRoom.id}`, {
            method: 'PUT',
            headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
            body: JSON.stringify({
              name: newRoomName,
              description: newRoomDesc,
              exits: { [reverseDirection]: fromRoom.id },
            }),
          })
        }

        const roomsResponse = await fetch(`${window.location.origin}/rooms`)
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
    },
    [rooms.length, newRoomForm]
  )

  const handleCreateStandaloneRoom = useCallback(async () => {
    await handleCreateRoom(null, '')
  }, [handleCreateRoom])

  const handleDeleteRoom = useCallback(async (roomId: number) => {
    try {
      const token = localStorage.getItem('token')
      const response = await fetch(`${window.location.origin}/rooms/${roomId}`, {
        method: 'DELETE',
        headers: { Authorization: `Bearer ${token}` },
      })
      if (!response.ok) throw new Error('Failed to delete room')
      const roomsResponse = await fetch(`${window.location.origin}/rooms`)
      const roomsData = await roomsResponse.json()
      setRooms(roomsData)
      setSelectedRoom(null)
    } catch (err) {
      console.error('Delete room error:', err)
      alert('Failed to delete room')
    }
  }, [])

  if (loading) return <div className="p-8 text-text">Loading map...</div>
  if (error) return <div className="p-8 text-danger">Error: {error}</div>

  return (
    <div className="flex h-screen bg-surface">
      <MapSidebar
        rooms={rooms}
        npcs={npcs}
        zLevels={zLevels}
        currentZLevel={currentZLevel}
        selectedRoom={selectedRoom}
        setCurrentZLevel={setCurrentZLevel}
        setSelectedRoom={setSelectedRoom}
        setShowCreateModal={setShowCreateModal}
      />

      <div className="flex-1 overflow-hidden relative">
        <MapToolbar currentZLevel={currentZLevel} zoom={zoom} setZoom={setZoom} />

        <div className="mt-[50px] h-[calc(100%-50px)] overflow-auto p-6">
          <div
            className="relative w-[3000px] h-[3000px]"
            style={{ transform: `scale(${zoom})`, transformOrigin: 'top left' }}
          >
            <ExitLines rooms={rooms} nodePositions={nodePositions} />

            {rooms.map(room => {
              const pos = nodePositions.get(room.id)
              if (!pos) return null
              const isSelected = selectedRoom?.id === room.id
              const roomNpcs = getNPCsInRoom(room.id)
              const roomItems = getEquipmentInRoom(room.id)
              return (
                <RoomNode
                  key={room.id}
                  room={room}
                  pos={pos}
                  isSelected={isSelected}
                  roomNpcs={roomNpcs}
                  roomItems={roomItems}
                  onSelect={setSelectedRoom}
                />
              )
            })}
          </div>
        </div>
      </div>

      <div className="w-[300px] bg-surface-muted border-l border-border flex flex-col">
        {editingRoom ? (
          <RoomEditor
            editForm={editForm}
            setEditForm={setEditForm}
            onSave={handleSaveRoom}
            onCancel={() => setEditingRoom(null)}
            saving={saving}
          />
        ) : selectedRoom ? (
          <RoomDetailPanel
            selectedRoom={selectedRoom}
            rooms={rooms}
            zLevels={zLevels}
            npcs={npcs}
            roomEquipment={roomEquipment}
            onSelectRoom={setSelectedRoom}
            onEditRoom={handleEditRoom}
            onDeleteRoom={handleDeleteRoom}
          />
        ) : (
          <div className="flex-1 flex flex-col justify-center items-center text-text-muted text-center p-4">
            <p className="m-0 mb-2">Click a room to see details</p>
            <p className="text-xs m-0 text-text-muted">
              👥 NPCs in room
              <br />
              📦 Items on ground
              <br />
              ▲/▼ Stairs
            </p>
          </div>
        )}
      </div>

      <CreateRoomModal
        isOpen={showCreateModal}
        onClose={() => setShowCreateModal(false)}
        newRoomForm={newRoomForm}
        setNewRoomForm={setNewRoomForm}
        onCreate={handleCreateStandaloneRoom}
        creating={creating}
      />
    </div>
  )
}
