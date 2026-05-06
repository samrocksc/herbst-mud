import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { useState } from 'react'
import { apiGet, apiPut, apiDelete } from '../../utils/apiFetch'
import { PageHeader } from '../../components/PageHeader'
import { Button } from '../../components/Button'

// ─── Types ──────────────────────────────────────────────────────────────────

type ItemInstance = Readonly<{
  id: number
  name: string
  description: string
  slot: string
  level: number
  weight: number
  isEquipped: boolean
  isImmovable: boolean
  color: string
  isVisible: boolean
  itemType: string
  ownerId: number | null
  roomId: number | null
  equipment_template_id: number | null
}>

type EditForm = Readonly<{
  name: string
  description: string
  slot: string
  itemType: string
  level: number
  weight: number
  color: string
  ownerId: number | null
  roomId: number | null
  isVisible: boolean
  isImmovable: boolean
  isEquipped: boolean
}>

const API = `${window.location.origin}`

const SLOT_OPTIONS = ['weapon', 'head', 'chest', 'legs', 'feet', 'hands', 'accessory', 'none'] as const
const ITEM_TYPE_OPTIONS = ['weapon', 'armor', 'consumable', 'key', 'misc'] as const

// ─── Route ─────────────────────────────────────────────────────────────────

export const Route = createFileRoute('/_auth/items/$itemId/instances/$instanceId')({
  component: ItemInstanceDetail,
})

// ─── Component ──────────────────────────────────────────────────────────────

function ItemInstanceDetail() {
  const { itemId, instanceId } = Route.useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const [editing, setEditing] = useState(false)
  const [confirmDelete, setConfirmDelete] = useState(false)

  const [form, setForm] = useState<EditForm>({
    name: '',
    description: '',
    slot: 'none',
    itemType: 'misc',
    level: 1,
    weight: 0,
    color: '',
    ownerId: null,
    roomId: null,
    isVisible: true,
    isImmovable: false,
    isEquipped: false,
  })

  // ── Queries ────────────────────────────────────────────────────────────────

  const { data: instance, isLoading, error } = useQuery<ItemInstance>({
    queryKey: ['item-instances', instanceId],
    queryFn: () => apiGet<ItemInstance>(`${API}/api/item-instances/${instanceId}`),
  })

  // ── Mutations ──────────────────────────────────────────────────────────────

  const updateMutation = useMutation({
    mutationFn: (body: EditForm) => apiPut(`${API}/api/item-instances/${instanceId}`, body),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances', instanceId] })
      queryClient.invalidateQueries({ queryKey: ['item-instances'] })
      setEditing(false)
    },
  })

  const deleteMutation = useMutation({
    mutationFn: () => apiDelete(`${API}/api/item-instances/${instanceId}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances'] })
      navigate({ to: '/items/$itemId', params: { itemId } })
    },
  })

  // ── Handlers ───────────────────────────────────────────────────────────────

  const startEditing = () => {
    if (!instance) return
    setForm({
      name: instance.name,
      description: instance.description,
      slot: instance.slot,
      itemType: instance.itemType,
      level: instance.level,
      weight: instance.weight,
      color: instance.color,
      ownerId: instance.ownerId,
      roomId: instance.roomId,
      isVisible: instance.isVisible,
      isImmovable: instance.isImmovable,
      isEquipped: instance.isEquipped,
    })
    setEditing(true)
  }

  const handleSave = () => {
    updateMutation.mutate(form)
  }

  const handleDelete = () => {
    if (confirmDelete) {
      deleteMutation.mutate()
    } else {
      setConfirmDelete(true)
    }
  }

  // ── Loading / Error ────────────────────────────────────────────────────────

  if (isLoading) {
    return (
      <div className="p-8">
        <PageHeader title="Loading..." backTo={`/items/${itemId}`} />
        <div className="text-text-muted">Loading item instance...</div>
      </div>
    )
  }

  if (error || !instance) {
    return (
      <div className="p-8">
        <PageHeader title="Error" backTo={`/items/${itemId}`} />
        <div className="text-danger">
          Failed to load instance: {error?.message ?? 'Unknown error'}
        </div>
      </div>
    )
  }

  // ── Render ─────────────────────────────────────────────────────────────────

  return (
    <div className="p-8">
      <PageHeader
        title={instance.name}
        backTo={`/items/${itemId}`}
        actions={
          !editing ? (
            <div className="flex items-center gap-2">
              <Button variant="primary" size="sm" onClick={startEditing}>
                Edit
              </Button>
              <Button
                variant="danger"
                size="sm"
                onClick={handleDelete}
                disabled={deleteMutation.isPending}
              >
                {confirmDelete
                  ? 'Confirm Delete?'
                  : deleteMutation.isPending
                    ? 'Deleting...'
                    : 'Delete'}
              </Button>
              {confirmDelete && (
                <Button variant="secondary" size="sm" onClick={() => setConfirmDelete(false)}>
                  Cancel
                </Button>
              )}
            </div>
          ) : undefined
        }
      />

      {editing ? (
        <div className="max-w-2xl">
          <div className="bg-surface-muted rounded-lg p-6 border border-border">
            <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Edit Instance</h2>

            <div className="grid grid-cols-2 gap-4 mb-4">
              <div>
                <label className="text-text-muted text-xs block mb-1">Name</label>
                <input
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div>
                <label className="text-text-muted text-xs block mb-1">Slot</label>
                <select
                  value={form.slot}
                  onChange={(e) => setForm({ ...form, slot: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                >
                  {SLOT_OPTIONS.map((s) => (
                    <option key={s} value={s}>{s}</option>
                  ))}
                </select>
              </div>

              <div>
                <label className="text-text-muted text-xs block mb-1">Item Type</label>
                <select
                  value={form.itemType}
                  onChange={(e) => setForm({ ...form, itemType: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                >
                  {ITEM_TYPE_OPTIONS.map((t) => (
                    <option key={t} value={t}>{t}</option>
                  ))}
                </select>
              </div>

              <div>
                <label className="text-text-muted text-xs block mb-1">Level</label>
                <input
                  type="number"
                  value={form.level}
                  onChange={(e) => setForm({ ...form, level: parseInt(e.target.value) || 0 })}
                  min={0}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div>
                <label className="text-text-muted text-xs block mb-1">Weight</label>
                <input
                  type="number"
                  value={form.weight}
                  onChange={(e) => setForm({ ...form, weight: parseInt(e.target.value) || 0 })}
                  min={0}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div>
                <label className="text-text-muted text-xs block mb-1">Color</label>
                <input
                  type="text"
                  value={form.color}
                  onChange={(e) => setForm({ ...form, color: e.target.value })}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div>
                <label className="text-text-muted text-xs block mb-1">Owner ID</label>
                <input
                  type="number"
                  value={form.ownerId ?? ''}
                  onChange={(e) => {
                    const val = e.target.value === '' ? null : parseInt(e.target.value) || null
                    setForm({ ...form, ownerId: val })
                  }}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div>
                <label className="text-text-muted text-xs block mb-1">Room ID</label>
                <input
                  type="number"
                  value={form.roomId ?? ''}
                  onChange={(e) => {
                    const val = e.target.value === '' ? null : parseInt(e.target.value) || null
                    setForm({ ...form, roomId: val })
                  }}
                  className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
                />
              </div>

              <div className="flex items-center gap-4 col-span-2 pt-1">
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.isVisible}
                    onChange={(e) => setForm({ ...form, isVisible: e.target.checked })}
                    className="cursor-pointer"
                  />
                  Visible
                </label>
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.isImmovable}
                    onChange={(e) => setForm({ ...form, isImmovable: e.target.checked })}
                    className="cursor-pointer"
                  />
                  Immovable
                </label>
                <label className="flex items-center gap-2 text-text-muted text-xs cursor-pointer">
                  <input
                    type="checkbox"
                    checked={form.isEquipped}
                    onChange={(e) => setForm({ ...form, isEquipped: e.target.checked })}
                    className="cursor-pointer"
                  />
                  Equipped
                </label>
              </div>
            </div>

            <div className="flex gap-2">
              <Button
                variant="primary"
                onClick={handleSave}
                disabled={updateMutation.isPending}
              >
                {updateMutation.isPending ? 'Saving...' : 'Save'}
              </Button>
              <Button variant="secondary" onClick={() => setEditing(false)}>
                Cancel
              </Button>
            </div>

            {updateMutation.isError && (
              <div className="mt-3 text-danger text-sm">
                Failed to save: {(updateMutation.error as Error)?.message}
              </div>
            )}
          </div>
        </div>
      ) : (
        <div className="max-w-2xl">
          <div className="bg-surface-muted rounded-lg p-6 border border-border mb-6">
            <h2 className="mt-0 mb-4 text-text text-lg font-semibold">Instance Stats</h2>
            <div className="grid grid-cols-2 gap-x-6 gap-y-3">
              <DetailField label="ID" value={String(instance.id)} />
              <DetailField label="Name" value={instance.name} />
              <DetailField label="Type" value={instance.itemType} />
              <DetailField label="Slot" value={instance.slot} />
              <DetailField label="Level" value={String(instance.level)} />
              <DetailField label="Weight" value={String(instance.weight)} />
              {instance.ownerId ? (
                <DetailField
                  label="Owner ID"
                  value={
                    <span className="text-primary">
                      {String(instance.ownerId)}
                    </span>
                  }
                />
              ) : (
                <DetailField label="Owner ID" value="None" />
              )}
              {instance.roomId && !instance.ownerId ? (
                <DetailField label="Room ID" value={String(instance.roomId)} />
              ) : (
                <DetailField label="Room ID" value={instance.roomId ? String(instance.roomId) : 'None'} />
              )}
              <DetailField label="Color" value={instance.color || 'none'} />
              <BoolBadge value={instance.isVisible} label="Visible" />
              <BoolBadge value={instance.isImmovable} label="Immovable" />
              <BoolBadge value={instance.isEquipped} label="Equipped" />
            </div>
          </div>

          {instance.ownerId && (
            <div className="bg-surface-muted rounded-lg p-4 border border-border mb-6">
              <span className="text-text-muted text-sm">
                Held by Character{' '}
                <span className="text-primary font-medium">
                  #{instance.ownerId}
                </span>
              </span>
            </div>
          )}

          {!instance.ownerId && instance.roomId && (
            <div className="bg-surface-muted rounded-lg p-4 border border-border mb-6">
              <span className="text-text-muted text-sm">
                In Room{' '}
                <span className="text-text font-medium">#{instance.roomId}</span>
              </span>
            </div>
          )}

          {deleteMutation.isError && (
            <div className="text-danger text-sm">
              Failed to delete: {(deleteMutation.error as Error)?.message}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

// ─── Helpers ────────────────────────────────────────────────────────────────

function DetailField({ label, value }: Readonly<{ label: string; value: React.ReactNode }>) {
  return (
    <div>
      <span className="text-text-muted text-xs block mb-0.5">{label}</span>
      <span className="text-text text-sm font-medium">{value}</span>
    </div>
  )
}

function BoolBadge({ value, label }: Readonly<{ value: boolean; label: string }>) {
  return (
    <div>
      <span className="text-text-muted text-xs block mb-0.5">{label}</span>
      <span
        className={
          value
            ? 'inline-block px-2 py-0.5 rounded text-xs font-medium bg-green-900/30 text-green-400 border border-green-700/40'
            : 'inline-block px-2 py-0.5 rounded text-xs font-medium bg-red-900/30 text-red-400 border border-red-700/40'
        }
      >
        {value ? 'Yes' : 'No'}
      </span>
    </div>
  )
}