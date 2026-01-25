import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/dashboard')({
  component: Dashboard,
})

function Dashboard() {
  return (
    <div className="dashboard">
      <h2>Admin Dashboard</h2>
      <div className="dashboard-content">
        <div className="card">
          <h3>Users</h3>
          <p>Manage user accounts</p>
          <button>View Users</button>
        </div>
        <div className="card">
          <h3>Characters</h3>
          <p>Manage game characters</p>
          <button>View Characters</button>
        </div>
        <div className="card">
          <h3>Rooms</h3>
          <p>Manage game rooms</p>
          <button>View Rooms</button>
        </div>
      </div>
    </div>
  )
}