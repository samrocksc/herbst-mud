import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/dashboard')({
  component: Dashboard,
})

function Dashboard() {
  return (
    <div className="page-header">
      <h1>🟣 Dashboard</h1>
      <p>Overview of your MUD administration panel.</p>
      <div className="dashboard-content">
        <div className="card">
          <h3>👥 Users</h3>
          <p>Manage user accounts and permissions.</p>
          <span className="stat">--</span>
        </div>
        <div className="card">
          <h3>🎮 Characters</h3>
          <p>View and manage game characters.</p>
          <span className="stat">--</span>
        </div>
        <div className="card">
          <h3>🏠 Rooms</h3>
          <p>Manage game rooms and areas.</p>
          <span className="stat">--</span>
        </div>
        <div className="card">
          <h3>⚔️ NPCs</h3>
          <p>Manage NPCs and enemies.</p>
          <span className="stat">--</span>
        </div>
        <div className="card">
          <h3>📦 Items</h3>
          <p>Manage game items and loot.</p>
          <span className="stat">--</span>
        </div>
        <div className="card">
          <h3>⭐ Skills</h3>
          <p>Manage skills and talents.</p>
          <span className="stat">--</span>
        </div>
      </div>
    </div>
  )
}