import { NODE_W, NODE_H } from "./constants";
import type { Room } from "./types";

export const OPPOSITE_DIR: Record<string, string> = {
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

const MIN_GAP_X = NODE_W + 20;
const MIN_GAP_Y = NODE_H + 20;

export const DIRECTION_OFFSETS: Record<string, { dx: number; dy: number }> = {
  north: { dx: 0, dy: -MIN_GAP_Y },
  south: { dx: 0, dy: MIN_GAP_Y },
  east: { dx: MIN_GAP_X, dy: 0 },
  west: { dx: -MIN_GAP_X, dy: 0 },
  northeast: { dx: MIN_GAP_X * 0.7, dy: -MIN_GAP_Y * 0.7 },
  northwest: { dx: -MIN_GAP_X * 0.7, dy: -MIN_GAP_Y * 0.7 },
  southeast: { dx: MIN_GAP_X * 0.7, dy: MIN_GAP_Y * 0.7 },
  southwest: { dx: -MIN_GAP_X * 0.7, dy: MIN_GAP_Y * 0.7 },
};

export const DirectionShortLabels: Record<string, string> = {
  north: "N",
  northeast: "NE",
  east: "E",
  southeast: "SE",
  south: "S",
  southwest: "SW",
  west: "W",
  northwest: "NW",
  up: "▲",
  down: "▼",
};

export const ALL_DIRECTIONS = [
  "north",
  "northeast",
  "east",
  "southeast",
  "south",
  "southwest",
  "west",
  "northwest",
  "up",
  "down",
];

const BASE_NODE_H = 85;
const EXIT_LINE_HEIGHT = 18;

/**
 * Estimate the rendered bounding box of a room node.
 * Width is fixed by Tailwind (w-[120px]). Height grows with the number of
 * visible exit lines stacked below the name/id badges plus a safety margin
 * so that diagonal/short orthogonal offsets do not produce overlaps.
 */
export function estimateNodeSize(room: Room): { w: number; h: number } {
  const exitCount = Object.keys(room.exits || {}).length;
  return {
    w: NODE_W,
    h: BASE_NODE_H + Math.max(0, exitCount) * EXIT_LINE_HEIGHT,
  };
}

/**
 * Return the top-left offset needed to place the target room next to the
 * source room in the given direction without the boxes overlapping.
 * The offset is relative to the source room's top-left position.
 * `margin` is the minimum empty gap between the two boxes.
 */
export function getNoOverlapOffset(
  dir: string,
  sourceBox: Readonly<{ w: number; h: number }>,
  targetBox: Readonly<{ w: number; h: number }>,
  margin = 30,
): { dx: number; dy: number } {
  switch (dir) {
    case "north":
      return { dx: 0, dy: -(targetBox.h + margin) };
    case "south":
      return { dx: 0, dy: sourceBox.h + margin };
    case "east":
      return { dx: sourceBox.w + margin, dy: 0 };
    case "west":
      return { dx: -(targetBox.w + margin), dy: 0 };
    case "northeast":
      return { dx: sourceBox.w + margin, dy: -(targetBox.h + margin) };
    case "northwest":
      return { dx: -(targetBox.w + margin), dy: -(targetBox.h + margin) };
    case "southeast":
      return { dx: sourceBox.w + margin, dy: sourceBox.h + margin };
    case "southwest":
      return { dx: -(targetBox.w + margin), dy: sourceBox.h + margin };
    default:
      return { dx: sourceBox.w + margin, dy: 0 };
  }
}
