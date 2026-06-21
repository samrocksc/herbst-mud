import { DIRECTION_OFFSETS } from "../components/map/DirectionUtils";
import { GRID, ORPHAN_COLS } from "../components/map/constants";
import type { Room } from "../components/map/types";

type Position = { x: number; y: number };
type Component = {
  rooms: readonly Room[]
  start: Room
};

function roundToGrid(n: number): number {
  return Math.round(n / GRID) * GRID;
}

function isOnFloor(r: Room, zLevels: Map<number, number>, currentZLevel: number): boolean {
  return (zLevels.get(r.id) ?? 0) === currentZLevel;
}

function getFloorExitPairs(
  room: Room,
  roomsById: Map<number, Room>,
  zLevels: Map<number, number>,
  currentZLevel: number,
): readonly { readonly dir: string; readonly tid: number }[] {
  return Object.entries(room.exits || {})
    .filter(([, tid]) => tid)
    .filter(([dir]) => dir !== "up" && dir !== "down")
    .filter(([, tid]) => {
      const target = roomsById.get(tid);
      return target != null && isOnFloor(target, zLevels, currentZLevel);
    })
    .map(([dir, tid]) => ({ dir, tid }));
}

function getComponent(
  startId: number,
  roomsById: Map<number, Room>,
  zLevels: Map<number, number>,
  currentZLevel: number,
  visited: ReadonlySet<number>,
): { members: readonly Room[]; visited: ReadonlySet<number> } {
  if (visited.has(startId)) {
    return { members: [], visited };
  }
  const room = roomsById.get(startId);
  if (room == null || !isOnFloor(room, zLevels, currentZLevel)) {
    return { members: [], visited };
  }
  const newVisited = new Set([...visited, startId]);
  const neighbors = getFloorExitPairs(room, roomsById, zLevels, currentZLevel);

  const result = neighbors.reduce<{ members: readonly Room[]; visited: ReadonlySet<number> }>(
    (acc, { tid }) => {
      const child = getComponent(tid, roomsById, zLevels, currentZLevel, acc.visited);
      return { members: [...acc.members, ...child.members], visited: child.visited };
    },
    { members: [], visited: newVisited },
  );

  return { members: [room, ...result.members], visited: result.visited };
}

function collectComponents(
  remaining: readonly Room[],
  roomsById: Map<number, Room>,
  zLevels: Map<number, number>,
  currentZLevel: number,
  visited: ReadonlySet<number>,
  acc: readonly Component[],
): readonly Component[] {
  if (remaining.length === 0) {
    return acc;
  }
  const [first, ...rest] = remaining;
  if (visited.has(first.id)) {
    return collectComponents(rest, roomsById, zLevels, currentZLevel, visited, acc);
  }
  const { members, visited: newVisited } = getComponent(first.id, roomsById, zLevels, currentZLevel, visited);
  const start = members.reduce(
    (best, r) => {
      const bestY = best.posY ?? Number.POSITIVE_INFINITY;
      const roomY = r.posY ?? Number.POSITIVE_INFINITY;
      if (roomY < bestY) return r;
      if (roomY === bestY && (r.posX ?? 0) < (best.posX ?? 0)) return r;
      return best;
    },
    members[0],
  );
  return collectComponents(rest, roomsById, zLevels, currentZLevel, newVisited, [...acc, { rooms: members, start }]);
}

function placeComponent(
  startId: number,
  x: number,
  y: number,
  placed: ReadonlyMap<number, Position>,
  roomsById: Map<number, Room>,
  zLevels: Map<number, number>,
  currentZLevel: number,
): ReadonlyMap<number, Position> {
  if (placed.has(startId)) {
    return placed;
  }
  const room = roomsById.get(startId);
  if (room == null || !isOnFloor(room, zLevels, currentZLevel)) {
    return placed;
  }
  const sx = roundToGrid(x);
  const sy = roundToGrid(y);
  const nextPlaced = new Map([...placed, [startId, { x: sx, y: sy }]]);
  const exits = getFloorExitPairs(room, roomsById, zLevels, currentZLevel);

  return exits.reduce<ReadonlyMap<number, Position>>(
    (acc, { dir, tid }) => {
      const off = DIRECTION_OFFSETS[dir] ?? { dx: 150, dy: 0 };
      return placeComponent(tid, sx + off.dx, sy + off.dy, acc, roomsById, zLevels, currentZLevel);
    },
    nextPlaced,
  );
}

function findComponentByRoom(components: readonly Component[], room: Room): Component | undefined {
  return components.find((c) => c.rooms.some((r) => r.id === room.id));
}

function pickMainComponent(components: readonly Component[], rootOnFloor: Room | undefined): Component {
  if (rootOnFloor != null) {
    const comp = findComponentByRoom(components, rootOnFloor);
    return comp ?? { rooms: [rootOnFloor], start: rootOnFloor };
  }
  return components.reduce((biggest, c) => (c.rooms.length > biggest.rooms.length ? c : biggest), components[0]);
}

function componentBottom(placed: ReadonlyMap<number, Position>, ids: ReadonlySet<number>): number {
  const positions = Array.from(placed.entries())
    .filter(([id]) => ids.has(id))
    .map(([, p]) => p.y);
  return positions.length === 0 ? 0 : Math.max(...positions);
}

export function computeRealignUpdates(
  rooms: readonly Room[],
  currentZLevel: number,
  zLevels: Map<number, number>,
): Array<{ roomId: number; posX: number; posY: number }> {
  const floorRooms = rooms.filter((r) => isOnFloor(r, zLevels, currentZLevel));
  if (floorRooms.length === 0) {
    return [];
  }
  const roomsById = new Map(rooms.map((r) => [r.id, r]));
  const components = collectComponents(floorRooms, roomsById, zLevels, currentZLevel, new Set(), []);
  const rootOnFloor = floorRooms.find((r) => r.isRootRoom);
  const main = pickMainComponent(components, rootOnFloor);

  const placed = placeComponent(
    main.start.id,
    roundToGrid(main.start.posX ?? 400),
    roundToGrid(main.start.posY ?? 300),
    new Map(),
    roomsById,
    zLevels,
    currentZLevel,
  );

  const mainIds = new Set(main.rooms.map((r) => r.id));
  const mainBottomY = componentBottom(placed, mainIds);
  const orphans = components.filter((c) => c !== main);

  const finalPlaced = orphans.reduce((acc, c, idx) => {
    const col = idx % ORPHAN_COLS;
    const row = Math.floor(idx / ORPHAN_COLS);
    const orphanX = roundToGrid(400 + col * 300);
    const orphanY = roundToGrid(mainBottomY + 400 + row * 250);
    return placeComponent(c.start.id, orphanX, orphanY, acc, roomsById, zLevels, currentZLevel);
  }, placed);

  return Array.from(finalPlaced.entries()).reduce(
    (acc, [roomId, pos]) => {
      const room = roomsById.get(roomId);
      if (room != null && (room.posX !== pos.x || room.posY !== pos.y)) {
        return [...acc, { roomId, posX: pos.x, posY: pos.y }];
      }
      return acc;
    },
    [] as Array<{ roomId: number; posX: number; posY: number }>,
  );
}
