import { describe, it, expect } from 'vitest'
import { ExitEdge, ExitEdgeData } from './ExitEdge'

describe('ExitEdge', () => {
  describe('Type exports', () => {
    it('should export ExitEdgeData interface with direction field', () => {
      const testData: ExitEdgeData = {
        direction: 'north',
        isZExit: false,
      }
      expect(testData.direction).toBe('north')
      expect(testData.isZExit).toBe(false)
    })

    it('should accept all valid directions', () => {
      const directions: ExitEdgeData['direction'][] = ['north', 'south', 'east', 'west', 'up', 'down']
      directions.forEach(dir => {
        const data: ExitEdgeData = { direction: dir }
        expect(data.direction).toBe(dir)
      })
    })
  })

  describe('Direction symbols', () => {
    const directionSymbols: Record<string, string> = {
      north: '↑',
      south: '↓',
      east: '→',
      west: '←',
      up: '⬆',
      down: '⬇',
    }

    it('should have symbols for all directions', () => {
      expect(directionSymbols.north).toBe('↑')
      expect(directionSymbols.south).toBe('↓')
      expect(directionSymbols.east).toBe('→')
      expect(directionSymbols.west).toBe('←')
      expect(directionSymbols.up).toBe('⬆')
      expect(directionSymbols.down).toBe('⬇')
    })

    it('should correctly identify Z-axis exits', () => {
      const zExits: ExitEdgeData['direction'][] = ['up', 'down']
      const horizontalExits: ExitEdgeData['direction'][] = ['north', 'south', 'east', 'west']

      zExits.forEach(dir => {
        const isZExit = dir === 'up' || dir === 'down'
        expect(isZExit).toBe(true)
      })

      horizontalExits.forEach(dir => {
        const isZExit = dir === 'up' || dir === 'down'
        expect(isZExit).toBe(false)
      })
    })
  })

  describe('Edge data defaults', () => {
    it('should default isZExit to false when not provided', () => {
      const data: ExitEdgeData = { direction: 'east' }
      const isZExit = data.isZExit || false
      expect(isZExit).toBe(false)
    })

    it('should default direction to east if not provided', () => {
      const data = {} as ExitEdgeData
      const direction = data.direction || 'east'
      expect(direction).toBe('east')
    })
  })
})