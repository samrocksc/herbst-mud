import { useState, useEffect } from 'react'
import type { Node } from '@xyflow/react'
import { DIRECTIONS, type Direction } from '../hooks/useRooms'

interface Exit {
  direction: string
  targetRoomId: number
}

interface RoomEditPanelProps {
  selectedNode: Node | null
  rooms: Array<{ id: number; name: string }>
  onUpdate: (id: string, data: { name: string; description: string; zLevel: number; exits: Record<string, number> }) => void
  onDelete: (id: string) => void
  onClose: () => void
}

export function RoomEditPanel({ selectedNode, rooms, onUpdate, onDelete, onClose }: RoomEditPanelProps) {
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [zLevel, setZLevel] = useState(0)
  const [exits, setExits] = useState<Exit[]>([])
  const [showExitForm, setShowExitForm] = useState(false)
  const [newExitDirection, setNewExitDirection] = useState<Direction>('north')
  const [newExitTarget, setNewExitTarget] = useState<number>(0)

  // Update form when selection changes
  useEffect(() => {
    if (selectedNode && selectedNode.data) {
      const data = selectedNode.data as {
        name?: string
        description?: string
        zLevel?: number
        exits?: Record<string, number>
      }
      setName(data.name || '')
      setDescription(data.description || '')
      setZLevel(data.zLevel ?? 0)
      
      // Convert exits object to array
      const exitsArray: Exit[] = data.exits
        ? Object.entries(data.exits).map(([direction, targetRoomId]) => ({
            direction,
            targetRoomId
          }))
        : []
      setExits(exitsArray)
    }
  }, [selectedNode])

  if (!selectedNode) {
    return null
  }

  const handleAddExit = () => {
    if (!newExitTarget) return
    
    // Check if direction already exists
    if (exits.some(e => e.direction === newExitDirection)) {
      alert(`Exit "${newExitDirection}" already exists`)
      return
    }
    
    setExits([...exits, { direction: newExitDirection, targetRoomId: newExitTarget }])
    setShowExitForm(false)
    setNewExitTarget(0)
  }

  const handleRemoveExit = (direction: string) => {
    setExits(exits.filter(e => e.direction !== direction))
  }

  const handleSave = () => {
    const exitsObj: Record<string, number> = {}
    exits.forEach(e => {
      exitsObj[e.direction] = e.targetRoomId
    })
    
    onUpdate(selectedNode.id, {
      name,
      description,
      zLevel,
      exits: exitsObj
    })
  }

  const getRoomName = (id: number) => {
    const room = rooms.find(r => r.id === id)
    return room ? room.name : `Room ${id}`
  }

  return (
    <div className="room-edit-panel" style={{
      width: '320px',
      padding: '16px',
      background: '#1a1a2e',
      borderRadius: '8px',
      border: '1px solid #444',
      overflowY: 'auto',
      maxHeight: 'calc(100vh - 200px)'
    }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
        <h3 style={{ margin: 0, color: '#e0e0e0' }}>Edit Room</h3>
        <button 
          onClick={onClose}
          style={{ 
            background: 'transparent', 
            border: 'none', 
            color: '#888', 
            fontSize: '20px', 
            cursor: 'pointer' 
          }}
          aria-label="Close panel"
        >
          ×
        </button>
      </div>

      {/* Name Field */}
      <div style={{ marginBottom: '12px' }}>
        <label style={{ display: 'block', marginBottom: '4px', color: '#aaa' }}>
          Name:
        </label>
        <input 
          type="text" 
          value={name}
          onChange={(e) => setName(e.target.value)}
          style={{ 
            width: '100%', 
            padding: '8px', 
            background: '#2a2a4a', 
            border: '1px solid #555', 
            borderRadius: '4px',
            color: '#fff',
            fontSize: '14px'
          }}
          placeholder="Room name"
        />
      </div>

      {/* Description Field */}
      <div style={{ marginBottom: '12px' }}>
        <label style={{ display: 'block', marginBottom: '4px', color: '#aaa' }}>
          Description:
        </label>
        <textarea 
          value={description}
          onChange={(e) => setDescription(e.target.value)}
          style={{ 
            width: '100%', 
            padding: '8px', 
            background: '#2a2a4a', 
            border: '1px solid #555', 
            borderRadius: '4px',
            color: '#fff',
            fontSize: '14px',
            minHeight: '80px',
            resize: 'vertical'
          }}
          placeholder="Room description"
        />
      </div>

      {/* Z-Level Selector */}
      <div style={{ marginBottom: '12px' }}>
        <label style={{ display: 'block', marginBottom: '4px', color: '#aaa' }}>
          Z-Level:
        </label>
        <select 
          value={zLevel}
          onChange={(e) => setZLevel(parseInt(e.target.value))}
          style={{ 
            width: '100%', 
            padding: '8px', 
            background: '#2a2a4a', 
            border: '1px solid #555', 
            borderRadius: '4px',
            color: '#fff',
            fontSize: '14px'
          }}
        >
          <option value={-2}>Z: -2 (Deep Underground)</option>
          <option value={-1}>Z: -1 (Underground)</option>
          <option value={0}>Z: 0 (Ground)</option>
          <option value={1}>Z: 1 (Upper Floor)</option>
          <option value={2}>Z: 2 (Tower)</option>
        </select>
      </div>

      {/* Exits Section */}
      <div style={{ marginBottom: '16px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '8px' }}>
          <label style={{ color: '#aaa' }}>Exits:</label>
          <button 
            onClick={() => setShowExitForm(!showExitForm)}
            style={{ 
              padding: '4px 8px',
              background: '#4a4a7a',
              border: '1px solid #666',
              borderRadius: '4px',
              color: '#fff',
              cursor: 'pointer',
              fontSize: '12px'
            }}
          >
            + Add Exit
          </button>
        </div>

        {/* Add Exit Form */}
        {showExitForm && (
          <div style={{ 
            padding: '8px', 
            background: '#2a2a4a', 
            borderRadius: '4px', 
            marginBottom: '8px',
            border: '1px solid #555'
          }}>
            <div style={{ marginBottom: '8px' }}>
              <select 
                value={newExitDirection}
                onChange={(e) => setNewExitDirection(e.target.value as Direction)}
                style={{ 
                  width: '100%', 
                  padding: '6px', 
                  background: '#1a1a2e', 
                  border: '1px solid #444', 
                  borderRadius: '4px',
                  color: '#fff'
                }}
              >
                {DIRECTIONS.map(dir => (
                  <option key={dir} value={dir}>{dir.charAt(0).toUpperCase() + dir.slice(1)}</option>
                ))}
              </select>
            </div>
            <div style={{ marginBottom: '8px' }}>
              <select 
                value={newExitTarget}
                onChange={(e) => setNewExitTarget(parseInt(e.target.value))}
                style={{ 
                  width: '100%', 
                  padding: '6px', 
                  background: '#1a1a2e', 
                  border: '1px solid #444', 
                  borderRadius: '4px',
                  color: '#fff'
                }}
              >
                <option value={0}>Select target room...</option>
                {rooms.filter(r => String(r.id) !== selectedNode.id).map(room => (
                  <option key={room.id} value={room.id}>
                    {room.name}
                  </option>
                ))}
              </select>
            </div>
            <div style={{ display: 'flex', gap: '8px' }}>
              <button 
                onClick={handleAddExit}
                disabled={!newExitTarget}
                style={{ 
                  flex: 1,
                  padding: '6px',
                  background: newExitTarget ? '#2d5a27' : '#333',
                  border: '1px solid #444',
                  borderRadius: '4px',
                  color: newExitTarget ? '#fff' : '#666',
                  cursor: newExitTarget ? 'pointer' : 'not-allowed'
                }}
              >
                Add
              </button>
              <button 
                onClick={() => setShowExitForm(false)}
                style={{ 
                  flex: 1,
                  padding: '6px',
                  background: '#4a3a3a',
                  border: '1px solid #444',
                  borderRadius: '4px',
                  color: '#fff',
                  cursor: 'pointer'
                }}
              >
                Cancel
              </button>
            </div>
          </div>
        )}

        {/* Exits List */}
        <div style={{ maxHeight: '120px', overflowY: 'auto' }}>
          {exits.length === 0 ? (
            <p style={{ color: '#666', fontSize: '12px', margin: '4px 0' }}>No exits configured</p>
          ) : (
            exits.map(exit => (
              <div 
                key={exit.direction}
                style={{ 
                  display: 'flex', 
                  justifyContent: 'space-between', 
                  alignItems: 'center',
                  padding: '6px 8px',
                  background: '#2a2a4a',
                  borderRadius: '4px',
                  marginBottom: '4px'
                }}
              >
                <span style={{ color: '#e0e0e0' }}>
                  <strong>{exit.direction}</strong> → {getRoomName(exit.targetRoomId)}
                </span>
                <button 
                  onClick={() => handleRemoveExit(exit.direction)}
                  style={{ 
                    background: 'transparent',
                    border: 'none',
                    color: '#e55',
                    cursor: 'pointer',
                    fontSize: '14px'
                  }}
                  aria-label={`Remove ${exit.direction} exit`}
                >
                  ×
                </button>
              </div>
            ))
          )}
        </div>
      </div>

      {/* Room ID (Info) */}
      <div style={{ 
        padding: '8px', 
        background: '#2a2a4a', 
        borderRadius: '4px', 
        marginBottom: '12px',
        fontSize: '12px',
        color: '#888'
      }}>
        <span>Node ID: {selectedNode.id}</span>
      </div>

      {/* Action Buttons */}
      <div style={{ display: 'flex', gap: '8px' }}>
        <button 
          onClick={handleSave}
          style={{ 
            flex: 2,
            padding: '10px',
            background: '#2d5a27',
            border: '1px solid #3a7a3a',
            borderRadius: '4px',
            color: '#fff',
            cursor: 'pointer',
            fontWeight: 'bold'
          }}
        >
          Save Changes
        </button>
        <button 
          onClick={() => onDelete(selectedNode.id)}
          style={{ 
            flex: 1,
            padding: '10px',
            background: '#5a2d2d',
            border: '1px solid #7a3a3a',
            borderRadius: '4px',
            color: '#fff',
            cursor: 'pointer'
          }}
        >
          Delete
        </button>
      </div>
    </div>
  )
}