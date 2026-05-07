import { createFileRoute, Link } from '@tanstack/react-router'
import { useState } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiGet, apiPost } from '../../utils/apiFetch'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { TemplateEditForm } from './-items.$itemId.editForm'
import { ItemDetailView } from './-items.$itemId.detailView'
import { SpawnModal } from './-items.$itemId.spawnModal'

export const Route = createFileRoute('/_auth/items/$itemId')({
  component: ItemDetail,
})

type ItemInstance = Readonly<{
  id: number; name: string; ownerId: number | null; roomId: number; equipment_template_id: string
}>

const instanceColumns: Column<ItemInstance>[] = [
  { header: 'ID', accessor: 'id', render: (_, row) => (
    <Link to="/items/$itemId/instances/$instanceId" params={{ itemId: row.equipment_template_id, instanceId: String(row.id) }}
      className="text-primary no-underline hover:underline font-mono text-xs">{row.id}</Link>
  )},
  { header: 'Name', accessor: 'name' },
  { header: 'Location', accessor: 'ownerId', render: (_, row) => {
    if (row.ownerId != null) return <Link to="/npcs" className="text-primary text-xs no-underline hover:underline">Held by Character #{row.ownerId}</Link>
    if (row.roomId > 0) return <Link to="/map" className="text-primary text-xs no-underline hover:underline">In Room #{row.roomId}</Link>
    return <span className="text-text-muted text-xs">Nowhere</span>
  }},
]

const API = `${window.location.origin}`

function ItemDetail() {
  const { itemId } = Route.useParams()
  const queryClient = useQueryClient()
  const [modalOpen, setModalOpen] = useState(false)
  const [spawnRoomId, setSpawnRoomId] = useState('')
  const [editing, setEditing] = useState(false)

  const templateQuery = useQuery({
    queryKey: ['item-template', itemId],
    queryFn: () => apiGet<ItemTemplate>(`${API}/api/equipment-templates/${itemId}`),
  })

  const instancesQuery = useQuery({
    queryKey: ['item-instances', 'template', itemId],
    queryFn: () => apiGet<ItemInstance[]>(`${API}/api/item-instances?templateId=${itemId}`),
  })

  const spawnMutation = useMutation({
    mutationFn: (roomId: number) => apiPost(`${API}/api/item-instances`, { equipment_template_id: itemId, room_id: roomId }),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances', 'template', itemId] }); setModalOpen(false); setSpawnRoomId('') },
  })

  const template = templateQuery.data
  const instances = instancesQuery.data ?? []

  if (templateQuery.isLoading) return <div className="p-8"><PageHeader title="Loading..." backTo="/items" /></div>
  if (templateQuery.error || !template) return <div className="p-8"><PageHeader title="Error" backTo="/items" /><div className="text-danger">Failed to load item</div></div>

  return (
    <div className="p-6 max-w-[1200px] mx-auto">
      <PageHeader title={template.name} backTo="/items" actions={
        <Button variant={editing ? 'secondary' : 'primary'} size="sm" onClick={() => setEditing(!editing)}>{editing ? 'Cancel' : 'Edit'}</Button>
      } />
      {editing ? <TemplateEditForm template={template} itemId={itemId} onDone={() => setEditing(false)} /> : <ItemDetailView template={template} />}
      <div className="bg-surface-muted rounded-lg p-6 border border-border mt-6">
        <div className="flex items-center justify-between mb-4">
          <h2 className="m-0 text-text text-lg font-semibold">Instances ({instances.length})</h2>
          <Button variant="secondary" size="sm" onClick={() => setModalOpen(true)}>+ Add Instance</Button>
        </div>
        {instancesQuery.isError && <div className="text-danger text-xs mb-3">Failed to load instances</div>}
        <DataTable<ItemInstance> columns={instanceColumns} data={instances} getKey={(row) => row.id} emptyMessage="No instances found." variant="dark" />
      </div>
      <SpawnModal open={modalOpen} onClose={() => { setModalOpen(false); setSpawnRoomId('') }} spawnRoomId={spawnRoomId} setSpawnRoomId={setSpawnRoomId}
        onSpawn={() => spawnMutation.mutate(Number(spawnRoomId))} isLoading={spawnMutation.isPending} error={spawnMutation.isError ? (spawnMutation.error as Error)?.message ?? 'Failed' : null} />
    </div>
  )
}

/** Re-exported type for use by child components. */
export type ItemTemplate = Readonly<{
  id: string; name: string; description: string; slot: string; level: number; weight: number
  item_type: string; stats: Record<string, unknown>; color: string; is_visible: boolean
  is_immovable: boolean; effect_type: string; effect_value: number; effect_duration: number
  is_container: boolean; container_capacity: number; is_locked: boolean; key_item_id: string
  reveal_condition: string; expires_at: string | null; armor_rating: number; armor_type: string
  rarity: string; skill_requirement: string; skill_requirement_level: number
  damage_dice_count: number; damage_dice_sides: number; damage_bonus: number
  damage_type: string; weapon_type: string; is_two_handed: boolean
}>