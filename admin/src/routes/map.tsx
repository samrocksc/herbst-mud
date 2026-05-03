import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useEffect, useState, useCallback, useMemo, useRef } from 'react'
import { MapSidebar } from '../components/map/MapSidebar'
import { MapToolbar } from '../components/map/MapToolbar'
import { RoomNode } from '../components/map/RoomNode'
import { ExitLines } from '../components/map/ExitLines'
import { RoomDetailPanel } from '../components/map/RoomDetailPanel'
import { RoomEditor } from '../components/map/RoomEditor'
import { CreateRoomModal } from '../components/map/CreateRoomModal'
import { DIRECTION_OFFSETS, OPPOSITE_DIR, ALL_DIRECTIONS } from '../components/map/DirectionUtils'
import { MenuIcon } from '../components/icons/MenuIcon'
import { Button } from '../components/Button'
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
  const [panOffset, setPanOffset] = useState({ x: 0, y: 0 })
  const viewportRef = useRef<HTMLDivElement>(null)

  /**
   * Zoom toward the center of the viewport.
   * Adjusts panOffset so the world point under the viewport center
   * stays stationary when the scale changes.
   */
  const handleZoom = useCallback((delta: number) => {
    const viewport = viewportRef.current
    if (!viewport) return

    setZoom(prevZoom => {
      const nextZoom = Math.min(Math.max(prevZoom + delta, 0.5), 2)
      if (nextZoom === prevZoom) return prevZoom

      const cx = viewport.clientWidth / 2
      const cy = viewport.clientHeight / 2

      setPanOffset(prev => ({
        x: cx - (cx - prev.x) * (nextZoom / prevZoom),
        y: cy - (cy - prev.y) * (nextZoom / prevZoom),
      }))

      return nextZoom
    })
  }, [])
  const [currentZLevel, setCurrentZLevel] = useState(0)
  const [saving, setSaving] = useState(false)
  const [creating, setCreating] = useState(false)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [isDragging, setIsDragging] = useState(false)
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
          const targetRoom = rooms.find(r => r.id === targetId)
          // Use server-provided position for targets that have custom positions
          const targetX = targetRoom?.posX != null ? targetRoom.posX : x + offset.dx
          const targetY = targetRoom?.posY != null ? targetRoom.posY : y + offset.dy
          positionRoom(targetId, targetX, targetY)
        }
      }
    }

    const startRoom = rooms.find(r => r.isStartingRoom) || rooms[0]
    if (startRoom) {
      const startX = startRoom.posX != null ? startRoom.posX : 400
      const startY = startRoom.posY != null ? startRoom.posY : 300
      positionRoom(startRoom.id, startX, startY)
    }

    let orphanX = 800
    for (const room of rooms) {
      const roomZ = zLevels.get(room.id) || 0
      if (roomZ === currentZLevel && !positioned.has(room.id)) {
        const orphanPosX = room.posX != null ? room.posX : orphanX
        const orphanPosY = room.posY != null ? room.posY : 300
        positions.set(room.id, { x: orphanPosX, y: orphanPosY })
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

  const [conflictInfo, setConflictInfo] = useState<{ current: Room; yourVersion: number } | null>(null)

  const handleSaveRoom = useCallback(async () => {
    if (!editingRoom) return
    setSaving(true)
    setConflictInfo(null)
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
        body: JSON.stringify({ name: editForm.name, description: editForm.description, exits, version: editingRoom.version }),
      })

      if (response.status === 409) {
        const conflictData = await response.json()
        setConflictInfo({ current: conflictData.current, yourVersion: conflictData.yourVersion })
        setSaving(false)
        return
      }

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
              body: JSON.stringify({ ...targetRoom, exits: updatedTargetExits, version: targetRoom.version }),
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
              body: JSON.stringify({ ...oldTargetRoom, exits: remainingExits || {}, version: oldTargetRoom.version }),
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

  /** Discard local edits and reload the room from the server's current state. */
  const handleReloadRoom = useCallback(async () => {
    if (!conflictInfo) return
    try {
      const response = await fetch(`${window.location.origin}/rooms`)
      if (!response.ok) return
      const roomsData = await response.json()
      setRooms(roomsData)
      const reloaded = roomsData.find((r: Room) => r.id === conflictInfo.current.id)
      if (reloaded) handleEditRoom(reloaded)
    } catch (err) {
      console.error('Reload error:', err)
    }
    setConflictInfo(null)
  }, [conflictInfo, handleEditRoom])

  /** Force-save by using the server's current version number. */
  const handleOverwriteRoom = useCallback(async () => {
    if (!conflictInfo || !editingRoom) return
    setSaving(true)
    setConflictInfo(null)
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
        body: JSON.stringify({
          name: editForm.name,
          description: editForm.description,
          exits,
          version: conflictInfo.current.version ?? 0,
        }),
      })
      if (!response.ok) throw new Error('Failed to overwrite room')

      const roomsResponse = await fetch(`${window.location.origin}/rooms`)
      const roomsData = await roomsResponse.json()
      setRooms(roomsData)
      setEditingRoom(null)
      setSelectedRoom(null)
    } catch (err) {
      console.error('Overwrite error:', err)
      alert('Failed to overwrite room')
    } finally {
      setSaving(false)
    }
  }, [conflictInfo, editingRoom, editForm])

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

  /** Snap a value to the nearest grid step (e.g., 20px). */
  const snapToGrid = useCallback((value: number, step: number = 20) => {
    return Math.round(value / step) * step
  }, [])

  /** Persist a room's new position after dragging. */
  const handleRoomDragEnd = useCallback(async (roomId: number, posX: number, posY: number) => {
    const snappedX = snapToGrid(posX)
    const snappedY = snapToGrid(posY)

    setRooms(prev => prev.map(r =>
      r.id === roomId ? { ...r, posX: snappedX, posY: snappedY } : r
    ))

    setIsDragging(false)

    try {
      const token = localStorage.getItem('token')
      const room = rooms.find(r => r.id === roomId)
      if (!room) return
      await fetch(`${window.location.origin}/rooms/${roomId}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json', Authorization: `Bearer ${token}` },
        body: JSON.stringify({ posX: snappedX, posY: snappedY, version: room.version }),
      })
    } catch (err) {
      console.error('Position save error:', err)
    }
  }, [rooms, snapToGrid])

  const handleWheel = useCallback((e: React.WheelEvent) => {
    e.preventDefault()
    const dx = e.shiftKey ? e.deltaY : e.deltaX
    const dy = e.shiftKey ? 0 : e.deltaY
    setPanOffset(prev => ({ x: prev.x - dx, y: prev.y - dy }))
  }, [])

  if (loading) return <div className="p-8 text-text">Loading map...</div>
  if (error) return <div className="p-8 text-danger">Error: {error}</div>

  return (
    <div className="flex h-screen bg-surface">
      {/* Mobile hamburger */}
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setSidebarOpen(true)}
        aria-label="Open map sidebar"
        className="fixed top-3 left-3 z-50 p-2 bg-surface border border-border text-text-muted hover:bg-surface-muted hover:text-text lg:hidden"
      >
        <MenuIcon stroke="currentColor" />
      </Button>

      {/* Sidebar overlay (mobile) / always-visible (desktop) */}
      <div className={['lg:block lg:relative lg:inset-auto lg:z-auto', sidebarOpen ? 'block' : 'hidden'].join(' ')}>
        <div className="fixed inset-y-0 left-0 z-40 lg:static">
          <MapSidebar
            rooms={rooms}
            npcs={npcs}
            zLevels={zLevels}
            currentZLevel={currentZLevel}
            selectedRoom={selectedRoom}
            setCurrentZLevel={setCurrentZLevel}
            setSelectedRoom={(room) => {
              setSelectedRoom(room)
              setSidebarOpen(false)
            }}
            setShowCreateModal={setShowCreateModal}
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

      <div className="flex-1 overflow-hidden relative">
        <MapToolbar currentZLevel={currentZLevel} zoom={zoom} onZoom={handleZoom} />

        <div ref={viewportRef} className="mt-[50px] h-[calc(100%-50px)] overflow-hidden p-6" onWheel={handleWheel}>
          <div
            className="relative w-[3000px] h-[3000px]"
            style={{ transform: `translate(${panOffset.x}px, ${panOffset.y}px) scale(${zoom})`, transformOrigin: 'top left' }}
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
                  rooms={rooms}
                  onSelect={setSelectedRoom}
                  isDragging={isDragging}
                  onDragEnd={handleRoomDragEnd}
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
