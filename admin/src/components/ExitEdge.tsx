import { EdgeProps, getBezierPath, EdgeLabelRenderer } from '@xyflow/react'
import { useMemo } from 'react'

interface ExitEdgeData {
  direction: string
  isZExit?: boolean
}

export type ExitEdgeType = EdgeProps & {
  data?: ExitEdgeData
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
}: ExitEdgeType) {
  const [edgePath, labelX, labelY] = getBezierPath({
    sourceX,
    sourceY,
    targetX,
    targetY,
    sourcePosition,
    targetPosition,
  })

  const direction = data?.direction || ''
  const isZExit = data?.isZExit || direction === 'up' || direction === 'down'

  const edgeStyle = useMemo(() => {
    return {
      stroke: isZExit ? '#9b59b6' : '#4a8c3f',
      strokeWidth: 2,
      strokeDasharray: isZExit ? '5,5' : undefined,
    }
  }, [isZExit])

  const arrowStyle = useMemo(() => {
    return {
      stroke: isZExit ? '#9b59b6' : '#4a8c3f',
      fill: 'none',
    }
  }, [isZExit])

  return (
    <>
      <path
        id={id}
        style={{
          ...edgeStyle,
          opacity: selected ? 1 : 0.7,
        }}
        className="react-flow__edge-path"
        d={edgePath}
      />
      <EdgeLabelRenderer>
        <div
          style={{
            position: 'absolute',
            transform: `translate(-50%, -50%) translate(${labelX}px,${labelY}px)`,
            pointerEvents: 'all',
          }}
        >
          <div
            style={{
              background: isZExit ? '#9b59b6' : '#4a8c3f',
              color: '#fff',
              padding: '2px 6px',
              borderRadius: '4px',
              fontSize: '10px',
              fontWeight: 'bold',
              whiteSpace: 'nowrap',
            }}
          >
            {direction}
          </div>
        </div>
      </EdgeLabelRenderer>
      {/* Arrow marker */}
      <svg
        style={{
          position: 'absolute',
          width: 20,
          height: 20,
          overflow: 'visible',
          pointerEvents: 'none',
        }}
      >
        <defs>
          <marker
            id={`arrow-${id}`}
            markerWidth="12"
            markerHeight="12"
            refX="10"
            refY="3"
            orient="auto"
            markerUnits="strokeWidth"
          >
            <path d="M0,0 L0,6 L9,3 z" fill={isZExit ? '#9b59b6' : '#4a8c3f'} />
          </marker>
        </defs>
      </svg>
    </>
  )
}

export default ExitEdge