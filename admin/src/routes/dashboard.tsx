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
}

function Dashboard() {
  const navigate = useNavigate()
  const [stats, setStats] = useState<Stats>({ rooms: 0, npcs: 0, items: 0, players: 0 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      navigate({ to: '/login' })
      return
    }

    // Fetch stats
    const fetchStats = async () => {
      try {
        const [roomsRes, npcsRes] = await Promise.all([
          fetch('http://localhost:8080/rooms'),
          fetch('http://localhost:8080/npcs')
        ])

        const roomsData = await roomsRes.json()
        const npcsData = await npcsRes.json()

        setStats({
          rooms: Array.isArray(roomsData) ? roomsData.length : 0,
          npcs: npcsData.npcs?.length || 0,
          items: 0, // Items endpoint not available
          players: 0 // Players endpoint not available
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
    <div className="min-h-screen bg-[#0a0a0f] text-white p-8">
      <div className="max-w-[1200px] mx-auto">
        {/* Header */}
        <div className="flex justify-between items-center mb-8 border-b border-[#333] pb-4">
          <h1 className="m-0 text-[#61dafb]">Herbst MUD Admin</h1>
          <button
            onClick={handleLogout}
            className="px-4 py-2 bg-[#e74c3c] border-none rounded text-white cursor-pointer hover:bg-[#c0392b]"
          >
            Logout
          </button>
        </div>

        {/* Welcome */}
        <div className="bg-[#1a1a2e] rounded-lg p-6 mb-8">
          <h2 className="m-0 mb-2">Welcome back!</h2>
          <p className="m-0 text-[#888]">Manage your MUD world from this admin panel.</p>
        </div>

        {/* Quick Stats */}
        <div className="grid grid-cols-[repeat(auto-fit,minmax(200px,1fr))] gap-4 mb-8">
          <div className="bg-[#1a1a2e] rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-[#27ae60]">
              {loading ? '--' : stats.rooms}
            </div>
            <div className="text-[#888] text-sm">Total Rooms</div>
          </div>
          <div className="bg-[#1a1a2e] rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-[#f39c12]">
              {loading ? '--' : stats.npcs}
            </div>
            <div className="text-[#888] text-sm">Active NPCs</div>
          </div>
          <div className="bg-[#1a1a2e] rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-[#3498db]">
              {loading ? '--' : stats.items}
            </div>
            <div className="text-[#888] text-sm">Items</div>
          </div>
          <div className="bg-[#1a1a2e] rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-[#9b59b6]">
              {loading ? '--' : stats.players}
            </div>
            <div className="text-[#888] text-sm">Players</div>
          </div>
        </div>

        {/* Navigation Cards */}
        <h3 className="mb-4">Admin Tools</h3>
        <div className="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4">
          <Link
            to="/map"
            className="block bg-[#1a1a2e] rounded-lg p-6 no-underline text-white border border-[#333] transition-colors hover:border-[#61dafb]"
          >
            <div className="text-2xl mb-2">🗺️</div>
            <div className="font-bold mb-1">Map Builder</div>
            <div className="text-[#888] text-sm">View and edit room layout, connections, and z-levels</div>
          </Link>

          <Link
            to="/npcs"
            className="block bg-[#1a1a2e] rounded-lg p-6 no-underline text-white border border-[#333] transition-colors hover:border-[#61dafb]"
          >
            <div className="text-2xl mb-2">👤</div>
            <div className="font-bold mb-1">NPC Manager</div>
            <div className="text-[#888] text-sm">Create, edit, and manage NPCs and their locations</div>
          </Link>

          <Link
            to="/items"
            className="block bg-[#1a1a2e] rounded-lg p-6 no-underline text-white border border-[#333] transition-colors hover:border-[#61dafb]"
          >
            <div className="text-2xl mb-2">📦</div>
            <div className="font-bold mb-1">Item Manager</div>
            <div className="text-[#888] text-sm">Create, edit, and manage items and equipment</div>
          </Link>

          <div className="bg-[#1a1a2e] rounded-lg p-6 border border-[#333] opacity-50">
            <div className="text-2xl mb-2">🎮</div>
            <div className="font-bold mb-1">Player Manager</div>
            <div className="text-[#666] text-sm">Coming soon</div>
          </div>
        </div>
      </div>
    </div>
  )
}