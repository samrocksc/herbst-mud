import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { SidebarPalette } from './SidebarPalette'

// Mock ReactFlow provider
vi.mock('@xyflow/react', () => ({
  ReactFlowProvider: ({ children }: { children: React.ReactNode }) => children,
  useReactFlow: () => ({
    screenToFlowPosition: (pos: { x: number; y: number }) => ({ x: pos.x, y: pos.y })
  })
}))

describe('SidebarPalette', () => {
  it('renders the palette header', () => {
    render(<SidebarPalette />)
    expect(screen.getByText('Drag to Add')).toBeInTheDocument()
  })

  it('renders draggable New Room option', () => {
    render(<SidebarPalette />)
    expect(screen.getByText('New Room')).toBeInTheDocument()
  })

  it('renders draggable New NPC option', () => {
    render(<SidebarPalette />)
    expect(screen.getByText('New NPC')).toBeInTheDocument()
  })

  it('renders draggable New Item option', () => {
    render(<SidebarPalette />)
    expect(screen.getByText('New Item')).toBeInTheDocument()
  })

  it('has draggable attributes on items', () => {
    render(<SidebarPalette />)
    const roomItem = screen.getByText('New Room').closest('div')
    expect(roomItem).toHaveAttribute('draggable', 'true')
  })
})