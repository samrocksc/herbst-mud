import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth/items')({
  component: ItemsPage,
})

function ItemsPage() {
  return (
    <div className="page-header">
      <h1>🟣 Items Management</h1>
      <p>Manage weapons, armor, consumables, and quest items for the game world.</p>
      <div className="placeholder-page">
        <h2>Coming Soon</h2>
        <p>Item management features will be available here.</p>
        <p>Features include: Add items, Edit items, Delete items, Bulk operations</p>
      </div>
    </div>
  )
}