import { createFileRoute, Link, Outlet, useLocation } from '@tanstack/react-router'
import { useState, useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { apiGet } from '../../utils/apiFetch'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'

export const Route = createFileRoute('/_auth/items')({
  component: ItemsIndex,
})

// ─── Types ──────────────────────────────────────────────────────────────────

type ItemTemplate = Readonly<{
  id: string
  name: string
  description: string
  slot: string
  level: number
  weight: number
  item_type: string
  color: string
  is_visible: boolean
  is_immovable: boolean
  effect_type: string
  effect_value: number
  effect_duration: number
  is_container: boolean
  container_capacity: number
  is_locked: boolean
}>

type ItemInstance = Readonly<{
  id: number
  equipment_template_id: string
  [key: string]: unknown
}>

const API = `${window.location.origin}`

// ─── Component ──────────────────────────────────────────────────────────────

function ItemsIndex() {
  const [searchQuery, setSearchQuery] = useState('')

  const templatesQuery = useQuery({
    queryKey: ['item-templates'],
    queryFn: () => apiGet<ItemTemplate[]>(`${API}/api/equipment-templates`),
  })

  const instancesQuery = useQuery({
    queryKey: ['item-instances'],
    queryFn: async () => {
      const data = await apiGet<ItemInstance[]>(`${API}/api/item-instances`)
      return Array.isArray(data) ? data : []
    },
  })

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

  const columns: Column<ItemTemplate>[] = [
    {
      header: 'ID',
      accessor: 'id',
      render: (_, row) => <span className="font-mono text-xs">{row.id}</span>,
    },
    {
      header: 'Name',
      accessor: 'name',
      className: 'font-bold',
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
  ]

  const location = useLocation()
  const isList = location.pathname === '/items'

  if (!isList) {
    return <Outlet />
  }

  return (
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader title="Items" showBack backTo="/dashboard" />

      {/* Search bar */}
      <div className="mb-4">
        <input
          type="text"
          placeholder="Search items by name..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-full max-w-sm p-2 bg-surface border border-border rounded text-text text-sm"
        />
      </div>

      {/* Loading */}
      {templatesQuery.isLoading && (
        <div className="p-8 text-text-muted text-center text-xs">Loading items...</div>
      )}

      {/* Error */}
      {templatesQuery.isError && (
        <div className="p-4 bg-danger/10 border border-danger rounded text-danger text-xs">
          Failed to load items: {templatesQuery.error?.message ?? 'Unknown error'}
        </div>
      )}

      {/* Data table */}
      {templatesQuery.isSuccess && (
        <DataTable<ItemTemplate>
          columns={columns}
          data={filteredItems}
          getKey={(row) => row.id}
          emptyMessage="No items found."
          variant="dark"
        />
      )}
    </div>
  )
}