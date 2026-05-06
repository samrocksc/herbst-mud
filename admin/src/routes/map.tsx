import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { logError } from '../utils/log'
import { useEffect, useState, useCallback, useMemo, useRef } from 'react'
import { apiGet, apiPost, apiPut, apiDelete } from '../utils/apiFetch'
import { MapSidebar } from '../components/map/MapSidebar'
import { MapToolbar } from '../components/map/MapToolbar'
import { RoomNode } from '../components/map/RoomNode'
import { ExitLines, resolveOverlaps } from '../components/map/ExitLines'
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

const API = `${window.location.origin}`
const GRID = 20
const ORPHAN_COLS = 5

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
  const [currentZLevel, setCurrentZLevel] = useState(0)
  const [saving, setSaving] = useState(false)
  const [creating, setCreating] = useState(false)
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [isDragging, setIsDragging] = useState(false)
  const [editForm, setEditForm] = useState({ name: '', description: '', exits: {} as Record<string, string> })
  const [newRoomForm, setNewRoomForm] = useState({ name: '', description: '' })
  const [pendingExit, setPendingExit] = useState<{ room: Room; dir: string } | null>(null)

  // Keep pre-drag positions for rollback on PUT failure
  const dragSnapshot = useRef<Map<number, { x: number; y: number }>>(new Map())

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) navigate({ to: '/login' })
  }, [navigate])

  // ── Data loading ──────────────────────────────────────────────────────────

  useEffect(() => {
    Promise.all([
      apiGet<Room[]>(`${API}/rooms`),
      apiGet<{ npcs: NPC[] }>(`${API}/npcs`),
    ])
      .then(([roomsData, npcsData]) => {
        setRooms(roomsData)
        setNpcs(npcsData.npcs || [])
        setLoading(false)
      })
      .catch(err => {
        setError(err instanceof Error ? err.message : String(err))
        setLoading(false)
      })
  }, [])

  useEffect(() => {
    if (!selectedRoom) return
    const roomId = selectedRoom.id
    apiGet<Equipment[]>(`${API}/rooms/${roomId}/equipment`)
      .then(data => setRoomEquipment(prev => ({ ...prev, [roomId]: data })))
      .catch(() => {})
  }, [selectedRoom?.id])

  // ── Z-levels ──────────────────────────────────────────────────────────────

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

  // ── Layout ────────────────────────────────────────────────────────────────

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

    // Orphan rooms: 5-column grid with row wrapping
    let orphanIdx = 0
    for (const room of rooms) {
      const roomZ = zLevels.get(room.id) || 0
      if (roomZ === currentZLevel && !positioned.has(room.id)) {
        const col = orphanIdx % ORPHAN_COLS
        const row = Math.floor(orphanIdx / ORPHAN_COLS)
        const x = room.posX != null ? room.posX : 800 + col * 180
        const y = room.posY != null ? room.posY : 300 + row * 120
        positions.set(room.id, { x, y })
        orphanIdx++
      }
    }

    return resolveOverlaps(positions, 50)
  }, [rooms, zLevels, currentZLevel])

  // ── Room editing ──────────────────────────────────────────────────────────

  const handleEditRoom = useCallback(
    (room: Room) => {
      setEditingRoom(room)
      setEditForm({
        name: room.name,
        description: room.description,
        exits: Object.fromEntries(
          Object.entries(room.exits || {}).map(([dir, id]) => [dir, String(id)])
        ),
      })
    },
    []
  )

  // ── Relayout ──────────────────────────────────────────────────────────────

  const handleRelayout = useCallback(() => {
    const current = nodePositions
    const clean = resolveOverlaps(new Map(current), 50)
    const updates: { roomId: number; posX: number; posY: number }[] = []

    for (const [roomId, pos] of clean) {
      // Snap to grid after force resolution
      const sx = Math.round(pos.x / GRID) * GRID
      const sy = Math.round(pos.y / GRID) * GRID
      const room = rooms.find((r) => r.id === roomId)
      if (room && (room.posX !== sx || room.posY !== sy)) {
        updates.push({ roomId, posX: sx, posY: sy })
      }
    }
    if (updates.length === 0) return

    setRooms((prev) =>
      prev.map((r) => {
        const u = updates.find((u) => u.roomId === r.id)
        return u ? { ...r, posX: u.posX, posY: u.posY } : r
      }),
    )

    for (const { roomId, posX, posY } of updates) {
      apiPut(`${API}/rooms/${roomId}`, { posX, posY }).catch((err) =>
        logError('Relayout save:', err),
      )
    }
  }, [nodePositions, rooms])

  // ── Save room ─────────────────────────────────────────────────────────────

  const handleSaveRoom = useCallback(async () => {
    if (!editingRoom) return
    setSaving(true)
    try {
      const exits: Record<string, number> = {}
      for (const [dir, id] of Object.entries(editForm.exits)) {
        const numId = parseInt(id, 10)
        if (!isNaN(numId)) exits[dir] = numId
      }

      await apiPut(`${API}/rooms/${editingRoom.id}`, {
        name: editForm.name,
        description: editForm.description,
        exits,
        version: editingRoom.version,
      })

      const oldExits = editingRoom.exits || {}

      for (const dir of ALL_DIRECTIONS) {
        const newTargetId = exits[dir]
        const oldTargetId = oldExits[dir]
        const reverseDir = OPPOSITE_DIR[dir]

        if (newTargetId && newTargetId !== oldTargetId) {
          const targetRoom = rooms.find(r => r.id === newTargetId)
          if (targetRoom && reverseDir) {
            const updatedTargetExits = { ...targetRoom.exits, [reverseDir]: editingRoom.id }
            await apiPut(`${API}/rooms/${newTargetId}`, {
              ...targetRoom,
              exits: updatedTargetExits,
              version: targetRoom.version,
            })
          }
        }

        if (oldTargetId && oldTargetId !== newTargetId) {
          const oldTargetRoom = rooms.find(r => r.id === oldTargetId)
          if (oldTargetRoom && reverseDir && oldTargetRoom.exits?.[reverseDir] === editingRoom.id) {
            const { [reverseDir]: _, ...remainingExits } = oldTargetRoom.exits
            await apiPut(`${API}/rooms/${oldTargetId}`, {
              ...oldTargetRoom,
              exits: remainingExits || {},
              version: oldTargetRoom.version,
            })
          }
        }
      }

      const roomsData = await apiGet<Room[]>(`${API}/rooms`)
      setRooms(roomsData)
      setEditingRoom(null)
      setSelectedRoom(null)
    } catch (err) {
      logError('Save error:', err)
      alert('Failed to save room')
    } finally {
      setSaving(false)
    }
  }, [editingRoom, editForm, rooms])

  // ── Create room ───────────────────────────────────────────────────────────

  const handleCreateRoom = useCallback(
    async (fromRoom: Room | null, direction: string) => {
      setCreating(true)
      try {
        const newRoomName = newRoomForm.name || `New Room ${rooms.length + 1}`
        const newRoomDesc = newRoomForm.description || 'A newly created room.'

        const newRoom = await apiPost<Room>(`${API}/rooms`, {
          name: newRoomName,
          description: newRoomDesc,
          exits: {},
        })

        if (fromRoom && direction) {
          const updatedExits = { ...fromRoom.exits, [direction]: newRoom.id }
          await apiPut(`${API}/rooms/${fromRoom.id}`, { ...fromRoom, exits: updatedExits })
          const reverseDirection = OPPOSITE_DIR[direction]
          await apiPut(`${API}/rooms/${newRoom.id}`, {
            name: newRoomName,
            description: newRoomDesc,
            exits: { [reverseDirection]: fromRoom.id },
          })
        }

        const roomsData = await apiGet<Room[]>(`${API}/rooms`)
        setRooms(roomsData)
        setSelectedRoom(roomsData.find((r: Room) => r.id === newRoom.id) || null)
        setShowCreateModal(false)
        setNewRoomForm({ name: '', description: '' })
      } catch (err) {
        logError('Create room error:', err)
        alert('Failed to create room')
      } finally {
        setCreating(false)
      }
    },
    [rooms.length, newRoomForm]
  )

  const handleCreateStandaloneRoom = useCallback(async () => {
    if (pendingExit) {
      await handleCreateRoom(pendingExit.room, pendingExit.dir)
      setPendingExit(null)
    } else {
      await handleCreateRoom(null, '')
    }
  }, [handleCreateRoom, pendingExit])

  // ── Delete room ───────────────────────────────────────────────────────────

  const handleDeleteRoom = useCallback(async (roomId: number) => {
    try {
      await apiDelete(`${API}/rooms/${roomId}`)
      const roomsData = await apiGet<Room[]>(`${API}/rooms`)
      setRooms(roomsData)
      setSelectedRoom(null)
    } catch (err) {
      logError('Delete room error:', err)
      alert('Failed to delete room')
    }
  }, [])

  // ── Zoom ──────────────────────────────────────────────────────────────────

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

  // ── Pan / scroll ──────────────────────────────────────────────────────────

  const handleWheel = useCallback((e: React.WheelEvent) => {
    // Ctrl+Scroll = zoom
    if (e.ctrlKey) {
      e.preventDefault()
      handleZoom(e.deltaY < 0 ? 0.1 : -0.1)
      return
    }
    e.preventDefault()
    const dx = e.shiftKey ? e.deltaY : e.deltaX
    const dy = e.shiftKey ? 0 : e.deltaY
    setPanOffset(prev => ({ x: prev.x - dx, y: prev.y - dy }))
  }, [handleZoom])

  // ── Drag ──────────────────────────────────────────────────────────────────

  const handleDragStart = useCallback((_roomId: number) => {
    // Snapshot all room positions before drag for rollback
    dragSnapshot.current = new Map(
      rooms.map(r => [r.id, { x: r.posX ?? 0, y: r.posY ?? 0 }])
    )
    setIsDragging(true)
  }, [rooms])

  const handleRoomDragEnd = useCallback(async (roomId: number, posX: number, posY: number) => {
    const snappedX = Math.round(posX / GRID) * GRID
    const snappedY = Math.round(posY / GRID) * GRID

    // Optimistic update
    setRooms(prev => prev.map(r =>
      r.id === roomId ? { ...r, posX: snappedX, posY: snappedY } : r
    ))
    setIsDragging(false)

    try {
      const room = rooms.find(r => r.id === roomId)
      if (!room) return
      await apiPut(`${API}/rooms/${roomId}`, {
        posX: snappedX,
        posY: snappedY,
        version: room.version,
      })
    } catch (err) {
      // Rollback: restore pre-drag positions
      const snap = dragSnapshot.current
      if (snap.size > 0) {
        setRooms(prev => prev.map(r => {
          const orig = snap.get(r.id)
          return orig ? { ...r, posX: orig.x, posY: orig.y } : r
        }))
      }
      logError('Position save error:', err)
    }
  }, [rooms])

  // ── Reset view ────────────────────────────────────────────────────────────

  const handleResetView = useCallback(() => {
    setZoom(1)
    setPanOffset({ x: 0, y: 0 })
  }, [])

  // ── Render ────────────────────────────────────────────────────────────────

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
        <MapToolbar
          currentZLevel={currentZLevel}
          zoom={zoom}
          onZoom={handleZoom}
          onResetView={handleResetView}
          onRelayout={handleRelayout}
        />

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
                  onDragStart={handleDragStart}
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
            onSelectRoom={setSelectedRoom}
            onEditRoom={handleEditRoom}
            onDeleteRoom={handleDeleteRoom}
            onAddRoom={(room, dir) => {
              setNewRoomForm({ name: '', description: '' })
              setPendingExit({ room, dir })
              setShowCreateModal(true)
            }}
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
        onClose={() => {
          setShowCreateModal(false)
          setPendingExit(null)
        }}
        newRoomForm={newRoomForm}
        setNewRoomForm={setNewRoomForm}
        onCreate={handleCreateStandaloneRoom}
        creating={creating}
      />
    </div>
  )
}
