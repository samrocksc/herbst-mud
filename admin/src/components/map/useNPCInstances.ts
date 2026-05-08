import { useState, useCallback } from 'react'
import { apiGet, apiPost, apiPut, apiDelete } from '../../utils/apiFetch'
import { showToast } from '../Toast'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import type { NPCInstanceView, NPCTemplate, SpawnFormData, EditFormData } from './NPCInstanceManager'

function parseRoomIds(raw: string): string[] {
  return raw.split(',').map((s) => s.trim()).filter((s) => s.length > 0)
}

export function useNPCInstances(roomId: number) {
  const queryClient = useQueryClient()
  const [showSpawn, setShowSpawn] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null)
  const [spawnForm, setSpawnForm] = useState<SpawnFormData>({ template_id: '', level: 0, hitpoints: 0, room_id: roomId, respawn_cooldown: 0, respawn_rooms: '' })
  const [editForm, setEditForm] = useState<EditFormData>({ level: 0, hitpoints: 0, room_id: roomId, starting_room_id: roomId })

  const instancesQuery = useQuery({
    queryKey: ['npc-instances', roomId],
    queryFn: async (): Promise<NPCInstanceView[]> =>
      apiGet<NPCInstanceView[]>(`${window.location.origin}/api/npc-instances?roomId=${roomId}`),
  })

  const templatesQuery = useQuery({
    queryKey: ['npc-templates'],
    queryFn: async (): Promise<NPCTemplate[]> =>
      apiGet<NPCTemplate[]>(`${window.location.origin}/api/npc-templates`),
  })

  const createMutation = useMutation({
    mutationFn: async (input: Record<string, unknown>): Promise<NPCInstanceView> =>
      apiPost<NPCInstanceView>(`${window.location.origin}/api/npc-instances`, input),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['npc-instances'] }) },
  })

  const updateMutation = useMutation({
    mutationFn: async (args: { id: number; update: Record<string, unknown> }): Promise<NPCInstanceView> =>
      apiPut<NPCInstanceView>(`${window.location.origin}/api/npc-instances/${args.id}`, args.update),
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['npc-instances'] }) },
  })

  const deleteMutation = useMutation({
    mutationFn: async (id: number): Promise<void> => { await apiDelete(`${window.location.origin}/api/npc-instances/${id}`) },
    onSuccess: () => { queryClient.invalidateQueries({ queryKey: ['npc-instances'] }); showToast('NPC instance deleted', 'success') },
  })

  const handleSpawn = useCallback(async () => {
    if (!spawnForm.template_id) return
    try {
      const payload: Record<string, unknown> = { template_id: spawnForm.template_id, room_id: spawnForm.room_id }
      if (spawnForm.level > 0) payload.level = spawnForm.level
      if (spawnForm.hitpoints > 0) payload.hitpoints = spawnForm.hitpoints
      if (spawnForm.respawn_cooldown > 0) payload.respawn_cooldown = spawnForm.respawn_cooldown
      const parsedRooms = parseRoomIds(spawnForm.respawn_rooms)
      if (parsedRooms.length > 0) payload.respawn_rooms = parsedRooms
      await createMutation.mutateAsync(payload)
      showToast('NPC spawned successfully', 'success')
      setShowSpawn(false)
      setSpawnForm({ template_id: '', level: 0, hitpoints: 0, room_id: roomId, respawn_cooldown: 0, respawn_rooms: '' })
    } catch (err) { showToast(`Spawn failed: ${(err as Error)?.message ?? 'Unknown error'}`, 'error') }
  }, [spawnForm, createMutation, roomId])

  const handleUpdate = useCallback(async () => {
    if (editingId === null) return
    try {
      const update: Record<string, unknown> = { room_id: editForm.room_id, starting_room_id: editForm.starting_room_id }
      if (editForm.level > 0) update.level = editForm.level
      if (editForm.hitpoints > 0) update.hitpoints = editForm.hitpoints
      await updateMutation.mutateAsync({ id: editingId, update })
      showToast('NPC updated', 'success')
      setEditingId(null)
      setEditForm({ level: 0, hitpoints: 0, room_id: roomId, starting_room_id: roomId })
    } catch (err) { showToast(`Update failed: ${(err as Error)?.message ?? 'Unknown error'}`, 'error') }
  }, [editingId, editForm, updateMutation, roomId])

  const handleDelete = useCallback(async (id: number) => {
    try { await deleteMutation.mutateAsync(id); setConfirmDeleteId(null) }
    catch (err) { showToast(`Delete failed: ${(err as Error)?.message ?? 'Unknown error'}`, 'error') }
  }, [deleteMutation])

  const startEdit = useCallback((inst: NPCInstanceView) => {
    setEditingId(inst.id); setConfirmDeleteId(null)
    setEditForm({ level: inst.level, hitpoints: inst.hitpoints, room_id: inst.room_id, starting_room_id: inst.room_id })
  }, [])

  return {
    instancesQuery, templatesQuery, createMutation, updateMutation, deleteMutation,
    showSpawn, setShowSpawn, editingId, setEditingId, confirmDeleteId, setConfirmDeleteId,
    spawnForm, setSpawnForm, editForm, setEditForm,
    handleSpawn, handleUpdate, handleDelete, startEdit,
  }
}