import { useState, memo } from 'react'
import type { Room } from './types'

type ExitLinesProps = {
  rooms: Room[]
  nodePositions: Map<number, { x: number; y: number }>
}

/** Must match Tailwind classes in RoomNode. */
const NODE_W = 120
const NODE_H = 65
const NODE_HALF_W = NODE_W / 2
const NODE_HALF_H = NODE_H / 2

type Segment = {
  key: string
  sx: number
  sy: number
  tx: number
  ty: number
  fwdDir: string
  phantom: boolean
}

function anchorOnEdge(
  angle: number,
  w: number,
  h: number,
): { dx: number; dy: number } {
  const a = ((angle % (2 * Math.PI)) + (2 * Math.PI)) % (2 * Math.PI)
  if (a >= Math.PI * 0.75 && a <= Math.PI * 1.25) return { dx: 0, dy: h / 2 }
  if (a >= Math.PI * 1.25 && a <= Math.PI * 1.75) return { dx: w / 2, dy: 0 }
  if (a >= Math.PI * 1.75 || a <= Math.PI * 0.25) return { dx: w, dy: h / 2 }
  return { dx: w / 2, dy: h }
}

function angleBetween(sx: number, sy: number, tx: number, ty: number): number {
  return Math.atan2(ty - sy, tx - sx)
}

export function resolveOverlaps(
  positions: Map<number, { x: number; y: number }>,
  minGap = 50,
  maxIterations = 50,
): Map<number, { x: number; y: number }> {
  const result = new Map<number, { x: number; y: number }>()
  for (const [id, pos] of positions) result.set(id, { ...pos })

  let moved = true
  for (let iter = 0; iter < maxIterations && moved; iter++) {
    moved = false
    const ids = Array.from(result.keys())
    for (let i = 0; i < ids.length; i++) {
      for (let j = i + 1; j < ids.length; j++) {
        const a = ids[i]
        const b = ids[j]
        const pa = result.get(a)!
        const pb = result.get(b)!

        const overlapX = (NODE_W + minGap) - Math.abs(pa.x - pb.x)
        const overlapY = (NODE_H + minGap) - Math.abs(pa.y - pb.y)

        if (overlapX > 0 && overlapY > 0) {
          const dx = pa.x - pb.x
          const dy = pa.y - pb.y
          if (overlapX < overlapY) {
            const push = overlapX / 2
            pa.x += dx > 0 ? push : -push
            pb.x += dx > 0 ? -push : push
          } else {
            const push = overlapY / 2
            pa.y += dy > 0 ? push : -push
            pb.y += dy > 0 ? -push : push
          }
          moved = true
        }
      }
    }
  }
  return result
}

export const ExitLines = memo(function ExitLines({ rooms, nodePositions }: ExitLinesProps) {
  const [hoveredKey, setHoveredKey] = useState<string | null>(null)

  const roomIds = new Set(rooms.map(r => r.id))
  const drawn = new Set<string>()
  const segments: Segment[] = []

  for (const room of rooms) {
    const pos = nodePositions.get(room.id)
    if (!pos) continue

    for (const [dir, targetId] of Object.entries(room.exits || {})) {
      if (dir === 'up' || dir === 'down') continue
      const targetPos = nodePositions.get(targetId)

      // Phantom exit: target room doesn't exist
      if (!roomIds.has(targetId)) {
        // Draw a short dashed line from the source room in the exit direction
        const sourceCx = pos.x + NODE_HALF_W
        const sourceCy = pos.y + NODE_HALF_H
        const angleMap: Record<string, number> = {
          north: -Math.PI / 2, south: Math.PI / 2,
          east: 0, west: Math.PI,
          northeast: -Math.PI / 4, northwest: -3 * Math.PI / 4,
          southeast: Math.PI / 4, southwest: 3 * Math.PI / 4,
        }
        const angle = angleMap[dir] ?? 0
        const len = 60
        segments.push({
          key: `phantom-${room.id}-${dir}`,
          sx: sourceCx, sy: sourceCy,
          tx: sourceCx + Math.cos(angle) * len,
          ty: sourceCy + Math.sin(angle) * len,
          fwdDir: dir,
          phantom: true,
        })
        continue
      }

      if (!targetPos) continue

      const [lo, hi] = room.id < targetId ? [room.id, targetId] : [targetId, room.id]
      const canon = `${lo}-${hi}`
      if (drawn.has(canon)) continue
      drawn.add(canon)

      const cx = pos.x + NODE_HALF_W
      const cy = pos.y + NODE_HALF_H
      const tcx = targetPos.x + NODE_HALF_W
      const tcy = targetPos.y + NODE_HALF_H
      const fwdAngle = angleBetween(cx, cy, tcx, tcy)
      const sourceEdge = anchorOnEdge(fwdAngle, NODE_W, NODE_H)
      const targetEdge = anchorOnEdge(fwdAngle + Math.PI, NODE_W, NODE_H)

      segments.push({
        key: canon,
        sx: pos.x + sourceEdge.dx,
        sy: pos.y + sourceEdge.dy,
        tx: targetPos.x + targetEdge.dx,
        ty: targetPos.y + targetEdge.dy,
        fwdDir: dir,
        phantom: false,
      })
    }
  }

  return (
    <svg className="absolute top-0 left-0 w-full h-full pointer-events-none">
      {segments.map(seg => {
        const isHovered = hoveredKey === seg.key

        if (seg.phantom) {
          return (
            <line
              key={seg.key}
              x1={seg.sx} y1={seg.sy} x2={seg.tx} y2={seg.ty}
              stroke="var(--color-danger, #ef4444)"
              strokeWidth={2}
              strokeDasharray="6 4"
              opacity={0.7}
            />
          )
        }

        const strokeColor = isHovered ? 'var(--color-accent)' : 'var(--color-text-muted)'
        const strokeWidth = isHovered ? 3 : 2

        return (
          <g
            key={seg.key}
            className="pointer-events-auto"
            onMouseEnter={() => setHoveredKey(seg.key)}
            onMouseLeave={() => setHoveredKey(null)}
            style={{ cursor: 'pointer' }}
          >
            <line x1={seg.sx} y1={seg.sy} x2={seg.tx} y2={seg.ty} stroke="transparent" strokeWidth={14} />
            <line
              x1={seg.sx} y1={seg.sy} x2={seg.tx} y2={seg.ty}
              stroke={strokeColor}
              strokeWidth={strokeWidth}
              style={{ transition: 'stroke 0.15s, stroke-width 0.15s' }}
            />
          </g>
        )
      })}
    </svg>
  )
})