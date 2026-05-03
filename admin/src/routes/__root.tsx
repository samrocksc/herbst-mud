import { createRootRoute, Outlet } from '@tanstack/react-router'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { useState } from 'react'
import { Sidebar } from '../components/Sidebar'
import { MenuIcon } from '../components/icons/MenuIcon'
import { Button } from '../components/Button'

const queryClient = new QueryClient()

export const Route = createRootRoute({
  component: RootComponent,
})

function RootComponent() {
  const [mobileSidebarOpen, setMobileSidebarOpen] = useState(false)

  return (
    <QueryClientProvider client={queryClient}>
      <div className="flex h-screen bg-white">
        {/* Mobile menu button — only visible on small screens */}
        <Button
          variant="ghost"
          size="sm"
          onClick={() => setMobileSidebarOpen(true)}
          aria-label="Open menu"
          className="fixed top-3 left-3 z-50 p-2 bg-surface border border-border text-text-muted hover:bg-surface-muted hover:text-text lg:hidden"
        >
          <MenuIcon stroke="currentColor" />
        </Button>

        {/* Sidebar — hidden on mobile by default, lg: always visible */}
        <div
          className={[
            'lg:block',
            mobileSidebarOpen ? 'block' : 'hidden',
          ].join(' ')}
        >
          <Sidebar />
        </div>

        {/* Mobile backdrop — closes sidebar when tapping outside */}
        {mobileSidebarOpen && (
          <div
            className="fixed inset-0 bg-black/30 z-30 lg:hidden"
            onClick={() => setMobileSidebarOpen(false)}
          />
        )}

        <main className="flex-1 overflow-auto bg-gray-50">
          <Outlet />
        </main>
      </div>
    </QueryClientProvider>
  )
}