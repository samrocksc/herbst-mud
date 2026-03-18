import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth')({
  component: AuthLayout,
})

function AuthLayout() {
  return (
    <div className="auth-layout">
      <header>
        <h1>Herbst MUD Admin</h1>
        <nav>
          <a href="/_auth/dashboard">Dashboard</a>
          <a href="/_auth/items">Items</a>
          <a href="/_auth/rooms">Rooms</a>
          <a href="/_auth/skills">Skills/Talents</a>
          <a href="/_auth/npcs">NPCs</a>
          <a href="/_auth/map">Map Builder</a>
          <a href="/">Logout</a>
        </nav>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  )
}