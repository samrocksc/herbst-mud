/* eslint-disable functional/immutable-data, functional/no-let, functional/no-loop-statements */
import { useMemo } from "react";
import { DIRECTION_OFFSETS } from "../components/map/DirectionUtils";
import { GRID, ORPHAN_COLS, NODE_W, NODE_H, NODE_GAP } from "../components/map/constants";
import type { Room } from "../components/map/types";

type Position = Readonly<{ x: number; y: number }>;

const ANCHOR_X = 400;
const ANCHOR_Y = 300;

const CARDINAL_DIRS = new Set(["north", "south", "east", "west"]);

/**
 * Vertical/horizontal step that places the target room at the parent's
 * edge plus the target's full width/height, with a gap between them.
 * The size is constant (pinned NODE_W × NODE_H) so all rooms reserve
 * the same canvas area regardless of label length.
 */
function stepOffset(dir: keyof typeof DIRECTION_OFFSETS): { dx: number; dy: number } {
  if (CARDINAL_DIRS.has(dir)) {
    return dir === "east" || dir === "west"
      ? { dx: Math.sign(DIRECTION_OFFSETS[dir].dx) * (NODE_W + NODE_GAP), dy: 0 }
      : { dx: 0, dy: Math.sign(DIRECTION_OFFSETS[dir].dy) * (NODE_H + NODE_GAP) };
  }
  // Ordinal: 45° diagonal at magnitude NODE_W + NODE_GAP (assuming NODE_W ≥ NODE_H).
  const signX = dir === "northeast" || dir === "southeast" ? 1 : -1;
  const signY = dir === "northeast" || dir === "northwest" ? -1 : 1;
  const m = NODE_W + NODE_GAP;
  const c = Math.round(m / Math.SQRT2 / GRID) * GRID;
  return { dx: signX * c, dy: signY * c };
}

function roundToGrid(n: number): number {
  return Math.round(n / GRID) * GRID;
}

function findFloorAnchor(
  rooms: ReadonlyMap<number, Room>,
  zLevels: ReadonlyMap<number, number>,
  currentZLevel: number,
): Room | undefined {
  const onFloor = Array.from(rooms.values()).filter(
    (r) => (zLevels.get(r.id) ?? 0) === currentZLevel,
  );
  if (onFloor.length === 0) return undefined;

  if (currentZLevel === 0) {
    const root = onFloor.find((r) => r.isRootRoom);
    if (root) return root;
    const start = onFloor.find((r) => r.isStartingRoom);
    if (start) return start;
    return onFloor[0];
  }

  const startId = Array.from(rooms.values()).find((r: Room) => r.isRootRoom)?.id;
  if (startId == null) return onFloor[0];

  const targetDir = currentZLevel > 0 ? "up" : "down";
  const visited = new Set<number>();
  const queue: Array<{ id: number; depth: number }> = [{ id: startId, depth: 0 }];
  visited.add(startId);

  while (queue.length > 0) {
    const { id, depth } = queue.shift()!;
    const room = rooms.get(id);
    if (!room) continue;
    if ((zLevels.get(id) ?? 0) === currentZLevel) return room;

    for (const [dir, targetId] of Object.entries(room.exits || {})) {
      if (!targetId || visited.has(targetId)) continue;
      visited.add(targetId);
      queue.push({ id: targetId, depth: depth + 1 });
      if (dir === targetDir && (zLevels.get(targetId) ?? 0) === currentZLevel) {
        return rooms.get(targetId);
      }
    }
  }
  return onFloor[0];
}

function isCardinalOrOrdinal(dir: string): boolean {
  return dir in DIRECTION_OFFSETS && dir !== "up" && dir !== "down";
}

/**
 * Pure BFS layout. The first encountered position is the final position —
 * revisits are no-ops. The anchor's position is sacred. All positions snap
 * to GRID. Step is sized per target (target size + 20px gap) so two rooms
 * with many exits never visually overlap, even when placed at the same
 * y or x by different BFS branches.
 */
export function computeBfsLayout(
  rooms: ReadonlyArray<Room>,
  zLevels: ReadonlyMap<number, number>,
  currentZLevel: number,
): Map<number, Position> {
  const roomsById = new Map(rooms.map((r) => [r.id, r]));
  const anchor = findFloorAnchor(roomsById, zLevels, currentZLevel);
  const positions = new Map<number, Position>();

  if (!anchor) return positions;

  const queue: Array<{ id: number; pos: Position }> = [
    { id: anchor.id, pos: { x: roundToGrid(ANCHOR_X), y: roundToGrid(ANCHOR_Y) } },
  ];
  positions.set(anchor.id, queue[0].pos);
  const visited = new Set<number>([anchor.id]);

  while (queue.length > 0) {
    const { id, pos } = queue.shift()!;
    const room = roomsById.get(id);
    if (!room) continue;

    for (const [dir, targetId] of Object.entries(room.exits || {})) {
      if (!targetId || !isCardinalOrOrdinal(dir)) continue;
      const target = roomsById.get(targetId);
      if (!target) continue;
      if ((zLevels.get(targetId) ?? 0) !== currentZLevel) continue;
      // First-encountered position is final. Revisits are no-ops so the
      // anchor and other cycle nodes stay readable.
      if (positions.has(targetId)) continue;

      const off = stepOffset(dir as keyof typeof DIRECTION_OFFSETS);
      const baseX = pos.x + off.dx;
      const baseY = pos.y + off.dy;

      const finalPos = {
        x: roundToGrid(baseX),
        y: roundToGrid(baseY),
      };
      positions.set(targetId, finalPos);
      visited.add(targetId);
      queue.push({ id: targetId, pos: finalPos });
    }
  }

  const orphanRooms = Array.from(roomsById.values()).filter(
    (r) => (zLevels.get(r.id) ?? 0) === currentZLevel && !positions.has(r.id),
  );
  orphanRooms.forEach((room, idx) => {
    const col = idx % ORPHAN_COLS;
    const row = Math.floor(idx / ORPHAN_COLS);
    const baseY = (() => {
      let maxY = ANCHOR_Y;
      positions.forEach((p) => {
        if (p.y > maxY) maxY = p.y;
      });
      return maxY + NODE_W;
    })();
    positions.set(room.id, {
      x: roundToGrid(ANCHOR_X + col * (NODE_W + GRID * 3)),
      y: roundToGrid(baseY + row * NODE_W),
    });
  });

  return positions;
}

export function useNodeLayout(rooms: Room[], currentZLevel: number) {
  const zLevels = useMemo(() => {
    const roomsById = new Map(rooms.map((r) => [r.id, r]));

    const assign = (id: number, z: number, visited: ReadonlySet<number>, zMap: Map<number, number>): Map<number, number> => {
      if (visited.has(id)) return zMap;
      const room = roomsById.get(id);
      if (!room) return zMap;

      const newVisited = new Set(visited).add(id);
      const newZMap = new Map(zMap).set(id, room.posZ != null ? room.posZ : z);

      const exits = room.exits || {};
      return Object.entries(exits).reduce((acc, [dir, targetId]) => {
        if (!targetId) return acc;
        const tz = dir === "up" ? z + 1 : dir === "down" ? z - 1 : z;
        return assign(targetId, tz, newVisited, acc);
      }, newZMap);
    };

    const start = rooms.find((r) => r.isRootRoom) || rooms.find((r) => r.isStartingRoom) || rooms[0];
    const initialMap = start ? assign(start.id, 0, new Set(), new Map()) : new Map<number, number>();

    return rooms.reduce(
      (acc, r) => (acc.has(r.id) ? acc : new Map(acc).set(r.id, r.posZ ?? 0)),
      initialMap,
    );
  }, [rooms]);

  const nodePositions = useMemo(
    () => computeBfsLayout(rooms, zLevels, currentZLevel),
    [rooms, zLevels, currentZLevel],
  );

  return { zLevels, nodePositions };
}
