import { memo } from 'react'
import {
  BaseEdge,
  EdgeLabelRenderer,
  getBezierPath,
  Position,
} from '@xyflow/react'

interface ExitEdgeData {
  direction: 'north' | 'south' | 'east' | 'west' | 'up' | 'down'
  isZExit?: boolean
  label?: string
}

export type ExitEdgeType = {
  id: string
  source: string
  target: string
  type: 'exit'
  data: ExitEdgeData
  animated?: boolean
}

interface ExitEdgeProps {
  id: string
  sourceX: number
  sourceY: number
  targetX: number
  targetY: number
  sourcePosition?: Position
  targetPosition?: Position
  data?: ExitEdgeData
  selected?: boolean
  animated?: boolean
}

const directionEmoji: Record<string, string> = {
  north: '↑',
  south: '↓',
  east: '→',
  west: '←',
  up: '⬆',
  down: '⬇',
}

export const directionColor = (direction: string, isZExit?: boolean): string => {
  if (isZExit) {
    return '#ff6b6b' // Red for Z-axis exits
  }
  switch (direction) {
    case 'north':
    case 'south':
      return '#4ecdc4' // Teal
    case 'east':
    case 'west':
      return '#ffe66d' // Yellow
    case 'up':
    case 'down':
      return '#ff6b6b' // Red for vertical
    default:
      return '#a8a8a8' // Gray
  }
}

export const ExitEdge = memo(({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition = Position.Bottom,
  targetPosition = Position.Top,
  data,
  selected,
  animated,
}: ExitEdgeProps) => {
  const direction = data?.direction || 'east'
  const isZExit = data?.isZExit || direction === 'up' || direction === 'down'
  const label = data?.label || directionEmoji[direction] || '→'

  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  })

  const edgeColor = directionColor(direction, isZExit)

  return (
    <>
      <BaseEdge
        id={id}
        path={edgePath}
        style={{
          stroke: edgeColor,
          strokeWidth: selected ? 4 : 2,
          transition: 'stroke-width 0.2s, stroke 0.2s',
        }}
        className={animated ? 'animated' : ''}
      />
      <EdgeLabelRenderer>
        <div
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px, ${labelY}px)`,
            pointerEvents: 'all',
          }}
          className="nodrag nopan"
        >
          <div
            style={{
              backgroundColor: isZExit ? '#2d1f1f' : '#1f2d2d',
              border: `2px solid ${edgeColor}`,
              borderRadius: '8px',
              padding: '4px 8px',
              fontSize: '14px',
              fontWeight: 'bold',
              color: edgeColor,
              whiteSpace: 'nowrap',
              boxShadow: selected
                ? `0 0 8px ${edgeColor}`
                : '0 2px 4px rgba(0,0,0,0.3)',
              transition: 'box-shadow 0.2s',
            }}
          >
            {label}
            {isZExit && <span style={{ fontSize: '10px', marginLeft: '4px' }}>Z</span>}
          </div>
        </div>
      </EdgeLabelRenderer>
    </>
  )
})

ExitEdge.displayName = 'ExitEdge'

// Helper to create an exit edge
export function createExitEdge(
  id: string,
  source: string,
  target: string,
  direction: ExitEdgeData['direction'],
  label?: string
): ExitEdgeType {
  const isZExit = direction === 'up' || direction === 'down'
  return {
    id,
    source,
    target,
    type: 'exit',
    data: {
      direction,
      isZExit,
      label: label || directionEmoji[direction],
    },
    animated: true,
  }
}

export default ExitEdge