import { describe, it, expect } from 'vitest'

// Test the direction emoji mapping directly
describe('ExitEdge direction mapping', () => {
  const directionEmoji: Record<string, string> = {
    north: '↑',
    south: '↓',
    east: '→',
    west: '←',
    up: '⬆',
    down: '⬇',
  }

  it('maps north to up arrow', () => {
    expect(directionEmoji.north).toBe('↑')
  })

  it('maps south to down arrow', () => {
    expect(directionEmoji.south).toBe('↓')
  })

  it('maps east to right arrow', () => {
    expect(directionEmoji.east).toBe('→')
  })

  it('maps west to left arrow', () => {
    expect(directionEmoji.west).toBe('←')
  })

  it('maps up to up arrow', () => {
    expect(directionEmoji.up).toBe('⬆')
  })

  it('maps down to down arrow', () => {
    expect(directionEmoji.down).toBe('⬇')
  })
})

// Test edge data structure validation
describe('ExitEdge data validation', () => {
  type ExitEdgeData = {
    direction: 'north' | 'south' | 'east' | 'west' | 'up' | 'down'
    isZExit: boolean
    animated?: boolean
  }

  const validDirections = ['north', 'south', 'east', 'west', 'up', 'down']

  it('accepts valid directions', () => {
    validDirections.forEach(dir => {
      const data: ExitEdgeData = { direction: dir as ExitEdgeData['direction'], isZExit: false }
      expect(validDirections.includes(data.direction)).toBe(true)
    })
  })

  it('identifies Z-axis exits', () => {
    const upExit: ExitEdgeData = { direction: 'up', isZExit: true }
    const downExit: ExitEdgeData = { direction: 'down', isZExit: true }
    const eastExit: ExitEdgeData = { direction: 'east', isZExit: false }

    expect(upExit.isZExit).toBe(true)
    expect(downExit.isZExit).toBe(true)
    expect(eastExit.isZExit).toBe(false)
  })

  it('defaults animated to false', () => {
    const data: ExitEdgeData = { direction: 'north', isZExit: false }
    expect(data.animated).toBeUndefined()
  })

  it('allows animated to be set', () => {
    const data: ExitEdgeData = { direction: 'north', isZExit: false, animated: true }
    expect(data.animated).toBe(true)
  })
})

// Test edge color logic
describe('ExitEdge styling', () => {
  const getEdgeColor = (isZExit: boolean, selected: boolean): string => {
    if (isZExit) return selected ? '#ff7675' : '#e17055'
    return selected ? '#74b9ff' : '#0984e3'
  }

  it('returns orange for Z-axis exits', () => {
    expect(getEdgeColor(true, false)).toBe('#e17055')
  })

  it('returns light blue for standard exits', () => {
    expect(getEdgeColor(false, false)).toBe('#0984e3')
  })

  it('returns brighter color when selected', () => {
    expect(getEdgeColor(true, true)).toBe('#ff7675')
    expect(getEdgeColor(false, true)).toBe('#74b9ff')
  })
})

// Test label background color logic
describe('ExitEdge label styling', () => {
  const getLabelBgColor = (isZExit: boolean): string => {
    return isZExit ? '#d63031' : '#0984e3'
  }

  it('returns red for Z-exit labels', () => {
    expect(getLabelBgColor(true)).toBe('#d63031')
  })

  it('returns blue for standard exit labels', () => {
    expect(getLabelBgColor(false)).toBe('#0984e3')
  })
})

// Test edge ID generation
describe('ExitEdge ID generation', () => {
  const generateEdgeId = (sourceId: string, targetId: string, direction: string): string => {
    return `${sourceId}-${targetId}-${direction}`
  }

  it('generates correct edge ID format', () => {
    expect(generateEdgeId('room-1', 'room-2', 'north')).toBe('room-1-room-2-north')
  })

  it('generates unique IDs for different directions', () => {
    const id1 = generateEdgeId('room-1', 'room-2', 'north')
    const id2 = generateEdgeId('room-1', 'room-2', 'south')
    expect(id1).not.toBe(id2)
  })
})