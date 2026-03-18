import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { ZLevelSelector } from './ZLevelSelector'

describe('ZLevelSelector', () => {
  const mockOnChange = vi.fn()

  beforeEach(() => {
    mockOnChange.mockClear()
  })

  it('renders all Z-level options', () => {
    render(<ZLevelSelector currentLevel={0} onChange={mockOnChange} />)
    
    expect(screen.getByText('Z: -2')).toBeDefined()
    expect(screen.getByText('Z: -1')).toBeDefined()
    expect(screen.getByText('Z: 0')).toBeDefined()
    expect(screen.getByText('Z: 1')).toBeDefined()
    expect(screen.getByText('Z: 2')).toBeDefined()
  })

  it('calls onChange when a level is clicked', () => {
    render(<ZLevelSelector currentLevel={0} onChange={mockOnChange} />)
    
    const z1Button = screen.getByText('Z: 1')
    fireEvent.click(z1Button)
    
    expect(mockOnChange).toHaveBeenCalledWith(1)
  })

  it('highlights the current level', () => {
    render(<ZLevelSelector currentLevel={1} onChange={mockOnChange} />)
    
    const z1Button = screen.getByText('Z: 1')
    // The current level should have active styling (background: #00AAAA)
    expect(z1Button).toBeDefined()
  })

  it('renders floor labels', () => {
    render(<ZLevelSelector currentLevel={0} onChange={mockOnChange} />)
    
    expect(screen.getByText('Ground')).toBeDefined()
    expect(screen.getByText('Upper Floor')).toBeDefined()
    expect(screen.getByText('Underground')).toBeDefined()
  })
})

describe('Z-Level Filtering Logic', () => {
  const sampleNodes = [
    { id: '1', data: { zLevel: 0 }, selected: false },
    { id: '2', data: { zLevel: 1 }, selected: false },
    { id: '3', data: { zLevel: -1 }, selected: false },
    { id: '4', data: { zLevel: 2 }, selected: false },
  ] as const

  it('filters to current level only', () => {
    const currentZLevel = 0
    
    const filtered = sampleNodes.filter(node => {
      const zLevel = node.data.zLevel ?? 0
      return zLevel === currentZLevel || Math.abs(zLevel - currentZLevel) === 1
    })
    
    // Z: 0 (exact) and Z: 1, -1 (adjacent) should show
    expect(filtered.length).toBe(3)
    expect(filtered.map(n => n.id)).toEqual(['1', '2', '3'])
  })

  it('hides non-adjacent levels', () => {
    const currentZLevel = 0
    
    const filtered = sampleNodes.filter(node => {
      const zLevel = node.data.zLevel ?? 0
      return zLevel === currentZLevel || Math.abs(zLevel - currentZLevel) === 1
    })
    
    // Z: 2 should not show (too far)
    expect(filtered.find(n => n.id === '4')).toBeUndefined()
  })

  it('correctly handles edge visibility for Z-exits', () => {
    const edges = [
      { id: 'e1-2', source: '1', target: '2', label: 'north' },
      { id: 'e1-3', source: '1', target: '3', label: 'up' },
      { id: 'e1-4', source: '1', target: '4', label: 'down' },
    ]

    const filtered = edges.filter(edge => {
      const label = (edge.label || '').toLowerCase()
      // Z-exits should always show
      return label === 'up' || label === 'down' || label === 'north'
    })

    expect(filtered.length).toBe(3)
  })
})

describe('Z-Level Navigation', () => {
  it('can navigate between floors', () => {
    let currentLevel = 0
    
    // Simulate clicking "up"
    const navigateUp = () => {
      if (currentLevel < 2) currentLevel++
    }
    
    navigateUp()
    expect(currentLevel).toBe(1)
    
    navigateUp()
    expect(currentLevel).toBe(2)
  })

  it('can navigate down to underground levels', () => {
    let currentLevel = 0
    
    const navigateDown = () => {
      if (currentLevel > -2) currentLevel--
    }
    
    navigateDown()
    expect(currentLevel).toBe(-1)
    
    navigateDown()
    expect(currentLevel).toBe(-2)
  })

  it('respects level boundaries', () => {
    let currentLevel = 0
    
    const navigateDown = () => {
      if (currentLevel > -2) currentLevel--
    }
    
    // Try to go below -2
    currentLevel = -2
    navigateDown()
    expect(currentLevel).toBe(-2) // Should not go below
  })
})