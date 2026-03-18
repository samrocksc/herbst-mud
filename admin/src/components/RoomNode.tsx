import { Handle, Position, NodeProps } from '@xyflow/react'

interface RoomNodeData {
  label: string
  roomId: string
  exits?: string[]
  onEdit?: (roomId: string) => void
}

export type RoomNodeType = NodeProps & {
  data: RoomNodeData
}

export function RoomNode({ data }: RoomNodeType) {
  const { label, roomId, exits = [], onEdit } = data

  return (
    <div 
      className="room-node"
      onClick={() => onEdit?.(roomId)}
      style={{
        background: '#2d5a27',
        border: '2px solid #4a8c3f',
        borderRadius: '8px',
        padding: '12px',
        minWidth: '120px',
        color: '#fff',
        cursor: 'pointer',
        position: 'relative',
      }}
    >
      <Handle 
        type="target" 
        position={Position.Left}
        style={{ background: '#4a8c3f', width: 10, height: 10 }}
        className="connection-handle"
      />
      <Handle 
        type="target" 
        position={Position.Top}
        style={{ background: '#4a8c3f', width: 10, height: 10 }}
        className="connection-handle"
      />
      <Handle 
        type="target" 
        position={Position.Bottom}
        style={{ background: '#4a8c3f', width: 10, height: 10 }}
        className="connection-handle"
      />
      
      <div style={{ fontWeight: 'bold', fontSize: '14px', marginBottom: 4 }}>
        {label}
      </div>
      
      {exits.length > 0 && (
        <div style={{ fontSize: '10px', color: '#aaa', marginTop: 4 }}>
          {exits.join(' ')}
        </div>
      )}
      
      <Handle 
        type="source" 
        position={Position.Right}
        style={{ background: '#4a8c3f', width: 10, height: 10 }}
        className="connection-handle"
      />
      <Handle 
        type="source" 
        position={Position.Top}
        style={{ background: '#4a8c3f', width: 10, height: 10 }}
        className="connection-handle"
      />
      <Handle 
        type="source" 
        position={Position.Bottom}
        style={{ background: '#4a8c3f', width: 10, height: 10 }}
        className="connection-handle"
      />
    </div>
  )
}

export default RoomNode