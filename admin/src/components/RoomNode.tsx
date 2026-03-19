import { memo, useState } from 'react'
import { Handle, Position } from '@xyflow/react'
import type { NodeProps } from '@xyflow/react'

export interface RoomNodeData {
  name: string
  description: string
  zLevel: number
  exits?: Record<string, number>
}

const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text
  return text.slice(0, maxLength) + '...'
}

function RoomNodeComponent({ data, selected }: NodeProps) {
  const roomData = data as unknown as RoomNodeData
  const showZLevel = roomData.zLevel !== 0
  const [isHovered, setIsHovered] = useState(false)
  
  const handleStyle = {
    background: '#6c5ce7',
    width: '10px',
    height: '10px',
    border: '2px solid #fff',
    transition: 'all 0.2s ease',
  }
  
  return (
    <div 
      style={{
        padding: '12px 16px',
        borderRadius: '8px',
        background: selected ? '#1a3a2f' : '#2d5a27',
        border: selected ? '2px solid #6c5ce7' : '2px solid #3a7a3a',
        minWidth: '160px',
        maxWidth: '200px',
        boxShadow: selected ? '0 0 12px rgba(108, 92, 231, 0.5)' : '0 2px 8px rgba(0,0,0,0.3)',
        cursor: 'grab',
        transition: 'all 0.2s ease',
      }}
      onMouseEnter={() => setIsHovered(true)}
      onMouseLeave={() => setIsHovered(false)}
    >
      {/* Handles for connections - shown on hover */}
      <Handle type="target" position={Position.Top} style={{ ...handleStyle, opacity: isHovered ? 1 : 0.3 }} />
      <Handle type="source" position={Position.Bottom} style={{ ...handleStyle, opacity: isHovered ? 1 : 0.3 }} />
      <Handle type="target" position={Position.Left} id="left-target" style={{ ...handleStyle, opacity: isHovered ? 1 : 0.3 }} />
      <Handle type="source" position={Position.Left} id="left-source" style={{ ...handleStyle, opacity: isHovered ? 1 : 0.3 }} />
      <Handle type="target" position={Position.Right} id="right-target" style={{ ...handleStyle, opacity: isHovered ? 1 : 0.3 }} />
      <Handle type="source" position={Position.Right} id="right-source" style={{ ...handleStyle, opacity: isHovered ? 1 : 0.3 }} />
      
      {/* Room name - bold */}
      <div style={{
        fontWeight: 'bold',
        color: '#fff',
        marginBottom: '4px',
        fontSize: '14px',
      }}>
        {roomData.name}
      </div>
      
      {/* Description - truncated */}
      <div style={{
        color: '#aaa',
        fontSize: '12px',
        marginBottom: showZLevel ? '6px' : '0',
      }}>
        {truncateText(roomData.description, 25)}
      </div>
      
      {/* Z-level badge (only shown if not 0) */}
      {showZLevel && (
        <div style={{
          display: 'inline-block',
          padding: '2px 8px',
          background: roomData.zLevel > 0 ? '#4a69bd' : '#b33939',
          borderRadius: '4px',
          fontSize: '11px',
          color: '#fff',
          fontWeight: 'bold',
        }}>
          Z: {roomData.zLevel}
        </div>
      )}
    </div>
  )
}

export const RoomNode = memo(RoomNodeComponent)