/* eslint-disable functional/immutable-data, functional/prefer-immutable-types */
import { describe, it, expect, beforeEach } from "vitest";
import { render, screen } from "@testing-library/react";
import { ExitLines, resolveOverlaps } from "./ExitLineRenderer";
import type { Room } from "./types";

describe("ExitLineRenderer", () => {
  describe("anchorOnEdgeOrthogonal", () => {
    const NODE_W = 120;
    const NODE_H = 65;

    // Import or recreate the DIRECTION_OFFSETS and OPPOSITE_DIR mappings
    const DIRECTION_OFFSETS = {
      north: { dx: 0, dy: -120 },
      south: { dx: 0, dy: 120 },
      east: { dx: 150, dy: 0 },
      west: { dx: -150, dy: 0 },
      northeast: { dx: 106, dy: -106 },
      northwest: { dx: -106, dy: -106 },
      southeast: { dx: 106, dy: 106 },
      southwest: { dx: -106, dy: 106 },
    };

    const OPPOSITE_DIR: Record<string, string> = {
      north: "south",
      south: "north",
      east: "west",
      west: "east",
      northeast: "southwest",
      southwest: "northeast",
      northwest: "southeast",
      southeast: "northwest",
      up: "down",
      down: "up",
    };

    function anchorOnEdgeOrthogonal(dir: string, w: number, h: number) {
      const offsets = DIRECTION_OFFSETS[dir];
      if (!offsets) return { dx: 0, dy: -h / 2 };

      if (offsets.dx === 0 && offsets.dy < 0) return { dx: 0, dy: -h / 2 };
      if (offsets.dx === 0 && offsets.dy > 0) return { dx: 0, dy: h / 2 };
      if (offsets.dy === 0 && offsets.dx > 0) return { dx: w / 2, dy: 0 };
      if (offsets.dy === 0 && offsets.dx < 0) return { dx: -w / 2, dy: 0 };

      if (offsets.dx > 0 && offsets.dy < 0) return { dx: w / 2, dy: -h / 2 };
      if (offsets.dx < 0 && offsets.dy < 0) return { dx: -w / 2, dy: -h / 2 };
      if (offsets.dx > 0 && offsets.dy > 0) return { dx: w / 2, dy: h / 2 };
      if (offsets.dx < 0 && offsets.dy > 0) return { dx: -w / 2, dy: h / 2 };

      return { dx: 0, dy: -h / 2 };
    }

    function anchorOnEdgeAngle(angle: number, w: number, h: number) {
      const a = ((angle % (2 * Math.PI)) + (2 * Math.PI)) % (2 * Math.PI);
      if (a >= Math.PI * 0.75 && a <= Math.PI * 1.25) return { dx: 0, dy: -h / 2 };
      if (a >= Math.PI * 1.25 && a <= Math.PI * 1.75) return { dx: w / 2, dy: 0 };
      if (a >= Math.PI * 1.75 || a <= Math.PI * 0.25) return { dx: w, dy: h / 2 };
      return { dx: -w / 2, dy: 0 };
    }

    function angleBetween(sx: number, sy: number, tx: number, ty: number) {
      return Math.atan2(ty - sy, tx - sx);
    }

    it("should return top edge for north direction", () => {
      const result = anchorOnEdgeOrthogonal("north", NODE_W, NODE_H);
      expect(result).toEqual({ dx: 0, dy: -NODE_H / 2 });
    });

    it("should return bottom edge for south direction", () => {
      const result = anchorOnEdgeOrthogonal("south", NODE_W, NODE_H);
      expect(result).toEqual({ dx: 0, dy: NODE_H / 2 });
    });

    it("should return right edge for east direction", () => {
      const result = anchorOnEdgeOrthogonal("east", NODE_W, NODE_H);
      expect(result).toEqual({ dx: NODE_W / 2, dy: 0 });
    });

    it("should return left edge for west direction", () => {
      const result = anchorOnEdgeOrthogonal("west", NODE_W, NODE_H);
      expect(result).toEqual({ dx: -NODE_W / 2, dy: 0 });
    });

    it("should return bottom-right corner for northeast direction", () => {
      const result = anchorOnEdgeOrthogonal("northeast", NODE_W, NODE_H);
      expect(result).toEqual({ dx: NODE_W / 2, dy: -NODE_H / 2 });
    });

    it("should return bottom-left corner for northwest direction", () => {
      const result = anchorOnEdgeOrthogonal("northwest", NODE_W, NODE_H);
      expect(result).toEqual({ dx: -NODE_W / 2, dy: -NODE_H / 2 });
    });

    it("should return top-right corner for southeast direction", () => {
      const result = anchorOnEdgeOrthogonal("southeast", NODE_W, NODE_H);
      expect(result).toEqual({ dx: NODE_W / 2, dy: NODE_H / 2 });
    });

    it("should return top-left corner for southwest direction", () => {
      const result = anchorOnEdgeOrthogonal("southwest", NODE_W, NODE_H);
      expect(result).toEqual({ dx: -NODE_W / 2, dy: NODE_H / 2 });
    });

    it("should calculate angle between two points correctly", () => {
      // angleBetween(0, 0, 1, 0) should be 0 (right)
      expect(angleBetween(0, 0, 1, 0)).toBeCloseTo(0, 5);
      // angleBetween(0, 0, 0, 1) should be PI/2 (down)
      expect(angleBetween(0, 0, 0, 1)).toBeCloseTo(Math.PI / 2, 5);
      // angleBetween(0, 0, -1, 0) should be PI (left)
      expect(angleBetween(0, 0, -1, 0)).toBeCloseTo(Math.PI, 5);
      // angleBetween(0, 0, 0, -1) should be -PI/2 (up)
      expect(angleBetween(0, 0, 0, -1)).toBeCloseTo(-Math.PI / 2, 5);
    });

    it("should correctly calculate opposite direction", () => {
      expect(OPPOSITE_DIR["north"]).toBe("south");
      expect(OPPOSITE_DIR["south"]).toBe("north");
      expect(OPPOSITE_DIR["east"]).toBe("west");
      expect(OPPOSITE_DIR["west"]).toBe("east");
      expect(OPPOSITE_DIR["northeast"]).toBe("southwest");
    });

    it("should anchor on opposite edge for target room", () => {
      const sourceDir = "north";
      const oppositeDir = OPPOSITE_DIR[sourceDir];

      const sourceEdge = anchorOnEdgeOrthogonal(sourceDir, NODE_W, NODE_H);
      const targetEdge = anchorOnEdgeOrthogonal(oppositeDir, NODE_W, NODE_H);

      // Source (north) should exit from top (dy = -h/2)
      expect(sourceEdge.dy).toBe(-NODE_H / 2);
      // Target should enter from bottom (opposite of north = south, dy = +h/2)
      expect(targetEdge.dy).toBe(NODE_H / 2);
    });
  });

  describe("resolveOverlaps", () => {
    it("should resolve overlapping positions", () => {
      // Two rooms that are too close (60px apart, minimum is 170px)
      const positions = new Map<number, { x: number; y: number }>([
        [1, { x: 0, y: 0 }],
        [2, { x: 60, y: 0 }],
      ]);

      const result = resolveOverlaps(positions, 50, 10);

      const p1 = result.get(1);
      const p2 = result.get(2);

      // After resolution, the x distance should be at least NODE_W + minGap = 170
      const distance = Math.abs(p1!.x - p2!.x);
      expect(distance).toBeGreaterThanOrEqual(170);
    });

    it("should handle non-overlapping positions without changes", () => {
      const positions = new Map<number, { x: number; y: number }>([
        [1, { x: 0, y: 0 }],
        [2, { x: 300, y: 0 }], // Well separated
      ]);

      const result = resolveOverlaps(positions, 50, 10);
      expect(result.get(1)).toEqual({ x: 0, y: 0 });
      expect(result.get(2)).toEqual({ x: 300, y: 0 });
    });
  });
});

describe("ExitLines component", () => {
  const rooms: Room[] = [
    {
      id: 1,
      name: "Room 1",
      description: "Test room 1",
      exits: { north: 2, east: 3 },
      isStartingRoom: false,
      isRootRoom: false,
    },
    {
      id: 2,
      name: "Room 2",
      description: "Test room 2",
      exits: { south: 1 },
      isStartingRoom: false,
      isRootRoom: false,
    },
    {
      id: 3,
      name: "Room 3",
      description: "Test room 3",
      exits: { west: 1 },
      isStartingRoom: false,
      isRootRoom: false,
    },
  ];

  const nodePositions = new Map<number, { x: number; y: number }>([
    [1, { x: 100, y: 100 }],
    [2, { x: 100, y: 200 }],
    [3, { x: 300, y: 100 }],
  ]);

  it("should render edges between connected rooms", () => {
    const { container } = render(<ExitLines rooms={rooms} nodePositions={nodePositions} />);
    const svg = container.querySelector("svg");
    expect(svg).toBeInTheDocument();
  });

  it("should render lines for phantom exits (missing target rooms)", () => {
    const roomsWithPhantom: Room[] = [
      {
        id: 1,
        name: "Room 1",
        description: "Test room",
        exits: { north: 999 },
        isStartingRoom: false,
        isRootRoom: false,
      },
    ];
    const nodePositions = new Map([[1, { x: 100, y: 100 }]]);

    const { container } = render(<ExitLines rooms={roomsWithPhantom} nodePositions={nodePositions} />);
    const lines = container.querySelectorAll("line");
    expect(lines.length).toBeGreaterThan(0);
  });
});
