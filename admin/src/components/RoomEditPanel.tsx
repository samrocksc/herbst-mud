import { useState, useEffect, useCallback } from 'react'
import type { Node } from '@xyflow/react'

export interface RoomEditPanelProps {
  selectedNode: Node | null
  onUpdateNode: (nodeId: string, data: Record<string, unknown>) => void
  onDeleteNode: (nodeId: string) => void
  onClose: () => void
  edges: { id: string; source: string; target: string; label?: string }[]
}

interface RoomData extends Record<string, unknown> {
  name: string
  description: string
  zLevel: number
  exits?: Record<string, number>
}

export function RoomEditPanel({ 
  selectedNode, 
  onUpdateNode, 
  onDeleteNode, 
  onClose,
  edges 
}: RoomEditPanelProps) {
  const [localData, setLocalData] = useState<RoomData>({
    name: '',
    description: '',
    zLevel: 0,
    exits: {},
  })
  const [isSaving, setIsSaving] = useState(false)
  const [showDeleteConfirm, setShowDeleteConfirm] = useState(false)

  // Sync local state with selected node
  useEffect(() => {
    if (selectedNode?.data) {
      setLocalData(selectedNode.data as unknown as RoomData)
    }
  }, [selectedNode])

  // Get exits for this room
  const roomExits = edges
    .filter(e => e.source === selectedNode?.id || e.target === selectedNode?.id)
    .map(e => {
      const isSource = e.source === selectedNode?.id
      return {
        id: e.id,
        direction: e.label || (isSource ? 'out' : 'in'),
        targetRoom: isSource ? e.target : e.source,
      }
    })

  const handleSave = useCallback(async () => {
    if (!selectedNode) return
    
    setIsSaving(true)
    try {
      // Update the node with local data
      onUpdateNode(selectedNode.id, localData)
      
      // TODO: API call to save room
      // await fetch(`/api/rooms/${selectedNode.id}`, {
      //   method: 'PUT',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify(localData),
      // })
      
      alert('Room saved successfully!')
    } catch (error) {
      console.error('Failed to save room:', error)
      alert('Failed to save room')
    } finally {
      setIsSaving(false)
    }
  }, [selectedNode, localData, onUpdateNode])

  const handleDelete = useCallback(() => {
    if (!selectedNode) return
    onDeleteNode(selectedNode.id)
    setShowDeleteConfirm(false)
  }, [selectedNode, onDeleteNode])

  const handleAddExit = useCallback(() => {
    // TODO: Implement add exit - show direction picker
    alert('Add exit - use the connection handles on the node')
  }, [])

  const handleRemoveExit = useCallback((_exitId: string) => {
    // TODO: Implement remove exit
    alert('Remove exit functionality coming soon')
  }, [])

  if (!selectedNode) return null

  return (
    <div className="room-edit-panel" style={{
      width: '300px',
      padding: '16px',
      background: '#1e1e2e',
      borderRadius: '8px',
      border: '1px solid #45475a',
      color: '#cdd6f4',
      maxHeight: 'calc(100vh - 200px)',
      overflowY: 'auto',
    }}>
      {/* Header */}
      <div style={{
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        marginBottom: '16px',
        paddingBottom: '12px',
        borderBottom: '1px solid #45475a',
      }}>
        <h3 style={{ margin: 0, color: '#cba6f7' }}>Room Details</h3>
        <button
          onClick={onClose}
          style={{
            background: 'transparent',
            border: 'none',
            color: '#a6adc8',
            fontSize: '18px',
            cursor: 'pointer',
            padding: '4px 8px',
          }}
          title="Close panel"
        >
          ✕
        </button>
      </div>

      {/* Name */}
      <div style={{ marginBottom: '16px' }}>
        <label style={{ 
          display: 'block', 
          marginBottom: '6px', 
          fontSize: '13px',
          color: '#a6adc8',
          fontWeight: 500,
        }}>
          Name
        </label>
        <input
          type="text"
          value={localData.name}
          onChange={(e) => setLocalData(prev => ({ ...prev, name: e.target.value }))}
          style={{
            width: '100%',
            padding: '10px 12px',
            background: '#313244',
            border: '1px solid #45475a',
            borderRadius: '6px',
            color: '#cdd6f4',
            fontSize: '14px',
          }}
        />
      </div>

      {/* Description */}
      <div style={{ marginBottom: '16px' }}>
        <label style={{ 
          display: 'block', 
          marginBottom: '6px', 
          fontSize: '13px',
          color: '#a6adc8',
          fontWeight: 500,
        }}>
          Description
        </label>
        <textarea
          value={localData.description}
          onChange={(e) => setLocalData(prev => ({ ...prev, description: e.target.value }))}
          rows={3}
          style={{
            width: '100%',
            padding: '10px 12px',
            background: '#313244',
            border: '1px solid #45475a',
            borderRadius: '6px',
            color: '#cdd6f4',
            fontSize: '14px',
            resize: 'vertical',
          }}
        />
      </div>

      {/* Z-Level */}
      <div style={{ marginBottom: '16px' }}>
        <label style={{ 
          display: 'block', 
          marginBottom: '6px', 
          fontSize: '13px',
          color: '#a6adc8',
          fontWeight: 500,
        }}>
          Z-Level
        </label>
        <select
          value={localData.zLevel}
          onChange={(e) => setLocalData(prev => ({ ...prev, zLevel: parseInt(e.target.value) }))}
          style={{
            width: '100%',
            padding: '10px 12px',
            background: '#313244',
            border: '1px solid #45475a',
            borderRadius: '6px',
            color: '#cdd6f4',
            fontSize: '14px',
          }}
        >
          <option value={-2}>Z: -2 (Deep Underground)</option>
          <option value={-1}>Z: -1 (Underground)</option>
          <option value={0}>Z: 0 (Ground)</option>
          <option value={1}>Z: 1 (Upper Floor)</option>
          <option value={2}>Z: 2 (Tower)</option>
        </select>
      </div>

      {/* Exits List */}
      <div style={{ marginBottom: '16px' }}>
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          alignItems: 'center',
          marginBottom: '8px',
        }}>
          <label style={{ 
            fontSize: '13px',
            color: '#a6adc8',
            fontWeight: 500,
          }}>
            Exits ({roomExits.length})
          </label>
          <button
            onClick={handleAddExit}
            style={{
              background: '#89b4fa',
              border: 'none',
              borderRadius: '4px',
              color: '#1e1e2e',
              padding: '4px 10px',
              fontSize: '12px',
              cursor: 'pointer',
              fontWeight: 500,
            }}
          >
            + Add Exit
          </button>
        </div>
        
        {roomExits.length > 0 ? (
          <div style={{
            background: '#313244',
            borderRadius: '6px',
            border: '1px solid #45475a',
            maxHeight: '120px',
            overflowY: 'auto',
          }}>
            {roomExits.map(exit => (
              <div
                key={exit.id}
                style={{
                  display: 'flex',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  padding: '8px 12px',
                  borderBottom: '1px solid #45475a',
                }}
              >
                <span style={{ color: '#f38ba8', fontWeight: 500 }}>
                  {exit.direction.toUpperCase()}
                </span>
                <span style={{ color: '#a6adc8', fontSize: '13px' }}>
                  → Room {exit.targetRoom}
                </span>
                <button
                  onClick={() => handleRemoveExit(exit.id)}
                  style={{
                    background: 'transparent',
                    border: 'none',
                    color: '#f38ba8',
                    cursor: 'pointer',
                    fontSize: '14px',
                    padding: '2px 6px',
                  }}
                  title="Remove exit"
                >
                  ✕
                </button>
              </div>
            ))}
          </div>
        ) : (
          <div style={{
            background: '#313244',
            borderRadius: '6px',
            border: '1px solid #45475a',
            padding: '16px',
            textAlign: 'center',
            color: '#6c7086',
            fontSize: '13px',
          }}>
            No exits yet. Connect rooms using handles.
          </div>
        )}
      </div>

      {/* Node ID */}
      <div style={{ 
        marginBottom: '20px', 
        padding: '8px', 
        background: '#313244', 
        borderRadius: '4px',
        fontSize: '12px',
        color: '#6c7086',
      }}>
        <strong>Node ID:</strong> {selectedNode.id}
      </div>

      {/* Action Buttons */}
      <div style={{ display: 'flex', gap: '8px' }}>
        <button
          onClick={handleSave}
          disabled={isSaving}
          style={{
            flex: 1,
            padding: '12px',
            background: isSaving ? '#6c7086' : '#a6e3a1',
            border: 'none',
            borderRadius: '6px',
            color: '#1e1e2e',
            fontSize: '14px',
            fontWeight: 600,
            cursor: isSaving ? 'not-allowed' : 'pointer',
          }}
        >
          {isSaving ? 'Saving...' : 'Save'}
        </button>
        
        {!showDeleteConfirm ? (
          <button
            onClick={() => setShowDeleteConfirm(true)}
            style={{
              padding: '12px 16px',
              background: 'transparent',
              border: '1px solid #f38ba8',
              borderRadius: '6px',
              color: '#f38ba8',
              fontSize: '14px',
              cursor: 'pointer',
            }}
          >
            Delete
          </button>
        ) : (
          <div style={{ display: 'flex', gap: '4px' }}>
            <button
              onClick={handleDelete}
              style={{
                padding: '12px',
                background: '#f38ba8',
                border: 'none',
                borderRadius: '6px',
                color: '#1e1e2e',
                fontSize: '12px',
                fontWeight: 600,
                cursor: 'pointer',
              }}
            >
              Confirm
            </button>
            <button
              onClick={() => setShowDeleteConfirm(false)}
              style={{
                padding: '12px',
                background: 'transparent',
                border: '1px solid #6c7086',
                borderRadius: '6px',
                color: '#a6adc8',
                fontSize: '12px',
                cursor: 'pointer',
              }}
            >
              Cancel
            </button>
          </div>
        )}
      </div>
    </div>
  )
}