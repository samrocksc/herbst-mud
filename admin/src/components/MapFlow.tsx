import { ReactFlow, Background, Controls, MiniMap } from '@xyflow/react'
import type { Node, Edge, Connection } from '@xyflow/react'
import { RoomNode } from './RoomNode'
import { ExitEdge, ExitEdgeData } from './ExitEdge'
import '@xyflow/react/dist/style.css'

const nodeTypes = {
  room: RoomNode,
}

const edgeTypes = {
  exit: ExitEdge,
}

export interface MapFlowProps {
  nodes: Node[]
  edges: Edge[]
  onNodesChange?: (changes: unknown) => void
  onEdgesChange?: (changes: unknown) => void
  onConnect?: (connection: Connection) => void
  onNodeClick?: (event: React.MouseEvent, node: Node) => void
}

export function MapFlow({
  nodes,
  edges,
  onNodesChange,
  onEdgesChange,
  onConnect,
  onNodeClick
}: MapFlowProps) {
  return (
    <div style={{ width: '100%', height: '600px', border: '1px solid #333', borderRadius: '8px', overflow: 'hidden' }}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        nodeTypes={nodeTypes}
        edgeTypes={edgeTypes}
        defaultEdgeType="exit"
        onNodesChange={onNodesChange as never}
        onEdgesChange={onEdgesChange as never}
        onConnect={onConnect as never}
        onNodeClick={onNodeClick as never}
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