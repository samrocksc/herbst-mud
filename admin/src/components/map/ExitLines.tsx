import type { Room } from './types'

type ExitLinesProps = {
  rooms: Room[]
  nodePositions: Map<number, { x: number; y: number }>
}

export function ExitLines({ rooms, nodePositions }: ExitLinesProps) {
  return (
    <svg className="absolute top-0 left-0 w-full h-full pointer-events-none">
      {rooms.map(room => {
        const pos = nodePositions.get(room.id)
        if (!pos) return null
        return Object.entries(room.exits || {}).map(([dir, targetId]) => {
          const targetPos = nodePositions.get(targetId)
          if (!targetPos || dir === 'up' || dir === 'down') return null
          return (
            <line
              key={`${room.id}-${targetId}-${dir}`}
              x1={pos.x + 60}
              y1={pos.y + 35}
              x2={targetPos.x + 60}
              y2={targetPos.y + 35}
              stroke="var(--color-border)"
              strokeWidth={2}
            />
          )
        })
      })}
    </svg>
  )
}
