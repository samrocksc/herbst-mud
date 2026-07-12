/**
 * Unit tests for the map layout system.
 *
 * After dropping posx/posy, room positions are derived from the exits graph
 * via `computeBfsLayout` in `useNodeLayout`. The BFS layout itself is
 * covered in `map-layout.test.ts`. This file now covers the remaining
 * invariant: the cardinal/ordinal direction math (DIRECTION_OFFSETS,
 * OPPOSITE_DIR, DirectionShortLabels, ALL_DIRECTIONS).
 *
 * The previous version of this file validated that stored positions fell
 * within canvas bounds and that stored dx/dy matched labeled directions.
 * Those tests assumed positions were persisted on the Room model, which
 * is no longer true. They have been deleted.
 */
import { describe, it, expect } from "vitest";
import {
  DIRECTION_OFFSETS,
  OPPOSITE_DIR,
  ALL_DIRECTIONS,
  DirectionShortLabels,
} from "./DirectionUtils";

describe("DirectionUtils", () => {
  describe("cardinal directions are purely horizontal or vertical", () => {
    it.each([
      ["north", { dx: 0, dyLessThan: 0 }],
      ["south", { dx: 0, dyGreaterThan: 0 }],
      ["east", { dxGreaterThan: 0, dy: 0 }],
      ["west", { dxLessThan: 0, dy: 0 }],
    ])("%s", (dir, expected) => {
      const o = DIRECTION_OFFSETS[dir as keyof typeof DIRECTION_OFFSETS];
      if ("dx" in expected) expect(o.dx).toBe(expected.dx);
      if ("dy" in expected) expect(o.dy).toBe(expected.dy);
      if ("dxGreaterThan" in expected) expect(o.dx).toBeGreaterThan(expected.dxGreaterThan);
      if ("dxLessThan" in expected) expect(o.dx).toBeLessThan(expected.dxLessThan);
      if ("dyGreaterThan" in expected) expect(o.dy).toBeGreaterThan(expected.dyGreaterThan);
      if ("dyLessThan" in expected) expect(o.dy).toBeLessThan(expected.dyLessThan);
    });
  });

  describe("ordinal directions are 45° diagonals (both axes non-zero, same sign convention)", () => {
    it.each([
      ["northeast", 1, -1],
      ["northwest", -1, -1],
      ["southeast", 1, 1],
      ["southwest", -1, 1],
    ])("%s", (dir, xSign, ySign) => {
      const o = DIRECTION_OFFSETS[dir as keyof typeof DIRECTION_OFFSETS];
      expect(Math.sign(o.dx)).toBe(xSign);
      expect(Math.sign(o.dy)).toBe(ySign);
      expect(o.dx).not.toBe(0);
      expect(o.dy).not.toBe(0);
    });
  });

  describe("opposite direction pairs", () => {
    it("cardinal", () => {
      expect(OPPOSITE_DIR.north).toBe("south");
      expect(OPPOSITE_DIR.south).toBe("north");
      expect(OPPOSITE_DIR.east).toBe("west");
      expect(OPPOSITE_DIR.west).toBe("east");
    });
    it("ordinal", () => {
      expect(OPPOSITE_DIR.northeast).toBe("southwest");
      expect(OPPOSITE_DIR.northwest).toBe("southeast");
      expect(OPPOSITE_DIR.southeast).toBe("northwest");
      expect(OPPOSITE_DIR.southwest).toBe("northeast");
    });
    it("vertical", () => {
      expect(OPPOSITE_DIR.up).toBe("down");
      expect(OPPOSITE_DIR.down).toBe("up");
    });
  });

  describe("direction short labels", () => {
    it.each([
      ["north", "N"],
      ["south", "S"],
      ["east", "E"],
      ["west", "W"],
      ["northeast", "NE"],
      ["northwest", "NW"],
      ["southeast", "SE"],
      ["southwest", "SW"],
      ["up", "▲"],
      ["down", "▼"],
    ])("%s → %s", (dir, label) => {
      expect(DirectionShortLabels[dir as keyof typeof DirectionShortLabels]).toBe(label);
    });
  });

  describe("ALL_DIRECTIONS list", () => {
    it("contains all 10 directions", () => {
      expect(ALL_DIRECTIONS).toHaveLength(10);
      expect(ALL_DIRECTIONS).toContain("north");
      expect(ALL_DIRECTIONS).toContain("south");
      expect(ALL_DIRECTIONS).toContain("east");
      expect(ALL_DIRECTIONS).toContain("west");
      expect(ALL_DIRECTIONS).toContain("northeast");
      expect(ALL_DIRECTIONS).toContain("northwest");
      expect(ALL_DIRECTIONS).toContain("southeast");
      expect(ALL_DIRECTIONS).toContain("southwest");
      expect(ALL_DIRECTIONS).toContain("up");
      expect(ALL_DIRECTIONS).toContain("down");
    });
  });
});
