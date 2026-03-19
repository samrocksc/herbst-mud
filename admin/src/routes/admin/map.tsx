import { createFileRoute } from '@tanstack/react-router'
import { useState, useCallback, useMemo, useEffect } from 'react'
import { MapFlow } from '../../components/MapFlow'
import { ZLevelSelector } from '../../components/ZLevelSelector'
import { DirectionPickerModal } from '../../components/DirectionPickerModal'
import { useRooms, useCreateRoom, useUpdateRoom, useDeleteRoom, Room, getOppositeDirection, DIRECTIONS, Direction } from '../../hooks/useRooms'
import type { Node, Edge, Connection } from '@xyflow/react'

export const Route = createFileRoute('/admin/map')({
  component: MapBuilder,
})

interface MapRoomData extends Record<string, unknown> {
  name: string
  description: string
  zLevel: number
  isStartingRoom?: boolean
  exits?: Record<string, number>
}

// Convert API Room to ReactFlow Node
function roomToNode(room: Room): Node {
  return {
    id: String(room.id),
    type: 'room',
    position: { 
      x: room.x ?? Math.random() * 400 + 100, 
      y: room.y ?? Math.random() * 400 + 100 
    },
    data: {
      name: room.name,
      description: room.description,
      zLevel: room.zLevel ?? 0,
      isStartingRoom: room.isStartingRoom,
      exits: room.exits ?? {}
    },
    selected: false
  }
}

// Convert Node data back to RoomInput for API
function nodeToRoomInput(node: Node): Partial<Room> {
  const data = node.data as MapRoomData
  return {
    name: data.name,
    description: data.description,
    isStartingRoom: data.isStartingRoom,
    zLevel: data.zLevel,
    x: node.position.x,
    y: node.position.y,
    exits: data.exits
  }
}

// Convert nodes/edges to exit structure and back
function nodesToExits(nodes: Node[]): Record<string, Record<string, number>> {
  const exits: Record<string, Record<string, number>> = {}
  
  for (const node of nodes) {
    const data = node.data as MapRoomData
    const nodeExits: Record<string, number> = {}
    
    // Find all edges originating from this node
    for (const n of nodes) {
      const targetData = n.data as MapRoomData
      const edge = nodes.find(e => {
        const eData = e.data as MapRoomData
        return (eData as any).sourceId === node.id && (eData as any).targetId === n.id
      })
    }
  }
  
  return exits
}

// Get edges from nodes (reconstruct exits as edges)
function getEdgesFromNodes(nodes: Node[]): Edge[] {
  const edges: Edge[] = []
  const nodeMap = new Map(nodes.map(n => [n.id, n]))
  
  for (const node of nodes) {
    const data = node.data as MapRoomData
    const exits = data.exits ?? {}
    
    for (const [direction, targetId] of Object.entries(exits)) {
      const targetNode = nodeMap.get(String(targetId))
      if (!targetNode) continue
      
      const targetData = targetNode.data as MapRoomData
      const isZExit = direction === 'up' || direction === 'down'
      
      edges.push({
        id: `e${node.id}-${targetNode.id}-${direction}`,
        source: node.id,
        target: String(targetId),
        label: direction,
        type: 'smoothstep',
        animated: isZExit,
        style: isZExit 
          ? { stroke: direction === 'up' ? '#e17055' : '#74b9ff', strokeWidth: 2 }
          : { stroke: '#666' },
        data: { isZExit, direction }
      })
    }
  }
  
  return edges
}

// Sample initial nodes (fallback if API unavailable)
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
]

const initialEdges: Edge[] = []

function MapBuilder() {
  const { data: rooms, isLoading, error } = useRooms()
  const createRoom = useCreateRoom()
  const updateRoom = useUpdateRoom()
  const deleteRoom = useDeleteRoom()
  
  const [nodes, setNodes] = useState<Node[]>(initialNodes)
  const [edges, setEdges] = useState<Edge[]>(initialEdges)
  const [currentZLevel, setCurrentZLevel] = useState(0)
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)
  const [isInitialized, setIsInitialized] = useState(false)
  
  // Direction picker modal state
  const [pendingConnection, setPendingConnection] = useState<{
    source: string
    target: string
    sourceName: string
    targetName: string
  } | null>(null)

  // Load rooms from API when available
  useEffect(() => {
    if (rooms && !isInitialized) {
      const loadedNodes = rooms.map(roomToNode)
      const loadedEdges = getEdgesFromNodes(loadedNodes)
      setNodes(loadedNodes)
      setEdges(loadedEdges)
      setIsInitialized(true)
    }
  }, [rooms, isInitialized])

  // Filter nodes by Z-level (show current level prominently, adjacent faintly)
  const filteredNodes = useMemo(() => {
    return nodes.map(node => {
      const zLevel = (node.data as MapRoomData).zLevel ?? 0
      const isCurrentLevel = zLevel === currentZLevel
      const isAdjacent = Math.abs(zLevel - currentZLevel) === 1
      
      return {
        ...node,
        selected: isCurrentLevel && node.selected,
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
      
      if (sourceZ === targetZ) return true
      
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
        setPendingConnection({
          source: connection.source,
          target: connection.target,
          sourceName: (sourceNode.data as MapRoomData).name,
          targetName: (targetNode.data as MapRoomData).name,
        })
      }
    }
  }, [nodes])

  // Handle direction selection - create bidirectional exits and sync to API
  const handleDirectionSelect = useCallback(async (sourceDirection: string, targetDirection: string) => {
    if (!pendingConnection) return
    
    // Add edges locally
    const newEdges: Edge[] = [
      {
        id: `e${pendingConnection.source}-${pendingConnection.target}-${sourceDirection}`,
        source: pendingConnection.source,
        target: pendingConnection.target,
        label: sourceDirection,
        type: 'smoothstep',
        animated: sourceDirection === 'up' || sourceDirection === 'down',
        style: sourceDirection === 'up' || sourceDirection === 'down' 
          ? { stroke: sourceDirection === 'up' ? '#e17055' : '#74b9ff', strokeWidth: 2 }
          : { stroke: '#666' },
        data: { isZExit: sourceDirection === 'up' || sourceDirection === 'down', direction: sourceDirection },
      },
      {
        id: `e${pendingConnection.target}-${pendingConnection.source}-${targetDirection}`,
        source: pendingConnection.target,
        target: pendingConnection.source,
        label: targetDirection,
        type: 'smoothstep',
        animated: targetDirection === 'up' || targetDirection === 'down',
        style: targetDirection === 'up' || targetDirection === 'down' 
          ? { stroke: targetDirection === 'up' ? '#e17055' : '#74b9ff', strokeWidth: 2 }
          : { stroke: '#666' },
        data: { isZExit: targetDirection === 'up' || targetDirection === 'down', direction: targetDirection },
      },
    ]
    
    // Update both nodes with the new exits
    const sourceNode = nodes.find(n => n.id === pendingConnection.source)
    const targetNode = nodes.find(n => n.id === pendingConnection.target)
    
    if (sourceNode && targetNode) {
      const sourceData = { ...(sourceNode.data as MapRoomData) }
      const targetData = { ...(targetNode.data as MapRoomData) }
      
      const sourceExits = sourceData.exits ?? {}
      const targetExits = targetData.exits ?? {}
      
      // Get the target IDs (convert to numbers for API)
      const targetIdNum = parseInt(pendingConnection.target)
      const sourceIdNum = parseInt(pendingConnection.source)
      
      sourceExits[sourceDirection] = targetIdNum
      targetExits[targetDirection] = sourceIdNum
      
      sourceData.exits = sourceExits
      targetData.exits = targetExits
      
      // Update local state first
      setNodes(nds => nds.map(n => {
        if (n.id === pendingConnection.source) return { ...n, data: sourceData }
        if (n.id === pendingConnection.target) return { ...n, data: targetData }
        return n
      }))
      setEdges(eds => [...eds, ...newEdges])
      
      // Sync to API
      try {
        await updateRoom.mutateAsync({ id: sourceIdNum, room: nodeToRoomInput({ ...sourceNode, data: sourceData }) })
        await updateRoom.mutateAsync({ id: targetIdNum, room: nodeToRoomInput({ ...targetNode, data: targetData }) })
      } catch (err) {
        console.error('Failed to save exits:', err)
      }
    }
    
    setPendingConnection(null)
  }, [pendingConnection, nodes, updateRoom])

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

  // Add new room - create in API first
  const addNewRoom = useCallback(async () => {
    const newRoomData = {
      name: `New Room ${nodes.length + 1}`,
      description: 'A new room',
      zLevel: currentZLevel,
      x: Math.random() * 400 + 100,
      y: Math.random() * 400 + 100,
      exits: {}
    }
    
    try {
      const created = await createRoom.mutateAsync(newRoomData)
      const newNode = roomToNode(created)
      setNodes(nds => [...nds, newNode])
    } catch (err) {
      console.error('Failed to create room:', err)
      // Fallback to local-only creation
      const newId = String(Math.max(...nodes.map(n => parseInt(n.id)), 0) + 1)
      const newNode: Node = {
        id: newId,
        type: 'room',
        position: { x: Math.random() * 400 + 100, y: Math.random() * 400 + 100 },
        data: { name: `Room ${newId}`, description: 'New room', zLevel: currentZLevel, exits: {} }
      }
      setNodes(nds => [...nds, newNode])
    }
  }, [nodes, currentZLevel, createRoom])

  // Save all changes
  const saveMap = useCallback(async () => {
    for (const node of nodes) {
      const roomId = parseInt(node.id)
      if (isNaN(roomId)) continue
      
      try {
        await updateRoom.mutateAsync({ 
          id: roomId, 
          room: nodeToRoomInput(node) 
        })
      } catch (err) {
        console.error(`Failed to save room ${roomId}:`, err)
      }
    }
  }, [nodes, updateRoom])

  // Delete selected room
  const handleDeleteRoom = useCallback(async () => {
    if (!selectedNode) return
    
    const roomId = parseInt(selectedNode.id)
    if (isNaN(roomId)) {
      // Local-only node, just remove
      setNodes(nds => nds.filter(n => n.id !== selectedNode.id))
      setSelectedNode(null)
      return
    }
    
    try {
      await deleteRoom.mutateAsync(roomId)
      setNodes(nds => nds.filter(n => n.id !== selectedNode.id))
      setEdges(eds => eds.filter(e => e.source !== selectedNode.id && e.target !== selectedNode.id))
      setSelectedNode(null)
    } catch (err) {
      console.error('Failed to delete room:', err)
    }
  }, [selectedNode, deleteRoom])

  // Update node in API when edited
  const handleNodeUpdate = useCallback(async (updatedNode: Node) => {
    const roomId = parseInt(updatedNode.id)
    if (isNaN(roomId)) return
    
    setNodes(nds => nds.map(n => 
      n.id === updatedNode.id ? updatedNode : n
    ))
    setSelectedNode(updatedNode)
    
    // Debounce API update
    try {
      await updateRoom.mutateAsync({ 
        id: roomId, 
        room: nodeToRoomInput(updatedNode) 
      })
    } catch (err) {
      console.error('Failed to update room:', err)
    }
  }, [updateRoom])

  const getRoomData = (node: Node | null): MapRoomData => {
    if (!node) return { name: '', description: '', zLevel: 0 }
    return node.data as MapRoomData
  }

  if (isLoading) {
    return (
      <div className="management-page">
        <div className="page-header">
          <h2>Map Builder</h2>
        </div>
        <div style={{ padding: '20px', textAlign: 'center' }}>
          <p>Loading rooms from API...</p>
        </div>
      </div>
    )
  }

  if (error) {
    console.error('Failed to load rooms:', error)
  }

  return (
    <div className="management-page">
      <div className="page-header">
        <h2>Map Builder</h2>
        <div className="map-actions">
          <button onClick={addNewRoom} disabled={createRoom.isPending}>
            {createRoom.isPending ? 'Creating...' : '+ Add Room'}
          </button>
          <button onClick={saveMap} disabled={updateRoom.isPending}>
            {updateRoom.isPending ? 'Saving...' : '💾 Save Map'}
          </button>
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
                  const updatedNode = { ...selectedNode, data: newData }
                  handleNodeUpdate(updatedNode)
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
                  const updatedNode = { ...selectedNode, data: newData }
                  handleNodeUpdate(updatedNode)
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
                  const updatedNode = { ...selectedNode, data: newData }
                  handleNodeUpdate(updatedNode)
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
            <div style={{ marginBottom: '12px' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: '8px' }}>
                <input 
                  type="checkbox"
                  checked={getRoomData(selectedNode).isStartingRoom ?? false}
                  onChange={(e) => {
                    const newData = { ...getRoomData(selectedNode), isStartingRoom: e.target.checked }
                    const updatedNode = { ...selectedNode, data: newData }
                    handleNodeUpdate(updatedNode)
                  }}
                />
                Starting Room
              </label>
            </div>
            <div style={{ marginBottom: '12px' }}>
              <label style={{ display: 'block', marginBottom: '4px' }}>Exits:</label>
              <div style={{ fontSize: '12px', color: '#aaa', fontFamily: 'monospace' }}>
                {JSON.stringify(getRoomData(selectedNode).exits ?? {}, null, 2)}
              </div>
            </div>
            <p style={{ fontSize: '12px', color: '#888' }}>Node ID: {selectedNode.id}</p>
            <button 
              onClick={handleDeleteRoom}
              disabled={deleteRoom.isPending}
              style={{ 
                marginTop: '8px', 
                background: '#c0392b', 
                color: '#fff', 
                border: 'none', 
                padding: '8px 16px', 
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              {deleteRoom.isPending ? 'Deleting...' : '🗑️ Delete Room'}
            </button>
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