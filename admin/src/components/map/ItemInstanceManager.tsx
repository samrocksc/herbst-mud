import { useState, useCallback } from 'react'
import { Button } from '../Button'
import { apiGet, apiPost, apiPut, apiDelete } from '../../utils/apiFetch'
import { showToast } from '../Toast'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { ItemEditRow } from './ItemEditRow'
import { ItemSpawnModal } from './ItemSpawnModal'

// ─── Types ──────────────────────────────────────────────────────────────────

export type EquipmentTemplate = Readonly<{
  equipment_template_id: string
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

export type ItemInstanceView = Readonly<{
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
  equipment_template_id: string
  effect_type: string
  effect_value: number
  effect_duration: number
  healing: number
  effect: string
  isContainer: boolean
  containerCapacity: number
  isLocked: boolean
}>

export type SpawnFormData = {
  template_id: string
  name: string
  description: string
  slot: string
  level: number
  weight: number
  color: string
  room_id: number
}

export type EditFormData = {
  name: string
  description: string
  slot: string
  level: number
  weight: number
  color: string
}

type ItemInstanceManagerProps = Readonly<{ roomId: number }>

// ─── Component ──────────────────────────────────────────────────────────────

export function ItemInstanceManager({ roomId }: ItemInstanceManagerProps) {
  const queryClient = useQueryClient()

  const instancesQuery = useQuery({
    queryKey: ['item-instances', roomId],
    queryFn: async (): Promise<ItemInstanceView[]> =>
      apiGet<ItemInstanceView[]>(`${window.location.origin}/api/item-instances?roomId=${roomId}`),
  })

  const templatesQuery = useQuery({
    queryKey: ['equipment-templates'],
    queryFn: async (): Promise<EquipmentTemplate[]> =>
      apiGet<EquipmentTemplate[]>(`${window.location.origin}/api/equipment-templates`),
  })

  const createMutation = useMutation({
    mutationFn: async (input: Record<string, unknown>): Promise<ItemInstanceView> =>
      apiPost<ItemInstanceView>(`${window.location.origin}/api/item-instances`, input),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances'] }) },
  })

  const updateMutation = useMutation({
    mutationFn: async (args: { id: number; update: Record<string, unknown> }): Promise<ItemInstanceView> =>
      apiPut<ItemInstanceView>(`${window.location.origin}/api/item-instances/${args.id}`, args.update),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances'] }) },
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: number): Promise<void> =>
      apiDelete(`${window.location.origin}/api/item-instances/${id}`),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances'] }) },
  })

  const emptySpawnForm = useCallback((): SpawnFormData => ({
    template_id: '', name: '', description: '', slot: 'none',
    level: 0, weight: 0, color: '', room_id: roomId,
  }), [roomId])

  const [showSpawn, setShowSpawn] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null)
  const [spawnForm, setSpawnForm] = useState<SpawnFormData>(emptySpawnForm())
  const [editForm, setEditForm] = useState<Partial<EditFormData>>({})

  const selectedTemplate = templatesQuery.data?.find(
    (t) => t.equipment_template_id === spawnForm.template_id
  )

  const applyTemplateDefaults = useCallback(
    (templateId: string) => {
      const t = templatesQuery.data?.find((tmpl) => tmpl.equipment_template_id === templateId)
      if (!t) return
      setSpawnForm((f) => ({
        ...f, template_id: t.equipment_template_id, name: t.name,
        description: t.description, slot: t.slot, level: t.level,
        weight: t.weight, color: t.color,
      }))
    },
    [templatesQuery.data]
  )

  const handleSpawn = useCallback(async () => {
    if (!spawnForm.template_id) return
    try {
      const payload: Record<string, unknown> = {
        equipment_template_id: spawnForm.template_id, room_id: spawnForm.room_id,
      }
      if (spawnForm.name.trim()) payload.name = spawnForm.name.trim()
      if (spawnForm.description.trim()) payload.description = spawnForm.description.trim()
      if (spawnForm.slot && spawnForm.slot !== 'none') payload.slot = spawnForm.slot
      if (spawnForm.level > 0) payload.level = spawnForm.level
      if (spawnForm.weight > 0) payload.weight = spawnForm.weight
      if (spawnForm.color.trim()) payload.color = spawnForm.color.trim()
      await createMutation.mutateAsync(payload)
      setShowSpawn(false)
      setSpawnForm(emptySpawnForm())
    } catch (err) {
      showToast(`Spawn failed: ${(err as Error)?.message ?? 'Unknown error'}`)
    }
  }, [spawnForm, createMutation, emptySpawnForm])

  const startEdit = useCallback((inst: ItemInstanceView) => {
    setEditingId(inst.id)
    setConfirmDeleteId(null)
    setEditForm({ name: inst.name, description: inst.description, slot: inst.slot, level: inst.level, weight: inst.weight, color: inst.color })
  }, [])

  const handleUpdate = useCallback(async () => {
    if (editingId === null) return
    try {
      const update: Record<string, unknown> = {}
      if (editForm.name !== undefined) update.name = editForm.name
      if (editForm.description !== undefined) update.description = editForm.description
      if (editForm.slot !== undefined) update.slot = editForm.slot
      if (editForm.level !== undefined) update.level = editForm.level
      if (editForm.weight !== undefined) update.weight = editForm.weight
      if (editForm.color !== undefined) update.color = editForm.color
      await updateMutation.mutateAsync({ id: editingId, update })
      setEditingId(null)
      setEditForm({})
    } catch (err) {
      showToast(`Update failed: ${(err as Error)?.message ?? 'Unknown error'}`)
    }
  }, [editingId, editForm, updateMutation])

  const handleDelete = useCallback(async (id: number) => {
    try {
      await deleteMutation.mutateAsync(id)
      setConfirmDeleteId(null)
      if (editingId === id) { setEditingId(null); setEditForm({}) }
    } catch (err) {
      showToast(`Delete failed: ${(err as Error)?.message ?? 'Unknown error'}`)
    }
  }, [deleteMutation, editingId])

  const handleOpenSpawn = useCallback(() => {
    setSpawnForm(emptySpawnForm())
    setShowSpawn(true)
  }, [emptySpawnForm])

  const instances = instancesQuery.data ?? []

  if (instancesQuery.isLoading) return (
    <div className="mb-3">
      <strong className="text-success text-xs">Items:</strong>
      <div className="text-text-muted text-[10px] mt-1">Loading...</div>
    </div>
  )

  if (instancesQuery.error) return (
    <div className="mb-3">
      <strong className="text-success text-xs">Items:</strong>
      <div className="text-danger text-[10px] mt-1">Error loading items</div>
    </div>
  )

  return (
    <div className="mb-3">
      <div className="flex items-center justify-between mb-1">
        <strong className="text-success text-xs">Items:</strong>
        <Button variant="primary" size="sm" className="!px-1.5 !py-0 !text-[10px]" onClick={handleOpenSpawn}>
          + Add Instance
        </Button>
      </div>
      {instances.length === 0 ? (
        <div className="text-text-muted text-[10px]">No item instances in this room.</div>
      ) : (
        <div className="mt-1 flex flex-col gap-1">
          {instances.map((inst) => (
            <div key={inst.id} className="p-1 bg-surface-muted rounded text-xs text-text">
              {editingId === inst.id ? (
                <ItemEditRow
                  inst={inst} editForm={editForm} setEditForm={setEditForm}
                  onSave={handleUpdate} onCancel={() => { setEditingId(null); setEditForm({}) }}
                  isPending={updateMutation.isPending}
                  error={updateMutation.error as Error | null}
                />
              ) : (
                <div className="flex justify-between items-center">
                  <div>
                    <span className="font-medium">{inst.name}</span>{' '}
                    <span className="text-text-muted">{inst.itemType} lv.{inst.level}</span>
                    {!inst.isVisible && <span className="text-warning ml-1 text-[10px]">(hidden)</span>}
                    {inst.isImmovable && <span className="text-danger ml-1 text-[10px]">(immovable)</span>}
                  </div>
                  <div className="flex gap-0.5">
                    <Button variant="ghost" size="sm" className="!px-0.5 !py-0" onClick={() => startEdit(inst)} aria-label={`Edit ${inst.name}`}>✏️</Button>
                    <Button
                      variant={confirmDeleteId === inst.id ? 'secondary' : 'ghost'}
                      size="sm" className="!px-0.5 !py-0"
                      onClick={() => confirmDeleteId === inst.id ? handleDelete(inst.id) : setConfirmDeleteId(inst.id)}
                      aria-label={confirmDeleteId === inst.id ? `Confirm delete ${inst.name}` : `Delete ${inst.name}`}
                    >
                      {confirmDeleteId === inst.id ? '❓' : '🗑'}
                    </Button>
                  </div>
                </div>
              )}
              {confirmDeleteId === inst.id && editingId !== inst.id && (
                <div className="mt-1 text-[10px] text-text">
                  Confirm delete?{' '}
                  <Button variant="danger" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={() => handleDelete(inst.id)}>Yes</Button>{' '}
                  <Button variant="ghost" size="sm" className="!px-1 !py-0 !text-[10px]" onClick={() => setConfirmDeleteId(null)}>No</Button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
      <ItemSpawnModal
        isOpen={showSpawn} onClose={() => setShowSpawn(false)}
        spawnForm={spawnForm} setSpawnForm={setSpawnForm}
        onSpawn={handleSpawn} isPending={createMutation.isPending}
        error={createMutation.error as Error | null}
        templates={templatesQuery.data ?? []} templatesLoading={templatesQuery.isLoading}
        selectedTemplate={selectedTemplate} applyTemplateDefaults={applyTemplateDefaults}
      />
    </div>
  )
}