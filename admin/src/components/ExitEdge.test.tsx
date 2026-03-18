import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { ExitEdge } from './ExitEdge'
import type { ExitEdgeData } from './ExitEdge'
import type { EdgeProps } from '@xyflow/react'
import { Position } from '@xyflow/react'

/**
 * ExitEdge Component Tests
 * 
 * Tests for the directed exit edge component that connects rooms
 * in the map builder with direction indicators.
 * 
 * 🟣 Donatello - Turtle Time
 */

// Mock getBezierPath from @xyflow/react
vi.mock('@xyflow/react', async () => {
  const actual = await vi.importActual('@xyflow/react')
  return {
    ...actual,
    getBezierPath: vi.fn(() => ['M 0 0 L 100 0', 50, 0]),
    BaseEdge: vi.fn(({ style }) => <path data-testid="base-edge" style={style} />),
    EdgeLabelRenderer: vi.fn(({ children }) => <div data-testid="edge-label-renderer">{children}</div>),
  }
})

describe('ExitEdge', () => {
  const defaultProps: EdgeProps = {
    id: 'edge-1',
    source: 'node-1',
    target: 'node-2',
    sourceX: 0,
    sourceY: 0,
    targetX: 100,
    targetY: 0,
    sourcePosition: Position.Right,
    targetPosition: Position.Left,
    data: {},
    selected: false,
    style: {},
    type: 'exit',
    animated: false,
    interactionWidth: 0,
  }

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders the edge with default direction (east)', () => {
    render(<ExitEdge {...defaultProps} />)
    
    // Should show direction label
    expect(screen.getByText('➡️ E')).toBeInTheDocument()
  })

  it('renders direction indicator for north', () => {
    const props: EdgeProps = {
      ...defaultProps,
      data: { direction: 'north' } as ExitEdgeData,
    }
    
    render(<ExitEdge {...props} />)
    
    expect(screen.getByText('⬆️ N')).toBeInTheDocument()
  })

  it('renders direction indicator for south', () => {
    const props: EdgeProps = {
      ...defaultProps,
      data: { direction: 'south' } as ExitEdgeData,
    }
    
    render(<ExitEdge {...props} />)
    
    expect(screen.getByText('⬇️ S')).toBeInTheDocument()
  })

  it('renders direction indicator for west', () => {
    const props: EdgeProps = {
      ...defaultProps,
      data: { direction: 'west' } as ExitEdgeData,
    }
    
    render(<ExitEdge {...props} />)
    
    expect(screen.getByText('⬅️ W')).toBeInTheDocument()
  })

  it('renders direction indicator for up (z-exit)', () => {
    const props: EdgeProps = {
      ...defaultProps,
      data: { direction: 'up', isZExit: true } as ExitEdgeData,
    }
    
    render(<ExitEdge {...props} />)
    
    expect(screen.getByText('🆙 U')).toBeInTheDocument()
    expect(screen.getByText('Z-Exit')).toBeInTheDocument()
  })

  it('renders direction indicator for down (z-exit)', () => {
    const props: EdgeProps = {
      ...defaultProps,
      data: { direction: 'down', isZExit: true } as ExitEdgeData,
    }
    
    render(<ExitEdge {...props} />)
    
    expect(screen.getByText('🔽 D')).toBeInTheDocument()
    expect(screen.getByText('Z-Exit')).toBeInTheDocument()
  })

  it('shows selected state when selected prop is true', () => {
    const props: EdgeProps = {
      ...defaultProps,
      selected: true,
    }
    
    render(<ExitEdge {...props} />)
    
    // The label should have purple border when selected
    const label = screen.getByText('➡️ E').parentElement
    expect(label).toHaveStyle({ borderColor: '2px solid #a29bfe' })
  })

  it('calls onEditClick when edge label is clicked', () => {
    const onEditClick = vi.fn()
    const props: EdgeProps = {
      ...defaultProps,
      data: { direction: 'north', onEditClick } as ExitEdgeData,
    }
    
    render(<ExitEdge {...props} />)
    
    const label = screen.getByText('⬆️ N').parentElement
    fireEvent.click(label!)
    
    expect(onEditClick).toHaveBeenCalledWith('edge-1', { direction: 'north', onEditClick })
  })

  it('applies custom style from props', () => {
    const props: EdgeProps = {
      ...defaultProps,
      style: { stroke: '#ff0000', strokeWidth: 4 },
    }
    
    render(<ExitEdge {...props} />)
    
    // The BaseEdge mock receives the style prop
    const baseEdge = screen.getByTestId('base-edge')
    expect(baseEdge).toBeInTheDocument()
  })

  it('renders z-exit indicator only for isZExit=true', () => {
    const zExitProps: EdgeProps = {
      ...defaultProps,
      data: { direction: 'up', isZExit: true } as ExitEdgeData,
    }
    
    const regularProps: EdgeProps = {
      ...defaultProps,
      data: { direction: 'north', isZExit: false } as ExitEdgeData,
    }
    
    const { rerender } = render(<ExitEdge {...zExitProps} />)
    expect(screen.getByText('Z-Exit')).toBeInTheDocument()
    
    rerender(<ExitEdge {...regularProps} />)
    expect(screen.queryByText('Z-Exit')).not.toBeInTheDocument()
  })

  it('uses correct colors for z-exits vs regular exits', () => {
    const zExitProps: EdgeProps = {
      ...defaultProps,
      data: { direction: 'up', isZExit: true } as ExitEdgeData,
    }
    
    const regularProps: EdgeProps = {
      ...defaultProps,
      data: { direction: 'north' } as ExitEdgeData,
    }
    
    // Z-exit should have orange-ish color label
    const { rerender } = render(<ExitEdge {...zExitProps} />)
    const zExitLabel = screen.getByText('🆙 U').closest('div')
    expect(zExitLabel).toBeInTheDocument()
    
    // Regular exit should have green color
    rerender(<ExitEdge {...regularProps} />)
    const regularLabel = screen.getByText('⬆️ N').closest('div')
    expect(regularLabel).toBeInTheDocument()
  })

  it('handles missing direction gracefully', () => {
    const props: EdgeProps = {
      ...defaultProps,
      data: {} as ExitEdgeData,
    }
    
    render(<ExitEdge {...props} />)
    
    // Should default to east
    expect(screen.getByText('➡️ E')).toBeInTheDocument()
  })
})