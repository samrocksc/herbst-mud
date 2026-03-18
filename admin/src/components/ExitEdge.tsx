import { BaseEdge, EdgeLabelRenderer, getBezierPath } from '@xyflow/react'
import type { EdgeProps } from '@xyflow/react'
import { useState, useCallback } from 'react'

/**
 * ExitEdge Component
 * 
 * A custom edge for displaying directed exits between rooms in the map builder.
 * Shows direction arrows (n/s/e/w/up/down) and supports click-to-edit functionality.
 * 
 * 🟣 Donatello - Turtle Time Map Builder
 */

export interface ExitEdgeData extends Record<string, unknown> {
  direction?: 'north' | 'south' | 'east' | 'west' | 'up' | 'down'
  isZExit?: boolean
  bidirectional?: boolean
  onEditClick?: (edgeId: string, data: ExitEdgeData) => void
}

const DIRECTION_CONFIG: Record<string, { symbol: string; label: string; color: string }> = {
  north: { symbol: '⬆️', label: 'N', color: '#74b9ff' },
  south: { symbol: '⬇️', label: 'S', color: '#74b9ff' },
  east: { symbol: '➡️', label: 'E', color: '#74b9ff' },
  west: { symbol: '⬅️', label: 'W', color: '#74b9ff' },
  up: { symbol: '🆙', label: 'U', color: '#e17055' },
  down: { symbol: '🔽', label: 'D', color: '#a29bfe' },
}

export function ExitEdge({
  id,
  sourceX,
  sourceY,
  targetX,
  targetY,
  sourcePosition,
  targetPosition,
  data,
  selected,
  style = {},
  markerEnd,
}: EdgeProps) {
  const [isHovered, setIsHovered] = useState(false)
  
  const edgeData = (data as ExitEdgeData) || {}
  const direction = edgeData.direction || 'east'
  const config = DIRECTION_CONFIG[direction] || DIRECTION_CONFIG.east
  
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    sourcePosition,
    targetX,
    targetY,
    targetPosition,
  })

  const handleClick = useCallback(() => {
    if (edgeData.onEditClick) {
      edgeData.onEditClick(id, edgeData)
    }
  }, [id, edgeData])

  const edgeStyle = {
    ...style,
    stroke: selected ? '#6c5ce7' : config.color,
    strokeWidth: selected ? 3 : 2,
    transition: 'stroke 0.2s ease',
  }

  // Calculate arrow marker position
  const arrowSize = 10
  const angle = Math.atan2(targetY - sourceY, targetX - sourceX)
  const arrowX = targetX - Math.cos(angle) * 15
  const arrowY = targetY - Math.sin(angle) * 15

  return (
    <>
      {/* Invisible wider path for easier hover/click */}
      <path
        d={edgePath}
        fill="none"
        strokeWidth={20}
        stroke="transparent"
        style={{ cursor: 'pointer' }}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
        onClick={handleClick}
      />
      
      {/* Visible edge path */}
      <BaseEdge
        id={id}
        path={edgePath}
        style={edgeStyle}
        markerEnd={markerEnd}
        interactionWidth={0}
      />
      
      {/* Direction arrow at end */}
      <polygon
        points={`
          ${arrowX},${arrowY}
          ${arrowX - arrowSize * Math.cos(angle - Math.PI / 6)},${arrowY - arrowSize * Math.sin(angle - Math.PI / 6)}
          ${arrowX - arrowSize * Math.cos(angle + Math.PI / 6)},${arrowY - arrowSize * Math.sin(angle + Math.PI / 6)}
        `}
        fill={selected ? '#6c5ce7' : config.color}
        style={{ transition: 'fill 0.2s ease' }}
      />
      
      {/* Label with direction */}
      <EdgeLabelRenderer>
        <div
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px, ${labelY}px)`,
            pointerEvents: 'all',
            cursor: 'pointer',
          }}
          className="nodrag nopan"
          onClick={handleClick}
          onMouseEnter={() => setIsHovered(true)}
          onMouseLeave={() => setIsHovered(false)}
        >
          <div
            style={{
              background: selected ? '#6c5ce7' : edgeData.isZExit ? '#e17055' : '#2d5a27',
              padding: '4px 8px',
              borderRadius: '4px',
              fontSize: '12px',
              fontWeight: 'bold',
              color: '#fff',
              border: isHovered || selected ? '2px solid #a29bfe' : '1px solid #444',
              boxShadow: isHovered || selected ? '0 0 8px rgba(108, 92, 231, 0.5)' : 'none',
              transition: 'all 0.2s ease',
              whiteSpace: 'nowrap',
            }}
          >
            {config.symbol} {config.label}
          </div>
        </div>
      </EdgeLabelRenderer>
      
      {/* Z-exit special indicator */}
      {edgeData.isZExit && (
        <EdgeLabelRenderer>
          <div
            style={{
              position: 'absolute',
              transform: `translate(-50%, -50%) translate(${labelX}px, ${labelY - 25}px)`,
              pointerEvents: 'none',
            }}
          >
            <div
              style={{
                background: '#e17055',
                padding: '2px 6px',
                borderRadius: '4px',
                fontSize: '10px',
                color: '#fff',
                border: '1px solid #ff7675',
              }}
            >
              Z-Exit
            </div>
          </div>
        </EdgeLabelRenderer>
      )}
    </>
  )
}