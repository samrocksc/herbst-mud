import { useState } from 'react'

export interface DirectionPickerModalProps {
  sourceName: string
  targetName: string
  onSelect: (sourceDirection: string, targetDirection: string) => void
  onCancel: () => void
}

const DIRECTIONS = [
  { id: 'north', label: 'North', symbol: '⬆️', opposite: 'south' },
  { id: 'south', label: 'South', symbol: '⬇️', opposite: 'north' },
  { id: 'east', label: 'East', symbol: '➡️', opposite: 'west' },
  { id: 'west', label: 'West', symbol: '⬅️', opposite: 'east' },
  { id: 'up', label: 'Up', symbol: '🆙', opposite: 'down' },
  { id: 'down', label: 'Down', symbol: '🔽', opposite: 'up' },
]

export function DirectionPickerModal({
  sourceName,
  targetName,
  onSelect,
  onCancel,
}: DirectionPickerModalProps) {
  const [selectedDirection, setSelectedDirection] = useState<string | null>(null)

  const handleConfirm = () => {
    if (selectedDirection) {
      const dir = DIRECTIONS.find(d => d.id === selectedDirection)
      if (dir) {
        onSelect(dir.id, dir.opposite)
      }
    }
  }

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      backgroundColor: 'rgba(0, 0, 0, 0.7)',
      display: 'flex',
      alignItems: 'center',
      justifyContent: 'center',
      zIndex: 1000,
    }}>
      <div style={{
        background: '#2a2a2a',
        borderRadius: '12px',
        padding: '24px',
        minWidth: '360px',
        maxWidth: '400px',
        boxShadow: '0 8px 32px rgba(0, 0, 0, 0.5)',
        border: '1px solid #444',
      }}>
        <h3 style={{ margin: '0 0 16px 0', color: '#fff', textAlign: 'center' }}>
          Select Exit Direction
        </h3>
        
        <div style={{ 
          display: 'flex', 
          justifyContent: 'space-between', 
          marginBottom: '20px',
          padding: '12px',
          background: '#1a1a1a',
          borderRadius: '8px',
        }}>
          <div style={{ color: '#aaa', fontSize: '14px' }}>
            From: <span style={{ color: '#6c5ce7', fontWeight: 'bold' }}>{sourceName}</span>
          </div>
          <div style={{ color: '#aaa', fontSize: '14px' }}>
            To: <span style={{ color: '#00cec9', fontWeight: 'bold' }}>{targetName}</span>
          </div>
        </div>

        <div style={{ 
          display: 'grid', 
          gridTemplateColumns: 'repeat(3, 1fr)', 
          gap: '12px',
          marginBottom: '24px',
        }}>
          {DIRECTIONS.map((dir) => (
            <button
              key={dir.id}
              onClick={() => setSelectedDirection(dir.id)}
              style={{
                padding: '12px',
                background: selectedDirection === dir.id ? '#6c5ce7' : '#3a3a3a',
                border: selectedDirection === dir.id ? '2px solid #a29bfe' : '2px solid #555',
                borderRadius: '8px',
                color: '#fff',
                cursor: 'pointer',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                gap: '4px',
                transition: 'all 0.2s ease',
              }}
            >
              <span style={{ fontSize: '20px' }}>{dir.symbol}</span>
              <span style={{ fontSize: '13px' }}>{dir.label}</span>
            </button>
          ))}
        </div>

        <div style={{ display: 'flex', gap: '12px', justifyContent: 'flex-end' }}>
          <button
            onClick={onCancel}
            style={{
              padding: '10px 20px',
              background: 'transparent',
              border: '1px solid #666',
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
              padding: '10px 20px',
              background: selectedDirection ? '#27ae60' : '#333',
              border: 'none',
              borderRadius: '6px',
              color: selectedDirection ? '#fff' : '#666',
              cursor: selectedDirection ? 'pointer' : 'not-allowed',
              fontSize: '14px',
              fontWeight: 'bold',
            }}
          >
            Connect (Bidirectional)
          </button>
        </div>
      </div>
    </div>
  )
}