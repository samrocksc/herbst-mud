import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { ZLevelSelector } from './ZLevelSelector'

describe('ZLevelSelector', () => {
  const mockOnChange = vi.fn()

  beforeEach(() => {
    mockOnChange.mockClear()
  })

  it('renders all z-level tabs', () => {
    render(<ZLevelSelector currentLevel={0} onChange={mockOnChange} />)
    
    expect(screen.getByText('Z: -2')).toBeInTheDocument()
    expect(screen.getByText('Z: -1')).toBeInTheDocument()
    expect(screen.getByText('Z: 0')).toBeInTheDocument()
    expect(screen.getByText('Z: 1')).toBeInTheDocument()
    expect(screen.getByText('Z: 2')).toBeInTheDocument()
  })

  it('highlights current level', () => {
    render(<ZLevelSelector currentLevel={0} onChange={mockOnChange} />)
    
    const level0 = screen.getByText('Z: 0').closest('button')
    expect(level0).toHaveClass('active')
  })

  it('calls onChange when level clicked', () => {
    render(<ZLevelSelector currentLevel={0} onChange={mockOnChange} />)
    
    fireEvent.click(screen.getByText('Z: 1'))
    expect(mockOnChange).toHaveBeenCalledWith(1)
  })

  it('displays level labels correctly', () => {
    render(<ZLevelSelector currentLevel={0} onChange={mockOnChange} />)
    
    expect(screen.getByText('Ground')).toBeInTheDocument()
    expect(screen.getByText('Underground')).toBeInTheDocument()
    expect(screen.getByText('Upper Floor')).toBeInTheDocument()
  })
})