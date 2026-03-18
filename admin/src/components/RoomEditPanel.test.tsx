import { describe, it, expect, vi } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { RoomEditPanel } from './RoomEditPanel'
import type { Node } from '@xyflow/react'

describe('RoomEditPanel', () => {
  const mockRooms = [
    { id: 1, name: 'Town Square' },
    { id: 2, name: 'Main Street' },
    { id: 3, name: 'Forest Path' }
  ]

  const mockNode: Node = {
    id: '1',
    type: 'room',
    position: { x: 0, y: 0 },
    data: {
      name: 'Town Square',
      description: 'The central hub of the village',
      zLevel: 0,
      exits: { north: 2, east: 3 }
    }
  }

  const defaultProps = {
    selectedNode: mockNode,
    rooms: mockRooms,
    onUpdate: vi.fn(),
    onDelete: vi.fn(),
    onClose: vi.fn()
  }

  it('renders null when no node is selected', () => {
    const { container } = render(
      <RoomEditPanel {...defaultProps} selectedNode={null} />
    )
    expect(container.firstChild).toBeNull()
  })

  it('displays room name in input', () => {
    render(<RoomEditPanel {...defaultProps} />)
    expect(screen.getByDisplayValue('Town Square')).toBeInTheDocument()
  })

  it('displays room description in textarea', () => {
    render(<RoomEditPanel {...defaultProps} />)
    expect(screen.getByDisplayValue('The central hub of the village')).toBeInTheDocument()
  })

  it('displays Z-level selector with correct value', () => {
    render(<RoomEditPanel {...defaultProps} />)
    const zLevelSection = screen.getByText(/Z-Level:/)
    const container = zLevelSection.parentElement
    const select = container?.querySelector('select')
    expect(select).toHaveValue('0')
  })

  it('displays existing exits', () => {
    render(<RoomEditPanel {...defaultProps} />)
    // The exits section should show the exits from mockNode.data.exits
    // north → 2, east → 3
    expect(screen.getByText(/north/)).toBeInTheDocument()
    expect(screen.getByText(/east/)).toBeInTheDocument()
  })

  it('shows target room name for exits', () => {
    render(<RoomEditPanel {...defaultProps} />)
    // north → Room 2 (Main Street)
    expect(screen.getByText(/Main Street/)).toBeInTheDocument()
  })

  it('calls onClose when close button clicked', () => {
    render(<RoomEditPanel {...defaultProps} />)
    const closeButton = screen.getByRole('button', { name: /close panel/i })
    fireEvent.click(closeButton)
    expect(defaultProps.onClose).toHaveBeenCalled()
  })

  it('updates name field on change', () => {
    render(<RoomEditPanel {...defaultProps} />)
    const nameInput = screen.getByDisplayValue('Town Square')
    fireEvent.change(nameInput, { target: { value: 'New Name' } })
    expect(screen.getByDisplayValue('New Name')).toBeInTheDocument()
  })

  it('shows add exit form when Add Exit clicked', () => {
    render(<RoomEditPanel {...defaultProps} />)
    const addButton = screen.getByRole('button', { name: /add exit/i })
    fireEvent.click(addButton)
    // After clicking, should show Add and Cancel buttons
    const addButtons = screen.getAllByRole('button', { name: /^add$/i })
    expect(addButtons.length).toBeGreaterThan(0)
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
  })

  it('calls onDelete when Delete button clicked', () => {
    render(<RoomEditPanel {...defaultProps} />)
    const deleteButton = screen.getByRole('button', { name: /delete/i })
    fireEvent.click(deleteButton)
    expect(defaultProps.onDelete).toHaveBeenCalledWith('1')
  })

  it('calls onUpdate with correct data when Save clicked', () => {
    render(<RoomEditPanel {...defaultProps} />)
    const saveButton = screen.getByRole('button', { name: /save changes/i })
    fireEvent.click(saveButton)
    expect(defaultProps.onUpdate).toHaveBeenCalledWith('1', {
      name: 'Town Square',
      description: 'The central hub of the village',
      zLevel: 0,
      exits: { north: 2, east: 3 }
    })
  })

  it('allows adding new exit via form', () => {
    render(<RoomEditPanel {...defaultProps} />)
    
    // Open add exit form
    fireEvent.click(screen.getByRole('button', { name: /add exit/i }))
    
    // Verify form is visible with direction and target dropdowns
    expect(screen.getByRole('button', { name: /cancel/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /^add$/i })).toBeInTheDocument()
    
    // Find all selects in the document - there are the Z-level select and exit form selects
    const allSelects = screen.getAllByRole('combobox')
    expect(allSelects.length).toBeGreaterThanOrEqual(2)
    
    // Cancel should close form
    fireEvent.click(screen.getByRole('button', { name: /cancel/i }))
    expect(screen.queryByRole('button', { name: /cancel/i })).not.toBeInTheDocument()
  })

  it('allows removing exits', () => {
    render(<RoomEditPanel {...defaultProps} />)
    
    // Find and click remove button for north exit
    const removeButtons = screen.getAllByRole('button', { name: /remove/i })
    fireEvent.click(removeButtons[0])
    
    // north exit should be removed
    expect(screen.queryByText(/north/)).not.toBeInTheDocument()
  })

  it('displays node ID in info section', () => {
    render(<RoomEditPanel {...defaultProps} />)
    expect(screen.getByText(/Node ID: 1/)).toBeInTheDocument()
  })

  it('shows "No exits configured" for room with no exits', () => {
    const nodeWithNoExits: Node = {
      ...mockNode,
      data: {
        name: 'Empty Room',
        description: 'No way out',
        zLevel: 0,
        exits: {}
      }
    }
    render(<RoomEditPanel {...defaultProps} selectedNode={nodeWithNoExits} />)
    expect(screen.getByText(/no exits configured/i)).toBeInTheDocument()
  })

  it('changes Z-level on selection', () => {
    render(<RoomEditPanel {...defaultProps} />)
    const zLevelSection = screen.getByText(/Z-Level:/)
    const container = zLevelSection.parentElement
    const select = container?.querySelector('select') as HTMLSelectElement
    fireEvent.change(select, { target: { value: '1' } })
    expect(select).toHaveValue('1')
  })
})