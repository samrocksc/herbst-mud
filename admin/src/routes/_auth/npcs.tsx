import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/npcs')({
  component: NPCsPage,
})

function NPCsPage() {
  return (
    <div className="page-header">
      <h1>👾 NPCs Management</h1>
      <p>Manage non-player characters, enemies, and friendly NPCs.</p>
      <div className="placeholder-page">
        <h2>Coming Soon</h2>
        <p>NPC management features will be available here.</p>
        <p>Features include: Create NPCs, Set stats, Define behaviors, Assign locations</p>
      </div>
    </div>
  )
}