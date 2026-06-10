/* eslint-disable functional/prefer-immutable-types, functional/immutable-data, functional/no-loop-statements */
import { useState, useCallback, useRef, useEffect } from "react";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { useRooms } from "./useRooms";
import { useNPCs } from "./useNPCs";
import { useRoomEquipment } from "./useRoomEquipment";
import { useNodeLayout } from "./useNodeLayout";
import { GRID, MIN_ZOOM, MAX_ZOOM, ZOOM_FINE_STEP } from "../components/map/constants";
import { DIRECTION_OFFSETS } from "../components/map/DirectionUtils";
import type { Room } from "../components/map/types";

type MapSearch = {
  room?: number
  floor?: number
}

function parseSearch(raw: Record<string, unknown>): MapSearch {
  return {
    room: raw.room != null ? Number(raw.room) || undefined : undefined,
    floor: raw.floor != null ? Number(raw.floor) || undefined : undefined,
  };
}

export function useMapState() {
  const navigate = useNavigate();
  const rawSearch = useSearch({ from: "/map" }) as Record<string, unknown>;
  const search = parseSearch(rawSearch);
  const { rooms, isLoading: roomsLoading, updateRoom, createRoom, createRoomAsync, deleteRoom, deleteRoomAsync, isCreating, cleanupOrphanExits, createBidirectionalExit } = useRooms();
  const npcsQuery = useNPCs();

  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null);
  const [zoom, setZoom] = useState(1);
  const [panOffset, setPanOffset] = useState({ x: 0, y: 0 });
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [editingRoom, setEditingRoom] = useState<Room | null>(null);
  const [toast, setToast] = useState<string | null>(null);
  const [cleanupConfirmOpen, setCleanupConfirmOpen] = useState(false);
  const [addRoomModal, setAddRoomModal] = useState<{
    open: boolean;
    fromRoom: Room | null;
    dir: string | null;
  }>({ open: false, fromRoom: null, dir: null });
  const [isAddingRoom, setIsAddingRoom] = useState(false);

  const currentZLevel = search.floor ?? 0;
  const initialSyncDone = useRef(false);
  const syncRunning = useRef(false);

  const viewportRef = useRef<HTMLDivElement>(null);

  const equipmentQuery = useRoomEquipment(selectedRoom?.id ?? null);

  const { zLevels, nodePositions } = useNodeLayout(rooms, currentZLevel);

  // Sync selected room from URL search param on mount and when rooms load
  useEffect(() => {
    if (roomsLoading || !rooms.length || syncRunning.current) return;
    syncRunning.current = true;
    const roomId = search.room;
    if (roomId != null) {
      const room = rooms.find(r => r.id === roomId);
      if (room && (selectedRoom == null || selectedRoom.id !== roomId)) {
        setSelectedRoom(room);
      }
    }
    if (roomId == null && selectedRoom != null && initialSyncDone.current) {
      setSelectedRoom(null);
    }
    initialSyncDone.current = true;
    syncRunning.current = false;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [rooms, search.room]);

  const updateSearchParams = useCallback((updates: { room?: number | null; floor?: number }) => {
    // Build the full next search state explicitly. Avoids the previous
    // else-if chain that could leave stale params in the URL and was
    // reported to occasionally misroute the navigation away from /map.
    const nextSearch: Record<string, number> = {};

    if (updates.room !== undefined) {
      if (updates.room !== null) nextSearch.room = updates.room;
      // updates.room === null → omit key, clears it
    } else if (search.room != null) {
      nextSearch.room = search.room;
    }

    if (updates.floor != null) {
      nextSearch.floor = updates.floor;
    } else if (currentZLevel !== 0) {
      nextSearch.floor = currentZLevel;
    }

    navigate({ to: "/map", search: nextSearch, replace: true });
  }, [navigate, search.room, currentZLevel]);

  const handleSetZLevel = useCallback((z: number) => {
    updateSearchParams({ floor: z });
  }, [updateSearchParams]);

  const handleRelayout = useCallback(() => {
    const clean = nodePositions;
    const updates: { roomId: number; posX: number; posY: number }[] = [];
    for (const [roomId, pos] of clean) {
      const sx = Math.round(pos.x / GRID) * GRID;
      const sy = Math.round(pos.y / GRID) * GRID;
      const room = rooms.find(r => r.id === roomId);
      if (room && (room.posX !== sx || room.posY !== sy)) {
        updates.push({ roomId, posX: sx, posY: sy });
      }
    }
    if (updates.length === 0) return;
    for (const { roomId, posX, posY } of updates) {
      updateRoom({ id: roomId, update: { posX, posY } });
    }
  }, [nodePositions, rooms, updateRoom]);

  const handleZoom = useCallback((delta: number) => {
    const viewport = viewportRef.current;
    if (!viewport) return;
    setZoom(prev => {
      const next = Math.min(Math.max(prev + delta, MIN_ZOOM), MAX_ZOOM);
      if (next === prev) return prev;
      const cx = viewport.clientWidth / 2;
      const cy = viewport.clientHeight / 2;
      setPanOffset(p => ({
        x: cx - (cx - p.x) * (next / prev),
        y: cy - (cy - p.y) * (next / prev),
      }));
      return next;
    });
  }, []);

  const handleWheel = useCallback((e: WheelEvent) => {
    if (e.ctrlKey) {
      e.preventDefault();
      handleZoom(e.deltaY < 0 ? ZOOM_FINE_STEP : -ZOOM_FINE_STEP);
      return;
    }
    e.preventDefault();
    const dx = e.shiftKey ? e.deltaY : e.deltaX;
    const dy = e.shiftKey ? 0 : e.deltaY;
    setPanOffset(p => ({ x: p.x - dx, y: p.y - dy }));
  }, [handleZoom]);

  const handleDragStart = useCallback((_roomId: number) => {
    setIsDragging(true);
  }, []);

  const handleRoomDragEnd = useCallback((roomId: number, posX: number, posY: number) => {
    const snappedX = Math.round(posX / GRID) * GRID;
    const snappedY = Math.round(posY / GRID) * GRID;
    setIsDragging(false);
    const room = rooms.find(r => r.id === roomId);
    if (!room) return;
    updateRoom({ id: roomId, update: { posX: snappedX, posY: snappedY, version: room.version } });
  }, [rooms, updateRoom]);

  const handleResetView = useCallback(() => {
    setZoom(1);
    setPanOffset({ x: 0, y: 0 });
  }, []);

  const handleSelectRoom = useCallback((room: Room | null) => {
    setSelectedRoom(room);
    setSidebarOpen(false);
    if (room) setEditingRoom(null);
    updateSearchParams({ room: room?.id ?? null });
  }, [updateSearchParams]);

  const handleEditRoom = useCallback((room: Room) => {
    setEditingRoom(room);
  }, []);

  const showToast = useCallback((msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(null), 3000);
  }, []);

  // Opens the modal to add a room
  const requestAddRoom = useCallback((fromRoom: Room, dir: string) => {
    setAddRoomModal({ open: true, fromRoom, dir });
  }, []);

  // Closes the modal without creating a room
  const cancelAddRoom = useCallback(() => {
    setAddRoomModal({ open: false, fromRoom: null, dir: null });
  }, []);

  // Creates the room with user-provided name/description and selects it
  const confirmAddRoom = useCallback(
    async (input: { name: string; description: string }) => {
      const fromRoom = addRoomModal.fromRoom;
      const dir = addRoomModal.dir;
      if (!fromRoom || !dir) return;

      setIsAddingRoom(true);
      const offset = DIRECTION_OFFSETS[dir];
      const posX = offset ? fromRoom.posX! + offset.dx : (fromRoom.posX ?? 0);
      const posY = offset ? fromRoom.posY! + offset.dy : (fromRoom.posY ?? 0);
      const parentZ = zLevels.get(fromRoom.id) ?? 0;
      const posZ = dir === "up" ? parentZ + 1 : dir === "down" ? parentZ - 1 : parentZ;
      try {
        const newRoom = await createRoomAsync({
          name: input.name,
          description: input.description,
          isStartingRoom: false,
          isRootRoom: false,
          exits: {},
          posX,
          posY,
          posZ,
        });
        await createBidirectionalExit({
          roomId: fromRoom.id,
          direction: dir,
          targetRoomId: newRoom.id,
        });
        // Auto-select the new room after creation
        handleSelectRoom(newRoom);
        setAddRoomModal({ open: false, fromRoom: null, dir: null });
      } catch {
        showToast("Failed to create room");
      } finally {
        setIsAddingRoom(false);
      }
    },
    [addRoomModal, createRoomAsync, createBidirectionalExit, zLevels, handleSelectRoom, showToast]
  );

  // Opens the modal (non-async, just triggers modal)
  const handleAddRoom = useCallback((fromRoom: Room, dir: string) => {
    requestAddRoom(fromRoom, dir);
  }, [requestAddRoom]);

  const handleAddFloor = useCallback(async () => {
    const sorted = Array.from(new Set(Array.from(zLevels.values()))).sort((a, b) => a - b);
    const maxZ = sorted[sorted.length - 1] ?? 0;
    const newZ = sorted.length === 0 ? 0 : maxZ + 1;
    
    // Guard: don't allow going beyond ±10
    if (newZ > 10 || newZ < -10) {
      showToast("Maximum floor range is -10 to +10");
      return;
    }

    try {
      // Check if any room is already a root — if not, make this new room the root
      const hasRoot = rooms.some(r => r.isRootRoom);
      await createRoomAsync({
        name: `Floor ${newZ}`,
        description: `The starting room of floor ${newZ}.`,
        isStartingRoom: false,
        isRootRoom: !hasRoot,
        exits: {},
        posX: 0,
        posY: 0,
        posZ: newZ,
        atmosphere: "air",
        tags: [],
      });
      updateSearchParams({ floor: newZ });
    } catch {
      showToast("Failed to create starter room");
    }
  }, [updateSearchParams, zLevels, rooms, createRoomAsync, showToast]);

  const [deleteFloorModalOpen, setDeleteFloorModalOpen] = useState(false);
  const [deleteRoomModalOpen, setDeleteRoomModalOpen] = useState(false);
  const [deletingRoomId, setDeletingRoomId] = useState<number | null>(null);
  const [isDeletingRoom, setIsDeletingRoom] = useState(false);
  const [deletingRoomDetails, setDeletingRoomDetails] = useState<{ affectedCharacterCount: number; orphanExitCount: number } | null>(null);

  const requestDeleteFloor = useCallback(() => {
    const roomsOnFloor = rooms.filter((r) => (zLevels.get(r.id) ?? 0) === currentZLevel);
    if (roomsOnFloor.length === 0) {
      // Empty floor — just navigate away
      const remaining = Array.from(new Set(Array.from(zLevels.values()))).filter((z) => z !== currentZLevel).sort((a, b) => a - b);
      const fallback = remaining[0] ?? 0;
      updateSearchParams({ floor: fallback });
      return;
    }
    setDeleteFloorModalOpen(true);
  }, [rooms, zLevels, currentZLevel, updateSearchParams]);

  const confirmDeleteFloor = useCallback(() => {
    const roomsOnFloor = rooms.filter((r) => (zLevels.get(r.id) ?? 0) === currentZLevel);
    for (const r of roomsOnFloor) {
      deleteRoom(r.id);
    }
    setDeleteFloorModalOpen(false);
    showToast(`Deleted ${roomsOnFloor.length} room(s) from floor ${currentZLevel}`);
  }, [rooms, zLevels, currentZLevel, deleteRoom, showToast]);

  const cancelDeleteFloor = useCallback(() => {
    setDeleteFloorModalOpen(false);
  }, []);

  const requestDeleteRoom = useCallback((roomId: number) => {
    const room = rooms.find((r) => r.id === roomId);
    if (!room) return;

    // Count affected characters
    const allNpcs = npcsQuery.data ?? [];
    const affectedCharacterCount = allNpcs.filter((n) => n.currentRoomId === roomId).length;

    // Count orphan exits (exits in other rooms pointing to this room)
    let orphanExitCount = 0;
    for (const r of rooms) {
      if (r.exits) {
        for (const targetId of Object.values(r.exits)) {
          if (targetId === roomId) {
            orphanExitCount++;
          }
        }
      }
    }

    setDeletingRoomDetails({ affectedCharacterCount, orphanExitCount });
    setDeletingRoomId(roomId);
    setDeleteRoomModalOpen(true);
  }, [rooms, npcsQuery.data]);

  const cancelDeleteRoom = useCallback(() => {
    setDeleteRoomModalOpen(false);
    setDeletingRoomId(null);
    setIsDeletingRoom(false);
    setDeletingRoomDetails(null);
  }, []);

  const confirmDeleteRoom = useCallback(async () => {
    if (deletingRoomId == null) return;
    setIsDeletingRoom(true);
    try {
      await deleteRoomAsync(deletingRoomId);
      setSelectedRoom(null);
      setDeleteRoomModalOpen(false);
      setDeletingRoomId(null);
      setIsDeletingRoom(false);
      setDeletingRoomDetails(null);
      showToast("Room deleted");
    } catch {
      showToast("Failed to delete room");
      // Keep modal open so user can retry or cancel
    } finally {
      setIsDeletingRoom(false);
    }
  }, [deletingRoomId, deleteRoomAsync, setSelectedRoom, showToast]);

  const handleRequestCleanupOrphanExits = useCallback(() => {
    setCleanupConfirmOpen(true);
  }, []);

  const handleConfirmCleanupOrphanExits = useCallback(async () => {
    setCleanupConfirmOpen(false);
    try {
      await cleanupOrphanExits();
      showToast("Cleanup complete");
    } catch {
      showToast("Cleanup failed");
    }
  }, [cleanupOrphanExits, showToast]);

  const handleCancelCleanupOrphanExits = useCallback(() => {
    setCleanupConfirmOpen(false);
  }, []);

  return {
    rooms, roomsLoading, selectedRoom, setSelectedRoom: handleSelectRoom,
    zoom, panOffset, currentZLevel, setCurrentZLevel: handleSetZLevel,
    sidebarOpen, setSidebarOpen, isDragging,
    editingRoom, setEditingRoom,
    toast, showToast,
    cleanupConfirmOpen, handleRequestCleanupOrphanExits, handleConfirmCleanupOrphanExits, handleCancelCleanupOrphanExits,
    addRoomModal, requestAddRoom, cancelAddRoom, confirmAddRoom, isAddingRoom,
    viewportRef, handleWheel, handleZoom, handleResetView,
    handleRelayout, handleDragStart, handleRoomDragEnd, handleEditRoom,
    nodePositions, zLevels,
    npcs: npcsQuery.data ?? [],
    roomEquipment: equipmentQuery.data ?? [],
    updateRoom, createRoom, deleteRoom, deleteRoomAsync, isCreating, cleanupOrphanExits,
    handleAddRoom, handleAddFloor, navigate, handleSetZLevel,
    deleteFloorModalOpen, requestDeleteFloor, confirmDeleteFloor, cancelDeleteFloor,
    deleteRoomModalOpen, requestDeleteRoom, cancelDeleteRoom, confirmDeleteRoom, isDeletingRoom, deletingRoomDetails,
  };
}
