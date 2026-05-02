import { createFileRoute, useNavigate, Link } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { StatCard } from '../components/StatCard'
import { StatGrid } from '../components/StatGrid'
import { PageHeader } from '../components/PageHeader'

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
        const [roomsRes, npcsRes, skillsRes, equipmentRes, charactersRes] = await Promise.all([
          fetch(`${window.location.origin}/rooms`),
          fetch(`${window.location.origin}/npcs`),
          fetch(`${window.location.origin}/skills`, { headers: { Authorization: `Bearer ${localStorage.getItem('token')}` } }),
          fetch(`${window.location.origin}/equipment`),
          fetch(`${window.location.origin}/characters`),
        ])
        const roomsData = await roomsRes.json()
        const npcsData = await npcsRes.json()
        const skillsData = await skillsRes.json()
        const equipmentData = await equipmentRes.json()
        const charactersData = await charactersRes.json()
        setStats({
          rooms: Array.isArray(roomsData) ? roomsData.length : 0,
          npcs: npcsData.npcs?.length || 0,
          items: Array.isArray(equipmentData) ? equipmentData.length : 0,
          players: Array.isArray(charactersData) ? charactersData.length : 0,
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
        <PageHeader
          title="Herbst MUD Admin"
          actions={
            <button onClick={handleLogout} className="bg-danger text-text-inverse border-border-dark hover:bg-danger-hover">
              Logout
            </button>
          }
        />

        <div className="bg-surface-muted rounded-lg p-6 mb-8">
          <h2 className="m-0 mb-2 text-text">Welcome back!</h2>
          <p className="m-0 text-text-muted">Manage your MUD world from this admin panel.</p>
        </div>

        <StatGrid>
          <StatCard label="Total Rooms" value={stats.rooms} accent="primary" loading={loading} />
          <StatCard label="Active NPCs" value={stats.npcs} accent="warning" loading={loading} />
          <StatCard label="Items" value={stats.items} accent="accent" loading={loading} />
          <StatCard label="Players" value={stats.players} accent="secondary" loading={loading} />
          <StatCard label="Skills" value={stats.skills} accent="success" loading={loading} />
        </StatGrid>

        <h3 className="mb-4 text-text">Admin Tools</h3>
        <div className="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4">
          <ToolCard to="/map" emoji="🗺️" title="Map Builder" desc="View and edit room layout, connections, and z-levels" />
          <ToolCard to="/npcs" emoji="👤" title="NPC Manager" desc="Create, edit, and manage NPCs and their locations" />
          <ToolCard to="/items" emoji="📦" title="Item Manager" desc="Create, edit, and manage items and equipment" />
          <ToolCard to="/export" emoji="💾" title="Export / Import" desc="Backup and restore game world data" />
          <ToolCard to="/players" emoji="🎮" title="Player Manager" desc="Manage players and reset passwords" />
          <ToolCard to="/abilities" emoji="⚡" title="Abilities Manager" desc="Create, edit, and manage abilities" />
          <ToolCard to="/weapon-skills" emoji="🎯" title="Weapon Skills Manager" desc="Manage weapon skill specializations" />
          <ToolCard to="/factions" emoji="⚔️" title="Factions Manager" desc="Manage factions, categories, and member standing" />
        </div>
      </div>
    </div>
  )
}

function ToolCard({ to, emoji, title, desc }: { to: string; emoji: string; title: string; desc: string }) {
  return (
    <Link
      to={to as any}
      className="block bg-surface-muted rounded-lg p-6 no-underline text-text border border-border transition-colors hover:border-primary hover:bg-surface-muted/70"
    >
      <div className="text-2xl mb-2">{emoji}</div>
      <div className="font-bold mb-1">{title}</div>
      <div className="text-text-muted text-sm">{desc}</div>
    </Link>
  )
}
