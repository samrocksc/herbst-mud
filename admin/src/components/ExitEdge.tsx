import { memo, useState } from 'react'
import {
  BaseEdge,
  EdgeLabelRenderer,
  getBezierPath,
  type EdgeProps,
} from '@xyflow/react'

export interface ExitEdgeData {
  direction: 'north' | 'south' | 'east' | 'west' | 'up' | 'down'
  isZExit: boolean
  animated?: boolean
}

export type ExitEdge = Edge<ExitEdgeData>

interface ExitEdgeProps extends EdgeProps<ExitEdgeData> {}

const directionEmoji: Record<string, string> = {
  north: '↑',
  south: '↓',
  east: '→',
  west: '←',
  up: '⬆',
  down: '⬇',
}

function ExitEdgeComponent({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  data,
  selected,
}: ExitEdgeProps) {
  const [isHovered, setIsHovered] = useState(false)

  const direction = data?.direction || 'east'
  const isZExit = data?.isZExit || false

  const edgePath = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  })

  const labelX = (sourceX + targetX) / 2
  const labelY = (sourceY + targetY) / 2

  const label = directionEmoji[direction] || direction.charAt(0).toUpperCase()

  // Style based on exit type
  const edgeColor = isZExit ? '#e17055' : '#74b9ff' // Orange for z-exits, blue for standard
  const edgeWidth = selected ? 4 : 2
  const labelBgColor = isZExit ? '#d63031' : '#0984e3'

  return (
    <>
      <BaseEdge
        id={id}
        path={edgePath}
        style={{
          stroke: edgeColor,
          strokeWidth: edgeWidth,
          ...(data?.animated ? { animation: 'flow 1s linear infinite' } : {}),
        }}
        className={data?.animated ? 'animated-edge' : undefined}
      />
      <EdgeLabelRenderer>
        <div
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            pointerEvents: 'all',
          }}
          className="edge-label-container"
        >
          <div
            onMouseEnter={() => setIsHovered(true)}
            onMouseLeave={() => setIsHovered(false)}
            style={{
              backgroundColor: labelBgColor,
              color: 'white',
              padding: isHovered ? '6px 10px' : '4px 8px',
              borderRadius: '4px',
              fontSize: isHovered ? '14px' : '12px',
              fontWeight: 'bold',
              whiteSpace: 'nowrap',
              boxShadow: selected ? '0 0 8px rgba(255,255,255,0.5)' : '0 2px 4px rgba(0,0,0,0.3)',
              transition: 'all 0.2s ease',
              cursor: 'pointer',
              border: selected ? '2px solid white' : 'none',
            }}
            title={`Exit: ${direction}${isZExit ? ' (Z-level)' : ''}`}
          >
            {label} {direction}
            {isZExit && ' ⬡'}
          </div>
        </div>
      </EdgeLabelRenderer>

      <style>{`
        @keyframes flow {
          0% { stroke-dashoffset: 20; }
          100% { stroke-dashoffset: 0; }
        }
        .animated-edge {
          stroke-dasharray: 10 10;
          animation: flow 0.5s linear infinite;
        }
      `}</style>
    </>
  )
}

export const ExitEdge = memo(ExitEdgeComponent)