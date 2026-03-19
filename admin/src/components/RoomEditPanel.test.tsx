import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'

// Test utility functions for room edit panel
describe('Room Edit Panel - Logic Tests', () => {
  
  // Test getRoomExits logic
  it('correctly extracts exits from edges', () => {
    const nodes = [
      { id: '1', data: { name: 'Room A' } },
      { id: '2', data: { name: 'Room B' } },
    ]
    const edges = [
      { id: 'e1-2', source: '1', target: '2', label: 'north' },
      { id: 'e2-1', source: '2', target: '1', label: 'south' },
    ]
    
    const getRoomExits = (nodeId: string) => {
      const outgoingEdges = edges.filter(e => e.source === nodeId)
      return outgoingEdges.map(edge => ({
        id: edge.id,
        direction: edge.label,
        targetId: edge.target,
        targetName: nodes.find(n => n.id === edge.target)?.data.name || 'Unknown'
      }))
    }
    
    const exits = getRoomExits('1')
    expect(exits).toHaveLength(1)
    expect(exits[0].direction).toBe('north')
    expect(exits[0].targetName).toBe('Room B')
  })
  
  // Test delete room logic
  it('correctly removes room and associated edges', () => {
    const nodes = [
      { id: '1', data: { name: 'Room A' } },
      { id: '2', data: { name: 'Room B' } },
      { id: '3', data: { name: 'Room C' } },
    ]
    const edges = [
      { id: 'e1-2', source: '1', target: '2' },
      { id: 'e2-1', source: '2', target: '1' },
      { id: 'e2-3', source: '2', target: '3' },
    ]
    
    const deleteNodeId = '2'
    
    const filteredNodes = nodes.filter(n => n.id !== deleteNodeId)
    const filteredEdges = edges.filter(e => e.source !== deleteNodeId && e.target !== deleteNodeId)
    
    expect(filteredNodes).toHaveLength(2)
    expect(filteredEdges).toHaveLength(0) // All edges connected to node 2
  })
  
  // Test room data update
  it('correctly updates room data', () => {
    const roomData = { name: 'Old Name', description: 'Old Desc', zLevel: 0 }
    
    const updateRoom = (data: typeof roomData, updates: Partial<typeof roomData>) => {
      return { ...data, ...updates }
    }
    
    const updated = updateRoom(roomData, { name: 'New Name', zLevel: 1 })
    expect(updated.name).toBe('New Name')
    expect(updated.zLevel).toBe(1)
    expect(updated.description).toBe('Old Desc') // Unchanged
  })
  
  // Test exit direction icons
  it('returns correct icon for exit direction', () => {
    const getExitIcon = (direction: string) => {
      if (direction === 'up') return '↑'
      if (direction === 'down') return '↓'
      return '→'
    }
    
    expect(getExitIcon('north')).toBe('→')
    expect(getExitIcon('up')).toBe('↑')
    expect(getExitIcon('down')).toBe('↓')
  })
  
  // Test exit colors
  it('returns correct color for exit direction', () => {
    const getExitColor = (direction: string) => {
      if (direction === 'up') return '#e17055'
      if (direction === 'down') return '#74b9ff'
      return '#6c5ce7'
    }
    
    expect(getExitColor('north')).toBe('#6c5ce7')
    expect(getExitColor('up')).toBe('#e17055')
    expect(getExitColor('down')).toBe('#74b9ff')
  })
})