import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback } from 'react'
import { MapFlow } from '../../components/MapFlow'
import type { Node, Edge, Connection } from '@xyflow/react'

export const Route = createFileRoute('/admin/map')({
  component: MapBuilder,
})

interface MapRoomData extends Record<string, unknown> {
  name: string
  description: string
  zLevel: number
}

function MapBuilder() {
  const [nodes, setNodes] = useState<Node[]>([
    { 
      id: '1', 
      type: 'room',
      position: { x: 250, y: 100 }, 
      data: { name: 'Town Square', description: 'The central hub', zLevel: 0 },
      selected: false 
    },
    { 
      id: '2', 
      type: 'room',
      position: { x: 250, y: 250 }, 
      data: { name: 'Main Street North', description: 'Street heading north', zLevel: 0 },
      selected: false 
    },
    { 
      id: '3', 
      type: 'room',
      position: { x: 250, y: 400 }, 
      data: { name: 'Main Street South', description: 'Street heading south', zLevel: 0 },
      selected: false 
    },
    { 
      id: '4', 
      type: 'room',
      position: { x: 450, y: 175 }, 
      data: { name: 'Forest Path', description: 'A path through the woods', zLevel: 0 },
      selected: false 
    },
    { 
      id: '5', 
      type: 'room',
      position: { x: 50, y: 175 }, 
      data: { name: 'Shop District', description: 'Where merchants sell goods', zLevel: 0 },
      selected: false 
    },
  ])

  const [edges, setEdges] = useState<Edge[]>([
    { id: 'e1-2', source: '1', target: '2', label: 'north', type: 'smoothstep' },
    { id: 'e1-3', source: '1', target: '3', label: 'south', type: 'smoothstep' },
    { id: 'e1-4', source: '1', target: '4', label: 'east', type: 'smoothstep' },
    { id: 'e1-5', source: '1', target: '5', label: 'west', type: 'smoothstep' },
  ])

  const [selectedNode, setSelectedNode] = useState<Node | null>(null)

  const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
    setSelectedNode(node)
    setNodes(nds => nds.map(n => ({
      ...n,
      selected: n.id === node.id
    })))
  }, [])

  const onConnect = useCallback((connection: Connection) => {
    if (connection.source && connection.target) {
      const direction = connection.sourceHandle || 'connected'
      setEdges(eds => [...eds, {
        id: `e${connection.source}-${connection.target}`,
        source: connection.source!,
        target: connection.target!,
        label: direction,
        type: 'smoothstep',
        animated: true
      }])
    }
  }, [])

  const addNewRoom = () => {
    const newId = String(nodes.length + 1)
    const newNode: Node = {
      id: newId,
      type: 'room',
      position: { x: Math.random() * 400 + 100, y: Math.random() * 400 + 100 },
      data: { name: `Room ${newId}`, description: 'New room', zLevel: 0 }
    }
    setNodes(nds => [...nds, newNode])
  }

  const getRoomData = (node: Node | null): MapRoomData => {
    if (!node) return { name: '', description: '', zLevel: 0 }
    return node.data as MapRoomData
  }

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Map Builder</h2>
        <div className="map-actions">
          <button onClick={addNewRoom}>Add Room</button>
          <button>Connect Rooms</button>
          <button>Save Map</button>
        </div>
      </div>

      <div className="map-container" style={{ display: 'flex', gap: '16px' }}>
        <div className="map-flow" style={{ flex: 1 }}>
          <MapFlow
            nodes={nodes}
            edges={edges}
            onConnect={onConnect}
            onNodeClick={onNodeClick}
          />
        </div>

        {selectedNode && (
          <div className="map-sidebar" style={{ 
            width: '280px', 
            padding: '16px', 
            background: '#222', 
            borderRadius: '8px',
            border: '1px solid #444'
          }}>
            <h3>Room Details</h3>
            <div style={{ marginBottom: '12px' }}>
              <label style={{ display: 'block', marginBottom: '4px' }}>Name:</label>
              <input 
                type="text" 
                value={getRoomData(selectedNode).name}
                onChange={(e) => {
                  const newData = { ...getRoomData(selectedNode), name: e.target.value }
                  setNodes(nds => nds.map(n => 
                    n.id === selectedNode.id 
                      ? { ...n, data: newData }
                      : n
                  ))
                  setSelectedNode({ ...selectedNode, data: newData })
                }}
                style={{ width: '100%', padding: '6px', background: '#333', border: '1px solid #555', color: '#fff' }}
              />
            </div>
            <div style={{ marginBottom: '12px' }}>
              <label style={{ display: 'block', marginBottom: '4px' }}>Description:</label>
              <textarea 
                value={getRoomData(selectedNode).description}
                onChange={(e) => {
                  const newData = { ...getRoomData(selectedNode), description: e.target.value }
                  setNodes(nds => nds.map(n => 
                    n.id === selectedNode.id 
                      ? { ...n, data: newData }
                      : n
                  ))
                  setSelectedNode({ ...selectedNode, data: newData })
                }}
                style={{ width: '100%', padding: '6px', background: '#333', border: '1px solid #555', color: '#fff', minHeight: '60px' }}
              />
            </div>
            <div style={{ marginBottom: '12px' }}>
              <label style={{ display: 'block', marginBottom: '4px' }}>Z-Level:</label>
              <select 
                value={getRoomData(selectedNode).zLevel}
                onChange={(e) => {
                  const zLevel = parseInt(e.target.value)
                  const newData = { ...getRoomData(selectedNode), zLevel }
                  setNodes(nds => nds.map(n => 
                    n.id === selectedNode.id 
                      ? { ...n, data: newData }
                      : n
                  ))
                  setSelectedNode({ ...selectedNode, data: newData })
                }}
                style={{ width: '100%', padding: '6px', background: '#333', border: '1px solid #555', color: '#fff' }}
              >
                <option value={-2}>Z: -2 (Deep Underground)</option>
                <option value={-1}>Z: -1 (Underground)</option>
                <option value={0}>Z: 0 (Ground)</option>
                <option value={1}>Z: 1 (Upper Floor)</option>
                <option value={2}>Z: 2 (Tower)</option>
              </select>
            </div>
            <p style={{ fontSize: '12px', color: '#888' }}>Node ID: {selectedNode.id}</p>
          </div>
        )}
      </div>

      <div className="map-legend" style={{ marginTop: '12px', padding: '8px', background: '#222', borderRadius: '4px' }}>
        <span><strong>Legend:</strong></span>
        <span style={{ marginLeft: '16px' }}>🟩 Room Node</span>
        <span style={{ marginLeft: '16px' }}>➡️ Exit Connection</span>
        <span style={{ marginLeft: '16px' }}>🟦 Selected</span>
      </div>
    </div>
  )
}