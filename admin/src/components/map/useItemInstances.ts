import { useState, useCallback } from 'react'
import { apiGet, apiPost, apiPut, apiDelete } from '../../utils/apiFetch'
import { showToast } from '../Toast'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { ItemInstanceView, EquipmentTemplate, SpawnFormData, EditFormData } from './ItemInstanceManager'

export function useItemInstances(roomId: number) {
  const queryClient = useQueryClient()

  const instancesQuery = useQuery({
    queryKey: ['item-instances', roomId],
    queryFn: async () =>
      apiGet<ItemInstanceView[]>(`${window.location.origin}/api/item-instances?roomId=${roomId}`),
  })

  const templatesQuery = useQuery({
    queryKey: ['equipment-templates'],
    queryFn: async () =>
      apiGet<EquipmentTemplate[]>(`${window.location.origin}/api/equipment-templates`),
  })

  const createMutation = useMutation({
    mutationFn: async (input: Record<string, unknown>) =>
      apiPost<ItemInstanceView>(`${window.location.origin}/api/item-instances`, input),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances'] }) },
  })

  const updateMutation = useMutation({
    mutationFn: async (args: { id: number; update: Record<string, unknown> }) =>
      apiPut<ItemInstanceView>(`${window.location.origin}/api/item-instances/${args.id}`, args.update),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['item-instances'] }) },
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: number) =>
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
      setEditingId(null); setEditForm({})
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
    setSpawnForm(emptySpawnForm()); setShowSpawn(true)
  }, [emptySpawnForm])

  return {
    instancesQuery, templatesQuery, createMutation, updateMutation, deleteMutation,
    showSpawn, setShowSpawn, editingId, setEditingId, confirmDeleteId, setConfirmDeleteId,
    spawnForm, setSpawnForm, editForm, setEditForm, selectedTemplate,
    applyTemplateDefaults, handleSpawn, handleUpdate, handleDelete, handleOpenSpawn,
  }
}