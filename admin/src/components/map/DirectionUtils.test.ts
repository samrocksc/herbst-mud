/* eslint-disable functional/no-return-void */
import { describe, it, expect } from "vitest";
import { getNoOverlapOffset, estimateNodeSize } from "./DirectionUtils";
import { NODE_W, NODE_H } from "./constants";
import type { Room } from "./types";

describe("estimateNodeSize", () => {
  it("returns the pinned NODE_W x NODE_H box regardless of exit count", () => {
    const room: Room = {
      id: 1,
      name: "Test",
      description: "",
      exits: { north: 2, south: 3, east: 4, west: 5 },
    };
    const size = estimateNodeSize(room);
    expect(size.w).toBe(NODE_W);
    expect(size.h).toBe(NODE_H);
  });

  it("uses the same pinned box when there are no exits", () => {
    const room: Room = { id: 1, name: "Test", description: "", exits: {} };
    expect(estimateNodeSize(room)).toEqual({ w: NODE_W, h: NODE_H });
  });
});

describe("getNoOverlapOffset", () => {
  const source = { w: 120, h: 80 };
  const target = { w: 120, h: 100 };
  const margin = 20;

  it("places a southern target below the source box plus margin", () => {
    const off = getNoOverlapOffset("south", source, target, margin);
    expect(off.dx).toBe(0);
    expect(off.dy).toBe(source.h + margin);
  });

  it("places a northern target above the target box plus margin", () => {
    const off = getNoOverlapOffset("north", source, target, margin);
    expect(off.dx).toBe(0);
    expect(off.dy).toBe(-(target.h + margin));
  });

  it("places an eastern target to the right of the source box plus margin", () => {
    const off = getNoOverlapOffset("east", source, target, margin);
    expect(off.dx).toBe(source.w + margin);
    expect(off.dy).toBe(0);
  });

  it("places a western target to the left of the target box plus margin", () => {
    const off = getNoOverlapOffset("west", source, target, margin);
    expect(off.dx).toBe(-(target.w + margin));
    expect(off.dy).toBe(0);
  });

  it("places a diagonal target outside both source and target dimensions", () => {
    const ne = getNoOverlapOffset("northeast", source, target, margin);
    expect(ne.dx).toBe(source.w + margin);
    expect(ne.dy).toBe(-(target.h + margin));

    const sw = getNoOverlapOffset("southwest", source, target, margin);
    expect(sw.dx).toBe(-(target.w + margin));
    expect(sw.dy).toBe(source.h + margin);
  });
});
