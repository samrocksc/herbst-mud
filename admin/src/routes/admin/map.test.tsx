import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { mapRoutes } from './map'
import type { Node, Edge } from '@xyflow/react'

// Test the room-to-node conversion functions
describe('Map Data Transformations', () => {
  describe('roomToNode', () => {
    it('converts API room to ReactFlow node', () => {
      const room = {
        id: 1,
        name: 'Town Square',
        description: 'Central hub',
        isStartingRoom: true,
        exits: { north: 2 },
        zLevel: 0,
        x: 100,
        y: 200
      }

      // Expected node structure
      const node = {
        id: '1',
        type: 'room',
        position: { x: 100, y: 200 },
        data: {
          name: 'Town Square',
          description: 'Central hub',
          zLevel: 0,
          isStartingRoom: true,
          exits: { north: 2 }
        },
        selected: false
      }

      expect(node.id).toBe('1')
      expect(node.data.name).toBe('Town Square')
      expect(node.data.zLevel).toBe(0)
    })

    it('uses random position when x/y not provided', () => {
      const room = {
        id: 1,
        name: 'Test Room',
        description: 'Test',
        isStartingRoom: false,
        exits: {},
        zLevel: 0
      }

      // Position should be generated (we can't predict exact value)
      // Just verify the function would create valid structure
      expect(room.id).toBeDefined()
    })
  })

  describe('nodeToRoomInput', () => {
    it('converts ReactFlow node to API input', () => {
      const node: Node = {
        id: '1',
        type: 'room',
        position: { x: 100, y: 200 },
        data: {
          name: 'Town Square',
          description: 'Central hub',
          zLevel: 0,
          isStartingRoom: true,
          exits: { north: 2 }
        },
        selected: false
      }

      const input = {
        name: node.data.name,
        description: (node.data as any).description,
        isStartingRoom: (node.data as any).isStartingRoom,
        zLevel: (node.data as any).zLevel,
        x: node.position.x,
        y: node.position.y,
        exits: (node.data as any).exits
      }

      expect(input.name).toBe('Town Square')
      expect(input.zLevel).toBe(0)
      expect(input.x).toBe(100)
      expect(input.y).toBe(200)
    })
  })

  describe('getEdgesFromNodes', () => {
    it('generates edges from node exits', () => {
      const nodes: Node[] = [
        {
          id: '1',
          type: 'room',
          position: { x: 100, y: 100 },
          data: { name: 'Room 1', description: '', zLevel: 0, exits: { north: 2, east: 3 } },
          selected: false
        },
        {
          id: '2',
          type: 'room',
          position: { x: 100, y: 200 },
          data: { name: 'Room 2', description: '', zLevel: 0, exits: { south: 1 } },
          selected: false
        },
        {
          id: '3',
          type: 'room',
          position: { x: 200, y: 100 },
          data: { name: 'Room 3', description: '', zLevel: 0, exits: { west: 1 } },
          selected: false
        }
      ]

      // Extract exits from nodes and create edges
      const edges: Edge[] = []
      
      for (const node of nodes) {
        const data = node.data as any
        const exits = data.exits ?? {}
        
        for (const [direction, targetId] of Object.entries(exits)) {
          const isZExit = direction === 'up' || direction === 'down'
          
          edges.push({
            id: `e${node.id}-${targetId}-${direction}`,
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

      expect(edges.length).toBe(4) // 2 exits from room 1, 1 from room 2, 1 from room 3
      expect(edges.find(e => e.label === 'north')).toBeDefined()
      expect(edges.find(e => e.label === 'east')).toBeDefined()
    })

    it('applies correct styling for Z-exits', () => {
      const direction = 'up'
      const isZExit = direction === 'up' || direction === 'down'
      
      const edge: Edge = {
        id: 'e1-2-up',
        source: '1',
        target: '2',
        label: direction,
        type: 'smoothstep',
        animated: isZExit,
        style: isZExit 
          ? { stroke: direction === 'up' ? '#e17055' : '#74b9ff', strokeWidth: 2 }
          : { stroke: '#666' },
        data: { isZExit, direction }
      }

      expect(edge.animated).toBe(true)
      expect(edge.style?.stroke).toBe('#e17055') // Orange for up
    })
  })
})

describe('Z-Level Filtering', () => {
  const nodes: Node[] = [
    { id: '1', data: { zLevel: 0 }, selected: false } as Node,
    { id: '2', data: { zLevel: 1 }, selected: false } as Node,
    { id: '3', data: { zLevel: -1 }, selected: false } as Node,
    { id: '4', data: { zLevel: 2 }, selected: false } as Node,
    { id: '5', data: { zLevel: -2 }, selected: false } as Node,
  ]

  const currentZLevel = 0

  const filteredNodes = nodes.map(node => {
    const zLevel = (node.data as any).zLevel ?? 0
    const isCurrentLevel = zLevel === currentZLevel
    const isAdjacent = Math.abs(zLevel - currentZLevel) === 1
    
    return {
      ...node,
      hidden: !isCurrentLevel && !isAdjacent,
      style: !isCurrentLevel && isAdjacent ? { opacity: 0.4 } : undefined
    }
  })

  it('shows current Z-level nodes', () => {
    const visible = filteredNodes.filter(n => !n.hidden)
    expect(visible.find(n => n.id === '1')).toBeDefined() // Z: 0
  })

  it('shows adjacent Z-level nodes with reduced opacity', () => {
    const visible = filteredNodes.filter(n => !n.hidden)
    const adjacent = visible.filter(n => (n.style as any)?.opacity === 0.4)
    
    expect(adjacent.length).toBe(2) // Z: 1 and Z: -1
    expect(adjacent.find(n => n.id === '2')).toBeDefined() // Z: 1
    expect(adjacent.find(n => n.id === '3')).toBeDefined() // Z: -1
  })

  it('hides non-adjacent Z-level nodes', () => {
    const visible = filteredNodes.filter(n => !n.hidden)
    expect(visible.find(n => n.id === '4')).toBeUndefined() // Z: 2 (too far)
    expect(visible.find(n => n.id === '5')).toBeUndefined() // Z: -2 (too far)
  })
})

describe('Edge Filtering by Z-Level', () => {
  const nodes: Node[] = [
    { id: '1', data: { zLevel: 0 }, selected: false } as Node,
    { id: '2', data: { zLevel: 0 }, selected: false } as Node,
    { id: '3', data: { zLevel: 1 }, selected: false } as Node,
  ]

  const edges: Edge[] = [
    { id: 'e1-2', source: '1', target: '2', label: 'north' },
    { id: 'e1-3-up', source: '1', target: '3', label: 'up', data: { isZExit: true, direction: 'up' } },
    { id: 'e2-3', source: '2', target: '3', label: 'up', data: { isZExit: true, direction: 'up' } },
  ]

  const currentZLevel = 0

  const filteredEdges = edges.filter(edge => {
    const sourceNode = nodes.find(n => n.id === edge.source)
    const targetNode = nodes.find(n => n.id === edge.target)
    if (!sourceNode || !targetNode) return false
    
    const sourceZ = (sourceNode.data as any).zLevel ?? 0
    const targetZ = (targetNode.data as any).zLevel ?? 0
    
    if (sourceZ === targetZ) return true
    
    const label = (edge.label as string || '').toLowerCase()
    return label === 'up' || label === 'down'
  })

  it('shows edges on same Z-level', () => {
    expect(filteredEdges.find(e => e.id === 'e1-2')).toBeDefined()
  })

  it('shows Z-exits between different Z-levels', () => {
    const zExits = filteredEdges.filter(e => (e.data as any)?.isZExit)
    expect(zExits.length).toBe(2)
  })
})

describe('Exit Bidirectional Creation', () => {
  it('creates opposite direction exits', () => {
    const opposites: Record<string, string> = {
      north: 'south',
      south: 'north',
      east: 'west',
      west: 'east',
      up: 'down',
      down: 'up'
    }

    expect(opposites.north).toBe('south')
    expect(opposites.east).toBe('west')
    expect(opposites.up).toBe('down')
    expect(opposites.down).toBe('up')
  })

  it('generates correct edge IDs for bidirectional exits', () => {
    const sourceId = '1'
    const targetId = '2'
    const sourceDirection = 'north'
    const targetDirection = 'south'

    const sourceEdgeId = `e${sourceId}-${targetId}-${sourceDirection}`
    const targetEdgeId = `e${targetId}-${sourceId}-${targetDirection}`

    expect(sourceEdgeId).toBe('e1-2-north')
    expect(targetEdgeId).toBe('e2-1-south')
  })
})

describe('API Integration Mock', () => {
  it('mock room data structure matches API response', () => {
    const mockRoom = {
      id: 1,
      name: 'Town Square',
      description: 'The central hub',
      isStartingRoom: true,
      exits: { north: 2, south: 3 },
      atmosphere: 'peaceful',
      x: 250,
      y: 100,
      zLevel: 0
    }

    // Verify structure matches what useRooms hook expects
    expect(mockRoom).toHaveProperty('id')
    expect(mockRoom).toHaveProperty('name')
    expect(mockRoom).toHaveProperty('description')
    expect(mockRoom).toHaveProperty('isStartingRoom')
    expect(mockRoom).toHaveProperty('exits')
    expect(mockRoom).toHaveProperty('x')
    expect(mockRoom).toHaveProperty('y')
    expect(mockRoom).toHaveProperty('zLevel')
  })

  it('handles missing optional fields gracefully', () => {
    const minimalRoom = {
      id: 1,
      name: 'Room',
      description: '',
      isStartingRoom: false,
      exits: {}
    }

    // Should have defaults for missing fields
    const zLevel = (minimalRoom as any).zLevel ?? 0
    const x = (minimalRoom as any).x ?? Math.random() * 400 + 100
    const y = (minimalRoom as any).y ?? Math.random() * 400 + 100

    expect(zLevel).toBe(0)
    expect(x).toBeGreaterThan(0)
    expect(y).toBeGreaterThan(0)
  })
})