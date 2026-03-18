import { useCallback, useState, useMemo } from 'react'
import {
  ReactFlow,
  Background,
  Controls,
  MiniMap,
  Connection,
  Edge,
  Node,
  useNodesState,
  useEdgesState,
  addEdge,
  MarkerType,
  BackgroundVariant,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'

import { RoomNode, RoomNodeType } from './RoomNode'
import { ExitEdge, ExitEdgeType } from './ExitEdge'

// Direction opposite mapping
const OPPOSITE_DIRECTION: Record<string, string> = {
  north: 'south',
  south: 'north',
  east: 'west',
  west: 'east',
  up: 'down',
  down: 'up',
}

// Calculate direction based on relative node positions
function calculateDirection(
  sourceX: number,
  sourceY: number,
  targetX: number,
  targetY: number
): string {
  const dx = targetX - sourceX
  const dy = targetY - sourceY

  if (Math.abs(dy) > Math.abs(dx)) {
    return dy > 0 ? 'south' : 'north'
  }
  return dx > 0 ? 'east' : 'west'
}

interface DirectionModalProps {
  connection: { source: string; target: string; sourceX: number; sourceY: number; targetX: number; targetY: number }
  onSelect: (direction: string) => void
  onCancel: () => void
}

function DirectionModal({ connection, onSelect, onCancel }: DirectionModalProps) {
  const directions = ['north', 'south', 'east', 'west', 'up', 'down']
  const calculatedDirection = calculateDirection(
    connection.sourceX,
    connection.sourceY,
    connection.targetX,
    connection.targetY
  )

  return (
    <div
      style={{
        position: 'fixed',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
        background: '#1a1a1a',
        border: '2px solid #4a8c3f',
        borderRadius: '12px',
        padding: '20px',
        zIndex: 1000,
        minWidth: '280px',
        boxShadow: '0 4px 20px rgba(0,0,0,0.5)',
      }}
    >
      <h3 style={{ marginTop: 0, color: '#fff' }}>Select Direction</h3>
      <p style={{ color: '#aaa', fontSize: '14px' }}>
        Auto-detected: <strong style={{ color: '#4a8c3f' }}>{calculatedDirection}</strong>
      </p>
      <div
        style={{
          display: 'grid',
          gridTemplateColumns: '1fr 1fr',
          gap: '8px',
          marginBottom: '16px',
        }}
      >
        {directions.map((dir) => (
          <label
            key={dir}
            style={{
              display: 'flex',
              alignItems: 'center',
              gap: '8px',
              color: '#fff',
              cursor: 'pointer',
              padding: '8px',
              borderRadius: '6px',
              background: dir === calculatedDirection ? '#2d5a27' : '#333',
            }}
          >
            <input
              type="radio"
              name="direction"
              value={dir}
              defaultChecked={dir === calculatedDirection}
              onChange={() => onSelect(dir)}
            />
            {dir.charAt(0).toUpperCase() + dir.slice(1)}
          </label>
        ))}
      </div>
      <div style={{ display: 'flex', gap: '12px', justifyContent: 'flex-end' }}>
        <button
          onClick={onCancel}
          style={{
            padding: '8px 16px',
            background: '#444',
            color: '#fff',
            border: 'none',
            borderRadius: '6px',
            cursor: 'pointer',
          }}
        >
          Cancel
        </button>
        <button
          onClick={() => onSelect(calculatedDirection)}
          style={{
            padding: '8px 16px',
            background: '#4a8c3f',
            color: '#fff',
            border: 'none',
            borderRadius: '6px',
            cursor: 'pointer',
          }}
        >
          Use Auto
        </button>
      </div>
    </div>
  )
}

interface MapFlowProps {
  initialNodes?: Node[]
  initialEdges?: Edge[]
}

export function MapFlow({ initialNodes = [], initialEdges = [] }: MapFlowProps) {
  const nodeTypes = useMemo(() => ({ roomNode: RoomNode }), [])
  const edgeTypes = useMemo(() => ({ exit: ExitEdge }), [])

  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges)
  const [directionModal, setDirectionModal] = useState<{
    source: string
    target: string
    sourceX: number
    sourceY: number
    targetX: number
    targetY: number
  } | null>(null)

  const onConnect = useCallback(
    (connection: Connection) => {
      if (!connection.source || !connection.target) return

      // Find node positions for direction calculation
      const sourceNode = nodes.find((n) => n.id === connection.source)
      const targetNode = nodes.find((n) => n.id === connection.target)

      if (!sourceNode || !targetNode) return

      // Show direction picker modal
      setDirectionModal({
        source: connection.source,
        target: connection.target,
        sourceX: sourceNode.position.x,
        sourceY: sourceNode.position.y,
        targetX: targetNode.position.x,
        targetY: targetNode.position.y,
      })
    },
    [nodes]
  )

  const handleDirectionSelect = useCallback(
    (direction: string) => {
      if (!directionModal) return

      const oppositeDir = OPPOSITE_DIRECTION[direction]

      // Create bidirectional exits
      const newEdges: Edge[] = [
        {
          id: `${directionModal.source}-${directionModal.target}-${direction}`,
          source: directionModal.source,
          target: directionModal.target,
          type: 'exit',
          data: { direction, isZExit: direction === 'up' || direction === 'down' },
          markerEnd: { type: MarkerType.ArrowClosed, color: '#4a8c3f' },
          style: { stroke: direction === 'up' || direction === 'down' ? '#9b59b6' : '#4a8c3f' },
        },
        {
          id: `${directionModal.target}-${directionModal.source}-${oppositeDir}`,
          source: directionModal.target,
          target: directionModal.source,
          type: 'exit',
          data: { direction: oppositeDir, isZExit: direction === 'up' || direction === 'down' },
          markerEnd: { type: MarkerType.ArrowClosed, color: '#4a8c3f' },
          style: { stroke: direction === 'up' || direction === 'down' ? '#9b59b6' : '#4a8c3f' },
        },
      ]

      setEdges((eds) => addEdge(newEdges, eds))
      setDirectionModal(null)
    },
    [directionModal, setEdges]
  )

  const handleCancelDirection = useCallback(() => {
    setDirectionModal(null)
  }, [])

  const onNodeEdit = useCallback((roomId: string) => {
    console.log('Edit room:', roomId)
    // TODO: Open room edit panel
  }, [])

  // Update initial nodes with edit handler
  const nodesWithHandlers = useMemo(() => {
    return nodes.map((node) => ({
      ...node,
      data: {
        ...node.data,
        onEdit: onNodeEdit,
      },
    }))
  }, [nodes, onNodeEdit])

  return (
    <div style={{ width: '100%', height: '600px', border: '1px solid #333', borderRadius: '8px', overflow: 'hidden' }}>
      {directionModal && (
        <DirectionModal
          connection={directionModal}
          onSelect={handleDirectionSelect}
          onCancel={handleCancelDirection}
        />
      )}
      <ReactFlow
        nodes={nodesWithHandlers}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        onConnect={onConnect}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        fitView
        attributionPosition="bottom-left"
      >
        <Background variant={BackgroundVariant.Dots} gap={20} size={1} color="#333" />
        <Controls />
        <MiniMap
          nodeColor="#2d5a27"
          maskColor="rgba(0, 0, 0, 0.5)"
          style={{ background: '#1a1a1a' }}
        />
      </ReactFlow>
    </div>
  )
}

export default MapFlow