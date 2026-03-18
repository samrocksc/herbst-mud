import { useState } from 'react'

export interface DirectionPickerProps {
  isOpen: boolean
  sourceId: string
  targetId: string
  sourceName: string
  targetName: string
  onSelect: (direction: string, reverseDirection: string) => void
  onCancel: () => void
}

const DIRECTIONS = [
  { id: 'north', label: 'North', symbol: 'N', reverse: 'south' },
  { id: 'south', label: 'South', symbol: 'S', reverse: 'north' },
  { id: 'east', label: 'East', symbol: 'E', reverse: 'west' },
  { id: 'west', label: 'West', symbol: 'W', reverse: 'east' },
  { id: 'up', label: 'Up', symbol: 'U', reverse: 'down' },
  { id: 'down', label: 'Down', symbol: 'D', reverse: 'up' },
]

export function DirectionPicker({
  isOpen,
  sourceId: _sourceId,
  targetId: _targetId,
  sourceName,
  targetName,
  onSelect,
  onCancel,
}: DirectionPickerProps) {
  const [selectedDirection, setSelectedDirection] = useState<string | null>(null)

  if (!isOpen) return null

  const handleConfirm = () => {
    if (!selectedDirection) return
    const direction = DIRECTIONS.find(d => d.id === selectedDirection)
    if (direction) {
      onSelect(direction.id, direction.reverse)
    }
  }

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      background: 'rgba(0, 0, 0, 0.7)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      zIndex: 1000,
    }}>
      <div style={{
        background: '#1a1a2e',
        border: '2px solid #6c5ce7',
        borderRadius: '12px',
        padding: '24px',
        minWidth: '320px',
        boxShadow: '0 8px 32px rgba(108, 92, 231, 0.4)',
      }}>
        <h3 style={{ 
          margin: '0 0 16px 0', 
          color: '#fff',
          textAlign: 'center',
        }}>
          Select Exit Direction
        </h3>
        
        <div style={{
          marginBottom: '16px',
          padding: '12px',
          background: '#16213e',
          borderRadius: '8px',
          color: '#aaa',
          fontSize: '14px',
        }}>
          <div><strong style={{ color: '#6c5ce7' }}>{sourceName}</strong> to <strong style={{ color: '#6c5ce7' }}>{targetName}</strong></div>
        </div>

        <div style={{
          display: 'grid',
          gridTemplateColumns: 'repeat(3, 1fr)',
          gap: '8px',
          marginBottom: '20px',
        }}>
          {DIRECTIONS.map(dir => (
            <button
              key={dir.id}
              onClick={() => setSelectedDirection(dir.id)}
              style={{
                padding: '12px 8px',
                background: selectedDirection === dir.id ? '#6c5ce7' : '#2d2d44',
                border: selectedDirection === dir.id ? '2px solid #a29bfe' : '2px solid #444',
                borderRadius: '8px',
                color: '#fff',
                cursor: 'pointer',
                transition: 'all 0.2s ease',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                gap: '4px',
              }}
            >
              <span style={{ fontSize: '20px', fontWeight: 'bold' }}>{dir.symbol}</span>
              <span style={{ fontSize: '12px' }}>{dir.label}</span>
            </button>
          ))}
        </div>

        <div style={{
          display: 'flex',
          gap: '12px',
          justifyContent: 'center',
        }}>
          <button
            onClick={onCancel}
            style={{
              padding: '10px 24px',
              background: 'transparent',
              border: '2px solid #666',
              borderRadius: '6px',
              color: '#aaa',
              cursor: 'pointer',
              fontSize: '14px',
            }}
          >
            Cancel
          </button>
          <button
            onClick={handleConfirm}
            disabled={!selectedDirection}
            style={{
              padding: '10px 24px',
              background: selectedDirection ? '#27ae60' : '#333',
              border: 'none',
              borderRadius: '6px',
              color: '#fff',
              cursor: selectedDirection ? 'pointer' : 'not-allowed',
              fontSize: '14px',
              fontWeight: 'bold',
            }}
          >
            Save Exit
          </button>
        </div>
      </div>
    </div>
  )
}