import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'

export const Route = createFileRoute('/_auth/dashboard')({
  component: Dashboard,
})

function Dashboard() {
  const [stats] = useState({
    onlinePlayers: 3,
    totalRooms: 42,
    totalItems: 156,
    totalNPCs: 28,
    uptime: '2d 14h 32m'
  })

  const [recentActivity] = useState([
    { id: 1, type: 'login', player: 'Player1', time: '2 min ago' },
    { id: 2, type: 'death', player: 'Player2', time: '15 min ago' },
    { id: 3, type: 'login', player: 'Player3', time: '32 min ago' },
  ])

  return (
    <div className="dashboard">
      <h2>Admin Dashboard</h2>
      
      <div className="dashboard-stats">
        <div className="stat-card">
          <div className="stat-value">{stats.onlinePlayers}</div>
          <div className="stat-label">Players Online</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.uptime}</div>
          <div className="stat-label">Server Uptime</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.totalRooms}</div>
          <div className="stat-label">Total Rooms</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.totalItems}</div>
          <div className="stat-label">Total Items</div>
        </div>
        <div className="stat-card">
          <div className="stat-value">{stats.totalNPCs}</div>
          <div className="stat-label">Total NPCs</div>
        </div>
      </div>

      <div className="dashboard-content">
        <div className="card">
          <h3>Quick Actions</h3>
          <div className="quick-actions">
            <a href="/_auth/items"><button>Manage Items</button></a>
            <a href="/_auth/rooms"><button>Manage Rooms</button></a>
            <a href="/_auth/skills"><button>Manage Skills</button></a>
            <a href="/_auth/npcs"><button>Manage NPCs</button></a>
            <a href="/_auth/map"><button>Map Builder</button></a>
          </div>
        </div>
        
        <div className="card">
          <h3>Recent Activity</h3>
          <ul className="activity-feed">
            {recentActivity.map(activity => (
              <li key={activity.id}>
                <span className={`activity-${activity.type}`}>{activity.type}</span>
                <span className="activity-player">{activity.player}</span>
                <span className="activity-time">{activity.time}</span>
              </li>
            ))}
          </ul>
        </div>
      </div>

      <div className="dashboard-content">
        <div className="card">
          <h3>System Health</h3>
          <div className="health-indicators">
            <div className="health-item">
              <span className="health-ok">●</span> Database: Connected
            </div>
            <div className="health-item">
              <span className="health-ok">●</span> Game Server: Running
            </div>
            <div className="health-item">
              <span className="health-ok">●</span> API: Healthy
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}