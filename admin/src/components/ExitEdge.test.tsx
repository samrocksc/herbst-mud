import { describe, it, expect } from 'vitest'
import { createExitEdge, directionColor } from './ExitEdge'

describe('ExitEdge', () => {
  describe('createExitEdge', () => {
    it('should create a basic exit edge', () => {
      const edge = createExitEdge('edge-1', 'room-1', 'room-2', 'north')
      
      expect(edge.id).toBe('edge-1')
      expect(edge.source).toBe('room-1')
      expect(edge.target).toBe('room-2')
      expect(edge.type).toBe('exit')
      expect(edge.data.direction).toBe('north')
      expect(edge.data.isZExit).toBe(false)
      expect(edge.animated).toBe(true)
    })

    it('should set isZExit for up/down directions', () => {
      const upEdge = createExitEdge('edge-up', 'room-1', 'room-2', 'up')
      const downEdge = createExitEdge('edge-down', 'room-1', 'room-2', 'down')
      
      expect(upEdge.data.isZExit).toBe(true)
      expect(downEdge.data.isZExit).toBe(true)
    })

    it('should allow custom label', () => {
      const edge = createExitEdge(
        'edge-1',
        'room-1',
        'room-2',
        'north',
        'To Town Square'
      )
      
      expect(edge.data.label).toBe('To Town Square')
    })
  })

  describe('directionColor', () => {
    it('should return teal for north/south', () => {
      expect(directionColor('north')).toBe('#4ecdc4')
      expect(directionColor('south')).toBe('#4ecdc4')
    })

    it('should return yellow for east/west', () => {
      expect(directionColor('east')).toBe('#ffe66d')
      expect(directionColor('west')).toBe('#ffe66d')
    })

    it('should return red for up/down', () => {
      expect(directionColor('up')).toBe('#ff6b6b')
      expect(directionColor('down')).toBe('#ff6b6b')
    })

    it('should return red for Z-axis exits', () => {
      expect(directionColor('north', true)).toBe('#ff6b6b')
    })
  })
})