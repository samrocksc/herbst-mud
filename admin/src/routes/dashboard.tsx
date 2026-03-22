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
          items: 0,
          players: 0
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
    <div className="min-h-screen bg-[#1a1612] text-[#e8dcc4] p-8">
      <div className="max-w-[1200px] mx-auto">
        <div className="flex justify-between items-center mb-8 border-b border-[#5a4a35] pb-4">
          <h1 className="m-0 text-[#4a7c4e]">Herbst MUD Admin</h1>
          <button onClick={handleLogout} className="px-4 py-2 bg-[#8b4444] border-none rounded text-[#e8dcc4] cursor-pointer hover:bg-[#a84444]">
            Logout
          </button>
        </div>

        <div className="bg-[#2d2416] rounded-lg p-6 mb-8">
          <h2 className="m-0 mb-2 text-[#e8dcc4]">Welcome back!</h2>
          <p className="m-0 text-[#a89070]">Manage your MUD world from this admin panel.</p>
        </div>

        <div className="grid grid-cols-[repeat(auto-fit,minmax(200px,1fr))] gap-4 mb-8">
          <div className="bg-[#2d2416] rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-[#4a7c4e]">{loading ? '--' : stats.rooms}</div>
            <div className="text-[#a89070] text-sm">Total Rooms</div>
          </div>
          <div className="bg-[#2d2416] rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-[#a87044]">{loading ? '--' : stats.npcs}</div>
            <div className="text-[#a89070] text-sm">Active NPCs</div>
          </div>
          <div className="bg-[#2d2416] rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-[#8b7355]">{loading ? '--' : stats.items}</div>
            <div className="text-[#a89070] text-sm">Items</div>
          </div>
          <div className="bg-[#2d2416] rounded-lg p-6 text-center">
            <div className="text-2xl font-bold text-[#5a9c5e]">{loading ? '--' : stats.players}</div>
            <div className="text-[#a89070] text-sm">Players</div>
          </div>
        </div>

        <h3 className="mb-4 text-[#e8dcc4]">Admin Tools</h3>
        <div className="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4">
          <Link to="/map" className="block bg-[#2d2416] rounded-lg p-6 no-underline text-[#e8dcc4] border border-[#5a4a35] transition-colors hover:border-[#4a7c4e]">
            <div className="text-2xl mb-2">🗺️</div>
            <div className="font-bold mb-1">Map Builder</div>
            <div className="text-[#a89070] text-sm">View and edit room layout, connections, and z-levels</div>
          </Link>

          <Link to="/npcs" className="block bg-[#2d2416] rounded-lg p-6 no-underline text-[#e8dcc4] border border-[#5a4a35] transition-colors hover:border-[#4a7c4e]">
            <div className="text-2xl mb-2">👤</div>
            <div className="font-bold mb-1">NPC Manager</div>
            <div className="text-[#a89070] text-sm">Create, edit, and manage NPCs and their locations</div>
          </Link>

          <Link to="/items" className="block bg-[#2d2416] rounded-lg p-6 no-underline text-[#e8dcc4] border border-[#5a4a35] transition-colors hover:border-[#4a7c4e]">
            <div className="text-2xl mb-2">📦</div>
            <div className="font-bold mb-1">Item Manager</div>
            <div className="text-[#a89070] text-sm">Create, edit, and manage items and equipment</div>
          </Link>

          <div className="bg-[#2d2416] rounded-lg p-6 border border-[#5a4a35] opacity-50">
            <div className="text-2xl mb-2">🎮</div>
            <div className="font-bold mb-1">Player Manager</div>
            <div className="text-[#a89070] text-sm">Coming soon</div>
          </div>
        </div>
      </div>
    </div>
  )
}