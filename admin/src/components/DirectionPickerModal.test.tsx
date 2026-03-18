import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { DirectionPickerModal } from './DirectionPickerModal'

describe('DirectionPickerModal', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders all direction options', () => {
    const onSelect = vi.fn()
    const onCancel = vi.fn()
    
    render(
      <DirectionPickerModal
        sourceName="Room A"
        targetName="Room B"
        onSelect={onSelect}
        onCancel={onCancel}
      />
    )
    
    expect(screen.getByText('Select Exit Direction')).toBeInTheDocument()
    expect(screen.getByText('Room A')).toBeInTheDocument()
    expect(screen.getByText('Room B')).toBeInTheDocument()
    
    // Check all directions are available
    expect(screen.getByText('North')).toBeInTheDocument()
    expect(screen.getByText('South')).toBeInTheDocument()
    expect(screen.getByText('East')).toBeInTheDocument()
    expect(screen.getByText('West')).toBeInTheDocument()
    expect(screen.getByText('Up')).toBeInTheDocument()
    expect(screen.getByText('Down')).toBeInTheDocument()
  })

  it('calls onSelect with direction when confirm clicked after selecting', () => {
    const onSelect = vi.fn()
    const onCancel = vi.fn()
    
    render(
      <DirectionPickerModal
        sourceName="Room A"
        targetName="Room B"
        onSelect={onSelect}
        onCancel={onCancel}
      />
    )
    
    // Select direction
    fireEvent.click(screen.getByText('North'))
    // Click confirm
    fireEvent.click(screen.getByText('Connect (Bidirectional)'))
    
    expect(onSelect).toHaveBeenCalledWith('north', 'south')
  })

  it('calls onCancel when Cancel button clicked', () => {
    const onSelect = vi.fn()
    const onCancel = vi.fn()
    
    render(
      <DirectionPickerModal
        sourceName="Room A"
        targetName="Room B"
        onSelect={onSelect}
        onCancel={onCancel}
      />
    )
    
    fireEvent.click(screen.getByText('Cancel'))
    
    expect(onCancel).toHaveBeenCalled()
    expect(onSelect).not.toHaveBeenCalled()
  })

  it('selects opposite direction for bidirectional exits', () => {
    const onSelect = vi.fn()
    
    render(
      <DirectionPickerModal
        sourceName="Room A"
        targetName="Room B"
        onSelect={onSelect}
        onCancel={() => {}}
      />
    )
    
    // Select direction and confirm
    fireEvent.click(screen.getByText('East'))
    fireEvent.click(screen.getByText('Connect (Bidirectional)'))
    
    // East from A to B, West from B to A
    expect(onSelect).toHaveBeenCalledWith('east', 'west')
  })
})