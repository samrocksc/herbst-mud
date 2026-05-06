import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { logError } from '../utils/log'
import { useEffect, useState, useCallback, useMemo, useRef } from 'react'
import { apiGet } from '../utils/apiFetch'
import { MapSidebar } from '../components/map/MapSidebar'
import { MapToolbar } from '../components/map/MapToolbar'
import { RoomNode } from '../components/map/RoomNode'
import { ExitLines, resolveOverlaps } from '../components/map/ExitLines'
import { DIRECTION_OFFSETS } from '../components/map/DirectionUtils'
import { MenuIcon } from '../components/icons/MenuIcon'
import { Button } from '../components/Button'
import type { Room, NPC, Equipment } from '../components/map/types'
import { useRooms } from '../hooks/useRooms'

export const Route = createFileRoute('/map')({
  component: MapBuilder,
})

const API = `${window.location.origin}`
const GRID = 20
const ORPHAN_COLS = 5

function MapBuilder() {
  const navigate = useNavigate()
  const { rooms, isLoading: roomsLoading, updateRoom } = useRooms()

  const [npcs, setNpcs] = useState<NPC[]>([])
  const [roomEquipment, setRoomEquipment] = useState<Record<number, Equipment[]>>({})
  const [error, setError] = useState<string | null>(null)
  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null)
  const [zoom, setZoom] = useState(1)
  const [panOffset, setPanOffset] = useState({ x: 0, y: 0 })
  const viewportRef = useRef<HTMLDivElement>(null)

  const currentZLevel = 0
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [isDragging, setIsDragging] = useState(false)

  const dragSnapshot = useRef<Map<number, { x: number; y: number }>>(new Map())

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) navigate({ to: '/login' })
  }, [navigate])

  useEffect(() => {
    apiGet<{ npcs: NPC[] }>(`${API}/npcs`)
      .then(npcsData => {
        setNpcs(npcsData.npcs || [])
      })
      .catch(err => {
        setError(err instanceof Error ? err.message : String(err))
      })
  }, [])

  useEffect(() => {
    if (!selectedRoom) return
    const roomId = selectedRoom.id
    apiGet<Equipment[]>(`${API}/rooms/${roomId}/equipment`)
      .then(data => setRoomEquipment(prev => ({ ...prev, [roomId]: data })))
      .catch(() => {})
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

  const handleRelayout = useCallback(() => {
    const current = nodePositions
    const clean = resolveOverlaps(new Map(current), 50)
    const updates: { roomId: number; posX: number; posY: number }[] = []

    for (const [roomId, pos] of clean) {
      const sx = Math.round(pos.x / GRID) * GRID
      const sy = Math.round(pos.y / GRID) * GRID
      const room = rooms.find((r) => r.id === roomId)
      if (room && (room.posX !== sx || room.posY !== sy)) {
        updates.push({ roomId, posX: sx, posY: sy })
      }
    }
    if (updates.length === 0) return

    for (const { roomId, posX, posY } of updates) {
      updateRoom({ id: roomId, update: { posX, posY } })
    }
  }, [nodePositions, rooms, updateRoom])

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

  const handleWheel = useCallback((e: React.WheelEvent) => {
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

  const handleDragStart = useCallback((_roomId: number) => {
    dragSnapshot.current = new Map(
      rooms.map(r => [r.id, { x: r.posX ?? 0, y: r.posY ?? 0 }])
    )
    setIsDragging(true)
  }, [rooms])

  const handleRoomDragEnd = useCallback(async (roomId: number, posX: number, posY: number) => {
    const snappedX = Math.round(posX / GRID) * GRID
    const snappedY = Math.round(posY / GRID) * GRID

    setIsDragging(false)

    try {
      const room = rooms.find(r => r.id === roomId)
      if (!room) return
      updateRoom({
        id: roomId,
        update: {
          posX: snappedX,
          posY: snappedY,
          version: room.version,
        }
      })
    } catch (err) {
      const snap = dragSnapshot.current
      if (snap.size > 0) {
        // The hook handles the state update, so we just need to trigger it
        // if the hook fails. For now, we rely on the hook's invalidateQueries.
      }
      logError('Position save error:', err)
    }
  }, [rooms, updateRoom])

  const handleResetView = useCallback(() => {
    setZoom(1)
    setPanOffset({ x: 0, y: 0 })
  }, [])

  if (roomsLoading) return <div className="p-8 text-text">Loading map...</div>
  if (error) return <div className="p-8 text-danger">Error: {error}</div>

  return (
    <div className="flex h-screen bg-surface">
      <Button
        variant="ghost"
        size="sm"
        onClick={() => setSidebarOpen(true)}
        aria-label="Open map sidebar"
        className="fixed top-3 left-3 z-50 p-2 bg-surface border border-border text-text-muted hover:bg-surface-muted hover:text-text lg:hidden"
      >
        <MenuIcon stroke="currentColor" />
      </Button>

      <div className={['lg:block lg:relative lg:inset-auto lg:z-auto', sidebarOpen ? 'block' : 'hidden'].join(' ')}>
        <div className="fixed inset-y-0 left-0 z-40 lg:static">
          <MapSidebar
            rooms={rooms}
            npcs={npcs}
            zLevels={zLevels}
            currentZLevel={currentZLevel}
            selectedRoom={selectedRoom}
            setCurrentZLevel={() => {}}
            setSelectedRoom={(room) => {
              setSelectedRoom(room)
              setSidebarOpen(false)
            }}
            setShowCreateModal={() => {}}
          />
        </div>
      </div>

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
    </div>
  )
}
