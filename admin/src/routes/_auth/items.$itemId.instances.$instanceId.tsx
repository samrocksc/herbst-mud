import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { apiGet, apiDelete } from '../../utils/apiFetch'
import { PageHeader } from '../../components/PageHeader'
import { Button } from '../../components/Button'
import { InstanceDetailView } from './-items.$itemId.instances.$instanceId.detailView'
import { InstanceEditForm } from './-items.$itemId.instances.$instanceId.editForm'

export const Route = createFileRoute('/_auth/items/$itemId/instances/$instanceId')({
  component: ItemInstanceDetail,
})

type ItemInstance = Readonly<{
  id: number; name: string; description: string; slot: string; level: number
  weight: number; isEquipped: boolean; isImmovable: boolean; color: string
  isVisible: boolean; itemType: string; ownerId: number | null; roomId: number | null
  equipment_template_id: number | null; armor_rating: number; armor_type: string
  rarity: string; skill_requirement: string; skill_requirement_level: number
  damage_dice_count: number; damage_dice_sides: number; damage_bonus: number
  damage_type: string; weapon_type: string; is_two_handed: boolean
  stats: Record<string, unknown>
}>

const API = `${window.location.origin}`

function ItemInstanceDetail() {
  const { itemId, instanceId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [editing, setEditing] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState(false)

  const { data: instance, isLoading, error } = useQuery<ItemInstance>({
    queryKey: ['item-instances', instanceId],
    queryFn: () => apiGet<ItemInstance>(`${API}/api/item-instances/${instanceId}`),
  })

  const deleteMutation = useMutation({
    mutationFn: () => apiDelete(`${API}/api/item-instances/${instanceId}`),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances'] }); navigate({ to: '/items/$itemId', params: { itemId } }) },
  })

  const handleDelete = () => { if (confirmDelete) deleteMutation.mutate(); else setConfirmDelete(true) }

  if (isLoading) return <div className="p-8"><PageHeader title="Loading..." backTo={`/items/${itemId}`} /></div>
  if (error || !instance) return <div className="p-8"><PageHeader title="Error" backTo={`/items/${itemId}`} /><div className="text-danger">Failed to load instance</div></div>

  return (
    <div className="p-8">
      <PageHeader title={instance.name} backTo={`/items/${itemId}`} actions={
        !editing ? (
          <div className="flex items-center gap-2">
            <Button variant="primary" size="sm" onClick={() => setEditing(true)}>Edit</Button>
            <Button variant="danger" size="sm" onClick={handleDelete} disabled={deleteMutation.isPending}>
              {confirmDelete ? 'Confirm Delete?' : deleteMutation.isPending ? 'Deleting...' : 'Delete'}</Button>
            {confirmDelete && <Button variant="secondary" size="sm" onClick={() => setConfirmDelete(false)}>Cancel</Button>}
          </div>
        ) : undefined
      } />
      {editing ? <InstanceEditForm instance={instance} instanceId={instanceId} onDone={() => setEditing(false)} /> : <InstanceDetailView instance={instance} />}
      {deleteMutation.isError && <div className="text-danger text-sm mt-2">Failed to delete: {(deleteMutation.error as Error)?.message}</div>}
    </div>
  )
}