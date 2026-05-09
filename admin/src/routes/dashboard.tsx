import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { StatCard } from '../components/StatCard'
import { StatGrid } from '../components/StatGrid'
import { PageHeader } from '../components/PageHeader'
import { Button } from '../components/Button'
import { showToast } from '../components/Toast'
import { apiGet } from '../utils/apiFetch'
import { ToolGrid } from './ToolGrid'

export const Route = createFileRoute('/dashboard')({
  component: Dashboard,
})

type Stats = { rooms: number; npcs: number; items: number; instances: number; players: number; skills: number }

function Dashboard() {
  const navigate = useNavigate()
  const [stats, setStats] = useState<Stats>({ rooms: 0, npcs: 0, items: 0, instances: 0, players: 0, skills: 0 })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    const token = localStorage.getItem('token')
    if (!token) { navigate({ to: '/login' }); return }

    const fetchStats = async () => {
      try {
        const [rooms, npcs, abilities, templates, instances, characters] = await Promise.all([
          apiGet<unknown[]>(`${window.location.origin}/rooms`),
          apiGet<unknown[]>(`${window.location.origin}/npcs`),
          apiGet<unknown[]>(`${window.location.origin}/api/abilities`),
          apiGet<unknown[]>(`${window.location.origin}/api/equipment-templates`),
          apiGet<unknown[]>(`${window.location.origin}/api/item-instances`),
          apiGet<unknown[]>(`${window.location.origin}/characters`),
        ])
        setStats({
          rooms: Array.isArray(rooms) ? rooms.length : 0,
          npcs: Array.isArray(npcs) ? npcs.length : 0,
          items: Array.isArray(templates) ? templates.length : 0,
          instances: Array.isArray(instances) ? instances.length : 0,
          players: Array.isArray(characters) ? characters.length : 0,
          skills: Array.isArray(abilities) ? abilities.length : 0,
        })
      } catch (err) {
        showToast(err instanceof Error ? err.message : 'Failed to load stats', 'error')
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
        <PageHeader title="Herbst MUD Admin" actions={<Button onClick={handleLogout} variant="danger">Logout</Button>} />
        <div className="bg-surface-muted rounded-lg p-6 mb-8">
          <h2 className="m-0 mb-2 text-text">Welcome back!</h2>
          <p className="m-0 text-text-muted">Manage your MUD world from this admin panel.</p>
        </div>
        <StatGrid>
          <StatCard label="Total Rooms" value={stats.rooms} accent="primary" loading={loading} />
          <StatCard label="Active NPCs" value={stats.npcs} accent="warning" loading={loading} />
          <StatCard label="Items" value={stats.items} accent="accent" loading={loading} />
          <StatCard label="Instances" value={stats.instances} accent="primary" loading={loading} />
          <StatCard label="Players" value={stats.players} accent="secondary" loading={loading} />
          <StatCard label="Skills" value={stats.skills} accent="success" loading={loading} />
        </StatGrid>
        <h3 className="mb-4 text-text">Admin Tools</h3>
        <ToolGrid />
      </div>
    </div>
  )
}