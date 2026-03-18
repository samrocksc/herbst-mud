import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useMemo, useEffect } from 'react'
import { MapFlow } from '../../components/MapFlow'
import { ZLevelSelector } from '../../components/ZLevelSelector'
import { SidebarPalette } from '../../components/SidebarPalette'
import { useRooms, useCreateRoom, useUpdateRoom, useDeleteRoom, type Room } from '../../hooks/useRooms'
import type { Node, Edge, Connection } from '@xyflow/react'

export const Route = createFileRoute('/admin/map')({
  component: MapBuilder,
})

interface MapRoomData extends Record<string, unknown> {
  id: number
  name: string
  description: string
  zLevel: number
  exits?: Record<string, number>
}

// Helper to transform API Room to ReactFlow Node
function roomToNode(room: Room): Node {
  return {
    id: String(room.id),
    type: 'room',
    position: { x: room.x ?? 100, y: room.y ?? 100 },
    data: {
      id: room.id,
      name: room.name,
      description: room.description,
      zLevel: room.zLevel ?? 0,
      exits: room.exits ?? {}
    } as MapRoomData,
    selected: false
  }
}

// Helper to transform API Room[] to Edge[] based on exits
function roomsToEdges(rooms: Room[]): Edge[] {
  const edges: Edge[] = []
  const directionLabels: Record<string, string> = {
    north: 'north',
    south: 'south',
    east: 'east',
    west: 'west',
    up: 'up',
    down: 'down'
  }

  for (const room of rooms) {
    if (!room.exits) continue
    for (const [direction, targetId] of Object.entries(room.exits)) {
      const edgeId = `e${room.id}-${targetId}-${direction}`
      const isZExit = direction === 'up' || direction === 'down'
      edges.push({
        id: edgeId,
        source: String(room.id),
        target: String(targetId),
        label: directionLabels[direction] ?? direction,
        type: 'smoothstep',
        animated: isZExit,
        style: isZExit
          ? { stroke: direction === 'up' ? '#e17055' : '#74b9ff', strokeWidth: 2 }
          : undefined
      })
    }
  }
  return edges
}

function MapBuilder() {
  const [nodes, setNodes] = useState<Node[]>(initialNodes)
  const [currentZLevel, setCurrentZLevel] = useState(0)
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  const [isCreatingRoom, setIsCreatingRoom] = useState(false)
  
  const createRoomMutation = useCreateRoom()

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

  const [edges, setEdges] = useState<Edge[]>(initialEdges)

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
      data: { name: `Room ${newId}`, description: 'New room', zLevel: currentZLevel }
    }
    setNodes(nds => [...nds, newNode])
  }

  // Handle room drop from palette - creates room via API
  const handleRoomDrop = useCallback(async (position: { x: number; y: number }) => {
    setIsCreatingRoom(true)
    try {
      // Create room via API
      const newRoom = await createRoomMutation.mutateAsync({
        name: 'New Room',
        description: 'A newly created room',
        zLevel: currentZLevel,
        x: Math.round(position.x),
        y: Math.round(position.y)
      })
      
      // Add to local state as ReactFlow node
      const newNode: Node = {
        id: String(newRoom.id),
        type: 'room',
        position: { x: newRoom.x || position.x, y: newRoom.y || position.y },
        data: { 
          name: newRoom.name, 
          description: newRoom.description, 
          zLevel: newRoom.zLevel ?? currentZLevel,
          exits: newRoom.exits || {}
        },
        selected: true
      }
      
      setNodes(nds => [...nds, newNode])
      setSelectedNode(newNode)
    } catch (error) {
      console.error('Failed to create room:', error)
      // Fallback: create locally without API
      const newId = String(nodes.length + 1)
      const newNode: Node = {
        id: newId,
        type: 'room',
        position,
        data: { name: 'New Room', description: 'New room (offline)', zLevel: currentZLevel }
      }
      setNodes(nds => [...nds, newNode])
    } finally {
      setIsCreatingRoom(false)
    }
  }, [createRoomMutation, currentZLevel, nodes.length])

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
          <button disabled={isCreatingRoom}>
            {isCreatingRoom ? 'Creating...' : 'Connect Rooms'}
          </button>
          <button>Save Map</button>
        </div>
      </div>

      {/* Z-Level Selector */}
      <ZLevelSelector 
        currentLevel={currentZLevel} 
        onChange={setCurrentZLevel}
      />

      <div className="map-container" style={{ display: 'flex', gap: '16px' }}>
        {/* Sidebar Palette for drag-and-drop */}
        <SidebarPalette />
        
        <div className="map-flow" style={{ flex: 1 }}>
          <MapFlow
            nodes={filteredNodes}
            edges={filteredEdges}
            onConnect={onConnect}
            onNodeClick={onNodeClick}
            onDrop={handleRoomDrop}
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
    </div>
  )
}