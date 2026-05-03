import { useState } from 'react'
import { OPPOSITE_DIR } from './DirectionUtils'
import type { Room } from './types'

type ExitLinesProps = {
  rooms: Room[]
  nodePositions: Map<number, { x: number; y: number }>
}

/** Half-width of a room node (used to offset line start/end from node center). */
const NODE_HALF_W = 60
/** Half-height of a room node. */
const NODE_HALF_H = 35
/** Size of the direction-arrow triangle. */
const ARROW_SIZE = 7
/** How far (px) to offset each arrow from midpoint along the line. */
const ARROW_OFFSET = 18

type Segment = {
  key: string
  sx: number
  sy: number
  tx: number
  ty: number
  midX: number
  midY: number
  fwdDir: string
  revDir: string
  fwdAngle: number
  revAngle: number
  fwdArrowX: number
  fwdArrowY: number
  revArrowX: number
  revArrowY: number
}

/**
 * Compute the angle (radians) from source center to target center.
 * Positive-x axis = 0, clockwise positive (SVG coords).
 */
function angleBetween(
  sx: number, sy: number,
  tx: number, ty: number,
): number {
  return Math.atan2(ty - sy, tx - sx)
}

/**
 * Build three vertices of a small isoceles triangle at (cx, cy)
 * pointing in the direction of `angle`.
 */
function arrowPoints(cx: number, cy: number, angle: number): string {
  const tipX = cx + ARROW_SIZE * Math.cos(angle)
  const tipY = cy + ARROW_SIZE * Math.sin(angle)
  const baseAngleL = angle + (2.4 * Math.PI) / 3
  const baseAngleR = angle - (2.4 * Math.PI) / 3
  const blX = cx + ARROW_SIZE * Math.cos(baseAngleL)
  const blY = cy + ARROW_SIZE * Math.sin(baseAngleL)
  const brX = cx + ARROW_SIZE * Math.cos(baseAngleR)
  const brY = cy + ARROW_SIZE * Math.sin(baseAngleR)
  return `${tipX},${tipY} ${blX},${blY} ${brX},${brY}`
}

export function ExitLines({ rooms, nodePositions }: ExitLinesProps) {
  const [hoveredKey, setHoveredKey] = useState<string | null>(null)

  /** Track which bidirectional pair we've already drawn to skip the reverse. */
  const drawn = new Set<string>()

  const segments: Segment[] = []

  for (const room of rooms) {
    const pos = nodePositions.get(room.id)
    if (!pos) continue

    for (const [dir, targetId] of Object.entries(room.exits || {})) {
      if (dir === 'up' || dir === 'down') continue
      const targetPos = nodePositions.get(targetId)
      if (!targetPos) continue

      /** Canonical key so A→B and B→A produce the same string. */
      const [lo, hi] = room.id < targetId
        ? [room.id, targetId]
        : [targetId, room.id]
      const canon = `${lo}-${hi}`
      if (drawn.has(canon)) continue
      drawn.add(canon)

      const sx = pos.x + NODE_HALF_W
      const sy = pos.y + NODE_HALF_H
      const tx = targetPos.x + NODE_HALF_W
      const ty = targetPos.y + NODE_HALF_H
      const midX = (sx + tx) / 2
      const midY = (sy + ty) / 2
      const dx = tx - sx
      const dy = ty - sy
      const len = Math.sqrt(dx * dx + dy * dy)
      const nx = len > 0 ? dx / len : 1
      const ny = len > 0 ? dy / len : 0

      const fwdAngle = angleBetween(sx, sy, tx, ty)
      const revDir = OPPOSITE_DIR[dir] ?? dir
      const revAngle = angleBetween(tx, ty, sx, sy)

      segments.push({
        key: canon,
        sx, sy, tx, ty, midX, midY,
        fwdDir: dir,
        revDir,
        fwdAngle,
        revAngle,
        fwdArrowX: midX + nx * ARROW_OFFSET,
        fwdArrowY: midY + ny * ARROW_OFFSET,
        revArrowX: midX - nx * ARROW_OFFSET,
        revArrowY: midY - ny * ARROW_OFFSET,
      })
    }
  }

  return (
    <svg className="absolute top-0 left-0 w-full h-full pointer-events-none">
      {segments.map(seg => {
        const isHovered = hoveredKey === seg.key
        const strokeColor = isHovered
          ? 'var(--color-accent)'
          : 'var(--color-text-muted)'
        const strokeWidth = isHovered ? 3 : 2
        const arrowFill = isHovered
          ? 'var(--color-accent)'
          : 'var(--color-text-muted)'

        return (
          <g
            key={seg.key}
            className="pointer-events-auto"
            onMouseEnter={() => setHoveredKey(seg.key)}
            onMouseLeave={() => setHoveredKey(null)}
            style={{ cursor: 'pointer' }}
          >
            {/* Invisible wider hit-area for easier hovering */}
            <line
              x1={seg.sx}
              y1={seg.sy}
              x2={seg.tx}
              y2={seg.ty}
              stroke="transparent"
              strokeWidth={14}
            />
            {/* Visible line */}
            <line
              x1={seg.sx}
              y1={seg.sy}
              x2={seg.tx}
              y2={seg.ty}
              stroke={strokeColor}
              strokeWidth={strokeWidth}
              style={{ transition: 'stroke 0.15s, stroke-width 0.15s' }}
            />
            {/* Forward direction arrow */}
            <polygon
              points={arrowPoints(seg.fwdArrowX, seg.fwdArrowY, seg.fwdAngle)}
              fill={arrowFill}
              style={{ transition: 'fill 0.15s' }}
            />
            {/* Reverse direction arrow (slightly transparent) */}
            <polygon
              points={arrowPoints(seg.revArrowX, seg.revArrowY, seg.revAngle)}
              fill={arrowFill}
              opacity={0.55}
              style={{ transition: 'fill 0.15s' }}
            />
          </g>
        )
      })}
    </svg>
  )
}