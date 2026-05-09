import { createFileRoute, Link, Outlet, useLocation } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { useEquipmentTemplates, useCreateTemplate, useDeleteTemplate } from '../../hooks/useEquipmentTemplates'
import { useItemInstances } from '../../hooks/useItemInstances'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { DeleteConfirmation } from '../../components/DeleteConfirmation'
import { showToast } from '../../components/Toast'
import { SLOT_OPTIONS, ITEM_TYPE_OPTIONS } from '../../components/itemConstants'
import type { EquipmentTemplate } from '../../hooks/useEquipmentTemplates'

export const Route = createFileRoute('/_auth/items')({
  component: ItemsIndex,
})

function ItemsIndex() {
  const [searchQuery, setSearchQuery] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [deleteId, setDeleteId] = useState<string | null>(null)

  const templatesQuery = useEquipmentTemplates()
  const instancesQuery = useItemInstances()
  const deleteMutation = useDeleteTemplate()

  const instanceCounts = useMemo(() => {
    const counts: Record<string, number> = {}
    for (const inst of instancesQuery.data ?? []) {
      const tid = inst.equipment_template_id
      if (tid) counts[tid] = (counts[tid] ?? 0) + 1
    }
    return counts
  }, [instancesQuery.data])

  const filteredItems = (templatesQuery.data ?? []).filter((item) =>
    item.name.toLowerCase().includes(searchQuery.toLowerCase()),
  )

  const handleDelete = (id: string) => {
    deleteMutation.mutate(id, {
      onSuccess: () => { setDeleteId(null); showToast('Item template deleted', 'success') },
    })
  }

  const columns: Column<EquipmentTemplate>[] = [
    {
      header: 'Name',
      accessor: 'name',
      render: (_, row) => (
        <Link
          to="/items/$itemId"
          params={{ itemId: row.id }}
          className="no-underline text-primary hover:underline font-bold"
        >
          {row.name}
        </Link>
      ),
    },
    { header: 'Slot', accessor: 'slot' },
    { header: 'Level', accessor: 'level', align: 'center' },
    { header: 'Type', accessor: 'item_type' },
    { header: 'Weight', accessor: 'weight', align: 'center' },
    {
      header: 'Instances',
      accessor: 'instances',
      align: 'center',
      render: (_, row) => (
        <span className="badge badge-neutral">{instanceCounts[row.id] ?? 0}</span>
      ),
    },
    {
      header: '',
      accessor: '_actions',
      align: 'right',
      render: (_, row) => (
        <div className="flex gap-2 justify-end">
          <Button variant="ghost" size="sm" onClick={(e) => { e.stopPropagation(); setDeleteId(row.id) }}>
            Delete
          </Button>
        </div>
      ),
    },
  ]

  const location = useLocation()
  const isList = location.pathname === '/items'

  if (!isList) return <Outlet />

  return (
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader title="Items" showBack backTo="/dashboard" actions={
        <Button variant="primary" onClick={() => setShowForm(true)}>+ Add Item</Button>
      } />

      <div className="mb-4">
        <input
          type="text"
          placeholder="Search items by name..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full max-w-sm p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>

      {templatesQuery.isLoading && <div className="p-8 text-text-muted text-center text-xs">Loading items...</div>}
      {templatesQuery.isError && (
        <div className="p-4 bg-danger/10 border border-danger rounded text-danger text-xs">
          Failed to load items: {templatesQuery.error?.message ?? 'Unknown error'}
        </div>
      )}
      {templatesQuery.isSuccess && (
        <DataTable<EquipmentTemplate>
          columns={columns}
          data={filteredItems}
          getKey={(row) => row.id}
          emptyMessage="No items found."
          variant="dark"
        />
      )}

      {showForm && (
        <CreateItemModal onClose={() => setShowForm(false)} />
      )}

      {deleteId && (
        <DeleteConfirmation
          open={!!deleteId}
          title="Delete Item Template"
          message="Are you sure? This will permanently delete this item template. Instances based on this template will not be deleted."
          onConfirm={() => handleDelete(deleteId)}
          onCancel={() => setDeleteId(null)}
          isLoading={deleteMutation.isPending}
        />
      )}
    </div>
  )
}

function CreateItemModal({ onClose }: { onClose: () => void }) {
  const { mutate: createTemplate, isPending } = useCreateTemplate()

  const [form, setForm] = useState({
    name: '',
    description: '',
    slot: '',
    item_type: 'misc',
    level: 1,
    weight: 0,
    color: '',
  })

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (!form.name.trim()) return
    createTemplate(form, {
      onSuccess: () => {
        showToast('Item template created', 'success')
        onClose()
      },
    })
  }

  return (
    <div className="modal-overlay" onClick={onClose}>
      <div className="modal-card" onClick={(e) => e.stopPropagation()}>
        <h3 className="mt-0">Create Item Template</h3>
        <form onSubmit={handleSubmit} className="space-y-3">
          <div>
            <label className="block text-sm text-text-muted mb-1">Name *</label>
            <input
              type="text"
              value={form.name}
              onChange={(e) => setForm({ ...form, name: e.target.value })}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              required
            />
          </div>
          <div>
            <label className="block text-sm text-text-muted mb-1">Description</label>
            <input
              type="text"
              value={form.description}
              onChange={(e) => setForm({ ...form, description: e.target.value })}
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>
          <div className="flex gap-3">
            <div className="flex-1">
              <label className="block text-sm text-text-muted mb-1">Slot</label>
              <select
                value={form.slot}
                onChange={(e) => setForm({ ...form, slot: e.target.value })}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              >
                {SLOT_OPTIONS.map((opt) => (
                  <option key={opt.value} value={opt.value}>{opt.label}</option>
                ))}
              </select>
            </div>
            <div className="flex-1">
              <label className="block text-sm text-text-muted mb-1">Type</label>
              <select
                value={form.item_type}
                onChange={(e) => setForm({ ...form, item_type: e.target.value })}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              >
                {ITEM_TYPE_OPTIONS.map((opt) => (
                  <option key={opt.value} value={opt.value}>{opt.label}</option>
                ))}
              </select>
            </div>
          </div>
          <div className="flex gap-3">
            <div className="flex-1">
              <label className="block text-sm text-text-muted mb-1">Level</label>
              <input
                type="number"
                value={form.level}
                onChange={(e) => setForm({ ...form, level: parseInt(e.target.value) || 1 })}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div className="flex-1">
              <label className="block text-sm text-text-muted mb-1">Weight</label>
              <input
                type="number"
                value={form.weight}
                onChange={(e) => setForm({ ...form, weight: parseInt(e.target.value) || 0 })}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>
          <div className="flex gap-2 justify-end mt-4">
            <Button variant="secondary" onClick={onClose}>Cancel</Button>
            <Button variant="primary" type="submit" disabled={isPending || !form.name.trim()}>
              {isPending ? 'Creating…' : 'Create'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  )
}