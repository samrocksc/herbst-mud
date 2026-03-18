import { createFileRoute } from '@tanstack/react-router'
import { useState, useEffect } from 'react'
import { Node, Edge } from '@xyflow/react'

import MapFlow from '../../components/MapFlow'
import { RoomNode, RoomNodeType } from '../../components/RoomNode'

export const Route = createFileRoute('/_auth/map')({
  component: MapBuilder,
})

// Demo rooms for initial state
const initialDemoNodes: Node[] = [
  {
    id: '1',
    type: 'roomNode',
    position: { x: 100, y: 100 },
    data: { label: 'Town Square', roomId: '1', exits: ['n', 'e', 's', 'w'] },
  },
  {
    id: '2',
    type: 'roomNode',
    position: { x: 100, y: 250 },
    data: { label: 'Main Street', roomId: '2', exits: ['n', 's'] },
  },
  {
    id: '3',
    type: 'roomNode',
    position: { x: 250, y: 100 },
    data: { label: 'Shop District', roomId: '3', exits: ['w'] },
  },
  {
    id: '4',
    type: 'roomNode',
    position: { x: 400, y: 100 },
    data: { label: 'Forest Path', roomId: '4', exits: ['e'] },
  },
  {
    id: '5',
    type: 'roomNode',
    position: { x: 250, y: 300 },
    data: { label: 'Tavern', roomId: '5', exits: ['n'] },
  },
]

const initialDemoEdges: Edge[] = [
  {
    id: 'e1-2',
    source: '1',
    target: '2',
    type: 'exit',
    data: { direction: 'south' },
  },
  {
    id: 'e1-3',
    source: '1',
    target: '3',
    type: 'exit',
    data: { direction: 'east' },
  },
]

interface RoomData {
  id: string
  name: string
  description?: string
  x?: number
  y?: number
  z?: number
  exits?: Array<{
    direction: string
    roomId: string
  }>
}

function MapBuilder() {
  const [nodes, setNodes] = useState<Node[]>(initialDemoNodes)
  const [edges, setEdges] = useState<Edge[]>(initialDemoEdges)
  const [isLoading, setIsLoading] = useState(false)
  const [selectedRoomId, setSelectedRoomId] = useState<string | null>(null)

  // TODO: Replace with actual API call when backend is ready
  // const { data: rooms } = useQuery({ queryKey: ['rooms'], queryFn: fetchRooms })

  const handleAddRoom = () => {
    const newId = String(Date.now())
    const newNode: Node = {
      id: newId,
      type: 'roomNode',
      position: { x: 200 + Math.random() * 100, y: 200 + Math.random() * 100 },
      data: {
        label: 'New Room',
        roomId: newId,
        exits: [],
      },
    }
    setNodes((nds) => [...nds, newNode])
  }

  const handleSaveMap = async () => {
    setIsLoading(true)
    try {
      // TODO: Implement actual save to API
      console.log('Saving map:', { nodes, edges })
      
      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 500))
      alert('Map saved successfully!')
    } catch (error) {
      console.error('Failed to save map:', error)
      alert('Failed to save map')
    } finally {
      setIsLoading(false)
    }
  }

  const selectedRoom = nodes.find((n) => n.id === selectedRoomId)

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>🗺️ Map Builder</h2>
        <div className="map-actions">
          <button onClick={handleAddRoom} className="btn-primary">
            + Add Room
          </button>
          <button onClick={handleSaveMap} disabled={isLoading} className="btn-success">
            {isLoading ? 'Saving...' : '💾 Save Map'}
          </button>
        </div>
      </div>

      <div className="map-layout">
        <div className="map-main">
          <MapFlow initialNodes={nodes} initialEdges={edges} />
        </div>

        {selectedRoom && (
          <div className="map-sidebar">
            <h3>Room Details</h3>
            <div className="room-form">
              <label>
                Name:
                <input
                  type="text"
                  value={selectedRoom.data.label}
                  onChange={(e) => {
                    setNodes((nds) =>
                      nds.map((n) =>
                        n.id === selectedRoom.id
                          ? { ...n, data: { ...n.data, label: e.target.value } }
                          : n
                      )
                    )
                  }}
                />
              </label>
              <div className="room-exits">
                <strong>Exits:</strong>
                <div className="exit-tags">
                  {edges
                    .filter((e) => e.source === selectedRoom.id)
                    .map((e) => (
                      <span key={e.id} className="exit-tag">
                        {(e.data as { direction?: string })?.direction || 'unknown'}
                      </span>
                    ))}
                </div>
              </div>
              <div className="room-actions">
                <button className="btn-small">Edit Details</button>
                <button className="btn-small">Add Items</button>
                <button className="btn-small">Add NPCs</button>
              </div>
              <button
                className="btn-danger"
                onClick={() => setSelectedRoomId(null)}
              >
                Close
              </button>
            </div>
          </div>
        )}
      </div>

      <div className="map-instructions">
        <h4>📝 Instructions</h4>
        <ul>
          <li><strong>Drag rooms</strong> to reposition them on the canvas</li>
          <li><strong>Drag from a handle</strong> (colored dots) to another room to create an exit</li>
          <li><strong>Click a room</strong> to edit its details</li>
          <li><strong>Scroll/zoom</strong> to navigate the map</li>
        </ul>
      </div>
    </div>
  )
}

export default MapBuilder