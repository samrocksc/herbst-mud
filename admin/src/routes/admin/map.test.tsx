import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { SidebarPalette } from '../../components/SidebarPalette'

// Mock ReactFlow components
vi.mock('@xyflow/react', () => ({
  ReactFlowProvider: ({ children }: { children: React.ReactNode }) => children,
  useReactFlow: () => ({
    screenToFlowPosition: (pos: { x: number; y: number }) => ({ x: pos.x, y: pos.y })
  })
}))

describe('Map Builder - Drag and Drop Room Creation', () => {
  it('renders SidebarPalette component', () => {
    render(<SidebarPalette />)
    expect(screen.getByText('Drag to Add')).toBeInTheDocument()
    expect(screen.getByText('New Room')).toBeInTheDocument()
  })

  it('renders draggable New Room, NPC, and Item options', () => {
    render(<SidebarPalette />)
    expect(screen.getByText('New Room')).toBeInTheDocument()
    expect(screen.getByText('New NPC')).toBeInTheDocument()
    expect(screen.getByText('New Item')).toBeInTheDocument()
  })

  it('has draggable attribute on palette items', () => {
    render(<SidebarPalette />)
    const newRoom = screen.getByText('New Room').closest('div')
    expect(newRoom).toHaveAttribute('draggable', 'true')
  })
})