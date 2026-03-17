import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/map')({
  component: MapPage,
})

function MapPage() {
  return (
    <div className="page-header">
      <h1>🗺️ Map Builder</h1>
      <p>Visual map builder for the game world.</p>
      <div className="placeholder-page">
        <h2>Coming Soon</h2>
        <p>Visual map builder will be available here.</p>
        <p>Features include: Visual grid view, Room connections, NPC/Item placement, Zone management</p>
      </div>
    </div>
  )
}