import { useCallback, useMemo } from 'react'
import { ReactFlow, Background, Controls, MiniMap, useReactFlow, ScreenFlowPosition, ReactFlowProvider } from '@xyflow/react'
import type { Node, Edge, Connection } from '@xyflow/react'
import { RoomNode } from './RoomNode'
import '@xyflow/react/dist/style.css'

const nodeTypes = {
  room: RoomNode,
}

export interface MapFlowProps {
  nodes: Node[]
  edges: Edge[]
  onNodesChange?: (changes: unknown) => void
  onEdgesChange?: (changes: unknown) => void
  onConnect?: (connection: Connection) => void
  onNodeClick?: (event: React.MouseEvent, node: Node) => void
  onDrop?: (position: { x: number; y: number }) => void
  onDragOver?: (event: React.DragEvent) => void
}

function MapFlowInner({
  nodes,
  edges,
  onNodesChange,
  onEdgesChange,
  onConnect,
  onNodeClick,
  onDrop,
  onDragOver
}: MapFlowProps) {
  const { screenToFlowPosition } = useReactFlow()

  const handleDrop = useCallback((event: React.DragEvent) => {
    event.preventDefault()
    const type = event.dataTransfer.getData('application/reactflow')
    if (!type || type !== 'room') return

    // Get drop position in flow coordinates
    const position = screenToFlowPosition({
      x: event.clientX,
      y: event.clientY
    })

    if (onDrop) {
      onDrop(position)
    }
  }, [screenToFlowPosition, onDrop])

  const handleDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
    if (onDragOver) {
      onDragOver(event)
    }
  }, [onDragOver])

  return (
    <div style={{ width: '100%', height: '600px', border: '1px solid #333', borderRadius: '8px', overflow: 'hidden' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        onNodesChange={onNodesChange as never}
        onEdgesChange={onEdgesChange as never}
        onConnect={onConnect as never}
        onNodeClick={onNodeClick as never}
        onDrop={handleDrop}
        onDragOver={handleDragOver}
        fitView
        attributionPosition="bottom-left"
      >
        <Background color="#444" gap={16} />
        <Controls />
        <MiniMap 
          nodeColor={(node) => {
            return node.selected ? '#6c5ce7' : '#2d5a27'
          }}
          maskColor="rgba(0, 0, 0, 0.3)"
        />
      </ReactFlow>
    </div>
  )
}

export function MapFlow(props: MapFlowProps) {
  return (
    <ReactFlowProvider>
      <MapFlowInner {...props} />
    </ReactFlowProvider>
  )
}