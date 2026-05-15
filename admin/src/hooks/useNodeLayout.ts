/* eslint-disable functional/immutable-data */
 
import { useMemo } from "react";
import { DIRECTION_OFFSETS } from "../components/map/DirectionUtils";
import { ORPHAN_COLS } from "../components/map/constants";
import type { Room } from "../components/map/types";

export function useNodeLayout(rooms: Room[], currentZLevel: number) {
  const zLevels = useMemo(() => {
    const roomsById = new Map(rooms.map(r => [r.id, r]));

    const assign = (id: number, z: number, visited: ReadonlySet<number>, zMap: Map<number, number>): Map<number, number> => {
      if (visited.has(id)) return zMap;
      const room = roomsById.get(id);
      if (!room) return zMap;

      const newVisited = new Set(visited).add(id);
      const newZMap = new Map(zMap).set(id, room.posZ != null && room.posZ !== 0 ? room.posZ : z);

      const exits = room.exits || {};
      return Object.entries(exits).reduce((acc, [dir, targetId]) => {
        if (!targetId) return acc;
        const tz = dir === "up" ? z + 1 : dir === "down" ? z - 1 : z;
        return assign(targetId, tz, newVisited, acc);
      }, newZMap);
    };

    const start = rooms.find(r => r.isRootRoom) || rooms.find(r => r.isStartingRoom) || rooms[0];
    const initialMap = start ? assign(start.id, 0, new Set(), new Map()) : new Map<number, number>();

    return rooms.reduce((acc, r) => acc.has(r.id) ? acc : new Map(acc).set(r.id, 0), initialMap);
  }, [rooms]);

  const nodePositions = useMemo(() => {
    const positioned = new Set<number>();
    const positions = new Map<number, { x: number; y: number }>();

    const posRoom = (id: number, x: number, y: number, pos: Set<number>, posMap: Map<number, { x: number; y: number }>): Map<number, { x: number; y: number }> => {
      const rz = zLevels.get(id) || 0;
      if (rz !== currentZLevel) return posMap;
      if (pos.has(id)) return posMap;

      const newPos = new Set(pos).add(id);
      const newPosMap = new Map(posMap).set(id, { x, y });

      const room = rooms.find(r => r.id === id);
      if (!room) return newPosMap;

      return Object.entries(room.exits || {}).reduce((acc, [dir, tid]) => {
        if (!tid || dir === "up" || dir === "down") return acc;
        if (acc.has(tid)) return acc;

        const off = DIRECTION_OFFSETS[dir] || { dx: 150, dy: 0 };
        const tr = rooms.find(r => r.id === tid);
        const tx = tr?.posX != null ? tr.posX : x + off.dx;
        const ty = tr?.posY != null ? tr.posY : y + off.dy;
        return posRoom(tid, tx, ty, newPos, acc);
      }, newPosMap);
    };

    const start = rooms.find(r => r.isRootRoom) || rooms.find(r => r.isStartingRoom) || rooms[0];
    if (!start) return new Map();
    const updatedPositions = posRoom(start.id, start.posX ?? 400, start.posY ?? 300, positioned, positions);

    // Place orphan rooms using map instead of for loop
    const orphanRooms = rooms
      .filter(r => (zLevels.get(r.id) || 0) === currentZLevel && !updatedPositions.has(r.id))
      .map((r, idx) => ({ room: r, col: idx % ORPHAN_COLS, row: Math.floor(idx / ORPHAN_COLS) }));

    return orphanRooms.reduce((acc, { room, col, row }) => {
      const newAcc = new Map(acc);
      newAcc.set(room.id, { x: room.posX ?? 800 + col * 180, y: room.posY ?? 300 + row * 120 });
      return newAcc;
    }, updatedPositions);
  }, [rooms, zLevels, currentZLevel]);

  return { zLevels, nodePositions };
}