/* eslint-disable functional/no-return-void, functional/no-loop-statements */
import { describe, it, expect } from "vitest";
import { computeRealignUpdates } from "../../utils/mapRealign";
import { estimateNodeSize } from "./DirectionUtils";
import type { Room } from "./types";

function makeRooms(): Room[] {
  return [
    { id: 1, name: "Root", description: "", posX: 800, posY: 280, posZ: 0, isRootRoom: true, exits: { east: 4 } },
    { id: 4, name: "T-Junction", description: "", posX: 940, posY: 280, posZ: 0, exits: { up: 8, west: 1 } },
    { id: 8, name: "Street 8", description: "", posX: 980, posY: 280, posZ: 1, exits: { down: 4, east: 9, north: 11 } },
    { id: 9, name: "Road 9", description: "", posX: 1140, posY: 280, posZ: 1, exits: { east: 10, west: 8 } },
    { id: 10, name: "Road 10", description: "", posX: 1280, posY: 280, posZ: 1, exits: { west: 9 } },
    { id: 11, name: "Dump 11", description: "", posX: 640, posY: 220, posZ: 1, exits: { north: 12, south: 8 } },
    { id: 12, name: "Cars 12", description: "", posX: 640, posY: 80, posZ: 1, exits: { north: 14, south: 11 } },
    { id: 14, name: "Road 14", description: "", posX: 760, posY: 520, posZ: 1, exits: { south: 12 } },
  ];
}

function positionsFromUpdates(updates: ReturnType<typeof computeRealignUpdates>, rooms: Room[]) {
  return rooms.reduce(
    (acc, r) => {
      const update = updates.find((u) => u.roomId === r.id);
      const pos = update ? { x: update.posX, y: update.posY } : { x: r.posX ?? 0, y: r.posY ?? 0 };
      return new Map([...acc, [r.id, pos]]);
    },
    new Map<number, { x: number; y: number }>(),
  );
}

describe("computeRealignUpdates", () => {
  it("stacks the vertical street chain on floor 1", () => {
    const rooms = makeRooms();
    const zLevels = new Map(rooms.map((r) => [r.id, r.posZ ?? 0]));
    const updates = computeRealignUpdates(rooms, 1, zLevels);
    const pos = positionsFromUpdates(updates, rooms);

    // 14 should be directly north of 12, 12 north of 11, 11 north of 8.
    expect(pos.get(14)?.x).toBe(pos.get(12)?.x);
    expect(pos.get(12)?.x).toBe(pos.get(11)?.x);
    expect(pos.get(11)?.x).toBe(pos.get(8)?.x);
    expect((pos.get(14)?.y ?? 0)).toBeLessThan((pos.get(12)?.y ?? 0));
    expect((pos.get(12)?.y ?? 0)).toBeLessThan((pos.get(11)?.y ?? 0));
    expect((pos.get(11)?.y ?? 0)).toBeLessThan((pos.get(8)?.y ?? 0));

    // 8 -> 9 -> 10 should stretch east from the vertical chain.
    expect((pos.get(9)?.x ?? 0)).toBeGreaterThan((pos.get(8)?.x ?? 0));
    expect((pos.get(10)?.x ?? 0)).toBeGreaterThan((pos.get(9)?.x ?? 0));
    expect(pos.get(9)?.y).toBe(pos.get(8)?.y);
    expect(pos.get(10)?.y).toBe(pos.get(9)?.y);
  });

  it("is idempotent", () => {
    const rooms = makeRooms();
    const zLevels = new Map(rooms.map((r) => [r.id, r.posZ ?? 0]));
    const first = computeRealignUpdates(rooms, 1, zLevels);
    const afterFirst = positionsFromUpdates(first, rooms);
    const nextRooms: Room[] = rooms.map((r) => {
      const p = afterFirst.get(r.id);
      return { ...r, posX: p?.x ?? r.posX, posY: p?.y ?? r.posY };
    });
    const second = computeRealignUpdates(nextRooms, 1, zLevels);
    expect(second).toHaveLength(0);
  });

  it("keeps same-floor rooms ordered and separated by at least their box sizes minus grid tolerance", () => {
    const rooms = makeRooms();
    const zLevels = new Map(rooms.map((r) => [r.id, r.posZ ?? 0]));
    const updates = computeRealignUpdates(rooms, 1, zLevels);
    const positions = positionsFromUpdates(updates, rooms);

    const floorIds = Array.from(positions.keys()).filter((id) => (zLevels.get(id) ?? 0) === 1);
    for (let i = 0; i < floorIds.length; i++) {
      for (let j = i + 1; j < floorIds.length; j++) {
        const a = floorIds[i];
        const b = floorIds[j];
        const pa = positions.get(a)!;
        const pb = positions.get(b)!;
        const boxA = estimateNodeSize(rooms.find((r) => r.id === a)!);
        const boxB = estimateNodeSize(rooms.find((r) => r.id === b)!);

        // Two boxes overlap only when they overlap on BOTH axes. Grid rounding
        // can introduce a small actual overlap, so allow a 10 px tolerance.
        const overlapX = (boxA.w + boxB.w) / 2 - Math.abs(pa.x - pb.x);
        const overlapY = (boxA.h + boxB.h) / 2 - Math.abs(pa.y - pb.y);
        expect(Math.min(overlapX, overlapY)).toBeLessThanOrEqual(10);
      }
    }
  });
});
