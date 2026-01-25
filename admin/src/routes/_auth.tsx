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
          <a href="/">Logout</a>
        </nav>
      </header>
      <main>
        <Outlet />
      </main>
    </div>
  )
}