import { createFileRoute, Link, Outlet } from '@tanstack/react-router'
import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useLocation } from '@tanstack/react-router'
import { apiGet } from '../../utils/apiFetch'
import { PageHeader } from '../../components/PageHeader'
import { DataTable, type Column } from '../../components/DataTable'
import { Button } from '../../components/Button'
import { TemplateEditForm } from './-items.$itemId.editForm'
import { ItemDetailView } from './-items.$itemId.detailView'
import type { ItemInstance } from '../../hooks/useItemInstances'
import type { EquipmentTemplate as ItemTemplate } from '../../hooks/useEquipmentTemplates'

export const Route = createFileRoute('/_auth/items/$itemId')({
  component: ItemDetail,
})

const instanceColumns: Column<ItemInstance>[] = [
  { header: 'ID', accessor: 'id', render: (_, row) => (
    <Link to="/items/$itemId/instances/$instanceId" params={{ itemId: row.equipment_template_id, instanceId: String(row.id) }}
      className="text-primary no-underline hover:underline font-mono text-xs">{row.id}</Link>
  )},
  { header: 'Name', accessor: 'name' },
  { header: 'Location', accessor: 'ownerId', render: (_, row) => {
    if (row.ownerId != null) return <Link to="/characters/$characterId" params={{ characterId: String(row.ownerId) }} className="text-primary text-xs no-underline hover:underline">Character #{row.ownerId}</Link>
    if (row.roomId > 0) return <Link to="/map" className="text-primary text-xs no-underline hover:underline">In Room #{row.roomId}</Link>
    return <span className="text-text-muted text-xs">Nowhere</span>
  }},
]

const API = `${window.location.origin}`

function ItemDetail() {
  const { itemId } = Route.useParams()
  const location = useLocation()
  const [editing, setEditing] = useState(false)

  const templateQuery = useQuery({
    queryKey: ['item-template', itemId],
    queryFn: () => apiGet<ItemTemplate>(`${API}/api/equipment-templates/${itemId}`),
  })

  const instancesQuery = useQuery({
    queryKey: ['item-instances', 'template', itemId],
    queryFn: () => apiGet<ItemInstance[]>(`${API}/api/item-instances?templateId=${itemId}`),
  })

  const template = templateQuery.data
  const instances = instancesQuery.data ?? []

  // Render outlet for child routes (spawn, instances, etc.)
  if (location.pathname !== `/items/${itemId}`) {
    return <Outlet />
  }

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
          <Link to="/items/$itemId/spawn" params={{ itemId }} className="no-underline">
            <Button variant="secondary" size="sm">+ Add Instance</Button>
          </Link>
        </div>
        {instancesQuery.isError && <div className="text-danger text-xs mb-3">Failed to load instances</div>}
        <DataTable<ItemInstance> columns={instanceColumns} data={instances} getKey={(row) => row.id} emptyMessage="No instances found." variant="dark" />
      </div>
    </div>
  )
}
