import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useEffect, useRef } from 'react'
import { MapFlow } from '../../components/MapFlow'
import { SidebarPalette } from '../../components/SidebarPalette'
import type { Node, Edge, Connection } from '@xyflow/react'
import { useRooms, type RoomNodeData } from '../../hooks/useRooms'
import { useReactFlow, ReactFlowProvider } from '@xyflow/react'

export const Route = createFileRoute('/admin/map')({
  component: MapBuilder,
})

function MapBuilderContent() {
  const { 
    nodes: apiNodes, 
    edges: apiEdges, 
    loading, 
    error, 
    updateNodePosition,
    updateRoom,
    refreshRooms 
  } = useRooms()
  
  const { screenToFlowPosition } = useReactFlow()
  
  const [nodes, setNodes] = useState<Node<RoomNodeData>[]>([])
  const [edges, setEdges] = useState<Edge[]>([])
  const [selectedNode, setSelectedNode] = useState<Node<RoomNodeData> | null>(null)
  const reactFlowWrapper = useRef<HTMLDivElement>(null)
  
  // Handle drag over - allow drop
  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    e.dataTransfer.dropEffect = 'copy'
  }, [])
  
  // Handle drop - create new room
  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault()
    
    const nodeType = e.dataTransfer.getData('application/x-herbstmud-node-type')
    if (nodeType !== 'room') return
    
    // Calculate position in flow coordinates
    const reactFlowBounds = reactFlowWrapper.current?.getBoundingClientRect()
    if (!reactFlowBounds) return
    
    const position = screenToFlowPosition({
      x: e.clientX - reactFlowBounds.left,
      y: e.clientY - reactFlowBounds.top,
    })
    
    // Create new room node
    const newId = String(nodes.length + 1)
    const newNode: Node<RoomNodeData> = {
      id: newId,
      type: 'room',
      position,
      data: { 
        name: `New Room`, 
        description: 'Click to edit', 
        zLevel: 0, 
        roomId: parseInt(newId),
        isStartingRoom: false,
        atmosphere: 'air',
        exits: {}
      }
    }
    
    // Add to nodes and select it
    setNodes(nds => [...nds, newNode])
    setSelectedNode(newNode)
    
    // TODO: Also create via API: createRoom({ name, description, position, zLevel })
  }, [nodes.length, screenToFlowPosition])

  // Sync API nodes to local state when loaded
  useEffect(() => {
    if (!loading && apiNodes.length > 0) {
      setNodes(apiNodes)
      setEdges(apiEdges)
    } else if (!loading && apiNodes.length === 0 && !error) {
      // No rooms in API - use fallback mock data for development
      setNodes([
        { 
          id: '1', 
          type: 'room',
          position: { x: 250, y: 100 }, 
          data: { name: 'Town Square', description: 'The central hub', zLevel: 0, roomId: 1, isStartingRoom: false, atmosphere: 'air', exits: {} },
          selected: false 
        },
        { 
          id: '2', 
          type: 'room',
          position: { x: 250, y: 250 }, 
          data: { name: 'Main Street North', description: 'Street heading north', zLevel: 0, roomId: 2, isStartingRoom: false, atmosphere: 'air', exits: {} },
          selected: false 
        },
        { 
          id: '3', 
          type: 'room',
          position: { x: 250, y: 400 }, 
          data: { name: 'Main Street South', description: 'Street heading south', zLevel: 0, roomId: 3, isStartingRoom: false, atmosphere: 'air', exits: {} },
          selected: false 
        },
        { 
          id: '4', 
          type: 'room',
          position: { x: 450, y: 175 }, 
          data: { name: 'Forest Path', description: 'A path through the woods', zLevel: 0, roomId: 4, isStartingRoom: false, atmosphere: 'air', exits: {} },
          selected: false 
        },
        { 
          id: '5', 
          type: 'room',
          position: { x: 50, y: 175 }, 
          data: { name: 'Shop District', description: 'Where merchants sell goods', zLevel: 0, roomId: 5, isStartingRoom: false, atmosphere: 'air', exits: {} },
          selected: false 
        },
      ])
      setEdges([
        { id: 'e1-2', source: '1', target: '2', label: 'north', type: 'smoothstep' },
        { id: 'e1-3', source: '1', target: '3', label: 'south', type: 'smoothstep' },
        { id: 'e1-4', source: '1', target: '4', label: 'east', type: 'smoothstep' },
        { id: 'e1-5', source: '1', target: '5', label: 'west', type: 'smoothstep' },
      ])
    }
  }, [loading, apiNodes, apiEdges, error])

  const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
    setSelectedNode(node as Node<RoomNodeData>)
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

  const addNewRoom = async () => {
    const newId = nodes.length + 1
    const newNode: Node<RoomNodeData> = {
      id: String(newId),
      type: 'room',
      position: { x: Math.random() * 400 + 100, y: Math.random() * 400 + 100 },
      data: { 
        name: `Room ${newId}`, 
        description: 'New room', 
        zLevel: 0, 
        roomId: newId,
        isStartingRoom: false,
        atmosphere: 'air',
        exits: {}
      }
    }
    setNodes(nds => [...nds, newNode])
  }

  const getRoomData = (node: Node<RoomNodeData> | null): RoomNodeData => {
    if (!node || !node.data) return { name: '', description: '', zLevel: 0, roomId: 0, isStartingRoom: false, atmosphere: 'air', exits: {} }
    return node.data
  }

  const handleSaveRoom = async () => {
    if (!selectedNode) return
    
    try {
      await updateRoom(selectedNode.data.roomId, {
        name: selectedNode.data.name,
        description: selectedNode.data.description,
      })
      alert('Room saved!')
    } catch (err) {
      console.error('Failed to save room:', err)
      alert('Failed to save room')
    }
  }

  const handleRefresh = () => {
    refreshRooms()
  }

  // Handle node position changes (dragging)
  const handleNodesChange = useCallback((changes: unknown) => {
    (changes as Array<unknown>).forEach((change: unknown) => {
      const c = change as { type?: string; id?: string; position?: { x: number; y: number }; dragging?: boolean }
      if (c.type === 'position' && c.position && c.dragging === false) {
        // Node drag finished - persist position
        updateNodePosition(c.id!, c.position)
      }
    })
  }, [updateNodePosition])

  // Loading state
  if (loading) {
    return (
      <div className="management-page">
        <div className="page-header">
          <h2>Map Builder</h2>
        </div>
        <div style={{ padding: '40px', textAlign: 'center' }}>
          <div className="loading-spinner"></div>
          <p>Loading rooms from API...</p>
        </div>
      </div>
    )
  }

  // Error state
  if (error) {
    return (
      <div className="management-page">
        <div className="page-header">
          <h2>Map Builder</h2>
          <button onClick={handleRefresh}>Retry</button>
        </div>
        <div style={{ padding: '40px', textAlign: 'center', color: '#ff6b6b' }}>
          <p>Error: {error}</p>
          <p style={{ fontSize: '12px', color: '#888' }}>Falling back to demo data</p>
        </div>
      </div>
    )
  }

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Map Builder</h2>
        <div className="map-actions">
          <button onClick={addNewRoom}>Add Room</button>
          <button>Connect Rooms</button>
          <button onClick={handleSaveRoom} disabled={!selectedNode}>Save Map</button>
          <button onClick={handleRefresh} style={{ marginLeft: '8px' }}>Refresh</button>
        </div>
      </div>

      <div className="map-container" style={{ display: 'flex', gap: '16px' }}>
        <SidebarPalette />
        <div 
          className="map-flow" 
          style={{ flex: 1 }} 
          ref={reactFlowWrapper}
          onDragOver={handleDragOver}
          onDrop={handleDrop}
        >
          <MapFlow
            nodes={nodes}
            edges={edges}
            onConnect={onConnect}
            onNodeClick={onNodeClick}
            onNodesChange={handleNodesChange}
          />
        </div>
      </div>

      <div className="map-legend" style={{ marginTop: '12px', padding: '8px', background: '#222', borderRadius: '4px' }}>
        <span><strong>Legend:</strong></span>
        <span style={{ marginLeft: '16px' }}>Room Node</span>
        <span style={{ marginLeft: '16px' }}>Exit Connection</span>
        <span style={{ marginLeft: '16px' }}>Selected</span>
        <span style={{ marginLeft: '16px', color: '#4ade80' }}>Starting Room</span>
      </div>

      {selectedNode && (
        <div className="map-sidebar" style={{ 
          position: 'fixed',
          right: '20px',
          top: '120px',
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
          <div style={{ marginBottom: '12px', fontSize: '12px', color: '#888' }}>
            <p>Room ID: {getRoomData(selectedNode).roomId}</p>
            <p>Starting Room: {getRoomData(selectedNode).isStartingRoom ? 'Yes' : 'No'}</p>
            <p>Atmosphere: {getRoomData(selectedNode).atmosphere}</p>
          </div>
          <p style={{ fontSize: '12px', color: '#666' }}>Node ID: {selectedNode.id}</p>
        </div>
      )}
    </div>
  )
}

function MapBuilder() {
  return (
    <ReactFlowProvider>
      <MapBuilderContent />
    </ReactFlowProvider>
  )
}