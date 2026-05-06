import { useState, useCallback } from 'react'
import { Button } from '../Button'
import { Modal } from '../Modal'
import { apiGet, apiPost, apiPut, apiDelete } from '../../utils/apiFetch'
import { logError } from '../../utils/log'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

// ─── Types ──────────────────────────────────────────────────────────────────

type EquipmentTemplate = Readonly<{
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

type ItemInstanceView = Readonly<{
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

type SpawnFormData = {
  template_id: string
  name: string
  description: string
  slot: string
  level: number
  weight: number
  color: string
  room_id: number
}

type ItemInstanceManagerProps = Readonly<{
  roomId: number
}>

// ─── Component ──────────────────────────────────────────────────────────────

export function ItemInstanceManager({ roomId }: ItemInstanceManagerProps) {
  const queryClient = useQueryClient()

  // Fetch item instances for this room
  const instancesQuery = useQuery({
    queryKey: ['item-instances', roomId],
    queryFn: async (): Promise<ItemInstanceView[]> => {
      return apiGet<ItemInstanceView[]>(
        `${window.location.origin}/api/item-instances?roomId=${roomId}`
      )
    },
  })

  // Fetch equipment templates for spawn form
  const templatesQuery = useQuery({
    queryKey: ['equipment-templates'],
    queryFn: async (): Promise<EquipmentTemplate[]> => {
      return apiGet<EquipmentTemplate[]>(
        `${window.location.origin}/api/equipment-templates`
      )
    },
  })

  // Create mutation
  const createMutation = useMutation({
    mutationFn: async (input: Record<string, unknown>): Promise<ItemInstanceView> => {
      return apiPost<ItemInstanceView>(
        `${window.location.origin}/api/item-instances`,
        input
      )
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances'] })
    },
  })

  // Update mutation
  const updateMutation = useMutation({
    mutationFn: async (args: {
      id: number
      update: Record<string, unknown>
    }): Promise<ItemInstanceView> => {
      return apiPut<ItemInstanceView>(
        `${window.location.origin}/api/item-instances/${args.id}`,
        args.update
      )
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances'] })
    },
  })

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: async (id: number): Promise<void> => {
      await apiDelete(`${window.location.origin}/api/item-instances/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['item-instances'] })
    },
  })

  const emptySpawnForm = (): SpawnFormData => ({
    template_id: '',
    name: '',
    description: '',
    slot: 'none',
    level: 0,
    weight: 0,
    color: '',
    room_id: roomId,
  })

  const [showSpawn, setShowSpawn] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null)
  const [spawnForm, setSpawnForm] = useState<SpawnFormData>(emptySpawnForm())
  const [editForm, setEditForm] = useState<Partial<SpawnFormData>>({})

  const selectedTemplate = templatesQuery.data?.find(
    (t) => t.equipment_template_id === spawnForm.template_id
  )

  const applyTemplateDefaults = useCallback(
    (templateId: string) => {
      const t = templatesQuery.data?.find((tmpl) => tmpl.equipment_template_id === templateId)
      if (!t) return
      setSpawnForm((f) => ({
        ...f,
        template_id: t.equipment_template_id,
        name: t.name,
        description: t.description,
        slot: t.slot,
        level: t.level,
        weight: t.weight,
        color: t.color,
      }))
    },
    [templatesQuery.data]
  )

  const handleSpawn = useCallback(async () => {
    if (!spawnForm.template_id) return
    try {
      const payload: Record<string, unknown> = {
        equipment_template_id: spawnForm.template_id,
        room_id: spawnForm.room_id,
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
      logError('Spawn item instance:', err)
    }
  }, [spawnForm, createMutation])

  const startEdit = useCallback(
    (inst: ItemInstanceView) => {
      setEditingId(inst.id)
      setConfirmDeleteId(null)
      setEditForm({
        name: inst.name,
        description: inst.description,
        slot: inst.slot,
        level: inst.level,
        weight: inst.weight,
        color: inst.color,
      })
    },
    []
  )

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
      logError('Update item instance:', err)
    }
  }, [editingId, editForm, updateMutation])

  const handleDelete = useCallback(
    async (id: number) => {
      try {
        await deleteMutation.mutateAsync(id)
        setConfirmDeleteId(null)
        if (editingId === id) {
          setEditingId(null)
          setEditForm({})
        }
      } catch (err) {
        logError('Delete item instance:', err)
      }
    },
    [deleteMutation, editingId]
  )

  const handleOpenSpawn = useCallback(() => {
    setSpawnForm(emptySpawnForm())
    setShowSpawn(true)
  }, [roomId])

  const instances = instancesQuery.data ?? []

  if (instancesQuery.isLoading) {
    return (
      <div className="mb-3">
        <strong className="text-success text-xs">Items:</strong>
        <div className="text-text-muted text-[10px] mt-1">Loading...</div>
      </div>
    )
  }

  if (instancesQuery.error) {
    return (
      <div className="mb-3">
        <strong className="text-success text-xs">Items:</strong>
        <div className="text-danger text-[10px] mt-1">
          Error loading items
        </div>
      </div>
    )
  }

  return (
    <div className="mb-3">
      <div className="flex items-center justify-between mb-1">
        <strong className="text-success text-xs">Items:</strong>
        <Button
          variant="primary"
          size="sm"
          className="!px-1.5 !py-0 !text-[10px]"
          onClick={handleOpenSpawn}
        >
          + Add Instance
        </Button>
      </div>

      {instances.length === 0 ? (
        <div className="text-text-muted text-[10px]">No item instances in this room.</div>
      ) : (
        <div className="mt-1 flex flex-col gap-1">
          {instances.map((inst) => (
            <div
              key={inst.id}
              className="p-1 bg-surface-muted rounded text-xs text-text"
            >
              {editingId === inst.id ? (
                <div className="space-y-1">
                  <div className="font-medium">{inst.name}</div>
                  <div className="flex gap-1">
                    <div className="flex-1">
                      <label className="text-text-muted block text-[10px] mb-0.5">Name</label>
                      <input
                        type="text"
                        value={editForm.name ?? inst.name}
                        onChange={(e) =>
                          setEditForm((f) => ({ ...f, name: e.target.value }))
                        }
                        className="w-full p-0.5 bg-surface border border-border rounded text-text text-[10px]"
                      />
                    </div>
                    <div className="flex-1">
                      <label className="text-text-muted block text-[10px] mb-0.5">Level</label>
                      <input
                        type="number"
                        value={editForm.level ?? inst.level}
                        onChange={(e) =>
                          setEditForm((f) => ({
                            ...f,
                            level: parseInt(e.target.value, 10) || 0,
                          }))
                        }
                        className="w-full p-0.5 bg-surface border border-border rounded text-text text-[10px]"
                      />
                    </div>
                  </div>
                  <div className="flex gap-1">
                    <div className="flex-1">
                      <label className="text-text-muted block text-[10px] mb-0.5">Slot</label>
                      <input
                        type="text"
                        value={editForm.slot ?? inst.slot}
                        onChange={(e) =>
                          setEditForm((f) => ({ ...f, slot: e.target.value }))
                        }
                        className="w-full p-0.5 bg-surface border border-border rounded text-text text-[10px]"
                      />
                    </div>
                    <div className="flex-1">
                      <label className="text-text-muted block text-[10px] mb-0.5">Weight</label>
                      <input
                        type="number"
                        value={editForm.weight ?? inst.weight}
                        onChange={(e) =>
                          setEditForm((f) => ({
                            ...f,
                            weight: parseInt(e.target.value, 10) || 0,
                          }))
                        }
                        className="w-full p-0.5 bg-surface border border-border rounded text-text text-[10px]"
                      />
                    </div>
                  </div>
                  <div className="flex gap-1">
                    <Button
                      variant="primary"
                      size="sm"
                      className="!px-1 !py-0 !text-[10px]"
                      onClick={handleUpdate}
                      disabled={updateMutation.isPending}
                    >
                      {updateMutation.isPending ? 'Saving...' : 'Save'}
                    </Button>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="!px-1 !py-0 !text-[10px]"
                      onClick={() => setEditingId(null)}
                    >
                      Cancel
                    </Button>
                  </div>
                </div>
              ) : (
                <div className="flex justify-between items-center">
                  <div>
                    <span className="font-medium">{inst.name}</span>{' '}
                    <span className="text-text-muted">
                      {inst.itemType} lv.{inst.level}
                    </span>
                    {!inst.isVisible && (
                      <span className="text-warning ml-1 text-[10px]">(hidden)</span>
                    )}
                    {inst.isImmovable && (
                      <span className="text-danger ml-1 text-[10px]">(immovable)</span>
                    )}
                  </div>
                  <div className="flex gap-0.5">
                    <Button
                      variant="ghost"
                      size="sm"
                      className="!px-0.5 !py-0"
                      onClick={() => startEdit(inst)}
                      aria-label={`Edit ${inst.name}`}
                    >
                      ✏️
                    </Button>
                    <Button
                      variant={confirmDeleteId === inst.id ? 'secondary' : 'ghost'}
                      size="sm"
                      className="!px-0.5 !py-0"
                      onClick={() => {
                        if (confirmDeleteId === inst.id) {
                          handleDelete(inst.id)
                        } else {
                          setConfirmDeleteId(inst.id)
                        }
                      }}
                      aria-label={
                        confirmDeleteId === inst.id
                          ? `Confirm delete ${inst.name}`
                          : `Delete ${inst.name}`
                      }
                    >
                      {confirmDeleteId === inst.id ? '❓' : '🗑'}
                    </Button>
                  </div>
                </div>
              )}

              {confirmDeleteId === inst.id && editingId !== inst.id && (
                <div className="mt-1 text-[10px] text-text">
                  Confirm delete?{' '}
                  <Button
                    variant="danger"
                    size="sm"
                    className="!px-1 !py-0 !text-[10px]"
                    onClick={() => handleDelete(inst.id)}
                  >
                    Yes
                  </Button>{' '}
                  <Button
                    variant="ghost"
                    size="sm"
                    className="!px-1 !py-0 !text-[10px]"
                    onClick={() => setConfirmDeleteId(null)}
                  >
                    No
                  </Button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      {/* Spawn modal */}
      <Modal
        isOpen={showSpawn}
        onClose={() => setShowSpawn(false)}
        title="Spawn Item Instance"
      >
        <div className="space-y-3">
          {createMutation.isError && (
            <div className="p-2 bg-danger/10 text-danger rounded text-xs">
              {(createMutation.error as Error)?.message ||
                'Failed to spawn item instance'}
            </div>
          )}

          <div>
            <label className="text-text-muted text-xs block mb-1">
              Equipment Template *
            </label>
            {templatesQuery.isLoading ? (
              <div className="text-text-muted text-xs">Loading templates...</div>
            ) : (
              <select
                value={spawnForm.template_id}
                onChange={(e) => applyTemplateDefaults(e.target.value)}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              >
                <option value="">-- Select template --</option>
                {(templatesQuery.data ?? []).map((t) => (
                  <option key={t.equipment_template_id} value={t.equipment_template_id}>
                    {t.name} ({t.slot}, lv.{t.level})
                  </option>
                ))}
              </select>
            )}
          </div>

          {selectedTemplate && (
            <div className="p-2 bg-surface-muted border border-border rounded text-xs text-text-muted space-y-0.5">
              <div>Slot: {selectedTemplate.slot}</div>
              <div>Level: {selectedTemplate.level} • Weight: {selectedTemplate.weight}</div>
              <div>Type: {selectedTemplate.item_type}</div>
            </div>
          )}

          <div>
            <label className="text-text-muted text-xs block mb-1">Name</label>
            <input
              type="text"
              value={spawnForm.name}
              onChange={(e) =>
                setSpawnForm((f) => ({ ...f, name: e.target.value }))
              }
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Description</label>
            <input
              type="text"
              value={spawnForm.description}
              onChange={(e) =>
                setSpawnForm((f) => ({ ...f, description: e.target.value }))
              }
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div className="flex gap-2">
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Slot</label>
              <input
                type="text"
                value={spawnForm.slot}
                onChange={(e) =>
                  setSpawnForm((f) => ({ ...f, slot: e.target.value }))
                }
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Level</label>
              <input
                type="number"
                value={spawnForm.level}
                onChange={(e) =>
                  setSpawnForm((f) => ({
                    ...f,
                    level: parseInt(e.target.value, 10) || 0,
                  }))
                }
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          <div className="flex gap-2">
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Weight</label>
              <input
                type="number"
                value={spawnForm.weight}
                onChange={(e) =>
                  setSpawnForm((f) => ({
                    ...f,
                    weight: parseInt(e.target.value, 10) || 0,
                  }))
                }
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">Color</label>
              <input
                type="text"
                value={spawnForm.color}
                onChange={(e) =>
                  setSpawnForm((f) => ({ ...f, color: e.target.value }))
                }
                placeholder="#8b5cf6"
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">Room ID</label>
            <input
              type="number"
              value={spawnForm.room_id}
              onChange={(e) =>
                setSpawnForm((f) => ({
                  ...f,
                  room_id: parseInt(e.target.value, 10) || 0,
                }))
              }
              className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
            />
          </div>

          <div className="flex gap-2 pt-2">
            <Button
              variant="primary"
              onClick={handleSpawn}
              disabled={createMutation.isPending || !spawnForm.template_id}
            >
              {createMutation.isPending ? 'Spawning...' : 'Spawn Instance'}
            </Button>
            <Button variant="secondary" onClick={() => setShowSpawn(false)}>
              Cancel
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  )
}
