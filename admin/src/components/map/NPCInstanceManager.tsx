import { useState, useCallback } from 'react'
import { Button } from '../Button'
import { Modal } from '../Modal'
import { apiGet, apiPost, apiPut, apiDelete } from '../../utils/apiFetch'
import { logError } from '../../utils/log'
import { SearchableSelect } from '../SearchableSelect'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'

// ─── Types ──────────────────────────────────────────────────────────────────

type NPCInstanceView = Readonly<{
  id: number
  name: string
  npc_template_id: string
  instance_number: number
  room_id: number
  starting_room_id: number
  level: number
  race: string
  hitpoints: number
  max_hitpoints: number
  stamina: number
  max_stamina: number
  mana: number
  max_mana: number
  isNPC: boolean
  is_instance: boolean
}>

type NPCTemplate = Readonly<{
  id: string
  name: string
  race: string
  level: number
  respawn_rooms: string[]
  respawn_cooldown: number
}>

type SpawnFormData = {
  template_id: string
  level: number
  hitpoints: number
  room_id: number
  respawn_cooldown: number
  respawn_rooms: string
}

type EditFormData = {
  level: number
  hitpoints: number
  room_id: number
  starting_room_id: number
}

type NPCInstanceManagerProps = Readonly<{
  roomId: number
}>

// ─── Helpers ────────────────────────────────────────────────────────────────

/** Parse a comma-separated string of room IDs into an array of trimmed strings. */
function parseRoomIds(raw: string): string[] {
  return raw
    .split(',')
    .map((s) => s.trim())
    .filter((s) => s.length > 0)
}

// ─── Component ──────────────────────────────────────────────────────────────

export function NPCInstanceManager({ roomId }: NPCInstanceManagerProps) {
  const queryClient = useQueryClient()

  // Fetch NPC instances for this room
  const instancesQuery = useQuery({
    queryKey: ['npc-instances', roomId],
    queryFn: async (): Promise<NPCInstanceView[]> => {
      return apiGet<NPCInstanceView[]>(
        `${window.location.origin}/api/npc-instances?roomId=${roomId}`
      )
    },
  })

  // Fetch NPC templates for spawn form
  const templatesQuery = useQuery({
    queryKey: ['npc-templates'],
    queryFn: async (): Promise<NPCTemplate[]> => {
      return apiGet<NPCTemplate[]>(
        `${window.location.origin}/api/npc-templates`
      )
    },
  })

  // Create mutation
  const createMutation = useMutation({
    mutationFn: async (input: Record<string, unknown>): Promise<NPCInstanceView> => {
      return apiPost<NPCInstanceView>(
        `${window.location.origin}/api/npc-instances`,
        input
      )
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-instances'] })
    },
  })

  // Update mutation
  const updateMutation = useMutation({
    mutationFn: async (args: {
      id: number
      update: Record<string, unknown>
    }): Promise<NPCInstanceView> => {
      return apiPut<NPCInstanceView>(
        `${window.location.origin}/api/npc-instances/${args.id}`,
        args.update
      )
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-instances'] })
    },
  })

  // Delete mutation
  const deleteMutation = useMutation({
    mutationFn: async (id: number): Promise<void> => {
      await apiDelete(`${window.location.origin}/api/npc-instances/${id}`)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['npc-instances'] })
    },
  })

  const defaultSpawnForm = (): SpawnFormData => ({
    template_id: '',
    level: 0,
    hitpoints: 0,
    room_id: roomId,
    respawn_cooldown: 0,
    respawn_rooms: '',
  })

  const defaultEditForm = (): EditFormData => ({
    level: 0,
    hitpoints: 0,
    room_id: roomId,
    starting_room_id: roomId,
  })

  const [showSpawn, setShowSpawn] = useState(false)
  const [editingId, setEditingId] = useState<number | null>(null)
  const [confirmDeleteId, setConfirmDeleteId] = useState<number | null>(null)
  const [spawnForm, setSpawnForm] = useState<SpawnFormData>(defaultSpawnForm())
  const [editForm, setEditForm] = useState<EditFormData>(defaultEditForm())

  const handleSpawn = useCallback(async () => {
    if (!spawnForm.template_id) return
    try {
      const payload: Record<string, unknown> = {
        template_id: spawnForm.template_id,
        room_id: spawnForm.room_id,
      }
      if (spawnForm.level > 0) payload.level = spawnForm.level
      if (spawnForm.hitpoints > 0) payload.hitpoints = spawnForm.hitpoints
      if (spawnForm.respawn_cooldown > 0) {
        payload.respawn_cooldown = spawnForm.respawn_cooldown
      }
      const parsedRooms = parseRoomIds(spawnForm.respawn_rooms)
      if (parsedRooms.length > 0) {
        payload.respawn_rooms = parsedRooms
      }
      await createMutation.mutateAsync(payload)
      setShowSpawn(false)
      setSpawnForm(defaultSpawnForm())
    } catch (err) {
      logError('Spawn NPC instance:', err)
    }
  }, [spawnForm, createMutation, roomId])

  const startEdit = useCallback(
    (inst: NPCInstanceView) => {
      setEditingId(inst.id)
      setConfirmDeleteId(null)
      setEditForm({
        level: inst.level,
        hitpoints: inst.hitpoints,
        room_id: inst.room_id,
        starting_room_id: inst.room_id,
      })
    },
    []
  )

  const handleUpdate = useCallback(async () => {
    if (editingId === null) return
    try {
      const update: Record<string, unknown> = {}
      if (editForm.level > 0) update.level = editForm.level
      if (editForm.hitpoints > 0) update.hitpoints = editForm.hitpoints
      update.room_id = editForm.room_id
      update.starting_room_id = editForm.starting_room_id
      await updateMutation.mutateAsync({ id: editingId, update })
      setEditingId(null)
      setEditForm(defaultEditForm())
    } catch (err) {
      logError('Update NPC instance:', err)
    }
  }, [editingId, editForm, updateMutation, roomId])

  const handleDelete = useCallback(
    async (id: number) => {
      try {
        await deleteMutation.mutateAsync(id)
        setConfirmDeleteId(null)
        if (editingId === id) {
          setEditingId(null)
          setEditForm(defaultEditForm())
        }
      } catch (err) {
        logError('Delete NPC instance:', err)
      }
    },
    [deleteMutation, editingId, roomId]
  )

  const handleOpenSpawn = useCallback(() => {
    setSpawnForm(defaultSpawnForm())
    setShowSpawn(true)
  }, [roomId])

  const templates = templatesQuery.data ?? []
  const instances = instancesQuery.data ?? []

  /** Find the currently-selected template to display its respawn info. */
  const selectedTemplate = spawnForm.template_id
    ? templates.find((t) => t.id === spawnForm.template_id) ?? null
    : null

  if (instancesQuery.isLoading) {
    return (
      <div className="mb-3">
        <strong className="text-warning text-xs">NPCs:</strong>
        <div className="text-text-muted text-[10px] mt-1">Loading...</div>
      </div>
    )
  }

  if (instancesQuery.error) {
    return (
      <div className="mb-3">
        <strong className="text-warning text-xs">NPCs:</strong>
        <div className="text-danger text-[10px] mt-1">
          Error loading NPCs
        </div>
      </div>
    )
  }

  return (
    <div className="mb-3">
      <div className="flex items-center justify-between mb-1">
        <strong className="text-warning text-xs">NPCs:</strong>
        <Button variant="primary" size="sm" className="!px-1.5 !py-0 !text-[10px]" onClick={handleOpenSpawn}>
          + Spawn
        </Button>
      </div>

      {instances.length === 0 ? (
        <div className="text-text-muted text-[10px]">No NPCs in this room.</div>
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
                      <label className="text-text-muted block text-[10px] mb-0.5">Level</label>
                      <input
                        type="number"
                        value={editForm.level}
                        onChange={(e) =>
                          setEditForm((f) => ({
                            ...f,
                            level: parseInt(e.target.value, 10) || 0,
                          }))
                        }
                        className="w-full p-0.5 bg-surface border border-border rounded text-text text-[10px]"
                      />
                    </div>
                    <div className="flex-1">
                      <label className="text-text-muted block text-[10px] mb-0.5">HP</label>
                      <input
                        type="number"
                        value={editForm.hitpoints}
                        onChange={(e) =>
                          setEditForm((f) => ({
                            ...f,
                            hitpoints: parseInt(e.target.value, 10) || 0,
                          }))
                        }
                        className="w-full p-0.5 bg-surface border border-border rounded text-text text-[10px]"
                      />
                    </div>
                  </div>
                  <div>
                    <label className="text-text-muted block text-[10px] mb-0.5">Room ID</label>
                    <input
                      type="number"
                      value={editForm.room_id}
                      onChange={(e) =>
                        setEditForm((f) => ({
                          ...f,
                          room_id: parseInt(e.target.value, 10) || 0,
                        }))
                      }
                      className="w-full p-0.5 bg-surface border border-border rounded text-text text-[10px]"
                    />
                  </div>
                  <div>
                    <label className="text-text-muted block text-[10px] mb-0.5">Starting Room</label>
                    <input
                      type="number"
                      value={editForm.starting_room_id}
                      onChange={(e) =>
                        setEditForm((f) => ({
                          ...f,
                          starting_room_id: parseInt(e.target.value, 10) || 0,
                        }))
                      }
                      className="w-full p-0.5 bg-surface border border-border rounded text-text text-[10px]"
                    />
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
                      {inst.race} lv.{inst.level} HP:{inst.hitpoints}/{inst.max_hitpoints}
                    </span>
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
        title="Spawn NPC Instance"
      >
        <div className="space-y-3">
          {createMutation.isError && (
            <div className="p-2 bg-danger/10 text-danger rounded text-xs">
              {(createMutation.error as Error)?.message ||
                'Failed to spawn NPC instance'}
            </div>
          )}

          <div>
            <label className="text-text-muted text-xs block mb-1">
              NPC Template *
            </label>
            {templatesQuery.isLoading ? (
              <div className="text-text-muted text-xs">Loading templates...</div>
            ) : (
              <SearchableSelect
                options={templates}
                value={spawnForm.template_id}
                onChange={(id) =>
                  setSpawnForm((f) => ({ ...f, template_id: id }))
                }
                placeholder="Search by name or ID..."
                disabled={createMutation.isPending}
              />
            )}
          </div>

          {/* Read-only template respawn info */}
          {selectedTemplate && (
            <div className="p-2 bg-surface-muted border border-border rounded text-xs text-text-muted space-y-0.5">
              <div>
                Respawn cooldown: {selectedTemplate.respawn_cooldown}s
              </div>
              <div>
                Respawn rooms: {selectedTemplate.respawn_rooms.length > 0
                  ? selectedTemplate.respawn_rooms.join(', ')
                  : 'none'}
              </div>
            </div>
          )}

          <div className="flex gap-2">
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">
                Level <span className="text-text-muted">(0 = default)</span>
              </label>
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
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">
                HP <span className="text-text-muted">(0 = default)</span>
              </label>
              <input
                type="number"
                value={spawnForm.hitpoints}
                onChange={(e) =>
                  setSpawnForm((f) => ({
                    ...f,
                    hitpoints: parseInt(e.target.value, 10) || 0,
                  }))
                }
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
          </div>

          <div>
            <label className="text-text-muted text-xs block mb-1">
              Room ID
            </label>
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

          {/* Respawn override fields */}
          <div className="flex gap-2">
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">
                Respawn Cooldown <span className="text-text-muted">(0 = template default)</span>
              </label>
              <input
                type="number"
                value={spawnForm.respawn_cooldown}
                onChange={(e) =>
                  setSpawnForm((f) => ({
                    ...f,
                    respawn_cooldown: parseInt(e.target.value, 10) || 0,
                  }))
                }
                min={0}
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
            <div className="flex-1">
              <label className="text-text-muted text-xs block mb-1">
                Respawn Rooms <span className="text-text-muted">(empty = template default)</span>
              </label>
              <input
                type="text"
                value={spawnForm.respawn_rooms}
                onChange={(e) =>
                  setSpawnForm((f) => ({
                    ...f,
                    respawn_rooms: e.target.value,
                  }))
                }
                placeholder="e.g. 101, 102, 103"
                className="w-full p-2 bg-surface border border-border rounded text-text text-sm"
              />
            </div>
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