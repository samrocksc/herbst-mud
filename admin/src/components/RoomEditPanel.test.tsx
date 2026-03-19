import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { RoomEditPanel } from './RoomEditPanel'
import type { Node } from '@xyflow/react'

describe('RoomEditPanel', () => {
  const mockSelectedNode: Node = {
    id: 'room-1',
    type: 'room',
    position: { x: 100, y: 100 },
    data: {
      name: 'Test Room',
      description: 'A test room description',
      zLevel: 0,
    },
  }

  const mockEdges = [
    { id: 'e1', source: 'room-1', target: 'room-2', label: 'north' },
    { id: 'e2', source: 'room-1', target: 'room-3', label: 'east' },
  ]

  const mockOnUpdateNode = vi.fn()
  const mockOnDeleteNode = vi.fn()
  const mockOnClose = vi.fn()

  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders room details when node is selected', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    expect(screen.getByText('Room Details')).toBeDefined()
    expect(screen.getByDisplayValue('Test Room')).toBeDefined()
    expect(screen.getByDisplayValue('A test room description')).toBeDefined()
  })

  it('displays the node ID', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    expect(screen.getByText(/Node ID:/)).toBeDefined()
    expect(screen.getByText('room-1')).toBeDefined()
  })

  it('displays exit count', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    expect(screen.getByText(/Exits \(2\)/)).toBeDefined()
  })

  it('displays exit list with directions', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    expect(screen.getByText('NORTH')).toBeDefined()
    expect(screen.getByText('EAST')).toBeDefined()
  })

  it('calls onClose when close button clicked', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    const closeButton = screen.getByTitle('Close panel')
    fireEvent.click(closeButton)
    expect(mockOnClose).toHaveBeenCalled()
  })

  it('shows delete confirmation when delete clicked', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    const deleteButton = screen.getByText('Delete')
    fireEvent.click(deleteButton)
    
    expect(screen.getByText('Confirm')).toBeDefined()
    expect(screen.getByText('Cancel')).toBeDefined()
  })

  it('calls onDeleteNode when delete confirmed', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    fireEvent.click(screen.getByText('Delete'))
    fireEvent.click(screen.getByText('Confirm'))
    
    expect(mockOnDeleteNode).toHaveBeenCalledWith('room-1')
  })

  it('calls onUpdateNode when save clicked', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    fireEvent.click(screen.getByText('Save'))
    
    expect(mockOnUpdateNode).toHaveBeenCalled()
  })

  it('updates local state when name changes', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    const nameInput = screen.getByDisplayValue('Test Room')
    fireEvent.change(nameInput, { target: { value: 'New Room Name' } })
    
    expect(screen.getByDisplayValue('New Room Name')).toBeDefined()
  })

  it('updates local state when description changes', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    const descInput = screen.getByDisplayValue('A test room description')
    fireEvent.change(descInput, { target: { value: 'New description' } })
    
    expect(screen.getByDisplayValue('New description')).toBeDefined()
  })

  it('updates local state when z-level changes', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={mockEdges}
      />
    )

    const zSelect = screen.getByDisplayValue('Z: 0 (Ground)')
    fireEvent.change(zSelect, { target: { value: '1' } })
    
    expect(screen.getByDisplayValue('Z: 1 (Upper Floor)')).toBeDefined()
  })

  it('shows empty state when no exits', () => {
    render(
      <RoomEditPanel
        selectedNode={mockSelectedNode}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={[]}
      />
    )

    expect(screen.getByText(/No exits yet/)).toBeDefined()
  })

  it('returns null when no node selected', () => {
    const { container } = render(
      <RoomEditPanel
        selectedNode={null}
        onUpdateNode={mockOnUpdateNode}
        onDeleteNode={mockOnDeleteNode}
        onClose={mockOnClose}
        edges={[]}
      />
    )

    expect(container.firstChild).toBeNull()
  })
})