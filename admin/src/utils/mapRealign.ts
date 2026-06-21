import { getNoOverlapOffset, estimateNodeSize } from "../components/map/DirectionUtils";
import { GRID, ORPHAN_COLS } from "../components/map/constants";
import type { Room } from "../components/map/types";

type Position = { readonly x: number; readonly y: number };
type Box = { readonly w: number; readonly h: number };
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
  const sourceBox = estimateNodeSize(room);

  return exits.reduce<ReadonlyMap<number, Position>>(
    (acc, { dir, tid }) => {
      const target = roomsById.get(tid);
      if (target == null) return acc;
      const targetBox = estimateNodeSize(target);
      const off = getNoOverlapOffset(dir, sourceBox, targetBox);
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

function resolveOverlaps(
  positions: ReadonlyMap<number, Position>,
  boxes: ReadonlyMap<number, Box>,
  margin = 20,
  maxIterations = 50,
): ReadonlyMap<number, Position> {
  const ids = Array.from(positions.keys());

  function step(remaining: number, current: ReadonlyMap<number, Position>): ReadonlyMap<number, Position> {
    if (remaining === 0) {
      return current;
    }

    const { result, moved } = ids.reduce(
      (acc, a, i) => ids.slice(i + 1).reduce((innerAcc, b) => {
        const pa = innerAcc.result.get(a)!;
        const pb = innerAcc.result.get(b)!;
        const boxA = boxes.get(a) ?? { w: 120, h: 65 };
        const boxB = boxes.get(b) ?? { w: 120, h: 65 };

        const halfW = (boxA.w + boxB.w) / 2 + margin;
        const halfH = (boxA.h + boxB.h) / 2 + margin;
        const overlapX = halfW - Math.abs(pa.x - pb.x);
        const overlapY = halfH - Math.abs(pa.y - pb.y);

        if (overlapX <= 0 || overlapY <= 0) {
          return innerAcc;
        }

        const dx = pa.x - pb.x;
        const dy = pa.y - pb.y;
        const useX = overlapX < overlapY;
        const push = useX ? overlapX / 2 : overlapY / 2;

        const nextA = useX
          ? { x: pa.x + (dx > 0 ? push : -push), y: pa.y }
          : { x: pa.x, y: pa.y + (dy > 0 ? push : -push) };
        const nextB = useX
          ? { x: pb.x + (dx > 0 ? -push : push), y: pb.y }
          : { x: pb.x, y: pb.y + (dy > 0 ? -push : push) };

        return {
          result: new Map([...innerAcc.result, [a, nextA], [b, nextB]]),
          moved: true,
        };
      }, acc),
      { result: current, moved: false } as { result: ReadonlyMap<number, Position>; moved: boolean },
    );

    if (!moved) {
      return result;
    }
    return step(remaining - 1, result);
  }

  return step(maxIterations, positions);
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

  const withOrphans = orphans.reduce((acc, c, idx) => {
    const col = idx % ORPHAN_COLS;
    const row = Math.floor(idx / ORPHAN_COLS);
    const orphanX = roundToGrid(400 + col * 300);
    const orphanY = roundToGrid(mainBottomY + 400 + row * 250);
    return placeComponent(c.start.id, orphanX, orphanY, acc, roomsById, zLevels, currentZLevel);
  }, placed);

  const boxes = new Map(Array.from(withOrphans.keys()).map((id) => [id, estimateNodeSize(roomsById.get(id)!)]));
  const resolved = resolveOverlaps(withOrphans, boxes);
  const finalPlaced = new Map(
    Array.from(resolved.entries()).map(([id, p]) => [id, { x: roundToGrid(p.x), y: roundToGrid(p.y) }]),
  );

  const updates = Array.from(finalPlaced.entries()).reduce(
    (acc, [roomId, pos]) => {
      const room = roomsById.get(roomId);
      if (room != null && (room.posX !== pos.x || room.posY !== pos.y)) {
        return [...acc, { roomId, posX: pos.x, posY: pos.y }];
      }
      return acc;
    },
    [] as Array<{ roomId: number; posX: number; posY: number }>,
  );
  return updates;
}
