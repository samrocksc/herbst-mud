import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/players')({
  component: PlayersPage,
})

function PlayersPage() {
  return (
    <div className="page-header">
      <h1>👥 Players Management</h1>
      <p>View and manage player accounts and their characters.</p>
      <div className="placeholder-page">
        <h2>Coming Soon</h2>
        <p>Player management features will be available here.</p>
        <p>Features include: View players, Ban/Unban, Character stats, Activity logs</p>
      </div>
    </div>
  )
}