import { memo } from 'react'
import {
  BaseEdge,
  EdgeLabelRenderer,
  getBezierPath,
  type EdgeProps,
} from '@xyflow/react'

export interface ExitEdgeData {
  direction: 'north' | 'south' | 'east' | 'west' | 'up' | 'down'
  isZExit?: boolean
}

export type ExitEdgeType = 'exit'

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
}: EdgeProps<ExitEdgeData>) {
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  })

  const direction = data?.direction || 'east'
  const isZExit = data?.isZExit || direction === 'up' || direction === 'down'
  
  // Direction arrow symbols
  const directionSymbols: Record<string, string> = {
    north: '↑',
    south: '↓',
    east: '→',
    west: '←',
    up: '⬆',
    down: '⬇',
  }

  const label = directionSymbols[direction] || '→'

  // Style based on exit type
  const edgeColor = isZExit ? '#f39c12' : selected ? '#6c5ce7' : '#4a9eff'
  const edgeWidth = selected ? 3 : 2
  const edgeStyle = {
    strokeWidth: edgeWidth,
    stroke: edgeColor,
  }

  // Animated style for selected edges
  const animatedStyle = selected ? { ...edgeStyle, strokeDasharray: '5,5' } : edgeStyle

  return (
    <>
      <BaseEdge id={id} path={edgePath} style={animatedStyle} />
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
            style={{
              backgroundColor: isZExit ? '#f39c12' : '#4a9eff',
              color: '#fff',
              padding: '2px 6px',
              borderRadius: '4px',
              fontSize: '12px',
              fontWeight: 'bold',
              boxShadow: '0 2px 4px rgba(0,0,0,0.3)',
              border: selected ? '2px solid #fff' : 'none',
            }}
          >
            {label}
          </div>
        </div>
      </EdgeLabelRenderer>
    </>
  )
}

export const ExitEdge = memo(ExitEdgeComponent)

export default ExitEdge