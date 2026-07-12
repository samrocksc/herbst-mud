export const GRID = 20;
export const STEP = 100;
export const ORPHAN_COLS = 5;
export const DRAG_THRESHOLD = 5;
export const CANVAS_W = 3000;
export const CANVAS_H = 3000;
export const MIN_ZOOM = 0.5;
export const MAX_ZOOM = 2.0;
export const ZOOM_STEP = 0.25;
export const ZOOM_FINE_STEP = 0.1;
// Pinned node box. Width and height are used for the BFS step so every
// room reserves the same canvas area regardless of label length. Long
// names and exit labels are truncated inside the box.
export const NODE_W = 140;
export const NODE_H = 100;
// BFS step: parent's half-width + target's half-width + gap, with the
// gap doubled so the line between them stays visually clear.
export const NODE_GAP = 20;
export const EXIT_LINE_HEIGHT = 18;
export const BASE_NODE_H = 65;
