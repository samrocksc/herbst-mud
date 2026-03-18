import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useMemo } from 'react'
import { MapFlow } from '../../components/MapFlow'
import { ZLevelSelector } from '../../components/ZLevelSelector'
import { DirectionPickerModal } from '../../components/DirectionPickerModal'
import type { Node, Edge, Connection } from '@xyflow/react'
import type { ExitEdgeData } from '../../components/ExitEdge'

export const Route = createFileRoute('/admin/map')({
  component: MapBuilder,
})

interface MapRoomData extends Record<string, unknown> {
  name: string
  description: string
  zLevel: number
}

// Sample rooms with different Z-levels for testing
const initialNodes: Node[] = [
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
    // Z-level 1 rooms (upper floor)
    {
      id: '6',
      type: 'room',
      position: { x: 250, y: 150 },
      data: { name: 'Town Square Upstairs', description: 'Upper level of town square', zLevel: 1 },
      selected: false
    },
    {
      id: '7',
      type: 'room',
      position: { x: 400, y: 200 },
      data: { name: 'Inn Upper Floor', description: 'Guest rooms upstairs', zLevel: 1 },
      selected: false
    },
    // Z-level -1 rooms (underground)
    {
      id: '8',
      type: 'room',
      position: { x: 250, y: 300 },
      data: { name: 'Town Square Cellar', description: 'Storage basement', zLevel: -1 },
      selected: false
    },
    {
      id: '9',
      type: 'room',
      position: { x: 100, y: 400 },
      data: { name: 'Sewers', description: 'Dark underground tunnels', zLevel: -1 },
      selected: false
    },
  ]

// Exit edges between rooms using ExitEdge component
const initialEdges: Edge<ExitEdgeData>[] = [
    { id: 'e1-2', source: '1', target: '2', type: 'exit', data: { direction: 'south' } },
    { id: 'e2-1', source: '2', target: '1', type: 'exit', data: { direction: 'north' } },
    { id: 'e1-3', source: '1', target: '3', type: 'exit', data: { direction: 'north' } },
    { id: 'e3-1', source: '3', target: '1', type: 'exit', data: { direction: 'south' } },
    { id: 'e1-4', source: '1', target: '4', type: 'exit', data: { direction: 'east' } },
    { id: 'e4-1', source: '4', target: '1', type: 'exit', data: { direction: 'west' } },
    { id: 'e1-5', source: '1', target: '5', type: 'exit', data: { direction: 'west' } },
    { id: 'e5-1', source: '5', target: '1', type: 'exit', data: { direction: 'east' } },
    // Z-exits (up/down connections between levels)
    { id: 'e1-6', source: '1', target: '6', type: 'exit', data: { direction: 'up', isZExit: true } },
    { id: 'e6-1', source: '6', target: '1', type: 'exit', data: { direction: 'down', isZExit: true } },
    { id: 'e1-8', source: '1', target: '8', type: 'exit', data: { direction: 'down', isZExit: true } },
    { id: 'e8-1', source: '8', target: '1', type: 'exit', data: { direction: 'up', isZExit: true } },
    { id: 'e8-9', source: '8', target: '9', type: 'exit', data: { direction: 'south' } },
    { id: 'e9-8', source: '9', target: '8', type: 'exit', data: { direction: 'north' } },
    { id: 'e6-7', source: '6', target: '7', type: 'exit', data: { direction: 'east' } },
    { id: 'e7-6', source: '7', target: '6', type: 'exit', data: { direction: 'west' } },
  ]

function MapBuilder() {
  const [nodes, setNodes] = useState<Node[]>(initialNodes)
  const [edges, setEdges] = useState<Edge<ExitEdgeData>[]>(initialEdges)
  const [currentZLevel, setCurrentZLevel] = useState(0)
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  
  // Direction picker modal state
  const [pendingConnection, setPendingConnection] = useState<{
    source: string
    target: string
    sourceName: string
    targetName: string
  } | null>(null)

  // Filter nodes by Z-level (show current level prominently, adjacent faintly)
  const filteredNodes = useMemo(() => {
    return nodes.map(node => {
      const zLevel = (node.data as MapRoomData).zLevel ?? 0
      const isCurrentLevel = zLevel === currentZLevel
      const isAdjacent = Math.abs(zLevel - currentZLevel) === 1
      
      // Return node with opacity based on level
      return {
        ...node,
        selected: isCurrentLevel && node.selected,
        // For nodes on other Z-levels, we don't show them unless adjacent
        hidden: !isCurrentLevel && !isAdjacent,
        style: !isCurrentLevel && isAdjacent ? { opacity: 0.4 } : undefined
      }
    })
  }, [nodes, currentZLevel])

  // Filter edges - hide edges between different Z-levels unless they're Z-exits
  const filteredEdges = useMemo(() => {
    return edges.filter(edge => {
      const sourceNode = nodes.find(n => n.id === edge.source)
      const targetNode = nodes.find(n => n.id === edge.target)
      if (!sourceNode || !targetNode) return false
      
      const sourceZ = (sourceNode.data as MapRoomData).zLevel ?? 0
      const targetZ = (targetNode.data as MapRoomData).zLevel ?? 0
      
      // Show if same Z-level
      if (sourceZ === targetZ) return true
      
      // Show Z-exit connections (different Z-levels with up/down)
      const label = (edge.label as string || '').toLowerCase()
      return label === 'up' || label === 'down'
    })
  }, [edges, nodes, currentZLevel])

  // Handle connection - show direction picker modal
  const onConnect = useCallback((connection: Connection) => {
    if (connection.source && connection.target) {
      const sourceNode = nodes.find(n => n.id === connection.source)
      const targetNode = nodes.find(n => n.id === connection.target)
      
      if (sourceNode && targetNode) {
        // Show direction picker
        setPendingConnection({
          source: connection.source,
          target: connection.target,
          sourceName: (sourceNode.data as MapRoomData).name,
          targetName: (targetNode.data as MapRoomData).name,
        })
      }
    }
  }, [nodes])

  // Handle direction selection - create bidirectional exits
  const handleDirectionSelect = useCallback((sourceDirection: string, targetDirection: string) => {
    if (!pendingConnection) return
    
    const isZExit = sourceDirection === 'up' || sourceDirection === 'down'
    
    const newEdges: Edge<ExitEdgeData>[] = [
      {
        id: `e${pendingConnection.source}-${pendingConnection.target}-${sourceDirection}`,
        source: pendingConnection.source,
        target: pendingConnection.target,
        type: 'exit',
        data: { 
          direction: sourceDirection as ExitEdgeData['direction'],
          isZExit: isZExit,
        },
      },
      {
        id: `e${pendingConnection.target}-${pendingConnection.source}-${targetDirection}`,
        source: pendingConnection.target,
        target: pendingConnection.source,
        type: 'exit',
        data: { 
          direction: targetDirection as ExitEdgeData['direction'],
          isZExit: targetDirection === 'up' || targetDirection === 'down',
        },
      },
    ]
    
    setEdges(eds => [...eds, ...newEdges])
    setPendingConnection(null)
  }, [pendingConnection])

  const handleConnectionCancel = useCallback(() => {
    setPendingConnection(null)
  }, [])

  const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
    setSelectedNode(node)
    setNodes(nds => nds.map(n => ({
      ...n,
      selected: n.id === node.id
    })))
  }, [])

  const addNewRoom = () => {
    const newId = String(nodes.length + 1)
    const newNode: Node = {
      id: newId,
      type: 'room',
      position: { x: Math.random() * 400 + 100, y: Math.random() * 400 + 100 },
      data: { name: `Room ${newId}`, description: 'New room', zLevel: currentZLevel }
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

      {/* Z-Level Selector */}
      <ZLevelSelector 
        currentLevel={currentZLevel} 
        onChange={setCurrentZLevel}
      />

      <div className="map-container" style={{ display: 'flex', gap: '16px' }}>
        <div className="map-flow" style={{ flex: 1 }}>
          <MapFlow
            nodes={filteredNodes}
            edges={filteredEdges}
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
        <span style={{ marginLeft: '16px' }}>🟠 Z: Up</span>
        <span style={{ marginLeft: '16px' }}>🔵 Z: Down</span>
      </div>

      {/* Direction Picker Modal */}
      {pendingConnection && (
        <DirectionPickerModal
          sourceName={pendingConnection.sourceName}
          targetName={pendingConnection.targetName}
          onSelect={handleDirectionSelect}
          onCancel={handleConnectionCancel}
        />
      )}
    </div>
  )
}