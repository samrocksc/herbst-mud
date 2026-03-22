import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useEffect, useState } from 'react'

export const Route = createFileRoute('/map')({
  component: MapBuilder,
})

interface Room {
  id: number
  name: string
  description: string
  isStartingRoom?: boolean
  exits: Record<string, number>
  x?: number
  y?: number
  zLevel?: number
}

function MapBuilder() {
  const navigate = useNavigate()
  const [rooms, setRooms] = useState<Room[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [selectedRoom, setSelectedRoom] = useState<Room | null>(null)
  const [zoom, setZoom] = useState(1)

  // Check authentication
  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate({ to: '/login' })
    }
  }, [navigate])

  // Load rooms from API
  useEffect(() => {
    fetch('http://localhost:8080/rooms')
      .then(res => res.json())
      .then(data => {
        setRooms(data)
        setLoading(false)
      })
      .catch(err => {
        setError(err.message)
        setLoading(false)
      })
  }, [])

  // Find starting room (Fountain Plaza) or first room
  const centerRoom = rooms.find(r => r.isStartingRoom) || rooms[0]

  // Build a graph-based layout
  const nodePositions: Map<number, { x: number; y: number }> = new Map()
  const visited = new Set<number>()

  const positionRoom = (roomId: number, x: number, y: number, fromDir?: string) => {
    if (visited.has(roomId)) return
    visited.add(roomId)
    nodePositions.set(roomId, { x, y })

    const room = rooms.find(r => r.id === roomId)
    if (!room) return

    const directionOffsets: Record<string, { dx: number; dy: number }> = {
      north: { dx: 0, dy: -120 },
      south: { dx: 0, dy: 120 },
      east: { dx: 150, dy: 0 },
      west: { dx: -150, dy: 0 },
      up: { dx: 100, dy: -80 },
      down: { dx: -100, dy: 80 }
    }

    for (const [dir, targetId] of Object.entries(room.exits || {})) {
      if (targetId && !visited.has(targetId)) {
        const offset = directionOffsets[dir] || { dx: 150, dy: 0 }
        positionRoom(targetId, x + offset.dx, y + offset.dy, dir)
      }
    }
  }

  // Position all rooms starting from center
  if (centerRoom) {
    positionRoom(centerRoom.id, 400, 300)
  }

  // Position remaining rooms that weren't connected
  let orphanX = 800
  for (const room of rooms) {
    if (!visited.has(room.id)) {
      nodePositions.set(room.id, { x: orphanX, y: 300 })
      orphanX += 150
    }
  }

  // Get all exits for a room
  const getExitLabels = (room: Room): string => {
    const exits = Object.entries(room.exits || {})
    if (exits.length === 0) return '-'
    return exits.map(([dir, id]) => `${dir.charAt(0).toUpperCase()}→${id}`).join(' ')
  }

  if (loading) {
    return <div style={{ padding: '2rem', color: '#fff' }}>Loading map...</div>
  }

  if (error) {
    return <div style={{ padding: '2rem', color: '#ff6b6b' }}>Error: {error}</div>
  }

  return (
    <div style={{ display: 'flex', height: '100vh', background: '#0a0a0f' }}>
      {/* Map Area */}
      <div style={{ flex: 1, overflow: 'hidden', position: 'relative' }}>
        {/* Header */}
        <div style={{
          position: 'absolute',
          top: 0,
          left: 0,
          right: 0,
          padding: '1rem',
          background: 'rgba(26, 26, 46, 0.95)',
          borderBottom: '1px solid #333',
          display: 'flex',
          justifyContent: 'space-between',
          alignItems: 'center',
          zIndex: 10
        }}>
          <h1 style={{ margin: 0, color: '#fff', fontSize: '1.5rem' }}>
            🗺️ Map Builder - {rooms.length} Rooms
          </h1>
          <div style={{ display: 'flex', gap: '1rem', alignItems: 'center' }}>
            <button
              onClick={() => setZoom(z => Math.min(z + 0.25, 2))}
              style={{ padding: '0.25rem 0.5rem', background: '#27ae60', border: 'none', borderRadius: '4px', color: '#fff', cursor: 'pointer' }}
            >
              +
            </button>
            <span style={{ color: '#888' }}>{Math.round(zoom * 100)}%</span>
            <button
              onClick={() => setZoom(z => Math.max(z - 0.25, 0.5))}
              style={{ padding: '0.25rem 0.5rem', background: '#e74c3c', border: 'none', borderRadius: '4px', color: '#fff', cursor: 'pointer' }}
            >
              −
            </button>
          </div>
        </div>

        {/* Map Canvas */}
        <div style={{
          marginTop: '60px',
          height: 'calc(100% - 60px)',
          overflow: 'auto',
          padding: '2rem'
        }}>
          <div style={{
            position: 'relative',
            width: '2000px',
            height: '2000px',
            transform: `scale(${zoom})`,
            transformOrigin: 'top left'
          }}>
            {/* Draw edges (connections) */}
            <svg style={{ position: 'absolute', top: 0, left: 0, width: '100%', height: '100%', pointerEvents: 'none' }}>
              {rooms.map(room => {
                const pos = nodePositions.get(room.id)
                if (!pos) return null

                return Object.entries(room.exits || {}).map(([dir, targetId]) => {
                  const targetPos = nodePositions.get(targetId)
                  if (!targetPos) return null

                  const isZExit = dir === 'up' || dir === 'down'
                  const color = isZExit
                    ? (dir === 'up' ? '#e17055' : '#74b9ff')
                    : '#666'

                  return (
                    <line
                      key={`${room.id}-${targetId}-${dir}`}
                      x1={pos.x + 50}
                      y1={pos.y + 30}
                      x2={targetPos.x + 50}
                      y2={targetPos.y + 30}
                      stroke={color}
                      strokeWidth={isZExit ? 3 : 2}
                      strokeDasharray={isZExit ? '5,5' : undefined}
                      markerEnd="url(#arrowhead)"
                    />
                  )
                })
              })}
            </svg>

            {/* Draw room nodes */}
            {rooms.map(room => {
              const pos = nodePositions.get(room.id)
              if (!pos) return null

              const isSelected = selectedRoom?.id === room.id

              return (
                <div
                  key={room.id}
                  onClick={() => setSelectedRoom(room)}
                  style={{
                    position: 'absolute',
                    left: pos.x,
                    top: pos.y,
                    width: '100px',
                    minHeight: '60px',
                    background: room.isStartingRoom ? '#27ae60' : isSelected ? '#6c5ce7' : '#2d5a27',
                    border: `2px solid ${isSelected ? '#a29bfe' : '#1a3a1a'}`,
                    borderRadius: '8px',
                    padding: '0.5rem',
                    cursor: 'pointer',
                    display: 'flex',
                    flexDirection: 'column',
                    justifyContent: 'center',
                    alignItems: 'center',
                    boxShadow: isSelected ? '0 0 15px rgba(108, 92, 231, 0.5)' : '0 2px 8px rgba(0,0,0,0.3)',
                    transition: 'all 0.2s',
                    zIndex: isSelected ? 10 : 1
                  }}
                >
                  <div style={{
                    color: '#fff',
                    fontWeight: 'bold',
                    fontSize: '0.75rem',
                    textAlign: 'center',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    width: '100%'
                  }}>
                    {room.name}
                  </div>
                  <div style={{
                    color: '#888',
                    fontSize: '0.65rem'
                  }}>
                    #{room.id}
                    {room.isStartingRoom && ' ⭐'}
                  </div>
                </div>
              )
            })}
          </div>
        </div>
      </div>

      {/* Room Details Panel */}
      <div style={{
        width: '320px',
        background: '#1a1a2e',
        borderLeft: '1px solid #333',
        padding: '1rem',
        overflowY: 'auto'
      }}>
        {selectedRoom ? (
          <>
            <h3 style={{ margin: '0 0 1rem 0', color: '#fff' }}>
              {selectedRoom.name}
              {selectedRoom.isStartingRoom && <span style={{ color: '#f39c12' }}> ⭐</span>}
            </h3>
            <div style={{ color: '#888', marginBottom: '0.5rem' }}>
              Room ID: {selectedRoom.id}
            </div>
            <div style={{ color: '#888', marginBottom: '1rem', fontSize: '0.9rem' }}>
              {selectedRoom.description}
            </div>

            <div style={{ marginBottom: '1rem' }}>
              <strong style={{ color: '#6c5ce7' }}>Exits:</strong>
              {Object.entries(selectedRoom.exits || {}).length === 0 ? (
                <div style={{ color: '#666', marginTop: '0.25rem' }}>None</div>
              ) : (
                <div style={{ marginTop: '0.5rem' }}>
                  {Object.entries(selectedRoom.exits || {}).map(([dir, targetId]) => {
                    const targetRoom = rooms.find(r => r.id === targetId)
                    return (
                      <div
                        key={dir}
                        onClick={() => targetRoom && setSelectedRoom(targetRoom)}
                        style={{
                          padding: '0.25rem 0.5rem',
                          margin: '0.25rem 0',
                          background: '#16213e',
                          borderRadius: '4px',
                          cursor: 'pointer',
                          color: '#fff',
                          fontSize: '0.85rem'
                        }}
                      >
                        <strong>{dir}</strong> → {targetRoom?.name || `Room ${targetId}`}
                      </div>
                    )
                  })}
                </div>
              )}
            </div>

            <button
              onClick={() => setSelectedRoom(null)}
              style={{
                width: '100%',
                padding: '0.5rem',
                background: '#16213e',
                border: '1px solid #333',
                borderRadius: '4px',
                color: '#888',
                cursor: 'pointer'
              }}
            >
              Close
            </button>
          </>
        ) : (
          <div style={{ color: '#666', textAlign: 'center', marginTop: '2rem' }}>
            <p>Click a room to see details</p>
            <p style={{ fontSize: '0.8rem', marginTop: '1rem' }}>
              Green rooms are starting points
            </p>
          </div>
        )}

        {/* Room List */}
        <div style={{ marginTop: '2rem', borderTop: '1px solid #333', paddingTop: '1rem' }}>
          <h4 style={{ margin: '0 0 0.5rem 0', color: '#888' }}>All Rooms</h4>
          <div style={{ maxHeight: '200px', overflowY: 'auto' }}>
            {rooms.slice(0, 20).map(room => (
              <div
                key={room.id}
                onClick={() => setSelectedRoom(room)}
                style={{
                  padding: '0.25rem 0.5rem',
                  cursor: 'pointer',
                  color: selectedRoom?.id === room.id ? '#6c5ce7' : '#fff',
                  background: selectedRoom?.id === room.id ? '#16213e' : 'transparent',
                  borderRadius: '4px',
                  fontSize: '0.85rem'
                }}
              >
                {room.name}
                {room.isStartingRoom && ' ⭐'}
              </div>
            ))}
            {rooms.length > 20 && (
              <div style={{ color: '#666', fontSize: '0.8rem', padding: '0.25rem' }}>
                ...and {rooms.length - 20} more
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}