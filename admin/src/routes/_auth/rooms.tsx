import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/rooms')({
  component: RoomsPage,
})

function RoomsPage() {
  return (
    <div className="page-header">
      <h1>🏠 Rooms Management</h1>
      <p>Manage game rooms, areas, and their connections.</p>
      <div className="placeholder-page">
        <h2>Coming Soon</h2>
        <p>Room management features will be available here.</p>
        <p>Features include: Add rooms, Edit room details, Connect exits, Set room properties</p>
      </div>
    </div>
  )
}