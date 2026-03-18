import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useMemo } from 'react'
import { MapFlow } from '../../components/MapFlow'
import { ZLevelSelector } from '../../components/ZLevelSelector'
import { RoomEditPanel } from '../../components/RoomEditPanel'
import type { Node, Edge, Connection } from '@xyflow/react'

export const Route = createFileRoute('/admin/map')({
  component: MapBuilder,
})

interface MapRoomData extends Record<string, unknown> {
  name: string
  description: string
  zLevel: number
  exits?: Record<string, number>
}

// Sample rooms with different Z-levels for testing
const initialNodes: Node[] = [
    { 
      id: '1', 
      type: 'room',
      position: { x: 250, y: 100 }, 
      data: { name: 'Town Square', description: 'The central hub', zLevel: 0, exits: { north: 2, east: 4, south: 3, west: 5 } },
      selected: false 
    },
    { 
      id: '2', 
      type: 'room',
      position: { x: 250, y: 250 }, 
      data: { name: 'Main Street North', description: 'Street heading north', zLevel: 0, exits: { south: 1 } },
      selected: false 
    },
    { 
      id: '3', 
      type: 'room',
      position: { x: 250, y: 400 }, 
      data: { name: 'Main Street South', description: 'Street heading south', zLevel: 0, exits: { north: 1 } },
      selected: false 
    },
    { 
      id: '4', 
      type: 'room',
      position: { x: 450, y: 175 }, 
      data: { name: 'Forest Path', description: 'A path through the woods', zLevel: 0, exits: { west: 1 } },
      selected: false 
    },
    { 
      id: '5', 
      type: 'room',
      position: { x: 50, y: 175 }, 
      data: { name: 'Shop District', description: 'Where merchants sell goods', zLevel: 0, exits: { east: 1 } },
      selected: false 
    },
    // Z-level 1 rooms (upper floor)
    {
      id: '6',
      type: 'room',
      position: { x: 250, y: 150 },
      data: { name: 'Town Square Upstairs', description: 'Upper level of town square', zLevel: 1, exits: { down: 1, east: 7 } },
      selected: false
    },
    {
      id: '7',
      type: 'room',
      position: { x: 400, y: 200 },
      data: { name: 'Inn Upper Floor', description: 'Guest rooms upstairs', zLevel: 1, exits: { west: 6 } },
      selected: false
    },
    // Z-level -1 rooms (underground)
    {
      id: '8',
      type: 'room',
      position: { x: 250, y: 300 },
      data: { name: 'Town Square Cellar', description: 'Storage basement', zLevel: -1, exits: { up: 1, south: 9 } },
      selected: false
    },
    {
      id: '9',
      type: 'room',
      position: { x: 100, y: 400 },
      data: { name: 'Sewers', description: 'Dark underground tunnels', zLevel: -1, exits: { north: 8 } },
      selected: false
    },
  ]

// Z-exits between levels
const initialEdges: Edge[] = [
    { id: 'e1-2', source: '1', target: '2', label: 'north', type: 'smoothstep' },
    { id: 'e1-3', source: '1', target: '3', label: 'south', type: 'smoothstep' },
    { id: 'e1-4', source: '1', target: '4', label: 'east', type: 'smoothstep' },
    { id: 'e1-5', source: '1', target: '5', label: 'west', type: 'smoothstep' },
    // Z-exits (up/down connections between levels)
    { id: 'e1-6', source: '1', target: '6', label: 'up', data: { isZExit: true, direction: 'up' }, type: 'smoothstep', animated: true, style: { stroke: '#e17055', strokeWidth: 2 } },
    { id: 'e1-8', source: '1', target: '8', label: 'down', data: { isZExit: true, direction: 'down' }, type: 'smoothstep', animated: true, style: { stroke: '#74b9ff', strokeWidth: 2 } },
    { id: 'e8-9', source: '8', target: '9', label: 'south', type: 'smoothstep' },
    { id: 'e6-7', source: '6', target: '7', label: 'east', type: 'smoothstep' },
  ]

function MapBuilder() {
  const [nodes, setNodes] = useState<Node[]>(initialNodes)
  const [edges, setEdges] = useState<Edge[]>(initialEdges)
  const [currentZLevel, setCurrentZLevel] = useState(0)
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)

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
      data: { name: `Room ${newId}`, description: 'New room', zLevel: currentZLevel, exits: {} }
    }
    setNodes(nds => [...nds, newNode])
  }

  // Convert nodes to rooms list for RoomEditPanel
  const roomsList = useMemo(() => {
    return nodes.map(node => ({
      id: parseInt(node.id, 10),
      name: (node.data as MapRoomData).name
    }))
  }, [nodes])

  // Handle room update from RoomEditPanel
  const handleRoomUpdate = useCallback((id: string, data: { name: string; description: string; zLevel: number; exits: Record<string, number> }) => {
    setNodes(nds => nds.map(n => 
      n.id === id 
        ? { ...n, data: { ...n.data, ...data } }
        : n
    ))
    
    // Update selected node to reflect changes
    if (selectedNode?.id === id) {
      setSelectedNode({ ...selectedNode, data: { ...selectedNode.data, ...data } })
    }
    
    // Update edges based on exits
    setEdges(eds => {
      // Remove old edges for this room
      const filtered = eds.filter(e => e.source !== id)
      
      // Add new edges from exits
      const newEdges = Object.entries(data.exits).map(([direction, targetId]) => ({
        id: `e${id}-${targetId}`,
        source: id,
        target: String(targetId),
        label: direction,
        type: 'smoothstep' as const
      }))
      
      return [...filtered, ...newEdges]
    })
  }, [selectedNode])

  // Handle room delete
  const handleRoomDelete = useCallback((id: string) => {
    if (confirm('Are you sure you want to delete this room?')) {
      setNodes(nds => nds.filter(n => n.id !== id))
      setEdges(eds => eds.filter(e => e.source !== id && e.target !== id))
      setSelectedNode(null)
    }
  }, [])

  // Handle panel close
  const handlePanelClose = useCallback(() => {
    setSelectedNode(null)
    setNodes(nds => nds.map(n => ({ ...n, selected: false })))
  }, [])

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

        <RoomEditPanel
          selectedNode={selectedNode}
          rooms={roomsList}
          onUpdate={handleRoomUpdate}
          onDelete={handleRoomDelete}
          onClose={handlePanelClose}
        />
      </div>

      <div className="map-legend" style={{ marginTop: '12px', padding: '8px', background: '#222', borderRadius: '4px' }}>
        <span><strong>Legend:</strong></span>
        <span style={{ marginLeft: '16px' }}>🟩 Room Node</span>
        <span style={{ marginLeft: '16px' }}>➡️ Exit Connection</span>
        <span style={{ marginLeft: '16px' }}>🟦 Selected</span>
        <span style={{ marginLeft: '16px' }}>🟠 Z: Up</span>
        <span style={{ marginLeft: '16px' }}>🔵 Z: Down</span>
      </div>
    </div>
  )
}