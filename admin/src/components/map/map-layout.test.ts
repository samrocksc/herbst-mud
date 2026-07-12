/**
 * Unit tests for the BFS layout algorithm.
 *
 * After dropping posx/posy, room positions are derived from the exits graph
 * via `computeBfsLayout` in `useNodeLayout`. The step is sized to the
 * TARGET room (target size + 20px gap) so two rooms with many exits
 * never visually overlap, even when placed at the same y or x by
 * different BFS branches.
 */
import { describe, it, expect } from "vitest";
import { computeBfsLayout } from "../../hooks/useNodeLayout";
import { GRID, NODE_W, NODE_H, NODE_GAP } from "./constants";
import type { Room } from "./types";

const STEP_CARDINAL_X = NODE_W + NODE_GAP;
const STEP_CARDINAL_Y = NODE_H + NODE_GAP;
const STEP_ORDINAL = Math.round((NODE_W + NODE_GAP) / Math.SQRT2 / GRID) * GRID;

const makeRoom = (id: number, exits: Record<string, number> = {}, opts: Partial<Room> = {}): Room => ({
  id,
  name: `Room ${id}`,
  description: "",
  exits,
  isRootRoom: opts.isRootRoom,
  isStartingRoom: opts.isStartingRoom,
  posZ: opts.posZ ?? 0,
});

const flatZLevels = (rooms: Room[]) =>
  new Map(rooms.map((r) => [r.id, r.posZ ?? 0]));

describe("computeBfsLayout", () => {
  it("places a target directly north of the source for a `north` exit", () => {
    const src = makeRoom(1, { north: 2 }, { isRootRoom: true });
    const tgt = makeRoom(2, { south: 1 });
    const layout = computeBfsLayout([src, tgt], flatZLevels([src, tgt]), 0);
    const p1 = layout.get(1)!;
    const p2 = layout.get(2)!;
    expect(p2.x).toBe(p1.x);
    expect(p1.y - p2.y).toBe(STEP_CARDINAL_Y);
  });

  it("places a target directly south of the source for a `south` exit", () => {
    const src = makeRoom(1, { south: 2 }, { isRootRoom: true });
    const tgt = makeRoom(2, { north: 1 });
    const layout = computeBfsLayout([src, tgt], flatZLevels([src, tgt]), 0);
    const p1 = layout.get(1)!;
    const p2 = layout.get(2)!;
    expect(p2.x).toBe(p1.x);
    expect(p2.y - p1.y).toBe(STEP_CARDINAL_Y);
  });

  it("places a target directly east of the source for an `east` exit", () => {
    const src = makeRoom(1, { east: 2 }, { isRootRoom: true });
    const tgt = makeRoom(2, { west: 1 });
    const layout = computeBfsLayout([src, tgt], flatZLevels([src, tgt]), 0);
    const p1 = layout.get(1)!;
    const p2 = layout.get(2)!;
    expect(p2.y).toBe(p1.y);
    expect(p2.x - p1.x).toBe(STEP_CARDINAL_X);
  });

  it("places a target directly west of the source for a `west` exit", () => {
    const src = makeRoom(1, { west: 2 }, { isRootRoom: true });
    const tgt = makeRoom(2, { east: 1 });
    const layout = computeBfsLayout([src, tgt], flatZLevels([src, tgt]), 0);
    const p1 = layout.get(1)!;
    const p2 = layout.get(2)!;
    expect(p2.y).toBe(p1.y);
    expect(p1.x - p2.x).toBe(STEP_CARDINAL_X);
  });

  it("places a 45° NE target for a `northeast` exit", () => {
    const src = makeRoom(1, { northeast: 2 }, { isRootRoom: true });
    const tgt = makeRoom(2, { southwest: 1 });
    const layout = computeBfsLayout([src, tgt], flatZLevels([src, tgt]), 0);
    const p1 = layout.get(1)!;
    const p2 = layout.get(2)!;
    const dx = p2.x - p1.x;
    const dy = p2.y - p1.y;
    expect(Math.abs(Math.abs(dx) - Math.abs(dy))).toBeLessThan(GRID);
    expect(dx).toBe(STEP_ORDINAL);
    expect(dy).toBe(-STEP_ORDINAL);
  });

  it("snaps every placed position to the grid", () => {
    const src = makeRoom(1, { east: 2, north: 3, south: 4 }, { isRootRoom: true });
    const layout = computeBfsLayout(
      [src, makeRoom(2), makeRoom(3), makeRoom(4)],
      flatZLevels([src, makeRoom(2), makeRoom(3), makeRoom(4)]),
      0,
    );
    layout.forEach((p) => {
      expect(p.x % GRID).toBe(0);
      expect(p.y % GRID).toBe(0);
    });
  });

  it("leaves a revisited room at its first-encountered position (cycles are no-ops)", () => {
    // Cycle: 1 -> 2 (east), 2 -> 3 (south), 3 -> 1 (west). The BFS places
    // 1 at the anchor, 2 east of 1, 3 south of 2. The 3 -> 1 revisit hits
    // the anchor and is skipped. Final layout: 1 at (anchor), 2 east,
    // 3 south of 2.
    const r1 = makeRoom(1, { east: 2 }, { isRootRoom: true });
    const r2 = makeRoom(2, { south: 3 });
    const r3 = makeRoom(3, { west: 1 });
    const layout = computeBfsLayout([r1, r2, r3], flatZLevels([r1, r2, r3]), 0);
    const p1 = layout.get(1)!;
    const p2 = layout.get(2)!;
    const p3 = layout.get(3)!;
    expect(p2.x - p1.x).toBe(STEP_CARDINAL_X);
    expect(p3.y - p2.y).toBe(STEP_CARDINAL_Y);
    // Room 1 is the anchor and is never moved.
    expect(p1.x).toBe(400);
    expect(p1.y).toBe(300);
  });

  it("places orphan rooms (no exits, not reachable) in a 5-column grid below the main layout", () => {
    const r1 = makeRoom(1, {}, { isRootRoom: true });
    const r2 = makeRoom(2, {});
    const r3 = makeRoom(3, {});
    const r4 = makeRoom(4, {});
    const r5 = makeRoom(5, {});
    const r6 = makeRoom(6, {});
    const layout = computeBfsLayout([r1, r2, r3, r4, r5, r6], flatZLevels([r1, r2, r3, r4, r5, r6]), 0);
    const anchorY = layout.get(1)!.y;
    for (let i = 2; i <= 6; i++) {
      expect(layout.get(i)!.y).toBeGreaterThan(anchorY);
    }
    expect(layout.get(6)!.y - layout.get(2)!.y).toBeGreaterThan(0);
    for (let i = 2; i <= 6; i++) {
      expect(layout.get(i)!.x % GRID).toBe(0);
      expect(layout.get(i)!.y % GRID).toBe(0);
    }
  });

  it("uses the root room as the anchor on floor 0", () => {
    const r1 = makeRoom(1, { east: 2 });
    const r2 = makeRoom(2, {}, { isRootRoom: true });
    const layout = computeBfsLayout([r1, r2], flatZLevels([r1, r2]), 0);
    expect(layout.get(2)!.x).toBe(400);
    expect(layout.get(2)!.y).toBe(300);
  });

  it("walks up the chain to anchor floor +1 at the room reachable via the shortest up chain", () => {
    const r1 = makeRoom(1, { up: 2 }, { isRootRoom: true, posZ: 0 });
    const r2 = makeRoom(2, { up: 3 }, { posZ: 1 });
    const r3 = makeRoom(3, {}, { posZ: 2 });
    const rooms = [r1, r2, r3];
    const layout = computeBfsLayout(rooms, flatZLevels(rooms), 1);
    expect(layout.get(2)!.x).toBe(400);
    expect(layout.get(2)!.y).toBe(300);
  });

  it("ignores `up` and `down` exits when placing rooms on the same floor", () => {
    const r1 = makeRoom(1, { up: 2 }, { isRootRoom: true, posZ: 0 });
    const r2 = makeRoom(2, {}, { posZ: 1 });
    const layout = computeBfsLayout([r1, r2], flatZLevels([r1, r2]), 0);
    expect(layout.has(2)).toBe(false);
  });

  it("places target rooms with a non-overlapping buffer regardless of exit count", () => {
    // Each target's box ends strictly past the source's box on the relevant
    // axis, with the constant NODE_W/NODE_H + NODE_GAP step.
    const src = makeRoom(1, { east: 2, north: 3, south: 4, west: 5 }, { isRootRoom: true });
    const layout = computeBfsLayout(
      [src, makeRoom(2), makeRoom(3), makeRoom(4), makeRoom(5)],
      flatZLevels([src, makeRoom(2), makeRoom(3), makeRoom(4), makeRoom(5)]),
      0,
    );
    const p1 = layout.get(1)!;
    const p2 = layout.get(2)!;
    const p3 = layout.get(3)!;
    const p4 = layout.get(4)!;
    const p5 = layout.get(5)!;
    expect(p2.x - p1.x).toBe(STEP_CARDINAL_X);
    expect(p1.y - p3.y).toBe(STEP_CARDINAL_Y);
    expect(p4.y - p1.y).toBe(STEP_CARDINAL_Y);
    expect(p1.x - p5.x).toBe(STEP_CARDINAL_X);
  });
});
