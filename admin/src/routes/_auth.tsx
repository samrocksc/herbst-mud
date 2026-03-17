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
          <a href="dashboard">Dashboard</a>
          <a href="items">Items</a>
          <a href="rooms">Rooms</a>
          <a href="skills">Skills/Talents</a>
          <a href="npcs">NPCs</a>
          <a href="map">Map Builder</a>
          <a href="/">Logout</a>
        </nav>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  )
}