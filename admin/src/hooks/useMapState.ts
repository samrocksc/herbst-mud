/* eslint-disable functional/prefer-immutable-types, functional/immutable-data, functional/no-loop-statements */
import { useState, useCallback, useRef, useEffect, useMemo } from "react";
import { useNavigate, useSearch } from "@tanstack/react-router";
import { useRooms } from "./useRooms";
import { useNPCs } from "./useNPCs";
import { useRoomEquipment } from "./useRoomEquipment";
import { useNodeLayout } from "./useNodeLayout";
import { computeRealignUpdates } from "../utils/mapRealign";
import { GRID, MIN_ZOOM, MAX_ZOOM, ZOOM_FINE_STEP } from "../components/map/constants";
import { DIRECTION_OFFSETS, OPPOSITE_DIR, ALL_DIRECTIONS } from "../components/map/DirectionUtils";
import type { Room } from "../components/map/types";

function angleBetweenRooms(a: { x: number; y: number }, b: { x: number; y: number }): number {
  return Math.atan2(b.y - a.y, b.x - a.x);
}

function computeTargetHandle(dir: string, sourceRoom: Room, targetRoom: Room): string {
  // Prefer the target room's reciprocal exit direction if it points back at the source.
  const reciprocal = Object.entries(targetRoom.exits || {}).find(([, tid]) => tid === sourceRoom.id)?.[0];
  if (reciprocal) return reciprocal;

  const orthogonal = ["north", "south", "east", "west"];
  if (orthogonal.includes(dir)) {
    return OPPOSITE_DIR[dir] ?? "south";
  }

  const angle = angleBetweenRooms(
    { x: sourceRoom.posX ?? 0, y: sourceRoom.posY ?? 0 },
    { x: targetRoom.posX ?? 0, y: targetRoom.posY ?? 0 },
  );
  const init = { nearest: dir, best: Infinity };
  return ALL_DIRECTIONS.reduce((acc, d) => {
    if (d === "up" || d === "down") return acc;
    const off = DIRECTION_OFFSETS[d];
    if (!off) return acc;
    const a = Math.atan2(off.dy, off.dx);
    const diff = Math.abs(((a - angle + Math.PI) % (2 * Math.PI)) - Math.PI);
    if (diff < acc.best) {
      return { nearest: d, best: diff };
    }
    return acc;
  }, init).nearest;
}

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

  // Update search params
  const updateSearchParams = useCallback((updates: { room?: number | null; floor?: number }) => {
    const nextSearch: Record<string, number> = {};
    if (updates.room !== undefined) {
      if (updates.room !== null) nextSearch.room = updates.room;
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

  // Node interaction handlers
  const handleSelectRoom = useCallback((room: Room | null) => {
    setSelectedRoom(room);
    setSidebarOpen(false);
    if (room) {
      setEditingRoom(null);
      updateSearchParams({ room: room.id });
    } else {
      updateSearchParams({ room: null });
    }
  }, [updateSearchParams]);

  const handleAddExitFromNode = useCallback((roomId: number, direction: string) => {
    const fromRoom = rooms.find(r => r.id === roomId);
    if (fromRoom) {
      requestAddRoom(fromRoom, direction);
    }
  }, [rooms]);

  // Room lifecycle
  const requestAddRoom = useCallback((fromRoom: Room, dir: string) => {
    setAddRoomModal({ open: true, fromRoom, dir });
  }, []);

  const cancelAddRoom = useCallback(() => {
    setAddRoomModal({ open: false, fromRoom: null, dir: null });
  }, []);

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
        handleSelectRoom(newRoom);
        setAddRoomModal({ open: false, fromRoom: null, dir: null });
      } catch {
        showToast("Failed to create room");
      } finally {
        setIsAddingRoom(false);
      }
    },
    [addRoomModal, createRoomAsync, createBidirectionalExit, zLevels, handleSelectRoom],
  );

  const showToast = useCallback((msg: string) => {
    setToast(msg);
    setTimeout(() => setToast(null), 3000);
  }, []);

  // Layout handler
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

  // Re-align rooms on the current floor to a grid based on exit graph
  const handleRealign = useCallback(() => {
    const updates = computeRealignUpdates(rooms, currentZLevel, zLevels);
    if (updates.length === 0) {
      showToast("Rooms already aligned");
      return;
    }
    updates.forEach(({ roomId, posX, posY }) => updateRoom({ id: roomId, update: { posX, posY } }));
    showToast(`Aligned ${updates.length} room(s) to grid`);
  }, [rooms, currentZLevel, zLevels, updateRoom, showToast]);

  const handleSetZLevel = useCallback((z: number) => {
    updateSearchParams({ floor: z });
  }, [updateSearchParams]);

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

  const handleResetView = useCallback(() => {
    setZoom(1);
    setPanOffset({ x: 0, y: 0 });
  }, []);

  // Floor management
  const handleAddFloor = useCallback(async () => {
    const sorted = Array.from(new Set(Array.from(zLevels.values()))).sort((a, b) => a - b);
    const maxZ = sorted[sorted.length - 1] ?? 0;
    const newZ = sorted.length === 0 ? 0 : maxZ + 1;

    if (newZ > 10 || newZ < -10) {
      showToast("Maximum floor range is -10 to +10");
      return;
    }

    try {
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

  const handleEditRoom = useCallback((room: Room) => {
    setEditingRoom(room);
  }, []);

  // Delete operations
  const [deleteFloorModalOpen, setDeleteFloorModalOpen] = useState(false);
  const [deleteRoomModalOpen, setDeleteRoomModalOpen] = useState(false);
  const [deletingRoomId, setDeletingRoomId] = useState<number | null>(null);
  const [isDeletingRoom, setIsDeletingRoom] = useState(false);
  const [deletingRoomDetails, setDeletingRoomDetails] = useState<{ affectedCharacterCount: number; orphanExitCount: number } | null>(null);

  const requestDeleteFloor = useCallback(() => {
    const roomsOnFloor = rooms.filter((r) => (zLevels.get(r.id) ?? 0) === currentZLevel);
    if (roomsOnFloor.length === 0) {
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

    const allNpcs = npcsQuery.data ?? [];
    const affectedCharacterCount = allNpcs.filter((n) => n.currentRoomId === roomId).length;

    const orphanExitCount = rooms.reduce((count, r) => {
      if (!r.exits) return count;
      return count + Object.values(r.exits).filter((targetId) => targetId === roomId).length;
    }, 0);

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

  // Build React Flow nodes for current z-level
  const reactFlowNodes = useMemo(() => {
    if (!rooms.length) return [];
    const getNPCsInRoom = (roomId: number) => npcsQuery.data?.filter(n => n.currentRoomId === roomId) || [];
    const getEquipmentInRoom = (_roomId: number) => equipmentQuery.data || [];

    const nodes: Array<{ id: string; position: { x: number; y: number }; type: string; data: Record<string, unknown> }> = [];
    for (const room of rooms) {
      const pos = nodePositions.get(room.id);
      if (!pos || (zLevels.get(room.id) ?? 0) !== currentZLevel) continue;

      nodes.push({
        id: `room-${room.id}`,
        position: { x: pos.x, y: pos.y },
        type: "default",
        data: {
          room,
          roomNpcs: getNPCsInRoom(room.id),
          roomItems: getEquipmentInRoom(room.id),
          rooms,
          onSelect: handleSelectRoom,
          onAddExit: handleAddExitFromNode,
        },
      });
    }
    return nodes;
  }, [rooms, nodePositions, zLevels, currentZLevel, npcsQuery.data, equipmentQuery.data]);

  // Build React Flow edges for current z-level
  const reactFlowEdges = useMemo(() => {
    if (!rooms.length) return [];
    const edges: Array<{ id: string; source: string; target: string; sourceHandle: string; targetHandle: string; type: string; data: Record<string, unknown>; style: Record<string, string | number> }> = [];
    const drawn = new Set<string>();

    for (const room of rooms) {
      const roomPos = nodePositions.get(room.id);
      if (!roomPos) continue;

      for (const [dir, targetId] of Object.entries(room.exits || {})) {
        if (dir === "up" || dir === "down") continue;
        const targetPos = nodePositions.get(targetId);
        if (!targetPos) continue;

        const [lo, hi] = room.id < targetId ? [room.id, targetId] : [targetId, room.id];
        const canon = `${lo}-${hi}`;
        if (drawn.has(canon)) continue;
        drawn.add(canon);

        // Use the source room's exit direction for the source side; prefer the
        // target room's reciprocal exit direction for the target side so the line
        // enters the same side a player would use to walk back.
        const targetRoom = rooms.find(r => r.id === targetId);
        if (!targetRoom) continue;
        const targetHandle = computeTargetHandle(dir, room, targetRoom);

        edges.push({
          id: `edge-${canon}`,
          source: `room-${room.id}`,
          target: `room-${targetId}`,
          sourceHandle: dir,
          targetHandle,
          type: "mapEdge",
          data: { dir, targetHandle },
          style: { stroke: "var(--color-success, #22c55e)", strokeWidth: 2 },
        });
      }
    }
    return edges;
  }, [rooms, nodePositions]);

  return {
    // State
    rooms, roomsLoading, selectedRoom, zoom, panOffset, currentZLevel,
    sidebarOpen, setSidebarOpen, editingRoom, setEditingRoom, toast, cleanupConfirmOpen, addRoomModal,
    isAddingRoom, deleteFloorModalOpen, deleteRoomModalOpen,
    deletingRoomId, isDeletingRoom, deletingRoomDetails,

    // React Flow data
    reactFlowNodes, reactFlowEdges,
    viewportRef, handleWheel, handleZoom, handleResetView,
    handleSelectRoom, handleAddExitFromNode,

    // Navigation
    navigate, updateSearchParams,

    // Handlers
    handleRelayout, handleRealign, handleSetZLevel, handleAddFloor, handleEditRoom,
    requestAddRoom, cancelAddRoom, confirmAddRoom,
    requestDeleteFloor, confirmDeleteFloor, cancelDeleteFloor,
    requestDeleteRoom, cancelDeleteRoom, confirmDeleteRoom,
    handleRequestCleanupOrphanExits, handleConfirmCleanupOrphanExits,
    handleCancelCleanupOrphanExits,

    // Data
    npcs: npcsQuery.data ?? [],
    roomEquipment: equipmentQuery.data ?? [],
    zLevels, nodePositions,
    updateRoom, createRoom, deleteRoom, deleteRoomAsync, isCreating, cleanupOrphanExits,
    showToast,
  };
}
