import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'

// Mock the TanStack router hooks
vi.mock('@tanstack/react-router', () => ({
  Link: ({ children, to, className }: { children: React.ReactNode; to: string; className?: string }) => (
    <a href={to} className={className}>
      {children}
    </a>
  ),
  useLocation: () => ({ pathname: '/dashboard' }),
}))

describe('Sidebar', () => {
  it('renders all navigation items', async () => {
    const { Sidebar } = await import('./Sidebar')
    render(<Sidebar />)
    
    // Verify navigation items are present
    expect(screen.getByText('Dashboard')).toBeInTheDocument()
    expect(screen.getByText('Items')).toBeInTheDocument()
    expect(screen.getByText('Rooms')).toBeInTheDocument()
    expect(screen.getByText('Skills')).toBeInTheDocument()
    expect(screen.getByText('Map')).toBeInTheDocument()
    expect(screen.getByText('NPCs')).toBeInTheDocument()
    expect(screen.getByText('Players')).toBeInTheDocument()
  })
  
  it('renders the header with title', async () => {
    const { Sidebar } = await import('./Sidebar')
    render(<Sidebar />)
    
    expect(screen.getByText('🦸 Herbst MUD')).toBeInTheDocument()
  })
  
  it('highlights active navigation item', async () => {
    const { Sidebar } = await import('./Sidebar')
    render(<Sidebar />)
    
    // The Dashboard link should have 'active' class since pathname is '/dashboard'
    const dashboardLink = screen.getByText('Dashboard').closest('a')
    expect(dashboardLink).toHaveClass('active')
  })
})