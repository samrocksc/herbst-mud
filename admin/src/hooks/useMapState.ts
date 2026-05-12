import { useState, useCallback, useRef, useEffect } from 'react'
import { useNavigate, useSearch } from '@tanstack/react-router'
import { useRooms } from './useRooms'
import { useNPCs } from './useNPCs'
import { useRoomEquipment } from './useRoomEquipment'
import { useNodeLayout } from './useNodeLayout'
import { GRID, MIN_ZOOM, MAX_ZOOM, ZOOM_FINE_STEP } from '../components/map/constants'
import { DIRECTION_OFFSETS } from '../components/map/DirectionUtils'
import type { Room } from '../components/map/types'

type MapSearch = {
  room?: number
  floor?: number
}

function parseSearch(raw: Record<string, unknown>): MapSearch {
  return {
    room: raw.room != null ? Number(raw.room) || undefined : undefined,
    floor: raw.floor != null ? Number(raw.floor) || undefined : undefined,
  }
}

export function useMapState() {
  const navigate = useNavigate()
  const rawSearch = useSearch({ from: '/map' }) as Record<string, unknown>
  const search = parseSearch(rawSearch)
  const { rooms, isLoading: roomsLoading, updateRoom, createRoom, createRoomAsync, deleteRoom, isCreating, cleanupOrphanExits, createBidirectionalExit } = useRooms()
  const npcsQuery = useNPCs()

  const [selectedRoom, setSelectedRoomState] = useState<Room | null>(null)
  const [zoom, setZoom] = useState(1)
  const [panOffset, setPanOffset] = useState({ x: 0, y: 0 })
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [isDragging, setIsDragging] = useState(false)
  const [editingRoom, setEditingRoom] = useState<Room | null>(null)
  const [toast, setToast] = useState<string | null>(null)

  const currentZLevel = search.floor ?? 0
  const initialSyncDone = useRef(false)

  const viewportRef = useRef<HTMLDivElement>(null)

  const equipmentQuery = useRoomEquipment(selectedRoom?.id ?? null)

  const { zLevels, nodePositions } = useNodeLayout(rooms, currentZLevel)

  // Sync selected room from URL search param on mount and when rooms load
  useEffect(() => {
    if (roomsLoading || !rooms.length) return
    const roomId = search.room
    if (roomId != null && (selectedRoom == null || selectedRoom.id !== roomId)) {
      const room = rooms.find(r => r.id === roomId)
      if (room) setSelectedRoomState(room)
    }
    if (roomId == null && selectedRoom != null && initialSyncDone.current) {
      setSelectedRoomState(null)
    }
    initialSyncDone.current = true
  }, [rooms, roomsLoading, search.room])

  const updateSearchParams = useCallback((updates: { room?: number | null; floor?: number }) => {
    const params: Record<string, number> = {}
    if (updates.room != null) params.room = updates.room
    else if (updates.room === null && search.room != null) { /* clear room */ }
    else if (search.room != null) params.room = search.room

    if (updates.floor != null) params.floor = updates.floor
    else if (currentZLevel !== 0) params.floor = currentZLevel

    navigate({ to: '/map', search: Object.keys(params).length > 0 ? params : undefined, replace: true })
  }, [navigate, search.room, currentZLevel])

  const handleSetZLevel = useCallback((z: number) => {
    updateSearchParams({ floor: z })
  }, [updateSearchParams])

  const handleRelayout = useCallback(() => {
    const clean = nodePositions
    const updates: { roomId: number; posX: number; posY: number }[] = []
    for (const [roomId, pos] of clean) {
      const sx = Math.round(pos.x / GRID) * GRID
      const sy = Math.round(pos.y / GRID) * GRID
      const room = rooms.find(r => r.id === roomId)
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
    setZoom(prev => {
      const next = Math.min(Math.max(prev + delta, MIN_ZOOM), MAX_ZOOM)
      if (next === prev) return prev
      const cx = viewport.clientWidth / 2
      const cy = viewport.clientHeight / 2
      setPanOffset(p => ({
        x: cx - (cx - p.x) * (next / prev),
        y: cy - (cy - p.y) * (next / prev),
      }))
      return next
    })
  }, [])

  const handleWheel = useCallback((e: WheelEvent) => {
    if (e.ctrlKey) {
      e.preventDefault()
      handleZoom(e.deltaY < 0 ? ZOOM_FINE_STEP : -ZOOM_FINE_STEP)
      return
    }
    e.preventDefault()
    const dx = e.shiftKey ? e.deltaY : e.deltaX
    const dy = e.shiftKey ? 0 : e.deltaY
    setPanOffset(p => ({ x: p.x - dx, y: p.y - dy }))
  }, [handleZoom])

  const handleDragStart = useCallback((_roomId: number) => {
    setIsDragging(true)
  }, [])

  const handleRoomDragEnd = useCallback((roomId: number, posX: number, posY: number) => {
    const snappedX = Math.round(posX / GRID) * GRID
    const snappedY = Math.round(posY / GRID) * GRID
    setIsDragging(false)
    const room = rooms.find(r => r.id === roomId)
    if (!room) return
    updateRoom({ id: roomId, update: { posX: snappedX, posY: snappedY, version: room.version } })
  }, [rooms, updateRoom])

  const handleResetView = useCallback(() => {
    setZoom(1)
    setPanOffset({ x: 0, y: 0 })
  }, [])

  const handleSelectRoom = useCallback((room: Room | null) => {
    setSelectedRoomState(room)
    setSidebarOpen(false)
    if (room) setEditingRoom(null)
    updateSearchParams({ room: room?.id ?? null })
  }, [updateSearchParams])

  const handleEditRoom = useCallback((room: Room) => {
    setEditingRoom(room)
  }, [])

  const showToast = useCallback((msg: string) => {
    setToast(msg)
    setTimeout(() => setToast(null), 3000)
  }, [])

  const handleAddRoom = useCallback(async (fromRoom: Room, dir: string) => {
    const offset = DIRECTION_OFFSETS[dir]
    const posX = offset ? fromRoom.posX! + offset.dx : (fromRoom.posX ?? 0)
    const posY = offset ? fromRoom.posY! + offset.dy : (fromRoom.posY ?? 0)
    const parentZ = zLevels.get(fromRoom.id) ?? 0
    const posZ = dir === 'up' ? parentZ + 1 : dir === 'down' ? parentZ - 1 : parentZ
    try {
      const newRoom = await createRoomAsync({
        name: 'New Room',
        description: 'A newly created room.',
        isStartingRoom: false,
        isRootRoom: false,
        exits: {},
        posX,
        posY,
        posZ,
      })
      await createBidirectionalExit({
        roomId: fromRoom.id,
        direction: dir,
        targetRoomId: newRoom.id,
      })
    } catch (err) {
      showToast('Failed to create room')
    }
  }, [createRoomAsync, createBidirectionalExit, showToast, zLevels])

  return {
    rooms, roomsLoading, selectedRoom, setSelectedRoom: handleSelectRoom,
    zoom, panOffset, currentZLevel, setCurrentZLevel: handleSetZLevel,
    sidebarOpen, setSidebarOpen, isDragging,
    editingRoom, setEditingRoom,
    toast, showToast,
    viewportRef, handleWheel, handleZoom, handleResetView,
    handleRelayout, handleDragStart, handleRoomDragEnd, handleEditRoom,
    nodePositions, zLevels,
    npcs: npcsQuery.data ?? [],
    roomEquipment: equipmentQuery.data ?? [],
    updateRoom, createRoom, deleteRoom, isCreating, cleanupOrphanExits,
    handleAddRoom, navigate,
  }
}