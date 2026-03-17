import { createFileRoute, Outlet } from '@tanstack/react-router'
import { Sidebar } from '../components/Sidebar'

export const Route = createFileRoute('/_auth')({
  component: AuthLayout,
})

function AuthLayout() {
  return (
    <div className="auth-layout">
      <Sidebar />
      <div className="auth-content">
        <main>
          <Outlet />
        </main>
      </div>
    </div>
  )
}