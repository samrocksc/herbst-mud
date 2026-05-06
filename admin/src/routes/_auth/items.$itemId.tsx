import { useState } from 'react'
import { createFileRoute, Link } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost } from '../../utils/apiFetch'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Modal } from '../../components/Modal'
import { Button } from '../../components/Button'

export const Route = createFileRoute('/_auth/items/$itemId')({
  component: ItemDetail,
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
  stats: Record<string, unknown>
  color: string
  is_visible: boolean
  is_immovable: boolean
  effect_type: string
  effect_value: number
  effect_duration: number
  is_container: boolean
  container_capacity: number
  is_locked: boolean
  key_item_id: string
  reveal_condition: string
  expires_at: string | null
}>

type ItemInstance = Readonly<{
  id: number
  name: string
  ownerId: number | null
  roomId: number
  equipment_template_id: string
}>

// ─── Columns ────────────────────────────────────────────────────────────────

const instanceColumns: Column<ItemInstance>[] = [
  {
    header: 'ID',
    accessor: 'id',
    render: (_, row) => (
      <Link
        to="/items/$itemId/instances/$instanceId"
        params={{ itemId: row.equipment_template_id, instanceId: String(row.id) }}
        className="text-primary no-underline hover:underline font-mono text-xs"
      >
        {row.id}
      </Link>
    ),
  },
  { header: 'Name', accessor: 'name' },
  {
    header: 'Location',
    accessor: 'ownerId',
    render: (_, row) => {
      if (row.ownerId != null) return (
        <Link
          to="/npcs"
          className="text-primary text-xs no-underline hover:underline"
        >
          Held by Character #{row.ownerId}
        </Link>
      )
      if (row.roomId > 0) return (
        <Link
          to="/map"
          className="text-primary text-xs no-underline hover:underline"
        >
          In Room #{row.roomId}
        </Link>
      )
      return <span className="text-text-muted text-xs">Nowhere</span>
    },
  },
]

// ─── Component ──────────────────────────────────────────────────────────────

function ItemDetail() {
  const { itemId } = Route.useParams()
  const queryClient = useQueryClient()
  const [modalOpen, setModalOpen] = useState(false)
  const [spawnRoomId, setSpawnRoomId] = useState('')

  const templateQuery = useQuery({
    queryKey: ['item-template', itemId],
    queryFn: () => apiGet<ItemTemplate>('/api/equipment-templates/' + itemId),
  })

  const instancesQuery = useQuery({
    queryKey: ['item-instances', 'template', itemId],
    queryFn: () => apiGet<ItemInstance[]>('/api/item-instances?templateId=' + itemId),
  })

  const spawnMutation = useMutation({
    mutationFn: (roomId: number) =>
      apiPost('/api/item-instances', { equipment_template_id: itemId, room_id: roomId }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances', 'template', itemId] })
      setModalOpen(false)
      setSpawnRoomId('')
    },
  })

  const template = templateQuery.data
  const instances = instancesQuery.data ?? []

  if (templateQuery.isLoading || instancesQuery.isLoading) {
    return (
      <div className="p-8">
        <PageHeader title="Loading..." backTo="/items" />
        <div className="text-text-muted text-sm">Loading item details...</div>
      </div>
    )
  }

  if (templateQuery.error || !template) {
    return (
      <div className="p-8">
        <PageHeader title="Error" backTo="/items" />
        <div className="text-danger">Failed to load item: {(templateQuery.error as Error)?.message ?? 'Unknown error'}</div>
      </div>
    )
  }

  return (
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader title={template.name} backTo="/items" />

      {/* Template Stats */}
      <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
        <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Item Stats</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 gap-4">
          <DetailField label="ID" value={template.id} />
          <DetailField label="Name" value={template.name} />
          <DetailField label="Slot" value={template.slot} />
          <DetailField label="Level" value={String(template.level)} />
          <DetailField label="Weight" value={String(template.weight)} />
          <DetailField label="Type" value={template.item_type} />
          <DetailField label="Color" value={template.color || '—'} />
          <DetailField label="Visible" value={template.is_visible ? 'Yes' : 'No'} />
          <DetailField label="Immovable" value={template.is_immovable ? 'Yes' : 'No'} />
          <DetailField label="Container" value={template.is_container ? `Yes (${template.container_capacity})` : 'No'} />
          <DetailField label="Locked" value={template.is_locked ? 'Yes' : 'No'} />
          <DetailField label="Key Item" value={template.key_item_id || '—'} />
          {template.effect_type && (
            <>
              <DetailField label="Effect" value={template.effect_type} />
              <DetailField label="Effect Value" value={String(template.effect_value)} />
              <DetailField label="Duration" value={String(template.effect_duration)} />
            </>
          )}
          {template.expires_at && (
            <DetailField label="Expires" value={new Date(template.expires_at).toLocaleDateString()} />
          )}
        </div>
        <p className="text-text text-sm mt-4">{template.description || 'No description.'}</p>
        {Object.keys(template.stats ?? {}).length > 0 && (
          <div className="mt-4 pt-4 border-t border-border">
            <span className="text-text-muted text-xs block mb-2">Stats</span>
            <div className="grid grid-cols-3 gap-x-4 gap-y-1">
              {Object.entries(template.stats as Record<string, number>).map(([stat, val]) => (
                <div key={stat} className="flex justify-between text-sm">
                  <span className="text-text-muted">{stat}</span>
                  <span className="text-text font-medium">{val}</span>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Instances */}
      <div className="bg-surface-muted rounded-lg p-6 border border-border">
        <div className="flex items-center justify-between mb-4">
          <h2 className="m-0 text-text text-lg font-semibold">Instances ({instances.length})</h2>
          <Button variant="secondary" size="sm" onClick={() => setModalOpen(true)}>
            + Add Instance
          </Button>
        </div>
        {instancesQuery.isError && (
          <div className="text-danger text-xs mb-3">Failed to load instances</div>
        )}
        <DataTable<ItemInstance>
          columns={instanceColumns}
          data={instances}
          getKey={(row) => row.id}
          emptyMessage="No instances of this item found."
          variant="dark"
        />
      </div>

      {/* Add Instance Modal */}
      <Modal isOpen={modalOpen} onClose={() => { setModalOpen(false); setSpawnRoomId('') }} title="Add Instance">
        <div className="space-y-4">
          <div>
            <label className="block text-text text-sm font-medium mb-1">Room ID</label>
            <input
              type="number"
              className="w-full bg-surface border border-border rounded-md px-3 py-2 text-text text-sm focus:outline-none focus:ring-2 focus:ring-primary"
              placeholder="Enter room ID to spawn instance in"
              value={spawnRoomId}
              onChange={(e) => setSpawnRoomId(e.target.value)}
            />
          </div>
          {spawnMutation.isError && (
            <div className="text-danger text-xs">{(spawnMutation.error as Error)?.message ?? 'Failed to spawn instance'}</div>
          )}
          <div className="flex justify-end gap-2">
            <Button variant="secondary" size="sm" onClick={() => { setModalOpen(false); setSpawnRoomId('') }}>
              Cancel
            </Button>
            <Button
              variant="primary"
              size="sm"
              disabled={!spawnRoomId || spawnMutation.isPending}
              onClick={() => spawnMutation.mutate(Number(spawnRoomId))}
            >
              {spawnMutation.isPending ? 'Spawning...' : 'Spawn'}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}

function DetailField({ label, value }: Readonly<{ label: string; value: string }>) {
  return (
    <div>
      <span className="text-text-muted text-xs block mb-0.5">{label}</span>
      <span className="text-text text-sm font-medium">{value}</span>
    </div>
  )
}