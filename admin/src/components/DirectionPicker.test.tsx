import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { DirectionPicker } from './DirectionPicker'

describe('DirectionPicker', () => {
  const defaultProps = {
    isOpen: true,
    sourceId: '1',
    targetId: '2',
    sourceName: 'Town Square',
    targetName: 'Main Street',
    onSelect: vi.fn(),
    onCancel: vi.fn(),
  }

  it('renders when open', () => {
    render(<DirectionPicker {...defaultProps} />)
    expect(screen.getByText('Select Exit Direction')).toBeDefined()
    expect(screen.getByText(/Town Square/)).toBeDefined()
    expect(screen.getByText(/Main Street/)).toBeDefined()
  })

  it('does not render when closed', () => {
    render(<DirectionPicker {...defaultProps} isOpen={false} />)
    expect(screen.queryByText('Select Exit Direction')).toBeNull()
  })

  it('calls onCancel when Cancel is clicked', () => {
    render(<DirectionPicker {...defaultProps} />)
    fireEvent.click(screen.getByText('Cancel'))
    expect(defaultProps.onCancel).toHaveBeenCalled()
  })

  it('calls onSelect with direction when Save Exit is clicked after selecting direction', () => {
    render(<DirectionPicker {...defaultProps} />)
    
    // Click on North direction
    fireEvent.click(screen.getByText('North'))
    
    // Click Save Exit
    fireEvent.click(screen.getByText('Save Exit'))
    
    expect(defaultProps.onSelect).toHaveBeenCalledWith('north', 'south')
  })

  it('displays all six directions', () => {
    render(<DirectionPicker {...defaultProps} />)
    expect(screen.getByText('North')).toBeDefined()
    expect(screen.getByText('South')).toBeDefined()
    expect(screen.getByText('East')).toBeDefined()
    expect(screen.getByText('West')).toBeDefined()
    expect(screen.getByText('Up')).toBeDefined()
    expect(screen.getByText('Down')).toBeDefined()
  })
})