import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/_auth/map')({
  component: MapBuilder,
})

interface MapRoom {
  id: string
  name: string
  x: number
  y: number
}

function MapBuilder() {
  const [rooms] = useState<MapRoom[]>([
    { id: '1', name: 'Town Square', x: 3, y: 3 },
    { id: '2', name: 'Main Street North', x: 3, y: 2 },
    { id: '3', name: 'Main Street South', x: 3, y: 4 },
    { id: '4', name: 'Forest Path', x: 5, y: 3 },
    { id: '5', name: 'Shop District', x: 1, y: 3 },
  ])

  const [selectedRoom, setSelectedRoom] = useState<string | null>(null)
  const gridSize = 8

  const getRoomAt = (x: number, y: number) => rooms.find(r => r.x === x && r.y === y)

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Map Builder</h2>
        <div className="map-actions">
          <button>Add Room</button>
          <button>Connect Rooms</button>
          <button>Save Map</button>
        </div>
      </div>

      <div className="map-container">
        <div className="map-grid" style={{ 
          display: 'grid', 
          gridTemplateColumns: `repeat(${gridSize}, 60px)`,
          gap: '4px'
        }}>
          {Array.from({ length: gridSize * gridSize }, (_, i) => {
            const x = i % gridSize
            const y = Math.floor(i / gridSize)
            const room = getRoomAt(x, y)
            const isSelected = room?.id === selectedRoom
            
            return (
              <div 
                key={i}
                className={`map-cell ${room ? 'has-room' : ''} ${isSelected ? 'selected' : ''}`}
                onClick={() => room && setSelectedRoom(room.id === selectedRoom ? null : room.id)}
                style={{
                  width: '60px',
                  height: '60px',
                  border: '1px solid #444',
                  backgroundColor: room ? '#2d5a27' : '#1a1a1a',
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  fontSize: '10px',
                  cursor: room ? 'pointer' : 'default',
                  position: 'relative'
                }}
              >
                {room && (
                  <>
                    <span style={{ textAlign: 'center', color: '#fff' }}>
                      {room.name.length > 8 ? room.name.slice(0, 6) + '..' : room.name}
                    </span>
                    {/* Exit indicators */}
                    {rooms.some(r => r.x === x - 1 && r.y === y) && (
                      <span style={{ position: 'absolute', left: 2, color: '#666' }}>◀</span>
                    )}
                    {rooms.some(r => r.x === x + 1 && r.y === y) && (
                      <span style={{ position: 'absolute', right: 2, color: '#666' }}>▶</span>
                    )}
                    {rooms.some(r => r.x === x && r.y === y - 1) && (
                      <span style={{ position: 'absolute', top: 2, color: '#666' }}>▲</span>
                    )}
                    {rooms.some(r => r.x === x && r.y === y + 1) && (
                      <span style={{ position: 'absolute', bottom: 2, color: '#666' }}>▼</span>
                    )}
                  </>
                )}
              </div>
            )
          })}
        </div>

        {selectedRoom && (
          <div className="map-sidebar">
            <h3>Room Details</h3>
            {rooms.filter(r => r.id === selectedRoom).map(room => (
              <div key={room.id}>
                <p><strong>Name:</strong> {room.name}</p>
                <p><strong>Position:</strong> ({room.x}, {room.y})</p>
                <div className="room-actions">
                  <button>Edit Room</button>
                  <button>Set Exits</button>
                  <button>Add Items</button>
                  <button>Add NPCs</button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      <div className="map-legend">
        <span><strong>Legend:</strong></span>
        <span>🟩 Room</span>
        <span>⬅️➡️⬆️⬇️ Exits</span>
        <span>🟦 Selected</span>
      </div>
    </div>
  )
}