import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useEffect, useState } from 'react'

export const Route = createFileRoute('/dashboard')({
  component: Dashboard,
})

interface Stats {
  rooms: number
  npcs: number
  items: number
  players: number
  skills: number
}

function Dashboard() {
  const navigate = useNavigate()
  const [stats, setStats] = useState<Stats>({ rooms: 0, npcs: 0, items: 0, players: 0, skills: 0 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate({ to: '/login' })
      return
    }

    const fetchStats = async () => {
      try {
        const [roomsRes, npcsRes, skillsRes] = await Promise.all([
          fetch('http://localhost:8080/rooms'),
          fetch('http://localhost:8080/npcs'),
          fetch('http://localhost:8080/skills')
        ])
        const roomsData = await roomsRes.json()
        const npcsData = await npcsRes.json()
        const skillsData = await skillsRes.json()
        setStats({
          rooms: Array.isArray(roomsData) ? roomsData.length : 0,
          npcs: npcsData.npcs?.length || 0,
          items: 0,
          players: 0,
          skills: skillsData.count || 0
        })
      } catch (err) {
        console.error('Failed to fetch stats:', err)
      } finally {
        setLoading(false)
      }
    }
    fetchStats()
  }, [navigate])

  const handleLogout = () => {
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
    localStorage.removeItem('email')
    localStorage.removeItem('isAdmin')
    navigate({ to: '/login' })
  }

  return (
    <div className="min-h-screen bg-surface text-text p-8">
      <div className="max-w-[1200px] mx-auto">
        <div className="flex justify-between items-center mb-8 border-b border-border pb-4">
          <h1 className="m-0 text-primary">Herbst MUD Admin</h1>
          <button onClick={handleLogout} className="px-4 py-2 bg-danger border-none rounded text-white cursor-pointer hover:bg-danger-hover">
            Logout
          </button>
        </div>

        <div className="bg-surface-muted rounded-lg p-6 mb-8">
          <h2 className="m-0 mb-2 text-text">Welcome back!</h2>
          <p className="m-0 text-text-muted">Manage your MUD world from this admin panel.</p>
        </div>

        <div className="grid grid-cols-[repeat(auto-fit,minmax(200px,1fr))] gap-4 mb-8">
          <div className="bg-surface-muted rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-primary">{loading ? '--' : stats.rooms}</div>
            <div className="text-text-muted text-sm">Total Rooms</div>
          </div>
          <div className="bg-surface-muted rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-warning">{loading ? '--' : stats.npcs}</div>
            <div className="text-text-muted text-sm">Active NPCs</div>
          </div>
          <div className="bg-surface-muted rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-accent">{loading ? '--' : stats.items}</div>
            <div className="text-text-muted text-sm">Items</div>
          </div>
          <div className="bg-surface-muted rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-primary-hover">{loading ? '--' : stats.skills}</div>
            <div className="text-text-muted text-sm">Skills</div>
          </div>

          <div className="bg-surface-muted rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-secondary">{loading ? '--' : stats.players}</div>
            <div className="text-text-muted text-sm">Players</div>
          </div>
        </div>

        <h3 className="mb-4 text-text">Admin Tools</h3>
        <div className="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4">
          <Link to="/map" className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary">
            <div className="text-2xl mb-2">🗺️</div>
            <div className="font-bold mb-1">Map Builder</div>
            <div className="text-text-muted text-sm">View and edit room layout, connections, and z-levels</div>
          </Link>

          <Link to="/npcs" className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary">
            <div className="text-2xl mb-2">👤</div>
            <div className="font-bold mb-1">NPC Manager</div>
            <div className="text-text-muted text-sm">Create, edit, and manage NPCs and their locations</div>
          </Link>

          <Link to="/items" className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary">
            <div className="text-2xl mb-2">📦</div>
            <div className="font-bold mb-1">Item Manager</div>
            <div className="text-text-muted text-sm">Create, edit, and manage items and equipment</div>
          </Link>

          <Link to="/export" className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary">
            <div className="text-2xl mb-2">💾</div>
            <div className="font-bold mb-1">Export / Import</div>
            <div className="text-text-muted text-sm">Backup and restore game world data</div>
          </Link>

          <Link to="/players" className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary">
            <div className="text-2xl mb-2">🎮</div>
            <div className="font-bold mb-1">Player Manager</div>
            <div className="text-text-muted text-sm">Manage players and reset passwords</div>
          </Link>

          <Link to="/skills" className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary">
            <div className="text-2xl mb-2">⚡</div>
            <div className="font-bold mb-1">Skills Manager</div>
            <div className="text-text-muted text-sm">Create, edit, and manage skills</div>
          </Link>

          <Link to="/talents" className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary">
            <div className="text-2xl mb-2">🎯</div>
            <div className="font-bold mb-1">Talents Manager</div>
            <div className="text-text-muted text-sm">Manage talent specializations</div>
          </Link>
        </div>
      </div>
    </div>
  )
}